package utils

import (
  "github.com/nsf/termbox-go"
  "strconv"; "os"
)

func Matches_render(glob *Global, line_nums bool) {
  var offset int = 2
  width, _ := termbox.Size()

  for i, str := range glob.Matches {
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
    
    bytes := strconv.Itoa(int(glob.Match_size[i]))
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
    if glob.Select_mode && glob.Select_index == i {
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
func Text_render(glob *Global) {
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

  if glob.Select_mode {
    return
  }
 
  termbox.SetCell(1, height-1, '_', termbox.ColorDefault | termbox.AttrBold, termbox.ColorDefault);  

  // Draw the text from the cursor struct
  for i, char := range glob.Text {
    termbox.SetCell(i+1, height-1, char, termbox.ColorDefault, termbox.ColorDefault);
    termbox.SetCell(i+2, height-1, '_', termbox.ColorDefault, termbox.ColorDefault);
  }
}
