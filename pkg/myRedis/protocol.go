package pkg

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	LF    = "\r\n"
	array = "*"
)

func protocolParser(input []byte) []string {
	inputStr := string(input)
	var commands []string
	if strings.HasPrefix(inputStr, array) {
		commands = ParseArray(inputStr)
	}
	return commands
}

func ParseArray(inputStr string) []string {
	var res []string
	end := strings.Index(inputStr, "\r\n")
	s := inputStr[end+2:]
	for i := 0; i < len(s); i++ {
		if string(s[i]) == "$" {
			j := strings.Index(s[i:], "\r\n")
			n, _ := strconv.Atoi(s[i+1 : j+i])
			word := s[j+i+2 : j+i+2+n]
			res = append(res, word)
		}
	} // *3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar
	fmt.Println(res)
	return res
}

func createBulkResponse(s string) string {

	return "$" + strconv.Itoa(len(s)) + LF + s + LF
}

func createNullResponse() string {

	return "$-1" + LF
}
