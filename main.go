package main

import (
  "fmt"; "os"
  "github.com/nsf/termbox-go"
  "where/utils"
)

func main() {
  // Termbox init
  err := termbox.Init()
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  defer termbox.Close()

  // Global definitions
  glob := &utils.Global {
    X: 3, 
    Text: "",
    Select_mode: false,
  }

  entry_points := []string { "/home" }
  
  // Last configurations
  termbox.SetOutputMode(termbox.Output256)

  // Main loop
  for {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

    // Finding matching subdirs
    if !glob.Select_mode {
      glob.Matches = utils.Find_subdirs(glob, entry_points, glob.Text, 5)
    }
    
    // Render function
    utils.Matches_render(glob, false)
    utils.Text_render(glob)
    
    termbox.Flush()
  
    // Input management
    event := termbox.PollEvent()
    utils.Input_mgr(glob, event)
  }
}

