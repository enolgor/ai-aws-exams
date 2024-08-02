package db

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/enolgor/ai-aws-exams/model"
	"github.com/stephenafamo/scan"
	"github.com/stephenafamo/scan/stdscan"
)

type DB struct {
	conn *sql.DB
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.conn.QueryContext(ctx, query, args...)
}

func (db *DB) Initialize(ctx context.Context) error {
	certifications, err := load_initial_certifications()
	if err != nil {
		return err
	}
	return db.insertInitialData(ctx, certifications)
}

func (db *DB) insertInitialData(ctx context.Context, certifications []model.Certification) error {
	for _, certification := range certifications {
		certification_id, err := db.insertCertification(ctx, certification)
		if err != nil {
			return err
		}
		for _, domain := range certification.Domains {
			domain.CertificationID = certification_id
			if err = db.insertDomain(ctx, *domain); err != nil {
				return err
			}
			for _, task := range domain.Tasks {
				task.CertificationID = certification_id
				task.DomainNumber = domain.Number
				if err = db.insertTask(ctx, *task); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (db *DB) insertCertification(ctx context.Context, certification model.Certification) (int, error) {
	query := sq.Insert("certifications").Columns("name", "description").Values(certification.Name, certification.Description).Suffix(`RETURNING "id"`)
	sql, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}
	return stdscan.One(ctx, db, scan.SingleColumnMapper[int], sql, args...)
}

func (db *DB) insertDomain(ctx context.Context, domain model.Domain) error {
	query := sq.Insert("domains").Columns("certification_id", "number", "name").Values(domain.CertificationID, domain.Number, domain.Name)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = db.conn.ExecContext(ctx, sql, args...)
	return err
}

func (db *DB) insertTask(ctx context.Context, task model.Task) error {
	query := sq.Insert("tasks").Columns("certification_id", "domain_number", "number", "name").Values(task.CertificationID, task.DomainNumber, task.Number, task.Name)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = db.conn.ExecContext(ctx, sql, args...)
	return err
}

func (db *DB) ListCertifications(ctx context.Context) ([]*model.Certification, error) {
	query := sq.Select("*").From("certifications")
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	return stdscan.All(ctx, db, scan.StructMapper[*model.Certification](), sql, args...)
}

func (db *DB) GetCertificationByName(ctx context.Context, name string) (*model.Certification, error) {
	query := sq.Select("*").From("certifications").Where(sq.Eq{"name": name}).Limit(1)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	return stdscan.One(ctx, db, scan.StructMapper[*model.Certification](), sql, args...)
}

func (db *DB) InsertQuestion(ctx context.Context, question model.Question) (int, error) {
	query := sq.Insert("questions").Columns("certification_id", "domain_number", "task_number", "question").Values(question.CertificationID, question.Domain, question.Task, question.Question).Suffix(`RETURNING "id"`)
	sql, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}
	return stdscan.One(ctx, db, scan.SingleColumnMapper[int], sql, args...)
}

func (db *DB) InsertAnswer(ctx context.Context, answer model.Answer) error {
	query := sq.Insert("answers").Columns("question_id", "answer", "correct", "explanation").Values(answer.QuestionID, answer.Answer, answer.Correct, answer.Explanation)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = db.conn.ExecContext(ctx, sql, args...)
	return err
}

func (db *DB) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := sq.Select("*").From("users").Where(sq.Eq{"email": email}).Limit(1)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	return stdscan.One(ctx, db, scan.StructMapper[*model.User](), sql, args...)
}

func (db *DB) InsertResponse(ctx context.Context, response model.Response) error {
	query := sq.Insert("responses").Columns("user_id", "question_id", "answer_id", "date").Values(response.UserID, response.QuestionID, response.AnswerID, response.Date)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = db.conn.ExecContext(ctx, sql, args...)
	return err
}

func (db *DB) GetDomainsByCertification(ctx context.Context, certificationID int) ([]*model.Domain, error) {
	query := sq.Select("*").From("domains").Where(sq.Eq{"certification_id": certificationID})
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	return stdscan.All(ctx, db, scan.StructMapper[*model.Domain](), sql, args...)
}

func (db *DB) GetTasksByCertificationAndDomain(ctx context.Context, certificationID int, domain int) ([]*model.Task, error) {
	query := sq.Select("*").From("tasks").Where(sq.Eq{"certification_id": certificationID, "domain_number": domain})
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	return stdscan.All(ctx, db, scan.StructMapper[*model.Task](), sql, args...)
}

func (db *DB) GetQuestions(ctx context.Context, certificate_id int, domain int, task int) ([]*model.Question, error) {
	query := sq.Select("*").From("questions").Where(sq.Eq{"certification_id": certificate_id, "domain_number": domain, "task_number": task})
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	return stdscan.All(ctx, db, scan.StructMapper[*model.Question](), sql, args...)
}

func (db *DB) GetResponsesByUserID(ctx context.Context, userID int) ([]*model.Response, error) {
	query := sq.Select("*").From("responses").Where(sq.Eq{"user_id": userID})
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	return stdscan.All(ctx, db, scan.StructMapper[*model.Response](), sql, args...)
}

func (db *DB) GetAnswersByQuestion(ctx context.Context, questionID int) ([]*model.Answer, error) {
	query := sq.Select("*").From("answers").Where(sq.Eq{"question_id": questionID})
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	return stdscan.All(ctx, db, scan.StructMapper[*model.Answer](), sql, args...)
}

func (db *DB) GetFullCertificationByID(ctx context.Context, certification int) (*model.Certification, error) {
	query := sq.Select("*").From("certifications").Where(sq.Eq{"id": certification}).Limit(1)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	cert, err := stdscan.One(ctx, db, scan.StructMapper[*model.Certification](), sql, args...)
	if err != nil {
		return nil, err
	}
	domains, err := db.GetDomainsByCertification(ctx, certification)
	if err != nil {
		return nil, err
	}
	cert.Domains = domains
	for _, domain := range domains {
		tasks, err := db.GetTasksByCertificationAndDomain(ctx, certification, domain.Number)
		if err != nil {
			return nil, err
		}
		domain.Tasks = tasks
		for _, task := range tasks {
			questions, err := db.GetQuestions(ctx, certification, domain.Number, task.Number)
			if err != nil {
				return nil, err
			}
			task.Questions = questions
			for _, question := range questions {
				answers, err := db.GetAnswersByQuestion(ctx, question.ID)
				if err != nil {
					return nil, err
				}
				question.Answers = answers
			}
		}
	}
	return cert, nil
}

func (db *DB) InsertUser(ctx context.Context, user model.User) (int, error) {
	query := sq.Insert("users").Columns("email").Values(user.Email).Suffix(`RETURNING "id"`)
	sql, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}
	return stdscan.One(ctx, db, scan.SingleColumnMapper[int], sql, args...)
}

