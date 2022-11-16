package gateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	nethttp "net/http"
	"os"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/suguer/EcsGateway/config"
	"github.com/suguer/EcsGateway/model"
	"github.com/suguer/EcsGateway/model/aliyun"
	"github.com/suguer/EcsGateway/private/http"
	"github.com/suguer/EcsGateway/private/utils"
)

type AliyunGateway struct {
	Gateway
	client *ecs20140526.Client
}

func (g *AliyunGateway) Init(c *model.Config) {
	g.Gateway.Init(c)
	g.ApiUrl = "https://ecs.aliyuncs.com"
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: tea.String(c.AppID),
		// 您的AccessKey Secret
		AccessKeySecret: tea.String(c.AppSecret),
		// 访问的 Region
		RegionId: tea.String("cn-hangzhou"),
	}
	g.client, _ = ecs20140526.NewClient(config)
}

func (g *AliyunGateway) DescribeRegions() ([]model.Region, error) {
	data := []model.Region{}
	DescribePriceRequest := &ecs20140526.DescribeRegionsRequest{}
	response, err := g.client.DescribeRegions(DescribePriceRequest)
	if err != nil {
		return data, err
	}
	for _, value := range response.Body.Regions.Region {
		data = append(data, model.Region{
			RegionId: *value.RegionId,
			Name:     *value.LocalName,
		})
	}
	model.Region{}.Save(data, "aliyun")
	return data, nil
}

func (g *AliyunGateway) DescribePrice(instance *model.Instance) (*model.PriceInfo, error) {
	data := &model.PriceInfo{}
	DescribePriceRequest := &ecs20140526.DescribePriceRequest{
		RegionId:     &instance.RegionId,
		ZoneId:       &instance.ZoneId,
		Period:       tea.ToInt32(&instance.Period),
		PriceUnit:    &instance.PriceUnit,
		ResourceType: tea.String("instance"),
		InstanceType: &instance.Hardware.InstanceType,
	}
	DescribePriceRequest.SystemDisk = &ecs20140526.DescribePriceRequestSystemDisk{
		Size:     tea.ToInt32(&instance.SystemDisk.Size),
		Category: tea.String(instance.GetSystemDiskCategory()),
	}
	response, err := g.client.DescribePrice(DescribePriceRequest)
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
		data.OriginalPrice += *response.Body.PriceInfo.Price.OriginalPrice
		data.TradePrice += *response.Body.PriceInfo.Price.TradePrice
		data.PriceInfoItem = append(data.PriceInfoItem, model.PriceInfoItem{
			ResourceType:  "ecs",
			OriginalPrice: *response.Body.PriceInfo.Price.OriginalPrice,
			TradePrice:    *response.Body.PriceInfo.Price.TradePrice,
		})
	}

	//计算IP
	if instance.Network.InternetMaxBandwidthOut > 0 {
		month := instance.Period
		if instance.PriceUnit == "Year" {
			month = instance.Period * 12
		}

		file, _ := ioutil.ReadFile("storage/aliyun_bandWidthPrice.json")
		var BandwidthPrice aliyun.BandwidthPrice
		json.Unmarshal(file, &BandwidthPrice)
		var BandwidthOriginalPrice float32 = 0
		Code := ""
		if _, ok := BandwidthPrice.PricingInfo[instance.RegionId]; ok {
			if instance.Network.InternetMaxBandwidthOut > 5 {
				BandwidthOriginalPrice = utils.ToFloat32(BandwidthPrice.PricingInfo[instance.RegionId].Months[5].Price) +
					(utils.ToFloat32(BandwidthPrice.PricingInfo[instance.RegionId].Months[6].Price)-utils.ToFloat32(BandwidthPrice.PricingInfo[instance.RegionId].Months[5].Price))*
						float32((instance.Network.InternetMaxBandwidthOut-5))
			} else {
				BandwidthOriginalPrice = utils.ToFloat32(BandwidthPrice.PricingInfo[instance.RegionId].Months[instance.Network.InternetMaxBandwidthOut].Price)
				fmt.Printf("BandwidthOriginalPrice: %v\n", BandwidthOriginalPrice)
			}
			BandwidthOriginalPrice = BandwidthOriginalPrice * float32(month)
		} else {
			Code = "unspport"
		}
		data.OriginalPrice += BandwidthOriginalPrice
		data.TradePrice += BandwidthOriginalPrice
		data.PriceInfoItem = append(data.PriceInfoItem, model.PriceInfoItem{
			ResourceType:  "eip",
			OriginalPrice: BandwidthOriginalPrice,
			TradePrice:    BandwidthOriginalPrice,
			Code:          Code,
		})
	}
	return data, nil
}

func (g *AliyunGateway) DescribeAvailableInstance(instance *model.Instance) ([]model.Instance, error) {
	InstanceData := []model.Instance{}
	request := &ecs20140526.DescribeAvailableResourceRequest{
		DestinationResource: tea.String("InstanceType"),
		InstanceChargeType:  tea.String("PrePaid"),
		RegionId:            &instance.RegionId,
		Cores:               tea.ToInt32(&instance.Hardware.CpuCount),
		Memory:              tea.Float32(float32(instance.Hardware.MemoryCapacityInGB / 1024)),
	}
	response, err := g.client.DescribeAvailableResource(request)
	if err != nil {
		return InstanceData, err
	}
	for _, AvailableZone := range response.Body.AvailableZones.AvailableZone {

		for _, SupportedResource := range AvailableZone.AvailableResources.AvailableResource[0].SupportedResources.SupportedResource {
			if *SupportedResource.Status != "Available" {
				continue
			}
			var temp model.Instance
			instance.DeepCopy(&temp)
			temp.Hardware.InstanceType = *SupportedResource.Value
			temp.ZoneId = *AvailableZone.ZoneId
			InstanceData = append(InstanceData, temp)
		}
	}
	return InstanceData, nil
}

func (g *AliyunGateway) buildParam(request *http.HttpRequest) {
	Timestamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	request.SetQuery("Format", "JSON")
	request.SetQuery("Version", "2014-05-26")
	request.SetQuery("AccessKeyId", g.Conf.AppID)
	request.SetQuery("SignatureNonce", utils.GetSignatureNonce())
	request.SetQuery("Timestamp", Timestamp)
	request.SetQuery("SignatureMethod", "HMAC-SHA1")
	request.SetQuery("SignatureVersion", "1.0")
	CanonicalizedQueryString := utils.GetSignString(request.GetQueryMap(), "full")
	stringToSign := "GET" + "&" + utils.EncodeURIComponent("/") + "&" + utils.EncodeURIComponent(CanonicalizedQueryString)
	Signature := utils.HMACSHA1(g.Conf.AppSecret+"&", stringToSign)
	request.SetQuery("Signature", Signature)
}

func (g *AliyunGateway) DownloadCache() {
	resp, _ := nethttp.Get("https://g.alicdn.com/aliyun/ecs-price-info/2.0.211/price/download/bandWidthPrice.json")
	body, _ := ioutil.ReadAll(resp.Body)
	file, err := os.OpenFile(config.StoragePath+"/aliyun_bandWidthPrice.json", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("文件打开/创建失败,原因是:", err)
		return
	}
	file.WriteString(string(body))
	file.Close()
}
