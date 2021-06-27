package echo

type (
	Router struct {
		tree   *node
		routes map[string]*Route
		echo   *Echo
	}

	node struct {
	}
)

func NewRouter(e *Echo) *Router {
	return &Router{

		routes: map[string]*Route{},
		echo:   e,
	}
}
