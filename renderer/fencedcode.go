package renderer

import (
	"github.com/yuin/goldmark/ast"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// ConfluenceFencedCodeBlockHTMLRender is a renderer.NodeRenderer implementation that
// renders FencedCodeBlock nodes.
type ConfluenceFencedCodeBlockHTMLRender struct {
	html.Config
}

// NewConfluenceFencedCodeBlockHTMLRender returns a new ConfluenceFencedCodeBlockHTMLRender.
func NewConfluenceFencedCodeBlockHTMLRender(opts ...html.Option) renderer.NodeRenderer {
	r := &ConfluenceFencedCodeBlockHTMLRender{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *ConfluenceFencedCodeBlockHTMLRender) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(gast.KindFencedCodeBlock, r.renderConfluenceFencedCode)
}

func (r *ConfluenceFencedCodeBlockHTMLRender) renderConfluenceFencedCode(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	n := node.(*gast.FencedCodeBlock)

	if entering {
		language := n.Language(source)

		s := `<ac:structured-macro ac:name="code" ac:schema-version="1">`
		s = s + `<ac:parameter ac:name="theme">Confluence</ac:parameter>`
		s = s + `<ac:parameter ac:name="linenumbers">true</ac:parameter>`

		if language != nil && isSupported(string(language)) {
			s = s + `<ac:parameter ac:name="language">` + string(language) + `</ac:parameter>`
		}

		s = s + `<ac:plain-text-body><![CDATA[ `
		_, _ = w.WriteString(s)
		r.writeLines(w, source, n)
	} else {
		s := ` ]]></ac:plain-text-body></ac:structured-macro>`
		_, _ = w.WriteString(s)
	}
	return ast.WalkContinue, nil
}

func (r *ConfluenceFencedCodeBlockHTMLRender) writeLines(w util.BufWriter, source []byte, n ast.Node) {
	l := n.Lines().Len()
	for i := 0; i < l; i++ {
		line := n.Lines().At(i)
		w.WriteString(string(line.Value(source)))
	}
}

var supportedLanguages = map[string]struct{}{
	"actionscript": struct{}{},
	"applescript":  struct{}{},
	"bash":         struct{}{},
	"c#":           struct{}{},
	"c++":          struct{}{},
	"css":          struct{}{},
	"coldfusion":   struct{}{},
	"delphi":       struct{}{},
	"diff":         struct{}{},
	"erlang":       struct{}{},
	"groovy":       struct{}{},
	"html":         struct{}{},
	"xml":          struct{}{},
	"java":         struct{}{},
	"java fx":      struct{}{},
	"javascript":   struct{}{},
	"php":          struct{}{},
	"plain text":   struct{}{},
	"powershell":   struct{}{},
	"python":       struct{}{},
	"ruby":         struct{}{},
	"sql":          struct{}{},
	"sass":         struct{}{},
	"scala":        struct{}{},
	"visual basic": struct{}{},
	"yaml":         struct{}{},
}

func isSupported(language string) bool {
	_, ok := supportedLanguages[language]

	return ok
}
