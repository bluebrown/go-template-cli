package textfunc

func Iter(n int) []int {
	l := make([]int, n)
	for i := 0; i < n; i++ {
		l[i] = i
	}
	return l
}
