/**
* @Author: mo tan
* @Description:
* @Date 2021/1/1 21:58
 */
package feign

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/maotan/go-truffle/cloud/serviceregistry"
	log "github.com/sirupsen/logrus"
	"net/url"
	"sync"
	"time"
)

const (
	DEFAULT_REFRESH_APP_URLS_INTERVALS = 120
)

var DefaultFeign = &Feign{
	appUrlMap:       make(map[string][]string),
	appNextUrlIndex: make(map[string]*uint32),
}

func Init(discoveryClient serviceregistry.DiscoveryClient){
	DefaultFeign.discoveryClient = discoveryClient
}

type Feign struct {
	// Discovery client to get Apps and instances
	discoveryClient serviceregistry.DiscoveryClient

	// assign app => urls
	appUrlMap map[string][]string

	// Counter for calculate next url'index
	appNextUrlIndex map[string]*uint32

	// seconds of updating app's urls periodically
	refreshAppUrlsIntervals int

	// ensure some daemon task only run one time
	once sync.Once

	mu sync.RWMutex
}

func (t *Feign) SetRefreshAppUrlsIntervals(intervals int) {
	t.refreshAppUrlsIntervals = intervals
}

// return resty.Client
func (t *Feign) App(app string) *resty.Client {
	defer func() {
		if err := recover(); err != nil {
			info := fmt.Sprintf("App(%s) catch panic err=%v", app, err)
			log.Info(info)
		}
	}()

	// daemon to update app urls periodically
	// only execute once globally
	t.once.Do(func() {
		if t.discoveryClient == nil {
			log.Info("no discovery client, no need to update appUrls periodically.")
			return
		}
		t.updateAppUrlsIntervals()
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
func (t *Feign) tryRefreshAppUrls(serviceId string) {
	if _, ok := t.GetAppUrls(serviceId); ok {
		return
	}

	if t.discoveryClient == nil  {
		log.Info("no discovery client, no need to update app'urls.")
		return
	}
	services,err := t.discoveryClient.GetServices()
	if (err != nil || len(services) == 0){
		log.Info("no discovery client, no need to update app'urls.")
		return
	}

	t.updateAppUrls(serviceId)
}

// update app urls periodically
func (t *Feign) updateAppUrlsIntervals() {
	if t.refreshAppUrlsIntervals <= 0 {
		t.refreshAppUrlsIntervals = DEFAULT_REFRESH_APP_URLS_INTERVALS
	}

	go func() {
		for {
			serviceArray, err := t.discoveryClient.GetServices()
			if err != nil || len(serviceArray)==0{
				log.Info("no discovery client, no need to update appUrls periodically.")
				return
			}
			currentServiceMap := make(map[string]bool)
			for _, serviceId := range(serviceArray){
				currentServiceMap[serviceId] = true
				t.updateAppUrls(serviceId)
			}
			t.DeleteUrls(currentServiceMap)

			time.Sleep(time.Second * time.Duration(t.refreshAppUrlsIntervals))
			log.Info("Update app urls interval...ok")
			for app, urls := range t.appUrlMap {
				info := fmt.Sprintf("app=> %s, urls => %v", app, urls)
				log.Info(info)
			}

		}
	}()
}

// Update app urls from registry apps
func (t *Feign) updateAppUrls(serviceId string) {
	instanceArray, error := t.discoveryClient.GetInstances(serviceId)
	if error != nil {
		return
	}

	var isAppAlreadyExist bool
	var curAppUrls []string
	var isUpdate bool
	// if app is already exist in t.appUrls, check whether app's urls are updated.
	// if app's urls are updated, t.appUrls
	curAppUrls, isAppAlreadyExist = t.GetAppUrls(serviceId)
	if !isAppAlreadyExist {
		if instanceArray == nil || len(instanceArray) == 0{
			return ;
		}
	} else {
		if len(curAppUrls) != len(instanceArray){
			isUpdate = true
		} else {
			for _, inst := range instanceArray {
				isExist := false
				for _, v := range curAppUrls {
					insHomePageUrl := fmt.Sprintf("http://%s:%d", inst.GetHost(), inst.GetPort())
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

	}

	// app are not exist in t.appUrls or app's urls has been update
	if !isAppAlreadyExist || isUpdate {
		tmpAppUrls := make([]string, 0)
		for _, insVo := range instanceArray {
			insHomePageUrl := fmt.Sprintf("http://%s:%d", insVo.GetHost(), insVo.GetPort())
			tmpAppUrls = append(tmpAppUrls, insHomePageUrl)
		}
		// update app's urls to feign
		t.UseUrls(serviceId, tmpAppUrls)
	}
}

// assign static app => urls
func (t *Feign) UseUrls(serviceId string, appUrls []string) *Feign {
	t.mu.Lock()
	defer t.mu.Unlock()

	//v := uint32(time.Now().UnixNano())
	tmpAppUrls := make([]string, 0)
	for _, appUrl := range appUrls {
		// reset app'urls
		_, err := url.Parse(appUrl)
		if err != nil {
			msg := fmt.Sprintf("Invalid url=%s, parse err=%s", appUrl, err.Error())
			log.Info(msg)
			continue
		}
		tmpAppUrls = append(tmpAppUrls, appUrl)

		if len(tmpAppUrls) == 0 {
			info := fmt.Sprintf("Empty valid urls for app=%s, skip to set app's urls", serviceId)
			log.Info(info)
			continue
		}

		t.appUrlMap[serviceId] = tmpAppUrls
		if t.appNextUrlIndex[serviceId] == nil {
			v := uint32(time.Now().UnixNano())
			t.appNextUrlIndex[serviceId] = &v
		}
	}
	return t
}
// delete unregistered service
func (t *Feign) DeleteUrls(currentMap map[string]bool) *Feign {
	t.mu.Lock()
	defer t.mu.Unlock()
	for serviceId, _ := range t.appUrlMap{
		_, ok := currentMap[serviceId]
		if !ok{
			log.Infof("feign remove service %s", serviceId)
			delete(t.appUrlMap, serviceId)
		}
	}
	return t
}

// get app's urls
func (t *Feign) GetAppUrls(app string) ([]string, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if _, ok := t.appUrlMap[app]; !ok {
		return nil, false
	}

	return t.appUrlMap[app], true
}