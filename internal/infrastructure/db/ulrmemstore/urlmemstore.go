package ulrmemstore

import (
	"context"
	"database/sql"
	"sync"
	"urlsortener/internal/entities/urlentities"

	"urlsortener/internal/entities/statsentities"

	"urlsortener/internal/usecases/app"
)


var _ app.URLStore = &UrlStore{}

type UrlStore struct {
	sync.Mutex
	originalUrl map[string]urlentities.URL
	shortUrl map[string]urlentities.URL
}

func NewUrlStore() *UrlStore {
	return &UrlStore{
		originalUrl: make(map[string]urlentities.URL),
		shortUrl: make(map[string]urlentities.URL),
	}
}
func (u *UrlStore) Create(ctx context.Context, original string) (*urlentities.URL, error){
	u.Lock()
	defer u.Unlock()

	select {
	case <- ctx.Done():
		return nil, ctx.Err()
	default:
	}

	url := urlentities.URL{
		ID: len(u.shortUrl) + len(u.originalUrl),
		OriginalURL: original,
	}
	
	u.originalUrl[original] = url
	return &url, nil
}

func(u *UrlStore) Update(ctx context.Context, url *urlentities.URL) error{
	u.Lock()
	defer u.Unlock()
	select {
	case <- ctx.Done():
		return  ctx.Err()
	default:
	}


	delete(u.originalUrl, url.OriginalURL)

	u.shortUrl[url.ShortUrl] = *url

	return nil

}

func(u *UrlStore) GetOrUrl(ctx context.Context, shortURL string)(*urlentities.URL, error){
	u.Lock()
	defer u.Unlock()
	select {
	case <- ctx.Done():
		return nil, ctx.Err()
	default:
	}

	url, exists := u.shortUrl[shortURL]
	if exists{
		return &url, nil
	}
	return nil, sql.ErrNoRows
}

func(u *UrlStore) Stats(ctx context.Context, shortURL string) (*statsentities.Stats, error) {
	u.Lock()
	defer u.Unlock()
	select {
	case <- ctx.Done():
		return nil, ctx.Err()
	default:
	}

	url, exists := u.shortUrl[shortURL]
	if exists {
		return &statsentities.Stats{
			ShortURL: shortURL,
			NumberRedirect: url.Count,
		}, nil
	}
	return nil, sql.ErrNoRows

}

func(u *UrlStore)CountRedirects(ctx context.Context, shortURL string) error {
	u.Lock()
	defer u.Unlock()
	select {
	case <- ctx.Done():
		return ctx.Err()
	default:
	}
	existURL, exist := u.shortUrl[shortURL]
	if exist{
		existURL.Count +=1
		u.shortUrl[shortURL] = existURL
		return nil
	}
	return sql.ErrNoRows
}

