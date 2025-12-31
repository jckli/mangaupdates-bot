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
