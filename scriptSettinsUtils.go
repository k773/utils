package utils

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
)

type SettingsManager struct {
	SettingsFilePath string
	Settings         interface{}
}

func NewScriptSettingsManager(settingsFilePath string, settings interface{}) SettingsManager {
	return SettingsManager{
		SettingsFilePath: settingsFilePath,
		Settings:         settings,
	}
}

func (sm *SettingsManager) LoadAsJson() (err error) {
	var fileData []byte

	if fileData, err = ioutil.ReadFile(sm.SettingsFilePath); err == nil {
		err = json.Unmarshal(fileData, &sm.Settings)
	}
	return err
}

func (sm *SettingsManager) SaveAsJson() (err error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)

	if err = encoder.Encode(sm.Settings); err == nil {
		err = ioutil.WriteFile(sm.SettingsFilePath, buffer.Bytes(), 0600)
	}
	return err
}

func (sm *SettingsManager) LoadAsXml() (err error) {
	var fileData []byte

	if fileData, err = ioutil.ReadFile(sm.SettingsFilePath); err == nil {
		err = xml.Unmarshal(fileData, &sm.Settings)
	}
	return err
}

func (sm *SettingsManager) SaveAsXml() (err error) {
	var marshaled []byte

	if marshaled, err = xml.MarshalIndent(sm.Settings, "", "    "); err == nil {
		err = ioutil.WriteFile(sm.SettingsFilePath, marshaled, 0600)
	}
	return err
}
