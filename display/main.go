package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"rswAES256/client"
	"rswAES256/config"
	"strings"
)

var filePath = "." + string(filepath.Separator) + filepath.Join("test")

// OS 독립적 파일 경로
func main() {
	fmt.Println("Put correct key!!!. Files can not back to original files when you put wrong key!!!")
	fmt.Println("검증되지 않은 키를 입력 시 파일을 되돌릴 수 없습니다!!! 올바른 키를 입력하세요!!!")
	for {
		fmt.Print("Put key for decryption your files : ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text()) // 앞뒤 공백 제거
		if input == "" {
			fmt.Println("Error: 입력값이 비어있습니다!")
		} else {

			// 설정파일 설정.
			c := config.NewConfig("./config.toml")
			fmt.Println(c)

			// 입력한 키를 client에 넣어줌.
			newClient := client.NewClient(c)
			newClient.Key = input

			//복호화 진행
			if err := newClient.AESDecryptDirectory(filePath); err != nil {
				fmt.Println("failed to decrypt(복호화 실패) : ", err)
			}

			break
		}
	}
}
