package gateway

import (
	"github.com/alibabacloud-go/tea/tea"
	"github.com/suguer/EcsGateway/model"
	"github.com/suguer/EcsGateway/pro"
	"github.com/ucloud/ucloud-sdk-go/services/uaccount"
	"github.com/ucloud/ucloud-sdk-go/services/uhost"
	"github.com/ucloud/ucloud-sdk-go/services/unet"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
)

type UcloudGateway struct {
	Gateway
	client     *uhost.UHostClient
	credential auth.Credential
	cfg        ucloud.Config
}

func (g *UcloudGateway) Init(c *model.Config) {
	g.Gateway.Init(c)
	g.cfg = ucloud.NewConfig()
	g.credential = auth.NewCredential()
	g.credential.PrivateKey = g.Conf.AppSecret
	g.credential.PublicKey = g.Conf.AppID
	g.client = uhost.NewClient(&g.cfg, &g.credential)
}

func (g *UcloudGateway) DescribeRegions() ([]model.Region, error) {
	data := []model.Region{}
	ucloud := pro.NewUcloudPro()
	uaccountClient := uaccount.NewClient(&g.cfg, &g.credential)
	GetRegionResponse, _ := uaccountClient.GetRegion(&uaccount.GetRegionRequest{})
	existRegion := make(map[string]int)
	for _, value := range GetRegionResponse.Regions {
		if _, ok := existRegion[value.Region]; !ok {
			Region := ucloud.GetRegionMapValue(value.Region)
			Name := ucloud.GetRegionMapChinese(value.Region)
			data = append(data, model.Region{
				RegionId: Region,
				Name:     Name,
			})
			existRegion[value.Region] = 1
		}

	}

	model.Region{}.Save(data, "ucloud")
	return data, nil
}

func (g *UcloudGateway) DescribeAvailableInstance(instance *model.Instance) ([]model.Instance, error) {
	InstanceData := []model.Instance{}
	ucloud := pro.NewUcloudPro()
	Region := ucloud.GetRegionMapKey(instance.RegionId)
	request := g.client.NewDescribeAvailableInstanceTypesRequest()
	request.Region = &Region
	response, err := g.client.DescribeAvailableInstanceTypes(request)
	if err != nil {
		return InstanceData, err
	}
	// fmt.Printf("err: %v\n", err)
	// fmt.Printf("response: %+v\n", response)
	for _, ait := range response.AvailableInstanceTypes {
		var temp model.Instance
		instance.DeepCopy(&temp)

		temp.SystemDisk.Category = ait.Disks[0].BootDisk[0].Name
		if len(ait.CpuPlatforms.Intel) > 0 {
			temp.Hardware.InstanceType = ait.Name
			temp.ZoneId = ait.Zone
			temp.Hardware.CpuPlatform = "Intel"
			InstanceData = append(InstanceData, temp)
		}
		if len(ait.CpuPlatforms.Amd) > 0 {
			temp.Hardware.InstanceType = ait.Name
			temp.ZoneId = ait.Zone
			temp.Hardware.CpuPlatform = "Amd"
			InstanceData = append(InstanceData, temp)
		}
	}
	return InstanceData, nil
}

func (g *UcloudGateway) DescribePrice(instance *model.Instance) (*model.PriceInfo, error) {
	data := &model.PriceInfo{}
	ucloud := pro.NewUcloudPro()
	Region := ucloud.GetRegionMapKey(instance.RegionId)
	GetUHostInstancePriceRequest := g.client.NewGetUHostInstancePriceRequest()
	GetUHostInstancePriceRequest.Count = tea.Int(1)
	GetUHostInstancePriceRequest.Region = &Region
	GetUHostInstancePriceRequest.Zone = &instance.ZoneId
	GetUHostInstancePriceRequest.ChargeType = &instance.PeriodUnit
	GetUHostInstancePriceRequest.Quantity = &instance.Period
	GetUHostInstancePriceRequest.CpuPlatform = &instance.Hardware.CpuPlatform
	GetUHostInstancePriceRequest.MachineType = &instance.Hardware.InstanceType
	GetUHostInstancePriceRequest.CPU = &instance.Hardware.CpuCount
	GetUHostInstancePriceRequest.Memory = &instance.Hardware.MemoryCapacityInMB
	GetUHostInstancePriceRequest.Disks = append(GetUHostInstancePriceRequest.Disks, uhost.UHostDisk{
		IsBoot: tea.String("True"),
		Size:   &instance.SystemDisk.Size,
		Type:   &instance.SystemDisk.Category,
	})
	response, err := g.client.GetUHostInstancePrice(GetUHostInstancePriceRequest)
	if err != nil {
		data.OriginalPrice += 0
		data.TradePrice += 0
		data.PriceInfoItem = append(data.PriceInfoItem, model.PriceInfoItem{
			ResourceType:  "ecs",
			OriginalPrice: 0,
			TradePrice:    0,
			Code:          err.Error(),
		})
	} else {
		data.OriginalPrice += float32(response.PriceSet[0].OriginalPrice)
		data.TradePrice += float32(response.PriceSet[0].Price)
		data.PriceInfoItem = append(data.PriceInfoItem, model.PriceInfoItem{
			ResourceType:  "ecs",
			OriginalPrice: float32(response.PriceSet[0].OriginalPrice),
			TradePrice:    float32(response.PriceSet[0].Price),
		})
	}
	//计算IP
	if instance.Network.InternetMaxBandwidthOut > 0 {
		ucloud := pro.NewUcloudPro()
		unetClient := unet.NewClient(&g.cfg, &g.credential)
		req := unetClient.NewGetEIPPriceRequest()
		req.Region = &Region
		req.Bandwidth = &instance.Network.InternetMaxBandwidthOut
		req.ChargeType = &instance.PeriodUnit
		req.Quantity = &instance.Period

		req.OperatorName = tea.String(ucloud.GetOperatorArray(Region))
		EIPPriceResponse, err := unetClient.GetEIPPrice(req)
		if err != nil {
			data.PriceInfoItem = append(data.PriceInfoItem, model.PriceInfoItem{
				ResourceType:  "eip",
				OriginalPrice: 0,
				TradePrice:    0,
				Code:          err.Error(),
			})
		} else {
			data.OriginalPrice += float32(EIPPriceResponse.PriceSet[0].OriginalPrice)
			data.TradePrice += float32(EIPPriceResponse.PriceSet[0].Price)
			data.PriceInfoItem = append(data.PriceInfoItem, model.PriceInfoItem{
				ResourceType:  "eip",
				OriginalPrice: float32(EIPPriceResponse.PriceSet[0].OriginalPrice),
				TradePrice:    float32(EIPPriceResponse.PriceSet[0].Price),
			})
		}
	}
	return data, nil
}
