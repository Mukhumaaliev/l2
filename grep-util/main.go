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
	var res int
	opts, files, pattern := ParseFlags()
	lines, err := ReadLines(files)
	if errors.Is(err, fmt.Errorf("не удалось открыть файл %q: %w", files[0], err)) {
		panic(err)
	}
	if opts.C {
		res = CountLine(pattern, lines, opts)
	}
	fmt.Print(res)
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

func CountWithoutPattern(pattern string, line string, opts options, count *int) (result string) {
	if opts.I {
		re := regexp.MustCompile(strings.ToLower(pattern))
		if !re.MatchString(strings.ToLower(line)) {
			result = line
			*count++
		}
	} else {
		re := regexp.MustCompile(pattern)
		if !re.MatchString(line) {
			result = line
			*count++
		}
	}
	return result
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

func ParseFlags() (opts options, files []string, sample string) {
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
	sample = args[0]
	files = args[1:]
	return
}
