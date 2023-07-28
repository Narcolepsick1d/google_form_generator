package model

type Payment struct {
	Id         string  `db:"guid"`
	UserId     string  `db:"user_id"`
	Status     string  `db:"status"`
	CreateAt   string  `db:"create_at"`
	Amount     float64 `db:"amount"`
	Quantity   int     `db:"quantity"`
	QuestionId string  `db:"question_id"`
}
