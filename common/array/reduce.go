package array

func Reduce[T1, T2 any](arr []T1, fn func(T1, int) []T2) (r []T2) {
	for i, v := range arr {
		r = append(r, fn(v, i)...)
	}
	return
}
