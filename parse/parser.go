package parse

import (
	"fmt"
	"encoding/json"
)

// BaseParser demo
type BaseParser struct {
	name string
}

// Unmarshal 反序列化Page
func (p *BaseParser) Unmarshal(page string) (*Response, error) {
	var response = &Response{}
	err := json.Unmarshal([]byte(page), response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Parse parse func
func (p *BaseParser) Parse(resp *Response) {
	data, err := resp.ToMapSA()
	if err != nil {
		fmt.Println("Resp To Map Failed: " + err.Error())
		return
	}
	for k, v := range data {
		fmt.Println("----------------------------------")
		fmt.Println(k, ":", v)
		fmt.Println("----------------------------------")
	}
}

// NewParser new parser
func NewParser() Parser {
	return &BaseParser{}
}