package markdown

import (
	"fmt"
)

type block struct {
	md     string
	Blocks []renderer
}

func (b *block) RenderHTML() string {
	tags := b.tags()
	var content string
	for _, block := range b.Blocks {
		content += block.renderHTML()
	}
	return fmt.Sprintf("%s%s%s", tags[0], content, tags[1])
}

type renderer interface {
	tags() []string
	RenderHTML() string
}

type document struct {
	block
}

func (document) tags() []string {
	return [2]string{"", ""}
}

type header struct {
	block
	level int
}

func (h *header) tags() []string {
	return [2]string{
		fmt.Sprintf("<h%d>", h.level),
		fmt.Sprintf("</h%d>", h.level),
	}
}

type paragraph struct {
	block
}

func (paragraph) tags() []string {
	return [2]string{"<p>", "</p>"}
}

type bold struct {
	block
}

func (bold) tags() []string {
	return [2]string{"<strong>", "</strong>"}
}

type italic struct {
	block
}

func (italic) tags() []string {
	return [2]string{"<em>", "</em>"}
}

type boldAndItalic struct {
	block
}

func (boldAndItalic) tags() []string {
	return [2]string{"<em><strong>", "</strong></em>"}
}

type orderedList struct {
	block
}

func (orderedList) tags() []string {
	return [2]string{"<ol>", "</ol>"}
}

type unorderedList struct {
	block
}

func (unorderedList) tags() []string {
	return [2]string{"<ul>", "</ul>"}
}

type listItem struct {
	block
}

func (listItem) tags() []string {
	return [2]string{"<li>", "</li>"}
}

type inlineCode struct {
	block
}

func (inlineCode) tags() []string {
	return [2]string{"<code>", "</code>"}
}

type codeBlock struct {
	block
}

func (codeBlock) tags() []string {
	return [2]string{"<pre><code>", "</code></pre>"}
}

type link struct {
	block
	href string
}

func (l *link) tags() []string {
	return [2]string{
		fmt.Sprintf("<a href=\"%s\">", l.href),
		"</a>",
	}
}

type image struct {
	src   string
	alt   string
	title string
}

func (i *image) renderHTML() string {
	return fmt.Sprintf(
		"<img src=\"%s\" alt=\"%s\" title=\"%s\">",
		i.src, i.alt, i.title,
	)
}

type plaintext string

func (p *plaintext) RenderHTML() string {
	return p
}

type horizontalRule struct{}

func (horizontalRule) RenderHTML() string {
	return "<hr/>\n"
}
