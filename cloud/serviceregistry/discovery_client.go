package serviceregistry

import "github.com/maotan/go-truffle/cloud"

type DiscoveryClient interface {

	/**
	 * Gets all ServiceInstances associated with a particular serviceId.
	 * @param serviceId The serviceId to query.
	 * @return A List of ServiceInstance.
	 */
	GetInstances(serviceId string) ([]cloud.ServiceInstance, error)

	/**
	 * @return All known service IDs.
	 */
	GetServices() ([]string, error)

	/** 获取本地注册的服务*/
	//GetRegistryServices() map[string]map[string]cloud.ServiceInstance

	/** 获取所有注册的实例信息 */
	//GetAllRegistryInstances() (map[string][]cloud.ServiceInstance, error)
}
