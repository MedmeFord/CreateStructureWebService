package main

import (
	"github.com/MedmeFord/CreateStructureWebService/pkg/models"
	"html/template"
	"path/filepath"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Инициализируем новую карту, которая будет хранить кэш.
	cache := map[string]*template.Template{}

	// Используем функцию filepath.Glob, чтобы получить срез всех файловых путей
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Перебираем файл шаблона от каждой страницы.
	for _, page := range pages {
		// Извлечение конечное названия файла (например, 'home.page.tmpl') из полного пути к файлу
		name := filepath.Base(page)
		// Обрабатываем итерируемый файл шаблона.
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Используем метод ParseGlob для добавления всех каркасных шаблонов.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		// Используем метод ParseGlob для добавления всех вспомогательных шаблонов.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		// Добавляем полученный набор шаблонов в кэш, используя название страницы
		cache[name] = ts
	}
	return cache, nil
}
