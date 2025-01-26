package testrunner

var Q4 = QuestionData{
	Title:        "Concat",
	Prompt:       "Concatenate all the elements into a string",
	numParams:    5,
	paramTypes:   []string{"int", "float", "list float", "bool", "string"},
	methodName:   "multipleParameters",
	NumCases:     5,
	VisibleCases: 3,
	cases: [][][]string{{{"1"}, {"2.5"}, {"1", "2", "3", "4", "5"}, {"True"}, {"blahblahblah"}},
		{{"1"}, {"2.5"}, {"1", "2", "3", "4", "5"}, {"False"}, {"atwatwa"}},
		{{"5"}, {"-2.5"}, {}, {"True"}, {""}},
		{{"-1400"}, {"-400.5"}, {"1", "2"}, {"False"}, {"dwa"}},
		{{"1784214214"}, {"2.142145"}, {"1"}, {"True"}, {"wiodjahdukjwhyu_W o-)A* 902jbjk mbak \\n"}}},
	expectedResults: []string{"12.512345Trueblahblahblah", "12.512345Falseatwatwa", "5-2.5True", "-1400-400.512Falsedwa", "17842142142.1421451Truewiodjahdukjwhyu_W o-)A* 902jbjk mbak \\n"},
	returnType:      "string",
	Templates: LanguageFunctionTemplates{
		Python:     `def multipleParameters(a:int, b:float, ls:list, c:bool, d:str):\n\n\n`,
		Javascript: `func returnList(a, b, ls, c, d){\n\n\n}`,
		Cpp:        `string returnList(int a, double b, vector<int> ls, bool c, string d){\n\n\n}`,
	},
}
