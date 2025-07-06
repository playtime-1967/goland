package main

import "fmt"

func RunSamples() {
	//Pointers
	i, j := 1, 1
	p := &i
	fmt.Println(*p)
	*p = 2
	fmt.Println(i)
	p = &j // point to j
	*p = 2
	fmt.Println(j) // see the new value of j

	//structs
	v := Vertex{1, 4}
	v.X = 3
	fmt.Println(v.X, v.Y)

	// Pointers to structs
	v2 := Vertex{1, 4}
	p2 := &v2
	p2.X = 3 //OR (*p2).X = 3
	fmt.Println(*p2)

	//Struct Literals
	var (
		v1 = Vertex{1, 2}  // has type Vertex
		v4 = Vertex{X: 1}  // Y:0 is implicit
		v5 = Vertex{}      // X:0 and Y:0
		p3 = &Vertex{1, 2} // has type *Vertex
	)
	fmt.Println(v1, v4, v5, p3)

	//arrays
	var a [2]string
	a[0] = "Hello"
	a[1] = "World"
	fmt.Println(a[0], a[1])
	fmt.Println(a)

	primes := [6]int{2, 3, 5, 7, 11, 13}
	fmt.Println(primes)

	//Slices

	slicer := primes[0:3]
	fmt.Println(slicer)

	//Slices are like references to arrays
	names := [4]string{
		"John",
		"Paul",
		"George",
		"Ringo",
	}
	fmt.Println(names)

	aa := names[0:2]
	b := names[1:3]

	b[0] = "XXX"
	fmt.Println(aa)
	fmt.Println(b)
	fmt.Println(names)
}

type Vertex struct {
	X, Y int
}
