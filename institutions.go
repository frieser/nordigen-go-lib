package nordigen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const institutionsPath = "institutions"
const countryParam = "country"

type Institution struct {
	Id                   string   `json:"id"`
	Name                 string   `json:"name"`
	Bic                  string   `json:"bic"`
	TransactionTotalDays string   `json:"transaction_total_days"`
	Countries            []string `json:"countries"`
	Logo                 string   `json:"logo"`
}

func (c Client) ListInstitutions(country string) ([]Institution, error) {
	req := http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: strings.Join([]string{institutionsPath, ""}, "/"),
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
	list := make([]Institution, 0)
	err = json.Unmarshal(body, &list)

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (c Client) GetInstitution(institutionID string) (Institution, error) {
	req := http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: strings.Join([]string{institutionsPath, institutionID, ""}, "/"),
		},
	}
	resp, err := c.c.Do(&req)

	if err != nil {
		return Institution{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return Institution{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return Institution{}, fmt.Errorf("expected %d status code: got %d", http.StatusOK, resp.StatusCode)
	}
	insttn := Institution{}
	err = json.Unmarshal(body, &insttn)

	if err != nil {
		return Institution{}, err
	}

	return insttn, nil
}
