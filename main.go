package main

import (
  "fmt"; "os"; "strings"
  "github.com/nsf/termbox-go"
  "path/filepath"; "strconv"
  "syscall"
)

type Cursor struct {
  matches []string
  x int; text string
  select_mode bool
  select_index int
}

func input_mgr(cursor *Cursor, event termbox.Event) {
  _, height := termbox.Size()

  switch event.Type {
  case termbox.EventKey:
    if event.Key == termbox.KeyEsc {
      termbox.Close()
      os.Exit(0)
    }

    if !cursor.select_mode {
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
    } else {
      if event.Key == termbox.KeyArrowDown {
        if cursor.select_index < height - 5 {
          cursor.select_index ++
        }
      }

      if event.Key == termbox.KeyArrowUp {
        if cursor.select_index > 0 {
          cursor.select_index --
        }
      }

      if event.Key == termbox.KeyEnter {
        termbox.Close()

        syscall.Chdir(cursor.matches[cursor.select_index])
        os.Exit(0)
      }
    }

    if event.Key == termbox.KeyTab {
      if cursor.select_mode == false {
        cursor.select_mode = true
      } else {
        cursor.select_mode = false
        cursor.text = ""
      }
    }
  }
}

func draw_matches(cursor *Cursor, line_nums bool) {
  var offset int = 2

  for i, str := range cursor.matches {
    if line_nums {
      rownum := strconv.Itoa(i)
      if i >= 10 {
        offset = 1
      }

      for x, char := range rownum {
        termbox.SetCell(x+offset, i+1, char, termbox.ColorYellow, termbox.ColorDefault)
      }
    }
    
    for j, char := range str {
      termbox.SetCell(j+4, i+1, char, termbox.ColorDefault, termbox.ColorDefault)
    }
    
    // Draw the matches to the screen
    if cursor.select_mode && cursor.select_index == i {
      for j, char := range str {
        termbox.SetCell(j+4, i+1, char, termbox.ColorBlack, termbox.ColorWhite)
      }
    }
    
  }
}

// Draw the Search field
func draw_text(cursor *Cursor) {
  width, height := termbox.Size();

  cwd, err := os.Getwd()
  if err != nil {
    return 
  }
  
  for i := 0; i < width; i++ {
    termbox.SetCell(i, height-2, ' ', termbox.ColorDefault, termbox.ColorWhite)
  }

  for i, char := range cwd {
    termbox.SetCell(i+1, height-2, char,termbox.ColorBlack | termbox.AttrBold, termbox.ColorWhite)
  }

  if cursor.select_mode {
    return
  }

  termbox.SetCell(1, height-1, '_', termbox.ColorDefault | termbox.AttrBold, termbox.ColorDefault); 
  
  // Draw underline
  /*for i := 0; i < 100; i++ {
    termbox.SetCell(i, 1, '━', termbox.ColorDefault, termbox.ColorDefault | termbox.AttrUnderline)
  }*/

  // Draw the text from the cursor struct
  for i, char := range cursor.text {
    termbox.SetCell(i+1, height-1, char, termbox.ColorDefault, termbox.ColorDefault);
    termbox.SetCell(i+2, height-1, '_', termbox.ColorDefault, termbox.ColorDefault);
  }
}

// Find subdirs with name matching input base path
func find_subdirs(entrypoints []string, query string, depth int) []string {
	var results []string
  _, height := termbox.Size()

	for _, root := range entrypoints {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
      if len(results) > height - 5 {
        return filepath.SkipDir
      }

      if err != nil {
				return err
			}

			// Calculate the current depth level by counting path separators.
			// We add 1 to account for the root directory itself.
			currentDepth := strings.Count(path[len(root):], string(os.PathSeparator)) + 1

			if info.IsDir() {
				if (depth == -1 && currentDepth >= 1) || currentDepth <= depth {
					if strings.Contains(filepath.Base(path), query) {
						results = append(results, path) 
          }
				} else {
					return filepath.SkipDir
				}
			}

			return nil
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
    select_mode: false,
  }

  entry_points := []string { "/home" }

  // Main loop
  for {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

    cursor.matches = find_subdirs(entry_points, cursor.text, 6)
    
    draw_text(cursor)
    draw_matches(cursor, false)
    
    termbox.Flush()
  
    event := termbox.PollEvent()
    input_mgr(cursor, event)
  }
}

