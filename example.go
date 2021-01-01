/**
* @Author: mo tan
* @Description:
* @Date 2021/1/1 21:58
 */
package main

import (
	"github.com/maotan/go-truffle/cloud"
	"github.com/maotan/go-truffle/cloud/serviceregistry"
	"github.com/maotan/go-truffle/feign"
	"github.com/maotan/go-truffle/routes"
	"github.com/maotan/go-truffle/util"
	"math/rand"
	"time"
)

func main() {

	host := "127.0.0.1"
	port := 8500
	token := ""
	registryDiscoveryClient, err := serviceregistry.NewConsulServiceRegistry(host, port, token)
	feign.Init(registryDiscoveryClient)

	ip, err := util.GetLocalIP()
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())

	si, _ := cloud.NewDefaultServiceInstance("go-user-server", ip, 5000,
		false, map[string]string{"user": "zyn2"}, "")

	registryDiscoveryClient.Register(si)

	//logger.ConfigLocalFileLogger()
	err = routes.Run()
	if err != nil{
		registryDiscoveryClient.Deregister()
	}
}