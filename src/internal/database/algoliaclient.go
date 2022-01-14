package database

import (
	"log"
	"os"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
)

// AlgoliaIndex ...
func AlgoliaIndex() *search.Index {
	var algoliaIndexName string
	if os.Getenv("AWS_ENV") == "development" {
		algoliaIndexName = os.Getenv("DEV_ALGOLIA_INDEX")
	} else if os.Getenv("AWS_ENV") == "production" {
		algoliaIndexName = os.Getenv("PROD_ALGOLIA_INDEX")
	}

	// Setting up algolia index
	client := search.NewClient(os.Getenv("ALGOLIA_APP_ID"), os.Getenv("ALGOLIA_API_KEY"))
	algoliaIndex := client.InitIndex(algoliaIndexName)

	_, err := algoliaIndex.SetSettings(search.Settings{
		SearchableAttributes: opt.SearchableAttributes(
			"ContentType",
			"Tags",
			"Title",
			"Url",
		),
	})

	if err != nil {
		log.Printf("Error while setting algolia index: %v ", err)
	}

	return algoliaIndex
}
