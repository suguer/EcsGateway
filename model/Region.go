package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/suguer/EcsGateway/private/utils"
)

type Region struct {
	RegionId string
	Name     string
	Platform []string
}

const StoragePath = ""
const RegionPath = StoragePath + "/RegionData.json"

func (Region) Save(temp []Region, platform string) {
	pwd, _ := os.Getwd()
	path := pwd + RegionPath
	read, err := ioutil.ReadFile(path)
	var data []Region
	if err == nil {
		err = json.Unmarshal(read, &data)
	}
	for _, r2 := range temp {
		isNeedInsert := true
		for i, r := range data {
			if r.RegionId == r2.RegionId {
				isNeedInsert = false
				if !utils.IsStringIn(platform, data[i].Platform) {
					data[i].Platform = append(data[i].Platform, platform)
				}
				break
			}
		}
		if isNeedInsert {
			r2.Platform = append(r2.Platform, platform)
			data = append(data, r2)
		}
	}
	content, _ := json.Marshal(data)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("文件打开/创建失败,原因是:", err)
		return
	}
	file.WriteString(string(content))
	file.Close()
}
