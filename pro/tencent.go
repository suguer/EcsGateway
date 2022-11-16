package pro

type TencentPro struct {
	RegionMap map[string]string
}

func NewTencentPro() *TencentPro {
	baidu := new(TencentPro)
	baidu.RegionMap = map[string]string{
		"ap-guangzhou":     "cn-guangzhou",
		"ap-shanghai":      "cn-shanghai",
		"ap-nanjing":       "cn-nanjing",
		"ap-beijing":       "cn-beijing",
		"ap-chengdu":       "cn-chengdu",
		"ap-hongkong":      "cn-hongkong",
		"ap-seoul":         "ap-northeast-2",
		"ap-tokyo":         "ap-northeast-1",
		"ap-singapore":     "ap-southeast-1",
		"ap-bangkok":       "ap-southeast-7",
		"ap-jakarta":       "ap-southeast-5",
		"na-siliconvalley": "us-west-1",
		"eu-frankfurt":     "eu-central-1",
		"ap-mumbai":        "ap-south-1",
		"na-ashburn":       "us-east-1",
	}
	return baidu
}

func (g *TencentPro) GetRegionMapValue(Region string) string {
	if _, ok := g.RegionMap[Region]; ok {
		return g.RegionMap[Region]
	}
	return Region
}

func (g *TencentPro) GetRegionMapKey(Region string) string {
	for k, v := range g.RegionMap {
		if v == Region {
			return k
		}
	}
	return Region
}
func (g *TencentPro) GetDiskCategory(category string) string {
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
