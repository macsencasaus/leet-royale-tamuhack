package testrunner

var Q1 = reflectionData{
	numParams:  2,
	paramTypes: []string{"int", "int"},
	methodName: "add",
	numCases:   5,
	cases: [][][]string{{{"1"}, {"2"}},
		{{"-10"}, {"7"}},
		{{"2147483647"}, {"-2147483648"}},
		{{"-543"}, {"543"}},
		{{"-500"}, {"4500"}}},
	expectedResults: []string{"3", "-3", "-1", "0", "4000"},
	returnType:      "int",
}

var Q1Template = map[string]string{
	"python":     `def add(a:int, b:int):\n\n\n`,
	"c++":        `int add(int a, int b){\n\n\n}`,
	"javascript": `func add(a, b){\n\n\n}`,
}
