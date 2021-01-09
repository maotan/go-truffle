/**
* @Author: mo tan
* @Description:
* @Date 2021/1/5 22:42
 */
package web

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" //mysql driver
	"github.com/ilibs/gosql/v2"
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
	store := cookie.NewStore([]byte("fds@dJD-0@"))
	router.Use(sessions.Sessions("my-session", store))
	router.Use(logger.LogerMiddleware())
	router.Use(truffle.Recover)

	router.GET("/actuator/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}

func DatabaseInit()  {
	databaseConf := yaml_config.YamlConf.DatabaseConf
	driver := databaseConf.Driver
	dsn := fmt.Sprintf("%s:%s@%s",
		databaseConf.Username, databaseConf.Password, databaseConf.Url)
	configs := make(map[string]*gosql.Config)
	configs["default"] = &gosql.Config{
		Enable:  true,
		Driver:  driver,
		Dsn:     dsn,
		ShowSql: true,
	}
	// connet database
	gosql.Connect(configs)
}
