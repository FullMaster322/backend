package repository

import (
	"backend/back/pkg/models"
	"context"
	"fmt"
)

func (repo *PGRepo) LectureExistsByName(name string) (bool, error) {
	var exists bool

	err := repo.pool.QueryRow(context.Background(), `
		SELECT EXISTS(SELECT 1 FROM lectures WHERE name = $1)
	`, name).Scan(&exists)

	if err != nil {
		return false, nil
	}

	return exists, nil
}

func (repo *PGRepo) GetLectures() ([]models.Lectures, error) {
	rows, err := repo.pool.Query(context.Background(), `
		SELECT id, name FROM lectures
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []models.Lectures
	for rows.Next() {
		var item models.Lectures
		err = rows.Scan(
			&item.ID,
			&item.NAME,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (repo *PGRepo) GetLectureById(id int) (*models.Lectures, error) {
	var item models.Lectures

	err := repo.pool.QueryRow(context.Background(),
		`SELECT id, name, description FROM lectures WHERE id = $1`, id,
	).Scan(&item.ID, &item.NAME, &item.DESCRIPTION)

	return &item, err
}

func (repo *PGRepo) CreateLecture(item models.Lectures) error {
	exists, err := repo.LectureExistsByName(item.NAME)

	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("lecture with name '%s' already exists", item.NAME)
	}

	_, err = repo.pool.Exec(context.Background(), `
        INSERT INTO lectures (name) 
        VALUES ($1)
    `, item.NAME)
	return err
}

func (repo *PGRepo) UpdateLectureById(id int, item models.Lectures) error {
	_, err := repo.pool.Exec(context.Background(), `
		UPDATE lectures
		SET NAME = $1
		WHERE id = $2
	`, item.NAME, id)
	return err
}

func (repo *PGRepo) DeleteLectureById(id int) error {
	_, err := repo.pool.Exec(context.Background(), `
		DELETE FROM lectures
		WHERE id = $1
	`, id)

	return err
}
