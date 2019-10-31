package zabbix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	UserPostTemplate = `{
		"jsonrpc": "2.0",
		"method": "user.create",
		"params": {
			"alias": "%v",
			"passwd": "%v",
			"usrgrps": [
				{
					"usrgrpid": "%v"
				}
			],
			"user_medias": [
				{
					"mediatypeid": "1",
					"sendto": [
						"%v"
					],
					"active": 0,
					"severity": 63,
					"period": "1-7,00:00-24:00"
				}
			]
		},
		"auth": "%v",
		"id": %v
	}`

	UserGetTemplate = `{
		"jsonrpc": "2.0",
		"method": "user.get",
		"params": {
			"output": "extend",
			"filter":{
				"alias":"%v"
			}
		},
		"auth": "%v",
		"id": %v
	}`
)

func (api *API) UserCreate(name, password, mail, groupID string) error {
	payload := strings.NewReader(fmt.Sprintf(UserPostTemplate, name, password, groupID, mail, api.Session, api.ID))
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

func (api *API) UserGet(name, password, mail string) (map[string]interface{}, error) {
	payload := strings.NewReader(fmt.Sprintf(UserGetTemplate, name, api.Session, api.ID))
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
		result, err := api.UserGroupGet("Guests")
		if err != nil {
			return nil, err
		}
		if result != nil && result["usrgrpid"] != nil {
			err = api.UserCreate(name, password, mail, result["usrgrpid"].(string))
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
