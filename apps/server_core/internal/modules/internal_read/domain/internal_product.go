package domain

type ProductCandidate struct {
	// InternalProductID is populated only after a positive CODPROD has been
	// established. ProductID is legacy compatibility input and is not a
	// substitute for this canonical identity.
	InternalProductID *InternalProductID
	ProductID         int
	Name              string
	EAN               *string
	ReferenceCode     *string
	ProductGroupID    *int
	ProductGroupName  *string
	BrandID           *int
	BrandName         *string
	IsActive          bool
	UsageType         *string
	Source            SourceMetadata
	QualityFlags      []QualityFlag
}
