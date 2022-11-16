package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func Printf(data any) {
	bs, _ := json.Marshal(data)
	var out bytes.Buffer
	json.Indent(&out, bs, "", "\t")
	fmt.Printf("%v\n", out.String())
}
