package collection

// CartesianProduct → {(x, y) | ∀ x ∈ s1, ∀ y ∈ s2}
// For example:
//		CartesianProduct(A, B), where A = {1, 2} and B = {7, 8}
//        => {(1, 7), (1, 8), (2, 7), (2, 8)}
func CartesianProduct(sets ...[]interface{}) [][]interface{} {
	if len(sets) == 0 {
		return nil
	}
	lens := func(i int) int { return len(sets[i]) }
	var product [][]interface{}
	for ix := make([]int, len(sets)); ix[0] < lens(0); nextIndex(ix, lens) {
		var r []interface{}
		for j, k := range ix {
			r = append(r, sets[j][k])
		}
		product = append(product, r)
	}
	return product
}

func nextIndex(ix []int, lens func(i int) int) {
	for j := len(ix) - 1; j >= 0; j-- {
		ix[j]++
		if j == 0 || ix[j] < lens(j) {
			return
		}
		ix[j] = 0
	}
}
