Managing personal and business documents can quickly become overwhelming as their number grows. Many of these files — often PDFs or scanned images — arrive in a continuous flow: scanned paperwork, email attachments, or downloaded forms, each named by date and a few identifying keywords such as topic, person, or year (e.g., `2023-12-01_Steuerbescheid.pdf`, `2024-05-02_Haftpflichtversicherung.pdf`). Over time, this results in a chronological archive that can be difficult to search just by filename.

While macOS provides built-in file tagging, its functionality and workflows may not fit all needs — especially for those who prefer a fast, terminal-based workflow or require flexible, folder-local catalogs. This is where **filemac** comes in.

**filemac** is a simple, native command-line tool to help you quickly tag, catalog, and retrieve documents stored in such folder structures. It enables you to assign and manipulate tags for files and URLs, all managed through a minimal `.cat` file in each directory. With commands to search, list, and update tags, `filemac` helps you cut through the noise and find documents using meaningful keywords, right from the terminal.

Whether your workflow involves dumping scans into date-named files or archiving diverse documents under project folders, `filemac` makes it easy to organize and query your personal archive without relying on heavyweight or proprietary solutions.

# filemac


Native macOS interactive file catalog/tagging shell for quickly tagging, organizing, and searching documents in your folders.

## Building

Requires Go 1.20+ (developed on Go 1.24, darwin/arm64).

To build the CLI binary:

```sh
go build -o filemac ./cmd/filemac
```

Or via Makefile:

```sh
make build
```

## Running


## Usage

After building, run:

```sh
./filemac
```

You'll enter an interactive shell (`filemac [cwd]>`) where you enter commands directly:

#### Catalog sync (do after adding/removing files!):
    i           # or: init
        Scan working directory: adds all visible files to .cat, removes vanished ones

#### Navigation/display:
    cd <path>   # Change directory (supports ~ expansion)
    ls          # List files in directory as a table (numbered)
        num    | type  | name
        1      | file  | example.pdf
    vc          # View catalog (.cat), numbers here are for tag operations
    vc -new     # Show only files and URLs with no tags yet (new entries)
    vl          # View .catlink file (linked folders)
    lt          # List all unique tags (from .cat, or from linked dirs if only .catlink)

#### Tag & catalog management:
    a <num> <tag>      # Add tag to catalog entry (as numbered in 'vc')
    ax <tag>           # Add tag to all catalog entries
    d <num> <tag>      # Remove tag from entry
    dx <tag>           # Remove tag from all
    r <num> <t1> <t2>  # Replace tag t1 with t2 in entry
    rx <t1> <t2>       # Replace tag t1 with t2 in all entries
    w [<num>]          # Walkthrough/interactive tag fixer
    link <path...>     # Create or overwrite .catlink file with absolute paths

#### Search:
    s <tag...>         # Search for ALL tags (AND, use !tag to exclude)
        e.g. s work 2023 !private
    sl                 # Interactive search loop (search, open file, repeat/quit)

#### Housekeeping:
    help               # Show command list
    quit, exit         # Exit shell

---

**Important:** Always run `i` (or `init`) when you add, remove, or move files using the OS, before performing tag operations! This will sync the working set with the `.cat` and keep numbers/tags consistent.

**History:** All loops (main and search) support up/down arrow for in-session history browsing and editing.

**.catlink**: replaces `.linkcat`. Use `link ...` to create/update, not by hand.

View (`vc`, `lt`, `vl`) never create or modify files. Tag add/remove only changes `.cat`. To create a `.cat`, use `init` first.

---

## Example session

```
$ ./filemac
filemac [~/docs]> ls
num    | type  | name
1      | file  | Urlaub.pdf
2      | file  | Rechnung-scan.pdf
filemac [~/docs]> i
.cat synchronized: 2 added, 0 removed
filemac [~/docs]> vc
num    | type  | name                | tags
1      | file  | Urlaub.pdf          |
2      | file  | Rechnung-scan.pdf   |
filemac [~/docs]> a 1 travel
tag 'travel' added to entry 1
filemac [~/docs]> ax processed
tag 'processed' added to 2 entries
filemac [~/docs]> vc
num    | type  | name                | tags
1      | file  | Urlaub.pdf          | travel, processed
2      | file  | Rechnung-scan.pdf   | processed
filemac [~/docs]> s travel
/Users/alex/docs/Urlaub.pdf
travel, processed

filemac [~/docs]> help
# ...command list as above...
filemac [~/docs]> quit
```
