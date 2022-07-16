package recipes

import (
	"archive/zip"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func Parse(data []byte) ([]Recipe, error) {
	archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		log.Fatal(err)
	}

	var recipes []Recipe

	for _, f := range archive.File {
		if !strings.HasSuffix(f.Name, ".melarecipe") {
			continue
		}

		recipeFile, err := f.Open()
		if err != nil {
			return []Recipe{}, fmt.Errorf("failed to open the recipe file: %w", err)
		}

		data, err := ioutil.ReadAll(recipeFile)
		if err != nil {
			return []Recipe{}, fmt.Errorf("failed to read the recipe file: %w", err)
		}

		var r Recipe
		err = json.Unmarshal(data, &r)
		if err != nil {
			return []Recipe{}, fmt.Errorf("failed to unmarshal recipe data: %w", err)
		}

		// create a new standard ID from inconsistent generated IDs
		h := sha1.New()
		h.Write([]byte(r.SourceID))
		r.ID = hex.EncodeToString(h.Sum(nil))

		recipes = append(recipes, r)
	}

	return recipes, nil
}
