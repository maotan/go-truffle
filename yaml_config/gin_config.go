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

type GinConfig struct {
	GinActive  string     `mapstructure:"gin-active" json:"ginActive" yaml:"gin-active"`
}

type BootstrapConfig struct{
	GinConf GinConfig `mapstructure:"gin" json:"gin" yaml:"gin"`
}

func (this *BootstrapConfig) DefaultGinConfig() {
	ginConf := GinConfig{GinActive: "local"}
	this.GinConf = ginConf
}

func (this *BootstrapConfig) InitGinConfig(path string) {
	this.DefaultGinConfig()
	file, _ := ioutil.ReadFile(path)
	if err := yaml.Unmarshal(file, this); err != nil {
		panic(err)
	}
}