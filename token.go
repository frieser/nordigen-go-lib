package nordigen

import (
	"bytes"
	"encoding/json"
	"io"
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

func (c Client) newToken() (*Token, error) {
	req := http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Path: strings.Join([]string{tokenPath, tokenNewPath}, "/"),
		},
	}

	data, err := json.Marshal(Secret{
		SecretId: c.secretId,
		AccessId: c.secretKey,
	})
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(data))
	resp, err := c.c.Do(&req)

	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{resp.StatusCode, string(body), err}
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
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{resp.StatusCode, string(body), err}
	}
	t := &Token{}
	err = json.Unmarshal(body, &t)

	if err != nil {
		return nil, err
	}

	return t, nil
}
