package model

type Certification struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Domains     []*Domain `json:"domains"`
}

type Domain struct {
	CertificationID int     `json:"certification" db:"certification_id"`
	Number          int     `json:"number" db:"number"`
	Name            string  `json:"name" db:"name"`
	Tasks           []*Task `json:"tasks"`
}

type Task struct {
	CertificationID int         `json:"certification_id" db:"certification_id"`
	DomainNumber    int         `json:"domain" db:"domain_number"`
	Number          int         `json:"number" db:"number"`
	Name            string      `json:"name" db:"name"`
	Questions       []*Question `json:"questions"`
}
