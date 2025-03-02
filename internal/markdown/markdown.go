package markdown

func MdToHTML(md string) string {

}

type block struct {
	md     string
	Blocks []renderer
}

func (b *block) RenderHTML() string {
	var result string
	result += b.renderOpenTag()
	for _, block := range b.Blocks {
		result += block.renderHTML()
	}
	result += b.renderCloseTag()
	return result
}

type renderer interface {
	renderOpenTag() string
	renderCloseTag() string
	RenderHTML() string
}

type document struct {
	block
}

type header struct {
	block
	level int
}

type paragraph struct {
	block
}

type bold struct {
	block
}

type italic struct {
	block
}

type boldAndItalic struct {
	block
}

type orderedList struct {
	block
}

type unorderedList struct {
	block
}

type listItem struct {
	block
}

// TODO: add remaining types found at https://www.markdownguide.org/basic-syntax
// and implement opening and closing tags for all of them.

type plaintext string

func (p *plaintext) RenderHTML() string {
	return p
}

type horizontalRule struct{}

func (horizontalRule) RenderHTML() string {
	return "<hr/>\n"
}
