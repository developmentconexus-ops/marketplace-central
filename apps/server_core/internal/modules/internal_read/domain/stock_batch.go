package domain

// StockFact is the provider-neutral fact returned by a batch stock read.
// A nil Quantity is an explicit unknown value, never zero.
type StockFact struct {
	ProductID    int64
	Quantity     *float64
	Source       SourceMetadata
	QualityFlags []QualityFlag
}
