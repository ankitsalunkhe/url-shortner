// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// Request defines model for request.
type Request struct {
	// Url URL to shorten
	Url string `json:"url"`
}

// CreateUrl defines model for CreateUrl.
type CreateUrl struct {
	Url string `json:"url"`
}

// DeleteUrl defines model for DeleteUrl.
type DeleteUrl struct {
	Message string `json:"message"`
}

// Error defines model for Error.
type Error struct {
	Message string `json:"message"`
}

// Ping defines model for Ping.
type Ping struct {
	Message *string `json:"message,omitempty"`
}

// PostURLJSONRequestBody defines body for PostURL for application/json ContentType.
type PostURLJSONRequestBody = Request

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetPing request
	GetPing(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PostURLWithBody request with any body
	PostURLWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostURL(ctx context.Context, body PostURLJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteURL request
	DeleteURL(ctx context.Context, shortUrl string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetURL request
	GetURL(ctx context.Context, shortUrl string, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetPing(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetPingRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostURLWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostURLRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostURL(ctx context.Context, body PostURLJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostURLRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteURL(ctx context.Context, shortUrl string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteURLRequest(c.Server, shortUrl)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetURL(ctx context.Context, shortUrl string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetURLRequest(c.Server, shortUrl)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetPingRequest generates requests for GetPing
func NewGetPingRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/ping")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewPostURLRequest calls the generic PostURL builder with application/json body
func NewPostURLRequest(server string, body PostURLJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostURLRequestWithBody(server, "application/json", bodyReader)
}

// NewPostURLRequestWithBody generates requests for PostURL with any type of body
func NewPostURLRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/url")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewDeleteURLRequest generates requests for DeleteURL
func NewDeleteURLRequest(server string, shortUrl string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "shortUrl", runtime.ParamLocationPath, shortUrl)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/url/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetURLRequest generates requests for GetURL
func NewGetURLRequest(server string, shortUrl string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "shortUrl", runtime.ParamLocationPath, shortUrl)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/url/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetPingWithResponse request
	GetPingWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetPingResponse, error)

	// PostURLWithBodyWithResponse request with any body
	PostURLWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostURLResponse, error)

	PostURLWithResponse(ctx context.Context, body PostURLJSONRequestBody, reqEditors ...RequestEditorFn) (*PostURLResponse, error)

	// DeleteURLWithResponse request
	DeleteURLWithResponse(ctx context.Context, shortUrl string, reqEditors ...RequestEditorFn) (*DeleteURLResponse, error)

	// GetURLWithResponse request
	GetURLWithResponse(ctx context.Context, shortUrl string, reqEditors ...RequestEditorFn) (*GetURLResponse, error)
}

type GetPingResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Ping
}

// Status returns HTTPResponse.Status
func (r GetPingResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetPingResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PostURLResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *CreateUrl
	JSON400      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r PostURLResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostURLResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteURLResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *DeleteUrl
	JSON400      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r DeleteURLResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteURLResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetURLResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetURLResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetURLResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetPingWithResponse request returning *GetPingResponse
func (c *ClientWithResponses) GetPingWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetPingResponse, error) {
	rsp, err := c.GetPing(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetPingResponse(rsp)
}

// PostURLWithBodyWithResponse request with arbitrary body returning *PostURLResponse
func (c *ClientWithResponses) PostURLWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostURLResponse, error) {
	rsp, err := c.PostURLWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostURLResponse(rsp)
}

func (c *ClientWithResponses) PostURLWithResponse(ctx context.Context, body PostURLJSONRequestBody, reqEditors ...RequestEditorFn) (*PostURLResponse, error) {
	rsp, err := c.PostURL(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostURLResponse(rsp)
}

// DeleteURLWithResponse request returning *DeleteURLResponse
func (c *ClientWithResponses) DeleteURLWithResponse(ctx context.Context, shortUrl string, reqEditors ...RequestEditorFn) (*DeleteURLResponse, error) {
	rsp, err := c.DeleteURL(ctx, shortUrl, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteURLResponse(rsp)
}

// GetURLWithResponse request returning *GetURLResponse
func (c *ClientWithResponses) GetURLWithResponse(ctx context.Context, shortUrl string, reqEditors ...RequestEditorFn) (*GetURLResponse, error) {
	rsp, err := c.GetURL(ctx, shortUrl, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetURLResponse(rsp)
}

// ParseGetPingResponse parses an HTTP response from a GetPingWithResponse call
func ParseGetPingResponse(rsp *http.Response) (*GetPingResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetPingResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Ping
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParsePostURLResponse parses an HTTP response from a PostURLWithResponse call
func ParsePostURLResponse(rsp *http.Response) (*PostURLResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostURLResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest CreateUrl
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseDeleteURLResponse parses an HTTP response from a DeleteURLWithResponse call
func ParseDeleteURLResponse(rsp *http.Response) (*DeleteURLResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteURLResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest DeleteUrl
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseGetURLResponse parses an HTTP response from a GetURLWithResponse call
func ParseGetURLResponse(rsp *http.Response) (*GetURLResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetURLResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Ping
	// (GET /ping)
	GetPing(ctx echo.Context) error
	// Create short url
	// (POST /url)
	PostURL(ctx echo.Context) error
	// Delete Long URL
	// (DELETE /url/{shortUrl})
	DeleteURL(ctx echo.Context, shortUrl string) error
	// Get Long URL
	// (GET /url/{shortUrl})
	GetURL(ctx echo.Context, shortUrl string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetPing converts echo context to params.
func (w *ServerInterfaceWrapper) GetPing(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPing(ctx)
	return err
}

// PostURL converts echo context to params.
func (w *ServerInterfaceWrapper) PostURL(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostURL(ctx)
	return err
}

// DeleteURL converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteURL(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "shortUrl" -------------
	var shortUrl string

	err = runtime.BindStyledParameterWithOptions("simple", "shortUrl", ctx.Param("shortUrl"), &shortUrl, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter shortUrl: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteURL(ctx, shortUrl)
	return err
}

// GetURL converts echo context to params.
func (w *ServerInterfaceWrapper) GetURL(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "shortUrl" -------------
	var shortUrl string

	err = runtime.BindStyledParameterWithOptions("simple", "shortUrl", ctx.Param("shortUrl"), &shortUrl, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter shortUrl: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetURL(ctx, shortUrl)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/ping", wrapper.GetPing)
	router.POST(baseURL+"/url", wrapper.PostURL)
	router.DELETE(baseURL+"/url/:shortUrl", wrapper.DeleteURL)
	router.GET(baseURL+"/url/:shortUrl", wrapper.GetURL)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xVzU7cMBB+lWjao0UWaC+5UYoQEocVlBPiYJLZjSGxzXhSFa387tU4CZtlaSkF0Z5I",
	"nGH8zfczu4LStd5ZtBygWAFh8M4GTC+HhJrxghp5KZ1ltCyP2vvGlJqNs/lNcFbOQlljq+XJk/NIbPoe",
	"Xf/f+EO3vkEooGb2ocjzpeGd0rW5DtVid28fFPC9l4LAZOwSYlRAeNcZwgqKy9To6qHIXd9gyRClqsJQ",
	"kvECBwo478oSQ1h0TTYOA1HBV2zwtbO0GIJe4tvMMzZ7xUxHRI7eep7O6usGM3ZZmeTPQu2I0WKViQRv",
	"OVaPPyqYS6d/qMuzSAXghPqoBixDZu46DPxL72+2ujg7FXIHVkE9DbnurnvU9tZw0E1nb2t8RUaMXbht",
	"LAc2O5ifTLTWtsqWyBPNBe6CXJs1zi7lTUAYTnjl27lUWiRpBAq+I4W+9+7ObGcm4jqPVnsDBeynIwVe",
	"c50Iyv0g/BJ5G10i/chW3hnLkBpRMsRJBQUcIyfbqM2dtTebyZ+PhAso4EO+Xm/5Q10+H8kLXdtquh8u",
	"S0f5oJp3vaSbt85d4J6EQfUvrrp/kXWfQjaYKR+dFGOv7MZcu8/PtV7YUcGnP2HiIYKfX1C9QdzhZEuk",
	"DTGSmK/S2QU1sddWFvA2o8NiTpx6TbpFRgpQXD62w7d6vEV8V7pK8mDki/gJFFjdpmwMt8I0G0wdqokK",
	"j3N09Tc2Wv+mvB/d/Z3Z6ZjGqMbwbMXjf+N0v7fwJoIzrAxhybKEuMbMkVkaq5vpvqlRVwn/Ck5dH6zt",
	"RjLKsFtp3RN+BzC+n2zHyBPNYow/AwAA//935IaOeQkAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
