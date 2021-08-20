package echo

import (
	stdContext "context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEcho(t *testing.T) {
	e := New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Router
	assert.NotNil(t, e.Router())

	// DefaultHTTPErrorHandler
	e.DefaultHTTPErrorHandler(errors.New("error"), c)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestEchoHandler(t *testing.T) {
	e := New()

	// HandlerFunc
	e.GET("/ok", func(c Context) error {
		return c.String(http.StatusOK, "OK")
	})

	c, b := request(http.MethodGet, "/ok", e)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, "OK", b)
}

func TestEchoGet(t *testing.T) {
	e := New()
	testMethod(t, http.MethodGet, "/", e)
}

func TestEchoRoutes(t *testing.T) {
	e := New()
	routes := []*Route{
		{http.MethodGet, "/users/:user/events", ""},
		{http.MethodGet, "/users/:user/events/public", ""},
		{http.MethodPost, "/repos/:owner/:repo/git/refs", ""},
		{http.MethodPost, "/repos/:owner/:repo/git/tags", ""},
	}
	for _, r := range routes {
		e.Add(r.Method, r.Path, func(c Context) error {
			return c.String(http.StatusOK, "OK")
		})
	}

	if assert.Equal(t, len(routes), len(e.Routes())) {
		for _, r := range e.Routes() {
			found := false
			for _, rr := range routes {
				if r.Method == rr.Method && r.Path == rr.Path {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Route %s %s not found", r.Method, r.Path)
			}
		}
	}
}

func TestEchoNotFound(t *testing.T) {
	e := New()
	req := httptest.NewRequest(http.MethodGet, "/files", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestEchoMethodNotAllowed(t *testing.T) {
	e := New()
	e.GET("/", func(c Context) error {
		return c.String(http.StatusOK, "Echo!")
	})
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
}

func TestEchoContext(t *testing.T) {
	e := New()
	c := e.pool.Get().(Context)
	assert.IsType(t, new(context), c)
	e.pool.Put(c)
}

func waitForServerStart(e *Echo, errChan <-chan error, isTLS bool) error {
	ctx, cancel := stdContext.WithTimeout(stdContext.Background(), 200*time.Millisecond)
	defer cancel()

	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			var addr net.Addr
			if isTLS {
				// addr = e.TLSListenerAddr()
			} else {
				addr = e.ListenerAddr()
			}
			if addr != nil && strings.Contains(addr.String(), ":") {
				return nil // was started
			}
		case err := <-errChan:
			if err == http.ErrServerClosed {
				return nil
			}
			return err
		}
	}
}

func TestEchoStart(t *testing.T) {
	e := New()
	errChan := make(chan error)

	go func() {
		err := e.Start(":0")
		if err != nil {
			errChan <- err
		}
	}()

	err := waitForServerStart(e, errChan, false)
	assert.NoError(t, err)

	assert.NoError(t, e.Close())
}

func testMethod(t *testing.T, method, path string, e *Echo) {
	p := reflect.ValueOf(path)
	h := reflect.ValueOf(func(c Context) error {
		return c.String(http.StatusOK, method)
	})
	i := interface{}(e)
	reflect.ValueOf(i).MethodByName(method).Call([]reflect.Value{p, h})
	_, body := request(method, path, e)
	assert.Equal(t, method, body)
}

func request(method, path string, e *Echo) (int, string) {
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

func TestHTTPError(t *testing.T) {
	t.Run("non-internal", func(t *testing.T) {
		err := NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"code": 12,
		})

		assert.Equal(t, "code=400, message=map[code:12]", err.Error())
	})
	// t.Run("internal", func(t *testing.T) {
	// 	err := NewHTTPError(http.StatusBadRequest, map[string]interface{}{
	// 		"code": 12,
	// 	})
	// 	err.SetInternal(errors.New("internal error"))
	// 	assert.Equal(t, "code=400, message=map[code:12], internal=internal error", err.Error())
	// })
}

func TestEchoClose(t *testing.T) {
	e := New()
	errCh := make(chan error)

	go func() {
		errCh <- e.Start(":0")
	}()

	err := waitForServerStart(e, errCh, false)
	assert.NoError(t, err)

	if err := e.Close(); err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, e.Close())

	err = <-errCh
	assert.Equal(t, err.Error(), "http: Server closed")
}

func TestEchoShutdown(t *testing.T) {
	e := New()
	errCh := make(chan error)

	go func() {
		errCh <- e.Start(":0")
	}()

	err := waitForServerStart(e, errCh, false)
	assert.NoError(t, err)

	if err := e.Close(); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := stdContext.WithTimeout(stdContext.Background(), 10*time.Second)
	defer cancel()
	assert.NoError(t, e.Shutdown(ctx))

	err = <-errCh
	assert.Equal(t, err.Error(), "http: Server closed")
}
