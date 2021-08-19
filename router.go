package echo

import (
	"net/http"
)

type (
	Router struct {
		tree   *node
		routes map[string]*Route
		echo   *Echo
	}

	node struct {
		kind kind
		// label  byte
		prefix string
		// parent *node
		// staticChildren children
		// ppath         string
		// pnames        []string
		methodHandler *methodHandler
		// paramChild     *node
		// anyChild       *node

		// isLeaf indicates that node does not have child routes
		// isLeaf bool

		// isHandler indicates that node has at least one handler registered to it
		isHandler bool
	}

	kind          uint8
	// children      []*node
	methodHandler struct {
		get HandlerFunc
	}
)

const (
	staticKind kind = iota
)

func NewRouter(e *Echo) *Router {
	return &Router{
		tree: &node{
			methodHandler: new(methodHandler),
		},
		routes: map[string]*Route{},
		echo:   e,
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
	ppath := path        // Pristine path

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

	r.insert(method, path, h, staticKind, ppath, pnames)
}

func (r *Router) insert(method, path string, h HandlerFunc, t kind, ppath string, pnames []string) {

	currentNode := r.tree // Current node as root
	// if currentNode == nil {
	// 	panic("echo: invalid method")
	// }

	search := path

	for {
		// searchLen := len(search)
		// prefixLen := len(currentNode.prefix)
		lcpLen := 0

		// LCP - Longest Common Prefix (https://en.wikipedia.org/wiki/LCP_array)
		// max := prefixLen
		// if searchLen < max {
		// 	max = searchLen
		// }
		// for ; lcpLen < max && search[lcpLen] == currentNode.prefix[lcpLen]; lcpLen++ {
		// }

		if lcpLen == 0 {
			// At root node
			// currentNode.label = search[0]
			currentNode.prefix = search
			if h != nil {
				// currentNode.kind = t
				currentNode.addHandler(method, h)
				// currentNode.ppath = ppath
				// currentNode.pnames = pnames
			}
			// currentNode.isLeaf = currentNode.staticChildren == nil && currentNode.paramChild == nil && currentNode.anyChild == nil
		}
		// TODO: else ...
		return
	}
}

func (n *node) addHandler(method string, h HandlerFunc) {
	switch method {
	case http.MethodGet:
		n.methodHandler.get = h
	}

	if h != nil {
		n.isHandler = true
	} else {
		// n.isHandler = n.methodHandler.isHandler()
	}
}

func (n *node) findHandler(method string) HandlerFunc {
	switch method {
	case http.MethodGet:
		return n.methodHandler.get
	default:
		return nil
	}
}

func (n *node) checkMethodNotAllowed() HandlerFunc {
	for _, m := range methods {
		if h := n.findHandler(m); h != nil {
			return MethodNotAllowedHandler
		}
	}
	return NotFoundHandler
}

func (r *Router) Find(method, path string, c Context) {
	ctx := c.(*context)
	// ctx.path = path
	currentNode := r.tree // Current node as root

	var (
		// previousBestMatchNode *node
		matchedHandler HandlerFunc
		// search stores the remaining path to check for match. By each iteration we move from start of path to end of the path
		// and search value gets shorter and shorter.
		search      = path
		searchIndex = 0
		// paramIndex  int           // Param counter
		// paramValues = ctx.pvalues // Use the internal slice so the interface can keep the illusion of a dynamic slice
	)

	for {
		prefixLen := 0 // Prefix length
		lcpLen := 0    // LCP (longest common prefix) length

		if currentNode.kind == staticKind {
			searchLen := len(search)
			prefixLen = len(currentNode.prefix)

			// LCP - Longest Common Prefix (https://en.wikipedia.org/wiki/LCP_array)
			max := prefixLen
			if searchLen < max {
				max = searchLen
			}
			for ; lcpLen < max && search[lcpLen] == currentNode.prefix[lcpLen]; lcpLen++ {
			}
		}

		search = search[lcpLen:]
		searchIndex = searchIndex + lcpLen

		if search == "" && currentNode.isHandler {
			// if previousBestMatchNode == nil {
			// 	previousBestMatchNode = currentNode
			// }
			if h := currentNode.findHandler(method); h != nil {
				matchedHandler = h
				break
			}
		}

		// Not found
		break

	}

	// if currentNode == nil && previousBestMatchNode == nil {
	// 	return // nothing matched at all
	// }

	if matchedHandler != nil {
		ctx.handler = matchedHandler
	} else {
		// use previous match as basis. although we have no matching handler we have path match.
		// so we can send http.StatusMethodNotAllowed (405) instead of http.StatusNotFound (404)
		// currentNode = previousBestMatchNode

		ctx.handler = currentNode.checkMethodNotAllowed()
	}
	// ctx.path = currentNode.ppath
	// ctx.pnames = currentNode.pnames

	return
}
