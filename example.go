package main

import (
	"fmt"
	"github.com/maotan/go-truffle/cloud"
	"github.com/maotan/go-truffle/cloud/serviceregistry"
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

	ip, err := util.GetLocalIP()
	if err != nil {
		//t.Error(err)
		panic(err)
	}

	fmt.Println(ip)
	rand.Seed(time.Now().UnixNano())

	si, _ := cloud.NewDefaultServiceInstance("go-user-server", "", 8010,
		false, map[string]string{"user": "zyn2"}, "")

	registryDiscoveryClient.Register(si)

	routes.Run()
}