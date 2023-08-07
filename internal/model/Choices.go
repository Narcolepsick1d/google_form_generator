package model

type Choices struct {
	Id          string
	Choice      string
	Probability string
	LabelId     string
}
type UpdateChoices struct {
	Id          string `db:"guid"`
	Choice      string `db:"choice"`
	Probability string `db:"probability"`
}
