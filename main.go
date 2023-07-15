package main

import (
  "fmt"; "os"; "strings"
  "github.com/nsf/termbox-go"
  "path/filepath"; "strconv"
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

func draw_matches(cursor *Cursor) {
  for i, str := range cursor.matches {
    for j, char := range str {
      termbox.SetCell(j+4, i+3, char, termbox.ColorDefault, termbox.ColorDefault)
    }
    
    rownum := strconv.Itoa(i)
    
    for x, char := range rownum {
      termbox.SetCell(x+1, i+3, char, termbox.ColorYellow, termbox.ColorDefault)
    }
  } 
}


// Draw the Search field
func draw_text(cursor Cursor) {
  termbox.SetCell(4, 1, '_', termbox.ColorDefault, termbox.ColorDefault); 

  // Draw the text from the cursor struct
  for i, char := range cursor.text {
    termbox.SetCell(i+4, 1, char, termbox.ColorDefault, termbox.ColorDefault);
    termbox.SetCell(i+5, 1, '_', termbox.ColorDefault, termbox.ColorDefault);
  }
}

// Find subdirs with name matching input base path
func find_subdirs(entrypoints []string, query string, depth int) []string {
  var results []string

  for _, root := range entrypoints {
    filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
      if err != nil {
        return err
      }

      if info.IsDir() {
        if strings.Contains(filepath.Base(path), filepath.Base(query)) {
          results = append(results, path)
        }
      }

      if depth > 0 || depth == -1 {
        if depth != -1 {
          depth--
        }
        return nil
      }

      return filepath.SkipDir
    })
  }

  return results
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

  entrypoints := []string { "/home/xyandzaxis/projects", "/home/xyandzaxis/.config" }

  // Main loop
  for {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

    cursor.matches = find_subdirs(entrypoints, cursor.text, 100)
    draw_text(*cursor)
    
    draw_matches(cursor)

    termbox.Flush()
  
    event := termbox.PollEvent()
    input_mgr(cursor, event)
  }
}

