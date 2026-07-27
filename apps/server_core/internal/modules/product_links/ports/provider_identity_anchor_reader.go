package ports

type ProviderIdentityAnchor struct {
	Anchor   string
	Supplied bool
}

type ProviderIdentityAnchorReader interface {
	ProviderIdentityAnchors(providerCode string) ([]ProviderIdentityAnchor, error)
}
