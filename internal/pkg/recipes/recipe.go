package recipes

import (
	"encoding/json"
	"log"
)

// Recipe wraps parsed recipe data from a mela archive
type Recipe struct {
	ID string `json:"recipe_id"`

	SourceID string  `json:"id"`
	Date     float64 `json:"date"`

	Title        string `json:"title"`
	Ingredients  string `json:"ingredients"`
	Instructions string `json:"instructions"`
	Description  string `json:"text"`
	Notes        string `json:"notes"`

	Images []string `json:"images"`

	Link string `json:"link"`

	Categories []string `json:"categories"`

	Yield string `json:"yield"`

	Favorite   bool `json:"favorite"`
	WantToCook bool `json:"wantToCook"`
}

func (r *Recipe) JSON() string {
	cp := *r
	cp.Images = []string{}
	json, err := json.MarshalIndent(cp, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return string(json)
}
