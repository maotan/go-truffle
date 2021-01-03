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

type ServerConfig struct {
	Port int
	Name string
}

type ConsulConfig struct {
	Host string
	Port int
	Token string
}

type DatabaseConfig struct {
	Name string		`mapstructure:"name" json:"name" yaml:"name"`
	Host string		`mapstructure:"host" json:"host" yaml:"host"`
	Port string		`mapstructure:"port" json:"port" yaml:"port"`
	Username string 	`mapstructure:"username" json:"username" yaml:"username"`
	Password string		`mapstructure:"password" json:"password" yaml:"password"`
}

type YamlConfig struct {
	//存放各个配置文件的路径 Path
	LogConf LogConfig `mapstructure:"log" json:"log" yaml:"log"`
	ServerConf ServerConfig `mapstructure:"server" json:"server" yaml:"server"`
	ConsulConf ConsulConfig `mapstructure:"consul" json:"consul" yaml:"consul"`
	DatabaseConf DatabaseConfig `mapstructure:"database" json:"database" yaml:"database"`
}

var YamlConf YamlConfig

func init()  {
	NewYamlConfig("config/")
}

func NewYamlConfig (configPath string) YamlConfig{
	//var config YamlConfig
	ginConfigPath := path.Join(configPath, "config.yaml")
	var ginConf GinConfig;
	ginConf.InitGinConfig(ginConfigPath)
	if ginConf.ActiveConf.Active == ""{
		panic("未配置config.yaml")
	}
	configFileName := fmt.Sprintf("%s%s%s", "config-", ginConf.ActiveConf.Active, ".yaml")
	fullPath := path.Join(configPath, configFileName)

	file, _ := ioutil.ReadFile(fullPath)
	if err := yaml.Unmarshal(file, &YamlConf); err != nil {
		panic(err)
	}
	return YamlConf
}

