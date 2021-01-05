/**
* @Author: mo tan
* @Description:
* @Date 2021/1/5 22:42
 */
package web

import (
	"github.com/maotan/go-truffle/cloud"
	"github.com/maotan/go-truffle/cloud/serviceregistry"
	"github.com/maotan/go-truffle/feign"
	"github.com/maotan/go-truffle/util"
	"github.com/maotan/go-truffle/yaml_config"
)

func WebInit() (errRes error, serviceRegistry serviceregistry.ServiceRegistry) {
	consulConf :=yaml_config.YamlConf.ConsulConf
	registryDiscoveryClient, err := serviceregistry.NewConsulServiceRegistry(consulConf.Host,
		consulConf.Port, consulConf.Token)
	feign.Init(registryDiscoveryClient)

	ip, err := util.GetLocalIP()
	if err != nil {
		panic(err)
	}
	serverConf := yaml_config.YamlConf.ServerConf
	si, _ := cloud.NewDefaultServiceInstance(serverConf.Name, ip, serverConf.Port,
		false, map[string]string{"user": "zyn2"}, "")
	registryDiscoveryClient.Register(si)
	return nil, registryDiscoveryClient
}
