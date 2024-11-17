package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
		return false
	}

	return false
}

func getByteNumber(fileName *string, stdin *string) int64 {
	if stdin != nil {
		byteStdin := []byte(*stdin)
		return int64(binary.Size(byteStdin))
	}

	if fileName == nil {
		panic("file name is required")
	}

	fileInfo, err := os.Stat(*fileName)
	if err != nil {
		panic(err.Error())
	}

	return fileInfo.Size()
}

func getNumberOfLines(fileName *string, stdin *string) int {
	if stdin != nil {
		count := 0
		lines := strings.NewReader(*stdin)
		scanner := bufio.NewScanner(lines)
		for scanner.Scan() {
			count++
		}
		return count
	}

	if fileName == nil {
		panic("file name is required if there not standard input")
	}

	file, err := os.Open(*fileName)
	if err != nil {
		panic(err.Error())
	}

	defer file.Close()
	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		count++
	}

	return count
}

func countWords(fileName string) int {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err.Error())
	}

	defer file.Close()
	wordCounter := 0
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		wordCounter++
	}

	return wordCounter
}

func counterCharacter(fileName string) int {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err.Error())
	}

	defer file.Close()
	characterCounter := 0
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		characterCounter++
	}

	return characterCounter
}

func ProcessFile(fileName string, cb func(string)) {
	if fileName == "" {
		panic("Must specify a file Name")
	}

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}

	pathsInTheCurrentDir := []string{}
	err = filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err.Error())
		} else {
			pathsInTheCurrentDir = append(pathsInTheCurrentDir, path)
		}

		return err
	})

	if err != nil {
		panic(err.Error())
	}

	if contains(pathsInTheCurrentDir, fileName) {
		panic("This file is found in the current directory")
	}

	cb(fileName)
}

func handleStandardInput() string {
	input := ""
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input += scanner.Text() + "\n"
	}
	return input
}

const (
	byteNumberFlag       = "c"
	linesCounterFlag     = "l"
	wordsCounterFlag     = "w"
	characterCounterFlag = "m"
)

func main() {

	args := os.Args[1:]

	if len(args) == 1 {
		flag_ := args[0][len(args[0])-1:]

		switch flag_ {
		case byteNumberFlag:
			stdin := handleStandardInput()
			fmt.Println(getByteNumber(nil, &stdin))
		case linesCounterFlag:
			stdin := handleStandardInput()
			fmt.Println(getNumberOfLines(nil, &stdin))
		case wordsCounterFlag:
			stdin := handleStandardInput()
			wordsSlice := strings.Fields(stdin)
			fmt.Println(len(wordsSlice))
		case characterCounterFlag:
			stdin := handleStandardInput()
			runes := utf8.RuneCountInString(stdin)
			fmt.Println(runes)
		default:
			ProcessFile(args[0], func(s string) {
				byteNumber := getByteNumber(&s, nil)
				lineNumber := getNumberOfLines(&s, nil)
				wordNumber := countWords(s)

				fmt.Println(byteNumber, lineNumber, wordNumber, s)
			})
		}
		return
	}

	byteNumber := flag.String(byteNumberFlag, "", "Flag to get the file byte number")
	lineNumber := flag.String(linesCounterFlag, "", "Flag to get the file line number")
	wordCounter := flag.String(wordsCounterFlag, "", "Flag to get the file word counter")
	characterCounter := flag.String(characterCounterFlag, "", "Flag to get the file character counter")

	flag.Parse()

	if *byteNumber != "" {
		ProcessFile(*byteNumber, func(s string) {
			fmt.Println(getByteNumber(&s, nil), s)
		})
		return
	}
	if *lineNumber != "" {
		ProcessFile(*lineNumber, func(s string) {
			fmt.Println(getNumberOfLines(&s, nil), s)
		})
		return
	}
	if *wordCounter != "" {
		ProcessFile(*wordCounter, func(s string) {
			fmt.Println(countWords(s), s)
		})
		return
	}
	if *characterCounter != "" {
		ProcessFile(*characterCounter, func(s string) {
			fmt.Println(counterCharacter(s), s)
		})
		return
	}

}