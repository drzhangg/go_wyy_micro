package demo

import "fmt"

type Options struct {
	Addrs []string
}

func main() {
	var arr []string = []string{"1","2"}
	fmt.Println(arr)


	var ma = make(map[string][]string)

	ma["one"] = []string{"1","2"}

	mm := map[string][]string{
		"one":[]string{"12","23"},
	}

	fmt.Println(ma)
	fmt.Println(mm)
}
