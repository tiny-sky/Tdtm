package main

import (
	"fmt"
	"io"
	"os"

	"github.com/tiny-sky/Tdtm/conf"
	"gopkg.in/yaml.v3"
)

var setting conf.Settings

func main() {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	file, err := os.Open(dir + "/conf.yml")
	if err != nil {
		return
	}
	defer file.Close()

	byteAll, err := io.ReadAll(file)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(byteAll, &setting)
	if err != nil {
		fmt.Printf("error %s", err)
	}
	fmt.Printf("%+v", setting)
}
