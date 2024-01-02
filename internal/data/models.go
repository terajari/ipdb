package data

import "database/sql"

type Models struct {
	Podcast IPodcast
}

func NewModels(db *sql.DB) Models {
	return Models{
		Podcast: NewPodcastModel(db),
	}
}
