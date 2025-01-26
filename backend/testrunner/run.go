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

type Result struct {
	Issue     ErrorStage
	NCasesRun int
	Stdout    [][]byte
	PFStatus  []bool
}

func generateMagic() int64 {
	return rand.Int63()
}

func magicToString(magic int64) string {
	return fmt.Sprintf("%d", magic)
}

func runCpp(fileContent []byte, magic int64) ([]byte, ErrorStage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	compFile := fmt.Sprintf("/tmp/testrunner-%d", magic)

	cmd := exec.CommandContext(ctx, "clang++", "--std=c++17", "-x", "c++", "-o", compFile, "-")
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

// Expects in the format `STDOUT` "\n" `MAGIC` "\n" `INFO`
func RunProblemTest(fileContent []byte, lang Language) (Result, error) {
	var magic int64
	magic = 9876543210
	magicString := fmt.Sprintf("\n%d\n", magic)

	var streamOut []byte
	var stage ErrorStage
	var err error

	switch lang {
	case CPP:
		streamOut, stage, err = runCpp(fileContent, magic)
		// TODO
	default:
		log.Fatal("Language not implemented")
	}

	if err != nil {
		log.Fatal(err)
	}

	var testCasesRun int
	var testCaseProgramOut [][]byte
	var testCaseInfo [][]byte
	var index int
	for idx := bytes.Index(streamOut[index:], []byte(magicString)); idx < len(streamOut); {
		testCasesRun++
		testCaseProgramOut = append(testCaseProgramOut, streamOut[index:idx])
		index = bytes.Index(streamOut[idx+1:], []byte(magicString))
		testCaseInfo = append(testCaseInfo, streamOut[idx:index])
	}

	// parse
	var testCaseStatus []bool

	return Result{stage, testCasesRun, testCaseProgramOut, testCaseStatus}, nil
}
