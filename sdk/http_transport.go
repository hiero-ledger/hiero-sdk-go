package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

// The generic HTTP pipe: one exchange in, one exchange out, no policy of any kind.
//
// Nothing here mentions the mirror node on purpose. Its only consumer today is mirror REST
// (mirror_http_client.go adds base-URL resolution and the retry policy), but the contract is
// plain HTTP — an opaque body with the content type as a header — so it can carry protobuf as
// readily as JSON. That matters because the likely next consumer is gRPC-Web, which is
// protobuf over HTTP POST. Naming this after its first consumer would have made that a
// rename of exported API rather than of nothing.
//
// This mirrors the split in the V3 spec: a generic `http.HttpClient` in a base namespace,
// with `mirrornode.http.MirrorNodeHttpClient` as the mirror-specific adapter above it.
//
// What this pipe still cannot do is stream, so it cannot carry a topic subscription or
// gRPC-Web server streaming. That is a separate contract, not a flag on this one: full
// buffering is what makes a retry replayable, and retry-after-partial-read has no meaning.
//
// Prototype for sdk-collaboration-hub#283. Deliberately unexported: the shipping version
// exports the transport so applications can inject a proxy / mTLS / tracing implementation,
// but a prototype must not commit v2 to a public shape. See
// .claude/docs/go-mirror-http-prototype.md.

const (
	httpMethodGet  = http.MethodGet
	httpMethodPost = http.MethodPost

	// Bodies are read into memory, so cap what a mirror node can make us allocate.
	httpMaxResponseBytes = 32 << 20

	// Redirects are followed but bounded, and credentials are not carried across origins.
	httpMaxRedirects = 5

	// Defaults matching http.DefaultTransport, so attaching this layer to an existing call
	// site cannot tighten a bound the caller was already living with. The knob is the new
	// capability here; the default deliberately is not.
	httpDefaultDialTimeout      = 30 * time.Second
	httpDefaultHandshakeTimeout = 10 * time.Second

	httpHeaderAuthorization = "Authorization"
	httpHeaderRetryAfter    = "Retry-After"
	httpHeaderUserAgent     = "User-Agent"
	httpHeaderContentType   = "Content-Type"
)

var errHttpTransportClosed = errors.New("mirror node HTTP transport is closed")

// httpRequest is one mirror REST exchange. The URL is absolute and already resolved;
// see mirror_http_path.go for how it is built from a base URL and a path.
type httpRequest struct {
	method      string
	url         string
	body        []byte
	contentType string
	headers     map[string]string
}

// httpResponse is a completed exchange. A non-2xx status is a completed exchange.
type httpResponse struct {
	statusCode int
	body       []byte
	headers    http.Header
}

// httpTransport performs single HTTP exchanges and nothing else. Status-code
// classification, retry and backoff live in mirrorHttpClient above it, never here — so a
// binding can swap in a native stack without reimplementing policy.
type httpTransport interface {
	// roundTrip performs one exchange. A non-2xx status is returned in the response, never as
	// an error; only a failure to obtain any response returns one, always a *httpError.
	roundTrip(ctx context.Context, req httpRequest) (httpResponse, error)

	// close releases resources, waiting at most closeTimeout for in-flight exchanges and then
	// aborting them. Idempotent, and safe to race with roundTrip.
	close(closeTimeout time.Duration) error
}

// httpErrorKind says whether retrying could plausibly help.
type httpErrorKind int

const (
	// httpTransient — the exchange may succeed if repeated.
	httpTransient httpErrorKind = iota
	// httpTerminal — repeating cannot fix it, so it must not consume the attempt budget.
	httpTerminal
	// httpClosed — the transport was closed; permanent for this instance.
	httpClosed
)

func (k httpErrorKind) String() string {
	switch k {
	case httpTransient:
		return "transient"
	case httpTerminal:
		return "terminal"
	case httpClosed:
		return "closed"
	default:
		return "unknown"
	}
}

