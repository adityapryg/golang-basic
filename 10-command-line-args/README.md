# Command-Line Arguments in Go

Command-line arguments in Go are the arguments that you pass to the program when you execute it. In Go, these arguments can be accessed using the `os.Args` slice, which contains all the arguments including the program name.

## Using os.Args

The `os.Args` variable is a slice of strings that holds the command-line arguments passed to the program. Here's an example:

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // os.Args[0] is the program name
    // os.Args[1:] holds the arguments passed
    fmt.Println("Program Name: ", os.Args[0])
    fmt.Println("Arguments: ", os.Args[1:])
}
```

### Running the Program

To run the above program and pass arguments, use:
```
go run main.go arg1 arg2
```

## Using the flag Package

The `flag` package provides a more sophisticated way to handle command-line arguments. It allows you to define flags, parse them, and access their values:

```go
package main

import (
    "flag"
    "fmt"
)

func main() {
    // Define a string flag with a default value
    name := flag.String("name", "World", "a name to say hello to")

    // Parse the flags from command-line
    flag.Parse()

    fmt.Printf("Hello, %s!
", *name)
}
```

### Running the Program with Flags

To run the program with the defined flag:
```
go run main.go -name=Alice
```

This will output:
```
Hello, Alice!
```

## Conclusion

Command-line arguments are a vital part of any command-line application, and Go provides simple yet powerful tools to handle them effectively.