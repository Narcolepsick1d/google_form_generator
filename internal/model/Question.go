package model

type Question struct {
	Id          string
	UserId      string
	Link        string
	CreatedTime string
}
type RespQuestion struct {
	Id     string `db:"guid"`
	Name   string `db:"name"`
	Choice string `db:"choice"`
}
