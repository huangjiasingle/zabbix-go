package zabbix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	hostGetTemplate = `{
		"jsonrpc": "2.0",
		"method": "host.get",
		"params": {
			"output": "extend",
			"filter": {
				"name": "%v"
			}
		},
		"auth": "%v",
		"id": %v
	}`

	hostAddTemplateTemplate = `{
		"jsonrpc": "2.0",
		"method": "host.massadd",
		"params": {
			"hosts":
				{
					"hostid": "%v"
				},
		   
			"templates":{
					"templateid": "%v"
				}
			},
		"auth": "%v",
		"id": %v
	}`

	hostDeleteTemplateTemplate = `{
		"jsonrpc": "2.0",
		"method": "host.massremove",
		"params": {
			"hostids": "%v",
		    "templateids_clear": "%v"
		}
		"auth": "%v",
		"id": %v
	}`
)

func (api *API) HostGet(name string) (map[string]interface{}, error) {
	payload := strings.NewReader(fmt.Sprintf(hostGetTemplate, name, api.Session, api.ID))
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

func (api *API) HostAddTemplate(hostid, templateid string) error {
	// fmt.Println(fmt.Sprintf(hostAddTemplateTemplate, hostid, templateid, api.Session, api.ID))
	payload := strings.NewReader(fmt.Sprintf(hostAddTemplateTemplate, hostid, templateid, api.Session, api.ID))
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

func (api *API) HostDeleteTemplate(hostid, templateid string) error {
	// fmt.Println(fmt.Sprintf(hostDeleteTemplateTemplate, hostid, templateid, api.Session, api.ID))
	payload := strings.NewReader(fmt.Sprintf(hostDeleteTemplateTemplate, hostid, templateid, api.Session, api.ID))
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
