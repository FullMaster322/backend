package repository

import (
	"backend/back/pkg/models"
	"testing"
)

const connStr = "postgres://postgres:psql@localhost:5432/msfbd"

func TestLecturesCRUD(t *testing.T) {

	db, err := New(connStr)
	if err != nil {
		t.Fatal(err.Error())
	}

	data, err := db.GetLectures()
	if err != nil {
		t.Fatal(err.Error())
	}

	length := len(data)

	lecture := models.Lectures{
		NAME: "test name",
	}

	err = db.CreateLecture(lecture)
	if err != nil {
		t.Fatal(err.Error())
	}

	data, err = db.GetLectures()
	if err != nil || length == len(data) {
		t.Fatal(err.Error())
	}
}
