package glob

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

/*
 * glob - an expanded version
 *
 * This implementation of globbing will still take advantage of the Glob
 * function in path/filepath, however this extends the pattern to include '**'
 *
 */

/*
Algorithm details

define Entry:
    path, index into glob

Base Case:
    add Entry{root, 0}

while num entries > 0
    given an entry (path, idx)
    given glob segment (gb) at idx

    if gb == **
        move cur entry idx + 1

        for each dir inside path
            add new Entry{dir, idx}
    else
        add gb to path
        check for any results from normal globbing
        if none
            remove entry
        else
            add an entry{result, idx + 1} unless idx + 1 is out of bounds

    keep current entry if it's idx is in bounds

 */

type matchEntry struct {
    path string
    idx int
}

func Glob(root string, pattern string) (matches []string, e error) {
    // TODO check if it even contains **, if not, just use normal filepath.Glob
    segments := strings.Split(pattern, string(os.PathSeparator))

    workingEntries := []matchEntry{
        matchEntry{path: root, idx: 0},
    }

    for len(workingEntries) > 0 {
        var temp []matchEntry
        for _, entry := range workingEntries {
            workingPath := entry.path
            idx := entry.idx
            segment := segments[entry.idx]

            if segment == "**" {
                entry.idx++

                isdir, err := isDir(entry.path)

                if !isdir || err != nil {
                    continue
                }

                d, err := os.Open(entry.path)
                if err != nil {
                    return
                }
                
                names, err := d.Readdirnames(-1)

                for _, name := range names {
                    path := filepath.Join(workingPath, name)
                    isdir, err = isDir(path)
                    
                    if err != nil {
                        continue
                    }

                    if isdir {
                        newEntry := matchEntry{
                            path: path,
                            idx: idx,
                        }

                        temp = append(temp, newEntry)
                    }
                }

            } else {
                path := filepath.Join(workingPath, segment)
                results, _ := filepath.Glob(path)

                for _, result := range results {
                    newEntry := matchEntry{
                        path: result,
                        idx: idx + 1,
                    }

                    if idx + 1 < len(segments) {
                        temp = append(temp, newEntry)                        
                    } else {
                        matches = append(matches, result)
                    }
                }
                entry.idx = len(segments)
            }

            if entry.idx < len(segments) {
                temp = append(temp, entry)
            }

            workingEntries = temp
        }
    }

    return
}

func isDir(path string) (bool, error) {
    fi, err := os.Stat(path)

    if err != nil {
        return false, err
    }

    return fi.IsDir(), nil
}

