package testrunner

import (
	"log"
	"testing"
)

func Test_generatePython(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		userInput   string
		magicNumber string
		q           QuestionData
	}{
		{
			name:        "bruh",
			userInput:   "def add(a: int, b: int) -> int:\n\treturn a + b",
			magicNumber: "AAAAA",
			q:           Q1,
		},
		{
			name:        "bruh",
			userInput:   "def sum(arr: list[int]) -> int:\n\treturn 0",
			magicNumber: "AAAAA",
			q:           Q2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generatePython(tt.userInput, tt.magicNumber, tt.q)
			log.Println("\n" + got)
		})
	}
}
