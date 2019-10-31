package zabbix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (api *API) HostGroupGet(name string) (map[string]interface{}, error) {
	payload := strings.NewReader(fmt.Sprintf(HostGroupGetTemplate, name, api.Session, api.ID))
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

	result := map[string]interface{}{}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	if len(result["result"].([]interface{})) != 0 {
		return result["result"].([]interface{})[0].(map[string]interface{}), nil
	}
	return nil, nil
}

func (api *API) HostGroupCreate(name string) error {
	payload := strings.NewReader(fmt.Sprintf(HostGroupPostTemplate, name, api.Session, api.ID))
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
