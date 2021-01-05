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
	"os"
	"path"
)

var YamlConf YamlConfig
var configDir = "config/"

func init()  {
	err:=GetConfig(&YamlConf)
	if err!=nil{
		panic(err)
	}
}

func getConfigActive(filePath string, out interface{}) (errRes error){
	_, err := os.Stat(filePath)
	if err != nil{
		return err
	}
	file, _ := ioutil.ReadFile(filePath)
	err = yaml.Unmarshal(file, out)
	return err
}

func GetConfig(out interface{}) (errRes error) {
	activeConfigPath := path.Join(configDir, "config.yaml")
	var ginConf GinConfig;
	err := getConfigActive(activeConfigPath, &ginConf)
	if err != nil {
		panic(err)
	}
	if (ginConf.ActiveConf.Active == "") {
		panic("config.yaml未配置正确")
	}

	configFileName := fmt.Sprintf("%s%s%s", "config-", ginConf.ActiveConf.Active, ".yaml")
	configFilePath := path.Join(configDir, configFileName)
	_, ef := os.Stat(configFilePath)
	if ef != nil {
		panic(ef)
	}
	file, _ := ioutil.ReadFile(configFilePath)
	err = yaml.Unmarshal(file, out)
	return err
}
