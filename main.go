package main

import (
  "os"
  "time"
)

import (
  "fmt"
  "github.com/nsf/termbox-go"
)

func main() {
  err := termbox.Init()
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  termbox.Flush()

  time.Sleep(time.Second);
  termbox.Close()

}

