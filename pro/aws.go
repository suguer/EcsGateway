package pro

import "github.com/suguer/EcsGateway/model"

type AwsPro struct {
	RegionMap   map[string]string
	RegionArray map[string]string
}

func NewAwsPro() *AwsPro {
	pro := new(AwsPro)
	pro.RegionMap = map[string]string{
		"ap-east-1":      "cn-hongkong",
		"ap-southeast-3": "ap-southeast-5",
		"ap-south-1":     "ap-south-1",
		"ap-northeast-2": "ap-northeast-2",
		"ap-southeast-1": "ap-southeast-1",
		"ap-southeast-2": "ap-southeast-2",
		"ap-northeast-1": "ap-northeast-1",
		"us-east-1":      "us-east-1",
		"eu-central-1":   "eu-central-1",
		"eu-west-1":      "eu-west-5", //特殊处理免与当前地域重复冲突了
		"eu-west-2":      "eu-west-1",
		"sa-east-1":      "sa-saopaulo",
	}
	pro.RegionArray = map[string]string{
		"ap-east-1":      "亚太地区(香港)",
		"ap-southeast-3": "亚太地区(雅加达)",
		"ap-south-1":     "亚太地区(孟买)",
		"ap-northeast-3": "亚太地区(大阪)",
		"ap-northeast-2": "亚太地区(首尔)",
		"ap-southeast-1": "亚太地区(新加坡)",
		"ap-southeast-2": "亚太地区(悉尼)",
		"ap-northeast-1": "亚太地区(东京)",
		"us-east-1":      "美国东部(弗吉尼亚北部)",
		"us-east-2":      "美国东部(俄亥俄州)",
		"us-west-1":      "美国西部(加利福尼亚北部)",
		"us-west-2":      "美国西部(俄勒冈州)",
		"af-south-1":     "非洲(开普敦)",
		"ca-central-1":   "加拿大(中部)",
		"eu-central-1":   "欧洲(法兰克福)",
		"eu-west-1":      "欧洲(爱尔兰)",
		"eu-west-2":      "欧洲(伦敦)",
		"eu-south-1":     "欧洲(米兰)",
		"eu-west-3":      "欧洲(巴黎)",
		"eu-north-1":     "欧洲(斯德哥尔摩)",
		"me-south-1":     "中东(巴林)",
		"me-central-1":   "中东(阿联酋)",
		"sa-east-1":      "南美洲(圣保罗)",
	}
	return pro
}

func (g *AwsPro) GetRegionMapValue(Region string) string {
	if _, ok := g.RegionMap[Region]; ok {
		return g.RegionMap[Region]
	}
	return Region
}

func (g *AwsPro) GetRegionMapKey(Region string) string {
	for k, v := range g.RegionMap {
		if v == Region {
			return k
		}
	}
	return Region
}
func (g *AwsPro) GetDiskCategory(category string) string {
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

func (g *AwsPro) GetRegionMapChinese(Region string) string {
	if _, ok := g.RegionArray[Region]; ok {
		return g.RegionArray[Region]
	}
	return Region
}

func (g *AwsPro) DescribeRegions() []model.Region {
	data := []model.Region{}
	for k, v := range g.RegionArray {
		data = append(data, model.Region{
			RegionId: k,
			Name:     v,
		})
	}
	return data
}
