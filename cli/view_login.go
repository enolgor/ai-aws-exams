package cli

import (
	"context"

	"github.com/rivo/tview"
)

var LoginView View = &loginView{}

type loginView struct {
	start_button *tview.Button
}

func (view *loginView) Build(cliapp *cliApp) (tview.Primitive, error) {
	var form *tview.Form
	certifications, err := cliapp.db.ListCertifications(context.Background())
	if err != nil {
		return nil, err
	}
	cliapp.cert_id = certifications[0].ID
	cert_names := make([]string, len(certifications))
	for i := range certifications {
		cert_names[i] = certifications[i].Name
	}
	shouldEnableStart := func() {
		if form == nil {
			return
		}
		if view.start_button == nil {
			return
		}
		if cliapp.user_email != "" && cliapp.cert_id != 0 {
			view.start_button.SetDisabled(false)
			return
		}
		view.start_button.SetDisabled(true)
	}
	on_cert_selected := func(_ string, idx int) {
		cliapp.cert_id = certifications[idx].ID
		shouldEnableStart()
	}
	cert_dropdown := tview.NewDropDown().
		SetLabel("Certification").
		SetOptions(cert_names, on_cert_selected).
		SetCurrentOption(0)
	on_username_changed := func(value string) {
		cliapp.user_email = value
		shouldEnableStart()
	}
	username_input := tview.NewInputField().
		SetLabel("Username").
		SetText(cliapp.user_email).
		SetFieldWidth(40).
		SetAcceptanceFunc(nil).
		SetChangedFunc(on_username_changed)
	on_start := func() { cliapp.Show(StartView) }
	on_quit := func() { cliapp.Stop() }
	form = tview.NewForm().
		AddFormItem(cert_dropdown).
		AddFormItem(username_input).
		AddButton("Start", on_start).
		AddButton("Quit", on_quit).
		SetCancelFunc(on_quit)
	view.start_button = form.GetButton(0)
	shouldEnableStart()
	return form, nil
}

func (view *loginView) OnMount(cliapp *cliApp) {
	if cliapp.user_email != "" {
		cliapp.cli.SetFocus(view.start_button)
	}
}

func (view *loginView) OnDismount(cliapp *cliApp) {

}
