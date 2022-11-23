package model

type PriceInfo struct {
	OriginalPrice float32
	TradePrice    float32
	PriceUnit     string
	PriceInfoItem []PriceInfoItem
}

type PriceInfoItem struct {
	ResourceType  string
	OriginalPrice float32
	TradePrice    float32
	Code          string
}
