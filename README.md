# go-chat
TCP chat server in Go

```go
messages := make(chan string)  // create the pipe

// Goroutine A:
messages <- "hello"   // put "hello" in the pipe (blocks until someone takes it)

// Goroutine B:
msg := <-messages     // take "hello" out of the pipe (blocks until something arrives)
```