package gateway

import (
	"errors"

	"github.com/suguer/EcsGateway/model"
	"github.com/suguer/EcsGateway/private/http"
)

type GatewayInterface interface {
	Init(c *model.Config)

	DownloadCache()

	DescribeRegions() ([]model.Region, error)

	DescribePrice(instance *model.Instance) (*model.PriceInfo, error)

	DescribeAvailableInstance(instance *model.Instance) ([]model.Instance, error)
}

type Gateway struct {
	Conf       *model.Config
	httpClient http.HttpClient
	ApiUrl     string
}

func (g *Gateway) send(request *http.HttpRequest) (*http.HttpResponse, error) {
	request.SetURL(g.ApiUrl + request.GetPath())
	response, err := g.httpClient.Send(request)
	return response, err
}

func (g *Gateway) DescribeRegions() ([]model.Region, error) {
	return []model.Region{}, errors.New("unsupport")
}

func (g *Gateway) DescribePrice(instance *model.Instance) (*model.PriceInfo, error) {
	return &model.PriceInfo{}, errors.New("unsupport")
}

func (g *Gateway) DescribeAvailableInstance(instance *model.Instance) ([]model.Instance, error) {
	return []model.Instance{}, errors.New("unsupport")
}

func (g *Gateway) DownloadCache() {
}

func (g *Gateway) Init(c *model.Config) {
	g.Conf = c
	g.httpClient = http.NewHttpClient()
}

func NewGatewayInterface(platform string, c *model.Config) (GatewayInterface, error) {
	var gateway GatewayInterface
	switch platform {
	case "aliyun":
		gateway = &AliyunGateway{}
	case "ucloud":
		gateway = &UcloudGateway{}
	case "tencent":
		gateway = &TencentGateway{}
	case "baidu":
		gateway = &BaiduGateway{}
	default:
		return nil, errors.New("unspport")
	}
	gateway.Init(c)
	return gateway, nil

}
