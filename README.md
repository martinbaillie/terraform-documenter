[![License](https://img.shields.io/badge/license-BSD-brightgreen.svg?style=flat-square)](/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/martinbaillie/terraform-documenter?style=flat-square)](https://goreportcard.com/report/github.com/martinbaillie/terraform-documenter)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/martinbaillie/terraform-documenter)
[![Build](https://img.shields.io/travis/martinbaillie/terraform-documenter/master.svg?style=flat-square)](https://travis-ci.org/martinbaillie/terraform-documenter)
[![Release](https://img.shields.io/github/release/martinbaillie/terraform-documenter.svg?style=flat-square)](https://github.com/martinbaillie/terraform-documenter/releases/latest)

# terraform-documenter

A quick utility for analysing a Terraform directory and outputting Github flavoured Markdown documentation for input and output variables.

The tool can be used to verify the existence of description fields on all input and output variables before outputting (`-verify`).

Additionally, if a source Markdown file is provided (`-source <file>`) then the full file will be formatted before being output, with any existing "Inputs/Outputs" section replaced or otherwise appended.

##### Libraries :clap:
- [Segment.io's Terraform Docs](https://godoc.org/github.com/segmentio/terraform-docs)
- [HCL's AST](https://godoc.org/github.com/hashicorp/hcl/hcl/ast)
- [Black Friday v2](https://godoc.org/gopkg.in/russross/blackfriday.v2)
- [Markdownfmt Renderer](https://godoc.org/github.com/jsternberg/markdownfmt/markdown)

### Usage

```bash
Usage of ./terraform-documenter:
  -dir string
        the Terraform directory to analyse (default ".")
  -source string
        path to optional source Markdown file
  -verify
        verify presence of input/output variable descriptions
```

### Git Hooks
`terraform-documenter` is best used in combination with Git hook scripts to enforce standards across Terraform module repositories. 

The `terraform` command itself has useful in-built linting (`terraform validate`) and formatting (`terraform fmt`).

Terraform repository with hooks and lots of issues:
![](images/issues.png?raw=true)

Terraform repository with hooks and issues gradually fixed:
![](images/fixes.png?raw=true)
