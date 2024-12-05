package readme

import (
	"bufio"
	"os"
)

// CreateFileReadme 는 readme 파일 생성
func CreateFileReadme(filename, message string) error {
	var file *os.File
	var err error

	if file, err = os.Create(filename); err != nil {
		return err
	} else if err = WriteFileReadme(file, message); err != nil {
		defer file.Close()
		return err
	} else {
		return nil
	}
}

// WriteFileReadme 는 readme 파일 작성
func WriteFileReadme(file *os.File, message string) error {
	_, err := file.Write([]byte(message))
	if err != nil {
		return err
	}
	return nil
}

// OpenFileReadme 파일을 수정 또는 생성 합니다.
func OpenFileReadme(filename, messages string) error {
	file, err := os.OpenFile("filename.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(messages)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}
