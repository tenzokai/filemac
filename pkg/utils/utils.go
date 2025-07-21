package utils


import (
    "fmt"
    "os"
    "sort"
)


// DirExists checks if a directory exists at the given path.
func DirExists(path string) bool {
    fi, err := os.Stat(path)
    return err == nil && fi.IsDir()
}

// MkdirIfMissing creates the specified directory (including parent dirs) if it does not exist.
func MkdirIfMissing(path string) error {
    if DirExists(path) {
        return nil
    }
    return os.MkdirAll(path, 0755)
}

// ListVisibleFiles lists non-dot files/dirs in the given directory.
func ListVisibleFiles(dir string) ([]string, error) {
    entries, err := os.ReadDir(dir)
    if err != nil {
        return nil, err
    }
    var names []string
    for _, entry := range entries {
        name := entry.Name()
        if len(name) > 0 && name[0] == '.' {
            continue
        }
        if entry.IsDir() {
            name += "/"
        }
        names = append(names, name)
    }
    sort.Strings(names)
    return names, nil
}

// AbsPath returns the absolute path for the given path, or original on error.
func AbsPath(path string) string {
    abspath, err := os.Getwd()
    if err != nil {
        return path
    }
    joined := path
    if !isAbs(path) {
        joined = abspath + string(os.PathSeparator) + path
    }
    norm, err := os.Stat(joined)
    if err == nil && norm.IsDir() {
        return joined
    }
    return joined // fallback
}

// isAbs checks if a path is absolute.
func isAbs(path string) bool {
    return len(path) > 0 && (path[0] == '/' || path[0] == '~')
}

// CmdCd changes into the specified directory, creates it if missing.
func CmdCd(path string) {
    if path == "" {
        fmt.Println("CmdCd: No path given!")
        return
    }
    err := MkdirIfMissing(path)
    if err != nil {
        fmt.Printf("CmdCd: Error creating directory: %v\n", err)
        return
    }
    err = os.Chdir(path)
    if err != nil {
        fmt.Printf("CmdCd: Failed to change directory: %v\n", err)
        return
    }
    fmt.Printf("Changed directory to: %s\n", path)
}

// CmdLs lists visible (non-dotfile) files/dirs in the current directory.
func CmdLs() {
    cwd, err := os.Getwd()
    if err != nil {
        fmt.Printf("CmdLs: error getting current dir: %v\n", err)
        return
    }
    names, err := ListVisibleFiles(cwd)
    if err != nil {
        fmt.Printf("CmdLs: error reading directory: %v\n", err)
        return
    }
    fmt.Printf("num    | type  | name\n")
    n := 1
    for _, name := range names {
        t := "file"
        if isURL(name) {
            t = "url"
        }
        fmt.Printf("%-6d | %-5s | %s\n", n, t, name)
        n++
    }
}

// isURL checks if name looks like a URL (for ls output)
func isURL(str string) bool {
    return (len(str) > 7 && (str[:7] == "http://" || str[:8] == "https://"))
}
