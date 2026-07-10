package domain

type ProductCandidate struct {
	ProductID        int
	Name             string
	EAN              *string
	ReferenceCode    *string
	ProductGroupID   *int
	ProductGroupName *string
	BrandID          *int
	BrandName        *string
	IsActive         bool
	UsageType        *string
	Source           SourceMetadata
	QualityFlags     []QualityFlag
}
