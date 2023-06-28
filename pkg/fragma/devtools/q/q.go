// Package q wraps github.com/ryboe/q for convenient usage during development
// without having to add/remove the module.
package q

import "github.com/ryboe/q"

func init() {
	q.CallDepth = 3
}

func Q(v ...any) {
	q.Q(v...)
}
