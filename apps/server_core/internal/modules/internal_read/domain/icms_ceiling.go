package domain

// UF is the numeric Sankhya state code used by TGFICM.UFORIG and UFDEST.
type UF int64

type ICMSCeiling struct {
	DestinationUF  UF
	CeilingPercent *float64
	Source         SourceMetadata
	QualityFlags   []QualityFlag
}
