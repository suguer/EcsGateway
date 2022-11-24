package gateway

import (
	"fmt"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/jdcloud-api/jdcloud-sdk-go/core"
	billing "github.com/jdcloud-api/jdcloud-sdk-go/services/billing/apis"
	billingClient "github.com/jdcloud-api/jdcloud-sdk-go/services/billing/client"
	billingModels "github.com/jdcloud-api/jdcloud-sdk-go/services/billing/models"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/apis"
	"github.com/jdcloud-api/jdcloud-sdk-go/services/vm/client"
	"github.com/suguer/EcsGateway/model"
	"github.com/suguer/EcsGateway/pro"
)

type JdCloudGateway struct {
	Gateway
	credentials *core.Credential
	client      *client.VmClient
}

func (g *JdCloudGateway) Init(c *model.Config) {
	g.Gateway.Init(c)
	g.credentials = core.NewCredentials(g.Conf.AppID, g.Conf.AppSecret)
	g.client = client.NewVmClient(g.credentials)
	g.client.DisableLogger()

}

func (g *JdCloudGateway) DescribeRegions() ([]model.Region, error) {
	baidu := pro.NewJdCloudPro()
	data := baidu.DescribeRegions()
	for k, r := range data {
		data[k].RegionId = baidu.GetRegionMapValue(r.RegionId)
	}
	model.Region{}.Save(data, "jdcloud")
	return data, nil
}

func (g *JdCloudGateway) DescribeAvailableInstance(instance *model.Instance) ([]model.Instance, error) {
	InstanceData := []model.Instance{}
	pro := pro.NewJdCloudPro()
	Region := pro.GetRegionMapKey(instance.RegionId)
	req := apis.NewDescribeInstanceTypesRequest(Region)
	response, err := g.client.DescribeInstanceTypes(req)
	if err != nil {
		return InstanceData, err
	}
	for _, it := range response.Result.InstanceTypes {
		if instance.Hardware.CpuCount > 0 && instance.Hardware.CpuCount != it.Cpu {
			continue
		}
		if instance.Hardware.MemoryCapacityInMB > 0 && instance.Hardware.MemoryCapacityInMB != it.MemoryMB {
			continue
		}
		for _, its := range it.State {
			if its.InStock == false {
				continue
			}
			var temp model.Instance
			instance.DeepCopy(&temp)
			temp.Hardware.InstanceType = it.InstanceType
			temp.Hardware.CpuCount = it.Cpu
			temp.Hardware.MemoryCapacityInMB = it.MemoryMB
			temp.ZoneId = its.Az
			InstanceData = append(InstanceData, temp)
		}
	}
	return InstanceData, nil
}

func (g *JdCloudGateway) DescribePrice(instance *model.Instance) (*model.PriceInfo, error) {
	data := &model.PriceInfo{}
	pro := pro.NewJdCloudPro()
	Region := pro.GetRegionMapKey(instance.RegionId)
	reqeust := billing.NewCalculateTotalPriceRequest(Region, 1, 1)
	TimeUnit := 3
	if instance.PeriodUnit == "Year" {
		TimeUnit = 4
	}
	reqeust.OrderList = append(reqeust.OrderList, billingModels.OrderPriceProtocol{
		AppCode:         tea.String("jcloud"),
		ServiceCode:     tea.String("vm"),
		Site:            tea.Int(0),
		Region:          tea.String(Region),
		BillingType:     tea.Int(3),
		TimeSpan:        tea.Int(instance.Period),
		TimeUnit:        tea.Int(TimeUnit),
		NetworkOperator: tea.Int(0),
		Formula: []billingModels.Formula{
			billingModels.Formula{
				Key:   &instance.Hardware.InstanceType,
				Value: tea.Float64(1),
			},
		},
	})
	client := billingClient.NewBillingClient(g.credentials)
	client.DisableLogger()

	response, err := client.CalculateTotalPrice(reqeust)
	fmt.Printf("response: %+v\n", response)
	if err != nil {
		return data, err
	}

	return data, nil
}
