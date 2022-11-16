package pro

import "github.com/suguer/EcsGateway/model"

type BaiduPro struct {
	RegionArray map[string]string
	RegionMap   map[string]string
}

func NewBaiduPro() *BaiduPro {
	baidu := new(BaiduPro)
	baidu.RegionArray = map[string]string{
		"bj":  "北京",
		"gz":  "广州",
		"su":  "苏州",
		"hkg": "香港",
		"fwh": "武汉",
		"bd":  "保定",
		"sin": "新加坡",
		"fsh": "上海",
	}
	baidu.RegionMap = map[string]string{
		"hkg": "cn-hongkong",
		"bj":  "cn-beijing",
		"gz":  "cn-guangzhou",
		"sin": "ap-southeast-1",
		"fsh": "cn-shanghai",
	}
	return baidu
}

func (g *BaiduPro) DescribeRegions() []model.Region {
	data := []model.Region{}
	for k, v := range g.RegionArray {
		data = append(data, model.Region{
			RegionId: k,
			Name:     v,
		})
	}
	return data
}

func (g *BaiduPro) GetRegionMapValue(Region string) string {
	if _, ok := g.RegionMap[Region]; ok {
		return g.RegionMap[Region]
	}
	return Region
}

func (g *BaiduPro) GetRegionMapKey(Region string) string {
	for k, v := range g.RegionMap {
		if v == Region {
			return k
		}
	}
	return Region
}

func (g *BaiduPro) GetDiskCategory(category string) string {
	var DiskCategoryMap = map[string]string{
		"cloud":            "hp1",
		"cloud_efficiency": "hp1",
		"cloud_ssd":        "cloud_hp1",
		"cloud_essd":       "enhanced_ssd_pl1",
		"ephemeral_ssd":    "local-ssd",
		"ephemeral":        "local",
	}

	for k, v := range DiskCategoryMap {
		if category == k {
			return v
		}
	}
	return category
}

func (g *BaiduPro) GetEndPoint(Region string) string {
	return "bcc." + Region + ".baidubce.com"
}