func (db *DB) ClearResponses(ctx context.Context, userID int, certificationID int) error {
	sql := `DELETE FROM responses
					WHERE user_id = ? AND question_id IN (
							SELECT q.id
							FROM questions q
							WHERE q.certification_id = ?
					);`
	_, err := db.conn.ExecContext(ctx, sql, userID, certificationID)
	return err
}

func (db *DB) GetRandomUnansweredQuestion(ctx context.Context, userID int, certificationID int) (*model.Question, error) {
	query := `SELECT q.id, q.certification_id, q.domain_number, q.task_number, q.question
					FROM questions q
					LEFT JOIN responses r ON q.id = r.question_id AND r.user_id = ?
					WHERE q.certification_id = ? AND r.question_id IS NULL
					ORDER BY RANDOM()
					LIMIT 1;
	`
	question, err := stdscan.One(ctx, db, scan.StructMapper[*model.Question](), query, userID, certificationID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	answers, err := db.GetAnswersByQuestion(context.Background(), question.ID)
	if err != nil {
		return nil, err
	}
	question.Answers = answers
	return question, err
}

func (db *DB) GetRandomIncorrectQuestion(ctx context.Context, userID int, certificationID int) (*model.Question, error) {
	query := `SELECT q.id, q.certification_id, q.domain_number, q.task_number, q.question
						FROM questions q
						JOIN responses r ON q.id = r.question_id
						JOIN answers a ON r.answer_id = a.id
						WHERE r.user_id = ? 
							AND q.certification_id = ? 
							AND a.correct = 0
							AND q.id NOT IN (
									SELECT q.id
									FROM questions q
									JOIN responses r ON q.id = r.question_id
									JOIN answers a ON r.answer_id = a.id
									WHERE r.user_id = ? 
										AND q.certification_id = ? 
										AND a.correct = 1
							)
						ORDER BY RANDOM()
						LIMIT 1;
	`
	question, err := stdscan.One(ctx, db, scan.StructMapper[*model.Question](), query, userID, certificationID, userID, certificationID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	answers, err := db.GetAnswersByQuestion(context.Background(), question.ID)
	if err != nil {
		return nil, err
	}
	question.Answers = answers
	return question, err
}

func (db *DB) CountCorrectQuestions(ctx context.Context, userID int, certificationID int) (int, error) {
	query := `SELECT COUNT(DISTINCT q.id) AS correctly_answered_count
						FROM questions q
						JOIN responses r ON q.id = r.question_id
						JOIN answers a ON r.answer_id = a.id
						WHERE r.user_id = ? 
							AND q.certification_id = ? 
							AND a.correct = 1;
	`
	count, err := stdscan.One(ctx, db, scan.SingleColumnMapper[int], query, userID, certificationID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return count, err
}

func (db *DB) CountIncorrectQuestions(ctx context.Context, userID int, certificationID int) (int, error) {
	query := `SELECT COUNT(DISTINCT q.id) AS incorrectly_answered_count
						FROM questions q
						JOIN responses r ON q.id = r.question_id
						JOIN answers a ON r.answer_id = a.id
						WHERE r.user_id = ? 
							AND q.certification_id = ? 
							AND q.id NOT IN (
									SELECT q.id
									FROM questions q
									JOIN responses r ON q.id = r.question_id
									JOIN answers a ON r.answer_id = a.id
									WHERE r.user_id = ? 
										AND q.certification_id = ? 
										AND a.correct = 1
							);
	`
	count, err := stdscan.One(ctx, db, scan.SingleColumnMapper[int], query, userID, certificationID, userID, certificationID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return count, err
}

func (db *DB) CountQuestions(ctx context.Context, certificationID int) (int, error) {
	query := `SELECT COUNT(*) AS total_questions_count
						FROM questions
						WHERE certification_id = ?;
	`
	return stdscan.One(ctx, db, scan.SingleColumnMapper[int], query, certificationID)
}

func (db *DB) GetLastResponse(ctx context.Context, user_id int, question_id int) (*model.Response, error) {
	query := `SELECT r.*
						FROM responses r
						WHERE r.question_id = ? AND r.user_id = ?
						ORDER BY r.date DESC
						LIMIT 1;
	`
	response, err := stdscan.One(ctx, db, scan.StructMapper[*model.Response](), query, question_id, user_id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return response, err
}
