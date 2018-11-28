package zabbix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	UserLoginTemplate     = `{"jsonrpc":"2.0","method":"user.login","params":{"user":"%v","password":"%v","userData":%v},"id":%v}`
	HostGroupGetTemplate  = `{"jsonrpc":"2.0","method":"hostgroup.get","params":{"output":"extend","filter":{"name":"%v"}},"auth":"%v","id":%v}`
	HostGroupPostTemplate = `{"jsonrpc":"2.0","method":"hostgroup.create","params":{"name":"%v"},"auth":"%v","id":%v}`
	TemplatePostTemplate  = `{"jsonrpc":"2.0","method":"template.create","params":{"host":"%v","groups":{"groupid":%v}},"auth":"%v","id":%v}`

	templateGetTemplate = `{
		"jsonrpc": "2.0",
		"method": "template.get",
		"params": {
			"output": "extend",
			"filter": {
				"host": "%v"
			},
			"sortfield": "name"
		},
		"auth": "%v",
		"id": %v
	}`
)

func (api *API) TemplateCreate(groupid, name string) error {
	payload := strings.NewReader(fmt.Sprintf(TemplatePostTemplate, name, groupid, api.Session, api.ID))
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

func (api *API) TemplateDelete() error {
	return nil
}

func (api *API) TemplateGet(name string) (map[string]interface{}, error) {
	payload := strings.NewReader(fmt.Sprintf(templateGetTemplate, name, api.Session, api.ID))
	req, err := http.NewRequest("POST", api.URL, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
GET:
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data := map[string]interface{}{}
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}
	if len(data["result"].([]interface{})) == 0 {
		result, err := api.HostGroupGet("Pass Middle Ware")
		if err != nil {
			return nil, err
		}
		if result != nil && result["groupid"] != nil {
			err = api.TemplateCreate(result["groupid"].(string), name)
			if err != nil {
				return nil, err
			}
			goto GET
		} else {
			err = api.HostGroupCreate("Pass Middle Ware")
			if err != nil {
				return nil, err
			}
			goto GET
		}
	}

	if len(data["result"].([]interface{})) != 0 {
		return data["result"].([]interface{})[0].(map[string]interface{}), nil
	}
	return nil, nil
}

func (api *API) TemplateUpdate() error {
	return nil
}
