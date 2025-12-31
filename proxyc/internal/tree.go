// Package internal implements internal functionality for the proxy client.
package internal

import (
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/relychan/goutils"
	"github.com/relychan/relixy/schema"
)

type nodeTyp uint8

const (
	ntStatic   nodeTyp = iota // /home
	ntRegexp                  // /{id:[0-9]+}
	ntParam                   // /{user}
	ntCatchAll                // /api/v1/*
)

// Route holds parameter values from the request path.
type Route struct {
	Node        *Node
	Handler     schema.RelixyHandler
	ParamValues map[string]string
}

// Node presents the route tree to organize the recursive route structure.
// Inspired by [Chi router]
//
// [Chi router]: https://github.com/go-chi/chi/tree/v5.2.3
type Node struct {
	handlers map[string]schema.RelixyHandler

	// regexp matcher for regexp nodes
	rex *regexp.Regexp

	// pattern is the full pattern of the leaf
	pattern string

	// prefix is the common prefix we ignore
	prefix string

	// parameter name if the node is param.
	paramName string

	// child nodes should be stored in-order for iteration,
	// in groups of the node type.
	children [ntCatchAll + 1]nodes

	// first byte of the child prefix
	tail byte

	// node type: static, regexp, param, catchAll
	typ nodeTyp

	// first byte of the prefix
	label byte
}

func (n *Node) InsertRoute(
	pattern string,
	operations *schema.RelixyPathItem,
	options *InsertRouteOptions,
) (*Node, error) {
	var parent *Node

	search := pattern

	for {
		// Handle key exhaustion
		if len(search) == 0 {
			// Insert or update the node's leaf handler
			err := n.setEndpoint(pattern, operations, options)

			return n, err
		}

		// We're going to be searching for a wild node next,
		// in this case, we need to get the tail
		var (
			label     = search[0]
			paramName string
			segTail   byte
			segEndIdx int
			segTyp    nodeTyp
			segRexpat string
			err       error
		)

		if label == '{' || label == '*' {
			segTyp, paramName, segRexpat, segTail, _, segEndIdx, err = patNextSegment(search)
			if err != nil {
				return nil, err
			}
		}

		var prefix string

		if segTyp == ntRegexp {
			prefix = segRexpat
		}

		// Look for the edge to attach to
		parent = n
		n = n.getEdge(segTyp, label, segTail, prefix)

		// No edge, create one
		if n == nil {
			child := &Node{
				label:     label,
				tail:      segTail,
				prefix:    search,
				paramName: paramName,
			}

			hn, err := parent.addChild(child, search)
			if err != nil {
				return nil, err
			}

			err = hn.setEndpoint(pattern, operations, options)

			return hn, err
		}

		// Found an edge to match the pattern
		if n.typ > ntStatic {
			// We found a param node, trim the param from the search path and continue.
			// This param/wild pattern segment would already be on the tree from a previous
			// call to addChild when creating a new node.
			search = search[segEndIdx:]

			continue
		}

		// Static nodes fall below here.
		// Determine longest prefix of the search key on match.
		commonPrefix := longestPrefix(search, n.prefix)
		if commonPrefix == len(n.prefix) {
			// the common prefix is as long as the current node's prefix we're attempting to insert.
			// keep the search going.
			search = search[commonPrefix:]

			continue
		}

		// Split the node
		child := &Node{
			typ:    ntStatic,
			prefix: search[:commonPrefix],
		}

		err = parent.replaceChild(search[0], segTail, child)
		if err != nil {
			return nil, fmt.Errorf("failed to replace child node of pattern %s: %w", pattern, err)
		}

		// Restore the existing node
		n.label = n.prefix[commonPrefix]
		n.prefix = n.prefix[commonPrefix:]

		_, err = child.addChild(n, n.prefix)
		if err != nil {
			return nil, err
		}

		// If the new key is a subset, set the method/handler on this node and finish.
		search = search[commonPrefix:]
		if len(search) == 0 {
			err := child.setEndpoint(pattern, operations, options)

			return child, err
		}

		// Create a new edge for the node
		subchild := &Node{
			typ:    ntStatic,
			label:  search[0],
			prefix: search,
		}

		hn, err := child.addChild(subchild, search)
		if err != nil {
			return nil, err
		}

		err = hn.setEndpoint(pattern, operations, options)

		return hn, err
	}
}

func (n *Node) FindRoute(path string, method string) *Route {
	route := &Route{
		ParamValues: map[string]string{},
	}

	// Find the routing handlers for the path
	rn := n.findRouteRecursive(path, route)
	if rn == nil {
		return nil
	}

	handler, ok := rn.handlers[method]
	if !ok {
		return nil
	}

	route.Node = rn
	route.Handler = handler

	return route
}

