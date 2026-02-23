# gitignorer

Generate `.gitignore` files from [github/gitignore](https://github.com/github/gitignore) templates.

## Usage
```sh
gitignorer <template1 template2 ...>   # print to stdout
gitignorer -w <template1 template2 ... >  # write to .gitignore
gitignorer -a <template1 template2 ...>   # append to .gitignore
gitignorer list                        # list available templates
```

## Install
```sh
go install github.com/bjluckow/gitignorer/cmd/gitignorer@latest
```
