package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// monthOrder сопоставляет сокращение месяца с его порядковым номером.
var monthOrder = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4,
	"may": 5, "jun": 6, "jul": 7, "aug": 8,
	"sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

// humanSuffixes — суффиксы человекочитаемых размеров и их множители.
var humanSuffixes = map[byte]float64{
	'k': 1e3, 'K': 1e3,
	'm': 1e6, 'M': 1e6,
	'g': 1e9, 'G': 1e9,
	't': 1e12, 'T': 1e12,
	'p': 1e15, 'P': 1e15,
}

// options хранит разобранные флаги командной строки.
type options struct {
	column      int
	numeric     bool
	reverse     bool
	unique      bool
	month       bool
	trimBlanks  bool
	checkSorted bool
	human       bool
}

// parseFlags разбирает аргументы командной строки и возвращает параметры и имена файлов.
func parseFlags() (opts options, files []string) {
	flag.IntVar(&opts.column, "k", 0, "сортировать по столбцу N (разделитель — табуляция)")
	flag.BoolVar(&opts.numeric, "n", false, "числовая сортировка")
	flag.BoolVar(&opts.reverse, "r", false, "обратный порядок")
	flag.BoolVar(&opts.unique, "u", false, "только уникальные строки")
	flag.BoolVar(&opts.month, "M", false, "сортировка по названию месяца")
	flag.BoolVar(&opts.trimBlanks, "b", false, "игнорировать хвостовые пробелы")
	flag.BoolVar(&opts.checkSorted, "c", false, "проверить, отсортированы ли данные")
	flag.BoolVar(&opts.human, "h", false, "сортировка по числу с суффиксами (K, M, G…)")
	flag.Parse()
	files = flag.Args()
	return
}

// readLines читает все строки из сканера в срез.
func readLines(scanner *bufio.Scanner) ([]string, error) {
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// openScanner возвращает сканер для первого файла из списка или для STDIN.
func openScanner(files []string) (*bufio.Scanner, func(), error) {
	if len(files) == 0 {
		return bufio.NewScanner(os.Stdin), func() {}, nil
	}

	f, err := os.Open(files[0])
	if err != nil {
		return nil, nil, fmt.Errorf("не удалось открыть файл %q: %w", files[0], err)
	}
	return bufio.NewScanner(f), func() { f.Close() }, nil
}

// extractKey извлекает ключ из строки согласно параметрам.
func extractKey(line string, opts options) string {
	key := line

	if opts.column > 0 {
		parts := strings.Split(line, "\t")
		if opts.column <= len(parts) {
			key = parts[opts.column-1]
		} else {
			key = ""
		}
	}

	if opts.trimBlanks {
		key = strings.TrimRight(key, " \t")
	}

	return key
}

// parseHuman разбирает строку с человекочитаемым числовым суффиксом (1K, 2M и т.д.).
// Возвращает числовое значение в байтах/единицах.
func parseHuman(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}

	last := s[len(s)-1]
	if mult, ok := humanSuffixes[last]; ok {
		num, err := strconv.ParseFloat(s[:len(s)-1], 64)
		if err != nil {
			return math.NaN()
		}
		return num * mult
	}

	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return math.NaN()
	}
	return num
}

// parseMonth возвращает числовой номер месяца по его сокращённому названию.
// Если строка не является месяцем, возвращает 0.
func parseMonth(s string) int {
	return monthOrder[strings.ToLower(strings.TrimSpace(s))]
}

// buildLessFunc создаёт функцию сравнения для sort.SliceStable на основе параметров.
func buildLessFunc(lines []string, opts options) func(i, j int) bool {
	getKey := func(line string) string {
		return extractKey(line, opts)
	}

	return func(i, j int) bool {
		a, b := getKey(lines[i]), getKey(lines[j])

		var less bool
		switch {
		case opts.numeric:
			na, ea := strconv.ParseFloat(strings.TrimSpace(a), 64)
			nb, eb := strconv.ParseFloat(strings.TrimSpace(b), 64)
			if ea == nil && eb == nil {
				less = na < nb
			} else {
				less = a < b
			}

		case opts.human:
			na, nb := parseHuman(a), parseHuman(b)
			if !math.IsNaN(na) && !math.IsNaN(nb) {
				less = na < nb
			} else {
				less = a < b
			}

		case opts.month:
			ma, mb := parseMonth(a), parseMonth(b)
			if ma != 0 && mb != 0 {
				less = ma < mb
			} else if ma != 0 {
				less = true
			} else if mb != 0 {
				less = false
			} else {
				less = a < b
			}

		default:
			less = a < b
		}

		if opts.reverse {
			return !less
		}
		return less
	}
}

// deduplicate удаляет строки с одинаковым ключом (для флага -u).
func deduplicate(lines []string, opts options) []string {
	if len(lines) == 0 {
		return lines
	}

	result := lines[:1]
	for i := 1; i < len(lines); i++ {
		if extractKey(lines[i], opts) != extractKey(lines[i-1], opts) {
			result = append(result, lines[i])
		}
	}
	return result
}

// checkOrder проверяет, отсортированы ли строки согласно заданным параметрам.
// Возвращает номер строки (1-based), нарушающей порядок, или 0 если всё ок.
func checkOrder(lines []string, opts options) int {
	less := buildLessFunc(lines, opts)
	for i := 1; i < len(lines); i++ {
		// lines[i] < lines[i-1] означает нарушение порядка
		if less(i, i-1) {
			return i + 1
		}
	}
	return 0
}

// writeLines записывает срез строк в стандартный вывод через буферизованный писатель.
func writeLines(lines []string) error {
	w := bufio.NewWriter(os.Stdout)
	for _, l := range lines {
		if _, err := fmt.Fprintln(w, l); err != nil {
			return err
		}
	}
	return w.Flush()
}

func main() {
	opts, files := parseFlags()

	scanner, cleanup, err := openScanner(files)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sort:", err)
		os.Exit(1)
	}
	defer cleanup()

	// Увеличиваем буфер сканера для обработки длинных строк.
	const maxLineSize = 10 * 1024 * 1024 // 10 МБ на строку
	scanner.Buffer(make([]byte, 64*1024), maxLineSize)

	lines, err := readLines(scanner)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sort: ошибка чтения:", err)
		os.Exit(1)
	}

	// Режим проверки (-c): только сообщаем, отсортированы ли данные.
	if opts.checkSorted {
		if n := checkOrder(lines, opts); n != 0 {
			fmt.Fprintf(os.Stderr, "sort: данные не отсортированы (строка %d нарушает порядок)\n", n)
			os.Exit(1)
		}
		return
	}

	sort.SliceStable(lines, buildLessFunc(lines, opts))

	if opts.unique {
		lines = deduplicate(lines, opts)
	}

	if err := writeLines(lines); err != nil {
		fmt.Fprintln(os.Stderr, "sort: ошибка записи:", err)
		os.Exit(1)
	}
}
