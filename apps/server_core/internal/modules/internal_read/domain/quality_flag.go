package domain

type QualityFlag string

const (
	QualityComplete         QualityFlag = "complete"
	QualityMissingProduct   QualityFlag = "missing_product"
	QualityMissingStock     QualityFlag = "missing_stock"
	QualityMissingCost      QualityFlag = "missing_cost"
	QualityMissingTax       QualityFlag = "missing_tax"
	QualityAmbiguousProduct QualityFlag = "ambiguous_product"
	QualityStaleSource      QualityFlag = "stale_source"
)
