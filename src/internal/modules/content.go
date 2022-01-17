package modules

import (
	"api/internal/database"
	"api/internal/types"
	"api/internal/utils"
	"errors"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/teris-io/shortid"
)

// SaveContent ...
func SaveContent(content []types.Content) error {

	contentLength := len(content)
	var wg sync.WaitGroup
	wg.Add(contentLength)

	for i := 0; i < contentLength; i++ {
		go func(i int) error {
			defer wg.Done()

			if len(content[i].PublishedAt) == 0 {
				content[i].PublishedAt = strconv.FormatInt(time.Now().Unix(), 10)
			}

			// if content status is inactive for more than a week, we can delete
			// the content.
			if content[i].ContentStatus == "inactive" {
				content[i].Expdate = time.Now().Add(time.Hour * 168).Unix()
			}

			// If PK data doesn't match Vertical#ContentType we need to
			// delete the old content
			if content[i].PK != content[i].Vertical+"#"+content[i].ContentType {
				err := database.DeleteContent(content[i])

				if err != nil {
					return err
				}
			}

			// Check if content has Img assigned. If not, try to find one using
			// html meta tags.
			if content[i].Img == "" {
				content[i].Img = utils.GetImageIfExist(content[i].Url)
			}

			// Change PK to the new one
			content[i].PK = content[i].Vertical + "#" + content[i].ContentType

			// Saving content
			err := database.SaveContent(content[i])

			if err != nil {
				return err
			}

			return nil
		}(i)
	}

	wg.Wait()

	//for _, item := range content {
	//	if len(item.PublishedAt) == 0 {
	//		item.PublishedAt = strconv.FormatInt(time.Now().Unix(), 10)
	//	}
	//
	//	// if content status is inactive for more than a week, we can delete
	//	// the content.
	//	if item.ContentStatus == "inactive" {
	//		item.Expdate = time.Now().Add(time.Hour * 168).Unix()
	//	}
	//
	//	// If PK data doesn't match Vertical#ContentType we need to
	//	// delete the old content
	//	if item.PK != item.Vertical+"#"+item.ContentType {
	//		err := database.DeleteContent(item)
	//
	//		if err != nil {
	//			return err
	//		}
	//	}
	//
	//	// Check if content has Img assigned. If not, try to find one using
	//	// html meta tags.
	//	if item.Img == "" {
	//		item.Img = utils.GetImageIfExist(item.Url)
	//	}
	//
	//	// Change PK to the new one
	//	item.PK = item.Vertical + "#" + item.ContentType
	//
	//	// Saving content
	//	err := database.SaveContent(item)
	//
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil
}

// CreateContent ...
func CreateContent(content types.Content) error {

	// Data sanitization
	content.PK = content.Vertical + "#" + content.ContentType
	content.SK, _ = shortid.Generate()
	content.PublishedAt = strconv.FormatInt(time.Now().Unix(), 10)
	if content.ContentStatus == "" {
		content.ContentStatus = "submitted"
	}

	// Check if content has Img assigned. If not, try to find one using
	// html meta tags.
	if content.Img == "" && content.ContentType != "Playlist" {
		content.Img = utils.GetImageIfExist(content.Url)
	}

	err := database.SaveContent(content)

	if err != nil {
		return err
	}

	return nil
}

