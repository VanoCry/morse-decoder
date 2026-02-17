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
	indexHTML, err := os.ReadFile("index.html")
	if err != nil {
		//Возврат стандартного html
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		defaultHTML := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Morse Converter</title>
</head>
<body>
    <h1>Morse Converter</h1>
    <form enctype="multipart/form-data" action="/upload" method="post">
        <input type="file" name="myFile" required>
        <select name="mode">
            <option value="auto">Auto</option>
            <option value="text-to-morse">Text to Morse</option>
            <option value="morse-to-text">Morse to Text</option>
        </select>
        <input type="submit" value="Convert">
    </form>
</body>
</html>`
		// Отправка HTML
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(defaultHTML))
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
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем режим
	modeStr := r.FormValue("mode")
	mode := 0
	switch modeStr {
	case "text-to-morse":
		mode = 1
	case "morse-to-text":
		mode = 2
	}

	// получаем файл из формы
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

	// Сохраняем файл (опционально)
	os.MkdirAll("saved_files", 0755)
	filename := fmt.Sprintf("converted_%d_%s.txt", mode, time.Now().Format("20060102_150405"))
	os.WriteFile(filepath.Join("saved_files", filename), []byte(result), 0644)

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
