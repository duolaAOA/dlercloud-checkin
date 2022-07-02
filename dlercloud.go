package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type DlerCloud struct {
	email    string
	password string

	token string
}

type CheckinInfo struct {
	Checkin   string `json:"checkin"`
	TodayUsed string `json:"today_used"`
	Used      string `json:"used"`
	Unused    string `json:"unused"`
	Traffic   string `json:"traffic"`
}

type UserInfo struct {
	Plan      string `json:"plan"`
	PlanTime  string `json:"plan_time"`
	Money     string `json:"money"`
	AffMoney  string `json:"aff_money"`
	TodayUsed string `json:"today_used"`
	Used      string `json:"used"`
	Unused    string `json:"unused"`
	Traffic   string `json:"traffic"`
	Integral  string `json:"integral"`
}

func NewClient(email string, password string) *DlerCloud {
	return &DlerCloud{email: email, password: password}
}

func (d *DlerCloud) login(ctx context.Context) error {
	var response = new(struct {
		Token string `json:"token"`
	})
	err := d.post(ctx, "login", map[string]interface{}{
		"email":  d.email,
		"passwd": d.password,
	}, response)
	if err != nil {
		return err
	}

	d.token = response.Token
	return nil
}

func (d *DlerCloud) TryToCheckin(ctx context.Context) (*CheckinInfo, error) {
	if err := d.login(ctx); err != nil {
		return nil, fmt.Errorf("not logged in")
	}

	var response = new(CheckinInfo)
	err := d.post(ctx, "checkin", map[string]interface{}{
		"access_token": d.token,
		"multiple":     1,
	}, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *DlerCloud) GetUserInfo(ctx context.Context) (*UserInfo, error) {
	if err := d.login(ctx); err != nil {
		return nil, fmt.Errorf("not logged in")
	}

	var response = new(UserInfo)
	err := d.post(ctx, "information", map[string]interface{}{
		"access_token": d.token,
	}, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *DlerCloud) getURL(path string) string {
	const urlFmt = `https://dler.cloud/api/v1/%s`
	return fmt.Sprintf(urlFmt, path)
}

func (d *DlerCloud) post(ctx context.Context, path string, body map[string]interface{}, dest interface{}) error {
	var (
		httpReq *http.Request
		err     error
	)

	if body != nil {
		form := make(url.Values, len(body))
		for k, v := range body {
			form.Set(k, fmt.Sprint(v))
		}
		httpReq, err = http.NewRequest(http.MethodPost, d.getURL(path), strings.NewReader(form.Encode()))
		if err != nil {
			return fmt.Errorf("failed to create HTTP request: %+v", err)
		}
		httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		httpReq, err = http.NewRequest(http.MethodPost, d.getURL(path), nil)
		if err != nil {
			return fmt.Errorf("failed to create HTTP request: %+v", err)
		}
	}

	httpReq = httpReq.WithContext(ctx)

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to do request: %+v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid response status code: %d", httpResp.StatusCode)
	}
	if httpResp.Body == nil {
		return fmt.Errorf("response body is nil")
	}

	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %+v", err)
	}
	resp := new(response)
	if err := json.Unmarshal(respBody, resp); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %+v", err)
	}
	if resp.Code != http.StatusOK {
		return fmt.Errorf("invalid result code %d with message: %s", resp.Code, resp.Message)
	}

	if dest == nil {
		return nil
	}
	if err := json.Unmarshal(resp.Data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal response data: %+v", err)
	}
	return nil
}

type response struct {
	Code    int             `json:"ret"`
	Message string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
}
