package model

type Choices struct {
	Id          string
	Choice      string
	Probability string
	LabelId     string
}
type UpdateChoices struct {
	Choice      string
	Probability string
}
