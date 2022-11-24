package gateway

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
	"github.com/aws/aws-sdk-go-v2/service/pricing/types"
	"github.com/suguer/EcsGateway/model"
	awsModel "github.com/suguer/EcsGateway/model/aws"
	"github.com/suguer/EcsGateway/private/utils"
	"github.com/suguer/EcsGateway/pro"
)

type AwsCloudGateway struct {
	Gateway
	cfg    aws.Config
	client *ec2.Client
	option ec2.Options
}

func (g *AwsCloudGateway) Init(c *model.Config) {
	g.Gateway.Init(c)
	g.option = ec2.Options{
		Region:      "ap-east-1",
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(g.Conf.AppID, g.Conf.AppSecret, "")),
	}
	g.client = ec2.New(g.option)
}

func (g *AwsCloudGateway) DescribeRegions() ([]model.Region, error) {
	pro := pro.NewAwsPro()
	data := pro.DescribeRegions()
	for k, r := range data {
		data[k].RegionId = pro.GetRegionMapValue(r.RegionId)
	}
	model.Region{}.Save(data, "aws")
	return data, nil
}

func (g *AwsCloudGateway) DescribeAvailableInstance(instance *model.Instance) ([]model.Instance, error) {

	pro := pro.NewAwsPro()
	Region := pro.GetRegionMapKey(instance.RegionId)

	InstanceData := []model.Instance{}
	client := pricing.New(pricing.Options{
		Region:      "ap-south-1",
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(g.Conf.AppID, g.Conf.AppSecret, "")),
	})
	request := &pricing.GetProductsInput{
		FormatVersion: tea.String("aws_v1"),
		ServiceCode:   tea.String("AmazonEC2"),
		MaxResults:    tea.Int32(100),
	}
	var Filters []types.Filter
	Filters = append(Filters, types.Filter{
		Field: tea.String("preInstalledSw"),
		Type:  types.FilterType("TERM_MATCH"),
		Value: tea.String("NA"),
	}, types.Filter{
		Field: tea.String("regionCode"),
		Type:  types.FilterType("TERM_MATCH"),
		Value: tea.String(Region),
	}, types.Filter{
		Field: tea.String("tenancy"),
		Type:  types.FilterType("TERM_MATCH"),
		Value: tea.String("Shared"),
	}, types.Filter{
		Field: tea.String("capacitystatus"),
		Type:  types.FilterType("TERM_MATCH"),
		Value: tea.String("UnusedCapacityReservation"),
	})
	if instance.Hardware.CpuCount > 0 {
		Filters = append(Filters, types.Filter{
			Field: tea.String("vcpu"),
			Type:  types.FilterType("TERM_MATCH"),
			Value: utils.IntToString(instance.Hardware.CpuCount),
		})
	}
	if instance.Hardware.MemoryCapacityInMB > 0 {
		temp := utils.Float32ToStringValue(float32(instance.Hardware.MemoryCapacityInMB)/1024) + " GiB"
		Filters = append(Filters, types.Filter{
			Field: tea.String("memory"),
			Type:  types.FilterType("TERM_MATCH"),
			Value: tea.String(temp),
		})
	}
	request.Filters = Filters
	response, err := client.GetProducts(context.TODO(), request)
	if err != nil {
		return InstanceData, err
	}
	for _, v := range response.PriceList {
		var GetProductItem awsModel.GetProductItem
		json.Unmarshal([]byte(v), &GetProductItem)
		var temp model.Instance
		isContinue := false
		for _, ti := range GetProductItem.Terms.OnDemand {
			for _, pd := range ti.PriceDimensions {
				if pd.PricePerUnit.USD == "0.0000000000" {
					isContinue = true
					break
				}
			}
		}
		if isContinue {
			continue
		}

		instance.DeepCopy(&temp)
		temp.Hardware.InstanceType = GetProductItem.Product.Attributes.InstanceType
		temp.Hardware.CpuCount = utils.StringToInt(GetProductItem.Product.Attributes.Vcpu)
		temp.Image = model.Image{
			OSType: GetProductItem.Product.Attributes.OperatingSystem,
		}
		InstanceData = append(InstanceData, temp)
	}
	return InstanceData, nil
}

