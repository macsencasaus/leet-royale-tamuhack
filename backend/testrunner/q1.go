package testrunner

var Q1 = QuestionData{
	Title:        "Add 2 Numbers",
	Prompt:       "Write a function that takes two numbers as input and returns their sum.",
	numParams:    2,
	paramTypes:   []string{"int", "int"},
	methodName:   "add",
	NumCases:     5,
	VisibleCases: 3,
	Cases: [][][]string{{{"1"}, {"2"}},
		{{"-10"}, {"7"}},
		{{"2147483647"}, {"-2147483648"}},
		{{"-543"}, {"543"}},
		{{"-500"}, {"4500"}}},
	ExpectedResults: []string{"3", "-3", "-1", "0", "4000"},
	returnType:      "int",
	Templates: LanguageFunctionTemplates{
		Python:     "def add(a:int, b:int):\n    \n",
		Javascript: "function add(a, b) {\n    \n}",
		Cpp:        "int add(int a, int b) {\n    \n}",
	},
}
