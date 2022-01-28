// example v0.0.1 11dde8c07720e169608bc087bb919a5761d7047f
// --
// This file has been generated by https://github.com/webrpc/webrpc using gen/golang
// Do not edit by hand. Update your webrpc schema and re-generate.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// WebRPC description and code-gen version
func WebRPCVersion() string {
	return "v1"
}

// Schema version of your RIDL schema
func WebRPCSchemaVersion() string {
	return "v0.0.1"
}

//
// Types
//

type Kind uint32

var (
	Kind_USER  Kind = 0
	Kind_ADMIN Kind = 1
)

var Kind_name = map[uint32]string{
	0: "USER",
	1: "ADMIN",
}

var Kind_value = map[string]uint32{
	"USER":  0,
	"ADMIN": 1,
}

func (x Kind) String() string {
	return Kind_name[uint32(x)]
}

func (x Kind) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString(`"`)
	buf.WriteString(Kind_name[uint32(x)])
	buf.WriteString(`"`)
	return buf.Bytes(), nil
}

func (x *Kind) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*x = Kind(Kind_value[j])
	return nil
}

func (x *Kind) UnmarshalText(b []byte) error {
	enum := string(b)
	*x = Kind(Kind_value[enum])
	return nil
}

func (x Kind) MarshalText() ([]byte, error) {
	return []byte(Kind_name[uint32(x)]), nil
}

type Empty struct {
}

type User struct {
	ID       uint64 `json:"id" db:"id"`
	Username string `json:"USERNAME" db:"username"`
	Role     string `json:"role" db:"-"`
}

type SearchFilter struct {
	Q string `json:"q"`
}

type Version struct {
	WebrpcVersion string `json:"webrpcVersion"`
	SchemaVersion string `json:"schemaVersion"`
	SchemaHash    string `json:"schemaHash"`
}

type ComplexType struct {
	Meta              map[string]interface{}       `json:"meta"`
	MetaNestedExample map[string]map[string]uint32 `json:"metaNestedExample"`
	NamesList         []string                     `json:"namesList"`
	NumsList          []int64                      `json:"numsList"`
	DoubleArray       [][]string                   `json:"doubleArray"`
	ListOfMaps        []map[string]uint32          `json:"listOfMaps"`
	ListOfUsers       []*User                      `json:"listOfUsers"`
	MapOfUsers        map[string]*User             `json:"mapOfUsers"`
	User              *User                        `json:"user"`
}

type ExampleService interface {
	Ping(ctx context.Context) error
	Status(ctx context.Context) (bool, error)
	Version(ctx context.Context) (*Version, error)
	GetUser(ctx context.Context, header map[string]string, userID uint64) (uint32, *User, error)
	FindUser(ctx context.Context, s *SearchFilter) (string, *User, error)
}

var WebRPCServices = map[string][]string{
	"ExampleService": {
		"Ping",
		"Status",
		"Version",
		"GetUser",
		"FindUser",
	},
}

//
// Server
//

type WebRPCServer interface {
	http.Handler
}

type ExampleServiceServer interface {
	ExampleService
}

type exampleServiceServer struct {
	service ExampleServiceServer
}

func NewExampleServiceServer(svc ExampleServiceServer) WebRPCServer {
	return &exampleServiceServer{
		service: svc,
	}
}

func (s *exampleServiceServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = context.WithValue(ctx, HTTPResponseWriterCtxKey, w)
	ctx = context.WithValue(ctx, HTTPRequestCtxKey, r)
	ctx = context.WithValue(ctx, ServiceNameCtxKey, "ExampleService")

	if r.Method != "POST" {
		RespondWithError(w, Errorf(ErrBadRoute, nil, "unsupported method %q (only POST is allowed)", r.Method))
		return
	}

	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		RespondWithError(w, Errorf(ErrBadRoute, nil, "unexpected Content-Type: %q", r.Header.Get("Content-Type")))
		return
	}

	switch r.URL.Path {
	case "/rpc/ExampleService/Ping":
		s.servePing(ctx, w, r)
		return
	case "/rpc/ExampleService/Status":
		s.serveStatus(ctx, w, r)
		return
	case "/rpc/ExampleService/Version":
		s.serveVersion(ctx, w, r)
		return
	case "/rpc/ExampleService/GetUser":
		s.serveGetUser(ctx, w, r)
		return
	case "/rpc/ExampleService/FindUser":
		s.serveFindUser(ctx, w, r)
		return
	default:
		RespondWithError(w, Errorf(ErrBadRoute, nil, "no handler for path %q", r.URL.Path))
		return
	}
}

