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
	"maps"
	"math"
	"net"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"sync"
	"syscall"
	"time"
)

// HttpMethod is the method of an HttpRequest.
type HttpMethod int32

const (
	HttpMethodGet HttpMethod = iota
	HttpMethodHead
	HttpMethodPost
	HttpMethodPut
	HttpMethodPatch
	HttpMethodDelete
	HttpMethodOptions
	HttpMethodTrace
	HttpMethodConnect
)

// String returns the HTTP method name, such as "GET".
func (m HttpMethod) String() string {
	switch m {
	case HttpMethodGet:
		return http.MethodGet
	case HttpMethodHead:
		return http.MethodHead
	case HttpMethodPost:
		return http.MethodPost
	case HttpMethodPut:
		return http.MethodPut
	case HttpMethodPatch:
		return http.MethodPatch
	case HttpMethodDelete:
		return http.MethodDelete
	case HttpMethodOptions:
		return http.MethodOptions
	case HttpMethodTrace:
		return http.MethodTrace
	case HttpMethodConnect:
		return http.MethodConnect
	default:
		return unknownString
	}
}

const (
	// Same as http.DefaultTransport.
	httpDefaultDialTimeout      = 30 * time.Second
	httpDefaultHandshakeTimeout = 10 * time.Second

	httpOpBuild = "build"
	httpOpSend  = "send"
	httpOpRead  = "read"

	httpHeaderRetryAfter  = "retry-after"
	httpHeaderXUserAgent  = "x-user-agent"
	httpHeaderContentType = "content-type"
)

// HttpRequest is a single HTTP request to an absolute URL.
type HttpRequest struct {
	method      HttpMethod
	url         string
	body        []byte
	contentType string
	headers     map[string]string
	deadline    time.Duration
}

var httpRequestURLPattern = regexp.MustCompile(`^https?://[^\s]+$`)

// NewHttpRequest creates an HttpRequest, returning an error if url is not an absolute http or https URL.
func NewHttpRequest(method HttpMethod, url string) (HttpRequest, error) {
	if !httpRequestURLPattern.MatchString(url) {
		return HttpRequest{}, fmt.Errorf("%q is not an absolute http or https URL", url)
	}

	return HttpRequest{method: method, url: url}, nil
}

// GetMethod returns the request method.
func (r HttpRequest) GetMethod() HttpMethod { return r.method }

// GetURL returns the absolute request URL.
func (r HttpRequest) GetURL() string { return r.url }

// GetBody returns a copy of the request body.
func (r HttpRequest) GetBody() []byte { return slices.Clone(r.body) }

// GetContentType returns the Content-Type of the body.
func (r HttpRequest) GetContentType() string { return r.contentType }

// GetHeaders returns a copy of the request headers, with lowercase names.
func (r HttpRequest) GetHeaders() map[string]string { return maps.Clone(r.headers) }

// GetDeadline returns the deadline for the whole exchange, body included, or 0 if none is set.
func (r HttpRequest) GetDeadline() time.Duration { return r.deadline }

// WithBody returns a copy of the request with the given body and Content-Type.
func (r HttpRequest) WithBody(contentType string, body []byte) HttpRequest {
	r.contentType = contentType
	r.body = slices.Clone(body)
	return r
}

// WithHeader returns a copy of the request with the header set. The name is lowercased.
func (r HttpRequest) WithHeader(name, value string) HttpRequest {
	r.headers = withHeader(r.headers, name, value)
	return r
}

// WithDeadline returns a copy of the request bounded by deadline, body included. It panics if deadline is negative.
func (r HttpRequest) WithDeadline(deadline time.Duration) HttpRequest {
	requireNonNegativeDuration("deadline", deadline)
	r.deadline = deadline
	return r
}

// HttpResponse is a completed HTTP exchange, whatever its status. Header names are lowercase.
type HttpResponse struct {
	statusCode int
	body       []byte
	headers    map[string][]string
}

// NewHttpResponse creates an HttpResponse, lowercasing the header names. body is not copied.
func NewHttpResponse(statusCode uint16, body []byte, headers map[string][]string) HttpResponse {
	return HttpResponse{statusCode: int(statusCode), body: body, headers: lowercaseHeaders(headers)}
}

// GetStatusCode returns the HTTP status code.
func (r HttpResponse) GetStatusCode() uint16 { return uint16(r.statusCode) }

