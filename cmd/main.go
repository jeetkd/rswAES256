package main

import (
	"flag"
	"fmt"
	"log"
	"rswAES256/client"
	"rswAES256/config"
	"time"
)

var path = "./config.toml"

// var pathDebug = "/Users/jeongseong-won/GolandProjects/rswAES256/cmd/config.toml"
var pathFlag = flag.String("config", path, "set toml path")

func main() {
	flag.Parse()
	//data := []byte("1234567890123456") // 암호화할 데이터 (128비트)
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

	// rootPath의 하위 파일들을 암호화 함.
	err = newClient.AESEncryptDirectory("./test/")
	if err != nil {
		log.Fatalln(err)
	}

	time.Sleep(time.Second * 5)
	// rootPath의 하위 파일들을 복호화(확장자가.jkd으로 암호화된것만 복호화)
	err = newClient.AESDecryptDirectory("./test/")
	if err != nil {
		log.Fatalln(err)
	}

	// AES 암호화
	//cipherText, err := encrypt.EncryptAES(data, []byte(*key))
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Printf("cipherText : %x\n", cipherText)

	// 파일 암호화
	//err = encrypt.EncryptFile("./test.txt", []byte(*key))
	//if err != nil {
	//	log.Fatalln(err)
	//}

	// 파일 확장자 변경
	//err = newClient.ChangeFileExtension("./test.txt")
	//if err != nil {
	//	log.Fatalln(err)
	//}

	// 파일 복호화
	//err = decrypt.DecryptFile("./test/abc.go", []byte(*key))
	//if err != nil {
	//	log.Fatalln()
	//}

	//AES 복호화
	//plainText, err := decrypt.DecryptAES(cipherText, []byte(*key))
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Println("plainText : ", string(plainText))

}
