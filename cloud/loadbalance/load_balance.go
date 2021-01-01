package loadbalance

import (
	"errors"
	"github.com/maotan/go-truffle/cloud"
	"math/rand"
	"time"
)

type LoadBalance interface {
	choose(serviceId string, instances []cloud.ServiceInstance) (cloud.ServiceInstance, error)
}

type FirstLoadBalance struct {
}

func (f FirstLoadBalance) choose(serviceId string, instances []cloud.ServiceInstance) (cloud.ServiceInstance, error) {

	if instances != nil &&  len(instances) > 0{
		length := len(instances)
		rand.Seed(time.Now().Unix())
		index := rand.Intn(length)
		return instances[index], nil
	}
	return nil, errors.New("no available instance")
}
