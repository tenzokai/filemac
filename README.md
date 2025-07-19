Managing personal and business documents can quickly become overwhelming as their number grows. Many of these files — often PDFs or scanned images — arrive in a continuous flow: scanned paperwork, email attachments, or downloaded forms, each named by date and a few identifying keywords such as topic, person, or year (e.g., `2023-12-01_Steuerbescheid.pdf`, `2024-05-02_Haftpflichtversicherung.pdf`). Over time, this results in a chronological archive that can be difficult to search just by filename.

While macOS provides built-in file tagging, its functionality and workflows may not fit all needs — especially for those who prefer a fast, terminal-based workflow or require flexible, folder-local catalogs. This is where **filemac** comes in.

**filemac** is a simple, native command-line tool to help you quickly tag, catalog, and retrieve documents stored in such folder structures. It enables you to assign and manipulate tags for files and URLs, all managed through a minimal `.cat` file in each directory. With commands to search, list, and update tags, `filemac` helps you cut through the noise and find documents using meaningful keywords, right from the terminal.

Whether your workflow involves dumping scans into date-named files or archiving diverse documents under project folders, `filemac` makes it easy to organize and query your personal archive without relying on heavyweight or proprietary solutions.


# filemac

Native macOS command-line tool for lightweight file tagging and cataloging.

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

After building, typical usage is:

```sh
./filemac -h
```

This displays all supported commands / flags.

Supported commands:

**View**
```sh
./filemac -cd mydir      # Change or create directory
./filemac -ls            # List directory contents
./filemac -lt            # List all tags
./filemac -vc            # View current catalog
./filemac -vl            # View .linkcat (linked folders)
```
**Manage tags and catalogs**
```sh
./filemac -a <num> <tag>     # Add a tag to entry
./filemac -ax <tag>          # Add tag to all entries
./filemac -d <num> <tag>     # Remove tag from entry
./filemac -dx <tag>          # Remove tag from all
./filemac -r <num> <t1> <t2> # Replace tag t1 with t2 in entry
./filemac -rx <t1> <t2>      # Replace tag t1 with t2 in all entries
./filemac -w [<num>]         # Walkthrough mode (interactive tag fixer)
./filemac -link <p1>...      # Create .linkcat referencing folders with a .cat file
```
**Search files**
```sh
./filemac -s travel !private    # Find items with tag 'travel' and NOT 'private'
./filemac -sl            		# Interactive search loop (advanced)
```

Search: To exclude items tagged a certain way, prefix the tag with `!` (AND logic for all terms).
E.g. `filemac -s contract !private` finds all items tagged `contract` but not `private`.

For more advanced selection/opening, use interactive search: `filemac -sl`.
You can repeatedly search using `s ...`, open any match (file), or quit. Result output:
```
    1 tag1 tag2 ...
    /abs/path1
    2 tag3 tag4 ...
    /abs/path2
```
Commands:
```
  s [!]<tag1> [!]<tag2> ...
  o <number-of-match> OR o <abs-path>
  q
```
