package cache

import "time"

// Cacher is the interface that wraps the basic cache operations
type Cacher interface {
	Get([]byte) ([]byte, error)
	Set([]byte, []byte, time.Duration) error
	Has([]byte) bool
	Delete([]byte) error
}
