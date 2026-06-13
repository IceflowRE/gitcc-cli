# GitCC CLI

![maintained](https://img.shields.io/badge/maintained-yes-brightgreen.svg)
![Programming Language](https://img.shields.io/badge/language-Go-orange.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/IceflowRE/gitcc/blob/main/LICENSE.md)

GitCC CLI checks commit messages for certain rules.

A native GitHub Actions written in typescript is also available at [IceflowRE/gitcc](https://github.com/IceflowRE/gitcc).

## Installation

Either download the released binary from the [releases](https://github.com/IceflowRE/gitcc-cli/releases) or if you have Go installed, you can run:

```bash
go install github.com/IceflowRE/gitcc-cli/v3/cmd/gitcc@latest
```

## Usage

Most of the time you want to check the latest commit:

```bash
gitcc commit --name regex -o summary="^(?:feat|fix)?: (.+)" -o description=".*"
```

Check a specific message:

```bash
gitcc message --name regex -o summary="^(?:feat|fix)?: (.+)" -o description=".*" "feat: add new feature"
```

Check the history of commits:

```bash
gitcc history --name regex -o summary="^(?:feat|fix)?: (.+)" -o description=".*"
```

if you append `--sha 310f0341a3f70b22527c00009fbe36594c72567d` it will only check the commits until the provided SHA.

For more commands and options, check `gitcc --help`.

## Shipped validators

### RegEx

Accepts two options `-o summary="..."` and `-o description="..."` for validating the summary and description of the commit message. The value is a regular expression that has to match the text.

### SimpleTag

Format: `[<tag>] <Good Description>` (e.g. `[ci] Fix testing suite installation`)

## Custom validators

You can provide your own custom validator written in Go.

You have to implement the [`Validator`](https://github.com/IceflowRE/gitcc-cli/blob/main/gitcc/validator.go#L9) interface and provide a function `NewValidator(options map[string]string) (*Validator, error)` in a `main` package.

You can use the [Regex Validator](https://github.com/IceflowRE/gitcc-cli/blob/main/gitcc/validators/regex/validator.go) as an example and the template below as a starting point.

### Installation

Compile it explicitly

```bash
gitcc validator compile NAME /path/to/validator
```

Compile it on demand, this should be avoided, as the validator will always be compiled if the file was changed, this could be an entrypoint for an attacker if the file is not protected.

```bash
gitcc commit --name NAME --path /path/to/validator --compile -o summary="^(?:feat|fix)?: (.+)" -o description=".*"
```

### Usage

To use a custom validator, either provide the name of the validator it was installed with `--name NAME` or provide the path to the validator with `--path /path/to/validator`. If both are provided `--path` is used.

### Template

```go
package main

import (
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/IceflowRE/gitcc-cli/v3/gitcc"
)

type Validator struct{}

func NewValidator(options map[string]string) (gitcc.Validator, error) {
	return &Validator{}, nil
}

func (v *Validator) Validate(commit *object.Commit) gitcc.Result {
	return gitcc.Result{
		Status:  gitcc.Valid,
		// Messages are ignored for valid results.
		Message: "This is a dummy validator that always returns valid. Please implement your own validator.",
	}
}
```

## pre-commit Config

```yaml
repos:
  - repo: local
    hooks:
      - id: gitcc-commit-msg
        name: GitCC
        entry: gitcc message
        language: system
        args: ["--name", "regex", "-o", "summary=REGEX", "-o", "description=REGEX", "--file"]
        stages: [commit-msg]
        pass_filenames: true
```

## Changelog

### 3.0.0

* First release of GitCC CLI (it is v3 to keep it in sync with the GitHub Action version)

## License

Copyright 2021-present Iceflower S (iceflower@iceflower.eu)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
