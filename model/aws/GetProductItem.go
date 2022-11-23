package aws

type Attributes struct {
	EnhancedNetworkingSupported string `json:"enhancedNetworkingSupported"`
	IntelTurboAvailable         string `json:"intelTurboAvailable"`
	Memory                      string `json:"memory"`
	DedicatedEbsThroughput      string `json:"dedicatedEbsThroughput"`
	Vcpu                        string `json:"vcpu"`
	Capacitystatus              string `json:"capacitystatus"`
	OperatingSystem             string `json:"operatingSystem"`
	RegionCode                  string `json:"regionCode"`
	InstanceType                string `json:"instanceType"`
}
type Terms struct {
	OnDemand map[string]TermItem
	Reserved map[string]TermItem
}
type TermItem struct {
	PriceDimensions map[string]PriceDimensions
}

type PriceDimensions struct {
	Unit         string `json:"unit"`
	PricePerUnit struct {
		USD string
	}
}
type GetProductItem struct {
	Product struct {
		ProductFamily string     `json:"productFamily"`
		Attributes    Attributes `json:"attributes"`
	}
	Terms Terms `json:"terms"`
}
