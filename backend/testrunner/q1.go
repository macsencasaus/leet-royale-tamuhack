package testrunner

var Q1 = QuestionData{
	Title:        "Add 2 Numbers",
	Prompt:       "Write a function that takes two numbers as input and returns their sum.",
	numParams:    2,
	paramTypes:   []string{"int", "int"},
	methodName:   "add",
	NumCases:     5,
	VisibleCases: 3,
	cases: [][][]string{{{"1"}, {"2"}},
		{{"-10"}, {"7"}},
		{{"2147483647"}, {"-2147483648"}},
		{{"-543"}, {"543"}},
		{{"-500"}, {"4500"}}},
	expectedResults: []string{"3", "-3", "-1", "0", "4000"},
	returnType:      "int",
	Templates: LanguageFunctionTemplates{
		Python:     `def add(a:int, b:int):\n\n\n`,
		Javascript: `func add(a, b){\n\n\n}`,
		Cpp:        `int add(int a, int b){\n\n\n}`,
	},
}
