package catalog

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var CatalogFilename = ".cat"

type CatEntry struct {
	Raw  string
	Name string
	Type string // "file" or "url"
	Tags []string
}

// Load catalog file, create if not exists
func LoadCatalog() ([]CatEntry, error) {
	// Only open if file exists
	_, err := os.Stat(CatalogFilename)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf(".cat not found in this directory")
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
		if line == "" {
			continue
		}
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
// Optionally print with preserved original indices (1-based positions from underlying catalog file).
func printCatalogWithIndices(entries []CatEntry, indices []int) {
       fmt.Printf("%-6s | %-6s | %-36s | %s\n", "num", "type", "name", "tags")
       for i, e := range entries {
               entryType := e.Type
               if strings.Contains(e.Name, ":") {
                       entryType = "url"
               } else {
                       entryType = "file"
               }
               num := i+1
               if indices != nil && i < len(indices) {
                       num = indices[i]+1
               }
               fmt.Printf("%-6d | %-6s | %-36s | %s\n", num, entryType, truncate(e.Name, 36), strings.Join(e.Tags, ", "))
       }
       if len(entries) == 0 {
               fmt.Printf("(catalog is empty)\n")
       }
}

// Backwards compat: printCatalog with no indices.
func printCatalog(entries []CatEntry) {
       printCatalogWithIndices(entries, nil)
}

// Helper for truncating name display
func truncate(s string, n int) string {
	if len([]rune(s)) <= n {
		return s
	}
	runes := []rune(s)
	return string(runes[:n-3]) + "..."
}

// CmdViewCat optionally takes "-new". If used, only entries with no tags are shown.
func CmdViewCat(args ...string) {
	entries, err := LoadCatalog()
	if err != nil {
		fmt.Println("catalog error:", err)
		return
	}
	fmt.Printf("args: %v args-lens: %v\n", args, len(args))
       if len(args) > 0 && args[0] == "-new" {
               var newEntries []CatEntry
               var indices []int
               for idx, e := range entries {
                       if len(e.Tags) == 0 {
                               newEntries = append(newEntries, e)
                               indices = append(indices, idx)
                       }
               }
               printCatalogWithIndices(newEntries, indices)
       } else {
               printCatalog(entries)
       }
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
		fmt.Println("No paths given for .catlink")
		return
	}
	f, err := os.Create(".catlink")
	if err != nil {
		fmt.Printf("error creating .catlink: %v\n", err)
		return
	}
	defer f.Close()
	for _, p := range paths {
		cleaned := strings.TrimSpace(p)
		if cleaned == "" {
			continue
		}
		abs := cleaned
		if !filepath.IsAbs(cleaned) {
			abs, _ = filepath.Abs(cleaned)
		}
		fmt.Fprintln(f, abs)
	}
	fmt.Printf(".catlink created with %d paths\n", len(paths))
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

// CmdViewLinkcat prints out contents of .catlink or a warning if missing.
func CmdViewLinkcat() {
	linkcat := ".catlink"
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error: could not determine current directory")
		return
	}
	f, err := os.Open(linkcat)
	if err != nil {
		fmt.Printf("No .catlink in %s\n", cwd)
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
		fmt.Println("(Empty .catlink file)")
	}
}

// CmdInitCatalog: synchronize .cat with directory files (add new, remove vanished)
func CmdInitCatalog() {
	cwd, _ := os.Getwd()
	files, err := os.ReadDir(cwd)
	if err != nil {
		fmt.Printf("init error: %v\n", err)
		return
	}
	fileSet := make(map[string]struct{})
	for _, f := range files {
		name := f.Name()
		if len(name) != 0 && name[0] != '.' && !f.IsDir() {
			fileSet[name] = struct{}{}
		}
	}

	// Load .cat
	entries, err := LoadCatalog()
	if err != nil {
		fmt.Printf("init error loading .cat: %v\n", err)
		return
	}
	// Keep entries whose file exists
	var kept []CatEntry
	seen := make(map[string]struct{})
	for _, e := range entries {
		if _, ok := fileSet[e.Name]; ok {
			kept = append(kept, e)
			seen[e.Name] = struct{}{}
		}
	}
	removed := len(entries) - len(kept)

	// Add new file entries
	added := 0
	for name := range fileSet {
		if _, ok := seen[name]; !ok {
			kept = append(kept, CatEntry{Name: name, Type: "file"})
			added++
		}
	}
	SaveCatalog(kept)
	fmt.Printf(".cat synchronized: %d added, %d removed\n", added, removed)
}

// LoadCatalogAt loads catalog entries from a given filepath.
func LoadCatalogAt(filename string) ([]CatEntry, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var entries []CatEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		entry := ParseCatalogLine(line)
		entries = append(entries, entry)
	}
	return entries, nil
}
