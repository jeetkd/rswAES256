package readme

import "os"

var s int

func CreateFileReadme(filename string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return file, nil
}

func WriteFileReadme(file *os.File, message string) error {
	_, err := file.Write([]byte(message))
	if err != nil {
		return err
	}
	return nil
}
