package utils

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"strconv"
)

func parseAPIError(status int, body []byte) error {
	var errResp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != "" {
		return fmt.Errorf(fmt.Sprintf("%s. Join the support server for extra assistance.", errResp.Error))
	}

	switch status {
	case fasthttp.StatusNotFound:
		return fmt.Errorf("Resource not found. Please report this in the support server.")
	case fasthttp.StatusInternalServerError:
		return fmt.Errorf("Internal server error. Please report this in the support server.")
	default:
		return fmt.Errorf("API request failed (Status: %d). Please report this in the support server.", status)
	}
}

func (c *Client) GetWatchlist(endpoint string, id string) (*[]TrackedManga, error) {
	path := fmt.Sprintf("/tsuuchi/%s/%s", endpoint, id)

	body, status, err := c.Get(path)
	if err != nil {
		return nil, err
	}

	if status == fasthttp.StatusNotFound {
		return nil, nil
	}

	if status != fasthttp.StatusOK {
		return nil, parseAPIError(status, body)
	}

	var res []TrackedManga
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &res, nil
}

func (c *Client) SearchManga(query string) ([]MangaSearchResult, error) {
	u := fasthttp.AcquireURI()
	defer fasthttp.ReleaseURI(u)
	u.SetPath("/tsuuchi/manga/search")
	u.QueryArgs().Set("q", query)
	endpoint := string(u.RequestURI())

	body, status, err := c.Get(endpoint)
	if err != nil {
		return nil, err
	}
	if status != fasthttp.StatusOK {
		return nil, parseAPIError(status, body)
	}

	var res []MangaSearchResult
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) GetMangaDetails(mangaID int64) (*MangaDetails, error) {
	body, status, err := c.Get(fmt.Sprintf("/tsuuchi/manga/%d", mangaID))
	if err != nil {
		return nil, err
	}
	if status != fasthttp.StatusOK {
		return nil, parseAPIError(status, body)
	}

	var res MangaDetails
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) AddMangaToWatchlist(endpoint, id string, mangaID int64) error {
	path := fmt.Sprintf("/tsuuchi/%s/%s/manga", endpoint, id)
	payload := map[string]any{
		"id": mangaID,
	}

	body, status, err := c.Post(path, payload)
	if err != nil {
		return err
	}
	if status != fasthttp.StatusOK && status != fasthttp.StatusCreated {
		return parseAPIError(status, body)
	}
	return nil
}

func (c *Client) RemoveMangaFromWatchlist(endpoint, id string, mangaID int64) error {
	path := fmt.Sprintf("/tsuuchi/%s/%s/manga/%d", endpoint, id, mangaID)

	body, status, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	if status != fasthttp.StatusOK {
		return parseAPIError(status, body)
	}

	return nil
}

func (c *Client) SearchGroups(query string) ([]GroupSearchResult, error) {
	u := fasthttp.AcquireURI()
	defer fasthttp.ReleaseURI(u)

	u.SetPath("/tsuuchi/group/search")
	u.QueryArgs().Set("q", query)

	body, status, err := c.Get(string(u.RequestURI()))
	if err != nil {
		return nil, err
	}
	if status != fasthttp.StatusOK {
		return nil, parseAPIError(status, body)
	}

	var res []GroupSearchResult
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) UpdateMangaGroup(endpoint, targetID string, mangaID int64, groupName string, groupID int64) error {
	path := fmt.Sprintf("/tsuuchi/%s/%s/manga/%d/group", endpoint, targetID, mangaID)
	payload := map[string]any{
		"group_name": groupName,
		"group_id":   groupID,
	}

	body, status, err := c.Patch(path, payload)
	if err != nil {
		return err
	}
	if status != fasthttp.StatusOK {
		return parseAPIError(status, body)
	}

	return nil
}

