package testrunner

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"time"
)

type ErrorStage string
type Language int
type TestCaseStatus string

const (
	Compile     ErrorStage = "Compile"
	CompileTime            = "CompileTime"
	Run                    = "Run"
	RunTime                = "RunTime"
	Success                = "Success"
)

const (
	CPP Language = iota
	Python
	Javascript
)

const (
	AC  TestCaseStatus = "Correct"
	WA                 = "Wrong"
	RE                 = "Runtime Error"
	TLE                = "Time Limit Exceeded"
)

type Result struct {
	Issue     ErrorStage
	NCasesRun int
	Stdout    [][]byte
	PFStatus  []TestCaseStatus
}

func (r Result) NCorrect() (int, int) {
	correct := 0
	for _, status := range r.PFStatus {
		if status == AC {
			correct++
		}
	}
	return correct, len(r.PFStatus)
}

func statusFromCode(letters []byte) TestCaseStatus {
	if bytes.Compare(letters, []byte("AC")) == 0 {
		return AC
	} else if bytes.Compare(letters, []byte("WA")) == 0 {
		return WA
	} else if bytes.Compare(letters, []byte("RE")) == 0 {
		return RE
	} else if bytes.Compare(letters, []byte("TLE")) == 0 {
		return TLE
	}
	log.Fatalf("Invalid test case status conversion from '%s'", letters)
	return AC
}

func generateMagic() int64 {
	return rand.Int63()
}

func magicToString(magic int64) string {
	return fmt.Sprintf("%d", magic)
}

func runCpp(fileContent []byte, magic int64) ([]byte, ErrorStage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	compFile := fmt.Sprintf("/tmp/testrunner-%d", magic)

	cmd := exec.CommandContext(ctx, "clang++", "--std=c++17", "-O3", "-fsanitize=address", "-Werror", "-x", "c++", "-o", compFile, "-")
	cmd.Stdin = bytes.NewReader(fileContent)

	out_ch := make(chan struct {
		out []byte
		err error
	})

	go func() {
		comOut, err := cmd.CombinedOutput()
		out_ch <- struct {
			out []byte
			err error
		}{comOut, err}
	}()

	select {
	case <-ctx.Done():
		return (<-out_ch).out, CompileTime, nil
	case out := <-out_ch:
		if out.err != nil {
			if _, ok := out.err.(*exec.ExitError); ok {
				return out.out, Compile, nil
			} else {
				return out.out, Compile, out.err
			}
		}
		// out.err is nil, progress normally
	}

	magicString := magicToString(magic)

	ctx, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	go func() {
		cmd = exec.CommandContext(ctx, compFile, magicString)
		runOut, err := cmd.CombinedOutput()
		out_ch <- struct {
			out []byte
			err error
		}{runOut, err}
	}()

	select {
	case <-ctx.Done():
		return (<-out_ch).out, RunTime, nil
	case out := <-out_ch:
		if out.err != nil {
			if v, ok := out.err.(*exec.ExitError); ok {
				if v.ExitCode() != 0 {
					return out.out, Run, nil
				} else {
					return out.out, Success, nil
				}
			} else {
				return out.out, Run, out.err
			}
		}
		return out.out, Success, nil
	}

	// this is unreachable
}

func runPython(fileContent []byte, magic int64) ([]byte, ErrorStage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	magicString := magicToString(magic)

	out_ch := make(chan struct {
		out []byte
		err error
	})

	go func() {
		cmd := exec.CommandContext(ctx, "python3", "-", magicString)
		cmd.Stdin = bytes.NewReader(fileContent)
		runOut, err := cmd.CombinedOutput()
		out_ch <- struct {
			out []byte
			err error
		}{runOut, err}
	}()

	select {
	case <-ctx.Done():
		return (<-out_ch).out, RunTime, nil
	case out := <-out_ch:
		if out.err != nil {
			if v, ok := out.err.(*exec.ExitError); ok {
				if v.ExitCode() != 0 {
					return out.out, Run, nil
				} else {
					return out.out, Success, nil
				}
			} else {
				return out.out, Run, out.err
			}
		}
		return out.out, Success, nil
	}

	// unreachable
}

func runJavascript(fileContent []byte, magic int64) ([]byte, ErrorStage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	magicString := magicToString(magic)

	out_ch := make(chan struct {
		out []byte
		err error
	})

	go func() {
		cmd := exec.CommandContext(ctx, "node", "-", "--", magicString)
		cmd.Stdin = bytes.NewReader(fileContent)
		runOut, err := cmd.CombinedOutput()
		out_ch <- struct {
			out []byte
			err error
		}{runOut, err}
	}()

	select {
	case <-ctx.Done():
		return (<-out_ch).out, RunTime, nil
	case out := <-out_ch:
		if out.err != nil {
			if v, ok := out.err.(*exec.ExitError); ok {
				if v.ExitCode() != 0 {
					return out.out, Run, nil
				} else {
					return out.out, Success, nil
				}
			} else {
				return out.out, Run, out.err
			}
		}
		return out.out, Success, nil
	}

	// unreachable
}

// Expects in the format `STDOUT` "\n" `MAGIC` "\n" `INFO` "\n" `MAGIC` "\n" ...
func RunProblemTest(fileContent []byte, lang Language, magic int64) (Result, error) {
	magicString := fmt.Sprintf("\n%d\n", magic)

	var streamOut []byte
	var stage ErrorStage
	var err error

	switch lang {
	case CPP:
		streamOut, stage, err = runCpp(fileContent, magic)
	case Javascript:
		streamOut, stage, err = runJavascript(fileContent, magic)
	case Python:
		streamOut, stage, err = runPython(fileContent, magic)
	default:
		log.Fatal("Language not implemented")
	}

	if err != nil {
		log.Fatal(err)
	}

	var testCasesRun int
	var testCaseProgramOut [][]byte
	var testCaseStatus []TestCaseStatus

	if stage == Compile || stage == CompileTime {
		testCaseProgramOut = append(testCaseProgramOut, streamOut)
		return Result{stage, 0, testCaseProgramOut, nil}, nil
	}

	var sections [][]byte = bytes.Split(streamOut, []byte(magicString))
	for i := 0; i < (len(sections)/2)*2; i += 2 {
		testCasesRun++
		testCaseProgramOut = append(testCaseProgramOut, sections[i])
		testCaseStatus = append(testCaseStatus, statusFromCode(sections[i+1]))
	}

	return Result{stage, testCasesRun, testCaseProgramOut, testCaseStatus}, nil
}

func RunTest(infile []byte, lang Language, question int) (Result, error) {
	magic := generateMagic()
	magicString := fmt.Sprintf("%d", magic)

	file := []byte(generate(string(infile), lang, magicString, question))
	return RunProblemTest(file, lang, magic)
}
