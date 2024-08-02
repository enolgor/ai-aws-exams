package db

import (
	"embed"
	"encoding/json"
	"io/fs"

	"github.com/enolgor/ai-aws-exams/model"
)

//go:embed initial_data/certifications/*
var initialDataFS embed.FS

func load_initial_certifications() ([]model.Certification, error) {
	var certifications []model.Certification
	files, err := fs.ReadDir(initialDataFS, "initial_data/certifications")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			data, err := initialDataFS.ReadFile("initial_data/certifications/" + file.Name())
			if err != nil {
				return nil, err
			}

			var certification model.Certification
			if err := json.Unmarshal(data, &certification); err != nil {
				return nil, err
			}

			certifications = append(certifications, certification)
		}
	}

	return certifications, nil
}
