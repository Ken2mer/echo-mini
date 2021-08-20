package echo

import (
	stdContext "context"
	"errors"
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
		pool   sync.Pool
		Server *http.Server
		// TLSServer   *http.Server
		Listener net.Listener
		// TLSListener net.Listener
		// AutoTLSManager   autocert.Manager
		DisableHTTP2 bool
		Debug        bool
		// HideBanner   bool
		// HidePort bool
		HTTPErrorHandler HTTPErrorHandler
		// Binder           Binder
		// Validator        Validator
		// Renderer         Renderer
		Logger Logger
		// IPExtractor      IPExtractor
		ListenerNetwork string
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

	MiddlewareFunc   func(HandlerFunc) HandlerFunc
	HandlerFunc      func(Context) error
	HTTPErrorHandler func(error, Context)

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

var (
	methods = [...]string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		PROPFIND,
		http.MethodPut,
		http.MethodTrace,
		REPORT,
	}
)

// Errors
var (
	ErrUnsupportedMediaType        = NewHTTPError(http.StatusUnsupportedMediaType)
	ErrNotFound                    = NewHTTPError(http.StatusNotFound)
	ErrUnauthorized                = NewHTTPError(http.StatusUnauthorized)
	ErrForbidden                   = NewHTTPError(http.StatusForbidden)
	ErrMethodNotAllowed            = NewHTTPError(http.StatusMethodNotAllowed)
	ErrStatusRequestEntityTooLarge = NewHTTPError(http.StatusRequestEntityTooLarge)
	ErrTooManyRequests             = NewHTTPError(http.StatusTooManyRequests)
	ErrBadRequest                  = NewHTTPError(http.StatusBadRequest)
	ErrBadGateway                  = NewHTTPError(http.StatusBadGateway)
	ErrInternalServerError         = NewHTTPError(http.StatusInternalServerError)
	ErrRequestTimeout              = NewHTTPError(http.StatusRequestTimeout)
	ErrServiceUnavailable          = NewHTTPError(http.StatusServiceUnavailable)
	ErrValidatorNotRegistered      = errors.New("validator not registered")
	ErrRendererNotRegistered       = errors.New("renderer not registered")
	ErrInvalidRedirectCode         = errors.New("invalid redirect status code")
	ErrCookieNotFound              = errors.New("cookie not found")
	ErrInvalidCertOrKeyType        = errors.New("invalid cert or key type, must be string or []byte")
	ErrInvalidListenerNetwork      = errors.New("invalid listener network")
)

// Error handlers
var (
	NotFoundHandler = func(c Context) error {
		return ErrNotFound
	}

	MethodNotAllowedHandler = func(c Context) error {
		return ErrMethodNotAllowed
	}
)

