package testrunner

import (
	"fmt"
	"reflect"
	"strings"
)

type QuestionData struct {
	Title        string
	Prompt       string
	FunctionName string
	Params       []Parameter
	ReturnType   TestType
	VisibleCases int
	Cases        []TestCase
	Templates    LanguageFunctionTemplates
}

type Parameter struct {
	Name string
	Type TestType
}

type TestCase struct {
	params []any
	result any
}

type TestCaseParams = []any

type TestType interface {
	String(Language) string
	Value(Language, any) string
}

type BoolType struct{}

func (BoolType) String(l Language) string {
	switch l {
	case CPP:
		return "bool"
	case Python:
		return "bool"
	case Javascript:
		return ""
	}
	return ""
}

func (BoolType) Value(l Language, v any) string {
	b, _ := v.(bool)

	switch l {
	case CPP, Javascript:
		if b {
			return "true"
		} else {
			return "false"
		}
	case Python:
		if b {
			return "True"
		} else {
			return "False"
		}
	}

	return ""
}

type IntType struct{}

func (IntType) String(l Language) string {
	switch l {
	case CPP:
		return "int"
	case Python:
		return "int"
	case Javascript:
		return ""
	}
	return ""
}

func (IntType) Value(_ Language, v any) string {
	i, _ := v.(int)
	return fmt.Sprintf("%d", i)
}

type StringType struct{}

func (StringType) String(l Language) string {
	switch l {
	case CPP:
		return "string"
	case Python:
		return "str"
	case Javascript:
		return ""
	}
	return ""
}

func (StringType) Value(_ Language, v any) string {
	s, _ := v.(string)
	return fmt.Sprintf(`"%s"`, s)
}

type ArrayType struct {
	inner TestType
}

func (at ArrayType) String(l Language) string {
	switch l {
	case CPP:
		return fmt.Sprintf("vector<%s>", at.inner.String(l))
	case Python:
		return fmt.Sprintf("list[%s]", at.inner.String(l))
	case Javascript:
		return ""
	}
	return ""
}

func (at ArrayType) Value(l Language, v any) string {
	a := toSliceOfAny(v)

	var sb strings.Builder

	switch l {
	case CPP:
		sb.WriteString(fmt.Sprintf("vector<%s>({", at.inner.String(CPP)))
	case Python, Javascript:
		sb.WriteString("[")
	}

	for i, e := range a {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(at.inner.Value(l, e))
	}

	switch l {
	case CPP:
		sb.WriteString("})")
	case Python, Javascript:
		sb.WriteString("]")
	}

	return sb.String()
}

var Questions = []QuestionData{
	// Q1,
	// Q2,
	// Q3,
	// Q4,
	// Q5,
	// Q6,
	// Q7,
	// Q8,
}

type LanguageFunctionTemplates struct {
	Python     string `json:"python"`
	Javascript string `json:"javascript"`
	Cpp        string `json:"cpp"`
}

var pythonToC = map[string]string{
	"bool": "bool", "int": "int", "string": "string", "float": "double", "list int": "vector<int>", "list float": "vector<double>",
	"list bool": "vector<bool>", "list string": "vector<string>",
}

// func generate(userInput string, language Language, magicNumber string, questionNumber int) string {
//
// 	r := Questions[questionNumber]
// 	if language == CPP {
// 		return generateC(userInput, magicNumber, r)
// 	} else if language == Python {
// 		return generatePython(userInput, magicNumber, r)
// 	} else if language == Javascript {
// 		return generateJavacript(userInput, magicNumber, r)
// 	}
// 	return "set language to either 'c++', 'python', or 'javascript'"
// }
//
// func isNotAList(typee string) bool {
// 	return !(len(typee) > 5 && typee[4] == ' ')
// }

func generateCPP(userInput, magic string, q QuestionData) string {
	var sb strings.Builder

	sb.WriteString("#include<bits/stdc++.h>\nusing namespace std;\n")

	return sb.String()
}

