package model

type Hardware struct {
	InstanceType       string
	CpuCount           int
	MemoryCapacityInGB int
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
type Instance struct {
	RegionId   string
	ZoneId     string
	ImageId    string
	Hardware   Hardware
	Network    Network
	SystemDisk Disk
	DataDisk   []Disk
	Period     int
	PriceUnit  string
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
