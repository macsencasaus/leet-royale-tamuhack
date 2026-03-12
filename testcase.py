def sum(arr: list[int]) -> int:
	return 0

class TestCase:
	def __init__(self, arr: list[int], result: int):
		self.arr = arr
		self.result = result


def magic(thingToPrint):
	print('\n', "AAAAA", '\n', thingToPrint, '\n', "AAAAA", sep='', end='')

def main():
	cases = [TestCase(arr=[1, 2, 3], result=6), TestCase(arr=[-1, -2, 3], result=0)]
	for case in cases:
		try:
			res = sum(arr=case.arr)
		except:
			magic("RE")
			return

		if type(res) is type(case.result) and res == case.result:
			magic("AC")
		else:
			magic("WA")

main()