func (s *exampleServiceServer) servePing(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var err error
	ctx = context.WithValue(ctx, MethodNameCtxKey, "Ping")

	// Call service method
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if rr := recover(); rr != nil {
				RespondWithError(w, ErrorInternal("internal service panic"))
				panic(rr)
			}
		}()
		err = s.service.Ping(ctx)
	}()

	if err != nil {
		RespondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (s *exampleServiceServer) serveStatus(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var err error
	ctx = context.WithValue(ctx, MethodNameCtxKey, "Status")

	// Call service method
	var ret0 bool
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if rr := recover(); rr != nil {
				RespondWithError(w, ErrorInternal("internal service panic"))
				panic(rr)
			}
		}()
		ret0, err = s.service.Status(ctx)
	}()
	respContent := struct {
		Ret0 bool `json:"status"`
	}{ret0}

	if err != nil {
		RespondWithError(w, err)
		return
	}
	respBody, err := json.Marshal(respContent)
	if err != nil {
		RespondWithError(w, Errorf(ErrInternal, err, "failed to marshal json response"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func (s *exampleServiceServer) serveVersion(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var err error
	ctx = context.WithValue(ctx, MethodNameCtxKey, "Version")

	// Call service method
	var ret0 *Version
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if rr := recover(); rr != nil {
				RespondWithError(w, ErrorInternal("internal service panic"))
				panic(rr)
			}
		}()
		ret0, err = s.service.Version(ctx)
	}()
	respContent := struct {
		Ret0 *Version `json:"version"`
	}{ret0}

	if err != nil {
		RespondWithError(w, err)
		return
	}
	respBody, err := json.Marshal(respContent)
	if err != nil {
		RespondWithError(w, Errorf(ErrInternal, err, "failed to marshal json response"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func (s *exampleServiceServer) serveGetUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var err error
	ctx = context.WithValue(ctx, MethodNameCtxKey, "GetUser")
	reqContent := struct {
		Arg0 map[string]string `json:"header"`
		Arg1 uint64            `json:"userID"`
	}{}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, Errorf(ErrInternal, err, "failed to read request data"))
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(reqBody, &reqContent)
	if err != nil {
		RespondWithError(w, Errorf(ErrInvalidArgument, err, "failed to unmarshal request data"))
		return
	}

	// Call service method
	var ret0 uint32
	var ret1 *User
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if rr := recover(); rr != nil {
				RespondWithError(w, ErrorInternal("internal service panic"))
				panic(rr)
			}
		}()
		ret0, ret1, err = s.service.GetUser(ctx, reqContent.Arg0, reqContent.Arg1)
	}()
	respContent := struct {
		Ret0 uint32 `json:"code"`
		Ret1 *User  `json:"user"`
	}{ret0, ret1}

	if err != nil {
		RespondWithError(w, err)
		return
	}
	respBody, err := json.Marshal(respContent)
	if err != nil {
		RespondWithError(w, Errorf(ErrInternal, err, "failed to marshal json response"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func (s *exampleServiceServer) serveFindUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var err error
	ctx = context.WithValue(ctx, MethodNameCtxKey, "FindUser")
	reqContent := struct {
		Arg0 *SearchFilter `json:"s"`
	}{}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, Errorf(ErrInternal, err, "failed to read request data"))
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(reqBody, &reqContent)
	if err != nil {
		RespondWithError(w, Errorf(ErrInvalidArgument, err, "failed to unmarshal request data"))
		return
	}

	// Call service method
	var ret0 string
	var ret1 *User
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if rr := recover(); rr != nil {
				RespondWithError(w, ErrorInternal("internal service panic"))
				panic(rr)
			}
		}()
		ret0, ret1, err = s.service.FindUser(ctx, reqContent.Arg0)
	}()
	respContent := struct {
		Ret0 string `json:"name"`
		Ret1 *User  `json:"user"`
	}{ret0, ret1}

	if err != nil {
		RespondWithError(w, err)
		return
	}
	respBody, err := json.Marshal(respContent)
	if err != nil {
		RespondWithError(w, Errorf(ErrInternal, err, "failed to marshal json response"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

//
// Server helpers
//

func RespondWithError(w http.ResponseWriter, err error) {
	e, ok := err.(Error)
	if !ok {
		e = Errorf(ErrInternal, err, err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(HTTPStatusFromError(err))
	respBody, _ := json.Marshal(e)
	w.Write(respBody)
}

//
// Client
//

type InjectHTTPRequestFunc func(*http.Request, context.Context) *http.Request

const ExampleServicePathPrefix = "/rpc/ExampleService/"

type ExampleServiceClient interface {
	ExampleService
}

type exampleServiceClient struct {
	client     HTTPClient
	urls       [5]string
	InjectFunc InjectHTTPRequestFunc
}

func NewExampleServiceClient(addr string, client HTTPClient, injectFunc InjectHTTPRequestFunc) ExampleServiceClient {
	prefix := urlBase(addr) + ExampleServicePathPrefix
	urls := [5]string{
		prefix + "Ping",
		prefix + "Status",
		prefix + "Version",
		prefix + "GetUser",
		prefix + "FindUser",
	}

	return &exampleServiceClient{
		client:     client,
		urls:       urls,
		InjectFunc: injectFunc,
	}
}

func (c *exampleServiceClient) Ping(ctx context.Context) error {
	_, err := clientRequest(ctx, c.client, c.urls[0], nil, nil, c.InjectFunction)
	return err
}

func (c *exampleServiceClient) Status(ctx context.Context) (bool, error) {
	out := struct {
		Ret0 bool `json:"status"`
	}{}
	_, err := clientRequest(ctx, c.client, c.urls[1], nil, &out, c.InjectFunction)
	return out.Ret0, err
}

func (c *exampleServiceClient) Version(ctx context.Context) (*Version, error) {
	out := struct {
		Ret0 *Version `json:"version"`
	}{}
	_, err := clientRequest(ctx, c.client, c.urls[2], nil, &out, c.InjectFunction)
	return out.Ret0, err
}

func (c *exampleServiceClient) GetUser(ctx context.Context, header map[string]string, userID uint64) (uint32, *User, error) {
	in := struct {
		Arg0 map[string]string `json:"header"`
		Arg1 uint64            `json:"userID"`
	}{header, userID}
	out := struct {
		Ret0 uint32 `json:"code"`
		Ret1 *User  `json:"user"`
	}{}
	_, err := clientRequest(ctx, c.client, c.urls[3], in, &out, c.InjectFunction)
	return out.Ret0, out.Ret1, err
}

func (c *exampleServiceClient) FindUser(ctx context.Context, s *SearchFilter) (string, *User, error) {
	in := struct {
		Arg0 *SearchFilter `json:"s"`
	}{s}
	out := struct {
		Ret0 string `json:"name"`
		Ret1 *User  `json:"user"`
	}{}
	_, err := clientRequest(ctx, c.client, c.urls[4], in, &out, c.InjectFunction)
	return out.Ret0, out.Ret1, err
}

//
// Client helpers
//

// HTTPClient is the interface used by generated clients to send HTTP requests.
// It is fulfilled by *(net/http).Client, which is sufficient for most users.
// Users can provide their own implementation for special retry policies.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// urlBase helps ensure that addr specifies a scheme. If it is unparsable
// as a URL, it returns addr unchanged.
func urlBase(addr string) string {
	// If the addr specifies a scheme, use it. If not, default to
	// http. If url.Parse fails on it, return it unchanged.
	url, err := url.Parse(addr)
	if err != nil {
		return addr
	}
	if url.Scheme == "" {
		url.Scheme = "http"
	}
	return url.String()
}

func clientRequest(ctx context.Context, client HTTPClient, url string, in, out interface{}, injectFunc InjectHTTPRequestFunc) (*http.Response, error) {
	reqBody, err := json.Marshal(in)
	if err != nil {
		return nil, Errorf(ErrInvalidArgument, err, "failed to marshal json request")
	}
	if err = ctx.Err(); err != nil {
		return nil, Errorf(ErrAborted, err, "aborted because context was done")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if headers, ok := GetClientRequestHeaders(ctx); ok {
		for k := range headers {
			for _, v := range headers[k] {
				req.Header.Add(k, v)
			}
		}
	}

	if injectFunc != nil {
		req = injectFunc(req, ctx)
	}

	resp, err := client.Do(req)
	if err != nil {
		return resp, Errorf(ErrFail, err, "request failed")
	}

	// auto-close body for non-streaming outputs
	if out != nil {
		defer func() {
			cerr := resp.Body.Close()
			if err == nil && cerr != nil {
				err = Errorf(ErrFail, cerr, "failed to close response body")
			}
		}()
	}

	if err = ctx.Err(); err != nil {
		return resp, Errorf(ErrAborted, err, "aborted because context was done")
	}

	if resp.StatusCode != 200 {
		return resp, errorFromResponse(resp)
	}

	if out != nil {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp, Errorf(ErrInternal, err, "failed to read response body")
		}

		err = json.Unmarshal(respBody, &out)
		if err != nil {
			return resp, Errorf(ErrInternal, err, "failed to unmarshal json response body")
		}
		if err = ctx.Err(); err != nil {
			return resp, Errorf(ErrAborted, err, "aborted because context was done")
		}
	}

	return resp, nil
}

// errorFromResponse builds a webrpc Error from a non-200 HTTP response.
func errorFromResponse(resp *http.Response) Error {
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Errorf(ErrInternal, err, "failed to read server error response body")
	}
	var respErr Error
	if err := json.Unmarshal(respBody, &respErr); err != nil {
		return Errorf(ErrInternal, err, "failed unmarshal error response")
	}
	return respErr
}

func WithClientRequestHeaders(ctx context.Context, h http.Header) (context.Context, error) {
	if _, ok := h["Accept"]; ok {
		return nil, errors.New("provided header cannot set Accept")
	}
	if _, ok := h["Content-Type"]; ok {
		return nil, errors.New("provided header cannot set Content-Type")
	}

	copied := make(http.Header, len(h))
	for k, vv := range h {
		if vv == nil {
			copied[k] = nil
			continue
		}
		copied[k] = make([]string, len(vv))
		copy(copied[k], vv)
	}

	return context.WithValue(ctx, HTTPClientRequestHeadersCtxKey, copied), nil
}

func GetClientRequestHeaders(ctx context.Context) (http.Header, bool) {
	h, ok := ctx.Value(HTTPClientRequestHeadersCtxKey).(http.Header)
	return h, ok
}

//
// Error helpers
//

type Error struct {
	Code    error  `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) Is(target error) bool {
	if errors.Is(target, e.Code) {
		return true
	}
	if e.Cause != nil && errors.Is(target, e.Cause) {
		return true
	}
	return false
}

func (e Error) Unwrap() error {
	if e.Cause != nil {
		return e.Cause
	} else {
		return e.Code
	}
}

func (e Error) MarshalJSON() ([]byte, error) {
	m, err := json.Marshal(e.Message)
	if err != nil {
		return nil, err
	}
	j := bytes.NewBufferString(`{`)
	j.WriteString(`"message": `)
	j.Write(m)
	j.WriteString(`}`)
	return j.Bytes(), nil
}

func (e *Error) UnmarshalJSON(b []byte) error {
	payload := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{}
	err := json.Unmarshal(b, &payload)
	if err != nil {
		return err
	}
	code := ErrorCodeFromString(payload.Code)
	if code == nil {
		code = ErrUnknown
	}
	*e = Error{
		Code:    code,
		Message: payload.Message,
	}
	return nil
}

var (
	// Fail indiciates a general error to processing a request.
	ErrFail = errors.New("fail")

	// Unknown error. For example when handling errors raised by APIs that do not
	// return enough error information.
	ErrUnknown = errors.New("unknown")

	// Internal errors. When some invariants expected by the underlying system
	// have been broken. In other words, something bad happened in the library or
	// backend service. Do not confuse with HTTP Internal Server Error; an
	// Internal error could also happen on the client code, i.e. when parsing a
	// server response.
	ErrInternal = errors.New("internal server error")

	// Unavailable indicates the service is currently unavailable. This is a most
	// likely a transient condition and may be corrected by retrying with a
	// backoff.
	ErrUnavailable = errors.New("unavailable")

	// Unsupported indicates the request was unsupported by the server. Perhaps
	// incorrect protocol version or missing feature.
	ErrUnsupported = errors.New("unsupported")

	// Canceled indicates the operation was cancelled (typically by the caller).
	ErrCanceled = errors.New("canceled")

	// InvalidArgument indicates client specified an invalid argument. It
	// indicates arguments that are problematic regardless of the state of the
	// system (i.e. a malformed file name, required argument, number out of range,
	// etc.).
	ErrInvalidArgument = errors.New("invalid argument")

	// DeadlineExceeded means operation expired before completion. For operations
	// that change the state of the system, this error may be returned even if the
	// operation has completed successfully (timeout).
	ErrDeadlineExceeded = errors.New("deadline exceeded")

	// NotFound means some requested entity was not found.
	ErrNotFound = errors.New("not found")

	// BadRoute means that the requested URL path wasn't routable to a webrpc
	// service and method. This is returned by the generated server, and usually
	// shouldn't be returned by applications. Instead, applications should use
	// NotFound or Unimplemented.
	ErrBadRoute = errors.New("bad route")

	// ErrMethodNotAllowed means that the requested URL path is available and the user
	// is authenticated and authorized. The input arguments are valid, but the
	// server needs to refuse the request for some reason
	ErrMethodNotAllowed = errors.New("method not allowed")

	// AlreadyExists means an attempt to create an entity failed because one
	// already exists.
	ErrAlreadyExists = errors.New("already exists")

	// PermissionDenied indicates the caller does not have permission to execute
	// the specified operation. It must not be used if the caller cannot be
	// identified (Unauthenticated).
	ErrPermissionDenied = errors.New("permission denied")

	// Unauthenticated indicates the request does not have valid authentication
	// credentials for the operation.
	ErrUnauthenticated = errors.New("unauthenticated")

	// ResourceExhausted indicates some resource has been exhausted, perhaps a
	// per-user quota, or perhaps the entire file system is out of space.
	ErrResourceExhausted = errors.New("resource exhausted")

	// Aborted indicates the operation was aborted, typically due to a concurrency
	// issue like sequencer check failures, transaction aborts, etc.
	ErrAborted = errors.New("aborted")

	// OutOfRange means operation was attempted past the valid range. For example,
	// seeking or reading past end of a paginated collection.
	ErrOutOfRange = errors.New("out of range")

	// Unimplemented indicates operation is not implemented or not
	// supported/enabled in this service.
	ErrUnimplemented = errors.New("unimplemented")

	// StreamClosed indicates that a connection stream has been closed.
	ErrStreamClosed = errors.New("stream closed")

	// StreamLost indiciates that a client or server connection has been interrupted
	// during an active transmission. It's a good idea to reconnect.
	ErrStreamLost = errors.New("stream lost")
)

func HTTPStatusFromError(err error) int {
	if errors.Is(err, ErrFail) {
		return 422 // Unprocessable Entity
	}
	if errors.Is(err, ErrUnknown) {
		return 400 // BadRequest
	}
	if errors.Is(err, ErrInternal) {
		return 500 // Internal Server Error
	}
	if errors.Is(err, ErrUnavailable) {
		return 503 // Service Unavailable
	}
	if errors.Is(err, ErrUnsupported) {
		return 500 // Internal Server Error
	}
	if errors.Is(err, ErrCanceled) {
		return 408 // RequestTimeout
	}
	if errors.Is(err, ErrInvalidArgument) {
		return 400 // BadRequest
	}
	if errors.Is(err, ErrDeadlineExceeded) {
		return 408 // RequestTimeout
	}
	if errors.Is(err, ErrNotFound) {
		return 404 // Not Found
	}
	if errors.Is(err, ErrBadRoute) {
		return 404 // Not Found
	}
	if errors.Is(err, ErrMethodNotAllowed) {
		return 405 // Method not allowed
	}
	if errors.Is(err, ErrAlreadyExists) {
		return 409 // Conflict
	}
	if errors.Is(err, ErrPermissionDenied) {
		return 403 // Forbidden
	}
	if errors.Is(err, ErrUnauthenticated) {
		return 401 // Unauthorized
	}
	if errors.Is(err, ErrResourceExhausted) {
		return 403 // Forbidden
	}
	if errors.Is(err, ErrAborted) {
		return 409 // Conflict
	}
	if errors.Is(err, ErrOutOfRange) {
		return 400 // Bad Request
	}
	if errors.Is(err, ErrUnimplemented) {
		return 501 // Not Implemented
	}
	if errors.Is(err, ErrStreamClosed) {
		return 200 // OK
	}
	if errors.Is(err, ErrStreamLost) {
		return 408 // RequestTimeout
	}
	return 500 // Invalid!
}

func ErrorCodeFromString(code string) error {
	switch code {
	case "fail":
		return ErrFail
	case "unknown":
		return ErrUnknown
	case "internal server error":
		return ErrInternal
	case "unavailable":
		return ErrUnavailable
	case "unsupported":
		return ErrUnsupported
	case "canceled":
		return ErrCanceled
	case "invalid argument":
		return ErrInvalidArgument
	case "deadline exceeded":
		return ErrDeadlineExceeded
	case "not found":
		return ErrNotFound
	case "bad route":
		return ErrBadRoute
	case "method not allowed":
		return ErrMethodNotAllowed
	case "already exists":
		return ErrAlreadyExists
	case "permission denied":
		return ErrPermissionDenied
	case "unauthenticated":
		return ErrUnauthenticated
	case "resource exhausted":
		return ErrResourceExhausted
	case "aborted":
		return ErrAborted
	case "out of range":
		return ErrOutOfRange
	case "unimplemented":
		return ErrUnimplemented
	case "stream closed":
		return ErrStreamClosed
	case "stream lost":
		return ErrStreamLost
	default:
		return nil
	}
}

func Errorf(code error, cause error, message string, args ...interface{}) Error {
	if ErrorCodeFromString(code.Error()) == nil {
		panic("invalid error code")
	}
	return Error{Code: code, Message: fmt.Sprintf(message, args...), Cause: cause}
}

func Failf(cause error, message string, args ...interface{}) Error {
	return Error{Code: ErrFail, Message: fmt.Sprintf(message, args...), Cause: cause}
}

func ErrorUnknown(message string, args ...interface{}) Error {
	return Errorf(ErrUnknown, nil, message, args...)
}

func ErrorInternal(message string, args ...interface{}) Error {
	return Errorf(ErrInternal, nil, message, args...)
}

func ErrorUnavailable(message string, args ...interface{}) Error {
	return Errorf(ErrUnavailable, nil, message, args...)
}

func ErrorUnsupported(message string, args ...interface{}) Error {
	return Errorf(ErrUnsupported, nil, message, args...)
}

func ErrorCanceled(message string, args ...interface{}) Error {
	return Errorf(ErrCanceled, nil, message, args...)
}

func ErrorInvalidArgument(message string, args ...interface{}) Error {
	return Errorf(ErrInvalidArgument, nil, message, args...)
}

func ErrorDeadlineExceeded(message string, args ...interface{}) Error {
	return Errorf(ErrDeadlineExceeded, nil, message, args...)
}

func ErrorNotFound(message string, args ...interface{}) Error {
	return Errorf(ErrNotFound, nil, message, args...)
}

func ErrorBadRoute(message string, args ...interface{}) Error {
	return Errorf(ErrBadRoute, nil, message, args...)
}

func ErrorMethodNotAllowed(message string, args ...interface{}) Error {
	return Errorf(ErrMethodNotAllowed, nil, message, args...)
}

func ErrorAlreadyExists(message string, args ...interface{}) Error {
	return Errorf(ErrAlreadyExists, nil, message, args...)
}

func ErrorPermissionDenied(message string, args ...interface{}) Error {
	return Errorf(ErrPermissionDenied, nil, message, args...)
}

func ErrorUnauthenticated(message string, args ...interface{}) Error {
	return Errorf(ErrUnauthenticated, nil, message, args...)
}

func ErrorResourceExhausted(message string, args ...interface{}) Error {
	return Errorf(ErrResourceExhausted, nil, message, args...)
}

func ErrorAborted(message string, args ...interface{}) Error {
	return Errorf(ErrAborted, nil, message, args...)
}

func ErrorOutOfRange(message string, args ...interface{}) Error {
	return Errorf(ErrOutOfRange, nil, message, args...)
}

func ErrorUnimplemented(message string, args ...interface{}) Error {
	return Errorf(ErrUnimplemented, nil, message, args...)
}

func ErrorStreamClosed(message string, args ...interface{}) Error {
	return Errorf(ErrStreamClosed, nil, message, args...)
}

func ErrorStreamLost(message string, args ...interface{}) Error {
	return Errorf(ErrStreamLost, nil, message, args...)
}

func GetErrorStack(err error) []error {
	errs := []error{err}
	for {
		unwrap, ok := err.(interface{ Unwrap() error })
		if !ok {
			break
		}
		werr := unwrap.Unwrap()
		if werr == nil {
			break
		}
		errs = append(errs, werr)
		err = werr
	}
	return errs
}

//
// Misc helpers
//

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "webrpc context value " + k.name
}

var (
	// For Client
	HTTPClientRequestHeadersCtxKey = &contextKey{"HTTPClientRequestHeaders"}

	// For Server
	HTTPResponseWriterCtxKey = &contextKey{"HTTPResponseWriter"} // http.ResponseWriter
	HTTPRequestCtxKey        = &contextKey{"HTTPRequest"}        // *http.Request
	ServiceNameCtxKey        = &contextKey{"ServiceName"}        // string
	MethodNameCtxKey         = &contextKey{"MethodName"}         // string
)
