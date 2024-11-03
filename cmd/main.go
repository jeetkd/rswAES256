package main

import (
	"flag"
	"fmt"
	"rswAES256/config"
)

var path = "./config.toml"

var pathFlag = flag.String("config", path, "set toml path")

func main() {
	flag.Parse()

	// 설정파일 설정.
	c := config.NewConfig(path)
	fmt.Println(c)
}
