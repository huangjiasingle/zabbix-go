package zabbix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	actionGetTemplate = `{
		"jsonrpc":"2.0",
		"method":"action.get",
		"params":{
			"output":"extend",
			"selectOperations":"extend",
			"selectRecoveryOperations":"extend",
			"selectFilter":"extend",
			"filter":{
				"eventsource":0,
				"name":"%v"
			}
		},
		"auth":"%v",
		"id":%v
	}`

	actionDeleteTemplate = `{
		"jsonrpc": "2.0",
		"method": "action.delete",
		"params": [
			%v
		],
		"auth": "%v",
		"id": %v
	}`

	actionPostTemplate = `{
		"jsonrpc": "2.0",
		"method": "action.create",
		"params": {
			"name": "%v",
			"eventsource": "0",
			"status": "0",
			"esc_period": "%v",  
			"def_shortdata": "Problem: {EVENT.NAME}",
			"def_longdata": "Problem started at {EVENT.TIME} on {EVENT.DATE}\r\nProblem name: {EVENT.NAME}\r\nHost: {HOST.NAME}\r\nSeverity: {EVENT.SEVERITY}\r\n\r\nOriginal problem ID: {EVENT.ID}\r\n{TRIGGER.URL}",
			"r_shortdata": "Resolved: {EVENT.NAME}",
			"r_longdata": "Problem has been resolved at {EVENT.RECOVERY.TIME} on {EVENT.RECOVERY.DATE}\r\nProblem name: {EVENT.NAME}\r\nHost: {HOST.NAME}\r\nSeverity: {EVENT.SEVERITY}\r\n\r\nOriginal problem ID: {EVENT.ID}\r\n{TRIGGER.URL}",
			"pause_suppressed": "1",
			"ack_shortdata": "Updated problem: {EVENT.NAME}",
			"ack_longdata": "{USER.FULLNAME} {EVENT.UPDATE.ACTION} problem at {EVENT.UPDATE.DATE} {EVENT.UPDATE.TIME}.\r\n{EVENT.UPDATE.MESSAGE}\r\n\r\nCurrent problem status is {EVENT.STATUS}, acknowledged: {EVENT.ACK.STATUS}.",
			"filter": {
				"evaltype": "0",
				"conditions": [{
					"conditiontype": "3",
					"operator": "2",
					"value": "%v"
				}]
			},
			"operations": [{
				"operationtype": "0",
				"esc_period": "0",
				"esc_step_from": "1",
				"esc_step_to": "1",
				"evaltype": "0",
				"opconditions": [],
				"opmessage": {
					"default_msg": "1",
					"subject": "Problem: {EVENT.NAME}",
					"message": "Problem started at {EVENT.TIME} on {EVENT.DATE}\r\nProblem name: {EVENT.NAME}\r\nHost: {HOST.NAME}\r\nSeverity: {EVENT.SEVERITY}\r\n\r\nOriginal problem ID: {EVENT.ID}\r\n{TRIGGER.URL}",
					"mediatypeid": "0"
				},
				"opmessage_grp": [],
				"opmessage_usr": [
					%v
				]
			}],
			"recovery_operations": [{
				"operationtype": "0",
				"evaltype": "0",
				"opconditions": [],
				"opmessage": {
					"operationid": "14",
					"default_msg": "1",
					"subject": "Resolved: {EVENT.NAME}",
					"message": "Problem has been resolved at {EVENT.RECOVERY.TIME} on {EVENT.RECOVERY.DATE}\r\nProblem name: {EVENT.NAME}\r\nHost: {HOST.NAME}\r\nSeverity: {EVENT.SEVERITY}\r\n\r\nOriginal problem ID: {EVENT.ID}\r\n{TRIGGER.URL}",
					"mediatypeid": "0"
				},
				"opmessage_grp": [],
				"opmessage_usr": [
					%v
				]
			}]
		},
		"auth": "%v",
		"id": %v
	}`
)

func (api *API) ActionCreate(name, interval, to string) error {
	users := ""
	ids := strings.Split(to, ",")
	fmt.Println(ids)
	for index, id := range ids {
		if index == 0 {
			users += fmt.Sprintf(`{"operationid":"14","userid":"%v"}`, id)
		} else {
			users += "," + fmt.Sprintf(`{"operationid":"14","userid":"%v"}`, id)
		}
	}

	fmt.Println(fmt.Sprintf(actionPostTemplate, name, interval, name, users, users, api.Session, api.ID))

	payload := strings.NewReader(fmt.Sprintf(actionPostTemplate, name, interval, name, users, users, api.Session, api.ID))
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

func (api *API) ActionDelete(ids []string) error {
	names := ""
	for _, id := range ids {
		names += fmt.Sprintf(`"%v"`, id)
	}
	fmt.Print(fmt.Sprintf(actionDeleteTemplate, names, api.Session, api.ID))
	payload := strings.NewReader(fmt.Sprintf(actionDeleteTemplate, names, api.Session, api.ID))
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

	result := map[string]interface{}{}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return err
	}
	if errmsg, ok := result["error"]; ok {
		return fmt.Errorf("%v", errmsg)
	}
	return nil
}

func (api *API) ActionGet(name string) (map[string]interface{}, error) {
	fmt.Println(fmt.Sprintf(actionGetTemplate, name, api.Session, api.ID))
	payload := strings.NewReader(fmt.Sprintf(actionGetTemplate, name, api.Session, api.ID))
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
