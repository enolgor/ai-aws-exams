package cli

import (
	"github.com/enolgor/ai-aws-exams/app"
	"github.com/enolgor/ai-aws-exams/conf"
	"github.com/enolgor/ai-aws-exams/db"
	"github.com/enolgor/ai-aws-exams/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CliApp interface {
	Run() error
}

//type View func(cliapp *cliApp) (tview.Primitive, error)

type View interface {
	Build(cliapp *cliApp) (tview.Primitive, error)
	OnMount(cliapp *cliApp)
	OnDismount(cliapp *cliApp)
}

type cliApp struct {
	db               *db.DB
	conf             *conf.Config
	user_email       string
	user_id          int
	cert_id          int
	cli              *tview.Application
	app              *app.App
	current_domain   *model.Domain
	current_task     *model.Task
	current_question *model.Question
	correct_count    int
	incorrect_count  int
	all_count        int
	selected_answer  int
	current_view     View
	keybinds         map[rune]func()
	builderr         error
}

func NewCliApp(db *db.DB, conf *conf.Config) (CliApp, error) {
	app := &cliApp{}
	app.db = db
	app.conf = conf
	app.user_email = conf.LastUsedEmail
	app.keybinds = make(map[rune]func())
	app.cli = tview.NewApplication()
	return app, nil
}

func (cliapp *cliApp) Run() error {
	cliapp.Show(LoginView)
	cliapp.cli.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if action, ok := cliapp.keybinds[event.Rune()]; ok {
			action()
		}
		return event
	})
	if err := cliapp.cli.Run(); err != nil {
		return err
	}
	return cliapp.builderr
}

func (cliapp *cliApp) Show(view View) {
	cliapp.builderr = nil
	primitive, err := view.Build(cliapp)
	if err != nil {
		cliapp.builderr = err
		cliapp.cli.Stop()
	}
	if primitive == nil {
		return
	}
	if cliapp.current_view != nil {
		cliapp.current_view.OnDismount(cliapp)
	}
	cliapp.current_view = view
	cliapp.cli.SetRoot(primitive, true).EnableMouse(true)
	view.OnMount(cliapp)
}

func (cliapp *cliApp) BindButtons(buttons ...*button) {
	for i := range buttons {
		if buttons[i].Shortcut == 0 {
			continue
		}
		cliapp.SetKeyBind(buttons[i].Shortcut, buttons[i].Func)
	}
}

func (cliapp *cliApp) UnbindButtons(buttons ...*button) {
	for i := range buttons {
		if buttons[i].Shortcut == 0 {
			continue
		}
		cliapp.ClearKeyBind(buttons[i].Shortcut)
	}
}

func (cliapp *cliApp) SetKeyBind(key rune, action func()) {
	cliapp.keybinds[key] = action
}

func (cliapp *cliApp) ClearKeyBind(key rune) {
	delete(cliapp.keybinds, key)
}

func (cliapp *cliApp) Stop() {
	cliapp.cli.Stop()
}
