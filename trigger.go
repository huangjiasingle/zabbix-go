package zabbix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	triggerPostTemplate = `{
		"jsonrpc": "2.0",
		"method": "trigger.create",
		"params": [{
			"description": "%v",
			"priority": 4,
			"expression": "%v",
			"comments": "%v"
		}],
		"auth": "%v",
		"id": %v
	}`

	triggerGetTemplate = `{
		"jsonrpc": "2.0",
		"method": "trigger.get",
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

func (api *API) TriggerCreate(name, expression, comments string) error {
	payload := strings.NewReader(fmt.Sprintf(triggerPostTemplate, name, expression, comments, api.Session, api.ID))
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

func (api *API) TriggerGet(name string) (map[string]interface{}, error) {
	payload := strings.NewReader(fmt.Sprintf(triggerGetTemplate, name, api.Session, api.ID))
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
	if len(data["result"].([]interface{})) != 0 {
		return data["result"].([]interface{})[0].(map[string]interface{}), nil
	}
	return nil, nil
}
