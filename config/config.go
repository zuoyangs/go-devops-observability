package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Jenkin_Instance struct {
	jenkinsURL string `yaml:"jenkinsURL"`
	username   string `yaml:"username"`
	password   string `yaml:"password"`
}

type YamlConfig struct {
	test2_jenkins Jenkin_Instance `yaml:"Jenkin_Instance"`
}

var yamlConfig YamlConfig

// 动态监听yaml配置文件
func WatchConfigYaml() {
	v := viper.New()
	v.SetConfigFile("./etc/config.yaml")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("配置文件不存在:%v\n", err)
		} else {
			fmt.Printf("配置文件存在,解析失败:%v\n", err)
		}
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		// 配置动态改变时，回调函数
		fmt.Printf("配置发生改变,重新解析配置: %v \n", e.Name)
		if err := v.Unmarshal(&yamlConfig); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&yamlConfig); err != nil {
		fmt.Println(err)
	}
}
