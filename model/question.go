package model

type Question struct {
	ID              int       `json:"id" db:"id"`
	CertificationID int       `json:"certification_id" db:"certification_id"`
	Domain          int       `json:"domain" db:"domain_number"`
	Task            int       `json:"task" db:"task_number"`
	Question        string    `json:"question" db:"question"`
	Answers         []*Answer `json:"answers"`
}

type Answer struct {
	ID          int    `json:"id" db:"id"`
	QuestionID  int    `json:"question_id" db:"question_id"`
	Answer      string `json:"answer" db:"answer"`
	Correct     bool   `json:"correct" db:"correct"`
	Explanation string `json:"explanation" db:"explanation"`
}
