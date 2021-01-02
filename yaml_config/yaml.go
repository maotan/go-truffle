/**
* @Author: mo tan
* @Description:
* @Date 2021/1/2 15:04
 */
package yaml_config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path"
)
type LogConfig struct {
	LogPath  string      `mapstructure:"log-path" json:"logPath" yaml:"log-path"`
	RotationCount uint   `mapstructure:"rotation-count" json:"rotationCount" yaml:"rotation-count"`
}

type GinConfig struct {
	Port int
	Name string
}

type ConsulConfig struct {
	Host string
	Port int
	Token string
}

type YamlConfig struct {
	//存放各个配置文件的路径 Path
	LogConf LogConfig `mapstructure:"log" json:"log" yaml:"log"`
	GinConf GinConfig `mapstructure:"gin" json:"gin" yaml:"gin"`
	ConsulConf ConsulConfig `mapstructure:"consul" json:"consul" yaml:"consul"`
}

var YamlConf YamlConfig

func init()  {
	NewYamlConfig("config/")
}

func NewYamlConfig (configPath string) YamlConfig{
	//var config YamlConfig
	ginConfigPath := path.Join(configPath, "config.yaml")
	var ginActiveConf GinActiveConfig;
	ginActiveConf.InitGinConfig(ginConfigPath)
	if ginActiveConf.ActiveConf.Active == ""{
		panic("未配置config.yaml")
	}
	configFileName := fmt.Sprintf("%s%s%s", "config-", ginActiveConf.ActiveConf.Active, ".yaml")
	fullPath := path.Join(configPath, configFileName)

	file, _ := ioutil.ReadFile(fullPath)
	if err := yaml.Unmarshal(file, &YamlConf); err != nil {
		panic(err)
	}
	return YamlConf
}

