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

	actionOperationTemplate = `{
		"operationtype": "0",
		"esc_period": "%v",
		"esc_step_from": "1",
		"esc_step_to": "%v",
		"evaltype": "0",
		"opconditions": [
			{
				"conditiontype":14,
				"value":0
			}
		],
		"opmessage": {
			"default_msg": "1",
			"subject": "Problem: {EVENT.NAME}",
			"message": "{EVENT.DATE} {EVENT.TIME}#{EVENT.NAME}#{EVENT.ID}#{ITEM.VALUE}#{ITEM.NAME}",
			"mediatypeid": "%v"
		},
		"opmessage_grp": [],
		"opmessage_usr": [
			%v
		]
	}`

	actionRecoveryOperationsTemplate = `{
		"operationtype": "0",
		"evaltype": "0",
		"opconditions": [],
		"opmessage": {
			"operationid": "14",
			"default_msg": "1",
			"subject": "Resolved: {EVENT.NAME}",
			"message": "{EVENT.RECOVERY.DATE} {EVENT.RECOVERY.TIME}#{EVENT.NAME}#{EVENT.ID}#{ITEM.VALUE}#{ITEM.NAME}",
			"mediatypeid": "%v"
		},
		"opmessage_grp": [],
		"opmessage_usr": [
			%v
		]
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
			"def_longdata": "告警触发: \r\nEventID: {EVENT.ID}\r\nStartAt: {EVENT.TIME} {EVENT.DATE}\r\n{TRIGGER.DESCRIPTION}\r\n{TRIGGER.NAME}: {ITEM.VALUE}",
			"r_shortdata": "Resolved: {EVENT.NAME}",
			"r_longdata": "告警恢复: \r\nEventID: {EVENT.ID}\r\nStartAt: {EVENT.TIME} {EVENT.DATE}\r\n{TRIGGER.DESCRIPTION}\r\n{TRIGGER.NAME}: {ITEM.VALUE}",
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
			"operations": [%v],
			"recovery_operations": [%v]
		},
		"auth": "%v",
		"id": %v
	}`
)

// ActionCreate create zabbix action
// toWX key mediatypeid of weixin  value is admin weixin user, toEmail mediatypeid of emal value of the user email
func (api *API) ActionCreate(name, interval, StepDuration string, alertNum int32, mediatypeidUserIDS map[string]string) error {
	operations := []string{}
	recoveryPperations := []string{}
	for mediatypeid, userIDS := range mediatypeidUserIDS {
		users := ""
		ids := strings.Split(userIDS, ",")
		for index, id := range ids {
			if index == 0 {
				users += fmt.Sprintf(`{"operationid":"14","userid":"%v"}`, id)
			} else {
				users += "," + fmt.Sprintf(`{"operationid":"14","userid":"%v"}`, id)
			}
		}
		operations = append(operations, fmt.Sprintf(actionOperationTemplate, StepDuration, alertNum, mediatypeid, users))
		recoveryPperations = append(recoveryPperations, fmt.Sprintf(actionRecoveryOperationsTemplate, mediatypeid, users))
	}
	payload := strings.NewReader(fmt.Sprintf(actionPostTemplate, name, interval, name, strings.Join(operations, ","), strings.Join(recoveryPperations, ","), api.Session, api.ID))
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

// ActionDelete delete zabbix action
func (api *API) ActionDelete(ids []string) error {
	names := ""
	for _, id := range ids {
		names += fmt.Sprintf(`"%v"`, id)
	}
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

// ActionGet get zabbix action
func (api *API) ActionGet(name string) (map[string]interface{}, error) {
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
