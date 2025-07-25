package main

import (
	"fmt"
	"image"
	"io"
	"math"
	"strings"
	"sync"
	"time"
)

func RunGoTour() {
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

	//Mutating Maps
	m2 := make(map[string]int)

	m2["Answer"] = 1
	m2["Answer"] = 2
	fmt.Println(m2)

	delete(m2, "Answer")
	m2["Answer"] = 3

	val, ok := m2["Answer"]
	fmt.Println("The value:", val, "Present?", ok)

	//Function values
	//Functions are values too. They can be passed around just like other values.
	fn := func(x, y float64) float64 {
		return math.Sqrt(x*x + y*y)
	}
	fmt.Println(fn(5, 12))
	fmt.Println(compute(fn))
	fmt.Println(compute(math.Pow))

	//Function closures
	pos, neg := adder(), adder()
	for i := 0; i < 3; i++ {
		fmt.Println(
			pos(i),
			neg(-2*i),
		)
	}

	//Methods
	//define methods on types
	//A method is a function with a special receiver argument.
	//You cannot declare a method with a receiver whose type is defined in another package
	vm := Wertex{1.1, 2.2}
	fmt.Println(vm.Abs())

	//Pointer receivers *T (T cannot itself be a pointer)
	vm.Scale(10)
	fmt.Println(vm)
	vm2 := &vm
	vm2.Scale(10) //Works: methods with pointer receivers take either a value or a pointer as the receiver when they are called:
	fmt.Println(vm2)

	//Pointers and functions
	Scale(&vm, 20)
	//Scale(vm, 20) Compile error!
	fmt.Println(vm)

	//Interfaces
	var abser Abser
	myf := MyFloat(math.Sqrt2)
	abser = myf

	wer := &Wertex{1, 2}
	abser = wer
	abser.Abs2()

	wer2 := Wertex{1, 2}
	//abser = wer2 Compile error!
	wer2.Abs2()

	//Interfaces are implemented implicitly
	var i2 I = &TT{S: "something"}
	i2.M()
	describe(i2)

	//Interface values with nil underlying values
	var i3 I
	var t3 *TT
	i3 = t3
	describe(i3)
	i3.M() //Nil

	//Nil interface value
	var i4 I
	describe(i4)
	//i4.M() // Compile error!: which concrete method to call

	//The empty interface :zero methods
	//An empty interface may hold values of any type.
	//For example, fmt.Print takes any number of arguments of type interface{}.
	var empty Empty
	describe(empty)
	empty = 10
	describe(empty)
	empty = "text"
	describe(empty)

	//Type assertions
	//A type assertion provides access to an interface value's underlying concrete value.
	var it interface{} = "Hello"
	st := it.(string)
	fmt.Println(st)

	//To test whether an interface value holds a specific type, a type assertion can return two values
	st2, ok2 := it.(string)
	fmt.Println(st2, ok2)

	st3, ok3 := it.(float64) //panic
	fmt.Println(st3, ok3)

	//Type switches
	do(110)
	do("Hi")
	do(true)

	//Stringers
	pp := &person{"Woolf", 45}
	fmt.Println(pp)

	//
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}

	//Errors
	if err := RunErrorTest(); err != nil {
		fmt.Println(err)
	}

	sq, err := Sqrt(10)
	fmt.Println(sq, err)

	sq, err = Sqrt(-10)
	fmt.Println(sq, err)
	if x, err := Sqrt(-10); err != nil {
		fmt.Println(err, x)
	}

	//Readers
	redr := strings.NewReader("Hello, Reader!")

	byt := make([]byte, 8)
	for {
		n, err := redr.Read(byt)
		fmt.Printf("n = %v err = %v b = %v\n", n, err, byt)
		fmt.Printf("b[:n] = %q\n", byt[:n])
		if err == io.EOF {
			break
		}
	}

	//infinite stream
	// myr := MyReader{}
	// __, _ := myr.Read(make([]byte, 1))
	// println(__)

	//Images
	mi := image.NewRGBA(image.Rect(0, 0, 100, 100))
	fmt.Println(mi.Bounds())
	fmt.Println(mi.At(0, 0).RGBA())

	//Generics
	//1- functions& parameters
	ig := []int{1, 10, 15, 20}
	println(Index(ig, 15))
	sg := []string{"foo", "bar", "zar"}
	println(Index(sg, "for"))

	//Goroutines
	//a lightweight thread managed by the Go runtime.

	go say("world")
	say("hello")

	//Channels
	//By default, sends and receives block until the other side is ready. This allows goroutines to
	// synchronize without explicit locks or condition variables.

	sc := []int{7, 2, 8, -9, 4, 0}
	cch := make(chan int)
	go sum(sc[:len(sc)/2], cch)
	go sum(sc[len(sc)/2:], cch)
	x, y := <-cch, <-cch // receive from cch
	fmt.Println(x, y, x+y)

	//1- Use a goroutine for sending or receiving
	ch := make(chan int)
	go func() {
		ch <- 42 //send to the ch in a separate thread of execution

	}()
	vv := <-ch //receive from the ch
	fmt.Println(vv)

	//2-Use a buffered channel
	ch2 := make(chan int, 1)
	ch2 <- 42
	vv2 := <-ch2
	fmt.Println(vv2)

	//Buffered Channels
	//Sends to a buffered channel block only when the buffer is full. Receives block when the buffer is empty.
	ch3 := make(chan int, 2)
	ch3 <- 7
	ch3 <- 8
	println(<-ch3)
	println(<-ch3)

	//Close channels
	//Range and Close
	println("------------------")
	ch4 := make(chan int, 10)
	go fibo(cap(ch4), ch4)
	//Closing is only necessary when the receiver must be told there are no more values coming,
	// such as to terminate a range loop.
	for i := range ch4 {
		println(i)
	}
	println("------------------")
	//Select
	//lets a goroutine wait on multiple communication operations.
	//A select blocks until one of its cases can run, then it executes that case
	ch5 := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println(<-ch5)
		}
		quit <- 0
	}()
	fibonacci(ch5, quit)

	//Default Selection
	selectFunc()
	println("------------------")
	//mutual exclusion
	counter := SafeCounter{v: make(map[string]int)}
	for i := 0; i < 10; i++ {
		go counter.Inc("mykey")
	}
	time.Sleep(time.Millisecond)
	println(counter.Value("mykey"))
}

