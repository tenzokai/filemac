package tags

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/tenzokai/filemac/pkg/catalog"
)

// List unique tags
func CmdListTags() {
	entries, err := catalog.LoadCatalog()
	if err != nil {
		fmt.Println("catalog error:", err)
		return
	}
	tagSet := make(map[string]struct{})
	for _, ent := range entries {
		for _, tag := range ent.Tags {
			tagSet[tag] = struct{}{}
		}
	}
	var tags []string
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	if len(tags) == 0 {
		fmt.Println("(no tags found)")
		return
	}
	for _, tag := range tags {
		fmt.Println(tag)
	}
}

func CmdAddTag(numStr string, tag string) {
	entries, err := catalog.LoadCatalog()
	if err != nil {
		fmt.Println("catalog error:", err)
		return
	}
	n, err := strconv.Atoi(numStr)
	if err != nil || n < 1 || n > len(entries) {
		fmt.Printf("invalid entry number: %v\n", numStr)
		return
	}
	entry := &entries[n-1]
	tag = strings.TrimSpace(tag)
	if tag == "" {
		fmt.Println("Empty tag not allowed")
		return
	}
	for _, t := range entry.Tags {
		if t == tag {
			fmt.Println("tag already present")
			return
		}
	}
	entry.Tags = append(entry.Tags, tag)
	if err := catalog.SaveCatalog(entries); err != nil {
		fmt.Println("error saving catalog:", err)
		return
	}
	fmt.Printf("tag '%s' added to entry %d\n", tag, n)
}

// saveCatalog writes all entries back to .cat
func CmdAddTagAll(tag string) {
	entries, err := catalog.LoadCatalog()
	if err != nil {
		fmt.Println("catalog error:", err)
		return
	}
	tag = strings.TrimSpace(tag)
	if tag == "" {
		fmt.Println("Empty tag not allowed")
		return
	}
	count := 0
	for i := range entries {
		found := false
		for _, t := range entries[i].Tags {
			if t == tag {
				found = true
				break
			}
		}
		if !found {
			entries[i].Tags = append(entries[i].Tags, tag)
			count++
		}
	}
	if count == 0 {
		fmt.Println("tag already present on all entries, nothing to add")
		return
	}
	if err := catalog.SaveCatalog(entries); err != nil {
		fmt.Println("error saving catalog:", err)
		return
	}
	fmt.Printf("tag '%s' added to %d entries\n", tag, count)
}
func CmdRemoveTag(numStr string, tag string) {
	entries, err := catalog.LoadCatalog()
	if err != nil {
		fmt.Println("catalog error:", err)
		return
	}
	n, err := strconv.Atoi(numStr)
	if err != nil || n < 1 || n > len(entries) {
		fmt.Printf("invalid entry number: %v\n", numStr)
		return
	}
	entry := &entries[n-1]
	tag = strings.TrimSpace(tag)
	if tag == "" {
		fmt.Println("Empty tag not allowed")
		return
	}
	found := false
	var newTags []string
	for _, t := range entry.Tags {
		if t == tag {
			found = true
			continue
		}
		newTags = append(newTags, t)
	}
	if !found {
		fmt.Println("tag not found on entry")
		return
	}
	entry.Tags = newTags
	if err := catalog.SaveCatalog(entries); err != nil {
		fmt.Println("error saving catalog:", err)
		return
	}
	fmt.Printf("tag '%s' removed from entry %d\n", tag, n)
}
func CmdRemoveTagAll(tag string) {
	entries, err := catalog.LoadCatalog()
	if err != nil {
		fmt.Println("catalog error:", err)
		return
	}
	tag = strings.TrimSpace(tag)
	if tag == "" {
		fmt.Println("Empty tag not allowed")
		return
	}
	count := 0
	for i := range entries {
		found := false
		var newTags []string
		for _, t := range entries[i].Tags {
			if t == tag {
				found = true
				continue
			}
			newTags = append(newTags, t)
		}
		if found {
			entries[i].Tags = newTags
			count++
		}
	}
	if count == 0 {
		fmt.Println("tag not found on any entry, nothing to remove")
		return
	}
	if err := catalog.SaveCatalog(entries); err != nil {
		fmt.Println("error saving catalog:", err)
		return
	}
	fmt.Printf("tag '%s' removed from %d entries\n", tag, count)
}
func CmdReplaceTag(numStr, t1, t2 string) {
	entries, err := catalog.LoadCatalog()
	if err != nil {
		fmt.Println("catalog error:", err)
		return
	}
	n, err := strconv.Atoi(numStr)
	if err != nil || n < 1 || n > len(entries) {
		fmt.Printf("invalid entry number: %v\n", numStr)
		return
	}
	t1 = strings.TrimSpace(t1)
	t2 = strings.TrimSpace(t2)
	if t1 == "" || t2 == "" {
		fmt.Println("tags must be non-empty")
		return
	}
	if t1 == t2 {
		fmt.Println("tags must be different")
		return
	}
	entry := &entries[n-1]
	found := false
	already := false
	for _, t := range entry.Tags {
		if t == t1 {
			found = true
		}
		if t == t2 {
			already = true
		}
	}
	if !found {
		fmt.Println("tag not found on entry")
		return
	}
	var newTags []string
	for _, t := range entry.Tags {
		if t == t1 {
			continue
		}
		newTags = append(newTags, t)
	}
	if !already {
		newTags = append(newTags, t2)
	}
	entry.Tags = newTags
	if err := catalog.SaveCatalog(entries); err != nil {
		fmt.Println("error saving catalog:", err)
		return
	}
	fmt.Printf("tag '%s' replaced with '%s' in entry %d\n", t1, t2, n)
}
func CmdReplaceTagAll(t1, t2 string) {
	entries, err := catalog.LoadCatalog()
	if err != nil {
		fmt.Println("catalog error:", err)
		return
	}
	t1 = strings.TrimSpace(t1)
	t2 = strings.TrimSpace(t2)
	if t1 == "" || t2 == "" {
		fmt.Println("tags must be non-empty")
		return
	}
	if t1 == t2 {
		fmt.Println("tags must be different")
		return
	}
	count := 0
	for i := range entries {
		found := false
		already := false
		for _, t := range entries[i].Tags {
			if t == t1 {
				found = true
			}
			if t == t2 {
				already = true
			}
		}
		if !found {
			continue
		}
		var newTags []string
		for _, t := range entries[i].Tags {
			if t == t1 {
				continue
			}
			newTags = append(newTags, t)
		}
		if !already {
			newTags = append(newTags, t2)
		}
		entries[i].Tags = newTags
		count++
	}
	if count == 0 {
		fmt.Println("tag not found on any entry")
		return
	}
	if err := catalog.SaveCatalog(entries); err != nil {
		fmt.Println("error saving catalog:", err)
		return
	}
	fmt.Printf("tag '%s' replaced with '%s' in %d entries\n", t1, t2, count)
}

