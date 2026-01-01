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
		return nil, fmt.Errorf("api returned status: %d", status)
	}

	var res []TrackedManga
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &res, nil
}

func (c *Client) SearchManga(query string) ([]MangaSearchResult, error) {
	body, status, err := c.Get("/tsuuchi/manga/search?q=" + query)
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
		return nil, fmt.Errorf("Details fetch failed with status: %d", status)
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
