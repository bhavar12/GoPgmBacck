package main

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func main1() {
	slice1 := []int{1, 2}
	slice2 := []int{3, 4}
	slice3 := slice1
	copy(slice1, slice2)
	fmt.Println(slice1, slice2, slice3) // {3,4} , {3,4} , {3,4}

	map_obj := make(map[string]struct{})
	for _, value := range []string{"test1", "test2", "test3"} {
		map_obj[value] = struct{}{}
	}
	fmt.Println(map_obj) // {test1:{},test2:{},test:{}}

	arr := [6]string{"This", "is", "a", "hello", "world", "program"}

	fmt.Println("Original Array:", arr)

	slicedArr := arr[1:4]

	fmt.Println("Sliced Array:", slicedArr) // is a hello

	fmt.Println("Length of the slice: %d", len(slicedArr)) // 3

	fmt.Println("Capacity of the slice: %d", cap(slicedArr)) // 5

	car1, car2 := createTeslaCar(), createTeslaCar()
	fmt.Println(car1(" red"), car2(" blue"))
	fmt.Println(car1(" car"))

	g := new(errgroup.Group)
	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.somestupidname.com/",
	}
	for _, url := range urls {
		// Launch a goroutine to fetch the URL.
		url := url // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			// Fetch the URL.
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
			}
			return err
		})
	}
	// Wait for all HTTP fetches to complete.
	if err := g.Wait(); err == nil {
		fmt.Println("Successfully fetched all URLs.")
	}
}

func createTeslaCar() func(string) string {
	car := "Tesla"
	return func(s string) string {
		car += s
		return car
	}
}
