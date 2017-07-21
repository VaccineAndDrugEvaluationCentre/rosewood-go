package carpenter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Settings struct {
	items map[string]interface{}
}

//todo: add support for loadDefaultSettings call back
func NewSettings() *Settings {
	s := Settings{}
	s.items = make(map[string]interface{}, 10) //initial capacity of 10
	return &s
}

func (s *Settings) Get(key string) (value interface{}, ok bool) {
	value, ok = s.items[key]
	return
}

func (s *Settings) Set(key string, value interface{}) interface{} {
	oldValue, _ := s.items[key]
	s.items[key] = value
	return oldValue
}

//todo: change path to io.reader
func (s *Settings) LoadSettings(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to load settings: %v", err)
	}
	err = json.Unmarshal(file, &s.items)
	if err != nil {
		return fmt.Errorf("failed to parse settings: %v", err)
	}
	return nil
}

//todo: change path to io.writer
func (s *Settings) SaveSettings(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to save settings: %v", err)
	}
	defer file.Close()
	e := json.NewEncoder(file)
	if err := e.Encode(s.items); err != nil {
		return fmt.Errorf("failed to save settings: %v", err)
	}
	return nil
}