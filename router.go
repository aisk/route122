package route122

import "strings"

// Router provides HTTP routing functionality
type Router struct {
	tree *routingNode
}

// Match represents the result of a successful route match
type Match struct {
	Handler any               // The handler that matches the request
	Params  map[string]string // Extracted path parameters
	Pattern string            // The matching pattern
}

// New creates a new Router instance
func New() *Router {
	return &Router{
		tree: &routingNode{},
	}
}

// Handle registers a new route pattern with its handler
func (r *Router) Handle(pattern string, handler any) error {
	if handler == nil {
		return &RouteError{
			Pattern: pattern,
			Message: "handler cannot be nil",
		}
	}

	p, err := parsePattern(pattern)
	if err != nil {
		return &RouteError{
			Pattern: pattern,
			Message: err.Error(),
		}
	}

	r.tree.addPattern(p, handler)
	return nil
}

// Match finds the handler that matches the given method, host, and path
func (r *Router) Match(method, host, path string) (Match, bool) {
	node, params := r.tree.match(host, method, path)
	if node == nil || node.handler == nil {
		return Match{}, false
	}

	return Match{
		Handler: node.handler,
		Params:  convertParams(params, node.pattern),
		Pattern: node.pattern.String(),
	}, true
}

func convertParams(wildcards []string, p *pattern) map[string]string {
	params := make(map[string]string)
	wildcardIndex := 0

	for _, seg := range p.segments {
		if seg.wild && seg.multi {
			remaining := wildcards[wildcardIndex:]
			params[seg.s] = strings.Join(remaining, "/")
			wildcardIndex += len(remaining)
			break
		} else if seg.wild {
			if wildcardIndex < len(wildcards) {
				params[seg.s] = wildcards[wildcardIndex]
				wildcardIndex++
			}
		}
	}

	return params
}

// RouteError represents an error that occurs during route registration
type RouteError struct {
	Pattern string
	Message string
}

func (e *RouteError) Error() string {
	return e.Message
}

