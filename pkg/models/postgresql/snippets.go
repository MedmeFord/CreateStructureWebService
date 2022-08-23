package postgresql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/MedmeFord/CreateStructureWebService/pkg/models"
)

// SnippetModel - Определяем тип который обертывает пул подключения sql.DB
type SnippetModel struct {
	DB *sql.DB
}

// Insert - Метод для создания новой заметки в базе дынных.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	// stmt := `INSERT INTO snippets (title, content_, created, expires)
	// VALUES($1, $2, current_timestamp, current_timestamp + $3::interval)`
	stmt := `INSERT INTO snippets (title, content_, created, expires)
		VALUES($1, $2, current_timestamp, current_timestamp + interval '1 year' * $3)`
	_, err := m.DB.Exec(stmt, title, content, expires) // много вопросов
	if err != nil {
		return 0, err
	}

	id, err := lastEventIdSnippets(m.DB)
	if err != nil {
		return 0, err
	}
	fmt.Println(id)
	return int(id), err // id - int64
}

// Get - Метод для возвращения данных заметки по её идентификатору ID.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT * FROM snippets
    WHERE expires > current_timestamp AND id = $1;`
	// Получение строчки
	row := m.DB.QueryRow(stmt, id)
	// Инициализируем указатель на новую структуру Snippet.
	s := &models.Snippet{}
	// Заполнение структуры для вывода
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	// Если все хорошо, возвращается объект Snippet.
	return s, nil
}

// Latest - Метод возвращает 10 наиболее часто используемые заметки.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {

	stmt := `SELECT * FROM snippets
    WHERE expires > current_timestamp ORDER BY created DESC LIMIT 10;`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []*models.Snippet // слайс структур

	for rows.Next() {
		// Создаем новый экземпляр структуры и заполняем ее из строки таблицы
		s := &models.Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Добавляем структуру в слайс структур
		snippets = append(snippets, s)
	}
	// вызываем метод Err() чтоб узнать не прозло ли ошибки пока работал цикл
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}

/*
support functioln
return - last insert id from snippets table
*/
func lastEventIdSnippets(db *sql.DB) (int, error) {
	var count int
	row := db.QueryRow("SELECT MAX(id) FROM snippets;")
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
