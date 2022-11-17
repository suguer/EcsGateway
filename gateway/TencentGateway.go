package gateway

import (
	"github.com/alibabacloud-go/tea/tea"
	"github.com/suguer/EcsGateway/model"
	"github.com/suguer/EcsGateway/pro"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	region "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/region/v20220627"
)

type TencentGateway struct {
	Gateway
	credential *common.Credential
	cpf        *profile.ClientProfile
}

func (g *TencentGateway) Init(c *model.Config) {
	g.Gateway.Init(c)
	g.credential = common.NewCredential(
		g.Conf.AppID,
		g.Conf.AppSecret,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	g.cpf = profile.NewClientProfile()
	g.cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"

}

func (g *TencentGateway) DescribeRegions() ([]model.Region, error) {
	data := []model.Region{}
	g.cpf.HttpProfile.Endpoint = "region.tencentcloudapi.com"
	client, _ := region.NewClient(g.credential, "", g.cpf)
	request := region.NewDescribeRegionsRequest()
	request.Product = common.StringPtr("cvm")
	response, err := client.DescribeRegions(request)
	if err != nil {
		return data, err
	}
	tencent := pro.NewTencentPro()
	for _, ri := range response.Response.RegionSet {
		if *ri.RegionState != "AVAILABLE" {
			continue
		}
		Region := tencent.GetRegionMapValue(*ri.Region)
		data = append(data, model.Region{
			RegionId: Region,
			Name:     *ri.RegionName,
		})
	}
	model.Region{}.Save(data, "tencent")
	return data, nil
}

func (g *TencentGateway) DescribeAvailableInstance(instance *model.Instance) ([]model.Instance, error) {
	InstanceData := []model.Instance{}
	tencent := pro.NewTencentPro()
	Region := tencent.GetRegionMapKey(instance.RegionId)
	client, _ := cvm.NewClient(g.credential, Region, g.cpf)
	request := cvm.NewDescribeInstanceTypeConfigsRequest()
	response, err := client.DescribeInstanceTypeConfigs(request)
	if err != nil {
		return InstanceData, err
	}
	for _, itc := range response.Response.InstanceTypeConfigSet {
		if instance.Hardware.CpuCount > 0 && instance.Hardware.CpuCount != int(*itc.CPU) {
			continue
		}
		if instance.Hardware.MemoryCapacityInMB > 0 && instance.Hardware.MemoryCapacityInMB != int(*itc.Memory)*1024 {
			continue
		}
		var temp model.Instance
		instance.DeepCopy(&temp)
		temp.Hardware.InstanceType = *itc.InstanceType
		temp.ZoneId = *itc.Zone
		InstanceData = append(InstanceData, temp)
	}
	return InstanceData, nil
}

func (g *TencentGateway) DescribePrice(instance *model.Instance) (*model.PriceInfo, error) {
	data := &model.PriceInfo{}
	tencent := pro.NewTencentPro()
	RegionId := tencent.GetRegionMapKey(instance.RegionId)
	client, _ := cvm.NewClient(g.credential, RegionId, g.cpf)
	request := cvm.NewInquiryPriceRunInstancesRequest()
	request.InstanceChargeType = tea.String("PREPAID")
	request.Placement = &cvm.Placement{
		Zone: &instance.ZoneId,
	}
	request.InstanceType = &instance.Hardware.InstanceType
	Period := instance.Period
	if instance.PriceUnit == "Year" {
		Period *= 12
	}
	request.InstanceChargePrepaid = &cvm.InstanceChargePrepaid{
		Period: tea.Int64(int64(Period)),
	}
	if instance.Network.InternetMaxBandwidthOut > 0 {
		request.InternetAccessible = &cvm.InternetAccessible{
			InternetChargeType:      tea.String("BANDWIDTH_PREPAID"),
			PublicIpAssigned:        tea.Bool(true),
			InternetMaxBandwidthOut: tea.Int64(int64(instance.Network.InternetMaxBandwidthOut)),
		}
	}
	OSType := "Windows"
	Image, _ := g.DescribeImage(instance.RegionId, OSType)
	request.ImageId = tea.String(Image.ImageId)
	response, err := client.InquiryPriceRunInstances(request)
	if err != nil {
		return data, err
	}
	if *response.Response.Price.InstancePrice.DiscountPrice > 0 {
		data.OriginalPrice += float32(*response.Response.Price.InstancePrice.OriginalPrice)
		data.TradePrice += float32(*response.Response.Price.InstancePrice.DiscountPrice)
		data.PriceInfoItem = append(data.PriceInfoItem, model.PriceInfoItem{
			ResourceType:  "ecs",
			OriginalPrice: float32(*response.Response.Price.InstancePrice.OriginalPrice),
			TradePrice:    float32(*response.Response.Price.InstancePrice.DiscountPrice),
		})
	}
	if *response.Response.Price.BandwidthPrice.DiscountPrice > 0 {
		data.OriginalPrice += float32(*response.Response.Price.BandwidthPrice.OriginalPrice)
		data.TradePrice += float32(*response.Response.Price.BandwidthPrice.DiscountPrice)
		data.PriceInfoItem = append(data.PriceInfoItem, model.PriceInfoItem{
			ResourceType:  "eip",
			OriginalPrice: float32(*response.Response.Price.BandwidthPrice.OriginalPrice),
			TradePrice:    float32(*response.Response.Price.BandwidthPrice.DiscountPrice),
		})
	}
	return data, nil
}

func (g *TencentGateway) DescribeImage(RegionId string, OSType string) (*model.Image, error) {
	Image := &model.Image{}
	tencent := pro.NewTencentPro()
	client, _ := cvm.NewClient(g.credential, tencent.GetRegionMapKey(RegionId), g.cpf)
	request := cvm.NewDescribeImagesRequest()
	request.Filters = append(request.Filters, &cvm.Filter{
		Name:   tea.String("image-type"),
		Values: tea.StringSlice([]string{"PUBLIC_IMAGE"}),
	})
	request.Filters = append(request.Filters, &cvm.Filter{
		Name:   tea.String("platform"),
		Values: tea.StringSlice([]string{OSType}),
	})
	request.Limit = tea.Uint64(1)
	response, err := client.DescribeImages(request)
	if err != nil {
		return Image, err
	}
	Image.ImageId = *response.Response.ImageSet[0].ImageId
	return Image, nil
}
