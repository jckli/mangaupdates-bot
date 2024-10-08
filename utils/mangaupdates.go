package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/jckli/mangaupdates-bot/mubot"
	"github.com/valyala/fasthttp"
)

var (
	username = os.Getenv("MU_USERNAME")
	password = os.Getenv("MU_PASSWORD")
)

func MuConvertOldId(oldID int64) string {
	const base = 36
	const digits = "0123456789abcdefghijklmnopqrstuvwxyz"

	if oldID == 0 {
		return "0"
	}

	var result strings.Builder
	for oldID > 0 {
		remainder := oldID % base
		result.WriteString(string(digits[remainder]))
		oldID /= base
	}

	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func MuConvertNewId(newID string) (int64, error) {
	const base = 36
	return strconv.ParseInt(newID, base, 64)
}

func MuLogin() (*MuLoginResponse, error) {
	login := MuLoginRequest{
		Username: username,
		Password: password,
	}
	loginJSON, err := json.Marshal(login)
	if err != nil {
		return nil, err
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("PUT")
	req.SetRequestURI("https://api.mangaupdates.com/v1/account/login")
	req.Header.Set("Content-Type", "application/json")
	req.SetBody(loginJSON)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = client.Do(req, resp)
	if err != nil {
		return nil, err
	}

	respBody := &MuLoginResponse{}
	if err = json.Unmarshal(resp.Body(), respBody); err != nil {
		return nil, err
	}

	return respBody, nil
}

func MuLogout(b *mubot.Bot) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("POST")
	req.SetRequestURI("https://api.mangaupdates.com/v1/account/logout")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+b.MuToken)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		return err
	}

	respBody := &MuLogoutResponse{}
	if err = json.Unmarshal(resp.Body(), respBody); err != nil {
		return err
	}

	return nil
}

func MuGetSeriesInfo(b *mubot.Bot, seriesId int64) (*MuSeriesInfoResponse, error) {
	var respBody *MuSeriesInfoResponse

	operation := func() error {
		resp, statusCode, err := muGetRequest(
			"https://api.mangaupdates.com/v1/series/"+strconv.FormatInt(seriesId, 10),
			b.MuToken,
		)

		if err != nil {
			return fmt.Errorf(
				"Failed to fetch series info: %s, %s, %d",
				err.Error(),
				string(resp),
				seriesId,
			)
		}

		if statusCode == 200 {
			respBody = &MuSeriesInfoResponse{}
			if err := json.Unmarshal(resp, respBody); err != nil {
				return fmt.Errorf("Failed to unmarshal series info: %s", err.Error())
			}
			return nil
		} else if statusCode >= 500 && statusCode < 600 {
			return fmt.Errorf("Series info server error: %d", statusCode)
		} else {
			return backoff.Permanent(fmt.Errorf("Series info client error: %d", statusCode))
		}
	}

	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.InitialInterval = 5 * time.Second
	backoffConfig.MaxInterval = 1 * time.Minute
	backoffConfig.MaxElapsedTime = 4 * time.Minute
	backoffConfig.Multiplier = 2
	backoffConfig.RandomizationFactor = 0.5

	retryPolicy := backoff.WithMaxRetries(backoffConfig, 10)

	err := backoff.Retry(operation, retryPolicy)
	if err != nil {
		return nil, fmt.Errorf(
			"Failed to get series info after retries: %v, Series ID: %d",
			err,
			seriesId,
		)
	}
	if respBody == nil {
		return nil, fmt.Errorf("MuGetSeriesInfo: respBody is nil after successful unmarshal")
	}

	return respBody, nil
}

func MuPostSearchGroups(b *mubot.Bot, groupName string) (*MuSearchGroupsResponse, error) {
	resp, err := muPostRequest(
		"https://api.mangaupdates.com/v1/groups/search",
		b.MuToken,
		MuSearchGroupsRequest{
			Search:  groupName,
			PerPage: 10,
		},
	)

	respBody := &MuSearchGroupsResponse{}
	if err = json.Unmarshal(resp, respBody); err != nil {
		return nil, err
	}

	return respBody, nil
}

func MuPostSearchSeries(b *mubot.Bot, seriesName string) (*MuSearchSeriesResponse, error) {
	resp, err := muPostRequest(
		"https://api.mangaupdates.com/v1/series/search",
		b.MuToken,
		MuSearchSeriesRequest{
			Search:  seriesName,
			PerPage: 10,
		},
	)

	respBody := &MuSearchSeriesResponse{}
	if err = json.Unmarshal(resp, respBody); err != nil {
		return nil, err
	}

	return respBody, nil
}
