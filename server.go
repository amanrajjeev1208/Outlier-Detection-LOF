// Very basic socket server
// https://golangr.com/

package main

import "net"
import "fmt"
import "bufio"
import "os"

func main() {
  fmt.Println("Start server...")

  // listen on port 8000
  ln, _ := net.Listen("tcp", ":8000")

  // accept connection
  conn, _ := ln.Accept()

  // run loop forever (or until ctrl-c)
  for {
    // get message, output
    message, _ := bufio.NewReader(conn).ReadString('\n')
    fmt.Print("Message Received:", string(message))
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Text to reply: ")
    text, _ := reader.ReadString('\n')
    // send to server
    fmt.Fprintf(conn, text + "\n")
  }
}