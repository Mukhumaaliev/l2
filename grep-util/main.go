package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type options struct {
	AfterNstring  int
	BeforeNstring int
	RoundNstring  int
	C             bool
	I             bool
	V             bool
	F             bool
	N             bool
}

func main() {
	var text []string
	opts, files, pattern := ParseFlags()
	lines, err := ReadLines(files)
	if errors.Is(err, fmt.Errorf("не удалось открыть файл %q: %w", files[0], err)) {
		panic(err)
	}
	if opts.C {
		count := CountLine(pattern, lines, opts)
		fmt.Print(count)
	}
	switch {
	case opts.V:
		text = TextWithoutPattern(pattern, lines, opts)
	case opts.AfterNstring > 0:
		text = AfterStrings(pattern, lines, opts)
	case opts.BeforeNstring > 0:
		text = BeforeStrings(pattern, lines, opts)
	case opts.RoundNstring > 0:
		text = RoundStrings(pattern, lines, opts)
	}
	if opts.N {
		fmt.Print(LineWithNumber(text))
	} else {
		for _, line := range text {
			fmt.Println(line)
		}
	}
}

func AfterStrings(pattern string, lines []string, opts options) (result []string) {
	for i, line := range lines {
		if check(pattern, line, opts) {
			end := i + opts.AfterNstring
			if end >= len(lines) {
				end = len(lines) - 1
			}
			for ind := i; ind <= end; ind++ {
				result = append(result, lines[ind])
			}
		}
	}
	return
}

func check(pattern string, line string, opts options) bool {
	if opts.F {
		return FixPatternString(pattern, line)
	}
	return SearchByPatternInLine(pattern, line, opts)
}

func BeforeStrings(pattern string, lines []string, opts options) (result []string) {
	for i, line := range lines {
		if check(pattern, line, opts) {
			start := i - opts.BeforeNstring
			if start < 0 {
				start = 0
			}
			for ind := start; ind <= i; ind++ {
				result = append(result, lines[ind])
			}
		}
	}
	return
}

func RoundStrings(pattern string, lines []string, opts options) (result []string) {
	for i, line := range lines {
		if check(pattern, line, opts) {
			start := i - opts.RoundNstring
			if start < 0 {
				start = 0
			}
			end := i + opts.RoundNstring
			if end >= len(lines) {
				end = len(lines) - 1
			}
			for ind := start; ind <= end; ind++ {
				result = append(result, lines[ind])
			}
		}
	}
	result = Unique(result)
	return
}

func Unique(lines []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
			result = append(result, line)
		}
	}
	return result
}

func LineWithNumber(lines []string) (result []string) {
	for num, line := range lines {
		result = append(result, fmt.Sprintf("%d %s\n", num+1, line))
	}
	return result
}

func CountFixPattern(pattern string, lines []string) (count int) {
	count = 0
	for _, line := range lines {
		if FixPatternString(pattern, line) {
			count++
		}
	}
	return
}
func CountLine(pattern string, lines []string, opts options) (count int) {
	count = 0
	for _, line := range lines {
		if !opts.V {
			found := SearchByPatternInLine(pattern, line, opts)
			if found {
				count++
			}
		} else {
			CountWithoutPattern(pattern, line, opts, &count)
		}
	}
	return
}
func TextWithoutPattern(pattern string, lines []string, opts options) (result []string) {
	for _, line := range lines {
		if opts.I {
			re := regexp.MustCompile(strings.ToLower(pattern))
			if !re.MatchString(strings.ToLower(line)) {
				result = append(result, line)
			}
		} else {
			re := regexp.MustCompile(pattern)
			if !re.MatchString(line) {
				result = append(result, line)
			}
		}
	}
	return result
}

func CountWithoutPattern(pattern string, line string, opts options, count *int) {
	if opts.I {
		re := regexp.MustCompile(strings.ToLower(pattern))
		if !re.MatchString(strings.ToLower(line)) {
			*count++
		}
	} else {
		re := regexp.MustCompile(pattern)
		if !re.MatchString(line) {
			*count++
		}
	}

}

func FixPatternString(pattern string, line string) (result bool) {
	result = false
	if line == pattern {
		result = true
	}
	return
}

func SearchByPatternInLine(pattern string, line string, opts options) bool {
	var result bool
	if opts.I {
		re := regexp.MustCompile(strings.ToLower(pattern))
		result = re.MatchString(strings.ToLower(line))
	} else {
		re := regexp.MustCompile(pattern)
		result = re.MatchString(line)
	}
	return result
}

func ReadLines(files []string) (lines []string, err error) {
	text, err := OpenFile(files)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл %q: %w", files[0], err)
	}
	for text.Scan() {
		lines = append(lines, text.Text())
	}
	return
}

func OpenFile(files []string) (text *bufio.Scanner, err error) {
	if len(files) == 0 {
		text := bufio.NewScanner(os.Stdin)
		return text, nil
	}
	f, err := os.Open(files[0])
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл %q: %w", files[0], err)
	}
	text = bufio.NewScanner(f)
	return text, nil
}

func ParseFlags() (opts options, files []string, pattern string) {
	flag.IntVar(&opts.AfterNstring, "A", 0, "после каждой найденной строки дополнительно вывести N строк после неё")
	flag.IntVar(&opts.BeforeNstring, "B", 0, "вывести N строк до каждой найденной строки")
	flag.IntVar(&opts.RoundNstring, "C", 0, "вывести N строк контекста вокруг найденной строки (включает и до, и после; эквивалентно -A N -B N)")
	flag.BoolVar(&opts.C, "c", false, "выводить только то количество строк, что совпадающих с шаблоном ")
	flag.BoolVar(&opts.I, "i", false, "игнорировать регистр")
	flag.BoolVar(&opts.V, "v", false, "инвертировать фильтр: выводить строки, не содержащие шаблон.")
	flag.BoolVar(&opts.F, "F", false, "воспринимать шаблон как фиксированную строку, а не регулярное выражение")
	flag.BoolVar(&opts.N, "n", false, "выводить номер строки перед каждой найденной строкой")
	flag.Parse()
	args := flag.Args()
	pattern = args[0]
	files = args[1:]
	return
}
