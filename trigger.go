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
			"expression": "%v"
		}],
		"auth": "%v",
		"id": %v
	}`
)

func (api *API) TriggerCreate(name, expression string) error {
	fmt.Println(fmt.Sprintf(triggerPostTemplate, name, expression, api.Session, api.ID))
	payload := strings.NewReader(fmt.Sprintf(triggerPostTemplate, name, expression, api.Session, api.ID))
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
