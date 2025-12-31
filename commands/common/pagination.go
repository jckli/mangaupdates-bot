package common

import (
	"fmt"
	"math"

	"github.com/disgoorg/disgo/discord"
)

const ItemsPerPage = 15

func GetPageSlice(items []string, page int) ([]string, int) {
	totalItems := len(items)
	totalPages := int(math.Ceil(float64(totalItems) / float64(ItemsPerPage)))

	if page < 1 {
		page = 1
	}
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}
	if totalPages == 0 {
		return []string{}, 0
	}

	start := (page - 1) * ItemsPerPage
	end := start + ItemsPerPage
	if end > totalItems {
		end = totalItems
	}

	return items[start:end], totalPages
}

func GeneratePaginationButtons(
	buttonPrefix string,
	page int,
	totalPages int,
) []discord.ContainerComponent {

	if totalPages <= 1 {
		return nil
	}

	return []discord.ContainerComponent{
		discord.ActionRowComponent{
			discord.NewSecondaryButton(
				"◀",
				fmt.Sprintf("%s/%d", buttonPrefix, page-1),
			).WithDisabled(page == 1),
			discord.NewSecondaryButton(fmt.Sprintf("%d/%d", page, totalPages), "page-counter").
				WithDisabled(true),
			discord.NewSecondaryButton(
				"▶",
				fmt.Sprintf("%s/%d", buttonPrefix, page+1),
			).WithDisabled(page == totalPages),
		},
	}
}
