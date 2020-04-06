# Check New Line in Imports

This tool is to check whether there're new line in imports.

## How To

## Install

* Clone repo, and then `make build`, then copy `bin/check-newline-in-imports` to `$GOPATH/bin` directory

## Use

* Normal use use max new line = 1
  `check-newline-in-imports -f=/home/user/source_code.go`

* Custom max new line
  `check-newline-in-imports -n=2 -f=/home/user/source_code.go`

* In common golang source tree
  `check-newline-in-imports -f="$( go list ./... | grep -v vendor | grep -P '^.+\.go$' )"`

If there're more new lines than expected, a message will be printed to `STDERR` and program will exit with code 1. Otherwise, nothing will be printed to screen and program will exit with code 0.

## Why
This tool exists because `goimports` or `goreturns` will add more and more new lines in imports when autoformatting golang source code in the case when the IDE (integrated development environment) being used auto add import statement to the source code.

Example case when `goimports` or `goreturns` will mess with import new lines:

Let's say your imports are like this at first:
```go
import (
    "error"
    "fmt"
    
    "github.com/a/b"
    "github.com/c/d"
)
```

then during coding, your IDE auto add an import statement `github.com/someone/somepackage`:
```go
import (
    "error"
    "fmt"
    "github.com/someone/somepackage"
    
    "github.com/a/b"
    "github.com/c/d"
)
```

now, when the file is saved `goimports` or `goreturns` will be executed by the IDE, and somehow this happens to your imports:
```go
import (
    "error"
    "fmt"
    
    "github.com/someone/somepackage"
    
    "github.com/a/b"
    "github.com/c/d"
)
```

Now imagine when the IDE auto add another import statement again, it will add it in the first import group, and then another auto formatting happens by `goimports` or `goreturns` will yield another new line!!

## Where

This tool is intended to be used during pre-commit to check `.go` source codes that have been cooked by `goimports` or `goreturns`.

## Notes

* This tool doesn't use [AST](https://golang.org/pkg/go/ast/), just a simple regex.
* This tool only cares with imports statement that are written like:
  ```go
  import (
  
      // bla bla bla package
  
  )
  ```
  not like:
  ```go
  import "bla1"
  import "bla2"
  // another import blas
  ```
* If this problem happens to you, and you develop a better tool that even better use AST, tell me to use that instead.