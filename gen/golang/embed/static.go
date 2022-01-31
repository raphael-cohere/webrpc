// Code generated by statik. DO NOT EDIT.

// Package contains static assets.
package embed

var	Asset = "PK\x03\x04\x14\x00\x08\x00\x00\x00W\x18?T\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0e\x00	\x00client.go.tmplUT\x05\x00\x01WQ\xf7a{{define \"client\"}}\n{{if .Services}}\n//\n// Client\n//\n\ntype InjectHTTPRequestFunc func(*http.Request, context.Context) *http.Request\n\n{{range .Services}}\nconst {{.Name | constPathPrefix}} = \"/rpc/{{.Name}}/\"\n{{end}}\n\n{{range .Services}}\n  type {{.Name}}Client interface {\n    {{.Name}}\n    {{- range .Methods}}\n      {{- if .StreamOutput }}\n        {{.Name}}({{.Inputs | methodInputs}}) ({{.Name}}StreamReader, error)\n      {{- end}}\n    {{end}}\n  }\n{{end}}\n{{range .Services}}\n  {{ $serviceName := .Name | clientServiceName}}\n  type {{$serviceName}} struct {\n    client         HTTPClient\n    urls           [{{.Methods | countMethods}}]string\n    InjectFunc InjectHTTPRequestFunc\n  }\n\n  func {{.Name | newClientServiceName }}(addr string, client HTTPClient, injectFunc InjectHTTPRequestFunc) {{.Name}}Client {\n    prefix := urlBase(addr) + {{.Name | constPathPrefix}}\n    urls := [{{.Methods | countMethods}}]string{\n      {{- range .Methods}}\n      prefix + \"{{.Name}}\",\n      {{- end}}\n    }\n\n    return &{{$serviceName}}{\n      client: client,\n      urls:   urls,\n      InjectFunc: injectFunc,\n    }\n  }\n\n  {{range $i, $method := .Methods}}\n    {{- if .StreamOutput}}\n    func (c *{{$serviceName}}) {{.Name}}({{.Inputs | methodInputs}}) ({{.Name}}StreamReader, error) {\n    {{- else}}\n    func (c *{{$serviceName}}) {{.Name}}({{.Inputs | methodInputs}}) ({{.Outputs | methodOutputs }}) {\n    {{- end}}\n      {{- $inputVar := \"nil\" -}}\n      {{- $outputVar := \"nil\" -}}\n      {{- if .Inputs | len}}\n      {{- $inputVar = \"in\"}}\n      in := struct {\n        {{- range $i, $input := .Inputs}}\n          Arg{{$i}} {{$input | methodArgType}} `json:\"{{$input.Name | downcaseName}}\"`\n        {{- end}}\n      }{ {{.Inputs | methodArgNames}} }\n      {{- end}}\n      {{- if .StreamOutput}}\n        resp, err := clientRequest(ctx, c.client, c.urls[{{$i}}], {{$inputVar}}, {{$outputVar}}, c.InjectFunc)\n        if err != nil {\n          return nil, err\n        }\n\n        return newClient{{.Name}}StreamReader(resp), nil\n      {{- else}}\n        {{- if .Outputs | len}}\n        {{- $outputVar = \"&out\"}}\n        out := struct {\n          {{- range $i, $output := .Outputs}}\n            Ret{{$i}} {{$output | methodArgType}} `json:\"{{$output.Name | downcaseName}}\"`\n          {{- end}}\n        }{}\n        {{- end}}\n        _, err := clientRequest(ctx, c.client, c.urls[{{$i}}], {{$inputVar}}, {{$outputVar}}, c.InjectFunc)\n        return {{argsList .Outputs \"out.Ret\"}}{{commaIfLen .Outputs}} err\n      {{- end}}\n    }\n  {{end}}\n{{end}}\n\n//\n// Client helpers\n//\n\n// HTTPClient is the interface used by generated clients to send HTTP requests.\n// It is fulfilled by *(net/http).Client, which is sufficient for most users.\n// Users can provide their own implementation for special retry policies.\ntype HTTPClient interface {\n  Do(req *http.Request) (*http.Response, error)\n}\n\n// urlBase helps ensure that addr specifies a scheme. If it is unparsable\n// as a URL, it returns addr unchanged.\nfunc urlBase(addr string) string {\n  // If the addr specifies a scheme, use it. If not, default to\n  // http. If url.Parse fails on it, return it unchanged.\n  url, err := url.Parse(addr)\n  if err != nil {\n    return addr\n  }\n  if url.Scheme == \"\" {\n    url.Scheme = \"http\"\n  }\n  return url.String()\n}\n\nfunc clientRequest(ctx context.Context, client HTTPClient, url string, in, out interface{}, injectFunc InjectHTTPRequestFunc) (*http.Response, error) {\n  reqBody, err := json.Marshal(in)\n  if err != nil {\n    return nil, Errorf(ErrInvalidArgument, err, \"failed to marshal json request\")\n  }\n  if err = ctx.Err(); err != nil {\n    return nil, Errorf(ErrAborted, err, \"aborted because context was done\")\n  }\n\n  req, err := http.NewRequest(\"POST\", url, bytes.NewBuffer(reqBody))\n  if err != nil {\n    return nil, err\n  }\n  req.Header.Set(\"Content-Type\", \"application/json\")\n  if headers, ok := GetClientRequestHeaders(ctx); ok {\n    for k := range headers {\n      for _, v := range headers[k] {\n        req.Header.Add(k, v)\n      }\n    }\n  }\n\n  if injectFunc != nil {\n    req = injectFunc(req, ctx)\n  }\n\n  resp, err := client.Do(req)\n  if err != nil {\n    return resp, Errorf(ErrFail, err, \"request failed\")\n  }\n\n  // auto-close body for non-streaming outputs\n  if out != nil {\n    defer func() {\n      cerr := resp.Body.Close()\n      if err == nil && cerr != nil {\n        err = Errorf(ErrFail, cerr, \"failed to close response body\")\n      }\n    }()\n  }\n\n  if err = ctx.Err(); err != nil {\n    return resp, Errorf(ErrAborted, err, \"aborted because context was done\")\n  }\n\n  if resp.StatusCode != 200 {\n    return resp, errorFromResponse(resp)\n  }\n\n  if out != nil {\n    respBody, err := ioutil.ReadAll(resp.Body)\n    if err != nil {\n      return resp, Errorf(ErrInternal, err, \"failed to read response body\")\n    }\n\n    err = json.Unmarshal(respBody, &out)\n    if err != nil {\n      return resp, Errorf(ErrInternal, err, \"failed to unmarshal json response body\")\n    }\n    if err = ctx.Err(); err != nil {\n      return resp, Errorf(ErrAborted, err, \"aborted because context was done\")\n    }\n  }\n\n  return resp, nil\n}\n\n// errorFromResponse builds a webrpc Error from a non-200 HTTP response.\nfunc errorFromResponse(resp *http.Response) Error {\n  respBody, err := ioutil.ReadAll(resp.Body)\n  if err != nil {\n    return Errorf(ErrInternal, err, \"failed to read server error response body\")\n  }\n  var respErr Error\n  if err := json.Unmarshal(respBody, &respErr); err != nil {\n    return Errorf(ErrInternal, err, \"failed unmarshal error response\")\n  }\n  return respErr\n}\n\nfunc WithClientRequestHeaders(ctx context.Context, h http.Header) (context.Context, error) {\n  if _, ok := h[\"Accept\"]; ok {\n    return nil, errors.New(\"provided header cannot set Accept\")\n  }\n  if _, ok := h[\"Content-Type\"]; ok {\n    return nil, errors.New(\"provided header cannot set Content-Type\")\n  }\n\n  copied := make(http.Header, len(h))\n  for k, vv := range h {\n    if vv == nil {\n      copied[k] = nil\n      continue\n    }\n    copied[k] = make([]string, len(vv))\n    copy(copied[k], vv)\n  }\n\n  return context.WithValue(ctx, HTTPClientRequestHeadersCtxKey, copied), nil\n}\n\nfunc GetClientRequestHeaders(ctx context.Context) (http.Header, bool) {\n  h, ok := ctx.Value(HTTPClientRequestHeadersCtxKey).(http.Header)\n  return h, ok\n}\n\n{{- if .Services | hasStreamOutput}}\n  //\n  // Client streaming helpers\n  //\n{{- end}}\n{{range .Services}}\n  {{- range .Methods}}\n  {{- if .StreamOutput}}\n    type client{{.Name}}StreamReader struct {\n      resp    *http.Response\n      reader  io.Reader\n      decoder *json.Decoder\n    }\n\n    func newClient{{.Name}}StreamReader(resp *http.Response) *client{{.Name}}StreamReader {\n      reader := httputil.NewChunkedReader(resp.Body)\n      decoder := json.NewDecoder(reader)\n      return &client{{.Name}}StreamReader{\n        resp: resp, reader: reader, decoder: decoder,\n      }\n    }\n\n    func (c *client{{.Name}}StreamReader) Read() ({{.Outputs | methodOutputsWithTypes}}, err error) {\n      for {\n        out := struct {\n          {{- range $i, $output := .Outputs}}\n            Ret{{$i}} {{$output | methodArgType}} `json:\"{{$output.Name | downcaseName}}\"`\n          {{- end}}\n          Error Error `json:\"error\"`\n          Ping  bool  `json:\"ping\"`\n        }{}\n\n        err = c.decoder.Decode(&out)\n\n        // Skip ping payloads\n        if err == nil && out.Ping {\n          continue\n        }\n\n        // Error checking\n        if err != nil {\n          if err == io.EOF {\n            return {{argsList .Outputs \"out.Ret\"}}{{commaIfLen .Outputs}} Errorf(ErrStreamClosed, err, err.Error())\n          }\n          return {{argsList .Outputs \"out.Ret\"}}{{commaIfLen .Outputs}} Errorf(ErrStreamLost, err, err.Error())\n        }\n        if out.Error.Code != nil || out.Error.Message != \"\" {\n          return {{argsList .Outputs \"out.Ret\"}}{{commaIfLen .Outputs}} out.Error\n        }\n\n        return {{argsList .Outputs \"out.Ret\"}}{{commaIfLen .Outputs}} nil\n      }\n    }\n  {{- end}}\n  {{- end}}\n{{- end}}\n{{end}}\n{{end}}\nPK\x07\x08\xb1{\x92\x13H\x1f\x00\x00H\x1f\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x98\x88-T\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x00	\x00helpers.go.tmplUT\x05\x00\x01\xb1[\xe0a{{define \"helpers\"}}\n\n//\n// Error helpers\n//\n\ntype Error struct {\n	Code    error  `json:\"code\"`\n	Message string `json:\"message\"`\n	Cause   error  `json:\"-\"`\n}\n\nfunc (e Error) Error() string {\n	return e.Message\n}\n\nfunc (e Error) Is(target error) bool {\n	if errors.Is(target, e.Code) {\n		return true\n	}\n	if e.Cause != nil && errors.Is(target, e.Cause) {\n		return true\n	}\n	return false\n}\n\nfunc (e Error) Unwrap() error {\n	if e.Cause != nil {\n		return e.Cause\n	} else {\n		return e.Code\n	}\n}\n\nfunc (e Error) MarshalJSON() ([]byte, error) {\n	m, err := json.Marshal(e.Message)\n	if err != nil {\n		return nil, err\n	}\n	j := bytes.NewBufferString(`{`)\n	j.WriteString(`\"message\": `)\n	j.Write(m)\n	j.WriteString(`}`)\n	return j.Bytes(), nil\n}\n\nfunc (e *Error) UnmarshalJSON(b []byte) error {\n	payload := struct {\n		Code    string `json:\"code\"`\n		Message string `json:\"message\"`\n	}{}\n	err := json.Unmarshal(b, &payload)\n	if err != nil {\n		return err\n	}\n	code := ErrorCodeFromString(payload.Code)\n	if code == nil {\n		code = ErrUnknown\n	}\n	*e = Error{\n		Code:    code,\n		Message: payload.Message,\n	}\n	return nil\n}\n\nvar (\n	// Fail indiciates a general error to processing a request.\n	ErrFail = errors.New(\"fail\")\n\n	// Unknown error. For example when handling errors raised by APIs that do not\n	// return enough error information.\n	ErrUnknown = errors.New(\"unknown\")\n\n	// Internal errors. When some invariants expected by the underlying system\n	// have been broken. In other words, something bad happened in the library or\n	// backend service. Do not confuse with HTTP Internal Server Error; an\n	// Internal error could also happen on the client code, i.e. when parsing a\n	// server response.\n	ErrInternal = errors.New(\"internal server error\")\n\n	// Unavailable indicates the service is currently unavailable. This is a most\n	// likely a transient condition and may be corrected by retrying with a\n	// backoff.\n	ErrUnavailable = errors.New(\"unavailable\")\n\n	// Unsupported indicates the request was unsupported by the server. Perhaps\n	// incorrect protocol version or missing feature.\n	ErrUnsupported = errors.New(\"unsupported\")\n\n	// Canceled indicates the operation was cancelled (typically by the caller).\n	ErrCanceled = errors.New(\"canceled\")\n\n	// InvalidArgument indicates client specified an invalid argument. It\n	// indicates arguments that are problematic regardless of the state of the\n	// system (i.e. a malformed file name, required argument, number out of range,\n	// etc.).\n	ErrInvalidArgument = errors.New(\"invalid argument\")\n\n	// DeadlineExceeded means operation expired before completion. For operations\n	// that change the state of the system, this error may be returned even if the\n	// operation has completed successfully (timeout).\n	ErrDeadlineExceeded = errors.New(\"deadline exceeded\")\n\n	// NotFound means some requested entity was not found.\n	ErrNotFound = errors.New(\"not found\")\n\n	// BadRoute means that the requested URL path wasn't routable to a webrpc\n	// service and method. This is returned by the generated server, and usually\n	// shouldn't be returned by applications. Instead, applications should use\n	// NotFound or Unimplemented.\n	ErrBadRoute = errors.New(\"bad route\")\n\n	// ErrMethodNotAllowed means that the requested URL path is available and the user\n	// is authenticated and authorized. The input arguments are valid, but the\n	// server needs to refuse the request for some reason\n	ErrMethodNotAllowed = errors.New(\"method not allowed\")\n\n	// AlreadyExists means an attempt to create an entity failed because one\n	// already exists.\n	ErrAlreadyExists = errors.New(\"already exists\")\n\n	// PermissionDenied indicates the caller does not have permission to execute\n	// the specified operation. It must not be used if the caller cannot be\n	// identified (Unauthenticated).\n	ErrPermissionDenied = errors.New(\"permission denied\")\n\n	// Unauthenticated indicates the request does not have valid authentication\n	// credentials for the operation.\n	ErrUnauthenticated = errors.New(\"unauthenticated\")\n\n	// ResourceExhausted indicates some resource has been exhausted, perhaps a\n	// per-user quota, or perhaps the entire file system is out of space.\n	ErrResourceExhausted = errors.New(\"resource exhausted\")\n\n	// Aborted indicates the operation was aborted, typically due to a concurrency\n	// issue like sequencer check failures, transaction aborts, etc.\n	ErrAborted = errors.New(\"aborted\")\n\n	// OutOfRange means operation was attempted past the valid range. For example,\n	// seeking or reading past end of a paginated collection.\n	ErrOutOfRange = errors.New(\"out of range\")\n\n	// Unimplemented indicates operation is not implemented or not\n	// supported/enabled in this service.\n	ErrUnimplemented = errors.New(\"unimplemented\")\n\n	// StreamClosed indicates that a connection stream has been closed.\n	ErrStreamClosed = errors.New(\"stream closed\")\n\n	// StreamLost indiciates that a client or server connection has been interrupted\n	// during an active transmission. It's a good idea to reconnect.\n	ErrStreamLost = errors.New(\"stream lost\")\n)\n\nfunc HTTPStatusFromError(err error) int {\n	if errors.Is(err, ErrFail) {\n		return 422 // Unprocessable Entity\n	}\n	if errors.Is(err, ErrUnknown) {\n		return 400 // BadRequest\n	}\n	if errors.Is(err, ErrInternal) {\n		return 500 // Internal Server Error\n	}\n	if errors.Is(err, ErrUnavailable) {\n		return 503 // Service Unavailable\n	}\n	if errors.Is(err, ErrUnsupported) {\n		return 500 // Internal Server Error\n	}\n	if errors.Is(err, ErrCanceled) {\n		return 408 // RequestTimeout\n	}\n	if errors.Is(err, ErrInvalidArgument) {\n		return 400 // BadRequest\n	}\n	if errors.Is(err, ErrDeadlineExceeded) {\n		return 408 // RequestTimeout\n	}\n	if errors.Is(err, ErrNotFound) {\n		return 404 // Not Found\n	}\n	if errors.Is(err, ErrBadRoute) {\n		return 404 // Not Found\n	}\n	if errors.Is(err, ErrMethodNotAllowed) {\n		return 405 // Method not allowed\n	}\n	if errors.Is(err, ErrAlreadyExists) {\n		return 409 // Conflict\n	}\n	if errors.Is(err, ErrPermissionDenied) {\n		return 403 // Forbidden\n	}\n	if errors.Is(err, ErrUnauthenticated) {\n		return 401 // Unauthorized\n	}\n	if errors.Is(err, ErrResourceExhausted) {\n		return 403 // Forbidden\n	}\n	if errors.Is(err, ErrAborted) {\n		return 409 // Conflict\n	}\n	if errors.Is(err, ErrOutOfRange) {\n		return 400 // Bad Request\n	}\n	if errors.Is(err, ErrUnimplemented) {\n		return 501 // Not Implemented\n	}\n	if errors.Is(err, ErrStreamClosed) {\n		return 200 // OK\n	}\n	if errors.Is(err, ErrStreamLost) {\n		return 408 // RequestTimeout\n	}\n	return 500 // Invalid!\n}\n\nfunc ErrorCodeFromString(code string) error {\n	switch code {\n	case \"fail\":\n		return ErrFail\n	case \"unknown\":\n		return ErrUnknown\n	case \"internal server error\":\n		return ErrInternal\n	case \"unavailable\":\n		return ErrUnavailable\n	case \"unsupported\":\n		return ErrUnsupported\n	case \"canceled\":\n		return ErrCanceled\n	case \"invalid argument\":\n		return ErrInvalidArgument\n	case \"deadline exceeded\":\n		return ErrDeadlineExceeded\n	case \"not found\":\n		return ErrNotFound\n	case \"bad route\":\n		return ErrBadRoute\n	case \"method not allowed\":\n		return ErrMethodNotAllowed\n	case \"already exists\":\n		return ErrAlreadyExists\n	case \"permission denied\":\n		return ErrPermissionDenied\n	case \"unauthenticated\":\n		return ErrUnauthenticated\n	case \"resource exhausted\":\n		return ErrResourceExhausted\n	case \"aborted\":\n		return ErrAborted\n	case \"out of range\":\n		return ErrOutOfRange\n	case \"unimplemented\":\n		return ErrUnimplemented\n	case \"stream closed\":\n		return ErrStreamClosed\n	case \"stream lost\":\n		return ErrStreamLost\n	default:\n		return nil\n	}\n}\n\nfunc Errorf(code error, cause error, message string, args ...interface{}) Error {\n	if ErrorCodeFromString(code.Error()) == nil {\n		panic(\"invalid error code\")\n	}\n	return Error{Code: code, Message: fmt.Sprintf(message, args...), Cause: cause}\n}\n\nfunc Failf(cause error, message string, args ...interface{}) Error {\n	return Error{Code: ErrFail, Message: fmt.Sprintf(message, args...), Cause: cause}\n}\n\nfunc ErrorUnknown(message string, args ...interface{}) Error {\n	return Errorf(ErrUnknown, nil, message, args...)\n}\n\nfunc ErrorInternal(message string, args ...interface{}) Error {\n	return Errorf(ErrInternal, nil, message, args...)\n}\n\nfunc ErrorUnavailable(message string, args ...interface{}) Error {\n	return Errorf(ErrUnavailable, nil, message, args...)\n}\n\nfunc ErrorUnsupported(message string, args ...interface{}) Error {\n	return Errorf(ErrUnsupported, nil, message, args...)\n}\n\nfunc ErrorCanceled(message string, args ...interface{}) Error {\n	return Errorf(ErrCanceled, nil, message, args...)\n}\n\nfunc ErrorInvalidArgument(message string, args ...interface{}) Error {\n	return Errorf(ErrInvalidArgument, nil, message, args...)\n}\n\nfunc ErrorDeadlineExceeded(message string, args ...interface{}) Error {\n	return Errorf(ErrDeadlineExceeded, nil, message, args...)\n}\n\nfunc ErrorNotFound(message string, args ...interface{}) Error {\n	return Errorf(ErrNotFound, nil, message, args...)\n}\n\nfunc ErrorBadRoute(message string, args ...interface{}) Error {\n	return Errorf(ErrBadRoute, nil, message, args...)\n}\n\nfunc ErrorMethodNotAllowed(message string, args ...interface{}) Error {\n	return Errorf(ErrMethodNotAllowed, nil, message, args...)\n}\n\nfunc ErrorAlreadyExists(message string, args ...interface{}) Error {\n	return Errorf(ErrAlreadyExists, nil, message, args...)\n}\n\nfunc ErrorPermissionDenied(message string, args ...interface{}) Error {\n	return Errorf(ErrPermissionDenied, nil, message, args...)\n}\n\nfunc ErrorUnauthenticated(message string, args ...interface{}) Error {\n	return Errorf(ErrUnauthenticated, nil, message, args...)\n}\n\nfunc ErrorResourceExhausted(message string, args ...interface{}) Error {\n	return Errorf(ErrResourceExhausted, nil, message, args...)\n}\n\nfunc ErrorAborted(message string, args ...interface{}) Error {\n	return Errorf(ErrAborted, nil, message, args...)\n}\n\nfunc ErrorOutOfRange(message string, args ...interface{}) Error {\n	return Errorf(ErrOutOfRange, nil, message, args...)\n}\n\nfunc ErrorUnimplemented(message string, args ...interface{}) Error {\n	return Errorf(ErrUnimplemented, nil, message, args...)\n}\n\nfunc ErrorStreamClosed(message string, args ...interface{}) Error {\n	return Errorf(ErrStreamClosed, nil, message, args...)\n}\n\nfunc ErrorStreamLost(message string, args ...interface{}) Error {\n	return Errorf(ErrStreamLost, nil, message, args...)\n}\n\nfunc GetErrorStack(err error) []error {\n	errs := []error{err}\n	for {\n		unwrap, ok := err.(interface{ Unwrap() error })\n		if !ok {\n			break\n		}\n		werr := unwrap.Unwrap()\n		if werr == nil {\n			break\n		}\n		errs = append(errs, werr)\n		err = werr\n	}\n	return errs\n}\n\n//\n// Misc helpers\n//\n\ntype contextKey struct {\n	name string\n}\n\nfunc (k *contextKey) String() string {\n	return \"webrpc context value \" + k.name\n}\n\nvar (\n	// For Client\n	HTTPClientRequestHeadersCtxKey = &contextKey{\"HTTPClientRequestHeaders\"}\n\n	// For Server\n	HTTPResponseWriterCtxKey = &contextKey{\"HTTPResponseWriter\"} // http.ResponseWriter\n	HTTPRequestCtxKey        = &contextKey{\"HTTPRequest\"}        // *http.Request\n	ServiceNameCtxKey        = &contextKey{\"ServiceName\"}        // string\n	MethodNameCtxKey         = &contextKey{\"MethodName\"}         // string\n)\n{{end}}\nPK\x07\x08\xd8\xdb\x18\xad\x8a+\x00\x00\x8a+\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x98\x88-T\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x11\x00	\x00proto.gen.go.tmplUT\x05\x00\x01\xb1[\xe0a{{- define \"proto\" -}}\n// {{.Name}} {{.SchemaVersion}} {{.SchemaHash}}\n// --\n// This file has been generated by https://github.com/webrpc/webrpc using gen/golang\n// Do not edit by hand. Update your webrpc schema and re-generate.\npackage {{.TargetOpts.PkgName}}\n\nimport (\n  \"context\"\n  \"encoding/json\"\n  \"fmt\"\n  \"io/ioutil\"\n  \"net/http\"\n  \"time\"\n  \"strings\"\n  \"bytes\"\n  \"errors\"\n  \"io\"\n  \"net/url\"\n)\n\n// WebRPC description and code-gen version\nfunc WebRPCVersion() string {\n  return \"{{.WebRPCVersion}}\"\n}\n\n// Schema version of your RIDL schema\nfunc WebRPCSchemaVersion() string {\n  return \"{{.SchemaVersion}}\"\n}\n\n{{template \"types\" .}}\n\n{{if .TargetOpts.Server}}\n  {{template \"server\" .}}\n{{end}}\n\n{{if .TargetOpts.Client}}\n  {{template \"client\" .}}\n{{end}}\n\n{{template \"helpers\" .}}\n\n{{- end}}\nPK\x07\x08\x0eD\x8eO\x1b\x03\x00\x00\x1b\x03\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x98\x88-T\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0e\x00	\x00server.go.tmplUT\x05\x00\x01\xb1[\xe0a{{define \"server\"}}\n{{if .Services}}\n//\n// Server\n//\n\ntype WebRPCServer interface {\n  http.Handler\n}\n\n{{- range .Services}}\n  {{$name := .Name}}\n  {{$serviceName := .Name | serverServiceName}}\n\n  type {{.Name}}Server interface {\n    {{.Name}}\n    {{- range .Methods}}\n      {{- if .StreamOutput }}\n            {{.Name}}({{ .Inputs | methodInputs }}, stream {{.Name}}StreamWriter) error\n      {{- end}}\n    {{- end}}\n  }\n\n  type {{$serviceName}} struct {\n    service {{.Name}}Server\n  }\n\n  func {{ .Name | newServerServiceName }}(svc {{.Name}}Server) WebRPCServer {\n    return &{{$serviceName}}{\n      service: svc,\n    }\n  }\n\n  func (s *{{$serviceName}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {\n    ctx := r.Context()\n    ctx = context.WithValue(ctx, HTTPResponseWriterCtxKey, w)\n    ctx = context.WithValue(ctx, HTTPRequestCtxKey, r)\n    ctx = context.WithValue(ctx, ServiceNameCtxKey, \"{{.Name}}\")\n\n    if r.Method != \"POST\" {\n      RespondWithError(w, Errorf(ErrBadRoute, nil, \"unsupported method %q (only POST is allowed)\", r.Method))\n      return\n    }\n\n    if !strings.HasPrefix(r.Header.Get(\"Content-Type\"), \"application/json\") {\n      RespondWithError(w, Errorf(ErrBadRoute, nil, \"unexpected Content-Type: %q\", r.Header.Get(\"Content-Type\")))\n      return\n    }\n\n    switch r.URL.Path {\n    {{- range .Methods}}\n    case \"/rpc/{{$name}}/{{.Name}}\":\n      s.{{.Name | serviceMethodName}}(ctx, w, r)\n      return\n    {{- end}}\n    default:\n      RespondWithError(w, Errorf(ErrBadRoute, nil, \"no handler for path %q\", r.URL.Path))\n      return\n    }\n  }\n\n  {{range .Methods}}\n    func (s *{{$serviceName}}) {{.Name | serviceMethodName}}(ctx context.Context, w http.ResponseWriter, r *http.Request) {\n      var err error\n      ctx = context.WithValue(ctx, MethodNameCtxKey, \"{{.Name}}\")\n\n\n      {{- if .Inputs|len}}\n      reqContent := struct {\n      {{- range $i, $input := .Inputs}}\n        Arg{{$i}} {{. | methodArgType}} `json:\"{{$input.Name | downcaseName}}\"`\n      {{- end}}\n      }{}\n\n      reqBody, err := ioutil.ReadAll(r.Body)\n      if err != nil {\n        RespondWithError(w, Errorf(ErrInternal, err, \"failed to read request data\"))\n        return\n      }\n      defer r.Body.Close()\n\n      err = json.Unmarshal(reqBody, &reqContent)\n      if err != nil {\n        RespondWithError(w, Errorf(ErrInvalidArgument, err, \"failed to unmarshal request data\"))\n        return\n      }\n      {{- end}}\n\n\n      // Call service method\n      {{- if .StreamOutput}}\n        sw, err := newServerStreamWriter(w)\n        if err != nil {\n          RespondWithError(w, Errorf(ErrUnsupported, err, \"http connection does not support streams\"))\n          return\n        }\n\n        streamWriter := &{{.Name | streamWriterName}}{sw}\n\n        // connection monitoring and keep-alive\n        go func() {\n          for {\n            select {\n            case <-time.After(StreamKeepAliveInterval):\n              streamWriter.Ping()\n            case <-r.Context().Done():\n              streamWriter.Close()\n              return\n            case <-streamWriter.Done():\n              return\n            }\n          }\n        }()\n\n        func() {\n          defer func() {\n            // In case of a panic, serve a error chunk and then panic.\n            if rr := recover(); rr != nil {\n              streamWriter.Error(ErrorInternal(\"internal service panic\"))\n              streamWriter.Close()\n              panic(rr)\n            }\n          }()\n          err = s.service.{{.Name}}(ctx{{.Inputs | commaIfLen}}{{argsList .Inputs \"reqContent.Arg\"}}, streamWriter)\n        }()\n\n        if err != nil {\n          streamWriter.Error(err) // the error to the client\n        }\n        streamWriter.Close() // always ensure we close the stream\n      {{- else}}\n        {{- range $i, $output := .Outputs}}\n        var ret{{$i}} {{$output | methodArgType}}\n        {{- end}}\n        func() {\n          defer func() {\n            // In case of a panic, serve a 500 error and then panic.\n            if rr := recover(); rr != nil {\n              RespondWithError(w, ErrorInternal(\"internal service panic\"))\n              panic(rr)\n            }\n          }()\n          {{argsList .Outputs \"ret\"}}{{.Outputs | commaIfLen}} err = s.service.{{.Name}}(ctx{{.Inputs | commaIfLen}}{{argsList .Inputs \"reqContent.Arg\"}})\n        }()\n        {{- if .Outputs | len}}\n        respContent := struct {\n        {{- range $i, $output := .Outputs}}\n          Ret{{$i}} {{$output | methodArgType}} `json:\"{{$output.Name | downcaseName}}\"`\n        {{- end}}\n        }{ {{argsList .Outputs \"ret\"}} }\n        {{- end}}\n\n        if err != nil {\n          RespondWithError(w, err)\n          return\n        }\n\n        {{- if .Outputs | len}}\n        respBody, err := json.Marshal(respContent)\n        if err != nil {\n          RespondWithError(w, Errorf(ErrInternal, err, \"failed to marshal json response\"))\n          return\n        }\n        {{- end}}\n\n        w.Header().Set(\"Content-Type\", \"application/json\")\n        w.WriteHeader(http.StatusOK)\n\n        {{- if .Outputs | len}}\n        w.Write(respBody)\n        {{- end}}\n      {{- end}}\n    }\n    {{- if .StreamOutput}}\n      type {{.Name | streamWriterName}} struct {\n        *serverStreamWriter\n      }\n\n      func (s *{{.Name | streamWriterName}}) Data({{.Outputs | methodOutputsWithTypes}}) error {\n        {{- range $i, $output := .Outputs}}\n          ret{{$i}} := {{$output.Name | downcaseName}}\n        {{- end}}\n\n        type data struct {\n          {{- range $i, $output := .Outputs}}\n          Ret{{$i}} {{$output | methodArgType}} `json:\"{{$output.Name}}\"`\n          {{- end}}\n        }\n\n        body := struct {\n          Data data `json:\"data\"`\n        }{data{ {{argsList .Outputs \"ret\"}}{{commaIfLen .Outputs}}}}\n\n        payload, err := json.Marshal(body.Data)\n        if err != nil {\n          werr := Errorf(ErrStreamLost, err, \"failed to marshal json response\")\n          s.Error(werr)\n          s.Close()\n          return werr\n        }\n\n        return s.Write(payload)\n      }\n    {{- end}}\n  {{end}}\n{{- end}}\n\n{{- if .Services | hasStreamOutput }}\n\n//\n// Server streaming helpers\n//\n\nconst StreamKeepAliveInterval = 30 * time.Second\n\ntype streamWriter interface {\n  Write(payload []byte) error\n  Error(err error) error\n  Ping() error\n  Close() error\n  Done() <-chan struct{}\n}\n\ntype serverStreamWriter struct {\n  w             http.ResponseWriter\n  flusher       http.Flusher\n  headerWritten bool\n  done          chan struct{}\n  mu            sync.Mutex\n}\n\nfunc newServerStreamWriter(w http.ResponseWriter) (*serverStreamWriter, error) {\n  flusher, ok := w.(http.Flusher)\n  if !ok {\n    return nil, errors.New(\"expected http.ResponseWriter to be an http.Flusher\")\n  }\n  return &serverStreamWriter{w: w, flusher: flusher}, nil\n}\n\nfunc (s *serverStreamWriter) Write(payload []byte) error {\n  select {\n  case <-s.Done():\n    return ErrStreamClosed\n  default:\n  }\n\n  s.mu.Lock()\n  defer s.mu.Unlock()\n\n  w := s.w\n  if !s.headerWritten {\n    // content-type is very improve here as proxy servers treat it differently\n    // w.Header().Set(\"Content-Type\", \"application/stream+json\")\n    w.Header().Set(\"Content-Type\", \"application/json\") // TODO: just for testing purposes..\n\n    w.Header().Set(\"Transfer-Encoding\", \"chunked\")\n    w.Header().Set(\"Connection\", \"keep-alive\")\n    w.Header().Set(\"Cache-Control\", \"no-cache\")\n    s.headerWritten = true\n  }\n\n  s.w.Write([]byte(fmt.Sprintf(\"%x\\r\\n\", len(payload))))\n  s.w.Write(payload)\n  s.w.Write([]byte(\"\\r\\n\"))\n  s.flusher.Flush()\n  return nil\n}\n\nfunc (s *serverStreamWriter) Error(err error) error {\n  e, ok := err.(Error)\n  if !ok {\n    e = Errorf(ErrInternal, err, err.Error())\n  }\n\n  body := struct {\n    Error Error `json:\"error\"`\n  }{e}\n\n  payload, err := json.Marshal(body)\n  if err != nil {\n    werr := Errorf(ErrStreamLost, err, \"failed to marshal json response\")\n    s.Close()\n    return werr\n  }\n\n  return s.Write(payload)\n}\n\nfunc (s *serverStreamWriter) Ping() error {\n  return s.Write([]byte(`{\"ping\":true}`))\n}\n\nfunc (s *serverStreamWriter) Close() error {\n  select {\n  case <-s.Done():\n    return nil\n  default:\n  }\n\n  s.mu.Lock()\n  fmt.Fprintf(s.w, \"0\\r\\n\")\n  s.flusher.Flush()\n  close(s.done)\n  s.mu.Unlock()\n  return nil\n}\n\nfunc (s *serverStreamWriter) Done() <-chan struct{} {\n  s.mu.Lock()\n  if s.done == nil {\n    s.done = make(chan struct{})\n  }\n  d := s.done\n  s.mu.Unlock()\n  return d\n}\n{{- end}}\n\n//\n// Server helpers\n//\n\nfunc RespondWithError(w http.ResponseWriter, err error) {\n  e, ok := err.(Error)\n  if !ok {\n    e = Errorf(ErrInternal, err, err.Error())\n  }\n  w.Header().Set(\"Content-Type\", \"application/json\")\n  w.WriteHeader(HTTPStatusFromError(err))\n  respBody, _ := json.Marshal(e)\n  w.Write(respBody)\n}\n\n{{end}}\n{{end}}\nPK\x07\x08\xb2\xca&\x00$\"\x00\x00$\"\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00\x98\x88-T\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0d\x00	\x00types.go.tmplUT\x05\x00\x01\xb1[\xe0a{{define \"types\"}}\n\n{{if .Messages}}\n//\n// Types\n//\n\n{{range .Messages}}\n  {{if .Type | isEnum}}\n    {{$enumName := .Name}}\n    {{$enumType := .EnumType}}\n    type {{$enumName}} {{$enumType}}\n\n    var (\n      {{- range .Fields}}\n        {{$enumName}}_{{.Name}} {{$enumName}} = {{.Value}}\n      {{- end}}\n    )\n\n    var {{$enumName}}_name = map[{{$enumType}}]string {\n      {{- range .Fields}}\n        {{.Value}}: \"{{.Name}}\",\n      {{- end}}\n    }\n\n    var {{$enumName}}_value = map[string]{{$enumType}} {\n      {{- range .Fields}}\n        \"{{.Name}}\": {{.Value}},\n      {{- end}}\n    }\n\n    func (x {{$enumName}}) String() string {\n      return {{$enumName}}_name[{{$enumType}}(x)]\n    }\n\n    func (x {{$enumName}}) MarshalJSON() ([]byte, error) {\n      buf := bytes.NewBufferString(`\"`)\n      buf.WriteString({{$enumName}}_name[{{$enumType}}(x)])\n      buf.WriteString(`\"`)\n      return buf.Bytes(), nil\n    }\n\n    func (x *{{$enumName}}) UnmarshalJSON(b []byte) error {\n      var j string\n      err := json.Unmarshal(b, &j)\n      if err != nil {\n        return err\n      }\n      *x = {{$enumName}}({{$enumName}}_value[j])\n      return nil\n    }\n\n    func (x *{{$enumName}}) UnmarshalText(b []byte) error {\n      enum := string(b)\n      *x = {{$enumName}}({{$enumName}}_value[enum])\n      return nil\n    }\n\n    func (x {{$enumName}}) MarshalText() ([]byte, error) {\n      return []byte({{$enumName}}_name[{{$enumType}}(x)]), nil\n    }\n  {{end}}\n  {{if .Type | isStruct  }}\n    type {{.Name}} struct {\n      {{- range .Fields}}\n        {{. | exportedField}} {{. | fieldOptional}}{{. | fieldTypeDef}} {{. | fieldTags}}\n      {{- end}}\n    }\n  {{end}}\n{{end}}\n{{end}}\n{{if .Services}}\n  {{range .Services}}\n    type {{.Name}} interface {\n      {{- range .Methods}}\n        {{- if not .StreamOutput }}\n          {{.Name}}({{.Inputs | methodInputs}}) ({{.Outputs | methodOutputs}})\n        {{- end}}\n      {{- end}}\n    }\n    {{- range .Methods}}\n      {{- if .StreamOutput }}\n\n          type {{.Name}}StreamWriter interface {\n            streamWriter\n            Data({{.Outputs | methodOutputsWithTypes}}) error\n          }\n\n          type {{.Name}}StreamReader interface {\n            Read() ({{.Outputs | methodOutputsWithTypes}}, err error)\n          }\n      {{- end}}\n    {{- end}}\n  {{end}}\n\n  var WebRPCServices = map[string][]string{\n    {{- range .Services}}\n      \"{{.Name}}\": {\n        {{- range .Methods}}\n          \"{{.Name}}\",\n        {{- end}}\n      },\n    {{- end}}\n  }\n{{end}}\n\n{{end}}\nPK\x07\x08Mr\xd7\x8f\xc5	\x00\x00\xc5	\x00\x00PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00W\x18?T\xb1{\x92\x13H\x1f\x00\x00H\x1f\x00\x00\x0e\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x00\x00\x00\x00client.go.tmplUT\x05\x00\x01WQ\xf7aPK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x98\x88-T\xd8\xdb\x18\xad\x8a+\x00\x00\x8a+\x00\x00\x0f\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\x8d\x1f\x00\x00helpers.go.tmplUT\x05\x00\x01\xb1[\xe0aPK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x98\x88-T\x0eD\x8eO\x1b\x03\x00\x00\x1b\x03\x00\x00\x11\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81]K\x00\x00proto.gen.go.tmplUT\x05\x00\x01\xb1[\xe0aPK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x98\x88-T\xb2\xca&\x00$\"\x00\x00$\"\x00\x00\x0e\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81\xc0N\x00\x00server.go.tmplUT\x05\x00\x01\xb1[\xe0aPK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00\x98\x88-TMr\xd7\x8f\xc5	\x00\x00\xc5	\x00\x00\x0d\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa4\x81)q\x00\x00types.go.tmplUT\x05\x00\x01\xb1[\xe0aPK\x05\x06\x00\x00\x00\x00\x05\x00\x05\x00\\\x01\x00\x002{\x00\x00\x00\x00"
