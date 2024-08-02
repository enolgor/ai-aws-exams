package cli

import (
	"context"

	"github.com/rivo/tview"
)

var RestartView View = &restartView{}

type restartView struct{}

func (view *restartView) Build(cliapp *cliApp) (tview.Primitive, error) {
	err := cliapp.db.ClearResponses(context.Background(), cliapp.user_id, cliapp.cert_id)
	if err != nil {
		return nil, err
	}
	cliapp.Show(QuestionView)
	return nil, nil
}

func (view *restartView) OnMount(cliapp *cliApp) {

}

func (view *restartView) OnDismount(cliapp *cliApp) {

}
