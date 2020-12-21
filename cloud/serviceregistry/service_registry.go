package serviceregistry

import "github.com/maotan/go-truffle/cloud"

type ServiceRegistry interface {
	Register(serviceInstance cloud.ServiceInstance) bool

	Deregister()
}
