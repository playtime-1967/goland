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

