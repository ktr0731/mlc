# mlc
Execute shell commands concurrency

===

## Equipments
- Go
- Linux, UNIX or macOS

## Installation
``` sh
$ go get github.com/lycoris0731/mlc
```

## Usage
``` sh
$ mlc command[, command...]
```

### Sample  
``` sh
$ mlc "command1" "command2"
```
Thanks to [mattn/go-shellwords](https://github.com/mattn/go-shellwords), we can use backtick and dollar expressions.  
``` sh
$ mlc "$(npm bin)/node-sass --watch ./scss/ --output ./app/css/"
```

## License
Please see [LICENSE](./LICENSE).
