// SPDX-License-Identifier: Apache-2.0

package hiero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/crypto/sha3"
)

// ABI represents the ethereum abi format
type ABI struct {
	Constructor        *Method
	Methods            map[string]*Method
	MethodsBySignature map[string]*Method
	Events             map[string]*Event
	Errors             map[string]*Error
}

func (a *ABI) GetMethod(name string) *Method {
	m := a.Methods[name]
	return m
}

func (a *ABI) GetMethodBySignature(methodSignature string) *Method {
	m := a.MethodsBySignature[methodSignature]
	return m
}

func (a *ABI) addError(e *Error) {
	if len(a.Errors) == 0 {
		a.Errors = map[string]*Error{}
	}
	a.Errors[e.Name] = e
}

func (a *ABI) addEvent(e *Event) {
	if len(a.Events) == 0 {
		a.Events = map[string]*Event{}
	}
	name := overloadedName(e.Name, func(s string) bool {
		_, ok := a.Events[s]
		return ok
	})
	a.Events[name] = e
}

func (a *ABI) addMethod(m *Method) {
	if len(a.Methods) == 0 {
		a.Methods = map[string]*Method{}
	}
	if len(a.MethodsBySignature) == 0 {
		a.MethodsBySignature = map[string]*Method{}
	}
	name := overloadedName(m.Name, func(s string) bool {
		_, ok := a.Methods[s]
		return ok
	})
	a.Methods[name] = m
	a.MethodsBySignature[m.Sig()] = m
}

func overloadedName(rawName string, isAvail func(string) bool) string {
	name := rawName
	ok := isAvail(name)
	for idx := 0; ok; idx++ {
		name = fmt.Sprintf("%s%d", rawName, idx)
		ok = isAvail(name)
	}
	return name
}

// NewABI returns a parsed ABI struct
func NewABI(s string) (*ABI, error) {
	return NewABIFromReader(bytes.NewReader([]byte(s)))
}

// NewABIFromReader returns an ABI object from a reader
func NewABIFromReader(r io.Reader) (*ABI, error) {
	var abi *ABI
	dec := json.NewDecoder(r)
	if err := dec.Decode(&abi); err != nil {
		return nil, err
	}
	return abi, nil
}

// UnmarshalJSON implements json.Unmarshaler interface
// nolint
func (a *ABI) UnmarshalJSON(data []byte) error {
	var fields []struct {
		Type            string
		Name            string
		Anonymous       bool
		StateMutability string
		Inputs          []*ArgumentStr
		Outputs         []*ArgumentStr
	}

	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	for _, field := range fields {
		switch field.Type {
		case "constructor":
			if a.Constructor != nil {
				return fmt.Errorf("multiple constructor declaration")
			}
			input, err := NewTupleTypeFromArgs(field.Inputs)
			if err != nil {
				return err
			}
			a.Constructor = &Method{
				Inputs: input,
			}

		case "function", "":
			c := field.StateMutability == "view" || field.StateMutability == "pure"

			inputs, err := NewTupleTypeFromArgs(field.Inputs)
			if err != nil {
				return err
			}
			outputs, err := NewTupleTypeFromArgs(field.Outputs)
			if err != nil {
				return err
			}
			method := &Method{
				Name:    field.Name,
				Const:   c,
				Inputs:  inputs,
				Outputs: outputs,
			}
			a.addMethod(method)

		case "event":
			input, err := NewTupleTypeFromArgs(field.Inputs)
			if err != nil {
				return err
			}
			event := &Event{
				Name:      field.Name,
				Anonymous: field.Anonymous,
				Inputs:    input,
			}
			a.addEvent(event)

		case "error":
			input, err := NewTupleTypeFromArgs(field.Inputs)
			if err != nil {
				return err
			}
			errObj := &Error{
				Name:   field.Name,
				Inputs: input,
			}
			a.addError(errObj)

		case "fallback":
		case "receive":
			// do nothing

		default:
			return fmt.Errorf("unknown field type '%s'", field.Type)
		}
	}
	return nil
}

// nolint
func NewABIFromList(humanReadableAbi []string) (*ABI, error) {
	res := &ABI{}
	for _, c := range humanReadableAbi {
		if strings.HasPrefix(c, "constructor") {
			typ, err := NewType("tuple" + strings.TrimPrefix(c, "constructor"))
			if err != nil {
				return nil, err
			}
			res.Constructor = &Method{
				Inputs: typ,
			}
		} else if strings.HasPrefix(c, "function ") {
			method, err := NewMethod(c)
			if err != nil {
				return nil, err
			}
			res.addMethod(method)
		} else if strings.HasPrefix(c, "event ") {
			evnt, err := NewEvent(c)
			if err != nil {
				return nil, err
			}
			res.addEvent(evnt)
		} else if strings.HasPrefix(c, "error ") {
			errTyp, err := NewError(c)
			if err != nil {
				return nil, err
			}
			res.addError(errTyp)
		} else {
			return nil, fmt.Errorf("either event or function expected")
		}
	}
	return res, nil
}

// Method is a callable function in the contract
type Method struct {
	Name    string
	Const   bool
	Inputs  *Type
	Outputs *Type
}

// Sig returns the signature of the method
func (m *Method) Sig() string {
	return buildSignature(m.Name, m.Inputs)
}

