package model

type Label struct {
	Id         string
	Entry      string
	Name       string
	QuestionId string
	IsMulti    string
}
type UpdateLabel struct {
	Id      string
	Entry   string
	Name    string
	IsMulti bool
}
