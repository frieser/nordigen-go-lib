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

func (c Client) ListRequisitions() (rs []Requisition, err error) {
	url := &url.URL{
		Path: strings.Join([]string{requisitionsPath, ""}, "/"),
	}

	err = c.fetchRequisitions(url, &rs)
	if err != nil {
		return []Requisition{}, err
	}

	return rs, nil
}

type requisitionsResponse struct {
	Count    int           `json:"count"`
	Next     string        `json:"next"`
	Previous string        `json:"previous"`
	Results  []Requisition `json:"results"`
}

// fetchRequisitions recursively
func (c Client) fetchRequisitions(u *url.URL, allRequisitions *[]Requisition) error {
	req := http.Request{
		Method: http.MethodGet,
		URL:    u,
	}

	resp, err := c.c.Do(&req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &APIError{resp.StatusCode, string(body), err}
	}

	if resp.StatusCode != http.StatusOK {
		return &APIError{resp.StatusCode, string(body), err}
	}

	var requisitions requisitionsResponse
	err = json.Unmarshal(body, &requisitions)
	if err != nil {
		return &APIError{resp.StatusCode, string(body), err}
	}

	*allRequisitions = append(*allRequisitions, requisitions.Results...)

	// If there is a next URL, make a recursive call
	if requisitions.Next != "" {
		next, err := url.Parse(requisitions.Next)
		if err != nil {
			panic("requisitions pagination url is invalid")
		}
		// The client expects the URL to be just the path, mangle it here and
		// append the query we got from the pagination URL.
		return c.fetchRequisitions(&url.URL{
			Path:     strings.Join([]string{requisitionsPath, ""}, "/"),
			RawQuery: next.RawQuery,
		}, allRequisitions)
	}

	return nil
}
