package cli

import (
	"github.com/enolgor/ai-aws-exams/model"
	"github.com/rivo/tview"
)

var QuestionView View = &questionView{}

type questionView struct {
	answersList               *tview.List
	answer_button             *button
	skip_button               *button
	hide_show_answer_button   *button
	quit_button               *button
	hide_show_response_button *button
	enable_hint               bool
	enable_last_response      bool
}

func (view *questionView) Build(cliapp *cliApp) (tview.Primitive, error) {
	view.enable_hint = false
	view.enable_last_response = false
	var err error
	var question *model.Question
	var previous_response *model.Response
	cliapp.current_domain, cliapp.current_task, question, previous_response, err = cliapp.app.GetNextRandomQuestion()
	if err != nil {
		return nil, err
	}
	if question == nil {
		cliapp.Show(AllAnsweredView)
		return nil, nil
	}
	cliapp.correct_count, cliapp.incorrect_count, cliapp.all_count, err = cliapp.app.GetStats()
	if err != nil {
		return nil, err
	}
	remaining := cliapp.all_count - cliapp.correct_count - cliapp.incorrect_count
	if remaining > 1 && cliapp.current_question != nil && question.ID == cliapp.current_question.ID {
		cliapp.Show(QuestionView)
		return nil, nil
	}
	cliapp.current_question = question
	cliapp.selected_answer = 0
	view.answersList = getAnswersList(cliapp.current_question, previous_response).SetSelectedFunc(func(idx int, _, _ string, _ rune) {
		cliapp.selected_answer = idx
		cliapp.Show(AnswerView)
	}).SetChangedFunc(func(idx int, _, _ string, _ rune) {
		cliapp.selected_answer = idx
	})
	view.answer_button = Button("Answer", func() { cliapp.Show(AnswerView) }, 'a')
	view.skip_button = Button("Skip", func() { cliapp.Show(QuestionView) }, 's')
	view.hide_show_answer_button = Button("Show Correct Answer", func() {
		view.HideShowAnswers(cliapp, previous_response)
	}, 'h')
	view.hide_show_response_button = nil
	if previous_response != nil {
		view.hide_show_response_button = Button("Show Previous Response", func() {
			view.HideShowPreviousResponse(cliapp, previous_response)
		}, 'i')
	}
	view.quit_button = QuitButton(cliapp)
	headerText := getHeaderText(cliapp)
	questionText := getQuestionText(cliapp, previous_response != nil)
	grid := defaultGrid(headerText, questionText, view.answersList, true, view.answer_button, view.skip_button, view.hide_show_answer_button, view.hide_show_response_button, view.quit_button)
	return grid, nil
}

func (view *questionView) OnMount(cliapp *cliApp) {
	cliapp.BindButtons(view.answer_button, view.skip_button, view.hide_show_answer_button, view.hide_show_response_button, view.quit_button)
}

func (view *questionView) OnDismount(cliapp *cliApp) {
	cliapp.UnbindButtons(view.answer_button, view.skip_button, view.hide_show_answer_button, view.hide_show_response_button, view.quit_button)
}

func (view *questionView) HideShowAnswers(cliapp *cliApp, previous_response *model.Response) {
	view.enable_hint = !view.enable_hint
	view.answersList.Clear()
	addAnswers(view.answersList, cliapp.current_question, view.enable_hint, view.enable_last_response, previous_response)
	if view.enable_hint {
		view.hide_show_answer_button._button.SetLabel("(h) Hide Correct Answer")
	} else {
		view.hide_show_answer_button._button.SetLabel("(h) Show Correct Answer")
	}
}

func (view *questionView) HideShowPreviousResponse(cliapp *cliApp, previous_response *model.Response) {
	view.enable_last_response = !view.enable_last_response
	view.answersList.Clear()
	addAnswers(view.answersList, cliapp.current_question, view.enable_hint, view.enable_last_response, previous_response)
	if view.enable_last_response {
		view.hide_show_response_button._button.SetLabel("(i) Hide Previous Response")
	} else {
		view.hide_show_response_button._button.SetLabel("(i) Show Previous Response")
	}
}

func getAnswersList(question *model.Question, lastResponse *model.Response) *tview.List {
	list := tview.NewList()
	addAnswers(list, question, false, false, lastResponse)
	return list
}

func addAnswers(list *tview.List, question *model.Question, enable_hint bool, enable_last_response bool, lastResponse *model.Response) {
	for i := range question.Answers {
		hint := ""
		if enable_hint && question.Answers[i].Correct {
			hint = "*"
		}
		last_response := ""
		if enable_last_response && lastResponse != nil && lastResponse.AnswerID == question.Answers[i].ID {
			last_response = " <-"
		}
		list.AddItem(question.Answers[i].Answer+last_response, hint, rune('0'+i+1), nil)
	}
}
