/**
* @Author: mo tan
* @Description:
* @Date 2021/1/5 22:42
 */
package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/maotan/go-truffle/cloud"
	"github.com/maotan/go-truffle/cloud/serviceregistry"
	"github.com/maotan/go-truffle/feign"
	"github.com/maotan/go-truffle/logger"
	"github.com/maotan/go-truffle/truffle"
	"github.com/maotan/go-truffle/util"
	"github.com/maotan/go-truffle/yaml_config"
)

func ConsulInit(metaMap map[string]string) (serviceRegistry serviceregistry.ServiceRegistry, errRes error) {
	consulConf :=yaml_config.YamlConf.ConsulConf
	registryDiscoveryClient, err := serviceregistry.NewConsulServiceRegistry(consulConf.Host,
		consulConf.Port, consulConf.Token)
	feign.Init(registryDiscoveryClient)

	ip, err := util.GetLocalIP()
	if err != nil {
		panic(err)
	}
	serverConf := yaml_config.YamlConf.ServerConf
	si, _ := cloud.NewDefaultServiceInstance(serverConf.Name, ip,
		serverConf.Port, false, metaMap, "")
	registryDiscoveryClient.Register(si)
	return  registryDiscoveryClient, nil
}

func RouterInit(router *gin.Engine)  {
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("my-session", store))
	router.Use(logger.LogerMiddleware())
	router.Use(truffle.Recover)

	router.GET("/actuator/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
