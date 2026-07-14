package ports

// CacheInvalidator evicts all cached reads for one fact class after a
// successful mutation. Implementations must not be called for failed writes.
type CacheInvalidator interface {
	InvalidateClass(factClass string)
}
