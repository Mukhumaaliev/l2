package unpack

import (
	"errors"
	"strings"
	"unicode"
)

// ErrInvalidString возвращается, если строка некорректна (например, состоит только из цифр).
var ErrInvalidString = errors.New("invalid string: no characters to repeat")

// Unpack распаковывает строку вида "a4bc2d5e" -> "aaaabccddddde".
//
// Правила:
//   - Цифра после обычного символа означает количество его повторений.
//   - Обратный слеш экранирует следующий символ (цифру или сам слеш),
//     делая его обычным символом.
//   - Строка, состоящая только из цифр (без экранирования), считается
//     некорректной и возвращает ErrInvalidString.
//   - Пустая строка возвращает пустую строку без ошибки.
func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	runes := []rune(s)
	var sb strings.Builder
	var prev rune
	hasPrev := false

	for i := 0; i < len(runes); i++ {
		ch := runes[i]

		switch {
		case ch == '\\':
			// Экранирование: следующий символ трактуется буквально.
			if i+1 >= len(runes) {
				// Одиночный слеш в конце строки — некорректный ввод.
				return "", ErrInvalidString
			}
			// Записываем предыдущий накопленный символ (без умножения).
			if hasPrev {
				sb.WriteRune(prev)
			}
			i++
			prev = runes[i]
			hasPrev = true

		case unicode.IsDigit(ch):
			// Цифра — множитель для предыдущего символа.
			if !hasPrev {
				// Строка начинается с цифры и предшествующего символа нет.
				return "", ErrInvalidString
			}
			count := int(ch - '0')
			for j := 0; j < count; j++ {
				sb.WriteRune(prev)
			}
			hasPrev = false
			prev = 0

		default:
			// Обычный символ: записываем предыдущий (если есть) и сохраняем текущий.
			if hasPrev {
				sb.WriteRune(prev)
			}
			prev = ch
			hasPrev = true
		}
	}

	// Записываем последний символ, если он не был обработан множителем.
	if hasPrev {
		sb.WriteRune(prev)
	}

	return sb.String(), nil
}