func generatePython(userInput, magic string, q QuestionData) string {
	newlines := "\n\n"

	var sb strings.Builder

	sb.WriteString(userInput)
	sb.WriteString(newlines)

	sb.WriteString("class TestCase:\n\tdef __init__(self, ")

	for _, param := range q.Params {
		sb.WriteString(param.Name + ": " + param.Type.String(Python))
		sb.WriteString(", ")
	}
	sb.WriteString("result: " + q.ReturnType.String(Python))

	sb.WriteString("):\n")

	for _, param := range q.Params {
		sb.WriteString(fmt.Sprintf("\t\tself.%s = %s\n", param.Name, param.Name))
	}
	sb.WriteString("\t\tself.result = result\n")

	sb.WriteString(newlines)

	sb.WriteString(fmt.Sprintf(`def magic(thingToPrint):
	print('\n', "%s", '\n', thingToPrint, '\n', "%s", sep='', end='')`,
		magic, magic))

	sb.WriteString(newlines)

	sb.WriteString("def main():\n")
	sb.WriteString("\tcases = [")

	for i, tc := range q.Cases {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("TestCase(")

		for j, p := range q.Params {
			v := p.Type.Value(Python, tc.params[j])
			sb.WriteString(fmt.Sprintf("%s=%s", p.Name, v))
			sb.WriteString(", ")
		}

		sb.WriteString(fmt.Sprintf("result=%s",
			q.ReturnType.Value(Python, tc.result)))

		sb.WriteString(")")
	}

	sb.WriteString("]\n")

	sb.WriteString("\tfor case in cases:\n")

	sb.WriteString("\t\ttry:\n")
	sb.WriteString(fmt.Sprintf("\t\t\tres = %s(", q.FunctionName))

	for i, p := range q.Params {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%[1]s=case.%[1]s", p.Name))
	}
	sb.WriteString(")\n")

	sb.WriteString("\t\texcept:\n")
	sb.WriteString("\t\t\tmagic(\"RE\")\n")
	sb.WriteString("\t\t\treturn\n\n")

	// TODO: implement deep compare
	sb.WriteString("\t\tif type(res) is type(case.result) and res == case.result:\n")
	sb.WriteString("\t\t\tmagic(\"AC\")\n")
	sb.WriteString("\t\telse:\n")
	sb.WriteString("\t\t\tmagic(\"WA\")\n")

	sb.WriteString("\nmain()\n")

	return sb.String()
}

