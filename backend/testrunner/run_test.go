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
        return 0;
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

func TestRunCpp_sadRun(t *testing.T) {
	file := `
    #include <iostream>
    int main(int argc, char** argv) {
        int* a = 0x0;
        std::cout << *a;
        return 0;
    }
    `
	magic := generateMagic()
	_, stage, err := runCpp([]byte(file), magic)
	if err != nil {
		t.Errorf("Err was non-nil: %v\n", err)
	}
	if stage != Run {
		t.Errorf("Program failed somewhere unexpected (stage): %#v\n", stage)
	}
}

func TestRunCpp_sadRunTime(t *testing.T) {
	file := `
    #include <iostream>
    int main(int argc, char** argv) {
        for (int i = 0; i > -1; i++) {
            std::cout << "";
        }
        return 0;
    }
    `
	magic := generateMagic()
	_, stage, err := runCpp([]byte(file), magic)
	if err != nil {
		t.Errorf("Err was non-nil: %v\n", err)
	}
	if stage != RunTime {
		t.Errorf("Program failed somewhere unexpected (stage): %#v\n", stage)
	}
}

func TestRunPython_happy(t *testing.T) {
	file := `
import sys
print(sys.argv[1], end="")
    `
	magic := generateMagic()
	magicString := magicToString(magic)
	out, stage, err := runPython([]byte(file), magic)
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

func TestRunPython_sadRun(t *testing.T) {
	file := `
import sys
print(sys.argv[3], end="")
    `
	magic := generateMagic()
	_, stage, err := runPython([]byte(file), magic)
	if err != nil {
		t.Errorf("Err was non-nil: %v\n", err)
	}
	if stage != Run {
		t.Errorf("Program failed somewhere unexpected (stage): %#v\n", stage)
	}
}

func TestRunJavascript_happy(t *testing.T) {
	file := `
process.stdout.write(process.argv[3])
    `
	magic := generateMagic()
	magicString := magicToString(magic)
	out, stage, err := runJavascript([]byte(file), magic)
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

func TestRunJavascript_sadRun(t *testing.T) {
	file := `
process.stdout.write(process.argv[5][0])
    `
	magic := generateMagic()
	_, stage, err := runJavascript([]byte(file), magic)
	if err != nil {
		t.Errorf("Err was non-nil: %v\n", err)
	}
	if stage != Run {
		t.Errorf("Program failed somewhere unexpected (stage): %#v\n", stage)
	}
}

func TestRunProblemTest_happy(t *testing.T) {
	file := `
    #include <iostream>
    int main(int argc, char** argv) {
        std::cout
            << "OUTPUT 1"
            << "\n" << argv[1] << "\n"
            << "AC"
            << "\n" << argv[1] << "\n"
            << "output 2 the great line of text and some more"
            << "\n" << argv[1] << "\n"
            << "AC"
            << "\n" << argv[1] << "\n"
        ;
    }
    `
	magic := generateMagic()
	res, err := RunProblemTest([]byte(file), CPP, magic)
	if err != nil {
		t.Errorf("Err was non-nil: %v\n", err)
	}
	if res.Issue != Success {
		t.Errorf("Program failed somewhere unexpected (stage): %#v\n", res.Issue)
	}
	if res.NCasesRun != 2 {
		t.Errorf("Incorrect number of parsed cases (%d) expected %d\n", res.NCasesRun, 2)
	}
	for i := 0; i < res.NCasesRun; i++ {
		if res.PFStatus[i] != AC {
			t.Errorf("Parsed incorrect case status '%s' expected '%s'\n", res.PFStatus[i], AC)
		}
	}
}

func Test_q1(t *testing.T) {
    file :=`
def add(a, b):
    return a + b
`
    res, err := RunTest([]byte(file), Python, 1)
    t.Log(res, err)
}
