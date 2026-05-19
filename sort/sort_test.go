package main

import (
	"bufio"
	"sort"
	"strings"
	"testing"
)

// TestReadLines проверяет чтение строк из сканера.
func TestReadLines(t *testing.T) {
	input := "banana\napple\ncherry\n"
	scanner := bufio.NewScanner(strings.NewReader(input))
	lines, err := readLines(scanner)
	if err != nil {
		t.Fatalf("readLines вернул ошибку: %v", err)
	}
	if len(lines) != 3 {
		t.Errorf("ожидалось 3 строки, получено %d", len(lines))
	}
}

// TestDefaultSort проверяет лексикографическую сортировку по умолчанию.
func TestDefaultSort(t *testing.T) {
	lines := []string{"banana", "apple", "cherry", "avocado"}
	opts := options{}
	sort.SliceStable(lines, buildLessFunc(lines, opts))
	want := []string{"apple", "avocado", "banana", "cherry"}
	for i, got := range lines {
		if got != want[i] {
			t.Errorf("позиция %d: ожидалось %q, получено %q", i, want[i], got)
		}
	}
}

// TestNumericSort проверяет числовую сортировку (-n).
func TestNumericSort(t *testing.T) {
	lines := []string{"10", "2", "100", "1", "20"}
	opts := options{numeric: true}
	sort.SliceStable(lines, buildLessFunc(lines, opts))
	want := []string{"1", "2", "10", "20", "100"}
	for i, got := range lines {
		if got != want[i] {
			t.Errorf("позиция %d: ожидалось %q, получено %q", i, want[i], got)
		}
	}
}

// TestReverseSort проверяет обратный порядок сортировки (-r).
func TestReverseSort(t *testing.T) {
	lines := []string{"apple", "cherry", "banana"}
	opts := options{reverse: true}
	sort.SliceStable(lines, buildLessFunc(lines, opts))
	want := []string{"cherry", "banana", "apple"}
	for i, got := range lines {
		if got != want[i] {
			t.Errorf("позиция %d: ожидалось %q, получено %q", i, want[i], got)
		}
	}
}

// TestNumericReverseSort проверяет числовую сортировку в обратном порядке (-nr).
func TestNumericReverseSort(t *testing.T) {
	lines := []string{"10", "2", "100", "1"}
	opts := options{numeric: true, reverse: true}
	sort.SliceStable(lines, buildLessFunc(lines, opts))
	want := []string{"100", "10", "2", "1"}
	for i, got := range lines {
		if got != want[i] {
			t.Errorf("позиция %d: ожидалось %q, получено %q", i, want[i], got)
		}
	}
}

// TestUnique проверяет удаление дублирующихся строк (-u).
func TestUnique(t *testing.T) {
	lines := []string{"apple", "apple", "banana", "banana", "cherry"}
	opts := options{}
	result := deduplicate(lines, opts)
	want := []string{"apple", "banana", "cherry"}
	if len(result) != len(want) {
		t.Fatalf("ожидалось %d строк, получено %d", len(want), len(result))
	}
	for i, got := range result {
		if got != want[i] {
			t.Errorf("позиция %d: ожидалось %q, получено %q", i, want[i], got)
		}
	}
}

