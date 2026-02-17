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
	indexHTML, err := os.ReadFile("../static/index.html")
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
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка метода запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusInternalServerError)
		return
	}

	// Парсинг multipart формы
	err := r.ParseMultipartForm(10 << 20) // 10 MB максимальный размер
	if err != nil {
		http.Error(w, "Ошибка при парсинге формы", http.StatusInternalServerError)
		return
	}

	// Получение режима конвертации
	modestr := r.FormValue("mode")
	mode := 0 // Режим "авто" по умолчанию

	switch modestr {
	case "text-to-morse":
		mode = 1
	case "morse-to-text":
		mode = 2
	case "auto":
		mode = 0
	}

	// Получение файла из формы "myFile"
	file, header, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Не удалось получить файл из формы", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Проверка расширения файла
	if filepath.Ext(header.Filename) != ".txt" {
		http.Error(w, "Пожалуйста, загрузите файл с расширением .txt", http.StatusInternalServerError)
		return // Или всё таки StatusBadRequest?
	}

	// Чтение данных из файла
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Ошибка при чтении файла", http.StatusInternalServerError)
		return
	}

	// Преобразование данных в строку
	inputString := string(data)

	// Передача данных в декодер пакета service
	result, err := service.TextDecoder(inputString, mode)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при конвертации: %v", err), http.StatusInternalServerError)
		return
	}

	// Генерация имя для локального файла
	// Названия режимов для наглядности
	modeNames := map[int]string{
		0: "auto",
		1: "text_to_morse",
		2: "morse_to_text",
	}
	modeName := modeNames[mode]

	// Генерация имени файла с использованием времени
	timestamp := time.Now().UTC().Format("20060102_150405")
	filename := fmt.Sprintf("converted_%s_%s.txt", modeName, timestamp)
	// Создание файла
	filenameDir := filepath.Join(filePathForSavedFiles, filename)
	localFile, err := os.Create(filenameDir)
	if err != nil {
		http.Error(w, "Не удалось создать локальный файл", http.StatusInternalServerError)
		return
	}
	defer localFile.Close()

	// Запись результата конвертации в файл
	_, err = localFile.WriteString(result)
	if err != nil {
		http.Error(w, "Не удалось записать в файл", http.StatusInternalServerError)
		return
	}

	modeMessage := ""
	switch mode {
	case 0:
		modeMessage = "Автоматический режим"
	case 1:
		modeMessage = "Текст -> Морзе"
	case 2:
		modeMessage = "Морзе -> Текст"
	}

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

	w.Write([]byte(resultHTML))
}
