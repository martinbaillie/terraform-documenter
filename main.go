package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/martinbaillie/terraform-documenter/markdown"
	"github.com/martinbaillie/terraform-documenter/terraform"
)

const desc = `A quick utility for analysing a Terraform directory and outputting
Github flavoured Markdown documentation for input and output variables.

The tool can be used to verify the existence of description fields on
all input and output variables before outputting (-verify).

Additionally, if a source Markdown file is provided (-source <file>) then
the full file will be formatted before being output, with any existing
"Inputs/Outputs" section replaced or otherwise appended.
`

var (
	dir    string
	source string
	verify bool
)

func init() {
	flag.StringVar(&dir, "dir", ".", "the Terraform directory to analyse")
	flag.StringVar(&source, "source", "", "path to optional source Markdown file")
	flag.BoolVar(&verify, "verify", false, "verify presence of input/output variable descriptions")

	flag.Usage = func() {
		fmt.Printf("%s\nUsage of %s:\n", desc, os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	var err error
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		fatal(err)
	}

	var tfMd []byte
	{
		var tfMdStr string
		if tfMdStr, err = terraform.Document(dir, verify); err != nil {
			fatal(err)
		}

		if tfMdStr == "" {
			os.Exit(0)
		}

		tfMd = []byte(tfMdStr)
	}

	var sourceMd []byte
	{
		if source != "" {
			if _, err = os.Stat(source); os.IsNotExist(err) {
				fatal(err)
			}

			if sourceMd, err = ioutil.ReadFile(source); err != nil {
				fatal(err)
			}
		}
	}

	fmt.Print(markdown.ReplaceFromHeader("Inputs", sourceMd, tfMd))
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
