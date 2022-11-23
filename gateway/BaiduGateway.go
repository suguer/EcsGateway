package gateway

import (
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/suguer/EcsGateway/model"
	"github.com/suguer/EcsGateway/private/utils"
	"github.com/suguer/EcsGateway/pro"
)

type BaiduGateway struct {
	Gateway
	client *bcc.Client
}

func (g *BaiduGateway) Init(c *model.Config) {
	g.Gateway.Init(c)
	g.client, _ = bcc.NewClient(g.Conf.AppID, g.Conf.AppSecret, "")

}

func (g *BaiduGateway) DescribeRegions() ([]model.Region, error) {
	pro := pro.NewBaiduPro()
	data := pro.DescribeRegions()
	for k, r := range data {
		data[k].RegionId = pro.GetRegionMapValue(r.RegionId)
	}
	model.Region{}.Save(data, "baidu")
	return data, nil
}

func (g *BaiduGateway) DescribeAvailableInstance(instance *model.Instance) ([]model.Instance, error) {
	InstanceData := []model.Instance{}
	baidu := pro.NewBaiduPro()
	Region := baidu.GetRegionMapKey(instance.RegionId)
	g.client.Config.Endpoint = baidu.GetEndPoint(Region)
	request := &api.ListFlavorSpecArgs{}
	response, err := g.client.ListFlavorSpec(request)
	if err != nil {
		return InstanceData, err
	}
	// content, _ := json.Marshal(response)
	// fmt.Printf("%v\n", string(content))
	exist_check := make(map[string]int)
	for _, zrds := range response.ZoneResources {
		for _, fg := range zrds.BccResources.FlavorGroups {
			for _, bf := range fg.Flavors {
				if bf.ProductType == "PostPaid" {
					continue
				}
				if instance.Hardware.CpuCount > 0 && instance.Hardware.CpuCount != bf.CpuCount {
					continue
				}
				if instance.Hardware.MemoryCapacityInMB > 0 && instance.Hardware.MemoryCapacityInMB != bf.MemoryCapacityInGB*1024 {
					continue
				}
				//防重插入检测
				if _, ok := exist_check[zrds.ZoneName+bf.Spec]; ok {
					continue
				}
				exist_check[zrds.ZoneName+bf.Spec] = 1
				var temp model.Instance
				instance.DeepCopy(&temp)
				temp.Hardware.InstanceType = bf.Spec
				temp.ZoneId = zrds.ZoneName
				InstanceData = append(InstanceData, temp)
			}
		}
	}
	return InstanceData, nil
}
func (g *BaiduGateway) DescribePrice(instance *model.Instance) (*model.PriceInfo, error) {
	data := &model.PriceInfo{}
	pro := pro.NewBaiduPro()
	Region := pro.GetRegionMapKey(instance.RegionId)
	g.client.Config.Endpoint = pro.GetEndPoint(Region)
	Period := instance.Period
	if instance.PriceUnit == "Year" {
		Period *= 12
	}
	request := &api.GetPriceBySpecArgs{
		Spec:           instance.Hardware.InstanceType,
		ZoneName:       instance.ZoneId,
		PurchaseLength: Period,
		PaymentTiming:  "Prepaid",
	}
	response, err := g.client.GetPriceBySpec(request)
	if err != nil {
		return data, err
	}
	data.OriginalPrice += utils.ToFloat32(response.Price[0].SpecPrices[0].SpecPrice)
	data.TradePrice += utils.ToFloat32(response.Price[0].SpecPrices[0].SpecPrice)
	data.PriceInfoItem = append(data.PriceInfoItem, model.PriceInfoItem{
		ResourceType:  "ecs",
		OriginalPrice: utils.ToFloat32(response.Price[0].SpecPrices[0].SpecPrice),
		TradePrice:    utils.ToFloat32(response.Price[0].SpecPrices[0].SpecPrice),
	})
	return data, nil
}
