package nordigen

import (
	"bytes"
	"context"
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

type TokenRefresh struct {
	Refresh string `json:"refresh"`
}

type Secret struct {
	SecretId string `json:"secret_id"`
	AccessId string `json:"secret_key"`
}

const tokenPath = "token"
const tokenNewPath = "new/"
const tokenRefreshPath = "refresh/"

func (c Client) newToken(ctx context.Context) (*Token, error) {
	data, err := json.Marshal(Secret{
		SecretId: c.secretId,
		AccessId: c.secretKey,
	})
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodPost,
		Body:   io.NopCloser(bytes.NewBuffer(data)),
		URL: &url.URL{
			Path: strings.Join([]string{tokenPath, tokenNewPath}, "/"),
		},
	}
	req = req.WithContext(ctx)

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	t := &Token{}
	if err := json.Unmarshal(body, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (c Client) refreshToken(ctx context.Context, refresh string) (*Token, error) {
	data, err := json.Marshal(TokenRefresh{Refresh: refresh})
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodPost,
		Body:   io.NopCloser(bytes.NewBuffer(data)),
		URL: &url.URL{
			Path: strings.Join([]string{tokenPath, tokenRefreshPath}, "/"),
		},
	}
	req = req.WithContext(ctx)

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	t := &Token{}
	if err := json.Unmarshal(body, t); err != nil {
		return nil, err
	}
	return t, nil
}
