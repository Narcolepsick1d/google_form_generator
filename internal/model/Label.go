package model

type Label struct {
	Id         string
	Entry      string
	Name       string
	QuestionId string
}
type UpdateLabel struct {
	Entry string
	Name  string
}