type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	//only one goroutine at a time can access the map
	c.v[key]++
	c.mu.Unlock()
}
func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock() //mutex will be unlocked at the end
	return c.v[key]
}

// ---------------
func selectFunc() {
	start := time.Now()
	tick := time.Tick(10 * time.Millisecond)
	boom := time.After(25 * time.Millisecond)
	elapsed := func() time.Duration {
		return time.Since(start).Round(time.Millisecond)
	}
	for {
		select {
		case <-tick:
			fmt.Printf("[%6s] tick.\n", elapsed())
		case <-boom:
			fmt.Printf("[%6s] BOOM!\n", elapsed())
			return
		default:
			fmt.Printf("[%6s]     .\n", elapsed())
			time.Sleep(5 * time.Millisecond)
		}
	}
}
func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			println("quit")
			return
		}
	}
}

// ----------
func fibo(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c) //if we don't close, the channel waiting forever to receive a new val
}

// /-------
func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}

func say(s string) {
	for i := 0; i < 3; i++ {
		time.Sleep(100 * time.Millisecond)
		println(s)
	}
}

type List[T any] struct {
	next  *List[T]
	value T
}

// ------
func Index[T comparable](s []T, x T) int {

	for i, v := range s {
		if v == x {
			return i
		}
	}
	return -1
}

// ------ Custom Reader
type MyReader struct{}

func (r MyReader) Read(b []byte) (int, error) {
	for {
		n, _ := strings.NewReader("A").Read(b)
		fmt.Printf("= %q\n", b[:n])
	}
}

// ------ Custom Error
type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}
	return x, nil
}

// ---
type MyError struct {
	When time.Time
	What string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("at %v, %s", e.When, e.What)
}

func RunErrorTest() error {
	return &MyError{time.Now(), "it didn't work"}
}

// -------------
type IPAddr [4]byte

func (addr IPAddr) String() string {
	return fmt.Sprintf("%v.%v.%v.%v", addr[0], addr[1], addr[2], addr[3])
}

// --------------
type person struct {
	Name string
	Age  int
}

func (p *person) String() string {
	return fmt.Sprintf("%v (%v years)", p.Name, p.Age)
}

// ----------------------------
func do(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Printf("twice %v is %v \n", v, v*2)
	case string:
		fmt.Printf("%q is %v  bytes long\n", v, len(v))
	default:
		fmt.Printf("I don't know about type %T\n", v)
	}
}

type Empty interface{}

// ---------------------
type I interface {
	M()
}
type TT struct {
	S string
}

func (t *TT) M() {
	if t == nil {
		fmt.Println("<nil>")
		return
	}
	fmt.Println(t.S)
}
func describe(i Empty) {
	fmt.Printf("(%v, %T) \n", i, i)
}

// ---------
type Abser interface {
	Abs2() float64
}
type MyFloat float64

func (f MyFloat) Abs2() float64 {
	if f < 0 {
		return float64(-f)
	}

	return float64(f)
}

func (w *Wertex) Abs2() float64 {
	return math.Sqrt(w.Lat*w.Lat + w.Long*w.Long)
}

///-------

func (v Wertex) Abs() float64 {
	return math.Sqrt(v.Long*v.Long + v.Lat*v.Lat)
}

// method
func (w *Wertex) Scale(f float64) {
	w.Lat = w.Lat * f
	w.Long = w.Long * f
}

// func
func Scale(w *Wertex, f float64) {
	w.Lat = w.Lat * f
	w.Long = w.Long * f
}

// ----------
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

func compute(myfunc func(float64, float64) float64) float64 {
	return myfunc(3, 4)
}

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

//
