# gosimfiles
A tool for searching for similarities between files.

## gosimfiles --help
```
usage: simfile --reffile=REFFILE [<flags>] <-->

A command-line tool to search for similarities between files.

Flags:
  --help             Show help.
  --reffile=REFFILE  Input file.
  --minsim=90.0      Minimal similarity (%).
  --slicelen=-1      Slice length of files before comparing.
  --progress         Show progress output
  --version          Show application version.

Args:
  <-->  One or more files to compare with reffile
```

## Basic usage
```
gosimfiles --reffile file-1.txt -- file-2.txt file-3.txt file-4.txt
```
