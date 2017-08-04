# Git Blame

This is a simple Go wrapper around the `git blame` command to provide a simple API for interacting with blame data.

## Install

```shell
go get -u github.com/jhaynie/gitblame
```

## Usage

Call the `Generate` method with the git repository directory, the commit SHA, the filename and a callback implementation.

```golang
callback := func(line BlameLine) error {
	// TODO: do something with the content
	return nil
}
if err := Generate(dir, sha, filename, callback); err != nil {
	// TODO: handle the error
}
```
