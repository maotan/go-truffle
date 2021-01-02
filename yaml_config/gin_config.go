/**
* @Author: mo tan
* @Description:
* @Date 2021/1/2 14:39
 */
package yaml_config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type ActiveConfig struct {
	Active  string     `mapstructure:"active" json:"active" yaml:"active"`
}

type GinActiveConfig struct{
	ActiveConf ActiveConfig `mapstructure:"gin" json:"gin" yaml:"gin"`
}

func (this *GinActiveConfig) DefaultGinConfig() {
	activeConf := ActiveConfig{Active: "local"}
	this.ActiveConf = activeConf
}

func (this *GinActiveConfig) InitGinConfig(path string) {
	this.DefaultGinConfig()
	file, _ := ioutil.ReadFile(path)
	if err := yaml.Unmarshal(file, this); err != nil {
		panic(err)
	}
}