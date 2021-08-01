package nordigen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type EndUserAgreement struct {
	Id                 string      `json:"id,omitempty"`
	Created            time.Time   `json:"created,omitempty"`
	Accepted           interface{} `json:"accepted,omitempty"`
	MaxHistoricalDays  int         `json:"max_historical_days,omitempty"`
	AccessValidForDays int         `json:"access_valid_for_days,omitempty"`
	EnduserId          string      `json:"enduser_id,omitempty"`
	AspspId            string      `json:"aspsp_id,omitempty"`
	AgreementVersion   string      `json:"agreement_version,omitempty"`
}

const agreementsPath = "agreements"
const endUserPath = "enduser"

func (c Client) CreateEndUserAgreement(eua EndUserAgreement) (EndUserAgreement, error) {
	req := http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Path: strings.Join([]string{agreementsPath, endUserPath, ""}, "/"),
		},
	}
	data, err := json.Marshal(eua)

	if err != nil {
		return EndUserAgreement{}, nil
	}
	req.Body = io.NopCloser(bytes.NewBuffer(data))

	resp, err := c.c.Do(&req)

	if err != nil {
		return EndUserAgreement{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return EndUserAgreement{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return EndUserAgreement{}, fmt.Errorf("expected %d status code: got %d", http.StatusCreated, resp.StatusCode)
	}
	err = json.Unmarshal(body, &eua)

	if err != nil {
		return EndUserAgreement{}, err
	}

	return eua, nil
}
