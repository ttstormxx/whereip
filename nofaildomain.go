package main

func RemoveFailDomain(misc []string, failint []int) []string {
	var filter_faildomain []string
	var domain_valid bool

	for i := 0; i < len(misc); i++ {
		domain_valid = true
		// fmt.Println("循环的i是", i)
		for _, value := range failint {
			if i == value {
				domain_valid = false
				break
			}
		}
		if domain_valid {
			// fmt.Println("这个域名是有效的", misc[i])
			filter_faildomain = append(filter_faildomain, misc[i])
		}
	}

	return filter_faildomain
}
