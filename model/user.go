package model

import "time"

type User struct {
	ID    int    `json:"id" db:"id"`
	Email string `json:"email" db:"email"`
}

type Response struct {
	UserID     int       `json:"user_id" db:"user_id"`
	QuestionID int       `json:"question_id" db:"question_id"`
	AnswerID   int       `json:"answer_id" db:"answer_id"`
	Date       time.Time `json:"date" db:"date"`
}
