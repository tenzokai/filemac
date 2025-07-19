# filemac Usage Guide

`filemac` is a native macOS CLI tool for lightweight tagging & cataloging of files and URLs in your folders.


## Command Structure

Commands are invoked as:
```sh
filemac -<command> [args...]
```

--

## Commands & Options

### Directory/Listing

#### `-cd <path>`
Change directory (to `<path>`). Will be created if missing.

#### `-ls`
List visible (non-dot) files/directories in current directory.

--

### Catalog Operations

#### `-vc`
View the `.cat` file:
- If missing, creates it.
- Output:
  ```
  num    | type  | name                                 | tags
  1      | file  | 2007-07-11_Jakob-Kindergeld.pdf      | jakob, kindergeld
  2      | url   | https://google.com                   | any, query
  ```
#### `-vl`
View the `.linkcat` file:
- Lists all referenced folder paths, one per line.
- If missing, prints: `No .linkcat in <abs-path-to-folder>`
  (Empty if none)

#### `-lt`
List all unique tags (sorted, one per line).

--

### Tag Operations

#### `-a <num> <tag>`
Add tag `<tag>` to entry `<num>`. No-op if already present.

#### `-ax <tag>`
Add tag `<tag>` to all entries. Does nothing to those already containing the tag.

#### `-d <num> <tag>`
Remove tag `<tag>` from entry `<num>`. Warns if not present.

#### `-dx <tag>`
Remove tag `<tag>` from all entries.

#### `-r <num> <tag1> <tag2>`
Replace tag `<tag1>` with `<tag2>` in entry `<num>`.

#### `-rx <tag1> <tag2>`
Replace tag `<tag1>` with `<tag2>` in all entries.

--

### Walkthrough Mode

#### `-w [<num>]`
Step through entries interactively, starting at specified number (default: 1).

Walkthrough steps:
- Show entry and current tags
- Prompt: `Enter tags:` (comma-separated; new set replaces old; blank to skip)
- Confirm: `Correct? y/n` (repeat input if "n")
- Type `stop` to abort early – last changed entry and unfinished count shown

--

### Linking/Linked Directory Search

#### `-link <path1> <path2> ...`
Creates `.linkcat` referencing folders line by line. Used for cross-folder catalog search.

#### `-s <tag1> <tag2> ...`

List all entries matching ALL tags (AND search). You can exclude tags by prefixing them with `!` (negation). For example:

   filemac -s travel !private
Matches all entries tagged with `travel` and NOT tagged with `private`. You can chain as many includes and excludes as needed (evaluation is AND across all terms).

- If directory has `.cat`, search it.
- Else if `.linkcat`, search each referenced folder’s `.cat`.
- If neither: prints warning.

Result example:
```
/absolute/path/to/file
tag1, tag2, tag3

/absolute/path/to/otherfile
tag2, tag4
```

## Notes & Conventions

- Catalog entries: `<name>*<tag1>*<tag2>...`
- Entry numbers correspond to the listing in `-vc`
- Tags are always deduplicated per entry
- Modifications write `.cat` atomically (full-rewrite)

## Example Workflow

```sh
# Change/create directory
filemac -cd docs

# List files in this directory
filemac -ls

# View catalog, see entry numbers
filemac -vc

# Add tag “travel” to entry 2
filemac -a 2 travel

# Remove tag “work” from entry 1
filemac -d 1 work

# Start step-by-step interactive tagging
filemac -w

# Add a tag to all
filemac -ax family

# Replace a tag across all
filemac -rx child kid

# Search for entries matching BOTH “family” and “2025”
filemac -s family 2025
```

---
**Platform:** Tested on macOS aarch64 (Apple Silicon). Requires Go if building from source.
#### `-sl`
Interactive search loop. Shows results like this:

   1 tag1 tag2 ...
   /abs/path/to/file1
   2 tag3 tag4 ...
   /abs/path/to/file2

Loop commands:
```
- s [!]<tag1> [!]<tag2> ...   (re)search with AND/NOT terms
- o <number>   open n-th result file
- o <path>     open result file by full path
- q            quit search loop
```
If you try to open a URL or non-file, a warning is shown.
