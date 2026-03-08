# gitignorer

A minimal CLI for generating `.gitignore` files from [github/gitignore](https://github.com/github/gitignore) templates. Caches templates locally to minimize GitHub API calls.

## Install
```bash
go install github.com/bjluckow/gitignorer@latest
```

## Subcommands

### `fetch`
Fetch and output one or more gitignore templates.
```bash
gitignorer fetch go python node
```

| Flag | Description |
|------|-------------|
| `-w` | Write output to `./.gitignore` |
| `-a` | Append output to `./.gitignore` |
| `-r` | Refresh cached templates |
```bash
# Write to .gitignore
gitignorer fetch -w go python

# Append to existing .gitignore
gitignorer fetch -a node

# Custom path via shell redirection
gitignorer fetch go > /path/to/.gitignore
```

### `list`
List all available templates.
```bash
gitignorer list
gitignorer list -r   # refresh index
```

### `cache`
Print the location of the local cache file.
```bash
gitignorer cache
```

## License

Apache 2.0