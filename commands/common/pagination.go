package common

import (
	"fmt"
	"math"

	"github.com/disgoorg/disgo/discord"
)

const ItemsPerPage = 15

func GetPageSlice[T any](items []T, page int) ([]T, int) {
	totalItems := len(items)
	const itemsPerPage = 25

	totalPages := int(math.Ceil(float64(totalItems) / float64(itemsPerPage)))

	if page < 1 {
		page = 1
	}
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	if totalPages == 0 {
		return []T{}, 0
	}

	start := (page - 1) * itemsPerPage
	end := start + itemsPerPage
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
