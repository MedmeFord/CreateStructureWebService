package postgresql

import (
	"database/sql"
	"github.com/MedmeFord/CreateStructureWebService/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	// Ниже будет SQL запрос, который мы хотим выполнить. Мы разделили его на две строки
	// для удобства чтения (поэтому он окружен обратными кавычками
	// вместо обычных двойных кавычек).
	// str := []rune(expires) + 'year'
	// str = append(str, 'year')

	// stmt := fmt.Sprintf("INSERT INTO snippets (title, content_, created, expires)VALUES(%s, %s,  current_timestamp, current_timestamp + interval '%s years", title, content, expires)

	// stmt := `INSERT INTO snippets (title, content_, created, expires)
	// VALUES($1, $2, current_timestamp, current_timestamp + $3::interval)`
	stmt := `INSERT INTO snippets (title, content, created, expires)
		VALUES($1, $2, current_timestamp, current_timestamp + interval '1 year' * $3)`
	_, err := m.DB.Exec(stmt, title, content, expires) // много вопросов
	if err != nil {
		return 0, err
	}
	// Используем метод LastInsertId(), чтобы получить последний ID
	// созданной записи из таблицу snippets.
	var id int
	err = m.DB.QueryRow("INSERT INTO snippets (title) VALUES ('John') RETURNING id").Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, err
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