func New() (e *Echo) {
	e = &Echo{
		Server: new(http.Server),
		// TLSServer: new(http.Server),
		// AutoTLSManager: autocert.Manager{
		// 	Prompt: autocert.AcceptTOS,
		// },
		// Logger:          log.New("echo"),
		// colorer:         color.New(),
		// maxParam:        new(int),
		ListenerNetwork: "tcp",
	}
	// e.Server.Handler = e
	// e.TLSServer.Handler = e
	e.HTTPErrorHandler = e.DefaultHTTPErrorHandler
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

func (e *Echo) GET(path string, h HandlerFunc, m ...MiddlewareFunc) *Route {
	return e.Add(http.MethodGet, path, h, m...)
}

func (e *Echo) add(host, method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
	// name := handlerName(handler)
	router := e.findRouter(host)
	router.Add(method, path, func(c Context) error {
		h := applyMiddleware(handler, middleware...)
		return h(c)
	})
	r := &Route{
		Method: method,
		Path:   path,
		// Name:   name,
	}
	e.router.routes[method+path] = r
	return r
}

func (e *Echo) Add(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) *Route {
	return e.add("", method, path, handler, middleware...)
}

func (e *Echo) Routes() []*Route {
	routes := make([]*Route, 0, len(e.router.routes))
	for _, v := range e.router.routes {
		routes = append(routes, v)
	}
	return routes
}

func (e *Echo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Acquire context
	c := e.pool.Get().(*context)
	c.Reset(r, w)
	h := NotFoundHandler

	// if e.premiddleware == nil {
	// 	e.findRouter(r.Host).Find(r.Method, GetPath(r), c)
	// 	h = c.Handler()
	// 	h = applyMiddleware(h, e.middleware...)
	// } else {
	h = func(c Context) error {
		e.findRouter(r.Host).Find(r.Method, GetPath(r), c)
		h := c.Handler()
		// h = applyMiddleware(h, e.middleware...)
		return h(c)
	}
	// h = applyMiddleware(h, e.premiddleware...)
	// }

	// Execute chain
	if err := h(c); err != nil {
		e.HTTPErrorHandler(err, c)
	}

	// Release context
	e.pool.Put(c)
}

func (e *Echo) Start(address string) error {
	// e.startupMutex.Lock()
	e.Server.Addr = address
	if err := e.configureServer(e.Server); err != nil {
		// e.startupMutex.Unlock()
		return err
	}
	// e.startupMutex.Unlock()
	return e.Server.Serve(e.Listener)
}

func (e *Echo) configureServer(s *http.Server) (err error) {
	// Setup
	// e.colorer.SetOutput(e.Logger.Output())
	// s.ErrorLog = e.StdLogger
	s.Handler = e
	if e.Debug {
		// e.Logger.SetLevel(log.DEBUG)
	}

	// if !e.HideBanner {
	// 	e.colorer.Printf(banner, e.colorer.Red("v"+Version), e.colorer.Blue(website))
	// }

	if s.TLSConfig == nil {
		if e.Listener == nil {
			e.Listener, err = newListener(s.Addr, e.ListenerNetwork)
			if err != nil {
				return err
			}
		}
		// if !e.HidePort {
		// 	e.colorer.Printf("⇨ http server started on %s\n", e.colorer.Green(e.Listener.Addr()))
		// }
		return nil
	}
	// if e.TLSListener == nil {
	// 	l, err := newListener(s.Addr, e.ListenerNetwork)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	e.TLSListener = tls.NewListener(l, s.TLSConfig)
	// }
	// if !e.HidePort {
	// 	e.colorer.Printf("⇨ https server started on %s\n", e.colorer.Green(e.TLSListener.Addr()))
	// }
	return nil
}

func (e *Echo) ListenerAddr() net.Addr {
	// e.startupMutex.RLock()
	// defer e.startupMutex.RUnlock()
	if e.Listener == nil {
		return nil
	}
	return e.Listener.Addr()
}

func (e *Echo) Close() error {
	// e.startupMutex.Lock()
	// defer e.startupMutex.Unlock()
	// if err := e.TLSServer.Close(); err != nil {
	// 	return err
	// }
	return e.Server.Close()
}

func (e *Echo) Shutdown(ctx stdContext.Context) error {
	// e.startupMutex.Lock()
	// defer e.startupMutex.Unlock()
	// if err := e.TLSServer.Shutdown(ctx); err != nil {
	// 	return err
	// }
	return e.Server.Shutdown(ctx)
}

func NewHTTPError(code int, message ...interface{}) *HTTPError {
	he := &HTTPError{Code: code, Message: http.StatusText(code)}
	if len(message) > 0 {
		he.Message = message[0]
	}
	return he
}

func (he *HTTPError) Error() string {
	if he.Internal == nil {
		return fmt.Sprintf("code=%d, message=%v", he.Code, he.Message)
	}
	return fmt.Sprintf("code=%d, message=%v, internal=%v", he.Code, he.Message, he.Internal)
}

func GetPath(r *http.Request) string {
	path := r.URL.RawPath
	if path == "" {
		path = r.URL.Path
	}
	return path
}

func (e *Echo) findRouter(host string) *Router {
	// if len(e.routers) > 0 {
	// 	if r, ok := e.routers[host]; ok {
	// 		return r
	// 	}
	// }
	return e.router
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func newListener(address, network string) (*tcpKeepAliveListener, error) {
	if network != "tcp" && network != "tcp4" && network != "tcp6" {
		return nil, ErrInvalidListenerNetwork
	}
	l, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	return &tcpKeepAliveListener{l.(*net.TCPListener)}, nil
}

func applyMiddleware(h HandlerFunc, middleware ...MiddlewareFunc) HandlerFunc {
	// for i := len(middleware) - 1; i >= 0; i-- {
	// 	h = middleware[i](h)
	// }
	return h
}
