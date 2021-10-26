# splits

This command line split file into parts.

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gonejack/splits)
![Build](https://github.com/gonejack/splits/actions/workflows/go.yml/badge.svg)
[![GitHub license](https://img.shields.io/github/license/gonejack/splits.svg?color=blue)](LICENSE)

### Install

```shell
> go get github.com/gonejack/splits
```

### Usage

1. split by parts
```shell
> splits -n 2 test.txt
```

2. split by size
```shell
> splits -b 100k test.txt
```

3. merge
```shell
> cat test.text.* > merged_test.txt
```

```
Flags:
  -h, --help                    Show context-sensitive help.
  -n, --chunks=3                split into n parts
  -b, --size-per-part=STRING    split into sized parts
  -v, --verbose                 Verbose printing.
      --about                   About.
```