// TestColumnSort проверяет сортировку по колонке (-k).
func TestColumnSort(t *testing.T) {
	lines := []string{
		"Charlie\t30\tEngineer",
		"Alice\t25\tDesigner",
		"Bob\t35\tManager",
	}
	opts := options{column: 2, numeric: true}
	sort.SliceStable(lines, buildLessFunc(lines, opts))
	if !strings.HasPrefix(lines[0], "Alice") {
		t.Errorf("первой строкой должна быть Alice (25), получено: %q", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Bob") {
		t.Errorf("последней строкой должна быть Bob (35), получено: %q", lines[2])
	}
}

// TestMonthSort проверяет сортировку по названию месяца (-M).
func TestMonthSort(t *testing.T) {
	lines := []string{"Dec", "Jan", "Jul", "Feb", "Mar"}
	opts := options{month: true}
	sort.SliceStable(lines, buildLessFunc(lines, opts))
	want := []string{"Jan", "Feb", "Mar", "Jul", "Dec"}
	for i, got := range lines {
		if got != want[i] {
			t.Errorf("позиция %d: ожидалось %q, получено %q", i, want[i], got)
		}
	}
}

// TestTrimBlanks проверяет игнорирование хвостовых пробелов (-b).
func TestTrimBlanks(t *testing.T) {
	lines := []string{"banana  ", "apple", "cherry   "}
	opts := options{trimBlanks: true}
	sort.SliceStable(lines, buildLessFunc(lines, opts))
	// После trim: "banana", "apple", "cherry" → "apple", "banana", "cherry"
	key0 := extractKey(lines[0], opts)
	if key0 != "apple" {
		t.Errorf("ожидался ключ %q, получено %q", "apple", key0)
	}
}

// TestHumanSort проверяет сортировку с суффиксами (-h).
func TestHumanSort(t *testing.T) {
	lines := []string{"2G", "500M", "1K", "10M"}
	opts := options{human: true}
	sort.SliceStable(lines, buildLessFunc(lines, opts))
	want := []string{"1K", "10M", "500M", "2G"}
	for i, got := range lines {
		if got != want[i] {
			t.Errorf("позиция %d: ожидалось %q, получено %q", i, want[i], got)
		}
	}
}

// TestCheckOrderSorted проверяет, что checkOrder возвращает 0 для отсортированных данных.
func TestCheckOrderSorted(t *testing.T) {
	lines := []string{"apple", "banana", "cherry"}
	opts := options{}
	if n := checkOrder(lines, opts); n != 0 {
		t.Errorf("ожидалось 0 (данные отсортированы), получено %d", n)
	}
}

// TestCheckOrderUnsorted проверяет, что checkOrder возвращает номер нарушающей строки.
func TestCheckOrderUnsorted(t *testing.T) {
	lines := []string{"apple", "cherry", "banana"}
	opts := options{}
	if n := checkOrder(lines, opts); n != 3 {
		t.Errorf("ожидалось 3 (строка 'banana' нарушает порядок), получено %d", n)
	}
}

// TestParseHuman проверяет разбор человекочитаемых числовых значений.
func TestParseHuman(t *testing.T) {
	cases := []struct {
		input string
		want  float64
	}{
		{"1K", 1e3},
		{"2M", 2e6},
		{"3G", 3e9},
		{"100", 100},
		{"1.5K", 1500},
	}
	for _, c := range cases {
		got := parseHuman(c.input)
		if got != c.want {
			t.Errorf("parseHuman(%q) = %v, ожидалось %v", c.input, got, c.want)
		}
	}
}

// TestParseMonth проверяет разбор названий месяцев.
func TestParseMonth(t *testing.T) {
	cases := []struct {
		input string
		want  int
	}{
		{"Jan", 1}, {"jan", 1}, {"JAN", 1},
		{"Dec", 12}, {"Jun", 6},
		{"unknown", 0},
	}
	for _, c := range cases {
		got := parseMonth(c.input)
		if got != c.want {
			t.Errorf("parseMonth(%q) = %d, ожидалось %d", c.input, got, c.want)
		}
	}
}

// TestExtractKeyNoColumn проверяет извлечение ключа без указания колонки.
func TestExtractKeyNoColumn(t *testing.T) {
	opts := options{}
	got := extractKey("hello world", opts)
	if got != "hello world" {
		t.Errorf("ожидалось %q, получено %q", "hello world", got)
	}
}

// TestExtractKeyWithColumn проверяет извлечение ключа по указанной колонке.
func TestExtractKeyWithColumn(t *testing.T) {
	opts := options{column: 2}
	got := extractKey("Alice\t42\tEngineer", opts)
	if got != "42" {
		t.Errorf("ожидалось %q, получено %q", "42", got)
	}
}

// TestExtractKeyColumnOutOfRange проверяет поведение при выходе колонки за диапазон.
func TestExtractKeyColumnOutOfRange(t *testing.T) {
	opts := options{column: 10}
	got := extractKey("only\ttwo\tfields", opts)
	if got != "" {
		t.Errorf("ожидалась пустая строка, получено %q", got)
	}
}

// TestDeduplicateEmpty проверяет deduplicate на пустом срезе.
func TestDeduplicateEmpty(t *testing.T) {
	result := deduplicate(nil, options{})
	if result != nil {
		t.Errorf("ожидался nil, получено %v", result)
	}
}

// TestSortStable проверяет стабильность сортировки (порядок равных элементов сохраняется).
func TestSortStable(t *testing.T) {
	lines := []string{"b\t2", "a\t1", "b\t1", "a\t2"}
	opts := options{column: 1}
	sort.SliceStable(lines, buildLessFunc(lines, opts))
	// "a" должен идти перед "b", а среди "a" и "b" сохраняется исходный порядок
	if !strings.HasPrefix(lines[0], "a") || !strings.HasPrefix(lines[1], "a") {
		t.Errorf("ожидалось, что первые две строки начинаются с 'a': %v", lines[:2])
	}
}
