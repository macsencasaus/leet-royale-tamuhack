package testrunner

import (
	// "fmt"
	"strconv"
	"strings"
)

// restictions for QuestionData
// -if put in a string with escape character, must do \\char instead of just \char
// -types can be bool, int, float, string, list type(has a space)
// -for boolean values, have them uppercase i.e. True, False
type QuestionData struct {
	Title           string
	numParams       int
	NumCases        int
    VisibleCases    int
	paramTypes      []string
	cases           [][][]string //[case][parameter][item] item is just 0 for non lists
	expectedResults []string     //if return a list, each test case starts with a number saying number of values then is the sequence of values
	methodName      string
	returnType      string
	Prompt          string
	Templates       LanguageFunctionTemplates
}

var Questions = []QuestionData{
	Q1, 
    Q2,
    Q3,
    Q4,
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

func collapseMe() {
	//following is test cases
	//python
	// fmt.Println(generate(`def add(a:int,b:int):
	// return a+b`, "python", "747474747", 1))

	// fmt.Println(generate(`def addLots(ls:list):
	// sum=0
	// for l in ls:
	// 	sum+=l
	// return [sum]`, "python", "747474747", 2))

	// fmt.Println(generate(`def returnList(ls:list):
	// return ls`, "python", "747474747", 3))

	// fmt.Println(generate(`def multipleParameters(a:int, b:float, ls:list, c:bool, d:str):
	// answer = str(a)+str(b)
	// for l in ls:
	// 	answer+=str(l)
	// answer+=str(c)+d
	// return answer
	// `,"python","747474747",4))

	//c++
	// fmt.Println(generate(`int add(int a, int b){
	// return a+b;}`, "c++", "747474747", 1))

	// fmt.Println(generate(`vector<int> addLots(vector<int> list){
	// int sum=0;
	// for(int i=0; i<list.size();i++)
	// sum+=list.at(i);
	// vector<int> answer = {sum};
	// return answer;}`, "c++", "747474747", 2))

	// fmt.Println(generate(`vector<int> returnList(vector<int> list){return list;}`, "c++", "747474747",3))

	// fmt.Println(generate(`string multipleParameters(int a, double b, vector<double> ls, bool c, string d){
	// string answer="";
	//     answer+=to_string(a);
	//     answer+=to_string(b);
	//     for(double i: ls)
	//             answer+=to_string(i);
	//     if(c)
	//         answer+=d;
	//     return answer;}
	// `, "c++", "747474747", 4))

	//javascript
	// fmt.Println(generate(`function add(a,b){
	// return a+b;}`, "javascript","747474747",1))

	// fmt.Println(generate(`function addLots(ls){
	// 	let total=0;
	// 	for(let i=0; i<ls.length; i++){
	// 		total+=ls[i];
	// 	}
	// 	return [total];}`, "javascript","747474747",2))

	// fmt.Println(generate(`function returnList(ls){
	// 	return ls;}`, "javascript","747474747",3))

	// fmt.Println(generate(`function multipleParameters(a, b, ls, c, d){
	// 	let answer=a+""+b
	// 		for(let i=0; i<ls.length;i++)
	// 			answer+=i;
	// 		answer+=c+d;
	// 		return answer;}
	// 	`, "javascript", "747474747", 4))
}

func generate(userInput string, language Language, magicNumber string, questionNumber int) string {

	r := Questions[questionNumber]
	if language == CPP {
		return generateC(userInput, magicNumber, r)
	} else if language == Python {
		return generatePython(userInput, magicNumber, r)
	} else if language == Javascript {
		return generateJavacript(userInput, magicNumber, r)
	}
	return "set language to either 'c++', 'python', or 'javascript'"
}

func isNotAList(typee string) bool {
	return !(len(typee) > 5 && typee[4] == ' ')
}

func generatePython(userInput, magicNumber string, r QuestionData) string {
	answer := userInput
	answer += "\n\n"
	//random number print method Prints Like:user_output \nmagic_number\n result \nmagic_number\n...
	answer += "def magic(thingToPrint):\n"
	answer += "\tprint('\\n',str('" + string(magicNumber) + "'),'\\n',thingToPrint,'\\n',str('" + string(magicNumber) + "'),'\\n',sep='',end='')\n\n"
	answer += "def main():\n"
	//constructs the expected results array
	if isNotAList(r.returnType) {
		answer += "\texpected_results = ["
		for i := 0; i < len(r.expectedResults); i++ {
			if r.returnType == "string" {
				answer += "\""
			}
			answer += r.expectedResults[i]
			if r.returnType == "string" {
				answer += "\""
			}
			if i+1 != len(r.expectedResults) {
				answer += ","
			}
		}
	} else { //because methods returns a list, constructs a 2D list  {[],[],[]} where each [] contains the list for a test case
		answer += "\texpected_results = ["
		currentIndex := 0
		for j := 0; j < r.NumCases; j++ { //j=current test case
			answer += "["
			lengthOfCurrentTestCase, _ := strconv.Atoi(r.expectedResults[currentIndex]) //get length of [] for the current test case
			currentIndex++
			for i := 0; i < lengthOfCurrentTestCase; i++ {
				if r.returnType == "list string" {
					answer += "\""
				}
				answer += r.expectedResults[currentIndex]
				currentIndex++
				if r.returnType == "list string" {
					answer += "\""
				}
				if i+1 != lengthOfCurrentTestCase {
					answer += ","
				}
			}
			answer += "]"
			if j+1 != r.NumCases {
				answer += ","
			}
		}
	}
	answer += "]\n"

	//initializes lists in order to put them in a bigger list later b/c easier than initializing all at once
	for i := 0; i < len(r.paramTypes); i++ { //i=currentParameter
		if !(isNotAList(r.paramTypes[i])) {
			//Lists have the name: a<testCaseNumber>_<parameterNumber>
			for j := 0; j < r.NumCases; j++ { //j==current_Case
				answer += "\ta" + strconv.Itoa(j) + "_" + strconv.Itoa(i) + " = ["
				for k := 0; k < len(r.cases[j][i]); k++ {
					if r.paramTypes[i] == "list string" {
						answer += "\""
					}
					answer += r.cases[j][i][k]
					if r.paramTypes[i] == "list string" {
						answer += "\""
					}
					if k+1 != len(r.cases[j][i]) {
						answer += ","
					}
				}
				answer += "]\n"
			}
		}
	}
	//constructs the array holding all of the parameters to be passed into the method
	answer += "\tcases = [["
	for i := 0; i < len(r.cases); i++ {
		for j := 0; j < len(r.cases[i]); j++ {
			if isNotAList(r.paramTypes[j]) { //if not a list just add the case data
				if r.paramTypes[j] == "string" {
					answer += "\""
				}
				answer += r.cases[i][j][0]
				if r.paramTypes[j] == "string" {
					answer += "\""
				}
			} else {
				answer += "a" + strconv.Itoa(i) + "_" + strconv.Itoa(j) //if are a list, add the name of the list constructed earlier
			}
			if j+1 != len(r.cases[i]) {
				answer += ","
			}
		}
		answer += "]"
		if i+1 != len(r.cases) {
			answer += ",["
		}
	}
	answer += "]\n"
	if isNotAList(r.returnType) { //This is used to determine how the test is validated.
		answer += "\tsimple_return = True\n"
	} else {
		answer += "\tsimple_return = False\n"
	}
	answer += "\tfor index, case in enumerate(cases):\n"
	answer += "\t\ttry:\n"
	//calls the method
	answer += "\t\t\tresult = "
	answer += r.methodName + "("
	for i := 0; i < r.numParams; i++ {
		answer += "case[" + strconv.Itoa(i) + "]"
		if i+1 != r.numParams {
			answer += ","
		}
	}
	answer += ")\n"
	answer += "\t\t\tif(simple_return):\n" //if return just a number, can do a simple comparison
	answer += "\t\t\t\tif result == expected_results[index]:\n"
	answer += "\t\t\t\t\tmagic('AC')\n"
	answer += "\t\t\t\telse:\n"
	answer += "\t\t\t\t\tmagic('WA')\n"
	answer += "\t\t\telse:\n"
	answer += "\t\t\t\tfailed=len(expected_results[index])!=len(result)\n"
	answer += "\t\t\t\tfor i in range(len(expected_results[index])):\n" //if return a list, compare every element in both lists
	answer += "\t\t\t\t\tif failed:\n"
	answer += "\t\t\t\t\t\tbreak\n"
	answer += "\t\t\t\t\tfailed = expected_results[index][i]!=result[i]\n"
	answer += "\t\t\t\tif failed:\n"
	answer += "\t\t\t\t\tmagic('WA')\n"
	answer += "\t\t\t\telse:\n"
	answer += "\t\t\t\t\tmagic('AC')\n"
	answer += "\t\texcept:\n"
	answer += "\t\t\tmagic('RE')\n"
	answer += "\nmain()"
	return answer
}

func generateC(userInput, magicNumber string, r QuestionData) string {
	answer := "#include<vector>\n#include<string>\n#include<iostream>\nusing namespace std;\n#include <tuple>\n\n"
	answer += userInput
	answer += "\n\n"
	//random number print method Prints Like:user_output \nmagic_number\n result \nmagic_number\n...
	answer += "void magic(string thingToPrint){\n"
	answer += "\tcout<<\"\\n\"<<" + magicNumber + "<<\"\\n\"<<thingToPrint<<\"\\n\"<<" + magicNumber + "<<\"\\n\";}\n\n"
	answer += "int main(){\n"
	//results array
	if isNotAList(r.returnType) {
		answer += "\tvector<" + pythonToC[r.returnType] + "> expected_results = {{"
		for i := 0; i < len(r.expectedResults); i++ {
			if r.returnType == "string" {
				answer += "\""
			}
			if r.returnType == "bool" {
				answer += strings.ToLower(string(r.expectedResults[i]))
			} else {
				answer += r.expectedResults[i]
			}
			if r.returnType == "string" {
				answer += "\""
			}
			if i+1 != len(r.expectedResults) {
				answer += "},{"
			} else {
				answer += "}"
			}
		}
	} else { //forming expected results as list of lists//results array
		answer += "\tvector<" + pythonToC[r.returnType] + "> expected_results = {"
		currentIndex := 0
		for j := 0; j < r.NumCases; j++ {
			answer += "{"
			lengthOfCurrentTestCase, _ := strconv.Atoi(r.expectedResults[currentIndex]) //get length of this test cases's array
			currentIndex++
			for i := 0; i < lengthOfCurrentTestCase; i++ {
				if r.returnType == "list string" {
					answer += "\""
				}
				if r.returnType == "list bool" {
					answer += strings.ToLower(string(r.expectedResults[currentIndex]))
				} else {
					answer += r.expectedResults[currentIndex]
				}
				if r.returnType == "list string" {
					answer += "\""
				}
				currentIndex++
				if i+1 != lengthOfCurrentTestCase {
					answer += ","
				}
			}
			answer += "}"
			if j+1 != r.NumCases {
				answer += ","
			}
		}
	}
	answer += "};\n"

	//initializes lists in order to put them in a bigger list later b/c easier than initializing all at once
	for i := 0; i < len(r.paramTypes); i++ { //i=current_paramter
		if !(isNotAList(r.paramTypes[i])) {
			//initializes lists and such with name: a<test_case_number>_<parameter_number>
			for j := 0; j < r.NumCases; j++ { //j==current_case
				answer += "\t" + pythonToC[r.paramTypes[i]] + " a" + strconv.Itoa(j) + "_" + strconv.Itoa(i) + " = {"
				for k := 0; k < len(r.cases[j][i]); k++ {
					if r.paramTypes[i] == "list string" {
						answer += "\""
					}
					if r.paramTypes[i] == "list bool" {
						answer += strings.ToLower(string(r.cases[j][i][k]))
					} else {
						answer += r.cases[j][i][k]
					}
					if r.paramTypes[i] == "list string" {
						answer += "\""
					}
					if k+1 != len(r.cases[j][i]) {
						answer += ","
					}
				}
				answer += "};\n"
			}

		}
	}
	//constructs the array holding all of the parameters to be passed into the method
	answer += "\tvector<tuple<"
	for i := 0; i < len(r.paramTypes); i++ {
		answer += pythonToC[r.paramTypes[i]]
		if i+1 != len(r.paramTypes) {
			answer += ", "
		}
	}
	answer += ">> cases = {make_tuple("
	for i := 0; i < len(r.cases); i++ {
		for j := 0; j < len(r.cases[i]); j++ {
			if isNotAList(r.paramTypes[j]) {
				if r.paramTypes[j] == "string" {
					answer += "\""
				}
				if r.paramTypes[j] == "bool" {
					answer += strings.ToLower(string(r.cases[i][j][0]))
				} else {
					answer += r.cases[i][j][0]
				}
				if r.paramTypes[j] == "string" {
					answer += "\""
				}
			} else {
				answer += "a" + strconv.Itoa(i) + "_" + strconv.Itoa(j)
			}
			if j+1 != len(r.cases[i]) {
				answer += ","
			}
		}
		answer += ")"
		if i+1 != len(r.cases) {
			answer += ", make_tuple("
		}
	}
	answer += "};\n"
	answer += "\tfor (int index=0; index<cases.size();index++){\n"
	answer += "\t\ttry{\n"
	answer += "\t\t\t" + pythonToC[r.returnType] + " result = "
	answer += r.methodName + "("
	for i := 0; i < r.numParams; i++ {
		answer += "get<" + strconv.Itoa(i) + ">(cases[index])"
		if i+1 != r.numParams {
			answer += ","
		}
	}
	answer += ");\n"
	if isNotAList(r.returnType) { //Do it this way instead of printing both b/c having both causes errors in c++. It is also cleaner this way.
		answer += "\t\t\t\tif (result == expected_results[index])\n"
		answer += "\t\t\t\t\tmagic(\"AC\");\n"
		answer += "\t\t\t\telse\n"
		answer += "\t\t\t\t\tmagic(\"WA\");}\n"
	} else {
		answer += "\t\t\t\tbool failed=expected_results[index].size()!=result.size();\n"
		answer += "\t\t\t\tfor (int i=0; i<expected_results[index].size();i++){\n"
		answer += "\t\t\t\t\tif (failed)\n"
		answer += "\t\t\t\t\t\tbreak;\n"
		answer += "\t\t\t\t\tfailed = expected_results[index][i]!=result[i];}\n"
		answer += "\t\t\t\tif (failed)\n"
		answer += "\t\t\t\t\tmagic(\"WA\");\n"
		answer += "\t\t\t\telse\n"
		answer += "\t\t\t\t\tmagic(\"AC\");}\n"
	}
	answer += "\t\tcatch(...){\n"
	answer += "\t\t\tmagic(\"RE\");}}}"
	return answer
}

func generateJavacript(userInput, magicNumber string, r QuestionData) string {
	answer := userInput
	answer += "\n\n"
	//random number print method Follows:user_output \nmagic_number\n result \nmagic_number\n user_output...
	answer += "function magic(thingToPrint){\n"
	answer += "\tconsole.log(\"\\n" + magicNumber + "\\n\"+thingToPrint+\"\\n" + magicNumber + "\");}\n\n"
	answer += "function main(){\n"
	//results array
	if isNotAList(r.returnType) {
		answer += "\tlet expected_results = ["
		for i := 0; i < len(r.expectedResults); i++ {
			if r.returnType == "string" {
				answer += "\""
			}
			if r.returnType == "bool" {
				answer += strings.ToLower(string(r.expectedResults[i]))
			} else {
				answer += r.expectedResults[i]
			}
			if r.returnType == "string" {
				answer += "\""
			}
			if i+1 != len(r.expectedResults) {
				answer += ","
			}
		}
	} else { //forming expected results as list of lists
		answer += "\tlet expected_results = ["
		currentIndex := 0
		for j := 0; j < r.NumCases; j++ {
			answer += "["
			lengthOfCurrentTestCase, _ := strconv.Atoi(r.expectedResults[currentIndex]) //get length of this test cases's array
			currentIndex++
			for i := 0; i < lengthOfCurrentTestCase; i++ {
				if r.returnType == "list string" {
					answer += "\""
				}
				if r.returnType == "list bool" {
					answer += strings.ToLower(string(r.expectedResults[currentIndex]))
				} else {
					answer += r.expectedResults[currentIndex]
				}
				if r.returnType == "list string" {
					answer += "\""
				}
				currentIndex++
				if i+1 != lengthOfCurrentTestCase {
					answer += ","
				}
			}
			answer += "]"
			if j+1 != r.NumCases {
				answer += ","
			}
		}
	}
	answer += "];\n"

	//parameter array
	for i := 0; i < len(r.paramTypes); i++ { //i=current_paramter
		if !(isNotAList(r.paramTypes[i])) {
			//initializes lists and such with name: a<test_case_number>_<parameter_number>
			for j := 0; j < r.NumCases; j++ { //j==current_Case
				answer += "\tlet a" + strconv.Itoa(j) + "_" + strconv.Itoa(i) + " = ["
				for k := 0; k < len(r.cases[j][i]); k++ {
					if r.paramTypes[i] == "list string" {
						answer += "\""
					}
					if r.paramTypes[i] == "list bool" {
						answer += strings.ToLower(string(r.cases[j][i][k]))
					} else {
						answer += r.cases[j][i][k]
					}
					if r.paramTypes[i] == "list string" {
						answer += "\""
					}
					if k+1 != len(r.cases[j][i]) {
						answer += ","
					}
				}
				answer += "];\n"
			}

		}
	}

	answer += "\tlet cases = [["
	for i := 0; i < len(r.cases); i++ {
		for j := 0; j < len(r.cases[i]); j++ {
			if isNotAList(r.paramTypes[j]) { //if regular
				if r.paramTypes[j] == "string" {
					answer += "\""
				}
				if r.paramTypes[j] == "bool" {
					answer += strings.ToLower(string(r.cases[i][j][0]))
				} else {
					answer += r.cases[i][j][0]
				}
				if r.paramTypes[j] == "string" {
					answer += "\""
				}
			} else {
				answer += "a" + strconv.Itoa(i) + "_" + strconv.Itoa(j)
			}
			if j+1 != len(r.cases[i]) {
				answer += ","
			}
		}
		answer += "]"
		if i+1 != len(r.cases) {
			answer += ",["
		}
	}
	answer += "];\n"
	if isNotAList(r.returnType) {
		answer += "\tlet simple_return = true;\n"
	} else {
		answer += "\tlet simple_return = false;\n"
	}
	answer += "\tfor (let index=0; index<cases.length; index++){\n"
	answer += "\t\ttry{\n"
	answer += "\t\t\tlet result = "
	answer += r.methodName + "("
	for i := 0; i < r.numParams; i++ {
		answer += "cases[index][" + strconv.Itoa(i) + "]"
		if i+1 != r.numParams {
			answer += ","
		}
	}
	answer += ")\n"
	answer += "\t\t\tif(simple_return)\n"
	answer += "\t\t\t\tif (result == expected_results[index])\n"
	answer += "\t\t\t\t\tmagic('AC');\n"
	answer += "\t\t\t\telse\n"
	answer += "\t\t\t\t\tmagic('WA');\n"
	answer += "\t\t\telse{\n"
	answer += "\t\t\t\tlet failed=expected_results[index].length!=result.length;\n"
	answer += "\t\t\t\tfor (let i=0; i<expected_results[index].length;i++){\n"
	answer += "\t\t\t\t\tif (failed)\n"
	answer += "\t\t\t\t\t\tbreak;\n"
	answer += "\t\t\t\t\tfailed = expected_results[index][i]!=result[i];}\n"
	answer += "\t\t\t\tif (failed)\n"
	answer += "\t\t\t\t\tmagic('WA');\n"
	answer += "\t\t\t\telse\n"
	answer += "\t\t\t\t\tmagic('AC');}}\n"
	answer += "\t\tcatch(error){\n"
	answer += "\t\t\tmagic('RE');}}}\n"
	answer += "\nmain();"
	return answer
}
