package unpack

import (
	"testing"
)

func TestUnpack(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		// --- Базовые случаи ---
		{
			name:  "пустая строка",
			input: "",
			want:  "",
		},
		{
			name:  "строка без цифр",
			input: "abcd",
			want:  "abcd",
		},
		{
			name:  "повторяющиеся символы",
			input: "a4bc2d5e",
			want:  "aaaabccddddde",
		},
		{
			name:  "множитель 1",
			input: "a1b1",
			want:  "ab",
		},
		{
			name:  "множитель 0 — символ исчезает",
			input: "a0b",
			want:  "b",
		},
		{
			name:    "только цифры — ошибка",
			input:   "45",
			want:    "",
			wantErr: true,
		},
		{
			name:    "строка начинается с цифры — ошибка",
			input:   "3abc",
			want:    "",
			wantErr: true,
		},

		// --- Escape-последовательности ---
		{
			name:  "экранированные цифры",
			input: `qwe\4\5`,
			want:  "qwe45",
		},
		{
			name:  "экранированная цифра как символ + множитель",
			input: `qwe\45`,
			want:  "qwe44444",
		},
		{
			name:  "экранированный слеш",
			input: `qwe\\`,
			want:  `qwe\`,
		},
		{
			name:  "экранированный слеш с множителем",
			input: `qwe\\3`,
			want:  `qwe\\\`,
		},
		{
			name:    "одиночный слеш в конце — ошибка",
			input:   `qwe\`,
			want:    "",
			wantErr: true,
		},

		// --- Unicode / кириллица ---
		{
			name:  "кириллические символы",
			input: "а4б2в",
			want:  "ааааббв",
		},
	}

	for _, tc := range tests {
		tc := tc // захват переменной для параллельных тестов
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := Unpack(tc.input)

			if tc.wantErr {
				if err == nil {
					t.Errorf("Unpack(%q): ожидалась ошибка, получено %q", tc.input, got)
				}
				return
			}

			if err != nil {
				t.Errorf("Unpack(%q): неожиданная ошибка: %v", tc.input, err)
				return
			}

			if got != tc.want {
				t.Errorf("Unpack(%q):\n  получено: %q\n  ожидалось: %q", tc.input, got, tc.want)
			}
		})
	}
}
