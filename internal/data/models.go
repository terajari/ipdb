package data

import "database/sql"

type Models struct {
	Podcast IPodcast
	User    IUser
	Token   IToken
}

func NewModels(db *sql.DB) Models {
	return Models{
		Podcast: NewPodcastModel(db),
		User:    NewUserModel(db),
		Token:   NewTokenModel(db),
	}
}
