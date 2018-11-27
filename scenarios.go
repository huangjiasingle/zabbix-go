package zabbix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	scenariosPostTemplate = `{
    "jsonrpc": "2.0",
    "method": "httptest.create",
    "params": {
        "name": "%v",
		"hostid": "%v",
		"delay": "%v",
        "steps": [
            {
                "name": "%v",
                "url": "%v",
                "status_codes": "200",
                "no": 1
            }
        ]
    },
    "auth": "%v",
    "id": %v
}`

	scenariosDeleteTemplate = `{
    "jsonrpc": "2.0",
    "method": "httptest.delete",
    "params": [
        "%v"
    ],
    "auth": "%v",
    "id": %v
}`

	scenariosGetTemplate = `{
    "jsonrpc": "2.0",
    "method": "httptest.get",
    "params": {
        "output": "extend",
        "selectSteps": "extend",
        "filter":{
        	"name":"%v",
        	"hostid": "%v"
        }
    },
    "auth": "%v",
    "id": %v
}`
)

func (api *API) ScenariosCreate(name, hostid, interval, url string) error {
	payload := strings.NewReader(fmt.Sprintf(scenariosPostTemplate, name, hostid, interval, name, url, api.Session, api.ID))
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

func (api *API) ScenariosDelete(id string) error {
	payload := strings.NewReader(fmt.Sprintf(scenariosDeleteTemplate, id, api.Session, api.ID))
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

func (api *API) ScenariosGet(name, hostid string) (map[string]interface{}, error) {
	payload := strings.NewReader(fmt.Sprintf(scenariosGetTemplate, name, hostid, api.Session, api.ID))
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
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("zabbix api return response code %v", res.StatusCode)
	}
	result := map[string]interface{}{}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	if errmsg, ok := result["error"]; ok {
		return nil, fmt.Errorf("%v", errmsg)
	}
	if len(result["result"].([]interface{})) != 0 {
		return result["result"].([]interface{})[0].(map[string]interface{}), nil
	}
	return nil, nil
}
