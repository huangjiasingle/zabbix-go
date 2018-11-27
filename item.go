package zabbix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	itemGetTemplate = `{
		"jsonrpc": "2.0",
		"method": "item.get",
		"params": {
			"output": "extend",
			"filter": {
				"name": "%v",
				"hostid":"%v"
			},
			"sortfield": "name"
		},
		"auth": "%v",
		"id": %v
	}`

	itemDeleteTemplate = `{
		"jsonrpc": "2.0",
		"method": "item.delete",
		"params": [
			%v
		],
		"auth": "%v",
		"id": %v
	}`
	itemPostTemplate = `{
		"jsonrpc": "2.0",
		"method": "item.create",
		"params": {
				"type": "3",
				"hostid": "%v",
				"name": "%v",
				"key_": "net.tcp.service[tcp,%v,%v]",
				"delay": "%v",
				"history": "90d",
				"trends": "365d",
				"value_type": "3"
			},
		"auth": "%v",
		"id": %v
	}`
)

func (api *API) ItemCreate(hostid, name, ip, interval string, port int32) error {
	fmt.Println(fmt.Sprintf(itemPostTemplate, hostid, name, ip, port, interval, api.Session, api.ID))
	payload := strings.NewReader(fmt.Sprintf(itemPostTemplate, hostid, name, ip, port, interval, api.Session, api.ID))
	req, err := http.NewRequest("POST", api.URL, payload)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("zabbix api return response code %v", res.StatusCode)
	}
	result := map[string]interface{}{}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return err
	}
	if errmsg, ok := result["error"]; ok {
		return fmt.Errorf("%v", errmsg)
	}
	return nil
}

func (api *API) ItemDelete(ids []string) error {
	name := ""
	for _, id := range ids {
		name += fmt.Sprintf(`"%v"`, id)
	}
	fmt.Print(fmt.Sprintf(itemDeleteTemplate, name, api.Session, api.ID))
	payload := strings.NewReader(fmt.Sprintf(itemDeleteTemplate, name, api.Session, api.ID))
	req, err := http.NewRequest("POST", api.URL, payload)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	result := map[string]interface{}{}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return err
	}
	if errmsg, ok := result["error"]; ok {
		return fmt.Errorf("%v", errmsg)
	}
	return nil
}

func (api *API) ItemGet(name, hostid string) (map[string]interface{}, error) {
	payload := strings.NewReader(fmt.Sprintf(itemGetTemplate, name, hostid, api.Session, api.ID))
	req, err := http.NewRequest("POST", api.URL, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
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

func (api *API) ItemUpdate() error {
	return nil
}
