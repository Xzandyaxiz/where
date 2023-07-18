package main

import (
  "fmt"; "os"; "strings"
  "github.com/nsf/termbox-go"
  "path/filepath"; "strconv"
)

const LightGray = 251

type Global struct {
  matches []string
  match_size []int64
  x int; text string
  select_mode bool
  select_index int
}

func input_mgr(glob *Global, event termbox.Event) {
  switch event.Type {
  case termbox.EventKey:
    if event.Key == termbox.KeyEsc {
      termbox.Close()
      os.Exit(0)
    }

    if !glob.select_mode {
      if event.Ch != 0 {
        glob.text += string(event.Ch)
      }

      if event.Key == termbox.KeySpace {
        glob.text += " "
      }
  
      if event.Key == termbox.KeyBackspace2 {
        if len(glob.text) > 0 {
          glob.text = glob.text[:len(glob.text) - 1]
        }
      }
    } else {
      if event.Key == termbox.KeyArrowDown {
        if glob.select_index < len(glob.matches) - 1 {
          glob.select_index ++
        }
      }

      if event.Key == termbox.KeyArrowUp {
        if glob.select_index > 0 {
          glob.select_index --
        }
      }

      if event.Key == termbox.KeyEnter {
        // nothing here yet...
      }
    }

    if event.Key == termbox.KeyTab {
      if glob.select_mode == false {
        glob.select_mode = true
      } else {
        glob.select_mode = false
        glob.select_index = 0
        glob.text = ""
      }
    }
  }
}

func draw_matches(glob *Global, line_nums bool) {
  var offset int = 2
  width, _ := termbox.Size()

  for i, str := range glob.matches {
    // Draw line nums if true
    if line_nums {
      rownum := strconv.Itoa(i)
      if i >= 10 {
        offset = 1
      }

      for x, char := range rownum {
        termbox.SetCell(x+offset, i+1, char, termbox.ColorYellow, termbox.ColorDefault)
      }
    }
    
    bytes := strconv.Itoa(int(glob.match_size[i]))
    match_size_startpos := width / 2 - len(bytes)

    // Print the match at the correct level
    for j, char := range str {
      if j+7 >= match_size_startpos {
        break
      }

      termbox.SetCell(j+5, i+1, char, termbox.ColorDefault, termbox.ColorDefault)
    }

    // Show the sizes right next to the printed matches
    for index, char := range bytes {
      termbox.SetCell(match_size_startpos + index, i+1, char, LightGray | termbox.AttrBold, termbox.ColorDefault)
    }
    
    // Draw the matches to the screen
    if glob.select_mode && glob.select_index == i {
      for j, char := range str {
        termbox.SetCell(j+5, i+1, char, termbox.ColorBlack, LightGray)
      }

      for j := len(str); j < width/2; j++ {
        termbox.SetCell(j+5, i+1, ' ', termbox.ColorDefault, LightGray)
      } 
    }
  }
}

// Draw the Search field
func draw_text(glob *Global) {
  width, height := termbox.Size();

  cwd, err := os.Getwd()
  if err != nil {
    return 
  }
  
  for i := 0; i < width; i++ {
    termbox.SetCell(i, height-2, ' ', termbox.ColorDefault, LightGray)
  }

  for i, char := range cwd {
    termbox.SetCell(i+1, height-2, char, termbox.ColorBlack | termbox.AttrBold, LightGray)
  }

  if glob.select_mode {
    return
  }
 
  termbox.SetCell(1, height-1, '_', termbox.ColorDefault | termbox.AttrBold, termbox.ColorDefault);  

  // Draw the text from the cursor struct
  for i, char := range glob.text {
    termbox.SetCell(i+1, height-1, char, termbox.ColorDefault, termbox.ColorDefault);
    termbox.SetCell(i+2, height-1, '_', termbox.ColorDefault, termbox.ColorDefault);
  }
}

// Find subdirs with name matching input base path
func find_subdirs(glob *Global, entrypoints []string, query string, depth int) []string {
	var results []string
  width, height := termbox.Size()

  glob.match_size = glob.match_size[:0]
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
					if strings.Contains(filepath.Base(path), query) && len(path) < width/2-1 {
						results = append(results, path) 

            fileInfo, _ := os.Stat(path)

            fileSize := fileInfo.Size()
            glob.match_size = append(glob.match_size, fileSize)
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

  glob := &Global {
    x: 3, 
    text: "",
    select_mode: false,
  }

  entry_points := []string { "/home" }

  termbox.SetOutputMode(termbox.Output256)

  // Main loop
  for {
    termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

    if !glob.select_mode {
      glob.matches = find_subdirs(glob, entry_points, glob.text, 5)
    }
    
    draw_matches(glob, false)
    draw_text(glob)
    
    termbox.Flush()
  
    event := termbox.PollEvent()
    input_mgr(glob, event)
  }
}