// addChild appends the new `child` node to the tree using the `pattern` as the trie key.
// For a URL router like chi's, we split the static, param, regexp and wildcard segments
// into different nodes. In addition, addChild will recursively call itself until every
// pattern segment is added to the url pattern tree as individual nodes, depending on type.
func (n *Node) addChild(child *Node, prefix string) (*Node, error) {
	search := prefix

	// handler leaf node added to the tree is the child.
	// this may be overridden later down the flow
	hn := child

	// Parse next segment
	segTyp, paramName, segRexpat, segTail, segStartIdx, segEndIdx, err := patNextSegment(search)
	if err != nil {
		return nil, fmt.Errorf("failed to add child node %s: %w", child.pattern, err)
	}

	// Add child depending on next up segment
	switch segTyp {
	case ntStatic:
		// Search prefix is all static (that is, has no params in path)
		// noop

	default:
		// Search prefix contains a param, regexp or wildcard
		if segTyp == ntRegexp {
			rex, err := regexp.Compile(segRexpat)
			if err != nil {
				return nil, fmt.Errorf("%w: %s", ErrInvalidRegexpPatternParamInRoute, segRexpat)
			}

			child.prefix = segRexpat
			child.rex = rex
		}

		if segStartIdx == 0 {
			// Route starts with a param
			child.typ = segTyp

			if segTyp == ntCatchAll {
				segStartIdx = -1
			} else {
				segStartIdx = segEndIdx
			}

			if segStartIdx < 0 {
				segStartIdx = len(search)
			}

			child.tail = segTail // for params, we set the tail

			if segStartIdx != len(search) {
				// add static edge for the remaining part, split the end.
				// its not possible to have adjacent param nodes, so its certainly
				// going to be a static node next.
				search = search[segStartIdx:] // advance search position

				nn := &Node{
					typ:       ntStatic,
					label:     search[0],
					prefix:    search,
					paramName: paramName,
				}

				hn, err = child.addChild(nn, search)
				if err != nil {
					return nil, fmt.Errorf("failed to add child node %s: %w", search, err)
				}
			}
		} else if segStartIdx > 0 {
			// Route has some param

			// starts with a static segment
			child.typ = ntStatic
			child.prefix = search[:segStartIdx]
			child.rex = nil

			// add the param edge node
			search = search[segStartIdx:]

			nn := &Node{
				typ:       segTyp,
				label:     search[0],
				tail:      segTail,
				paramName: paramName,
			}

			hn, err = child.addChild(nn, search)
			if err != nil {
				return nil, fmt.Errorf("failed to add child node %s: %w", search, err)
			}
		}
	}

	n.children[child.typ] = append(n.children[child.typ], child)
	n.children[child.typ].Sort()

	return hn, nil
}

func (n *Node) replaceChild(label, tail byte, child *Node) error {
	for i := range len(n.children[child.typ]) {
		if n.children[child.typ][i].label == label && n.children[child.typ][i].tail == tail {
			n.children[child.typ][i] = child
			n.children[child.typ][i].label = label
			n.children[child.typ][i].tail = tail

			return nil
		}
	}

	return fmt.Errorf("%w: %s", ErrReplaceMissingChildNode, child.pattern)
}

