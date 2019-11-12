# What is "psh"

When type like "ps u -h", expect then PSZ and VSS show
human-readable formatted byte number with unit.
(ex: "2.6Gi", "123Bi")

But don't work expected. The reality, they are
no-formatted byte number without unit.

So "psh” is my answer to solve this.

# Usage

```bash
$ ./psh # can use ps option, ex: psh ax
```

# Installation

```bash
$ go get github.com/tanakakz/psh
```

# License

MIT

# Author

Kazutoshi Tanaka (a.k.a. tanakakz)
