# Defer, Panic, and Recover in Go

## Introduction
This document provides a comprehensive overview of the Defer, Panic, and Recover mechanisms in Go programming language, which are used for managing the flow of control during error handling and resource management.

## Defer
### Explanation  
The `defer` statement in Go is used to ensure that a function call is performed later in the program's execution, usually for purposes of cleanup. Defer statements are executed in Last In, First Out (LIFO) order.

### Examples  
```go
func main() {
    defer fmt.Println("World")
    fmt.Println("Hello")
}
// Output: Hello
// Output: World
```

### Use Cases  
- Releasing resources (like file handles)
- Unlocking mutexes in concurrent programming

### Pitfalls  
- Deferred calls can increase the overall execution time if used excessively.

### Best Practices  
- Use defer for cleanup tasks and make sure they’re the last thing in the scope of function.

## Panic
### Explanation  
A `panic` is a built-in function in Go that stops the ordinary flow of control and begins panicking. It is often used for unrecoverable errors.

### Examples  
```go
func main() {
    panic("something went wrong")
}
```

### Use Cases  
- Handling critical errors that cannot be recovered.

### Pitfalls  
- Overusing panic can lead to unmanageable code and crashes.

### Best Practices  
- Only use panic for truly exceptional conditions.

## Recover
### Explanation  
`recover` is a built-in function that allows you to regain control of a panicking goroutine. This can only be called within a deferred function.

### Examples  
```go
func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from", r)
        }
    }()
    panic("panic example")
}
// Output: Recovered from panic example
```

### Use Cases  
- Cleaning up state when a goroutine panics.

### Pitfalls  
- Using recover without a proper understanding may hide underlying issues.

### Best Practices  
- Always pair recover with defer for graceful error handling.

## Exercises  
1. Write a program that demonstrates the use of defer to close a file.
2. Create a function that panics, and show how to recover from it without crashing the program.
3. Experiment with nesting defer statements and observe the execution order.
