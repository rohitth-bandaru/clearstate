package store

import "errors"

var ErrConflict = errors.New("conflict: resource was modified (optimistic lock)")
var ErrNotFound = errors.New("not found")
