// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package openapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// RequestURL defines model for RequestURL.
type RequestURL struct {
	OriginalURL *string `json:"originalURL,omitempty"`
}

// ResponseURL defines model for ResponseURL.
type ResponseURL struct {
	ShortURL *string `json:"shortURL,omitempty"`
}

// Stats defines model for Stats.
type Stats struct {
	OriginalURL *string `json:"originalURL,omitempty"`
	ShortURL    *string `json:"shortURL,omitempty"`
	Statistics  *int64  `json:"statistics,omitempty"`
}

// CreateshortURLJSONBody defines parameters for CreateshortURL.
type CreateshortURLJSONBody = RequestURL

// CreateshortURLJSONRequestBody defines body for CreateshortURL for application/json ContentType.
type CreateshortURLJSONRequestBody = CreateshortURLJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create short URL
	// (POST /)
	CreateshortURL(w http.ResponseWriter, r *http.Request)
	// Get url statistics
	// (GET /stats/{short-url})
	GetStats(w http.ResponseWriter, r *http.Request, shortUrl string)
	// Redirect to original URL by short URL
	// (GET /{short-url})
	RedirectURL(w http.ResponseWriter, r *http.Request, shortUrl string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// CreateshortURL operation middleware
func (siw *ServerInterfaceWrapper) CreateshortURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateshortURL(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// GetStats operation middleware
func (siw *ServerInterfaceWrapper) GetStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "short-url" -------------
	var shortUrl string

	err = runtime.BindStyledParameter("simple", false, "short-url", chi.URLParam(r, "short-url"), &shortUrl)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "short-url", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetStats(w, r, shortUrl)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// RedirectURL operation middleware
func (siw *ServerInterfaceWrapper) RedirectURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "short-url" -------------
	var shortUrl string

	err = runtime.BindStyledParameter("simple", false, "short-url", chi.URLParam(r, "short-url"), &shortUrl)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "short-url", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.RedirectURL(w, r, shortUrl)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/", wrapper.CreateshortURL)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/stats/{short-url}", wrapper.GetStats)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/{short-url}", wrapper.RedirectURL)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8RUzW7bPBB8FWK/76hYap32oFvbQxCgQAsbPhU+0NJKZiCR7HIVwDD07gVJSZYdt01/",
	"gt6k5f4MZ2Z5hMK01mjU7CA/giv22MrwucKvHTrerD76P0vGIrHCcGZI1UrLZjisDLWSIYeOGkiADxYh",
	"B8ekdA19P0XM7gELhj6BFTprtMOr3d3eEP926zXLeJc/gJz8AoYEHEtWjlXhzrKV5re3p3ylGWuka6B9",
	"SOnK+PISXUHKsjIacvhkUb/7fC821IgACTWS76m48R021Kxn4UckFwtfLbJF5sEZi1paBTksF9liCQlY",
	"yfuANA0sGceBH4sk/dT7EnL4QCgZJxISoOiG96Y8+OzCaEYdCqW1jSpCafrg/PDRRv7rf8IKcvgvPfks",
	"HUyWzhwWKPAzFGEJOVOHIRBdEtC+zrK/OPlkvzD6nPVwb9FRI4rAQ+mJvI3zzzN3shQDNT7nzbWce81I",
	"WjbCIT0iCSQy0Qaua1tJh4nvKLGIjLOsHeRfYD2KsPUlqTebS48h86ajpvcDa7yi4R1y3AUvOckWGcl3",
	"vMTn71kZEjUyK10LNxQpf+jNAglo2eJIzE1cg3O1khnzl0u6fUEl4w2vadgVBTpXdY2YaIky3j6VSBsW",
	"lel0+V0R1TNEvMNomtl7MJMxAI0aPke9FZaKsBjW74cCrkfXBBkJy84XvqSCy2z5lKJ/wPhIkmAjxgc+",
	"ELE7/HSX+il8OTksY9iEWYuBv6lDn1yW3c32J8ov5M50lw2iYbf9twAAAP//oXTP8nsHAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
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
	var res = make(map[string]func() ([]byte, error))
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
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
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