// httpError carries the retryability decision alongside the cause, so the layer above
// does not have to re-derive it by inspecting error strings.
type httpError struct {
	op   string
	kind httpErrorKind
	err  error
}

func (e *httpError) Error() string {
	return fmt.Sprintf("mirror node HTTP %s: %s failure: %v", e.op, e.kind, e.err)
}

func (e *httpError) Unwrap() error {
	return e.err
}

// httpErrorKindOf reports the classification of err, and whether it carried one.
func httpErrorKindOf(err error) (httpErrorKind, bool) {
	var httpErr *httpError
	if errors.As(err, &httpErr) {
		return httpErr.kind, true
	}
	return httpTransient, false
}

// classifyTransportError decides whether repeating the exchange could help. A name that does
// not resolve and a certificate that does not validate are permanent: retrying cannot fix
// either, and re-running certificate validation to burn the budget is a mild security smell.
func classifyTransportError(err error) httpErrorKind {
	if errors.Is(err, context.Canceled) {
		return httpTerminal
	}

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) && dnsErr.IsNotFound {
		return httpTerminal
	}

	var certErr *tls.CertificateVerificationError
	if errors.As(err, &certErr) {
		return httpTerminal
	}

	var unknownAuthority x509.UnknownAuthorityError
	if errors.As(err, &unknownAuthority) {
		return httpTerminal
	}

	var hostnameErr x509.HostnameError
	if errors.As(err, &hostnameErr) {
		return httpTerminal
	}

	var invalidCert x509.CertificateInvalidError
	if errors.As(err, &invalidCert) {
		return httpTerminal
	}

	// Connection refused, reset, truncated reads and per-attempt timeouts all get another go.
	return httpTransient
}

// defaultHttpTransport owns a private http.Transport rather than sharing
// http.DefaultTransport, so pool limits, TLS and proxy settings are ours and not whatever
// else in the process touched the default.
type defaultHttpTransport struct {
	httpClient *http.Client
	userAgent  string

	// abortCtx is cancelled by close once the grace period expires, which is what turns
	// "wait, then abort" into an actual abort rather than a hope.
	abortCtx context.Context
	abort    context.CancelFunc

	mu       sync.Mutex
	closed   bool
	inFlight sync.WaitGroup
}

// newDefaultHttpTransport builds a transport over a private http.Transport.
// A connectTimeout of 0 keeps net/http's own dial and handshake bounds.
func newDefaultHttpTransport(connectTimeout time.Duration) *defaultHttpTransport {
	dialTimeout, handshakeTimeout := httpDefaultDialTimeout, httpDefaultHandshakeTimeout
	if connectTimeout > 0 {
		dialTimeout, handshakeTimeout = connectTimeout, connectTimeout
	}

	dialer := &net.Dialer{
		Timeout:   dialTimeout,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   8,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   handshakeTimeout,
		ExpectContinueTimeout: time.Second,
	}

	return newHttpTransportOver(transport)
}

// newHttpTransportOver wraps a caller-supplied RoundTripper instead of building one.
// This is how a proxy, custom TLS, tracing or metrics layer is added without reimplementing
// the policy above it — and how a test substitutes a fake without touching any global.
// The supplied RoundTripper owns dialing, so it also owns the connect timeout.
func newHttpTransportOver(base http.RoundTripper) *defaultHttpTransport {
	ctx, cancel := context.WithCancel(context.Background())

	return &defaultHttpTransport{
		httpClient: &http.Client{
			Transport:     base,
			CheckRedirect: checkHttpRedirect,
			// No Client.Timeout: the per-attempt deadline travels on the context so a caller
			// can cancel one mirror read, which Client.requestTimeout cannot express.
		},
		userAgent: getUserAgent(),
		abortCtx:  ctx,
		abort:     cancel,
	}
}

