package postgresql

import (
	"database/sql"
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
	return nil, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
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
