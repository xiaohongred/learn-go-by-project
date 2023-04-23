package main

import "fmt"

func f() (result int) {
	defer func() {
		result++
	}()
	return 0
}

func f2() (r int) {
	t := 5
	defer func() {
		t = t + 5
	}()
	return t
}

func f1(arr []int) {
	arr = append(arr, []int{1, 2, 3}...)
}

func f3() (r int) {
	defer func(r int) {
		r = r + 5
	}(r)
	return 1
}

func example1b() {
	sl := []string{"a", "b", "c", "d"}
	fmt.Println(sl)

	func(slParam []string) {
		slParam = append(slParam, "er")
		slParam = append(slParam, "luohongsheng")
		fmt.Println(slParam)
		fmt.Println(len(slParam), cap(slParam))
	}(sl)

	fmt.Println(sl)
	fmt.Println(len(sl), cap(sl))
}

func main() {
	example1b()
}
