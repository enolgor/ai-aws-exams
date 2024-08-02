package cli

import (
	"fmt"

	"github.com/rivo/tview"
)

func defaultGrid(header, question *tview.TextView, main tview.Primitive, focus bool, buttondefs ...*button) *tview.Grid {
	footer := getFooterView(buttondefs...)
	grid := tview.NewGrid().
		SetRows(5, 2, 0, 3).
		AddItem(header, 0, 0, 1, 1, 0, 0, false).
		AddItem(question, 1, 0, 1, 1, 0, 0, false).
		AddItem(main, 2, 0, 1, 1, 0, 0, focus).
		AddItem(footer, 3, 0, 1, 1, 0, 0, !focus).
		SetBorders(false)
	return grid
}

func getHeaderText(cliapp *cliApp) *tview.TextView {
	unanswered := cliapp.all_count - cliapp.correct_count - cliapp.incorrect_count
	return tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true).
		SetText(fmt.Sprintf(
			"[green]Correct: %d[white]   [red]Incorrect: %d[white]   [blue] Unanswered: %d[white]   Total: %d\n\nDomain: %d - %s\nTask:   %d - %s",
			cliapp.correct_count,
			cliapp.incorrect_count,
			unanswered,
			cliapp.all_count,
			cliapp.current_domain.Number,
			cliapp.current_domain.Name,
			cliapp.current_task.Number,
			cliapp.current_task.Name,
		))
}

func getQuestionText(cliapp *cliApp, previous_incorrect bool) *tview.TextView {
	showIncorrect := ""
	if previous_incorrect {
		showIncorrect = " [red]*Incorrect last time"
	}
	return tview.NewTextView().
		SetTextAlign(tview.AlignLeft).
		SetDynamicColors(true).
		SetText("[yellow]" + cliapp.current_question.Question + showIncorrect)
}

type button struct {
	Label    string
	Func     func()
	Shortcut rune
	_button  *tview.Button
}

func QuitButton(cliapp *cliApp) *button {
	return Button("Quit", cliapp.Stop, 'q')
}

func Button(label string, f func(), shortcut rune) *button {
	return &button{label, f, shortcut, nil}
}

func getFooterView(buttondefs ...*button) *tview.Form {
	buttons := []*button{}
	for i := range buttondefs {
		if buttondefs[i] != nil {
			buttons = append(buttons, buttondefs[i])
		}
	}
	form := tview.NewForm()
	for i := range buttons {
		shortcut := ""
		if buttons[i].Shortcut != 0 {
			shortcut = "(" + string(buttons[i].Shortcut) + ") "
		}
		form = form.AddButton(fmt.Sprintf("%s%s", shortcut, buttons[i].Label), buttons[i].Func)
	}
	for i := range buttons {
		buttons[i]._button = form.GetButton(i)
	}
	return form
}
