/**
* @Author: mo tan
* @Description:
* @Date 2021/1/1 21:58
 */
package truffle

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ParseResponse(response *http.Response) (map[string]interface{}, error){
	var result map[string]interface{}
	body,err := ioutil.ReadAll(response.Body)
	if err == nil {
		err = json.Unmarshal(body, &result)
	}

	return result,err
}