// func generatePython(userInput, magicNumber string, r QuestionData) string {
// 	answer := userInput
// 	answer += "\n\n"
// 	//random number print method Prints Like:user_output \nmagic_number\n result \nmagic_number\n...
// 	answer += "def magic(thingToPrint):\n"
// 	answer += "\tprint('\\n',str('" + string(magicNumber) + "'),'\\n',thingToPrint,'\\n',str('" + string(magicNumber) + "'),'\\n',sep='',end='')\n\n"
// 	answer += "def main():\n"
// 	//constructs the expected results array
// 	if isNotAList(r.ReturnType) {
// 		answer += "\texpected_results = ["
// 		for i := 0; i < len(r.ExpectedResults); i++ {
// 			if r.ReturnType == "string" {
// 				answer += "\""
// 			}
// 			answer += r.ExpectedResults[i]
// 			if r.ReturnType == "string" {
// 				answer += "\""
// 			}
// 			if i+1 != len(r.ExpectedResults) {
// 				answer += ","
// 			}
// 		}
// 	} else { //because methods returns a list, constructs a 2D list  {[],[],[]} where each [] contains the list for a test case
// 		answer += "\texpected_results = ["
// 		currentIndex := 0
// 		for j := 0; j < r.NumCases; j++ { //j=current test case
// 			answer += "["
// 			lengthOfCurrentTestCase, _ := strconv.Atoi(r.ExpectedResults[currentIndex]) //get length of [] for the current test case
// 			currentIndex++
// 			for i := 0; i < lengthOfCurrentTestCase; i++ {
// 				if r.ReturnType == "list string" {
// 					answer += "\""
// 				}
// 				answer += r.ExpectedResults[currentIndex]
// 				currentIndex++
// 				if r.ReturnType == "list string" {
// 					answer += "\""
// 				}
// 				if i+1 != lengthOfCurrentTestCase {
// 					answer += ","
// 				}
// 			}
// 			answer += "]"
// 			if j+1 != r.NumCases {
// 				answer += ","
// 			}
// 		}
// 	}
// 	answer += "]\n"
//
// 	//initializes lists in order to put them in a bigger list later b/c easier than initializing all at once
// 	for i := 0; i < len(r.paramTypes); i++ { //i=currentParameter
// 		if !(isNotAList(r.paramTypes[i])) {
// 			//Lists have the name: a<testCaseNumber>_<parameterNumber>
// 			for j := 0; j < r.NumCases; j++ { //j==current_Case
// 				answer += "\ta" + strconv.Itoa(j) + "_" + strconv.Itoa(i) + " = ["
// 				for k := 0; k < len(r.Cases[j][i]); k++ {
// 					if r.paramTypes[i] == "list string" {
// 						answer += "\""
// 					}
// 					answer += r.Cases[j][i][k]
// 					if r.paramTypes[i] == "list string" {
// 						answer += "\""
// 					}
// 					if k+1 != len(r.Cases[j][i]) {
// 						answer += ","
// 					}
// 				}
// 				answer += "]\n"
// 			}
// 		}
// 	}
// 	//constructs the array holding all of the parameters to be passed into the method
// 	answer += "\tcases = [["
// 	for i := 0; i < len(r.Cases); i++ {
// 		for j := 0; j < len(r.Cases[i]); j++ {
// 			if isNotAList(r.paramTypes[j]) { //if not a list just add the case data
// 				if r.paramTypes[j] == "string" {
// 					answer += "\""
// 				}
// 				answer += r.Cases[i][j][0]
// 				if r.paramTypes[j] == "string" {
// 					answer += "\""
// 				}
// 			} else {
// 				answer += "a" + strconv.Itoa(i) + "_" + strconv.Itoa(j) //if are a list, add the name of the list constructed earlier
// 			}
// 			if j+1 != len(r.Cases[i]) {
// 				answer += ","
// 			}
// 		}
// 		answer += "]"
// 		if i+1 != len(r.Cases) {
// 			answer += ",["
// 		}
// 	}
// 	answer += "]\n"
// 	if isNotAList(r.ReturnType) { //This is used to determine how the test is validated.
// 		answer += "\tsimple_return = True\n"
// 	} else {
// 		answer += "\tsimple_return = False\n"
// 	}
// 	answer += "\tfor index, case in enumerate(cases):\n"
// 	answer += "\t\ttry:\n"
// 	//calls the method
// 	answer += "\t\t\tresult = "
// 	answer += r.MethodName + "("
// 	for i := 0; i < r.numParams; i++ {
// 		answer += "case[" + strconv.Itoa(i) + "]"
// 		if i+1 != r.numParams {
// 			answer += ","
// 		}
// 	}
// 	answer += ")\n"
// 	answer += "\t\t\tif(simple_return):\n" //if return just a number, can do a simple comparison
// 	answer += "\t\t\t\tif result == expected_results[index]:\n"
// 	answer += "\t\t\t\t\tmagic('AC')\n"
// 	answer += "\t\t\t\telse:\n"
// 	answer += "\t\t\t\t\tmagic('WA')\n"
// 	answer += "\t\t\telse:\n"
// 	answer += "\t\t\t\tfailed=len(expected_results[index])!=len(result)\n"
// 	answer += "\t\t\t\tfor i in range(len(expected_results[index])):\n" //if return a list, compare every element in both lists
// 	answer += "\t\t\t\t\tif failed:\n"
// 	answer += "\t\t\t\t\t\tbreak\n"
// 	answer += "\t\t\t\t\tfailed = expected_results[index][i]!=result[i]\n"
// 	answer += "\t\t\t\tif failed:\n"
// 	answer += "\t\t\t\t\tmagic('WA')\n"
// 	answer += "\t\t\t\telse:\n"
// 	answer += "\t\t\t\t\tmagic('AC')\n"
// 	answer += "\t\texcept:\n"
// 	answer += "\t\t\tmagic('RE')\n"
// 	answer += "\nmain()"
// 	return answer
// }

