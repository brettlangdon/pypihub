package pypihub

func uniqueSlice(s []string) []string {
	var m = make(map[string]bool)
	for _, v := range s {
		m[v] = true
	}

	var o = make([]string, 0)
	for v := range m {
		o = append(o, v)
	}
	return o
}

func removeEmpty(s []string) []string {
	var n = make([]string, 0)
	for _, c := range s {
		if c != "" {
			n = append(n, c)
		}
	}
	return n
}
