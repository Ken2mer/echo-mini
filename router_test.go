package echo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	handlerFunc = func(c Context) error {
		c.Set("path", c.Path())
		return nil
	}
)

func TestRouterStatic(t *testing.T) {
	e := New()
	r := e.router
	path := "/folders/a/files/echo.gif"
	r.Add(http.MethodGet, path, handlerFunc)
	c := e.NewContext(nil, nil).(*context)

	r.Find(http.MethodGet, path, c)
	c.handler(c)

	assert.Equal(t, path, c.Get("path"))
}
