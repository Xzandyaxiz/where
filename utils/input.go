package utils

import (
  "os"; "github.com/nsf/termbox-go"
)

func Input_mgr(glob *Global, event termbox.Event) {
  switch event.Type {
  case termbox.EventKey:
    if event.Key == termbox.KeyEsc {
      termbox.Close()
      os.Exit(0)
    }

    if !glob.Select_mode {
      if event.Ch != 0 {
        glob.Text += string(event.Ch)
      }

      if event.Key == termbox.KeySpace {
        glob.Text += " "
      }
  
      if event.Key == termbox.KeyBackspace2 {
        if len(glob.Text) > 0 {
          glob.Text = glob.Text[:len(glob.Text) - 1]
        }
      }
    } else {
      if event.Key == termbox.KeyArrowDown {
        if glob.Select_index < len(glob.Matches) - 1 {
          glob.Select_index ++
        }
      }

      if event.Key == termbox.KeyArrowUp {
        if glob.Select_index > 0 {
          glob.Select_index --
        }
      }

      if event.Key == termbox.KeyEnter {
        // nothing here yet...
      }
    }

    if event.Key == termbox.KeyTab {
      if glob.Select_mode == false {
        glob.Select_mode = true
      } else {
        glob.Select_mode = false
        glob.Select_index = 0
        glob.Text = ""
      }
    }
  }
}
