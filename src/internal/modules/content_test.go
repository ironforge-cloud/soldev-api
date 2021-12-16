package modules

import (
	"api/internal/types"
	"testing"
)

func TestReviewNewContent(t *testing.T) {
	data := []types.Content{
		{
			PK:            "Rust#courses",
			SK:            "keep-new",
			ContentStatus: "active",
			Tags:          nil,
			SpecialTag:    "New",
			Url:           "https://solhack.com/courses/how-rusty-is-your-rust-lang/",
			Title:         "How Rusty is Your Rust Lang?",
			Author:        "Solhack",
			ContentType:   "courses",
			Vertical:      "Rust",
			PublishedAt:   "1638196466076",
		},
		{
			PK:            "Rust#courses",
			SK:            "no-date",
			ContentStatus: "active",
			Tags:          nil,
			SpecialTag:    "New",
			Url:           "https://solhack.com/courses/how-rusty-is-your-rust-lang/",
			Title:         "How Rusty is Your Rust Lang?",
			Author:        "Solhack",
			ContentType:   "courses",
			Vertical:      "Rust",
		},
		{
			PK:            "Rust#courses",
			SK:            "remove-new",
			ContentStatus: "active",
			Tags:          nil,
			SpecialTag:    "New",
			Url:           "https://solhack.com/courses/how-rusty-is-your-rust-lang/",
			Title:         "How Rusty is Your Rust Lang?",
			Author:        "Solhack",
			ContentType:   "courses",
			Vertical:      "Rust",
			PublishedAt:   "1257894000",
		},
	}

	// Review content
	contentList, err := ReviewNewContent(data)
	if err != nil {
		t.Error("error reviewing content", err)
	}

	for _, item := range contentList {
		if item.SK == "keep-new" {
			if item.SpecialTag != "New" {
				t.Errorf("expected %q but got %q", "New", item.SpecialTag)
			}
		} else if item.SK == "remove-new" {
			if item.SpecialTag != "" {
				t.Errorf("expected %q but got %q", "", item.SpecialTag)
			}
		} else if item.SK == "no-date" {
			if item.SpecialTag != "New" {
				t.Errorf("expected %q but got %q", "New", item.SpecialTag)
			}
		}
	}
}
