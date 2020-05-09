# lbdl

A headless torrent client written in Go.

## The Basics

`lbdl` is designed to run as a very simple standalone torrent client. Start all the torrents and magnet links found in the configured paths, and run until all of them are downloaded.

## Compiling

```bash
./make all
```

## Running (Linux)

```bash
cd build
./lbdl.x86_64.linux 2> lbdl.log
```

## Arguments

* `-d <DIR>` Downloads directory
* `-m <FILE>` Magnet link list file
* `-t <DIR>` Torrent file directory
