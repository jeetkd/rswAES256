package main

import (
	"flag"
	"fmt"
	"log"
	"rswAES256/client"
	"rswAES256/config"
	"rswAES256/decrypt"
	"rswAES256/encrypt"
)

var path = "./config.toml"

var pathFlag = flag.String("config", path, "set toml path")

func main() {
	flag.Parse()
	data := []byte("1234567890123456") // 암호화할 데이터 (128비트)
	//key := []byte("0123456789abcdef")  // 암호 키 (128비트)

	// 설정파일 설정.
	c := config.NewConfig(path)
	fmt.Println(c)

	// 서버로부터 key를 가져옴.
	newClient := client.NewClient(c)
	key, err := newClient.GetKey()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Key from server : ", *key)

	// AES 암호화
	cipherText, err := encrypt.EncryptAES(data, []byte(*key))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("cipherText : %x\n", cipherText)

	// 파일 암호화
	err = encrypt.EncryptFile("test.txt", []byte(*key))
	if err != nil {
		log.Fatalln(err)
	}

	//AES 복호화
	plainText, err := decrypt.DecryptAES(cipherText, []byte(*key))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("plainText : ", string(plainText))

}