func (n *Node) setEndpoint(
	pattern string,
	operations *schema.RelixyPathItem,
	options *InsertRouteOptions,
) error {
	paramKeys, err := patParamKeys(pattern)
	if err != nil {
		return fmt.Errorf("failed to extract param keys in the pattern %s: %w", pattern, err)
	}

	params := operations.Parameters

	params = schema.ExtractCommonParametersOfOperation(params, operations.Get)
	params = schema.ExtractCommonParametersOfOperation(params, operations.Post)
	params = schema.ExtractCommonParametersOfOperation(params, operations.Put)
	params = schema.ExtractCommonParametersOfOperation(params, operations.Patch)
	params = schema.ExtractCommonParametersOfOperation(params, operations.Delete)

	// validates and add unknown parameters from the request pattern
	for _, key := range paramKeys {
		if slices.ContainsFunc(params, func(param schema.Parameter) bool {
			return param.In == schema.InPath && param.Name == key
		}) {
			continue
		}

		params = append(params, schema.Parameter{
			Name:     key,
			In:       schema.InPath,
			Required: goutils.ToPtr(true),
			Schema: &schema.RelixySchema{
				Type: []schema.PrimitiveType{schema.String},
			},
		})
	}

	operations.Parameters = params
	n.pattern = pattern

	n.handlers = map[string]schema.RelixyHandler{}

	if operations.Get != nil {
		method := http.MethodGet

		handler, err := NewProxyHandler(operations.Get, &schema.NewRelixyHandlerOptions{
			Method:     method,
			Parameters: operations.Parameters,
			GetEnv:     options.GetEnv,
		})
		if err != nil {
			return fmt.Errorf("failed to create handler for GET %s: %w", pattern, err)
		}

		n.handlers[method] = handler
	}

	if operations.Post != nil {
		method := http.MethodPost

		handler, err := NewProxyHandler(operations.Post, &schema.NewRelixyHandlerOptions{
			Method:     method,
			Parameters: operations.Parameters,
			GetEnv:     options.GetEnv,
		})
		if err != nil {
			return fmt.Errorf("failed to create handler for POST %s: %w", pattern, err)
		}

		n.handlers[http.MethodPost] = handler
	}

	if operations.Put != nil {
		method := http.MethodPut

		handler, err := NewProxyHandler(operations.Put, &schema.NewRelixyHandlerOptions{
			Method:     method,
			Parameters: operations.Parameters,
			GetEnv:     options.GetEnv,
		})
		if err != nil {
			return fmt.Errorf("failed to create handler for PUT %s: %w", pattern, err)
		}

		n.handlers[method] = handler
	}

	if operations.Patch != nil {
		method := http.MethodPatch

		handler, err := NewProxyHandler(operations.Patch, &schema.NewRelixyHandlerOptions{
			Method:     method,
			Parameters: operations.Parameters,
			GetEnv:     options.GetEnv,
		})
		if err != nil {
			return fmt.Errorf("failed to create handler for PATCH %s: %w", pattern, err)
		}

		n.handlers[method] = handler
	}

	if operations.Delete != nil {
		method := http.MethodDelete

		handler, err := NewProxyHandler(operations.Delete, &schema.NewRelixyHandlerOptions{
			Method:     method,
			Parameters: operations.Parameters,
			GetEnv:     options.GetEnv,
		})
		if err != nil {
			return fmt.Errorf("failed to create handler for DELETE %s: %w", pattern, err)
		}

		n.handlers[method] = handler
	}

	return nil
}

func (n *Node) getEdge(ntyp nodeTyp, label, tail byte, prefix string) *Node {
	nds := n.children[ntyp]

	for i := range nds {
		if nds[i].label == label && nds[i].tail == tail {
			if ntyp == ntRegexp && nds[i].prefix != prefix {
				continue
			}

			return nds[i]
		}
	}

	return nil
}

// Recursive edge traversal by checking all nodeTyp groups along the way.
// It's like searching through a multi-dimensional radix trie.
func (n *Node) findRouteRecursive(path string, route *Route) *Node { //nolint:gocognit,cyclop
	nn := n
	search := path

	for t, nds := range nn.children {
		ntyp := nodeTyp(t) //nolint:gosec

		if len(nds) == 0 {
			continue
		}

		var (
			xn    *Node
			label byte
		)

		xsearch := search

		if search != "" {
			label = search[0]
		}

		switch ntyp {
		case ntStatic:
			xn = nds.findEdge(label)
			if xn == nil || !strings.HasPrefix(xsearch, xn.prefix) {
				continue
			}

			xsearch = xsearch[len(xn.prefix):]
		case ntParam, ntRegexp:
			// short-circuit and return no matching route for empty param values
			if xsearch == "" {
				continue
			}

			// serially loop through each node grouped by the tail delimiter
			for idx := range nds {
				xn = nds[idx]

				// label for param nodes is the delimiter byte
				p := strings.IndexByte(xsearch, xn.tail)

				if p < 0 {
					if xn.tail != '/' {
						continue
					}

					p = len(xsearch)
				} else if ntyp == ntRegexp && p == 0 {
					continue
				}

				if ntyp == ntRegexp && xn.rex != nil {
					if !xn.rex.MatchString(xsearch[:p]) {
						continue
					}
				} else if strings.IndexByte(xsearch[:p], '/') != -1 {
					// avoid a match across path segments
					continue
				}

				route.ParamValues[xn.paramName] = xsearch[:p]
				xsearch = xsearch[p:]

				if len(xsearch) == 0 {
					if xn.isLeaf() {
						return xn
					}
				}

				// recursively find the next node on this branch
				fin := xn.findRouteRecursive(xsearch, route)
				if fin != nil {
					return fin
				}

				// not found on this branch, reset vars
				xsearch = search
			}
		default:
			// catch-all nodes
			xn = nds[0]
			xsearch = ""
		}

		if xn == nil {
			continue
		}

		// did we find it yet?
		if len(xsearch) == 0 {
			if xn.isLeaf() {
				return xn
			}
		}

		// recursively find the next node..
		fin := xn.findRouteRecursive(xsearch, route)
		if fin != nil {
			return fin
		}
	}

	return nil
}

