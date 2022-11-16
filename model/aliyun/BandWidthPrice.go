package aliyun

type BandwidthPrice struct {
	Currency    string `json:"currency"`
	Description string `json:"description"`
	PricingInfo map[string]BandwidthPriceInfo
}

type BandwidthPriceInfo struct {
	Hours  []BandwidthPriceHours  `json:"hours"`
	Months []BandwidthPriceMonths `json:"months"`
}
type BandwidthPriceHours struct {
	Price  string `json:"price"`
	Period string `json:"period"`
}
type BandwidthPriceMonths struct {
	Value  string `json:"value"`
	Price  string `json:"price"`
	Period string `json:"period"`
	Unit   string `json:"unit"`
}
