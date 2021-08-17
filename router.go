package echo

type (
	Router struct {
		tree   *node
		routes map[string]*Route
		echo   *Echo
	}

	node struct {
	}

	kind          uint8
)

const (
	staticKind kind = iota
)

func NewRouter(e *Echo) *Router {
	return &Router{
		routes: map[string]*Route{},
		echo: e,
	}
}

func (r *Router) Add(method, path string, h HandlerFunc) {
	// Validate path
	// if path == "" {
	// 	path = "/"
	// }
	// if path[0] != '/' {
	// 	path = "/" + path
	// }

	pnames := []string{} // Param names
	// ppath := path        // Pristine path

	for i, lcpIndex := 0, len(path); i < lcpIndex; i++ {
		if path[i] == ':' {
			j := i + 1

			// r.insert(method, path[:i], nil, staticKind, "", nil)
			for ; i < lcpIndex && path[i] != '/'; i++ {
			}

			pnames = append(pnames, path[j:i])
			path = path[:j] + path[i:]
			i, lcpIndex = j, len(path)

			if i == lcpIndex {
				// path node is last fragment of route path. ie. `/users/:id`
				// r.insert(method, path[:i], h, paramKind, ppath, pnames)
			} else {
				// r.insert(method, path[:i], nil, paramKind, "", nil)
			}
		} else if path[i] == '*' {
			// r.insert(method, path[:i], nil, staticKind, "", nil)
			pnames = append(pnames, "*")
			// r.insert(method, path[:i+1], h, anyKind, ppath, pnames)
		}
	}

	// r.insert(method, path, h, staticKind, ppath, pnames)
}
