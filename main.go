package main

import (
  "fmt"; "os"
  "github.com/nsf/termbox-go"
)

type Cursor struct {
  matches []string
  x int; text string
}

func input_mgr(cursor *Cursor, event termbox.Event) {
  switch event.Type {
  case termbox.EventKey:
    if event.Key == termbox.KeyEsc {
      termbox.Close()
      os.Exit(0)
    }

    if event.Ch != 0 {
      cursor.text += string(event.Ch)
    }

    if event.Key == termbox.KeySpace {
      cursor.text += " "
    }
  
    if event.Key == termbox.KeyBackspace2 {
      if len(cursor.text) > 0 {
        cursor.text = cursor.text[:len(cursor.text) - 1]
      }
    }
  }
}


// Draw the Search field
func draw_text(cursor Cursor) {
  termbox.SetCell(2, 1, '>', termbox.ColorWhite, termbox.ColorDefault)
  termbox.SetCell(4, 1, '_', termbox.ColorDefault, termbox.ColorDefault);
  
  // Draw the text from the cursor struct
  for i, char := range cursor.text {
    termbox.SetCell(i+4, 1, char, termbox.ColorDefault, termbox.ColorDefault);
    termbox.SetCell(i+5, 1, '_', termbox.ColorDefault, termbox.ColorDefault);
  }
}

func main() {
  err := termbox.Init()
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  defer termbox.Close()

  cursor := &Cursor {
    x: 3, 
    text: "",
  }

  // Main loop
  for {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
    draw_text(*cursor)

    termbox.Flush()
  
    event := termbox.PollEvent()
    input_mgr(cursor, event)
  }
}

