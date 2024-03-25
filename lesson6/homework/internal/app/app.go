package app

import (
	"context"
	"homework6/internal/ads"
	"homework6/internal/validator"
)

type App interface {
	CreateAd(ctx context.Context, title string, text string, authorID int64) (*ads.Ad, error)
	GetAdByID(ctx context.Context, id int64) (*ads.Ad, error)
	ChangeAdStatus(ctx context.Context, id int64, userID int64, published bool) error
	UpdateAd(ctx context.Context, id int64, userID int64, title string, text string) error
}

type AdApp struct {
	repository Repository
}

type Repository interface {
	CreateAd(ctx context.Context, ad *ads.Ad) error
	ReadAd(ctx context.Context, id int64) (*ads.Ad, error)
	UpdateAd(ctx context.Context, id int64, ad *ads.Ad) error
	GetId(ctx context.Context) (int64, error)
}

func NewApp(repo Repository) App {
	return &AdApp{repository: repo}
}

func (a *AdApp) CreateAd(ctx context.Context, title string, text string, authorID int64) (*ads.Ad, error) {
	id, err := a.repository.GetId(ctx)
	if err != nil {
		return nil, err
	}
	ad := ads.New(id, title, text, authorID)
	err = validator.Validate(*ad)
	if err != nil {
		return nil, ErrInvalid
	}

	err = a.repository.CreateAd(ctx, ad)
	if err != nil {
		return nil, err
	}
	return ad, nil
}

func (a *AdApp) ChangeAdStatus(ctx context.Context, id int64, userID int64, published bool) error {
	ad, err := a.repository.ReadAd(ctx, id)
	if err != nil {
		return ErrNotFound
	}

	if ad.AuthorID != userID {
		return ErrNotAuthor
	}

	ad.Published = published
	return nil
}

func (a *AdApp) UpdateAd(ctx context.Context, id int64, userID int64, title string, text string) error {
	ad, err := a.repository.ReadAd(ctx, id)
	if err != nil {
		return ErrNotFound
	}

	if ad.AuthorID != userID {
		return ErrNotAuthor
	}

	ad.Title = title
	ad.Text = text
	err = validator.Validate(*ad)
	if err != nil {
		return ErrInvalid
	}

	err = a.repository.UpdateAd(ctx, id, ad)
	if err != nil {
		return err
	}
	return nil
}

func (a *AdApp) GetAdByID(ctx context.Context, id int64) (*ads.Ad, error) {
	ad, err := a.repository.ReadAd(ctx, id)
	if err != nil {
		return nil, ErrNotFound
	}
	return ad, nil
}