func (c *Client) GetGroupDetails(groupID int64) (*GroupDetails, error) {
	body, status, err := c.Get(fmt.Sprintf("/tsuuchi/group/%d", groupID))
	if err != nil {
		return nil, err
	}
	if status != fasthttp.StatusOK {
		return nil, parseAPIError(status, body)
	}

	var res GroupDetails
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetMangaGroups(mangaID int64) ([]GroupSearchResult, error) {
	path := fmt.Sprintf("/tsuuchi/manga/%d/groups", mangaID)
	body, status, err := c.Get(path)
	if err != nil {
		return nil, err
	}
	if status != fasthttp.StatusOK {
		return nil, parseAPIError(status, body)
	}

	var res []GroupSearchResult
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) SetupServer(serverID, serverName, channelID string) error {
	path := fmt.Sprintf("/tsuuchi/server/%s/setup", serverID)

	payload := map[string]any{
		"server_name": serverName,
		"channel_id":  channelID,
	}

	body, status, err := c.Post(path, payload)
	if err != nil {
		return err
	}
	if status != fasthttp.StatusOK {
		return parseAPIError(status, body)
	}
	return nil
}

func (c *Client) SetupUser(userID, username string) error {
	path := fmt.Sprintf("/tsuuchi/user/%s/setup", userID)

	payload := map[string]any{
		"username": username,
	}

	body, status, err := c.Post(path, payload)
	if err != nil {
		return err
	}
	if status != fasthttp.StatusOK {
		return parseAPIError(status, body)
	}
	return nil
}

func (c *Client) DeleteServer(serverID string) error {
	path := fmt.Sprintf("/tsuuchi/server/%s", serverID)
	body, status, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	if status != fasthttp.StatusOK {
		return parseAPIError(status, body)
	}
	return nil
}

func (c *Client) DeleteUser(userID string) error {
	path := fmt.Sprintf("/tsuuchi/user/%s", userID)
	body, status, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	if status != fasthttp.StatusOK {
		return parseAPIError(status, body)
	}
	return nil
}

func (c *Client) GetServerConfig(serverID string) (*ServerConfig, error) {
	path := fmt.Sprintf("/tsuuchi/server/%s/config", serverID)

	body, status, err := c.Get(path)
	if err != nil {
		return nil, err
	}

	if status == fasthttp.StatusNotFound {
		return nil, nil
	}

	if status != fasthttp.StatusOK {
		return nil, parseAPIError(status, body)
	}

	var config ServerConfig
	if err := json.Unmarshal(body, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Client) GetUserConfig(userID string) (*UserConfig, error) {
	path := fmt.Sprintf("/tsuuchi/user/%s/config", userID)

	body, status, err := c.Get(path)
	if err != nil {
		return nil, err
	}

	if status == fasthttp.StatusNotFound {
		return nil, nil
	}

	if status != fasthttp.StatusOK {
		return nil, parseAPIError(status, body)
	}

	var config UserConfig
	if err := json.Unmarshal(body, &config); err != nil {
		return nil, fmt.Errorf("failed to decode user config: %w", err)
	}

	return &config, nil
}

func (c *Client) SetServerRole(serverID string, roleID string, roleType string) error {
	path := fmt.Sprintf("/tsuuchi/server/%s/role", serverID)

	rID, _ := strconv.ParseInt(roleID, 10, 64)
	payload := SetRoleRequest{
		RoleID:   rID,
		RoleType: roleType,
	}

	body, status, err := c.Post(path, payload)
	if err != nil {
		return err
	}
	if status != fasthttp.StatusOK {
		return parseAPIError(status, body)
	}
	return nil
}

func (c *Client) RemoveServerRole(serverID string, roleType string) error {
	path := fmt.Sprintf("/tsuuchi/server/%s/role", serverID)

	payload := SetRoleRequest{
		RoleType: roleType,
	}

	body, status, err := c.Delete(path, payload)
	if err != nil {
		return err
	}
	if status != fasthttp.StatusOK {
		return parseAPIError(status, body)
	}
	return nil
}

func (c *Client) UpdateServerChannel(serverID string, channelID string) error {
	path := fmt.Sprintf("/tsuuchi/server/%s/channel", serverID)

	cID, _ := strconv.ParseInt(channelID, 10, 64)
	payload := map[string]int64{
		"channel_id": cID,
	}

	body, status, err := c.Patch(path, payload)
	if err != nil {
		return err
	}
	if status != fasthttp.StatusOK {
		return parseAPIError(status, body)
	}
	return nil
}
