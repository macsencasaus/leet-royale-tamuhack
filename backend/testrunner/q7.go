package testrunner

var Q7 = QuestionData{
	Title:        "Max Plus Size",
	Prompt:       `You are given an array a1,a2,…,an of positive integers.

You can color some elements of the array red, but there cannot be two adjacent red elements (i.e., for 1≤i≤n−1, at least one of ai and ai+1 must not be red).

Your score is the maximum value of a red element plus the number of red elements. Find the maximum score you can get.

Input
You are given an array arr of length n (1≤n≤100). The array contains integers a1,a2,…,an (1≤ai≤1000).

Output
For each test case, return a single integer: the maximum possible score you can get after coloring some elements red according to the statement.
Problem Source:https://codeforces.com/contest/2019/problem/A`,
	numParams:    1,
	paramTypes:   []string{"list int"},
	methodName:   "MaxPlusSize",
	NumCases:     20,
	VisibleCases: 4,
	cases: [][][]string{{{"5","4","5"}},{{"4","5","4"}},{{"3","3","3","3","4","1","2","3","4","5"}},{{"17","89","92","42","29","92","14","70","45"}},{{"7","11","8"}},{{"6"}},{{"10","6","4","7"}},{{"5","12"}},{{"1","1"}},{{"3","7","1"}},{{"5","7","3","2"}},{{"2","1","5","4"}},{{"2","1","2","1","1"}},{{"2"}},{{"7","5","6","8"}},{{"7","2"}},{{"2","2","2","2","1"}},{{"1","1","1","1","1"}},{{"8","6","6","1","2"}},{{"1"}}},
	expectedResults: []string{"7","6","10","97","12","7","12","13","2","8","9","7","5","3","10","8","5","4","11","2"},
	returnType:      "int",
	Templates: LanguageFunctionTemplates{
		Python:     `def MaxPlusSize(arr:list):\n\n\n`,
		Javascript: `function MaxPlusSize(arr){\n\n\n}`,
		Cpp:        `int MaxPlusSize(vector<int> arr){\n\n\n}`,
	},
}