// GetContent ...
// TODO: The sort functionality in this method can be improved. Too many iterations
func GetContent(vertical string, contentType string, tags string, specialTag string) ([]types.Content, error) {

	contentList, videoContent, err := database.QueryContent(vertical, contentType, "", "")

	if err != nil {
		return nil, err
	}

	// If videoContent we don't need to sort. DynamoDB GSI video-gsi is
	// taking care of that already using the sort key
	if videoContent {
		return contentList, nil
	}

	// Sort structure
	type categorizedContent struct {
		best    []types.Content
		hot     []types.Content
		new     []types.Content
		old     []types.Content
		noBadge []types.Content
	}

	var contentByCategory categorizedContent
	for _, content := range contentList {
		if content.SpecialTag == "New" {
			contentByCategory.new = append(contentByCategory.new, content)
		} else if content.SpecialTag == "Hot" {
			contentByCategory.hot = append(contentByCategory.hot, content)
		} else if content.SpecialTag == "Old" {
			contentByCategory.old = append(contentByCategory.old, content)
		} else if content.SpecialTag == "Best" {
			contentByCategory.best = append(contentByCategory.best, content)
		} else {
			contentByCategory.noBadge = append(contentByCategory.noBadge, content)
		}
	}

	var contentFilteredBySpecialTags []types.Content

	// If the request includes special tags we need to filter the data
	if len(specialTag) > 0 {
		specialTagSlice := strings.Split(specialTag, ",")

		for _, tag := range specialTagSlice {
			if tag == "New" {
				contentFilteredBySpecialTags = append(contentFilteredBySpecialTags, contentByCategory.new...)
			} else if tag == "Hot" {
				contentFilteredBySpecialTags = append(contentFilteredBySpecialTags, contentByCategory.hot...)
			} else if tag == "Old" {
				contentFilteredBySpecialTags = append(contentFilteredBySpecialTags, contentByCategory.old...)
			} else if tag == "Best" {
				contentFilteredBySpecialTags = append(contentFilteredBySpecialTags, contentByCategory.best...)
			}
		}
	} else {
		contentFilteredBySpecialTags = mergeContent(contentByCategory.best, contentByCategory.new, contentByCategory.hot, contentByCategory.noBadge, contentByCategory.old)
	}

	// Content needs to be sorted using Position. Position in this scope means weight,
	// in other words higher number goes first.
	contentSortedByPosition := sortContentByPosition(contentFilteredBySpecialTags)

	// If there are no tags to filter by, we can return here
	if len(tags) == 0 {
		return contentSortedByPosition, nil
	}

	var contentResponse []types.Content

	// If tags is not empty I need to filter QueryContent response
	if len(tags) > 0 {
		tagsSlice := strings.Split(tags, ",")

		for _, content := range contentSortedByPosition {
			for _, requestTag := range tagsSlice {
				// if content is added we need to break this range too
				shouldBreak := false

				// If item includes
				for _, contentTag := range content.Tags {
					if contentTag == requestTag {
						contentResponse = append(contentResponse, content)
						shouldBreak = true
						break
					}
				}
				if shouldBreak {
					break
				}
			}
		}
	}

	return contentResponse, nil
}

// GetContentByID finds a specific content in the DB
func GetContentByID(vertical string, contentType string, id string) (types.Content, error) {
	content, err := database.GetContentDetails(vertical+"#"+contentType, id)

	if err != nil {
		return types.Content{}, err
	}

	return content, nil
}

// Live finds content labeled as live in the DB
func Live(vertical string) ([]types.Content, error) {
	contentList, err := database.GetLive(vertical)
	if err != nil {
		return nil, err
	}

	return contentList, nil
}

// Promoted finds content labeled as promoted in the DB
func Promoted(vertical string) ([]types.Content, error) {
	contentList, err := database.GetPromoted(vertical)
	if err != nil {
		return nil, err
	}

	return contentList, nil
}

// GetContentByStatus ...
func GetContentByStatus(status string) ([]types.Content, error) {
	content, err := database.QueryContentByStatus(status)
	if err != nil {
		return nil, err
	}

	if content == nil {
		return nil, errors.New("404")
	}

	return content, nil
}

// ReviewNewContent finds content with the specialTag "New" and checks if should continue having it or not
func ReviewNewContent(content []types.Content) ([]types.Content, error) {

	// Validate that the new content PublishedAt timestamp is not 1 month old
	for index, item := range content {
		// It's possible some old content doesn't have a date
		if len(item.PublishedAt) == 0 {
			continue
		}

		// PublishedAt is a string, we need to convert it first
		contentDate, _ := strconv.ParseInt(item.PublishedAt, 10, 64)
		// Adding 15 days to the content date
		limitDate := time.Unix(contentDate, 0).Add(time.Hour * 168).Unix()

		// Check if we need to remove to update the DB record to remove the tag
		if limitDate < time.Now().Unix() {
			content[index].SpecialTag = ""
		}
	}

	return content, nil
}

// DoesContentExist checks if there are content in the database with the
// specified url
func DoesContentExist(url string) (bool, error) {

	contentList, err := database.GetContentByUrl(url)

	if err != nil {
		return false, err
	}

	// If there are items in the response, then, there's content in
	// associated with that url
	if len(contentList) > 0 {
		return true, nil
	}

	return false, nil
}

// FilterContentByList ...
func FilterContentByList(filter string) ([]types.Content, error) {
	// Get list of content with specialTag
	queryResponse, err := database.ScanContent()

	if err != nil {
		return nil, err
	}

	var content []types.Content
	for _, item := range queryResponse {
		if item.Lists != filter {
			continue
		}

		content = append(content, item)
	}

	content = sortContentByPosition(content)

	return content, nil
}

// mergeContent ...
func mergeContent(args ...[]types.Content) []types.Content {
	mergedSlice := make([]types.Content, 0)
	for _, oneSlice := range args {
		mergedSlice = append(mergedSlice, oneSlice...)
	}

	return mergedSlice
}

// sortContentByPosition ...
func sortContentByPosition(contentList []types.Content) []types.Content {
	data := contentList

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].Position > data[j].Position
	})

	return data
}
