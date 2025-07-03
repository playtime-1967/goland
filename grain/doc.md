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