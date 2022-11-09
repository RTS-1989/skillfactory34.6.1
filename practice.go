package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
)

func getInfoFromFile(fileName, reg string) ([][]string, error) {
	contentBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(reg)
	subMatches := re.FindAllStringSubmatch(string(contentBytes), -1)

	return subMatches, nil
}

func writeInfoToFile(firstValue, secondValue int, mathOperator string, writer *bufio.Writer) error {
	result := compute(firstValue, secondValue, mathOperator)
	stringToWrite := fmt.Sprintf("%d%s%d=%d\n", firstValue, mathOperator, secondValue, result)
	_, err := writer.WriteString(stringToWrite)
	if err != nil {
		return err
	}

	return nil
}

func compute(firstValue, secondValue int, mathOperator string) int {
	var result int
	switch mathOperator {
	case "+":
		result = firstValue + secondValue
	case "-":
		result = firstValue - secondValue
	case "*":
		result = firstValue * secondValue
	case "/":
		result = firstValue / secondValue
	}

	return result
}

func runWriteFiles(processedInfo [][]string, fileToWrite string) error {
	if err := os.Truncate(fileToWrite, 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}
	file, err := os.OpenFile(fileToWrite, os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	for n, entry := range processedInfo {
		mathOperand := entry[2]
		firstValue, err := strconv.Atoi(entry[1])
		if err != nil {
			return err
		}
		secondValue, err := strconv.Atoi(entry[3])
		if err != nil {
			return err
		}
		if n%5 == 0 {
			err = writer.Flush()
			if err != nil {
				return err
			}
		}
		err = writeInfoToFile(firstValue, secondValue, mathOperand, writer)
		if err != nil {
			return err
		}
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var sourceFileName, resultFileName string
	reg := `([\d]+)([\+|\-|*|\/])([\d]+)`

	switch len(os.Args) {
	case 1:
		log.Println("You should add file with source info and file for results to write!")
		return
	case 2:
		log.Println("You should add file for results to write!")
		return
	case 3:
		sourceFileName, resultFileName = os.Args[1], os.Args[2]
	}

	_, err := os.Open(resultFileName)
	if err != nil {
		_, err = os.Create(resultFileName)
		if err != nil {
			log.Println("Couldn't create file")
			return
		}
	}

	info, err := getInfoFromFile(sourceFileName, reg)

	if err != nil {
		log.Println(err)
		return
	}

	err = runWriteFiles(info, resultFileName)
	if err != nil {
		log.Println(err)
		return
	}
}
