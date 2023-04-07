package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Printf("result: %v", sum(3, 4))
}

func sum(a int, b int) int {
	fmt.Println("Add percent of test coverage")
	fmt.Println("Add percent of test coverageB")
	fmt.Println("Add percent of test coverageC")
	left := []string{"leftA", "leftB", "leftC"}
	right := []string{"rightA", "rightB", "rightC"}

	//Test gosec detection
	name := fmt.Sprintf("%s_%s", left[rand.Intn(len(left))], right[rand.Intn(len(right))])
	fmt.Printf("name: %v", name)

	//Test gosec detection
	name2 := fmt.Sprintf("%s_%s", left[rand.Intn(len(left))], right[rand.Intn(len(right))])
	fmt.Printf("name: %v", name2)
	return a + b
}
