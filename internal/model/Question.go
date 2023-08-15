package model

type Question struct {
	Id          string
	UserId      string
	Link        string
	CreatedTime string
}
type RespQuestionDb struct {
	Url         string `db:"url"`
	LabelId     string `db:"label_id"`
	Id          string `db:"guid"`
	EntryId     string `db:"entry_id"`
	Name        string `db:"name"`
	Choice      string `db:"choice"`
	Probability int    `db:"probability"`
	IsMulti     bool   `db:"is_multi"`
}
type RespQuestion struct {
	Id      string
	EntryId string
	Name    string
	IsMulti bool
	Choices []Choices
	Url     string
}
