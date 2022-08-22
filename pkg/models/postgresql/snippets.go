package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/MedmeFord/CreateStructureWebService/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, contents, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
		VALUES($1, $2, current_timestamp, current_timestamp + interval '1 year' * $3)`
	_, err := m.DB.Exec(stmt, title, contents, expires) // много вопросов
	if err != nil {
		return 0, err
	}

	id, err := lastEventIdSnippets(m.DB)
	if err != nil {
		return 0, err
	}
	fmt.Println(id)
	return int(id), err
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, created, created, expires FROM snippets
	WHERE expires > current_timestamp AND id = $1`

	row := m.DB.QueryRow(stmt, id)
	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > current_timestamp ORDER BY created DESC LIMIT 10`

	rows := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []*models.Snippet

	for rows.Next() {
		s := &models.Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func lastEventIdSnippets(db *sql.DB) (int, error) {
	var count int
	row := db.QueryRow("SELECT MAX(id) FROM snippets;")
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
