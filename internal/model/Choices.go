package model

type Choices struct {
	Id          string
	Choice      string
	Probability int
	LabelId     string
}
type UpdateChoices struct {
	Id          string `db:"guid"`
	Choice      string `db:"choice"`
	Probability int    `db:"probability"`
}
