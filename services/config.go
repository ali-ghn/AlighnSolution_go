package services

import (
	"encoding/json"
	"os"
)

type IConfigHelper interface {
	GetSection(key string) (interface{}, error)
}

type ConfigHelper struct {
	Filename string
	Config   string
}

func NewConfigHelper(filename string) (*ConfigHelper, error) {
	if filename == "" {
		filename = "appSettings.json"
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &ConfigHelper{
		Filename: filename,
		Config:   string(content),
	}, nil
}

func (ch ConfigHelper) GetSection(key string) (interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(ch.Config), &data)
	if err != nil {
		return nil, err
	}
	return data[key], nil
}
