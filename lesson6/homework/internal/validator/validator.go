package validator

import (
	"errors"
	"homework6/internal/ads"
)

func Validate(ad ads.Ad) error {
	if ad.Title == "" {
		return errors.New("title is empty")
	}

	if len(ad.Title) > 100 {
		return errors.New("title should be shorter than 100 symbols")
	}

	if ad.Text == "" {
		return errors.New("text is empty")
	}

	if len(ad.Text) > 500 {
		return errors.New("text should be shorter than 500 symbols")
	}

	return nil
}
