package model

type Question struct {
	Id          string
	UserId      string
	Link        string
	CreatedTime string
}
type RespQuestionDb struct {
	LabelId string `db:"label_id"`
	Id      string `db:"guid"`
	EntryId string `db:"entry_id"`
	Name    string `db:"name"`
	Choice  string `db:"choice"`
}
type RespQuestion struct {
	Id      string
	EntryId string
	Name    string
	Choices []Choices
}
