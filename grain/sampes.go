package main

import (
	"fmt"
	"strings"
)

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

	//Slice literals
	q := []int{2, 3, 5, 7, 11, 13}
	fmt.Println(q)

	r := []bool{true, false, true, true, false, true}
	fmt.Println(r)

	s := []struct {
		i int
		b bool
	}{
		{2, true},
		{5, false},
	}
	fmt.Println(s)

	//Slice defaults
	s2 := []int{2, 3, 5, 7, 11, 13}

	s2 = s2[1:4]
	s2 = s2[:2]
	s2 = s2[1:]
	fmt.Println(s2[0:4], s2[:], s2[1:], s2[:2])

	//Slice length and capacity
	s3 := []int{2, 3, 5, 7, 11, 13}
	printSlice(s3)

	// Slice the slice to give it zero length.
	s3 = s3[:0]
	printSlice(s3)

	// Extend its length.
	s3 = s3[:4]
	printSlice(s3)

	// Drop its first two values.
	s3 = s3[0:6]
	printSlice(s3)

	s3 = s3[:]
	printSlice(s3)

	//Nil slices
	//s4 := []int{}
	var s4 []int
	fmt.Println(s4, len(s4), cap(s4))
	if s4 == nil {
		fmt.Println("nil!")
	}

	//Creating a slice with make
	aaa := make([]int, 5)
	printSlice2("a", aaa)

	bbb := make([]int, 0, 5)
	printSlice2("b", bbb)

	c := bbb[:2]
	printSlice2("c", c)

	d := c[2:5]
	printSlice2("d", d)

	//Slices of slices
	//Slices can contain any type, including other slices.
	bcc := [][]string{
		[]string{"a", "c"},
	}
	fmt.Println(bcc[0][0], bcc[0][1])

	// Create a tic-tac-toe board.
	board := [][]string{
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
	}

	// The players take turns.
	board[0][0] = "X"
	board[2][2] = "O"
	board[1][2] = "X"
	board[1][0] = "O"
	board[0][2] = "X"

	for i := 0; i < len(board); i++ {
		fmt.Printf("%s\n", strings.Join(board[i], " "))
	}

	//Appending to a slice- append new elements to a slice
	var sbc []int
	sbc = append(sbc, 1)
	sbc = append(sbc, 2, 3, 4)
	printSlice(sbc)

	//Range
	//When ranging over a slice, two values are returned for each iteration. The first is the index,
	// and the second is a copy of the element at that index.
	pow := []int{1, 2, 4, 8, 16, 32}
	for i, v := range pow {
		fmt.Printf("2**%d = %d\n", i, v)
	}

	//Maps
	var m map[string]Wertex = make(map[string]Wertex)
	m["a1"] = Wertex{1.2, 1.4}
	fmt.Println(m)
	fmt.Printf("%v \n", m)

	var n map[int]string = make(map[int]string)
	n[0] = "test"
	fmt.Println(n)

	//Map literals
	var mm = map[string]Wertex{
		"Bell Labs": {
			40.68433, -74.39967,
		},
		"Google": {
			37.42202, -122.08408,
		},
	}
	fmt.Println(mm)
}

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func printSlice2(s string, x []int) {
	fmt.Printf("%s len=%d cap=%d %v\n",
		s, len(x), cap(x), x)
}

type Vertex struct {
	X, Y int
}

type Wertex struct {
	Long, Lat float64
}

//
