package nordigen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const aspspsPath = "aspsps"
const countryParam = "country"

type Aspsps struct {
	Id                   string   `json:"id"`
	Name                 string   `json:"name"`
	Bic                  string   `json:"bic"`
	TransactionTotalDays string   `json:"transaction_total_days"`
	Countries            []string `json:"countries"`
	Logo                 string   `json:"logo"`
}

func (c Client) ListAspsps(country string) ([]Aspsps, error) {
	req := http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: strings.Join([]string{aspspsPath, ""}, "/"),
		},
	}
	q := req.URL.Query()
	q.Add(countryParam, country)
	req.URL.RawQuery = q.Encode()

	resp, err := c.c.Do(&req)

	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %d status code: got %d", http.StatusOK, resp.StatusCode)
	}
	list := make([]Aspsps, 0)
	err = json.Unmarshal(body, &list)

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (c Client) GetAspsps(aspspsID string) (Aspsps, error) {
	req := http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: strings.Join([]string{aspspsPath, aspspsID, ""}, "/"),
		},
	}
	resp, err := c.c.Do(&req)

	if err != nil {
		return Aspsps{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return Aspsps{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return Aspsps{}, fmt.Errorf("expected %d status code: got %d", http.StatusOK, resp.StatusCode)
	}
	aspsps := Aspsps{}
	err = json.Unmarshal(body, &aspsps)

	if err != nil {
		return Aspsps{}, err
	}

	return aspsps, nil
}
