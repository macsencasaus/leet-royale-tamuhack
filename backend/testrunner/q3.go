package testrunner

var Q3 = QuestionData{
	Title:      "List",
	Prompt:     "Return the list.",
	numParams:  1,
	paramTypes: []string{"list int"},
	methodName: "returnList",
	NumCases:   5,
    VisibleCases: 3,
	cases: [][][]string{{{"1", "2"}},
		{{"-10"}},
		{{"2147483647", "-2147483648", "1515125"}},
		{{}},
		{{"-500", "4500"}}},
	expectedResults: []string{"2", "1", "2", "1", "-10", "3", "2147483647", "-2147483648", "1515125", "0", "2", "-500", "4500"},
	returnType:      "list int",
	Templates: LanguageFunctionTemplates{
		Python:     `def returnList(ls:list):\n\n\n`,
		Javascript: `func returnList(ls){\n\n\n}`,
		Cpp:        `vector<int> returnList(vector<int> list){\n\n\n}`,
	},
}
