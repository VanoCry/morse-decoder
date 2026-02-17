// Package service предоставляет функции для работы с азбукой Морзе.
// Он позволяет преобразовывать текст в код Морзе и обратно,
// а также автоматически определять тип входных данных.
package service

import (
	"errors"
	"morse_decoder_vanocry/pkg/morse"
)

var ErrUnknownMode = errors.New("конвертация: неизвестный режим")

// TextDecoder выполняет преобразование между текстом и азбукой Морзе.
//
// Параметры:
//   - s: входная строка для преобразования
//   - mode: режим работы (0 - авто, 1 - текст->морзе, 2 - морзе->текст)
//
// Возвращает:
//   - string: результат преобразования
//   - error: ошибка в случае некорректного ввода или неизвестного режима
func TextDecoder(s string, mode int) (string, error) {
	switch mode {
	case 0:
		// Автоопределение
		if IsMorseCode(s) {
			return morse.ToText(s), nil
		}
		return morse.ToMorse(s), nil
	case 1:
		return morse.ToMorse(s), nil
	case 2:
		return morse.ToText(s), nil
	default:
		return "", ErrUnknownMode
	}
}

// IsMorseCode выполняет определение кода морзе.
//
// Параметры:
//   - s: входная строка для определения на наличие кода морзе
//
// Возвращает:
//   - bool: логическое значение наличия кода морзе (true - код морзе, false - обычный текст)
func IsMorseCode(s string) bool {
	for _, r := range s {
		if r != '.' && r != '-' && r != ' ' && r != '/' {
			return false // Найден символ не из морзе
		}
	}
	// Если строка не пустая и содержит хотя бы один символ морзе
	if len(s) > 0 {
		return true
	}
	return false
}
