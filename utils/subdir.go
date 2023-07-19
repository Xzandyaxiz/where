package utils

import (
  "path/filepath"
  "os"; "strings"
  "github.com/nsf/termbox-go"
)

// Find subdirs with name matching input base path
func Find_subdirs(glob *Global, entrypoints []string, query string, depth int) []string {
	var results []string
  width, height := termbox.Size()

  glob.Match_size = glob.Match_size[:0]
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
            glob.Match_size = append(glob.Match_size, fileSize)
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
