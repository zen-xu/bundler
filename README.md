# Bundler

`Bundler` a tool for bundling resources into a single executable.

## Installation

```bash
go get github.com/zen-xu/bundler
```

## Usage

```bash
bundler config.yaml
```

This will output an executable bundle `config.bundle`

You can just run the bundle

```bash
./config.bundle <options>
```

Or you can unpack it for getting resources

```bash
./config.bundle -u
```

## Configuration Options

```yaml
command: echo hello  # bash command
archive_paths:  # the paths need to be bundled
  - data/
  - csv/*.csv   # support glob
ignore_paths:   # the paths should ignore to be bundled
  - data/secrets
  - csv/demo.csv # support glob
```