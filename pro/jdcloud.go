package pro

import "github.com/suguer/EcsGateway/model"

type JdCloudPro struct {
	RegionArray map[string]string
	RegionMap   map[string]string
}

func NewJdCloudPro() *JdCloudPro {
	baidu := new(JdCloudPro)
	baidu.RegionArray = map[string]string{
		"cn-north-1": "华北-北京",
		"cn-east-1":  "华东-宿迁",
		"cn-east-2":  "华东-上海",
		"cn-south-1": "华南-广州",
	}
	baidu.RegionMap = map[string]string{
		"cn-north-1": "cn-beijing",
		"cn-south-1": "cn-guangzhou",
		"cn-east-2":  "cn-shanghai",
	}
	return baidu
}

func (g *JdCloudPro) DescribeRegions() []model.Region {
	data := []model.Region{}
	for k, v := range g.RegionArray {
		data = append(data, model.Region{
			RegionId: k,
			Name:     v,
		})
	}
	return data
}

func (g *JdCloudPro) GetRegionMapValue(Region string) string {
	if _, ok := g.RegionMap[Region]; ok {
		return g.RegionMap[Region]
	}
	return Region
}

func (g *JdCloudPro) GetRegionMapKey(Region string) string {
	for k, v := range g.RegionMap {
		if v == Region {
			return k
		}
	}
	return Region
}

func (g *JdCloudPro) GetDiskCategory(category string) string {
	var DiskCategoryMap = map[string]string{
		// "cloud":            "hp1",
		// "cloud_efficiency": "hp1",
		"cloud_ssd":  "ssd.gp1",
		"cloud_essd": "ssd.io1",
		// "ephemeral_ssd":    "local-ssd",
		// "ephemeral":        "local",
	}

	for k, v := range DiskCategoryMap {
		if category == k {
			return v
		}
	}
	return category
}
