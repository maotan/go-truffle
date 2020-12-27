package feign

import (
	"fmt"
	"github.com/maotan/go-truffle/cloud"
	"github.com/maotan/go-truffle/cloud/serviceregistry"
	"gopkg.in/resty.v1"
	"log"
	"net/url"
	"sync"
	"time"
)

const (
	DEFAULT_REFRESH_APP_URLS_INTERVALS = 30
)

var DefaultFeign = &Feign{
	appUrls:         make(map[string][]string),
	appNextUrlIndex: make(map[string]*uint32),
}

func Init(discoveryClient serviceregistry.DiscoveryClient){
	DefaultFeign.discoveryClient = discoveryClient
}


type Feign struct {
	// Discovery client to get Apps and instances
	discoveryClient serviceregistry.DiscoveryClient

	// assign app => urls
	appUrls map[string][]string

	// Counter for calculate next url'index
	appNextUrlIndex map[string]*uint32

	// seconds of updating app's urls periodically
	refreshAppUrlsIntervals int

	// ensure some daemon task only run one time
	once sync.Once

	mu sync.RWMutex
}

// use discovery client to get all registry app => instances
func (t *Feign) UseDiscoveryClient(client serviceregistry.DiscoveryClient) *Feign {
	t.discoveryClient = client
	return t
}

// assign static app => urls
func (t *Feign) UseUrls(appUrls map[string][]string) *Feign {
	t.mu.Lock()
	defer t.mu.Unlock()

	//v := uint32(time.Now().UnixNano())
	//appNextUrlIndex[t.app] = &v
	for app, urls := range appUrls {

		// reset app'urls
		tmpAppUrls := make([]string, 0)
		for _, u := range urls {
			_, err := url.Parse(u)
			if err != nil {
				log.Print("Invalid url=%s, parse err=%s", u, err.Error())
				continue
			}

			tmpAppUrls = append(tmpAppUrls, u)
		}

		if len(tmpAppUrls) == 0 {
			log.Print("Empty valid urls for app=%s, skip to set app's urls", app)
			continue
		}

		t.appUrls[app] = tmpAppUrls
		if t.appNextUrlIndex[app] == nil {
			v := uint32(time.Now().UnixNano())
			t.appNextUrlIndex[app] = &v
		}
	}

	return t
}

func (t *Feign) SetRefreshAppUrlsIntervals(intervals int) {
	t.refreshAppUrlsIntervals = intervals
}

// return resty.Client
func (t *Feign) App(app string) *resty.Client {
	defer func() {
		if err := recover(); err != nil {
			log.Print("App(%s) catch panic err=%v", app, err)
		}
	}()

	// daemon to update app urls periodically
	// only execute once globally
	t.once.Do(func() {
		if t.discoveryClient == nil {
			log.Print("no discovery client, no need to update appUrls periodically.")
			return
		}
		instances, err := t.discoveryClient.GetInstances(app)
		if err == nil || len(instances)==0{
			log.Print("no discovery client, no need to update appUrls periodically.")
			return
		}

		t.updateAppUrlsIntervals(instances)
	})

	// try update app's urls.
	// if app's urls is exist, do nothing
	t.tryRefreshAppUrls(app)

	lbc := &Lbc{
		feign: t,
		app:   app,
	}
	return lbc.pick().client
}

// try update app's urls
// if app's urls is exist, do nothing
func (t *Feign) tryRefreshAppUrls(app string) {
	if _, ok := t.GetAppUrls(app); ok {
		return
	}

	if t.discoveryClient == nil  {
		log.Print("no discovery client, no need to update app'urls.")
		return
	}
	instances,err := t.discoveryClient.GetServices()
	if (err != nil || len(instances) == 0){
		log.Print("no discovery client, no need to update app'urls.")
		return
	}

	t.updateAppUrls()
}

// update app urls periodically
func (t *Feign) updateAppUrlsIntervals(instances []cloud.ServiceInstance) {
	if t.refreshAppUrlsIntervals <= 0 {
		t.refreshAppUrlsIntervals = DEFAULT_REFRESH_APP_URLS_INTERVALS
	}

	go func() {
		for {
			t.updateAppUrls()

			time.Sleep(time.Second * time.Duration(t.refreshAppUrlsIntervals))
			log.Print("Update app urls interval...ok")
			for app, urls := range t.appUrls {
				log.Print("app=> %s, urls => %v", app, urls)
			}

		}
	}()
}

// Update app urls from registry apps
func (t *Feign) updateAppUrls() {
	registryServices, _ := t.discoveryClient.GetRegistryServices()
	tmpAppUrls := make(map[string][]string)

	for serviceId, serviceInstMap := range registryServices {
		var isAppAlreadyExist bool
		var curAppUrls []string
		var isUpdate bool

		// if app is already exist in t.appUrls, check whether app's urls are updated.
		// if app's urls are updated, t.appUrls
		if curAppUrls, isAppAlreadyExist = t.GetAppUrls(serviceId); isAppAlreadyExist {
			for _, inst := range serviceInstMap {
				isExist := false
				for _, v := range curAppUrls {
					//insHomePageUrl := strings.TrimRight(inst.GetHost(), "/")
					insHomePageUrl := fmt.Sprintf("%s:%d/", inst.GetHost(), inst.GetPort())
					if v == insHomePageUrl {
						isExist = true
						break
					}
				}

				if !isExist {
					isUpdate = true
					break
				}
			}
		}

		// app are not exist in t.appUrls or app's urls has been update
		if !isAppAlreadyExist || isUpdate {
			tmpAppUrls[serviceId] = make([]string, 0)
			for _, insVo := range serviceInstMap {
				insHomePageUrl := fmt.Sprintf("%s:%d/", insVo.GetHost(), insVo.GetPort())
				tmpAppUrls[serviceId] = append(tmpAppUrls[serviceId], insHomePageUrl)
			}
		}
	}

	// update app's urls to feign
	t.UseUrls(tmpAppUrls)
}

// get app's urls
func (t *Feign) GetAppUrls(app string) ([]string, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if _, ok := t.appUrls[app]; !ok {
		return nil, false
	}

	return t.appUrls[app], true
}