// GetBody returns a copy of the response body.
func (r HttpResponse) GetBody() []byte { return slices.Clone(r.body) }

// GetHeaders returns a copy of the response headers, with lowercase names.
func (r HttpResponse) GetHeaders() map[string][]string {
	return lowercaseHeaders(r.headers)
}

// HttpTransport sends single HTTP exchanges. Retries and status handling are left to the caller.
type HttpTransport interface {
	// RoundTrip sends req. A non-2xx status is a response, not an error. If req's deadline expires
	// the error is HttpTransportTimeoutError; if cancellation fires it is HttpTransportCancelledError.
	RoundTrip(req HttpRequest, cancellation Cancellation) (HttpResponse, error)

	// Close waits up to closeTimeout for in-flight requests, then aborts them. It is idempotent.
	Close(closeTimeout time.Duration)
}

// httpErrorKind is the retry verdict for an HttpTransportErrorID.
type httpErrorKind int

const (
	httpTransient httpErrorKind = iota
	httpTerminal
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

// HttpTransportErrorID identifies why an HTTP exchange failed. The zero value is unknown and not retried.
type HttpTransportErrorID int

const (
	_ HttpTransportErrorID = iota
	// HttpTransportConnectionError means the connection was refused, reset or unreachable. It is retried.
	HttpTransportConnectionError
	// HttpTransportTimeoutError means the request deadline expired. It is retried.
	HttpTransportTimeoutError
	// HttpTransportUnknownHostError means the host name did not resolve.
	HttpTransportUnknownHostError
	// HttpTransportTLSError means the TLS handshake or certificate verification failed.
	HttpTransportTLSError
	// HttpTransportClientClosedError means the transport or the Client was closed.
	HttpTransportClientClosedError
	// HttpTransportCancelledError means the caller's cancellation fired.
	HttpTransportCancelledError
	// HttpTransportResponseTooLargeError means the response body exceeded the configured maximum size.
	HttpTransportResponseTooLargeError
)

// String returns the error ID name, such as "timeout-error".
func (id HttpTransportErrorID) String() string {
	switch id {
	case HttpTransportConnectionError:
		return "connection-error"
	case HttpTransportTimeoutError:
		return "timeout-error"
	case HttpTransportUnknownHostError:
		return "unknown-host-error"
	case HttpTransportTLSError:
		return "tls-error"
	case HttpTransportClientClosedError:
		return "client-closed-error"
	case HttpTransportCancelledError:
		return "cancelled-error"
	case HttpTransportResponseTooLargeError:
		return "response-too-large-error"
	default:
		return "unknown-error"
	}
}

func (id HttpTransportErrorID) kind() httpErrorKind {
	switch id {
	case HttpTransportConnectionError, HttpTransportTimeoutError:
		return httpTransient
	case HttpTransportClientClosedError:
		return httpClosed
	default:
		return httpTerminal
	}
}

// HttpTransportError is returned by an HttpTransport when no response was received.
type HttpTransportError struct {
	ID  HttpTransportErrorID
	Err error
	op  string
}

// Error() implements the Error interface
func (e *HttpTransportError) Error() string {
	if e.op == "" {
		return fmt.Sprintf("mirror node HTTP: %s: %v", e.ID, e.Err)
	}
	return fmt.Sprintf("mirror node HTTP %s: %s: %v", e.op, e.ID, e.Err)
}

func (e *HttpTransportError) Unwrap() error {
	return e.Err
}

// IsRetryable reports whether repeating the exchange could help.
func (e *HttpTransportError) IsRetryable() bool {
	return e.ID.kind() == httpTransient
}

// ErrHttpTransportRetryable can be wrapped by a custom HttpTransport to mark an error as retryable.
// Any other error that is not an HttpTransportError is not retried.
var ErrHttpTransportRetryable = errors.New("retryable HTTP transport failure")

// httpErrorKindOf returns the retry verdict for err and whether err was classified. Unclassified errors are terminal.
func httpErrorKindOf(err error) (httpErrorKind, bool) {
	var httpErr *HttpTransportError
	if errors.As(err, &httpErr) {
		return httpErr.ID.kind(), true
	}
	if errors.Is(err, ErrHttpTransportRetryable) {
		return httpTransient, true
	}
	return httpTerminal, false
}

// classifyTransportError maps a net/http error to an HttpTransportErrorID, reporting false if it is not recognised.
func classifyTransportError(err error) (HttpTransportErrorID, bool) {
	switch {
	case errors.Is(err, context.Canceled):
		return HttpTransportCancelledError, true
	case errors.Is(err, context.DeadlineExceeded):
		return HttpTransportTimeoutError, true
	}

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return HttpTransportUnknownHostError, true
	}

	if isTLSFailure(err) {
		return HttpTransportTLSError, true
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return HttpTransportTimeoutError, true
	}

	var opErr *net.OpError
	switch {
	case errors.As(err, &opErr),
		errors.Is(err, syscall.ECONNREFUSED),
		errors.Is(err, syscall.ECONNRESET),
		errors.Is(err, syscall.EHOSTUNREACH),
		errors.Is(err, io.EOF),
		errors.Is(err, io.ErrUnexpectedEOF):
		return HttpTransportConnectionError, true
	}

	return 0, false
}

