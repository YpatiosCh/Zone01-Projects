package piscine

func Unmatch(arr []int) int {
	for _, res := range arr {
		pair := 0
		for _, el := range arr {
			if el == res {
				pair++
			}
		}
		if pair == 1 || pair%2 == 1 {
			return res
		}
	}
	return -1
}
