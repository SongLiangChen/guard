package config

import (
	"io/ioutil"

	"gopkg.in/ini.v1"
)

var instance *ini.File

func InitConfig() (err error) {
	var (
		sources  = []string{}
		basePath = "./data/ini/"
		fs, _    = ioutil.ReadDir(basePath)
	)

	for i := 0; i < len(fs); i++ {
		sources = append(sources, basePath+"/"+fs[i].Name())
	}

	others := make([]interface{}, 0)
	for i := 1; i < len(sources); i++ {
		others = append(others, sources[i])
	}
	if instance, err = ini.Load(sources[0], others...); err != nil {
		return err
	}

	return nil
}

func GetConfig() *ini.File {
	return instance
}