func (n *Node) isLeaf() bool {
	return n.handlers != nil
}

// patNextSegment returns the next segment details from a pattern:
// node type, param key, regexp string, param tail byte, param starting index, param ending index.
func patNextSegment( //nolint:revive
	pattern string,
) (nodeTyp, string, string, byte, int, int, error) { //nolint:revive
	ps := strings.Index(pattern, "{")
	ws := strings.Index(pattern, "*")

	if ps < 0 && ws < 0 {
		return ntStatic, "", "", 0, 0, len(pattern), nil // we return the entire thing
	}

	// Sanity check
	if ps >= 0 && ws >= 0 && ws < ps {
		return ntStatic, "", "", 0, 0, 0, ErrWildcardMustBeLast
	}

	var tail byte = '/' // Default endpoint tail to / byte

	if ps >= 0 { //nolint:nestif
		// Param/Regexp pattern is next
		nt := ntParam

		// Read to closing } taking into account opens and closes in curl count (cc)
		cc := 0
		pe := ps

		for i, c := range pattern[ps:] {
			if c == '{' {
				cc++
			} else if c == '}' {
				cc--
				if cc == 0 {
					pe = ps + i

					break
				}
			}
		}

		if pe == ps {
			return ntStatic, "", "", 0, 0, 0,
				fmt.Errorf("%s: %w", pattern, ErrMissingClosingBracket)
		}

		key := pattern[ps+1 : pe]
		pe++ // set end to next position

		if pe < len(pattern) {
			tail = pattern[pe]
		}

		key, rexpat, isRegexp := strings.Cut(key, ":")
		if isRegexp {
			nt = ntRegexp
		}

		if len(rexpat) > 0 {
			if rexpat[0] != '^' {
				rexpat = "^" + rexpat
			}

			if rexpat[len(rexpat)-1] != '$' {
				rexpat += "$"
			}
		}

		return nt, key, rexpat, tail, ps, pe, nil
	}

	// Wildcard pattern as finale
	if ws < len(pattern)-1 {
		return ntStatic, "", "", 0, 0, 0, ErrWildcardMustBeLast
	}

	return ntCatchAll, "*", "", 0, ws, len(pattern), nil
}

func patParamKeys(pattern string) ([]string, error) {
	pat := pattern
	paramKeys := []string{}

	for {
		ptyp, paramKey, _, _, _, e, err := patNextSegment(pat)
		if err != nil {
			return nil, err
		}

		if ptyp == ntStatic {
			return paramKeys, nil
		}

		if slices.Contains(paramKeys, paramKey) {
			return nil, fmt.Errorf(
				"%s: %w, '%s'",
				pattern,
				ErrDuplicatedParamKey,
				paramKey,
			)
		}

		paramKeys = append(paramKeys, paramKey)
		pat = pat[e:]
	}
}

// longestPrefix finds the length of the shared prefix
// of two strings.
func longestPrefix(k1, k2 string) int {
	maxLength := min(len(k1), len(k2))

	var i int

	for ; i < maxLength; i++ {
		if k1[i] != k2[i] {
			break
		}
	}

	return i
}

type nodes []*Node

// Sort the list of nodes by label.
func (ns nodes) Sort()              { sort.Sort(ns); ns.tailSort() }
func (ns nodes) Len() int           { return len(ns) }
func (ns nodes) Swap(i, j int)      { ns[i], ns[j] = ns[j], ns[i] }
func (ns nodes) Less(i, j int) bool { return ns[i].label < ns[j].label }

// tailSort pushes nodes with '/' as the tail to the end of the list for param nodes.
// The list order determines the traversal order.
func (ns nodes) tailSort() {
	for i := len(ns) - 1; i >= 0; i-- {
		if ns[i].typ > ntStatic && ns[i].tail == '/' {
			ns.Swap(i, len(ns)-1)

			return
		}
	}
}

func (ns nodes) findEdge(label byte) *Node {
	num := len(ns)
	idx := 0
	i, j := 0, num-1

	for i <= j {
		idx = i + (j-i)/2

		switch {
		case label > ns[idx].label:
			i = idx + 1
		case label < ns[idx].label:
			j = idx - 1
		default:
			i = num // breaks cond
		}
	}

	if ns[idx].label != label {
		return nil
	}

	return ns[idx]
}
