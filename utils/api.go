package utils

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)

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
		return nil, fmt.Errorf("API returned status: %d", status)
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
		return nil, fmt.Errorf("Search failed with status: %d", status)
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
		return nil, fmt.Errorf("Manga details fetch failed with status: %d", status)
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

	_, status, err := c.Post(path, payload)
	if err != nil {
		return err
	}
	if status == fasthttp.StatusConflict || status == fasthttp.StatusBadRequest {
		return fmt.Errorf("This manga is already on the list.")
	}
	if status != fasthttp.StatusOK && status != fasthttp.StatusCreated {
		return fmt.Errorf("API returned status: %d. Please report this to me in my support server.", status)
	}
	return nil
}

func (c *Client) RemoveMangaFromWatchlist(endpoint, id string, mangaID int64) error {
	path := fmt.Sprintf("/tsuuchi/%s/%s/manga/%d", endpoint, id, mangaID)

	_, status, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	if status == fasthttp.StatusNotFound {
		return fmt.Errorf("This manga is not found in the list.")
	}
	if status != fasthttp.StatusOK {
		return fmt.Errorf("API returned status: %d. Please report this to me in my support server.", status)
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
		return nil, fmt.Errorf("API returned status: %d", status)
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

	_, status, err := c.Patch(path, payload)
	if err != nil {
		return err
	}
	if status != fasthttp.StatusOK {
		return fmt.Errorf("API returned status: %d", status)
	}

	return nil
}

func (c *Client) GetGroupDetails(groupID int64) (*GroupDetails, error) {
	body, status, err := c.Get(fmt.Sprintf("/tsuuchi/group/%d", groupID))
	if err != nil {
		return nil, err
	}
	if status != fasthttp.StatusOK {
		return nil, fmt.Errorf("Group details fetch failed with status: %d", status)
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
		return nil, fmt.Errorf("API returned status: %d", status)
	}

	var res []GroupSearchResult
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res, nil
}
