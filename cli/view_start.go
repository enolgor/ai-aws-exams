package cli

import (
	"context"

	"github.com/enolgor/ai-aws-exams/app"
	"github.com/enolgor/ai-aws-exams/model"
	"github.com/rivo/tview"
)

var StartView View = &startView{}

type startView struct{}

func (view *startView) Build(cliapp *cliApp) (tview.Primitive, error) {
	var err error
	cliapp.conf.LastUsedEmail = cliapp.user_email
	if err = cliapp.conf.Save(); err != nil {
		return nil, err
	}
	user, _ := cliapp.db.GetUserByEmail(context.Background(), cliapp.user_email)
	if user == nil {
		cliapp.user_id, err = cliapp.db.InsertUser(context.Background(), model.User{Email: cliapp.user_email})
		if err != nil {
			return nil, err
		}
	} else {
		cliapp.user_id = user.ID
	}
	cliapp.app, err = app.NewApp(cliapp.db, cliapp.user_id, cliapp.cert_id)
	if err != nil {
		return nil, err
	}
	cliapp.Show(QuestionView)
	return nil, nil
}

func (view *startView) OnMount(cliapp *cliApp) {

}

func (view *startView) OnDismount(cliapp *cliApp) {

}
