package cli

import "github.com/rivo/tview"

var AllAnsweredView View = &allAnsweredView{}

type allAnsweredView struct {
	restart_button *button
	quit_button    *button
}

func (view *allAnsweredView) Build(cliapp *cliApp) (tview.Primitive, error) {
	view.restart_button = Button("Restart", func() { cliapp.Show(RestartView) }, 'r')
	view.quit_button = QuitButton(cliapp)
	return tview.NewGrid().SetRows(3, 0, 3).
		AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetDynamicColors(true).SetText("\n[green]CONGRATULATIONS!!"), 0, 0, 1, 1, 0, 0, false).
		AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("\nYou answered all questions for this certification."), 1, 0, 1, 1, 0, 0, false).
		AddItem(getFooterView(view.restart_button, view.quit_button), 2, 0, 1, 1, 0, 0, true).
		SetBorders(false), nil
}

func (view *allAnsweredView) OnMount(cliapp *cliApp) {
	cliapp.BindButtons(view.restart_button, view.quit_button)
}

func (view *allAnsweredView) OnDismount(cliapp *cliApp) {
	cliapp.UnbindButtons(view.restart_button, view.quit_button)
}
