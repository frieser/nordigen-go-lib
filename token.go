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
)

type Token struct {
	Access         string `json:"access"`
	AccessExpires  int    `json:"access_expires"`
	Refresh        string `json:"refresh"`
	RefreshExpires int    `json:"refresh_expires"`
}

type Secret struct {
	SecretId string `json:"secret_id"`
	AccessId string `json:"secret_key"`
}

const tokenPath = "token"
const tokenNewPath = "new/"
const tokenRefreshPath = "refresh"

func (c Client) newToken(secretId, secretKey string) (*Token, error) {
	req := http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme: "https",
			Host: baseUrl,
			Path: strings.Join([]string{apiPath, tokenPath, tokenNewPath}, "/"),
		},
	}
	req.Header = http.Header{}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	data, err := json.Marshal(Secret{
		SecretId: secretId,
		AccessId: secretKey,
	})
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(data))
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
	t := &Token{}
	err = json.Unmarshal(body, &t)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (c Client) refreshToken(refresh string) (*Token, error) {
	req := http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Path: strings.Join([]string{tokenPath, tokenRefreshPath}, "/"),
		},
	}
	data, err := json.Marshal(refresh)

	if err != nil {
		return &Token{}, err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(data))

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
	t := &Token{}
	err = json.Unmarshal(body, &t)

	if err != nil {
		return nil, err
	}

	return t, nil
}