// checkHttpRedirect bounds redirects and drops credentials when the origin changes.
func checkHttpRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= httpMaxRedirects {
		return fmt.Errorf("stopped after %d redirects", httpMaxRedirects)
	}

	previous := via[len(via)-1].URL
	if req.URL.Scheme != previous.Scheme || req.URL.Host != previous.Host {
		req.Header.Del(httpHeaderAuthorization)
	}

	return nil
}

func (t *defaultHttpTransport) roundTrip(ctx context.Context, req httpRequest) (httpResponse, error) {
	if err := t.beginRequest(); err != nil {
		return httpResponse{}, err
	}
	defer t.inFlight.Done()

	// Abort on close as well as on the caller's own deadline.
	ctx, cancel := httpJoinContexts(ctx, t.abortCtx)
	defer cancel()

	httpReq, err := t.buildRequest(ctx, req)
	if err != nil {
		return httpResponse{}, &httpError{op: "build", kind: httpTerminal, err: err}
	}

	resp, err := t.httpClient.Do(httpReq)
	if err != nil {
		if t.isClosed() && errors.Is(err, context.Canceled) {
			return httpResponse{}, &httpError{op: "send", kind: httpClosed, err: errHttpTransportClosed}
		}
		return httpResponse{}, &httpError{op: "send", kind: classifyTransportError(err), err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, httpMaxResponseBytes))
	if err != nil {
		// A truncated read is a failure to obtain a response, so it is transport-shaped.
		return httpResponse{}, &httpError{op: "read", kind: classifyTransportError(err), err: err}
	}

	return httpResponse{statusCode: resp.StatusCode, body: body, headers: resp.Header}, nil
}

func (t *defaultHttpTransport) buildRequest(ctx context.Context, req httpRequest) (*http.Request, error) {
	var bodyReader io.Reader
	if len(req.body) > 0 {
		// A fresh reader per attempt is what makes the body replayable on retry.
		bodyReader = bytes.NewReader(req.body)
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.method, req.url, bodyReader)
	if err != nil {
		return nil, err
	}

	for name, value := range req.headers {
		httpReq.Header.Set(name, value)
	}
	if req.contentType != "" {
		httpReq.Header.Set(httpHeaderContentType, req.contentType)
	}
	// The SDK token is always present; a caller-supplied agent is appended, not replaced.
	if caller := httpReq.Header.Get(httpHeaderUserAgent); caller != "" {
		httpReq.Header.Set(httpHeaderUserAgent, t.userAgent+" "+caller)
	} else {
		httpReq.Header.Set(httpHeaderUserAgent, t.userAgent)
	}

	return httpReq, nil
}

// beginRequest registers an in-flight exchange, rejecting it if the transport is closed. The
// registration happens under the same lock as the closed check so close cannot miss it.
func (t *defaultHttpTransport) beginRequest() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return &httpError{op: "send", kind: httpClosed, err: errHttpTransportClosed}
	}
	t.inFlight.Add(1)

	return nil
}

func (t *defaultHttpTransport) isClosed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.closed
}

// close stops accepting work, waits up to closeTimeout for in-flight exchanges, then aborts
// whatever is left. It never reports a failure: a bounded shutdown is the guarantee, not an
// outcome.
func (t *defaultHttpTransport) close(closeTimeout time.Duration) error {
	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return nil
	}
	t.closed = true
	t.mu.Unlock()

	drained := make(chan struct{})
	go func() {
		t.inFlight.Wait()
		close(drained)
	}()

	select {
	case <-drained:
	case <-time.After(closeTimeout):
	}

	t.abort()
	t.httpClient.CloseIdleConnections()

	return nil
}

// httpJoinContexts returns a context cancelled when either input is. context.WithCancel
// only takes one parent, and an exchange has two reasons to stop: the caller and close.
func httpJoinContexts(caller, abort context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(caller)

	stop := make(chan struct{})
	go func() {
		select {
		case <-abort.Done():
			cancel()
		case <-stop:
		}
	}()

	return ctx, func() {
		close(stop)
		cancel()
	}
}
