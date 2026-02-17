package handlers

import (
	"fmt"
	"io"
	"morse_decoder_vanocry/internal/service"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var filePathForSavedFiles string = "../saved_files/"

func RootHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка пути
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// Чтение файла index.html
	indexHTML, err := os.ReadFile("../cmd/index.html")
	if err != nil {
		http.Error(w, "Не удалось прочитать файл index.html", http.StatusInternalServerError)
		return
	}
	// Установка заголовка Content-Type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Отправка HTML
	w.WriteHeader(http.StatusOK)
	w.Write(indexHTML)
}

// UploadHandler обрабатывает загрузку файлов
// Переписано под тесты
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсим форму
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем режим
	mode := 0 // auto
	switch r.FormValue("mode") {
	case "text-to-morse":
		mode = 1
	case "morse-to-text":
		mode = 2
	}

	// Получаем файл из формы
	file, _, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Читаем содержимое файла
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Конвертируем
	result, err := service.TextDecoder(string(data), mode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем папку если нужно
	os.MkdirAll(filePathForSavedFiles, 0755)

	// Сохраняем результат в файл
	filename := fmt.Sprintf("converted_%d_%s.txt", mode, time.Now().Format("20060102_150405"))
	fullPath := filepath.Join(filePathForSavedFiles, filename)
	if err := os.WriteFile(fullPath, []byte(result), 0644); err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

/* Возврат html отключен для прохождения тестов
	// Возвращаем HTML страницу с результатом
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Встраиваем результат в HTML
	resultHTML := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta http-equiv="refresh" content="0;url=../static/upload.html?mode=%d&filename=%s&savedAs=%s&filepath=%s">
</head>
<body>
    <script>
        sessionStorage.setItem('conversionResult', JSON.stringify({
            mode: '%s',
            filename: '%s',
            savedAs: '%s',
            filepath: '%s',
            result: %q
        }));
        window.location.href = '../static/upload.html';
    </script>
</body>
</html>`,
		mode, header.Filename, filename, filenameDir,
		modeMessage, header.Filename, filename, filenameDir, result)

	w.Write([]byte(resultHTML))*/
