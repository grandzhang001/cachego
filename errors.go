package cachego

import "fmt"

type err string

// Error returns the string error value.
func (e err) Error() string {
	return string(e)
}

const (
	// ErrCacheExpired returns an error when the cache key was expired.
	ErrCacheExpired = err("cache expired")

	// ErrDelete return an error when deletion fails.
	ErrDelete = err("unable to delete")

	// ErrFlush returns an error when flush fails.
	ErrFlush = err("unable to flush")

	// ErrSave returns an error when save fails.
	ErrSave = err("unable to save")
)

func wrap(err, additionalErr error) error {
	return fmt.Errorf("%s: %w", additionalErr, err)
}