// func generateC(userInput, magicNumber string, r QuestionData) string {
// 	answer := "#include<vector>\n#include<string>\n#include<iostream>\nusing namespace std;\n#include <tuple>\n\n"
// 	answer += userInput
// 	answer += "\n\n"
// 	//random number print method Prints Like:user_output \nmagic_number\n result \nmagic_number\n...
// 	answer += "void magic(string thingToPrint){\n"
// 	answer += "\tcout<<\"\\n\"<<" + magicNumber + "<<\"\\n\"<<thingToPrint<<\"\\n\"<<" + magicNumber + "<<\"\\n\";}\n\n"
// 	answer += "int main(){\n"
// 	//results array
// 	if isNotAList(r.ReturnType) {
// 		answer += "\tvector<" + pythonToC[r.ReturnType] + "> expected_results = {{"
// 		for i := 0; i < len(r.ExpectedResults); i++ {
// 			if r.ReturnType == "string" {
// 				answer += "\""
// 			}
// 			if r.ReturnType == "bool" {
// 				answer += strings.ToLower(string(r.ExpectedResults[i]))
// 			} else {
// 				answer += r.ExpectedResults[i]
// 			}
// 			if r.ReturnType == "string" {
// 				answer += "\""
// 			}
// 			if i+1 != len(r.ExpectedResults) {
// 				answer += "},{"
// 			} else {
// 				answer += "}"
// 			}
// 		}
// 	} else { //forming expected results as list of lists//results array
// 		answer += "\tvector<" + pythonToC[r.ReturnType] + "> expected_results = {"
// 		currentIndex := 0
// 		for j := 0; j < r.NumCases; j++ {
// 			answer += "{"
// 			lengthOfCurrentTestCase, _ := strconv.Atoi(r.ExpectedResults[currentIndex]) //get length of this test cases's array
// 			currentIndex++
// 			for i := 0; i < lengthOfCurrentTestCase; i++ {
// 				if r.ReturnType == "list string" {
// 					answer += "\""
// 				}
// 				if r.ReturnType == "list bool" {
// 					answer += strings.ToLower(string(r.ExpectedResults[currentIndex]))
// 				} else {
// 					answer += r.ExpectedResults[currentIndex]
// 				}
// 				if r.ReturnType == "list string" {
// 					answer += "\""
// 				}
// 				currentIndex++
// 				if i+1 != lengthOfCurrentTestCase {
// 					answer += ","
// 				}
// 			}
// 			answer += "}"
// 			if j+1 != r.NumCases {
// 				answer += ","
// 			}
// 		}
// 	}
// 	answer += "};\n"
//
// 	//initializes lists in order to put them in a bigger list later b/c easier than initializing all at once
// 	for i := 0; i < len(r.paramTypes); i++ { //i=current_paramter
// 		if !(isNotAList(r.paramTypes[i])) {
// 			//initializes lists and such with name: a<test_case_number>_<parameter_number>
// 			for j := 0; j < r.NumCases; j++ { //j==current_case
// 				answer += "\t" + pythonToC[r.paramTypes[i]] + " a" + strconv.Itoa(j) + "_" + strconv.Itoa(i) + " = {"
// 				for k := 0; k < len(r.Cases[j][i]); k++ {
// 					if r.paramTypes[i] == "list string" {
// 						answer += "\""
// 					}
// 					if r.paramTypes[i] == "list bool" {
// 						answer += strings.ToLower(string(r.Cases[j][i][k]))
// 					} else {
// 						answer += r.Cases[j][i][k]
// 					}
// 					if r.paramTypes[i] == "list string" {
// 						answer += "\""
// 					}
// 					if k+1 != len(r.Cases[j][i]) {
// 						answer += ","
// 					}
// 				}
// 				answer += "};\n"
// 			}
//
// 		}
// 	}
// 	//constructs the array holding all of the parameters to be passed into the method
// 	answer += "\tvector<tuple<"
// 	for i := 0; i < len(r.paramTypes); i++ {
// 		answer += pythonToC[r.paramTypes[i]]
// 		if i+1 != len(r.paramTypes) {
// 			answer += ", "
// 		}
// 	}
// 	answer += ">> cases = {make_tuple("
// 	for i := 0; i < len(r.Cases); i++ {
// 		for j := 0; j < len(r.Cases[i]); j++ {
// 			if isNotAList(r.paramTypes[j]) {
// 				if r.paramTypes[j] == "string" {
// 					answer += "\""
// 				}
// 				if r.paramTypes[j] == "bool" {
// 					answer += strings.ToLower(string(r.Cases[i][j][0]))
// 				} else {
// 					answer += r.Cases[i][j][0]
// 				}
// 				if r.paramTypes[j] == "string" {
// 					answer += "\""
// 				}
// 			} else {
// 				answer += "a" + strconv.Itoa(i) + "_" + strconv.Itoa(j)
// 			}
// 			if j+1 != len(r.Cases[i]) {
// 				answer += ","
// 			}
// 		}
// 		answer += ")"
// 		if i+1 != len(r.Cases) {
// 			answer += ", make_tuple("
// 		}
// 	}
// 	answer += "};\n"
// 	answer += "\tfor (int index=0; index<cases.size();index++){\n"
// 	answer += "\t\ttry{\n"
// 	answer += "\t\t\t" + pythonToC[r.ReturnType] + " result = "
// 	answer += r.MethodName + "("
// 	for i := 0; i < r.numParams; i++ {
// 		answer += "get<" + strconv.Itoa(i) + ">(cases[index])"
// 		if i+1 != r.numParams {
// 			answer += ","
// 		}
// 	}
// 	answer += ");\n"
// 	if isNotAList(r.ReturnType) { //Do it this way instead of printing both b/c having both causes errors in c++. It is also cleaner this way.
// 		answer += "\t\t\t\tif (result == expected_results[index])\n"
// 		answer += "\t\t\t\t\tmagic(\"AC\");\n"
// 		answer += "\t\t\t\telse\n"
// 		answer += "\t\t\t\t\tmagic(\"WA\");}\n"
// 	} else {
// 		answer += "\t\t\t\tbool failed=expected_results[index].size()!=result.size();\n"
// 		answer += "\t\t\t\tfor (int i=0; i<expected_results[index].size();i++){\n"
// 		answer += "\t\t\t\t\tif (failed)\n"
// 		answer += "\t\t\t\t\t\tbreak;\n"
// 		answer += "\t\t\t\t\tfailed = expected_results[index][i]!=result[i];}\n"
// 		answer += "\t\t\t\tif (failed)\n"
// 		answer += "\t\t\t\t\tmagic(\"WA\");\n"
// 		answer += "\t\t\t\telse\n"
// 		answer += "\t\t\t\t\tmagic(\"AC\");}\n"
// 	}
// 	answer += "\t\tcatch(...){\n"
// 	answer += "\t\t\tmagic(\"RE\");}}}"
// 	return answer
// }