func isTLSFailure(err error) bool {
	var certErr *tls.CertificateVerificationError
	var unknownAuthority x509.UnknownAuthorityError
	var hostnameErr x509.HostnameError
	var invalidCert x509.CertificateInvalidError
	var alert tls.AlertError
	var recordHeader tls.RecordHeaderError

	return errors.As(err, &certErr) ||
		errors.As(err, &unknownAuthority) ||
		errors.As(err, &hostnameErr) ||
		errors.As(err, &invalidCert) ||
		errors.As(err, &alert) ||
		errors.As(err, &recordHeader)
}

// DefaultHttpTransport is the HttpTransport the SDK builds when none is injected. It uses its own http.Transport.
type DefaultHttpTransport struct {
	httpClient       *http.Client
	userAgent        string
	maxRedirects     int
	maxResponseBytes int64
	defaultHeaders   map[string]string

	// abortCtx is cancelled by Close once the grace period ends.
	abortCtx context.Context
	abort    context.CancelFunc

	mu       sync.Mutex
	closed   bool
	inFlight sync.WaitGroup
}

// NewDefaultHttpTransport builds a transport over a private http.Transport.
func NewDefaultHttpTransport(config HttpTransportConfiguration) *DefaultHttpTransport {
	dialTimeout, handshakeTimeout := httpDefaultDialTimeout, httpDefaultHandshakeTimeout
	if connectTimeout := config.GetConnectTimeout(); connectTimeout > 0 {
		dialTimeout, handshakeTimeout = connectTimeout, connectTimeout
	}

	dialer := &net.Dialer{
		Timeout:   dialTimeout,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		// Honour HTTP_PROXY and NO_PROXY like http.DefaultTransport.
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   8,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   handshakeTimeout,
		ExpectContinueTimeout: time.Second,
	}

	return newHttpTransportOver(transport, config)
}

// newHttpTransportOver wraps base, which then owns dialing and the connect timeout.
func newHttpTransportOver(base http.RoundTripper, config HttpTransportConfiguration) *DefaultHttpTransport {
	ctx, cancel := context.WithCancel(context.Background())

	t := &DefaultHttpTransport{
		userAgent:        getUserAgent(),
		maxRedirects:     int(config.GetMaxRedirects()),
		maxResponseBytes: config.GetMaxResponseBytes(),
		defaultHeaders:   config.GetDefaultHeaders(),
		abortCtx:         ctx,
		abort:            cancel,
	}
	t.httpClient = &http.Client{
		Transport:     base,
		CheckRedirect: t.checkRedirect,
		// No Timeout: the per-attempt deadline is set on the request context.
	}

	return t
}

// checkRedirect limits redirects and, on a cross-origin redirect, drops every header except x-user-agent.
func (t *DefaultHttpTransport) checkRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= t.maxRedirects {
		return fmt.Errorf("stopped after %d redirects", t.maxRedirects)
	}

	previous := via[len(via)-1].URL
	if req.URL.Scheme != previous.Scheme || req.URL.Host != previous.Host {
		for name := range req.Header {
			if !strings.EqualFold(name, httpHeaderXUserAgent) {
				req.Header.Del(name)
			}
		}
	}

	return nil
}

