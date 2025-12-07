package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"quoter_back/schemas"

	"gorm.io/datatypes"
)

func (h *Handler) PublicGetList(c echo.Context) error {
	var quotes []struct {
		ID     string         `json:"id"`
		Text   string         `json:"text"`
		Author string         `json:"author"`
		Tags   datatypes.JSON `json:"tags" gorm:"type:json"`
	}

	if err := h.DB.Table("quotes").
		Select("quotes.id, quotes.quote_text AS text, quotes.author_name AS author, COALESCE(quotes.tags, '[]') AS tags").
		Joins("JOIN moderations ON quotes.id = moderations.quote_id").
		Where("moderations.status = ?", "approved").
		Scan(&quotes).Error; err != nil {
		log.Printf("Error fetching approved quotes: %v", err)
		return c.JSON(http.StatusInternalServerError, schemas.ErrorMessage{Error: "An error occurred while fetching quotes"})
	}

	// Convert tags from JSON to []string
	var response []struct {
		ID     string   `json:"id"`
		Text   string   `json:"text"`
		Author string   `json:"author"`
		Tags   []string `json:"tags"`
	}
	for _, quote := range quotes {
		var tags []string
		if err := json.Unmarshal(quote.Tags, &tags); err != nil {
			log.Printf("Error unmarshaling tags: %v", err)
			return c.JSON(http.StatusInternalServerError, schemas.ErrorMessage{Error: "An error occurred while processing tags"})
		}
		response = append(response, struct {
			ID     string   `json:"id"`
			Text   string   `json:"text"`
			Author string   `json:"author"`
			Tags   []string `json:"tags"`
		}{
			ID:     quote.ID,
			Text:   quote.Text,
			Author: quote.Author,
			Tags:   tags,
		})
	}

	if response == nil {
		response = make([]struct {
			ID     string   `json:"id"`
			Text   string   `json:"text"`
			Author string   `json:"author"`
			Tags   []string `json:"tags"`
		}, 0)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) PublicGetRandom(c echo.Context) error {
	var quote struct {
		ID     string `json:"id"`
		Text   string `json:"text"`
		Author string `json:"author"`
	}

	if err := h.DB.Table("quotes").
		Select("quotes.id, quotes.quote_text AS text, quotes.author_name AS author").
		Joins("JOIN moderations ON quotes.id = moderations.quote_id").
		Where("moderations.status = ?", "approved").
		Order("RANDOM()").
		Limit(1).
		Scan(&quote).Error; err != nil {
		log.Printf("Error fetching random approved quote: %v", err)
		return c.JSON(http.StatusInternalServerError, schemas.ErrorMessage{Error: "An error occurred while fetching a random quote"})
	}

	return c.JSON(http.StatusOK, quote)
}