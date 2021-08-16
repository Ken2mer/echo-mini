package echo

import (
	"fmt"
	"net"
	"net/http"
	"sync"
)

type (
	Echo struct {
		// common

		// startupMutex sync.RWMutex
		// StdLogger        *stdLog.Logger
		// colorer          *color.Color
		// premiddleware    []MiddlewareFunc
		// middleware       []MiddlewareFunc
		// maxParam *int
		router  *Router
		// routers map[string]*Router
		// notFoundHandler  HandlerFunc
		pool sync.Pool
		// Server *http.Server
		// TLSServer   *http.Server
		Listener    net.Listener
		TLSListener net.Listener
		// AutoTLSManager   autocert.Manager
		DisableHTTP2 bool
		Debug        bool
		// HideBanner   bool
		// HidePort bool
		// HTTPErrorHandler HTTPErrorHandler
		// Binder           Binder
		// Validator        Validator
		// Renderer         Renderer
		Logger Logger
		// IPExtractor      IPExtractor
		// ListenerNetwork string
	}

	Route struct {
		Method string `json:"method"`
		Path   string `json:"path"`
		Name   string `json:"name"`
	}

	HTTPError struct {
		Code     int         `json:"-"`
		Message  interface{} `json:"message"`
		Internal error       `json:"-"` // Stores the error returned by an external dependency
	}

	Map map[string]interface{}
)

// MIME types
const (
	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJSONCharsetUTF8       = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; " + charsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; " + charsetUTF8
	MIMETextXML                          = "text/xml"
	MIMETextXMLCharsetUTF8               = MIMETextXML + "; " + charsetUTF8
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEApplicationMsgpack               = "application/msgpack"
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; " + charsetUTF8
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"
)

const (
	charsetUTF8 = "charset=UTF-8"
	// PROPFIND Method can be used on collection and property resources.
	PROPFIND = "PROPFIND"
	// REPORT Method can be used to get information about a resource, see rfc 3253
	REPORT = "REPORT"
)

// Headers
const (
	HeaderAccept              = "Accept"
	HeaderAcceptEncoding      = "Accept-Encoding"
	HeaderAllow               = "Allow"
	HeaderAuthorization       = "Authorization"
	HeaderContentDisposition  = "Content-Disposition"
	HeaderContentEncoding     = "Content-Encoding"
	HeaderContentLength       = "Content-Length"
	HeaderContentType         = "Content-Type"
	HeaderCookie              = "Cookie"
	HeaderSetCookie           = "Set-Cookie"
	HeaderIfModifiedSince     = "If-Modified-Since"
	HeaderLastModified        = "Last-Modified"
	HeaderLocation            = "Location"
	HeaderUpgrade             = "Upgrade"
	HeaderVary                = "Vary"
	HeaderWWWAuthenticate     = "WWW-Authenticate"
	HeaderXForwardedFor       = "X-Forwarded-For"
	HeaderXForwardedProto     = "X-Forwarded-Proto"
	HeaderXForwardedProtocol  = "X-Forwarded-Protocol"
	HeaderXForwardedSsl       = "X-Forwarded-Ssl"
	HeaderXUrlScheme          = "X-Url-Scheme"
	HeaderXHTTPMethodOverride = "X-HTTP-Method-Override"
	HeaderXRealIP             = "X-Real-IP"
	HeaderXRequestID          = "X-Request-ID"
	HeaderXRequestedWith      = "X-Requested-With"
	HeaderServer              = "Server"
	HeaderOrigin              = "Origin"

	// Access control
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	// Security
	HeaderStrictTransportSecurity         = "Strict-Transport-Security"
	HeaderXContentTypeOptions             = "X-Content-Type-Options"
	HeaderXXSSProtection                  = "X-XSS-Protection"
	HeaderXFrameOptions                   = "X-Frame-Options"
	HeaderContentSecurityPolicy           = "Content-Security-Policy"
	HeaderContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
	HeaderXCSRFToken                      = "X-CSRF-Token"
	HeaderReferrerPolicy                  = "Referrer-Policy"
)

func New() (e *Echo) {
	e = &Echo{
		// Server: new(http.Server),
		// TLSServer: new(http.Server),
		// AutoTLSManager: autocert.Manager{
		// 	Prompt: autocert.AcceptTOS,
		// },
		// Logger:          log.New("echo"),
		// colorer:         color.New(),
		// maxParam:        new(int),
		// ListenerNetwork: "tcp",
	}
	// e.Server.Handler = e
	// e.TLSServer.Handler = e
	// e.HTTPErrorHandler = e.DefaultHTTPErrorHandler
	// e.Binder = &DefaultBinder{}
	// e.Logger.SetLevel(log.ERROR)
	// e.StdLogger = stdLog.New(e.Logger.Output(), e.Logger.Prefix()+": ", 0)
	e.pool.New = func() interface{} {
		return e.NewContext(nil, nil)
	}
	e.router = NewRouter(e)
	// e.routers = map[string]*Router{}
	return
}

func (e *Echo) NewContext(r *http.Request, w http.ResponseWriter) Context {
	return &context{
		request:  r,
		response: NewResponse(w, e),
		// store:    make(Map),
		echo: e,
		// pvalues:  make([]string, *e.maxParam),
		// handler: NotFoundHandler,
	}
}

func (e *Echo) Router() *Router {
	return e.router
}

// func (e *Echo) Routers() map[string]*Router {
// 	return e.routers
// }

func (e *Echo) DefaultHTTPErrorHandler(err error, c Context) {
	he, ok := err.(*HTTPError)
	if !ok {
		he = &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}

	err = c.JSON(he.Code, he.Message)
	if err != nil {
		e.Logger.Error(err)
	}
}

func (he *HTTPError) Error() string {
	if he.Internal == nil {
		return fmt.Sprintf("code=%d, message=%v", he.Code, he.Message)
	}
	return fmt.Sprintf("code=%d, message=%v, internal=%v", he.Code, he.Message, he.Internal)
}
