package zabbix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	MediaTypeGetTemplate = `{
		"jsonrpc": "2.0",
		"method": "mediatype.get",
		"params": {
			"output": "extend",
			"filter":{
							"description": "%v"
			}
		},
		"auth": "%v",
		"id": %v
	}`
)

func (api *API) MediaTypeGet(name string) (map[string]interface{}, error) {
	// fmt.Println(fmt.Sprintf(MediaTypeGetTemplate, name, api.Session, api.ID))
	payload := strings.NewReader(fmt.Sprintf(MediaTypeGetTemplate, name, api.Session, api.ID))
	req, err := http.NewRequest("POST", api.URL, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data := map[string]interface{}{}
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}
	if errmsg, ok := data["error"]; ok {
		return nil, fmt.Errorf("%v", errmsg)
	}
	if len(data["result"].([]interface{})) != 0 {
		return data["result"].([]interface{})[0].(map[string]interface{}), nil
	}
	return nil, nil
}
