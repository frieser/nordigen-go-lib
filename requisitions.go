package nordigen

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const requisitionsPath = "requisitions"

type Requisition struct {
	Id       string    `json:"id,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Redirect string    `json:"redirect,omitempty"`
	Status   string    `json:"status,omitempty"`
	// There is an issue in the api, the status is still a string
	// like in v1
	//Status        Status    `json:"status,omitempty"`
	InstitutionId string   `json:"institution_id,omitempty"`
	Agreement     string   `json:"agreement,omitempty"`
	Reference     string   `json:"reference,omitempty"`
	Accounts      []string `json:"accounts,omitempty"`
	UserLanguage  string   `json:"user_language,omitempty"`
	Link          string   `json:"link,omitempty"`
}

type Status struct {
	Short       string `json:"short,omitempty"`
	Long        string `json:"long,omitempty"`
	Description string `json:"description,omitempty"`
}

func (c Client) CreateRequisition(r Requisition) (Requisition, error) {
	req := http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Path: strings.Join([]string{requisitionsPath, ""}, "/"),
		},
	}
	data, err := json.Marshal(r)

	if err != nil {
		return Requisition{}, err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(data))

	resp, err := c.c.Do(&req)

	if err != nil {
		return Requisition{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return Requisition{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return Requisition{}, &APIError{resp.StatusCode, string(body), err}
	}
	err = json.Unmarshal(body, &r)

	if err != nil {
		return Requisition{}, err
	}

	return r, nil
}

func (c Client) GetRequisition(id string) (r Requisition, err error) {
	req := http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: strings.Join([]string{requisitionsPath, id, ""}, "/"),
		},
	}
	resp, err := c.c.Do(&req)

	if err != nil {
		return Requisition{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return Requisition{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Requisition{}, &APIError{resp.StatusCode, string(body), err}
	}
	err = json.Unmarshal(body, &r)

	if err != nil {
		return Requisition{}, err
	}

	return r, nil
}

func (c Client) ListRequisitions() (r []Requisition, err error) {
	req := http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: strings.Join([]string{requisitionsPath, ""}, "/"),
		},
	}
	resp, err := c.c.Do(&req)

	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{resp.StatusCode, string(body), err}
	}
	
	list := make([]Requisition, 0)
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
