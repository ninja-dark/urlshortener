package generator



func Encode(i int) string {
	dict := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	base := len(dict)
	digits := []int{}
	for i > 0 {
		r := i % base
		digits = append([]int{r}, digits...)
		i = i / base
	}

	r := []rune{}
	for _, d := range digits {
		r = append(r, dict[d])
	}
	return string(r)
}

func Decode(s string) int {
	dict := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	base := len(dict)
	d := 0 
	for _, ch := range s {
		for i, a := range dict {
			if a == ch {
				d = d*base + i
			}
		}
	}
	return d
}