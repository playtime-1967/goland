Workspace:
mkdir my-workspace && cd my-workspace
go work init
go work use ./{module-name} ./{module-name}

Make a module:
go mod init {module-name}


run mod:
go run main.go


install a pkg
go get github.com/gin-gonic/gin
go get github.com/gocql/gocql

go.mod vs go.sum 
is essential for managing dependencies in Go modules
# go.mod – Dependency Manifest
This file declares:
Your module name
The direct dependencies your code requires
Their versions
# go.sum – Dependency Checksum File
This file ensures security and reproducibility.
It stores:
Cryptographic checksums for all dependencies (including transitive ones).
These are verified during go build, go get, etc., to ensure no tampering or mismatch.



# Exported names
In Go, a name is exported if it begins with a capital letter.
pizza and pi do not start with a capital letter, so they are not exported.


# gin.HandlerFunc
gin.HandlerFunc is a func(c *gin.Context) — a typical HTTP handler in Gin.
c *gin.Context is the actual Gin handler that handles the HTTP request 

if initialization; condition { }
if err:= func(); err!=nil{ }


# exception handling
In Go, error handling is explicit and always returned as a value (error).
err := repo.Create(user) if err != nil {     // handle it }

# * in Go- Avoid 
a pointer to an object- reference a memory location.
copying large structs, Share mutable state.
*gocql.Session: share a single Cassandra session across functions without copying it

# return &cassandraUserRepo{session: session}
& creates a pointer to the struct.

#  Receiver Methods
func (u *User) IsEmpty() bool {
	return u.Name == ""
}

interfaces in Go are already reference-like types


# Cassandra consistency level
consistency level controls how many replicas must acknowledge a read or write operation before it's considered successful.
gocql.One means:
Only 1 replica needs to respond for the query to succeed.
It’s fast, but less reliable if replicas are out of sync.
gocql.Quorum: a majority of replicas must respond — safer for consistency.
gocql.All: all replicas must agree — safest, but slower.

# Defer
A defer statement defers the execution of a function until the surrounding function returns.
Stacking defers: Deferred function calls are pushed onto a stack. 


# Slices
a dynamically-sized, flexible view into the elements of an array. 
includes the first element, but excludes the last one.
a[low : high]
a[1:4]
Slices are like references to arrays; A slice does not store any data, it just describes a section of an underlying array.
Changing the elements of a slice modifies the corresponding elements of its underlying array.


# Slice literals
A slice literal is like an array literal without the length.
creates an array, then builds a slice that references it:
[]bool{true, true, false}

# Slice length and capacity
The length of a slice is the number of elements it contains
The capacity of a slice is the number of elements in the underlying array, counting from the first element in the slice.
The length and capacity of a slice s can be obtained using the expressions len(s) and cap(s).
You can extend a slice's length by re-slicing it, provided it has sufficient capacity.


# fmt.Printf
%d	Decimal integer
%s	String	
%v	Default format (very useful for slices, structs, etc.)	
%T	Type of the value	fmt.Printf("%T", x) → int
%f	Decimal floating point
%t	Boolean	fmt.Printf("%t", true)
%q Double-quoted string fmt.Printf("%q", "hello") → "hello"


# Choosing a value or pointer receiver
There are two reasons to use a pointer receiver.
The first is so that the method can modify the value that its receiver points to.
The second is to avoid copying the value on each method call. This can be more efficient if the receiver is a large struct.

# Interfaces are implemented implicitly
A type implements an interface by implementing its methods. There is no explicit declaration of intent, no "implements" keyword.
Implicit interfaces decouple the definition of an interface from its implementation, which could then appear in any package without prearrangement.


# Interface values with nil underlying values
