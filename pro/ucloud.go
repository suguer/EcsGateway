package pro

type UcloudPro struct {
	RegionMap map[string]string
}

func NewUcloudPro() *UcloudPro {
	pro := new(UcloudPro)
	pro.RegionMap = map[string]string{
		"hk":           "cn-hongkong",
		"cn-bj2":       "cn-beijing",
		"cn-gd":        "cn-guangzhou",
		"cn-sh2":       "cn-shanghai",
		"jpn-tky":      "ap-northeast-1",
		"kr-seoul":     "ap-northeast-2",
		"sg":           "ap-southeast-1",
		"ph-mnl":       "ap-southeast-6",
		"ge-fra":       "eu-central-1",
		"idn-jakarta":  "ap-southeast-5",
		"ind-mumbai":   "ap-south-1",
		"th-bkk":       "ap-southeast-7",
		"uae-dubai":    "me-east-1",
		"uk-london":    "eu-west-1",
		"rus-mosc":     "eu-moscow",
		"bra-saopaulo": "sa-saopaulo",
	}
	return pro
}

func (g *UcloudPro) GetRegionMapValue(Region string) string {
	if _, ok := g.RegionMap[Region]; ok {
		return g.RegionMap[Region]
	}
	return Region
}

func (g *UcloudPro) GetRegionMapKey(Region string) string {
	for k, v := range g.RegionMap {
		if v == Region {
			return k
		}
	}
	return Region
}
func (g *UcloudPro) GetDiskCategory(category string) string {
	var DiskCategoryMap = map[string]string{
		"cloud":            "CLOUD_BASIC",
		"cloud_efficiency": "CLOUD_PREMIUM",
		"cloud_ssd":        "CLOUD_SSD",
		"ephemeral_ssd":    "LOCAL_SSD",
		"ephemeral":        "LOCAL_BASIC",
	}
	for k, v := range DiskCategoryMap {
		if category == k {
			return v
		}
	}
	return category
}

func (g *UcloudPro) GetRegionMapChinese(Region string) string {
	RegionMapChinese := map[string]string{
		"us-ca":        "洛杉矶",
		"us-ws":        "华盛顿",
		"rus-mosc":     "莫斯科",
		"tw-tp":        "台北",
		"bra-saopaulo": "圣保罗",
		"afr-nigeria":  "拉各斯",
		"vn-sng":       "胡志明市",
		"cn-qz":        "福建",
		"cn-wlcb":      "华北二",
	}
	if _, ok := RegionMapChinese[Region]; ok {
		return RegionMapChinese[Region]
	}
	return Region
}
func (g *UcloudPro) GetOperatorArray(Region string) string {
	var OperatorArray = map[string]string{
		"cn-sh1":      "Bgp",
		"cn-sh2":      "Bgp",
		"cn-gd":       "Bgp",
		"cn-bj1":      "Bgp",
		"cn-bj2":      "Bgp",
		"cn-east-01":  "Duplet",
		"cn-east-02":  "Bgp",
		"cn-south-01": "Duplet",
		"cn-south-02": "Bgp",
		"cn-north-01": "Bgp",
		"cn-north-02": "Bgp",
		"cn-north-03": "Bgp",
		"cn-north-04": "Bgp",
	}
	if _, ok := OperatorArray[Region]; ok {
		return OperatorArray[Region]
	}
	return "International"
}
