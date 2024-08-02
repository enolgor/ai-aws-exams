package app

import (
	"context"
	"errors"
	"math/rand/v2"

	"github.com/enolgor/ai-aws-exams/db"
	"github.com/enolgor/ai-aws-exams/model"
)

type App struct {
	db            *db.DB
	user          int
	certification model.Certification
}

func NewApp(db *db.DB, user int, certification int) (*App, error) {
	cerification, err := db.GetFullCertificationByID(context.Background(), certification)
	if err != nil {
		return nil, err
	}
	return &App{db, user, *cerification}, nil
}

func (app *App) GetCertification() model.Certification {
	return app.certification
}

func (app *App) GetNextRandomQuestion() (*model.Domain, *model.Task, *model.Question, *model.Response, error) {
	var previous_response *model.Response
	question, err := app.db.GetRandomUnansweredQuestion(context.Background(), app.user, app.certification.ID)
	if err != nil {
		return nil, nil, nil, previous_response, err
	}
	if question == nil {
		question, previous_response, err = app.getRandomIncorrectQuestion()
	}
	if err != nil {
		return nil, nil, nil, previous_response, err
	}
	if question == nil {
		return nil, nil, nil, previous_response, nil
	}
	rand.Shuffle(len(question.Answers), func(i, j int) { question.Answers[i], question.Answers[j] = question.Answers[j], question.Answers[i] })
	domain, task := app.getDomainAndTaskFromQuestion(question.ID)
	return domain, task, question, previous_response, nil
}

func (app *App) getRandomIncorrectQuestion() (*model.Question, *model.Response, error) {
	question, err := app.db.GetRandomIncorrectQuestion(context.Background(), app.user, app.certification.ID)
	if err != nil {
		return nil, nil, err
	}
	if question == nil {
		return nil, nil, nil
	}
	previous_response, err := app.db.GetLastResponse(context.Background(), app.user, question.ID)
	if err != nil {
		return nil, nil, err
	}
	return question, previous_response, nil
}

func (app *App) getDomainAndTaskFromQuestion(questionID int) (domain *model.Domain, task *model.Task) {
	for i := range app.certification.Domains {
		for j := range app.certification.Domains[i].Tasks {
			for _, q := range app.certification.Domains[i].Tasks[j].Questions {
				if q.ID == questionID {
					domain = app.certification.Domains[i]
					task = app.certification.Domains[i].Tasks[j]
					break
				}
			}
		}
	}
	return
}

func (app *App) GetStats() (int, int, int, error) {
	var errs error
	correct, err := app.db.CountCorrectQuestions(context.Background(), app.user, app.certification.ID)
	errs = errors.Join(errs, err)
	incorrect, err := app.db.CountIncorrectQuestions(context.Background(), app.user, app.certification.ID)
	errs = errors.Join(errs, err)
	all, err := app.db.CountQuestions(context.Background(), app.certification.ID)
	errs = errors.Join(errs, err)
	return correct, incorrect, all, errs
}
