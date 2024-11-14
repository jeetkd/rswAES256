package client

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func (c *Client) GetKey() (*string, error) {
	res, err := http.Get(c.config.Network.Uri)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// body 닫기
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Println("Error closing response body:", err)
		}
	}(res.Body)

	// body 읽어옴
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// body에 있는 데이터를 client에 넣어줌.
	err = json.Unmarshal(body, &c)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &c.Key, nil
}
