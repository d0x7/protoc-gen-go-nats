package protoc_plugin

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/pkg/errors"
)

// Predefined Errors
var (
	ErrMarshallingFailed   = errors.New("Failed to marshal proto message")
	ErrUnmarshallingFailed = errors.New("Failed to unmarshal proto message")
)

// NATS Status Codes
// NATS status codes are used to represent the status of a NATS message.
// These codes are based on HTTP status codes, and those are not included here.
// We're trying to avoid any known HTTP status codes to prevent confusion, regardless if official or unofficial.
// For example, Shopify, Cloudflare, Microsoft and others are using some
// codes within the 520-559 range, which is why ours start from 560.
var (
	StatusProtobufProcessingFailed = "560"
)

// NATS Error

type NATSError struct {
	Code, Description string
	Wrapped           error
	Headers           map[string][]string
}

func (n NATSError) Error() string {
	if n.Wrapped != nil {
		return n.Description + ": " + n.Wrapped.Error()
	}
	return n.Description
}

func (n NATSError) Cause() error {
	return n.Wrapped
}

// GetWrapped returns the wrapped error as a byte slice, or nil if there is no wrapped error.
// It's therefore safe to be used directly in a NATS response, for example, like this:
// ```request.Error(natsErr.code, natsErr.description, natsErr.GetWrapped())```
func (n NATSError) GetWrapped() []byte {
	if n.Wrapped != nil {
		return []byte(n.Wrapped.Error())
	}
	return nil
}

func (n NATSError) ensureHeader() {
	if n.Headers == nil {
		n.Headers = make(nats.Header)
	}
}

func (n NATSError) GetOptHeaders() micro.RespondOpt {
	if n.Headers == nil || len(n.Headers) == 0 {
		return func(m *nats.Msg) {}
	}

	return func(m *nats.Msg) {
		if m.Header == nil {
			m.Header = n.Headers
			return
		}

		for k, v := range n.Headers {
			m.Header[k] = v
		}
	}
}

func (n NATSError) RespondWith(req micro.Request) error {
	return req.Error(n.Code, n.Description, n.GetWrapped(), n.GetOptHeaders())
}

func (n NATSError) AddHeader(header, value string) NATSError {
	n.ensureHeader()
	n.Headers[header] = append(n.Headers[header], value)
	return n
}

func (n NATSError) SetHeader(header, value string) NATSError {
	n.ensureHeader()
	n.Headers[header] = []string{value}
	return n
}

func (n NATSError) GetHeaders() micro.Headers {
	n.ensureHeader()
	return n.Headers
}

func (n NATSError) WithHeaders(headers map[string][]string) error {
	n.Headers = headers
	return n
}

func NewNATSErr(code, description string) NATSError {
	return WrapNATSErr(nil, code, description)
}

func WrapNATSErr(err error, code, description string) NATSError {
	return NATSError{
		Code:        code,
		Description: description,
		Wrapped:     err,
	}
}
