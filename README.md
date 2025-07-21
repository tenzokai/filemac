## filemac

**filemac** is a small macOS terminal tool to help you tag, organize, and search your personal documents.

Many people collect PDF scans, forms, or letters in folders — often with filenames like `2024-05-02_Haftpflicht.pdf`. Over time, these folders grow and become hard to search just by name.

With `filemac`, you can turn any folder into a searchable catalog:
- run `init` to create or update a `.cat` file for all files in the folder
- assign tags to each file (`a`, `ax`, `d`, `r`, etc.)
- use `s` to search for files by tags

You can manage many such folders. Each has its own `.cat` file — no global index, no cloud, no vendor lock-in.

To search across multiple folders, use `link` in any folder to create a `.catlink` file:
- the folder itself doesn't need a `.cat`
- any search (`s`) in this folder will include all linked `.cat` folders
- this allows you to define flexible, local views across your archive


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
    lt          # List all unique tags (from .cat, or from linked dirs if only .catlink exists)

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

Use `link ...` to create/update, not by hand.

View (`vc`, `vc -new`, `lt`, `vl`) never create or modify files. Tag add/remove only changes `.cat`. To create a `.cat`, use `init` first.

---

## Example session


```
$ ./filemac
filemac [~/docs]> i
.cat synchronized: 3 added, 0 removed

filemac [~/docs]> vc
num    | type  | name                 | tags
1      | file  | Rechnung.pdf         |
2      | file  | Urlaub.pdf           |
3      | file  | Versicherung.pdf     |

filemac [~/docs]> a 1 steuer
filemac [~/docs]> ax 2024

filemac [~/docs]> cd ~/search-hub
filemac [~/search-hub]> link ~/docs ~/insurance ~/receipts
linked 3 folders

filemac [~/search-hub]> s steuer
/Users/alex/docs/Rechnung.pdf
steuer, 2024
```

## Outlook

A future version will support semi-automatic tagging. Users will provide a list of tags, and filemac will use a locally stored pre-trained classification model — expected to support several European languages — along with a lightweight inference interface. This will enable automatic tagging of common file types such as txt, pdf, docx, xlsx, and others.

tenzoki, July 2025
