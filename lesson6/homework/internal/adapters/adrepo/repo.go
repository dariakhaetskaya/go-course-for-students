package adrepo

import (
	"context"
	"homework6/internal/ads"
	"homework6/internal/app"
)

type AdRepo struct {
	ads map[int64]*ads.Ad
}

func (a AdRepo) GetId(ctx context.Context) (int64, error) {
	if e := ctx.Err(); e != nil {
		return 0, e
	}
	return int64(len(a.ads)), nil
}

func (a AdRepo) CreateAd(ctx context.Context, ad *ads.Ad) error {
	if e := ctx.Err(); e != nil {
		return e
	}

	a.ads[ad.ID] = ad
	return nil
}

func (a AdRepo) ReadAd(ctx context.Context, id int64) (*ads.Ad, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}

	if ad, ok := a.ads[id]; !ok {
		return nil, app.ErrNotFound
	} else {
		return ad, nil
	}
}

func (a AdRepo) UpdateAd(ctx context.Context, id int64, ad *ads.Ad) error {
	if e := ctx.Err(); e != nil {
		return e
	}

	if _, ok := a.ads[id]; !ok {
		return app.ErrNotFound
	} else {
		a.ads[id] = ad
		return nil
	}
}

func New() app.Repository {
	return &AdRepo{
		make(map[int64]*ads.Ad),
	}
}
