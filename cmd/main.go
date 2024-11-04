package main

import (
	"flag"
	"fmt"
	"log"
	"rswAES256/config"
	"rswAES256/encrypt"
)

var path = "./config.toml"

var pathFlag = flag.String("config", path, "set toml path")

func main() {
	flag.Parse()
	data := []byte("1234567890123456") // 128비트
	key := []byte("0123456789abcdef")  // 128비트

	// 설정파일 설정.
	c := config.NewConfig(path)
	fmt.Println(c)

	// AES 암호화
	cipherText, err := encrypt.EncryptAES(data, key)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("cipherText : %x\n", cipherText)
}