func CmdSearch(tags []string) {
	// Split search terms into includes and excludes
	var includes, excludes []string
	for _, t := range tags {
		t = strings.TrimSpace(t)
		if strings.HasPrefix(t, "!") && len(t) > 1 {
			excludes = append(excludes, t[1:])
		} else if t != "" {
			includes = append(includes, t)
		}
	}

	hasTags := func(entryTags []string) bool {
		// Must have all includes
		for _, want := range includes {
			found := false
			for _, have := range entryTags {
				if want == have {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		// Must have none of the excludes
		for _, notag := range excludes {
			for _, have := range entryTags {
				if notag == have {
					return false
				}
			}
		}
		return true
	}

	// Helper for printing result lines
	printEntry := func(e catalog.CatEntry, dir string) {
		path := e.Name
		if e.Type == "file" && dir != "" && !filepath.IsAbs(path) {
			path = filepath.Join(dir, path)
			path, _ = filepath.Abs(path)
		}
		fmt.Println(path)
		fmt.Println("   " + strings.Join(e.Tags, ", "))
		fmt.Println("")
	}

	catExists := func(file string) bool {
		s, err := os.Stat(file)
		return err == nil && !s.IsDir()
	}

	cwd, _ := os.Getwd()
	catPath := filepath.Join(cwd, ".cat")
	linkPath := filepath.Join(cwd, ".linkcat")

	if catExists(catPath) {
		entries, err := catalog.LoadCatalog()
		if err != nil {
			fmt.Println("catalog error:", err)
			return
		}
		for _, e := range entries {
			if hasTags(e.Tags) {
				printEntry(e, cwd)
			}
		}
		return
	} else if catExists(linkPath) {
		// Search all referenced .cat files
		f, err := os.Open(linkPath)
		if err != nil {
			fmt.Println("could not open .linkcat:", err)
			return
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			dir := strings.TrimSpace(scanner.Text())
			if dir == "" {
				continue
			}
			catfile := filepath.Join(dir, ".cat")
			if !catExists(catfile) {
				continue
			}
			entries, err := readCatalogAt(catfile)
			if err != nil {
				fmt.Printf("error reading %s: %v\n", catfile, err)
				continue
			}
			for _, e := range entries {
				if hasTags(e.Tags) {
					printEntry(e, dir)
				}
			}
		}
		return
	}
	// Neither .cat nor .linkcat
	fmt.Printf("(No .cat or .linkcat found in %s)\n", cwd)
}

// Like LoadCatalog but at arbitrary filename
func readCatalogAt(filename string) ([]catalog.CatEntry, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var entries []catalog.CatEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		entry := catalog.ParseCatalogLine(line)
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

// CmdSearchLoop: interactive search and open
func CmdSearchLoop() {
	reader := bufio.NewReader(os.Stdin)

	type Match struct {
		Tags []string
		Path string
		Type string // file/url
	}

	var matches []Match
	var printResults = func() {
		if len(matches) == 0 {
			fmt.Println("No matches.")
			return
		}
		for idx, m := range matches {
			fmt.Printf("%d ", idx+1)
			fmt.Println(strings.Join(m.Tags, " "))
			fmt.Println(m.Path)
		}
	}

	var performSearch = func(q []string) {
		// search logic from CmdSearch, but stores matches as []Match
		var includes, excludes []string
		for _, t := range q {
			t = strings.TrimSpace(t)
			if strings.HasPrefix(t, "!") && len(t) > 1 {
				excludes = append(excludes, t[1:])
			} else if t != "" {
				includes = append(includes, t)
			}
		}
		matches = matches[:0]
		cwd, _ := os.Getwd()
		catfile := filepath.Join(cwd, ".cat")
		linkfile := filepath.Join(cwd, ".linkcat")
		catExists := func(file string) bool {
			s, err := os.Stat(file)
			return err == nil && !s.IsDir()
		}
		hasTags := func(entryTags []string) bool {
			for _, want := range includes {
				found := false
				for _, have := range entryTags {
					if want == have {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
			for _, notag := range excludes {
				for _, have := range entryTags {
					if notag == have {
						return false
					}
				}
			}
			return true
		}
		// catalog check
		if catExists(catfile) {
			entries, err := catalog.LoadCatalog()
			if err != nil {
				fmt.Println("catalog error:", err)
				return
			}
			for _, e := range entries {
				if hasTags(e.Tags) {
					abspath := e.Name
					if e.Type == "file" && !filepath.IsAbs(e.Name) {
						abspath = filepath.Join(cwd, e.Name)
						abspath, _ = filepath.Abs(abspath)
					}
					matches = append(matches, Match{Tags: e.Tags, Path: abspath, Type: e.Type})
				}
			}
			return
		} else if catExists(linkfile) {
			f, err := os.Open(linkfile)
			if err != nil {
				fmt.Println("could not open .linkcat:", err)
				return
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				dir := strings.TrimSpace(scanner.Text())
				if dir == "" {
					continue
				}
				catfile := filepath.Join(dir, ".cat")
				if !catExists(catfile) {
					continue
				}
				entries, err := readCatalogAt(catfile)
				if err != nil {
					continue
				}
				for _, e := range entries {
					if hasTags(e.Tags) {
						abspath := e.Name
						if e.Type == "file" && !filepath.IsAbs(e.Name) {
							abspath = filepath.Join(dir, e.Name)
							abspath, _ = filepath.Abs(abspath)
						}
						matches = append(matches, Match{Tags: e.Tags, Path: abspath, Type: e.Type})
					}
				}
			}
			return
		}
		fmt.Println("(No .cat or .linkcat found in current dir)")
	}

	fmt.Println("Interactive search loop. Enter:")
	fmt.Println("    s <tag> [!notag] ...   to (re)search")
	fmt.Println("    o <n> | o <path>       to open match")
	fmt.Println("    q                      to quit")
	for {
		fmt.Print("SearchLoop> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		cmd := parts[0]
		switch cmd {
		case "q":
			fmt.Println("Quitting.")
			return
		case "s":
			q := parts[1:]
			performSearch(q)
			printResults()
		case "o":
			if len(parts) < 2 {
				fmt.Println("Usage: o <match-number|path>")
				continue
			}
			arg := parts[1]
			idx, err := strconv.Atoi(arg)
			var tgt Match
			if err == nil && idx >= 1 && idx <= len(matches) {
				tgt = matches[idx-1]
			} else {
				// try path match
				found := false
				for _, m := range matches {
					if m.Path == arg {
						tgt = m
						found = true
						break
					}
				}
				if !found {
					fmt.Println("Path not in last results.")
					continue
				}
			}
			if tgt.Type == "file" {
				fmt.Printf("Opening %s...\n", tgt.Path)
				err := openFile(tgt.Path)
				if err != nil {
					fmt.Println("Error opening:", err)
				}
			} else if tgt.Type == "url" {
				fmt.Println("Cannot open URLs from here.")
			} else {
				fmt.Printf("Entry is not a file.")
			}
		default:
			fmt.Println("Commands: s ... | o <n|path> | q")
		}
	}
}

func openFile(path string) error {
	cmd := exec.Command("open", path)
	return cmd.Run()
}
