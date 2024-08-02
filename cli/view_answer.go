package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/enolgor/ai-aws-exams/model"
	"github.com/rivo/tview"
)

var AnswerView View = &answerView{}

type answerView struct {
	next_button *button
	quit_button *button
}

func (view *answerView) Build(cliapp *cliApp) (tview.Primitive, error) {
	var err error
	if err = cliapp.db.InsertResponse(context.Background(), model.Response{
		QuestionID: cliapp.current_question.ID,
		AnswerID:   cliapp.current_question.Answers[cliapp.selected_answer].ID,
		UserID:     cliapp.user_id,
		Date:       time.Now(),
	}); err != nil {
		return nil, err
	}
	cliapp.correct_count, cliapp.incorrect_count, cliapp.all_count, err = cliapp.app.GetStats()
	if err != nil {
		return nil, err
	}
	headerText := getHeaderText(cliapp)
	questionText := getQuestionText(cliapp, false)
	response_text := responseText(cliapp)
	view.next_button = Button("Next", func() { cliapp.Show(QuestionView) }, 'n')
	view.quit_button = QuitButton(cliapp)
	return defaultGrid(headerText, questionText, response_text, false, view.next_button, view.quit_button), nil
}

func (view *answerView) OnMount(cliapp *cliApp) {
	cliapp.BindButtons(view.next_button, view.quit_button)
}

func (view *answerView) OnDismount(cliapp *cliApp) {
	cliapp.UnbindButtons(view.next_button, view.quit_button)
}

func responseText(cliapp *cliApp) *tview.TextView {
	style := "red"
	if cliapp.current_question.Answers[cliapp.selected_answer].Correct {
		style = "green"
	}
	return tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignLeft).SetText(fmt.Sprintf(
		"[%s]%s\n\n[white]%s",
		style,
		cliapp.current_question.Answers[cliapp.selected_answer].Answer,
		cliapp.current_question.Answers[cliapp.selected_answer].Explanation,
	))
}
