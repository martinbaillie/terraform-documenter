package terraform

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	multierror "github.com/hashicorp/go-multierror"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/segmentio/terraform-docs/doc"
	"github.com/segmentio/terraform-docs/print"
)

const (
	varTypeInput  = "input"
	varTypeOutput = "output"
)

// Document will analyse the Terraform `dir` and return generated Github
// flavoured Markdown documentation for input and out variables.
//
// Document can optionally verify the existence of variable descriptions.
func Document(dir string, verify bool) (string, error) {
	doc, err := getDoc(dir)
	if err != nil {
		return "", err
	}

	if verify {
		if err := verifyAllDescriptions(doc); err != nil {
			return "", err
		}
	}

	return print.Markdown(doc, true)
}

func getDoc(dir string) (*doc.Doc, error) {
	filenames, err := filepath.Glob(path.Join(dir, "*.tf"))
	if err != nil {
		return nil, err
	}

	files := make(map[string]*ast.File, len(filenames))
	for _, filename := range filenames {
		buf, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		f, err := hcl.ParseBytes(buf)
		if err != nil {
			return nil, err
		}

		files[filename] = f
	}

	return doc.Create(files), nil
}

func verifyAllDescriptions(doc *doc.Doc) (err error) {
	for _, i := range doc.Inputs {
		err = verifyDescription(i.Name, i.Description, varTypeInput, err)
	}

	for _, i := range doc.Outputs {
		err = verifyDescription(i.Name, i.Description, varTypeOutput, err)
	}

	if merr, ok := err.(*multierror.Error); ok {
		merr.ErrorFormat = variableErrFormatFunc
	}

	return err
}

func verifyDescription(name, description, varType string, err error) error {
	if description == "" {
		err = multierror.Append(err,
			fmt.Errorf("%s '%s' is missing a description field", varType, name))
	}
	return err
}

// errFormatFunc is a basic formatter that outputs the number of errors
// that occurred along with a bullet point list of the errors.
func variableErrFormatFunc(es []error) string {
	points := make([]string, len(es))
	for i, err := range es {
		points[i] = fmt.Sprintf("  * %s", err)
	}
	return strings.Join(points, "\n")
}
