package utils

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

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
	resp, err := muGetRequest(
		"https://api.mangaupdates.com/v1/series/"+strconv.FormatInt(seriesId, 10),
		b.MuToken,
	)

	respBody := &MuSeriesInfoResponse{}
	if err = json.Unmarshal(resp, respBody); err != nil {
		return nil, err
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
