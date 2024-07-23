package main

import (
	"fmt"
	"github.com/sahithibanda01/auto-detection/autodetect"
)

func main() {
	err := autodetect.DetectDirectoriesToCache()
	if err != nil {
		fmt.Errorf("unable to detect or inject the build configurations: ", err)
	}
}
