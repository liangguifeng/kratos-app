package config

import (
	"fmt"
	"github.com/liangguifeng/kratos-app/internal/setup"
)

// MapConfig 加载配置到结构体中
func MapConfig(section string, v interface{}, isListenChange bool) error {
	return mapNacosConfig(section, v, isListenChange)
}

// mapNacosConfig
func mapNacosConfig(section string, v interface{}, isListenChange bool) error {
	var uniqueSectionDict = make(map[string]bool)
	if _, ok := uniqueSectionDict[section]; ok {
		return fmt.Errorf("repeate section config")
	}
	uniqueSectionDict[section] = true

	fields, err := setup.GetReflectFields(section, v)
	if err != nil {
		return err
	}
	if len(fields) == 0 {
		return nil
	}

	err = setup.SaveWatchConfigField(v, fields)
	if err != nil {
		return fmt.Errorf("save err: %v", err)
	}

	if isListenChange {
		for nacosKeyName, field := range fields {
			setup.WatchConfigFields[nacosKeyName] = field
		}
	}

	return nil
}
