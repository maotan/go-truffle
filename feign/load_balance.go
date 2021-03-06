/**
* @Author: mo tan
* @Description:
* @Date 2021/1/1 21:58
 */
package feign

import (
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
)

// Feign load balance client
type Lbc struct {
	feign *Feign
	app   string
	index uint32

	// resty.Client
	client *resty.Client
}

// pick a server to send request
func (t *Lbc) pickUrl() string {
	urls, ok := t.feign.GetAppUrls(t.app)
	if !ok || len(urls) == 0 {
		log.Infof("Failed to pick server, reason: no available urls for app=%s", t.app)

		// no need to panic
		// coz return empty string won't panic while calls like "DefaultFeign.App("APP_NAME").R().Post()"
		// it will only get request failed
		return ""
	}

	idx := atomic.AddUint32(t.feign.appNextUrlIndex[t.app], 1)
	idx %= uint32(len(urls))
	t.index = idx
	atomic.CompareAndSwapUint32(t.feign.appNextUrlIndex[t.app], uint32(len(urls)), 0)

	return urls[idx]
}

// pick target url and new resty.Client
func (t *Lbc) pick() *Lbc {
	t.client = resty.New()
	t.client.HostURL = t.pickUrl()
	log.Infof("Picked url=%s", t.client.HostURL)
	return t
}
