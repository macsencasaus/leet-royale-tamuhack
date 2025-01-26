package testrunner

var Q8 = QuestionData{
	Title:        "Frog 1",
	Prompt:       `There are n stones, numbered 1,2,...,N(2≤N≤10^5). For each i(1≤i≤N), the height of Stone i is hi(1≤hi≤10^4). These heights are stored in arr.
There is a frog who is initially on Stone 1. He will repeat the following action some number of times to reach Stone N:
If the frog is currently on Stone i, jump to Stone i+1 or Stone i+2. Here, a cost of abs(hi​ - hj​) is incurred, where j is the stone to land on.
Return the minimum possible total cost incurred before the frog reaches Stone N.
Problem Source: https://atcoder.jp/contests/dp/tasks/dp_a`,
	numParams:    2,
	paramTypes:   []string{"int","list int"},
	methodName:   "Frog1",
	NumCases:     7,
	VisibleCases: 3,
	cases: [][][]string{{{"4"},{"10","30","40","20"}},{{"2"},{"10","10"}},{{"6"},{"30","10","60","10","60","50"}},{{"2"},{"1","1"}},{{"10"},{"1000","900","700","1000","855","388","535","999","434","987"}},{{"100"},{"691","157","802","416","988","516","973","930","903","697","280","475","685","951","338","637","838","677","188","471","39","18","726","135","301","693","36","478","708","704","251","996","640","896","660","333","92","820","52","417","201","335","214","132","990","329","933","889","587","179","696","624","208","735","752","954","576","27","619","461","875","739","884","449","587","377","87","496","385","473","884","605","409","584","798","377","28","813","542","764","91","867","344","502","949","497","493","655","661","232","602","537","466","167","636","377","619","151","210","179"}},{{"77"},{"406","382","14","159","450","96","747","681","278","652","102","237","816","47","148","161","719","803","337","135","472","141","454","539","823","242","509","389","913","818","52","852","825","549","882","135","82","137","670","469","85","290","426","729","438","879","849","472","278","186","147","656","130","511","752","500","515","677","827","917","501","847","361","405","212","318","872","75","683","559","611","826","673","146","335","427","688"}}},
	expectedResults: []string{"30","0","40","0","1141","9682","9084"},
	returnType:      "int",
	Templates: LanguageFunctionTemplates{
		Python:     `def Frog1(n, arr:list):\n\n\n`,
		Javascript: `function Frog1(n, arr){\n\n\n}`,
		Cpp:        `int Frog1(int n, vector<int> arr){\n\n\n}`,
	},
}