package main

import "fmt"

func Run() {
	fmt.Println("\n#1 Counting til 10")
	for i := 1; i <= 10; i++ {
		fmt.Printf("%d ", i)
	}

	fmt.Println("\n\n#2 Counting til first even number over 5")
	var i = 0
	shouldStop := false

	for !shouldStop {
		fmt.Printf("%d ", i)
		shouldStop = i > 5 && i%2 == 0
		i++
	}

	fmt.Println("\n\n#3 Array creation")
	array := make([]int, 0)

	fmt.Printf("\nArray created: %v", array)

	for i := 0; i < 10; i++ {
		array = append(array, i)
	}
	fmt.Printf("Array filled: %v\n\n", array)

	fmt.Println("#4 Map creation")

	mapa := make(map[string]interface{}, 0)
	fmt.Printf("\nMap created: %v", mapa)
	mapa["0"] = "Saudade Burra"
	mapa["1"] = "Kamasutra"
	mapa["2"] = "Sombra desconhecida"
	mapa["3"] = false
	mapa["4"] = 1000
	fmt.Println("\nMap filled")
	for i, v := range mapa {
		fmt.Printf("%s: %v\n", i, v)
	}

	fmt.Println("\n#5 Array Creation 2")
	array2 := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("\nArray created: %v", array2)
	for i := range array2 {
		array2[i] = i * i
	}

	fmt.Printf("\nArray filled: %v", array2)

	fmt.Println("#6 Map creation")

	mapa2 := map[int]interface{}{
		0: "Né possy vi",
	}
	fmt.Printf("\n\nmap created %v\n\n", mapa2)
}