// func generateJavacript(userInput, magicNumber string, r QuestionData) string {
// 	answer := userInput
// 	answer += "\n\n"
// 	//random number print method Follows:user_output \nmagic_number\n result \nmagic_number\n user_output...
// 	answer += "function magic(thingToPrint){\n"
// 	answer += "\tconsole.log(\"\\n" + magicNumber + "\\n\"+thingToPrint+\"\\n" + magicNumber + "\");}\n\n"
// 	answer += "function main(){\n"
// 	//results array
// 	if isNotAList(r.ReturnType) {
// 		answer += "\tlet expected_results = ["
// 		for i := 0; i < len(r.ExpectedResults); i++ {
// 			if r.ReturnType == "string" {
// 				answer += "\""
// 			}
// 			if r.ReturnType == "bool" {
// 				answer += strings.ToLower(string(r.ExpectedResults[i]))
// 			} else {
// 				answer += r.ExpectedResults[i]
// 			}
// 			if r.ReturnType == "string" {
// 				answer += "\""
// 			}
// 			if i+1 != len(r.ExpectedResults) {
// 				answer += ","
// 			}
// 		}
// 	} else { //forming expected results as list of lists
// 		answer += "\tlet expected_results = ["
// 		currentIndex := 0
// 		for j := 0; j < r.NumCases; j++ {
// 			answer += "["
// 			lengthOfCurrentTestCase, _ := strconv.Atoi(r.ExpectedResults[currentIndex]) //get length of this test cases's array
// 			currentIndex++
// 			for i := 0; i < lengthOfCurrentTestCase; i++ {
// 				if r.ReturnType == "list string" {
// 					answer += "\""
// 				}
// 				if r.ReturnType == "list bool" {
// 					answer += strings.ToLower(string(r.ExpectedResults[currentIndex]))
// 				} else {
// 					answer += r.ExpectedResults[currentIndex]
// 				}
// 				if r.ReturnType == "list string" {
// 					answer += "\""
// 				}
// 				currentIndex++
// 				if i+1 != lengthOfCurrentTestCase {
// 					answer += ","
// 				}
// 			}
// 			answer += "]"
// 			if j+1 != r.NumCases {
// 				answer += ","
// 			}
// 		}
// 	}
// 	answer += "];\n"
//
// 	//parameter array
// 	for i := 0; i < len(r.paramTypes); i++ { //i=current_paramter
// 		if !(isNotAList(r.paramTypes[i])) {
// 			//initializes lists and such with name: a<test_case_number>_<parameter_number>
// 			for j := 0; j < r.NumCases; j++ { //j==current_Case
// 				answer += "\tlet a" + strconv.Itoa(j) + "_" + strconv.Itoa(i) + " = ["
// 				for k := 0; k < len(r.Cases[j][i]); k++ {
// 					if r.paramTypes[i] == "list string" {
// 						answer += "\""
// 					}
// 					if r.paramTypes[i] == "list bool" {
// 						answer += strings.ToLower(string(r.Cases[j][i][k]))
// 					} else {
// 						answer += r.Cases[j][i][k]
// 					}
// 					if r.paramTypes[i] == "list string" {
// 						answer += "\""
// 					}
// 					if k+1 != len(r.Cases[j][i]) {
// 						answer += ","
// 					}
// 				}
// 				answer += "];\n"
// 			}
//
// 		}
// 	}
//
// 	answer += "\tlet cases = [["
// 	for i := 0; i < len(r.Cases); i++ {
// 		for j := 0; j < len(r.Cases[i]); j++ {
// 			if isNotAList(r.paramTypes[j]) { //if regular
// 				if r.paramTypes[j] == "string" {
// 					answer += "\""
// 				}
// 				if r.paramTypes[j] == "bool" {
// 					answer += strings.ToLower(string(r.Cases[i][j][0]))
// 				} else {
// 					answer += r.Cases[i][j][0]
// 				}
// 				if r.paramTypes[j] == "string" {
// 					answer += "\""
// 				}
// 			} else {
// 				answer += "a" + strconv.Itoa(i) + "_" + strconv.Itoa(j)
// 			}
// 			if j+1 != len(r.Cases[i]) {
// 				answer += ","
// 			}
// 		}
// 		answer += "]"
// 		if i+1 != len(r.Cases) {
// 			answer += ",["
// 		}
// 	}
// 	answer += "];\n"
// 	if isNotAList(r.ReturnType) {
// 		answer += "\tlet simple_return = true;\n"
// 	} else {
// 		answer += "\tlet simple_return = false;\n"
// 	}
// 	answer += "\tfor (let index=0; index<cases.length; index++){\n"
// 	answer += "\t\ttry{\n"
// 	answer += "\t\t\tlet result = "
// 	answer += r.MethodName + "("
// 	for i := 0; i < r.numParams; i++ {
// 		answer += "cases[index][" + strconv.Itoa(i) + "]"
// 		if i+1 != r.numParams {
// 			answer += ","
// 		}
// 	}
// 	answer += ")\n"
// 	answer += "\t\t\tif(simple_return)\n"
// 	answer += "\t\t\t\tif (result == expected_results[index])\n"
// 	answer += "\t\t\t\t\tmagic('AC');\n"
// 	answer += "\t\t\t\telse\n"
// 	answer += "\t\t\t\t\tmagic('WA');\n"
// 	answer += "\t\t\telse{\n"
// 	answer += "\t\t\t\tlet failed=expected_results[index].length!=result.length;\n"
// 	answer += "\t\t\t\tfor (let i=0; i<expected_results[index].length;i++){\n"
// 	answer += "\t\t\t\t\tif (failed)\n"
// 	answer += "\t\t\t\t\t\tbreak;\n"
// 	answer += "\t\t\t\t\tfailed = expected_results[index][i]!=result[i];}\n"
// 	answer += "\t\t\t\tif (failed)\n"
// 	answer += "\t\t\t\t\tmagic('WA');\n"
// 	answer += "\t\t\t\telse\n"
// 	answer += "\t\t\t\t\tmagic('AC');}}\n"
// 	answer += "\t\tcatch(error){\n"
// 	answer += "\t\t\tmagic('RE');}}}\n"
// 	answer += "\nmain();"
// 	return answer
// }

func toSliceOfAny(v any) []any {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice {
		panic("value is not a slice")
	}

	result := make([]any, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		result[i] = rv.Index(i).Interface()
	}
	return result
}
