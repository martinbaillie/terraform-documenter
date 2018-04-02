package markdown

import (
	"bytes"

	renderer "github.com/jsternberg/markdownfmt/markdown"
	bf "gopkg.in/russross/blackfriday.v2"
)

const (
	// Extensions for GitHub Flavoured Markdown.
	extensions = bf.NoIntraEmphasis |
		bf.Tables |
		bf.FencedCode |
		bf.Autolink |
		bf.Strikethrough |
		bf.SpaceHeadings |
		bf.NoEmptyLineBeforeBlock |
		bf.CommonExtensions
)

// ReplaceFromHeader will replace the original markdown with the replacement
// markdown after the specified header.
func ReplaceFromHeader(header string, original, replacement []byte) string {
	var (
		r               = renderer.NewRenderer(nil)
		options         = []bf.Option{bf.WithRenderer(r), bf.WithExtensions(extensions)}
		originalRoot    = bf.New(options...).Parse(original)
		replacementRoot = bf.New(options...).Parse(replacement)

		buf bytes.Buffer
	)

	// Cut at the header then append the replacement.
	cutAtHeader(header, originalRoot)
	originalRoot.AppendChild(replacementRoot)

	// Walk the AST rendering modified markdown into a buffer and return.
	originalRoot.Walk(func(current *bf.Node, entering bool) bf.WalkStatus {
		return r.RenderNode(&buf, current, entering)
	})
	return buf.String()
}

// cutAtHeader walks the AST until it finds the header. Once found, it
// unlinks all nodes from then on.
func cutAtHeader(header string, root *bf.Node) {
	var (
		unlinkNodes    []*bf.Node
		unlinkFromHere bool
	)

	root.Walk(func(current *bf.Node, entering bool) bf.WalkStatus {
		if entering && current.Type == bf.Heading {
			if child := current.FirstChild; child != nil {
				if string(child.Literal) == header {
					unlinkFromHere = true
				}
			}
		}

		if unlinkFromHere {
			unlinkNodes = append(unlinkNodes, current)
		}

		return bf.GoToNext
	})

	for _, n := range unlinkNodes {
		n.Unlink()
	}
}
