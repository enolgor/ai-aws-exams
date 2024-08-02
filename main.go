package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/enolgor/ai-aws-exams/cli"
	"github.com/enolgor/ai-aws-exams/conf"
	"github.com/enolgor/ai-aws-exams/db"
	"github.com/enolgor/ai-aws-exams/model"
	"github.com/enolgor/ai-aws-exams/prompts"
)

var appConfigDir string
var runcli bool
var parse_questions bool
var prompt bool
var initialize bool
var certification string
var file string
var questions bool
var definition bool
var config *conf.Config

func init() {
	var err error
	appConfigDir, err = os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	appConfigDir = path.Join(appConfigDir, ".ai-aws-exams")
	if err := os.MkdirAll(appConfigDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if config, err = conf.Load(appConfigDir); err != nil {
		log.Fatal(err)
	}
	flag.BoolVar(&initialize, "initialize", false, "initialize database")
	flag.BoolVar(&parse_questions, "parse-questions", false, "parse questions file")
	flag.BoolVar(&prompt, "prompt", false, "gpt prompt")
	flag.BoolVar(&questions, "questions", false, "questions prompt")
	flag.BoolVar(&definition, "definition", false, "definition prompt")
	flag.StringVar(&certification, "certification", "", "certification")
	flag.StringVar(&file, "file", "", "file")
	flag.Parse()
	if !runcli && !parse_questions && !prompt && !initialize {
		runcli = true
	}
	if parse_questions && (file == "" || certification == "") {
		log.Fatal("Must specify certification and file")
	}
	if prompt && (certification == "") {
		log.Fatal("Must specify certification")
	}
	if prompt && (!questions && !definition) {
		log.Fatal("Must specify question or definition")
	}
}

func main() {
	if prompt {
		exec_prompt()
		return
	}
	dbfile := path.Join(appConfigDir, "database.db")
	db, err := db.NewSQLiteDB(dbfile)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Println("closing db")
		db.Close()
	}()

	if initialize {
		exec_initialize(db)
		return
	}
	if parse_questions {
		exec_parse(db)
		return
	}
	if runcli {
		exec_cli(db)
		return
	}

}

func exec_initialize(db *db.DB) {
	if err := db.Initialize(context.Background()); err != nil {
		log.Fatal(err)
	}
	log.Println("Initialized db")
}

func exec_parse(db *db.DB) {
	cert, err := db.GetCertificationByName(context.Background(), certification)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	var questions []model.Question
	if err := dec.Decode(&questions); err != nil {
		log.Fatal(err)
	}
	for _, question := range questions {
		question.CertificationID = cert.ID
		id, err := db.InsertQuestion(context.Background(), question)
		if err != nil {
			log.Fatal(err)
		}
		for _, answer := range question.Answers {
			answer.QuestionID = id
			if err := db.InsertAnswer(context.Background(), *answer); err != nil {
				log.Fatal(err)
			}
		}
	}
	log.Println("Parsed questions")
}

func exec_prompt() {
	var str string
	var err error
	if questions {
		str, err = prompts.QuestionsPrompt(certification)
	}
	if definition {
		str, err = prompts.DefinitionPrompt(certification)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(str)
}

func exec_cli(db *db.DB) {
	app, err := cli.NewCliApp(db, config)
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
