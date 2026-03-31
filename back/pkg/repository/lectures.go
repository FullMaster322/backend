package repository

import (
	"backend/back/pkg/models"
	"context"
	"fmt"
	"strings"
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

func prepareQuery(q string) string {
	words := strings.Fields(q)
	for i, w := range words {
		words[i] = w + ":*" // частичный поиск
	}
	return strings.Join(words, " & ")
}

func (repo *PGRepo) SearchLectures(query string) ([]map[string]interface{}, error) {
	rows, err := repo.pool.Query(context.Background(), `
		SELECT 
			id,
			name,
			ts_rank(
				to_tsvector('russian', name || ' ' || description),
				plainto_tsquery('russian', $1)
			) AS rank,
			ts_headline(
				'russian',
				description,
				plainto_tsquery('russian', $1),
				'MaxWords=20, MinWords=10'
			) AS snippet
		FROM lectures
		WHERE
			name ILIKE '%' || $1 || '%'
			OR description ILIKE '%' || $1 || '%'
			OR to_tsvector('russian', name || ' ' || description) @@ plainto_tsquery('russian', $1)
		ORDER BY rank DESC
		LIMIT 20
	`, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var name, snippet string
		var rank float32
		if err := rows.Scan(&id, &name, &rank, &snippet); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"id":      id,
			"name":    name,
			"snippet": snippet,
			"rank":    rank,
		})
	}
	return results, nil
}
