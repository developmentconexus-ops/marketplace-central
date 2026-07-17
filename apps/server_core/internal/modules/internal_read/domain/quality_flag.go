package domain

type QualityFlag string

const (
	QualityComplete         QualityFlag = "complete"
	QualityMissingProduct   QualityFlag = "missing_product"
	QualityAmbiguousProduct QualityFlag = "ambiguous_product"
	QualityMissingStock     QualityFlag = "missing_stock"
	QualityMissingPrice     QualityFlag = "missing_price"
	QualityMissingCost      QualityFlag = "missing_cost"
	QualityAmbiguousPrice   QualityFlag = "ambiguous_price"
	QualityMissingTax       QualityFlag = "missing_tax"
	QualityStaleSource      QualityFlag = "stale_source"
	QualityInvalidEAN       QualityFlag = "invalid_ean"
	QualityEANCollision     QualityFlag = "ean_collision"
)

func HasQualityFlag(flags []QualityFlag, target QualityFlag) bool {
	for _, flag := range flags {
		if flag == target {
			return true
		}
	}
	return false
}
