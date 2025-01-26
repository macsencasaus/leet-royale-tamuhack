package testrunner

var Q2 = reflectionData{
	numParams:  1,
	paramTypes: []string{"list int"},
	methodName: "addLots",
	numCases:   5,
	cases: [][][]string{{{"1", "2"}},
		{{"-10", "7"}},
		{{"2147483647", "-2147483648"}},
		{{"-543", "543"}},
		{{"-500", "4500"}}},
	expectedResults: []string{"1", "3", "1", "-3", "1", "-1", "1", "0", "1", "4000"},
	returnType:      "list int",
}

var Q2Template = map[string]string{
	"python":     `def addLots(ls:list):\n\n\n`,
	"c++":        `vector<int> addLots(vector<int> list){\n\n\n}`,
	"javascript": `func addLots(ls){\n\n\n}`,
}
