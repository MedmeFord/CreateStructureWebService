package main

import (
	"database/sql"
	"flag"
	"github.com/MedmeFord/CreateStructureWebService/pkg/models/postgresql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *postgresql.SnippetModel
}

func main() {
	addr := flag.String("addr", "127.0.0.1:4000", "Сетевоой адресс HTTP") // флаг командной строки

	// очень жестко(pg_hba.conf + нужно правильно прописать адрес
	dsn := flag.String("dsn", "postgresql://web:123@127.0.0.1:5432/snipetbox?sslmode=disable", "Название postSQL источника данных")

	flag.Parse()                                                                  // извлечение флага из командной строки(меняет по адресу addr)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                  // создание логгера INFO в stdout
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) // логгер ошибок ERROR

	db, err := openDB(*dsn) // инициализируем пул подключений к базе
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog: infoLog,
		infoLog:  errorLog,
		snippets: &postgresql.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(), // создает маршрутизатор и тп для декомпозиции
	}
	infoLog.Printf("Запуск веб-сервера на http://%s", *addr)
	err = srv.ListenAndServe() // Запуск нового веб-сервера
	errorLog.Fatal(err)
}

type neuterdFileSystem struct {
	fs http.FileSystem
}

func (new_fs neuterdFileSystem) Open(path string) (http.File, error) {
	f, err := new_fs.fs.Open(path) // открываем вызываемый путь
	if err != nil {
		return nil, err
	}
	s, err := f.Stat() // os.File предоставляет доступ к информации о файле/пути os.FileInfo
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := new_fs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}

// Функция openDB() обертывает sql.Open() и возвращает пул соединений sql.DB
// для заданной строки подключения (DSN).
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil { // проверка того что все настроено правильно
		return nil, err
	}
	return db, nil
}

// type neuteredFileSystem struct {
// 	fs http.FileSystem
// }
//
// func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
// 	f, err := nfs.fs.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	s, err := f.Stat()
//
