package catalog

import (
    "bufio"
    "fmt"
    "os"
    "strings"
"strconv")

var CatalogFilename = ".cat"

type CatEntry struct {
    Raw    string
    Name   string
    Type   string // "file" or "url"
    Tags   []string
}

// Load catalog file, create if not exists
func LoadCatalog() ([]CatEntry, error) {
    // Ensure file exists
    _, err := os.Stat(CatalogFilename)
    if os.IsNotExist(err) {
        file, cerr := os.Create(CatalogFilename)
        if cerr != nil {
            return nil, cerr
        }
        file.Close()
        return []CatEntry{}, nil
    } else if err != nil {
        return nil, err
    }

    f, err := os.Open(CatalogFilename)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var entries []CatEntry
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" { continue }
        entry := ParseCatalogLine(line)
        entries = append(entries, entry)
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    return entries, nil
}

func ParseCatalogLine(line string) CatEntry {
    // Parse using only '*' as delimiter: first '*' separates name from tags.
    var e CatEntry
    e.Raw = line
    e.Tags = nil

    namePart := line
    tagsPart := ""
    if idx := strings.Index(line, "*"); idx != -1 {
        namePart = line[:idx]
        tagsPart = line[idx+1:]
    }
    e.Name = namePart
    if tagsPart != "" {
        for _, tag := range strings.Split(tagsPart, "*") {
            tag = strings.TrimSpace(tag)
            if tag != "" {
                e.Tags = append(e.Tags, tag)
            }
        }
    }

    // Type logic per user: if name contains ':' (anywhere), type=url, else file
    if strings.Contains(e.Name, ":") {
        e.Type = "url"
    } else {
        e.Type = "file"
    }
    return e
}

// Print the catalog as a table
func printCatalog(entries []CatEntry) {
    fmt.Printf("%-6s | %-6s | %-36s | %s\n", "num", "type", "name", "tags")
    for i, e := range entries {
        entryType := e.Type
        if strings.Contains(e.Name, ":") {
            entryType = "url"
        } else {
            entryType = "file"
        }
        //fmt.Printf("// DEBUG: i=%d name='%s' type='%s'\n", i+1, e.Name, entryType)
        fmt.Printf("%-6d | %-6s | %-36s | %s\n", i+1, entryType, truncate(e.Name,36), strings.Join(e.Tags, ", "))
    }
    if len(entries) == 0 {
        fmt.Printf("(catalog is empty)\n")
    }
}

// Helper for truncating name display
func truncate(s string, n int) string {
    if len([]rune(s)) <= n { return s }
    runes := []rune(s)
    return string(runes[:n-3]) + "..."
}

// CLI handler for -vc
func CmdViewCat() {
    entries, err := LoadCatalog()
    if err != nil {
        fmt.Println("catalog error:", err)
        return
    }
    printCatalog(entries)
}

func CmdWalkthrough(num string) {
    entries, err := LoadCatalog()
    if err != nil {
        fmt.Println("catalog error:", err)
        return
    }
    if len(entries) == 0 {
        fmt.Println("No entries in catalog.")
        return
    }
    startIdx := 0
    if num != "" {
        n, err := strconv.Atoi(num)
        if err == nil && n >= 1 && n <= len(entries) {
            startIdx = n - 1
        }
    }
    reader := bufio.NewReader(os.Stdin)
    lastChangedIdx := -1
    skipped := 0
    for i := startIdx; i < len(entries); i++ {
        fmt.Printf("\nEntry #%d / %d:\n", i+1, len(entries))
        fmt.Printf("  Name: %s\n  Tags: %s\n", entries[i].Name, strings.Join(entries[i].Tags, ", "))
    walkthroughInput:
        for {
            fmt.Print("Enter tags (comma-separated) or blank to skip, 'stop' to abort: ")
            line, _ := reader.ReadString('\n')
            line = strings.TrimSpace(line)
            if line == "stop" {
                fmt.Printf("Aborted at entry #%d (%d left).\n", i+1, len(entries)-i)
                if lastChangedIdx != -1 {
                    fmt.Printf("Last changed entry: #%d\n", lastChangedIdx+1)
                }
                return
            }
            if line == "" {
                skipped++
                break walkthroughInput
            }
            // parse, apply tags
            var tags []string
            for _, t := range strings.Split(line, ",") {
                t = strings.TrimSpace(t)
                if t != "" {
                    tags = append(tags, t)
                }
            }
            fmt.Printf("New tags: %s\n", strings.Join(tags, ", "))
            fmt.Print("Correct? y/n: ")
            yn, _ := reader.ReadString('\n')
            yn = strings.ToLower(strings.TrimSpace(yn))
            if yn == "y" {
                entries[i].Tags = tags
                lastChangedIdx = i
                // Save after each update
                saveErr := SaveCatalog(entries)
                if saveErr != nil {
                    fmt.Printf("Error saving: %v\n", saveErr)
                } else {
                    fmt.Println("Saved.")
                }
                break walkthroughInput
            } // else repeat entry
        }
    }
    fmt.Printf("Finished walkthrough. Changed %d, skipped %d entries.\n", lastChangedIdx+1, skipped)
}
func CmdLink(paths []string) {
    if len(paths) == 0 {
        fmt.Println("No paths given for .linkcat")
        return
    }
    f, err := os.Create(".linkcat")
    if err != nil {
        fmt.Printf("error creating .linkcat: %v\n", err)
        return
    }
    defer f.Close()
    for _, p := range paths {
        cleaned := strings.TrimSpace(p)
        if cleaned == "" {
            continue
        }
        fmt.Fprintln(f, cleaned)
    }
    fmt.Printf(".linkcat created with %d paths\n", len(paths))
}
// SaveCatalog writes all entries back to .cat
func SaveCatalog(entries []CatEntry) error {
    f, err := os.Create(CatalogFilename)
    if err != nil {
        return err
    }
    defer f.Close()
    for _, e := range entries {
        line := e.Name
        if len(e.Tags) > 0 {
            line += "*" + strings.Join(e.Tags, "*")
        }
        if _, err := f.WriteString(line + "\n"); err != nil {
            return err
        }
    }
    return nil
}
// CmdViewLinkcat prints out contents of .linkcat or a warning if missing.
func CmdViewLinkcat() {
    linkcat := ".linkcat"
    cwd, err := os.Getwd()
    if err != nil {
        fmt.Println("Error: could not determine current directory")
        return
    }
    f, err := os.Open(linkcat)
    if err != nil {
        fmt.Printf("No .linkcat in %s\n", cwd)
        return
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    count := 0
    for scanner.Scan() {
        line := scanner.Text()
        fmt.Println(line)
        count++
    }
    if count == 0 {
        fmt.Println("(Empty .linkcat file)")
    }
}