// ID returns the id of the method
func (m *Method) ID() []byte {
	k := acquireKeccak()
	k.Write([]byte(m.Sig()))
	dst := k.Sum(nil)[:4]
	releaseKeccak(k)
	return dst
}

// Encode encodes the inputs with this function
func (m *Method) Encode(args interface{}) ([]byte, error) {
	data, err := Encode(args, m.Inputs)
	if err != nil {
		return nil, err
	}
	data = append(m.ID(), data...)
	return data, nil
}

// Decode decodes the output with this function
func (m *Method) Decode(data []byte) (map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty response")
	}
	respInterface, err := Decode(m.Outputs, data)
	if err != nil {
		return nil, err
	}
	resp := respInterface.(map[string]interface{})
	return resp, nil
}

func NewMethod(name string) (*Method, error) {
	name, inputs, outputs, err := parseMethodSignature(name)
	if err != nil {
		return nil, err
	}
	m := &Method{Name: name, Inputs: inputs, Outputs: outputs}
	return m, nil
}

var (
	funcRegexpWithReturn    = regexp.MustCompile(`(\w*)\s*\((.*)\)(.*)\s*returns\s*\((.*)\)`)
	funcRegexpWithoutReturn = regexp.MustCompile(`(\w*)\s*\((.*)\)(.*)`)
)

// Event is a triggered log mechanism
type Event struct {
	Name      string
	Anonymous bool
	Inputs    *Type
}

// NewEvent creates a new solidity event object using the signature
func NewEvent(name string) (*Event, error) {
	name, typ, err := parseEventOrErrorSignature("event ", name)
	if err != nil {
		return nil, err
	}
	return NewEventFromType(name, typ), nil
}

// NewEventFromType creates a new solidity event object using the name and type
func NewEventFromType(name string, typ *Type) *Event {
	return &Event{Name: name, Inputs: typ}
}

// Error is a solidity error object
type Error struct {
	Name   string
	Inputs *Type
}

// NewError creates a new solidity error object
func NewError(name string) (*Error, error) {
	name, typ, err := parseEventOrErrorSignature("error ", name)
	if err != nil {
		return nil, err
	}
	return &Error{Name: name, Inputs: typ}, nil
}

// ArgumentStr encodes a type object
type ArgumentStr struct {
	Name         string
	Type         string
	Indexed      bool
	Components   []*ArgumentStr
	InternalType string
}

var keccakPool = sync.Pool{
	New: func() interface{} {
		return sha3.NewLegacyKeccak256()
	},
}

func acquireKeccak() hash.Hash {
	return keccakPool.Get().(hash.Hash)
}

func releaseKeccak(k hash.Hash) {
	k.Reset()
	keccakPool.Put(k)
}

type Log struct {
	Removed          bool
	LogIndex         uint64
	TransactionIndex uint64
	TransactionHash  Hash
	BlockHash        Hash
	BlockNumber      uint64
	Address          Address
	Topics           []Hash
	Data             []byte
}

// nolint
func parseMethodSignature(name string) (string, *Type, *Type, error) {
	name = strings.Replace(name, "\n", " ", -1)
	name = strings.Replace(name, "\t", " ", -1)

	name = strings.TrimPrefix(name, "function ")
	name = strings.TrimSpace(name)

	var funcName, inputArgs, outputArgs string

	if strings.Contains(name, "returns") {
		matches := funcRegexpWithReturn.FindAllStringSubmatch(name, -1)
		if len(matches) == 0 {
			return "", nil, nil, fmt.Errorf("no matches found")
		}
		funcName = strings.TrimSpace(matches[0][1])
		inputArgs = strings.TrimSpace(matches[0][2])
		outputArgs = strings.TrimSpace(matches[0][4])
	} else {
		matches := funcRegexpWithoutReturn.FindAllStringSubmatch(name, -1)
		if len(matches) == 0 {
			return "", nil, nil, fmt.Errorf("no matches found")
		}
		funcName = strings.TrimSpace(matches[0][1])
		inputArgs = strings.TrimSpace(matches[0][2])
	}

	input, err := NewType("tuple(" + inputArgs + ")")
	if err != nil {
		return "", nil, nil, err
	}
	output, err := NewType("tuple(" + outputArgs + ")")
	if err != nil {
		return "", nil, nil, err
	}
	return funcName, input, output, nil
}

func buildSignature(name string, typ *Type) string {
	types := make([]string, len(typ.tuple))
	for i, input := range typ.tuple {
		// nolint
		types[i] = strings.Replace(input.Elem.String(), "tuple", "", -1)
	}
	return fmt.Sprintf("%v(%v)", name, strings.Join(types, ","))
}

func parseEventOrErrorSignature(prefix string, name string) (string, *Type, error) {
	if !strings.HasPrefix(name, prefix) {
		return "", nil, fmt.Errorf("prefix '%s' not found", prefix)
	}
	name = strings.TrimPrefix(name, prefix)

	if !strings.HasSuffix(name, ")") {
		return "", nil, fmt.Errorf("failed to parse input, expected 'name(types)'")
	}
	indx := strings.Index(name, "(")
	if indx == -1 {
		return "", nil, fmt.Errorf("failed to parse input, expected 'name(types)'")
	}

	funcName, signature := name[:indx], name[indx:]
	signature = "tuple" + signature

	typ, err := NewType(signature)
	if err != nil {
		return "", nil, err
	}
	return funcName, typ, nil
}
