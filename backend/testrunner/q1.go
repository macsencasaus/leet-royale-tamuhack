package testrunner

var Q1 = QuestionData{
	Title:        "Add 2 Numbers",
	Prompt:       "Write a function that takes two numbers as input and returns their sum.",
	FunctionName:   "add",
	Params:       []Parameter{{"a", IntType{}}, {"b", IntType{}}},
	ReturnType:   IntType{},
	VisibleCases: 3,
	Cases: []TestCase{
		{params: TestCaseParams{1, 2}, result: 3},
		{params: TestCaseParams{3, 4}, result: 7},
	},
	Templates: LanguageFunctionTemplates{
		Python:     "def add(a:int, b:int):\n    \n",
		Javascript: "function add(a, b) {\n    \n}",
		Cpp:        "int add(int a, int b) {\n    \n}",
	},
}
