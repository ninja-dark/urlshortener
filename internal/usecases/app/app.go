package app

import (
	"context"
	"fmt"

	"urlsortener/internal/entities/urlentities" 
	"urlsortener/internal/usecases/app/generator"
	"urlsortener/internal/entities/statsentities" 


	
)

type URLStore interface { 
	Create(ctx context.Context, originalURL string) (*urlentities.URL, error)
	Update(ctx context.Context, url *urlentities.URL) error
	GetOrUrl(ctx context.Context, shortURL string) (*urlentities.URL, error)
	Stats(ctx context.Context, shortURL string) (*statsentities.Stats, error)
	CountRedirects(ctx context.Context, shortURL string)  error
}

type Urls struct{
	store URLStore
}

func NewUrls( store URLStore) *Urls {
	return &Urls{
		store:store,
	}

}

func (u *Urls) CreateURL(ctx context.Context, url string) (*urlentities.URL, error) {
	url1, err := u.store.Create(ctx, url)
	if err != nil{
		return nil, fmt.Errorf("error: %w" ,err)
	}
	
	short := generator.Encode(int(url1.ID))

	url1.ShortUrl = short

	if err = u.store.Update(ctx, url1); err !=nil{
		return nil, fmt.Errorf("can't save URL error %w", err)
	}
	return url1, nil

}

func(u *Urls) RedirectUrl(ctx context.Context, shortURL string) (*urlentities.URL, error){
	url, err := u.store.GetOrUrl(ctx, shortURL)
	if err != nil {
		return nil, fmt.Errorf("can't get url %w", err)
	}
	u.CountRedirects(ctx, shortURL)
	
	return url, nil
}

func(u *Urls) Stats(ctx context.Context, shortURL string) (*statsentities.Stats, error) {
	s, err := u.store.Stats(ctx, shortURL)
	if err != nil{
		return nil, fmt.Errorf("can't get url %w", err)
	}
	return s, nil
}

func(u *Urls) CountRedirects(ctx context.Context, shortURL string){
	e := u.store.CountRedirects(ctx, shortURL)
	if e != nil {
		fmt.Println("error")
	}
}