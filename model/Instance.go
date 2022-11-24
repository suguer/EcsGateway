package model

type Hardware struct {
	InstanceType       string
	CpuCount           int
	MemoryCapacityInMB int
	CpuPlatform        string
}
type Network struct {
	InternetChargeType      string
	InternetMaxBandwidthOut int
}
type Disk struct {
	Category string
	Size     int
}
type Image struct {
	OSName  string
	ImageId string
	OSType  string
}
type Instance struct {
	Platform string
	RegionId string
	ZoneId   string
	Image    Image
	// ImageId  string
	// ImageSystem string
	Hardware   Hardware
	Network    Network
	SystemDisk Disk
	DataDisk   []Disk
	Period     int
	PeriodUnit string
}

func (m *Instance) DeepCopy(desc *Instance) {
	*desc = *m
}

var DiskCategory = []string{
	"cloud",
	"cloud_efficiency",
	"cloud_ssd",
	"ephemeral_ssd",
	"ephemeral",
	"cloud_essd",
}

func (instance *Instance) GetSystemDiskCategory() string {
	for _, v := range DiskCategory {
		if v == instance.SystemDisk.Category {
			return instance.SystemDisk.Category
		}
	}
	return "cloud_efficiency"
}
