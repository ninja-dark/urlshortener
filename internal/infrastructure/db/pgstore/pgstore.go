package pgstore

import (
	"context"
	"database/sql"
	"time"

	"urlsortener/internal/entities/urlentities"

	"urlsortener/internal/entities/statsentities"

	"urlsortener/internal/usecases/app"


	_ "github.com/jackc/pgx/v4/stdlib" // Postgresql driver
)

var _ app.URLStore = &UrlStore{}

type DBPgURL struct {
	ID          int
	CreatedAt   time.Time
	OriginalURL string
	ShortURL    string
	Redirects   int
}

type DBPgURLStat struct {
	ShortURL    string
	Redirects   int
}


type UrlStore struct{
	db *sql.DB
}

func NewStore(dsn string) (*UrlStore, error){
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE OF NOT EXISTS urls (
		id uuid NIT NULL,
		created_at timestamptz NOT NULL,
		origina_url varchar,
		short_url varchar,
		redirects integer
	) `)
	if err != nil {
		db.Close()
		return nil, err
	}
	u := &UrlStore{
		db: db,
	}
	return u, nil

}

func (s *UrlStore) Close() error {
	return s.db.Close()
}

func (s *UrlStore)Create(ctx context.Context, origURL string) (*urlentities.URL, error){
	dbUrl := &DBPgURL{
		CreatedAt: time.Now(),
		OriginalURL: origURL,
	}
	var id int

	err:= s.db.QueryRowContext(ctx, `INSERT INTO urls (created_at, original_url) values ($1, $2) RETURNING id`, dbUrl.CreatedAt, dbUrl.OriginalURL).Scan(&id)
	
	if err != nil{
		return nil, err
	}

	return &urlentities.URL{
		ID: id,
		OriginalURL: dbUrl.OriginalURL,
	}, nil
}

func (s *UrlStore) Update(ctx context.Context, url *urlentities.URL) error {
	dbUrl := &DBPgURL{
		ID: url.ID,
		ShortURL: url.ShortUrl,
	}

	_, err := s.db.ExecContext(ctx, `UPDATE urls SET short_url = $2 WHERE id = $1`,
		dbUrl.ID, dbUrl.ShortURL)
	

	if err != nil{
		return  err
	}
	return nil
}

func (s *UrlStore) GetOrUrl(ctx context.Context, shortURL string) (*urlentities.URL, error){
	dbUrl := &DBPgURL{}
	row, err := s.db.QueryContext(ctx, `SELECT id, created_at, original_url, short_url, redirects FROM urls WHERE short_url &1`, shortURL)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next(){
		if err := row.Scan(
			&dbUrl.ID,
			&dbUrl.CreatedAt,
			&dbUrl.OriginalURL,
			&dbUrl.ShortURL,
			&dbUrl.Redirects,
		); err != nil {
			return nil, err
		}
	}

	return &urlentities.URL{
		ID: dbUrl.ID,
		OriginalURL: dbUrl.OriginalURL,
		ShortUrl: dbUrl.ShortURL,
		Count: dbUrl.Redirects,
	}, nil
}

func (s *UrlStore) Stats(ctx context.Context, shortURL string) (*statsentities.Stats, error){
	urlStats:= &DBPgURLStat{}

	row, err := s.db.QueryContext(ctx, `SELECT short_url, redirects FROM urls WHERE short_url &1`, shortURL)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next(){
		if err := row.Scan(
			&urlStats.ShortURL,
			urlStats.Redirects,
		); err != nil {
			return nil, err
		}
	}
	return &statsentities.Stats{
		ShortURL: urlStats.ShortURL,
		NumberRedirect: urlStats.Redirects,
	}, nil
}

func (s *UrlStore)CountRedirects(ctx context.Context, shortURL string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE urls SET redirects = redirects + 1  WHERE shortURL = $1`, shortURL)
	if err != nil{
		return  err
	}
	return nil
}