func (g *AwsCloudGateway) DescribePrice(instance *model.Instance) (*model.PriceInfo, error) {
	data := &model.PriceInfo{
		PriceUnit:     "USD",
		OriginalPrice: 0,
		TradePrice:    0,
	}
	pro := pro.NewAwsPro()
	Region := pro.GetRegionMapKey(instance.RegionId)
	client := pricing.New(pricing.Options{
		Region:      "ap-south-1",
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(g.Conf.AppID, g.Conf.AppSecret, "")),
	})
	request := &pricing.GetProductsInput{
		FormatVersion: tea.String("aws_v1"),
		ServiceCode:   tea.String("AmazonEC2"),
		MaxResults:    tea.Int32(100),
	}
	var Filters []types.Filter
	temp := utils.Float32ToStringValue(float32(instance.Hardware.MemoryCapacityInMB)/1024) + " GiB"
	Filters = append(Filters,
		types.Filter{
			Field: tea.String("preInstalledSw"),
			Type:  types.FilterType("TERM_MATCH"),
			Value: tea.String("NA"),
		},
		types.Filter{
			Field: tea.String("regionCode"),
			Type:  types.FilterType("TERM_MATCH"),
			Value: tea.String(Region),
		},
		types.Filter{
			Field: tea.String("vcpu"),
			Type:  types.FilterType("TERM_MATCH"),
			Value: utils.IntToString(instance.Hardware.CpuCount),
		},
		types.Filter{
			Field: tea.String("memory"),
			Type:  types.FilterType("TERM_MATCH"),
			Value: tea.String(temp),
		},
		types.Filter{
			Field: tea.String("operatingSystem"),
			Type:  types.FilterType("TERM_MATCH"),
			Value: tea.String(instance.Image.OSType),
		},
		types.Filter{
			Field: tea.String("instanceType"),
			Type:  types.FilterType("TERM_MATCH"),
			Value: tea.String(instance.Hardware.InstanceType),
		},
		types.Filter{
			Field: tea.String("tenancy"),
			Type:  types.FilterType("TERM_MATCH"),
			Value: tea.String("Shared"),
		},
		types.Filter{
			Field: tea.String("capacitystatus"),
			Type:  types.FilterType("TERM_MATCH"),
			Value: tea.String("UnusedCapacityReservation"),
		})
	request.Filters = Filters
	response, err := client.GetProducts(context.TODO(), request)
	fmt.Printf("instance: %v\n", instance)
	fmt.Printf("response: %v\n", response)
	if err != nil {
		return data, err
	}
	for _, v := range response.PriceList {
		// fmt.Printf("v: %v\n", v)
		var GetProductItem awsModel.GetProductItem
		json.Unmarshal([]byte(v), &GetProductItem)
		for _, ti := range GetProductItem.Terms.OnDemand {
			for _, pd := range ti.PriceDimensions {
				Period := instance.Period * 30 * 24
				if instance.PeriodUnit == "Year" {
					Period = instance.Period * 365 * 24
				}
				// fmt.Printf("Period: %v\n", Period)
				PricePerUnitPrice := utils.ToFloat32(pd.PricePerUnit.USD)
				data.OriginalPrice = PricePerUnitPrice * float32(Period)
				data.TradePrice = PricePerUnitPrice * float32(Period)
				data.PriceInfoItem = append(data.PriceInfoItem, model.PriceInfoItem{
					TradePrice:    PricePerUnitPrice * float32(Period),
					OriginalPrice: PricePerUnitPrice * float32(Period),
					ResourceType:  "ecs",
				})
			}
		}
	}
	return data, nil
}
