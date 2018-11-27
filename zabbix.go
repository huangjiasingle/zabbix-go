package zabbix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type API struct {
	JsonRPC string
	URL     string
	Session string
	ID      int
}

func NewAPI(url, user, password string, id int) (*API, error) {
	payload := strings.NewReader(fmt.Sprintf(UserLoginTemplate, user, password, false, id))
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if errmsg, ok := result["error"]; ok {
		return nil, fmt.Errorf("%v", errmsg)
	}
	return &API{JsonRPC: result["jsonrpc"].(string), URL: url, Session: result["result"].(string), ID: id}, nil
}
