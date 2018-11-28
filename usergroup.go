package zabbix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	UserGroupGetTemplate = `{
		"jsonrpc": "2.0",
		"method": "usergroup.get",
		"params": {
			"output": "extend",
			"status": 0,
			"filter":{
				"name":"%v"
			}
		},
		"auth": "%v",
		"id": %v
	}`
	UserGroupPostTemplate = `{
		"jsonrpc": "2.0",
		"method": "usergroup.create",
		"params": {
			"name": "%v",
			"rights": {
				"permission": 3,
				"id": "2"
			}
		},
		"auth": "%v",
		"id": %v
	}`
)

func (api *API) UserGroupCreate(name string) error {
	fmt.Println(fmt.Sprintf(UserGroupPostTemplate, name, api.Session, api.ID))
	payload := strings.NewReader(fmt.Sprintf(UserGroupPostTemplate, name, api.Session, api.ID))
	req, err := http.NewRequest("POST", api.URL, payload)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

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

func (api *API) UserGroupGet(name string) (map[string]interface{}, error) {
	payload := strings.NewReader(fmt.Sprintf(UserGroupGetTemplate, name, api.Session, api.ID))
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
	data := map[string]interface{}{}
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	if len(data["result"].([]interface{})) == 0 {
		err = api.UserGroupCreate(name)
		if err != nil {
			return nil, err
		}
		goto GET
	}
	if len(data["result"].([]interface{})) != 0 {
		return data["result"].([]interface{})[0].(map[string]interface{}), nil
	}
	return nil, nil
}