// RoundTrip performs one exchange, bounded by the request's deadline and by the cancellation.
func (t *DefaultHttpTransport) RoundTrip(req HttpRequest, cancellation Cancellation) (HttpResponse, error) {
	if err := t.beginRequest(); err != nil {
		return HttpResponse{}, err
	}
	defer t.inFlight.Done()

	ctx := cancellation.GetContext()
	if req.deadline > 0 {
		var cancelDeadline context.CancelFunc
		ctx, cancelDeadline = context.WithTimeout(ctx, req.deadline)
		defer cancelDeadline()
	}
	ctx, cancel := httpJoinContexts(ctx, t.abortCtx)
	defer cancel()

	httpReq, err := t.buildRequest(ctx, req)
	if err != nil {
		return HttpResponse{}, fmt.Errorf("mirror node HTTP %s: %w", httpOpBuild, err)
	}

	resp, err := t.httpClient.Do(httpReq)
	if err != nil {
		return HttpResponse{}, t.failure(httpOpSend, err, cancellation)
	}
	defer resp.Body.Close()

	// Read one byte past the cap so an oversized body is an error, not a truncation.
	limit := t.maxResponseBytes
	if limit < math.MaxInt64 {
		limit++
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, limit))
	if err != nil {
		return HttpResponse{}, t.failure(httpOpRead, err, cancellation)
	}
	if int64(len(body)) > t.maxResponseBytes {
		return HttpResponse{}, &HttpTransportError{
			op:  httpOpRead,
			ID:  HttpTransportResponseTooLargeError,
			Err: fmt.Errorf("%w of %d bytes", errHttpResponseTooLarge, t.maxResponseBytes),
		}
	}

	return HttpResponse{statusCode: resp.StatusCode, body: body, headers: lowercaseHeaders(resp.Header)}, nil
}

func lowercaseHeaders(header map[string][]string) map[string][]string {
	lowered := make(map[string][]string, len(header))
	for name, values := range header {
		key := strings.ToLower(name)
		lowered[key] = append(lowered[key], values...)
	}

	return lowered
}

// failure wraps err in an HttpTransportError, or returns it without an ID if it is not recognised.
func (t *DefaultHttpTransport) failure(op string, err error, cancellation Cancellation) error {
	if t.isClosed() && errors.Is(err, context.Canceled) {
		return &HttpTransportError{op: op, ID: HttpTransportClientClosedError, Err: errHttpTransportClosed}
	}
	if cancellation.IsCancelled() {
		return &HttpTransportError{op: op, ID: HttpTransportCancelledError, Err: err}
	}
	if id, ok := classifyTransportError(err); ok {
		return &HttpTransportError{op: op, ID: id, Err: err}
	}

	return fmt.Errorf("mirror node HTTP %s: %w", op, err)
}

func (t *DefaultHttpTransport) buildRequest(ctx context.Context, req HttpRequest) (*http.Request, error) {
	var bodyReader io.Reader
	if len(req.body) > 0 {
		// A fresh reader per attempt is what makes the body replayable on retry.
		bodyReader = bytes.NewReader(req.body)
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.method.String(), req.url, bodyReader)
	if err != nil {
		return nil, err
	}

	// Later headers override earlier ones: defaults, request headers, Content-Type, x-user-agent.
	for name, value := range t.defaultHeaders {
		httpReq.Header.Set(name, value)
	}
	for name, value := range req.headers {
		httpReq.Header.Set(name, value)
	}
	if req.contentType != "" {
		httpReq.Header.Set(httpHeaderContentType, req.contentType)
	}
	httpReq.Header.Set(httpHeaderXUserAgent, t.userAgent)

	return httpReq, nil
}

// beginRequest registers an in-flight request, or returns an error if the transport is closed.
func (t *DefaultHttpTransport) beginRequest() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return &HttpTransportError{op: httpOpSend, ID: HttpTransportClientClosedError, Err: errHttpTransportClosed}
	}
	t.inFlight.Add(1)

	return nil
}

func (t *DefaultHttpTransport) isClosed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.closed
}

// Close waits up to closeTimeout for in-flight requests, then aborts them. It is idempotent.
func (t *DefaultHttpTransport) Close(closeTimeout time.Duration) {
	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return
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
}

// httpJoinContexts returns a context that is cancelled when either caller or abort is.
func httpJoinContexts(caller, abort context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(caller)
	stop := context.AfterFunc(abort, cancel)

	return ctx, func() {
		stop()
		cancel()
	}
}
