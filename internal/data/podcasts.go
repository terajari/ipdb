package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/terajari/ipdb/internal/validator"
)

type Podcast struct {
	Id            int64     `json:"id"`
	Title         string    `json:"title"`
	Platform      string    `json:"platform"`
	Url           string    `json:"url"`
	Host          string    `json:"host"`
	Program       string    `json:"program"`
	GuestSpeakers []string  `json:"guest_speakers"`
	Year          int64     `json:"year"`
	Language      string    `json:"language"`
	Tags          []string  `json:"tags"`
	CreatedAt     time.Time `json:"created_at"`
}

type PodcastModel struct {
	Db *sql.DB
}

type IPodcast interface {
	Insert(*Podcast) error
	FindById(int64) (*Podcast, error)
	GetPodcasts() ([]*Podcast, error)
	UpdatePodcast(*Podcast) error
	DeleteById(int64) error
	GetAll(string, []string, Filters) (*[]Podcast, Metadata, error)
}

func NewPodcastModel(db *sql.DB) IPodcast {
	return &PodcastModel{Db: db}
}

func ValidatePodcast(v *validator.Validator, podcast *Podcast) {
	v.Check(podcast.Title != "", "title", "must be provided")
	v.Check(len(podcast.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(podcast.Platform != "", "platform", "must be provided")
	v.Check(len(podcast.Platform) <= 500, "platform", "must not be more than 500 bytes long")
	v.Check(podcast.Url != "", "url", "must be provided")
	v.Check(len(podcast.Url) <= 500, "url", "must not be more than 500 bytes long")
	v.Check(podcast.Host != "", "host", "must be provided")
	v.Check(len(podcast.Host) <= 500, "host", "must not be more than 500 bytes long")
	v.Check(podcast.Program != "", "program", "must be provided")
	v.Check(podcast.Language != "", "language", "must be provided")
	v.Check(len(podcast.Language) <= 500, "language", "must not be more than 500 bytes long")
	v.Check(podcast.Year != 0, "year", "must be provided")
	v.Check(podcast.Year >= 2003, "year", "must be greater than 2003")
	v.Check(podcast.Year <= int64(time.Now().Year()), "year", "must not be in the future")
	v.Check(len(podcast.Tags) >= 1, "tags", "must contain at least 1 tag")
	v.Check(len(podcast.Tags) <= 10, "tags", "must not contain more than 10 tags")
	v.Check(validator.Unique[string](podcast.Tags...), "tags", "must not contain duplicate tags")
	v.Check(len(podcast.GuestSpeakers) >= 1, "guest_speakers", "must contain at least 1 guest_speaker")
	v.Check(len(podcast.GuestSpeakers) <= 10, "guest_speakers", "must not contain more than 10 guest_speakers")
	v.Check(validator.Unique[string](podcast.GuestSpeakers...), "guest_speakers", "must not contain duplicate guest_speakers")
}

func (pm PodcastModel) Insert(podcast *Podcast) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO podcasts 
		(title, platform, url, host, program, guest_speakers, year, language, tags)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	args := []any{
		podcast.Title,
		podcast.Platform,
		podcast.Url,
		podcast.Host,
		podcast.Program,
		pq.Array(podcast.GuestSpeakers),
		podcast.Year,
		podcast.Language,
		pq.Array(podcast.Tags),
	}
	return pm.Db.QueryRowContext(ctx, query, args...).Scan(&podcast.Id, &podcast.CreatedAt)
}

func (pm PodcastModel) FindById(id int64) (*Podcast, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, title, platform, url, host, program, guest_speakers, year, language, tags, created_at 
		FROM podcasts
		WHERE id = $1
	`

	var podcast Podcast
	if err := pm.Db.QueryRowContext(ctx, query, id).Scan(
		&podcast.Id,
		&podcast.Title,
		&podcast.Platform,
		&podcast.Url,
		&podcast.Host,
		&podcast.Program,
		pq.Array(&podcast.GuestSpeakers),
		&podcast.Year,
		&podcast.Language,
		pq.Array(&podcast.Tags),
		&podcast.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &podcast, nil
}

func (pm PodcastModel) UpdatePodcast(podcast *Podcast) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE podcasts
		SET title = $1, platform = $2, url = $3, host = $4, program = $5, guest_speakers = $6, year = $7, language = $8, tags = $9
		WHERE id = $10
		RETURNING id, created_at
	`
	args := []any{
		podcast.Title,
		podcast.Platform,
		podcast.Url,
		podcast.Host,
		podcast.Program,
		pq.Array(podcast.GuestSpeakers),
		podcast.Year,
		podcast.Language,
		pq.Array(podcast.Tags),
		podcast.Id,
	}

	return pm.Db.QueryRowContext(ctx, query, args...).Scan(&podcast.Id, &podcast.CreatedAt)
}

func (pm PodcastModel) GetPodcasts() ([]*Podcast, error) {
	return nil, nil
}

func (pm PodcastModel) DeleteById(id int64) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		DELETE FROM podcasts WHERE id = $1
	`

	_, err := pm.Db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (pm PodcastModel) GetAll(platform string, tags []string, filters Filters) (*[]Podcast, Metadata, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, title, platform, url, host, program, guest_speakers, year, language, tags, created_at
		FROM podcasts
		WHERE (to_tsvector('simple', platform) @@ plainto_tsquery('simple', $1) OR $1 = '')
		OR (tags @> $2 OR $2 = '{}')
		ORDER BY id
		LIMIT $3 OFFSET $4
	`

	rows, err := pm.Db.QueryContext(ctx, query, platform, pq.Array(tags), filters.Limit(), filters.Offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	podcasts := []Podcast{}

	for rows.Next() {
		var podcast Podcast
		if err := rows.Scan(
			&podcast.Id,
			&podcast.Title,
			&podcast.Platform,
			&podcast.Url,
			&podcast.Host,
			&podcast.Program,
			pq.Array(&podcast.GuestSpeakers),
			&podcast.Year,
			&podcast.Language,
			pq.Array(&podcast.Tags),
			&podcast.CreatedAt,
		); err != nil {
			return nil, Metadata{}, err
		}
		podcasts = append(podcasts, podcast)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(len(podcasts), filters.Page, filters.PageSize)

	return &podcasts, metadata, nil
}
