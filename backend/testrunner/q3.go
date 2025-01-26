package testrunner

var Q3 = reflectionData{
	numParams:  1,
	paramTypes: []string{"list int"},
	methodName: "returnList",
	numCases:   5,
	cases: [][][]string{{{"1", "2"}},
		{{"-10"}},
		{{"2147483647", "-2147483648", "1515125"}},
		{{}},
		{{"-500", "4500"}}},
	expectedResults: []string{"2", "1", "2", "1", "-10", "3", "2147483647", "-2147483648", "1515125", "0", "2", "-500", "4500"},
	returnType:      "list int",
}

var Q3Template = map[string]string{
	"python":     `def returnList(ls:list):\n\n\n`,
	"c++":        `vector<int> returnList(vector<int> list){\n\n\n}`,
	"javascript": `func returnList(ls){\n\n\n}`,
}
