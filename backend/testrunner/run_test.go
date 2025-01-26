package testrunner

import (
	"bytes"
	"testing"
)

func TestRunCpp_happy(t *testing.T) {
	file := `
    #include <iostream>
    int main(int argc, char** argv) {
        std::cout << argv[1];
    }
    `
	magic := generateMagic()
    magicString := magicToString(magic)
	out, stage, err := runCpp([]byte(file), magic)
	if err != nil {
		t.Errorf("Err was non-nil: %v\n", err)
	}
	if stage != Success {
		t.Errorf("Program failed somewhere unexpected (stage): %#v\n", stage)
	}
	if bytes.Compare(out, []byte(magicString)) != 0 {
		t.Errorf("Magic bytes and output were not the same:\n (Expected)       vs. (Actual)        \n'%16s' '%16s'\n", magicString, out)
	}
}

func TestRunCpp_sadCompile(t *testing.T) {
	file := `
    int main(int argc, char** argv) {
        std::cout << argv[1];
    }
    `
	magic := generateMagic()
	out, stage, err := runCpp([]byte(file), magic)
	if err != nil {
		t.Errorf("Err was non-nil: %v\n", err)
	}
	if stage != Compile {
		t.Errorf("Program failed somewhere unexpected (stage): %#v\n", stage)
	}
    search := []byte("use of undeclared identifier")
	if bytes.Index(out, search) == -1 {
		t.Errorf("Output did not have error:\n%s\n>>Should have contained '%s'\n", out, search)
	}
}
