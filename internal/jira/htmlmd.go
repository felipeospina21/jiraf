package jira

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// HTMLToMarkdown converts simple HTML (as returned by Jira renderedFields) to markdown.
func HTMLToMarkdown(s string) string {
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return s
	}
	var b strings.Builder
	walkNode(&b, doc, state{})
	return strings.TrimSpace(b.String())
}

type state struct {
	ordered  bool
	listIdx  int
	depth    int
	preBlock bool
}

func walkNode(b *strings.Builder, n *html.Node, s state) {
	if n.Type == html.TextNode {
		text := n.Data
		if !s.preBlock {
			text = strings.ReplaceAll(text, "\n", "")
			text = strings.ReplaceAll(text, "\r", "")
			text = strings.ReplaceAll(text, "\t", "")
		}
		b.WriteString(text)
		return
	}

	if n.Type == html.ElementNode {
		tag := n.DataAtom.String()
		if tag == "" {
			tag = n.Data
		}

		switch tag {
		case "h1", "h2", "h3", "h4", "h5", "h6":
			level := int(tag[1] - '0')
			b.WriteString("\n" + strings.Repeat("#", level) + " ")
			walkChildren(b, n, s)
			b.WriteString("\n\n")
			return
		case "p":
			walkChildren(b, n, s)
			b.WriteString("\n\n")
			return
		case "br":
			b.WriteString("\n")
			return
		case "strong", "b":
			b.WriteString("**")
			walkChildren(b, n, s)
			b.WriteString("**")
			return
		case "em", "i":
			b.WriteString("*")
			walkChildren(b, n, s)
			b.WriteString("*")
			return
		case "code":
			b.WriteString("`")
			walkChildren(b, n, s)
			b.WriteString("`")
			return
		case "pre":
			b.WriteString("\n```\n")
			s.preBlock = true
			walkChildren(b, n, s)
			b.WriteString("\n```\n\n")
			return
		case "a":
			href := attr(n, "href")
			var lb strings.Builder
			walkChildren(&lb, n, s)
			text := lb.String()
			if text == href || strings.TrimRight(text, "/") == strings.TrimRight(href, "/") {
				text = "link"
			}
			b.WriteString(fmt.Sprintf("[%s](%s)", text, href))
			return
		case "ol":
			b.WriteString("\n")
			cs := state{ordered: true, listIdx: 0, depth: s.depth + 1}
			walkChildren(b, n, cs)
			if s.depth == 0 {
				b.WriteString("\n")
			}
			return
		case "ul":
			b.WriteString("\n")
			cs := state{ordered: false, depth: s.depth + 1}
			walkChildren(b, n, cs)
			if s.depth == 0 {
				b.WriteString("\n")
			}
			return
		case "li":
			indent := strings.Repeat("   ", s.depth-1)
			if s.ordered {
				s.listIdx++
				b.WriteString(fmt.Sprintf("%s%d. ", indent, s.listIdx))
			} else {
				b.WriteString(indent + "- ")
			}
			walkChildren(b, n, s)
			b.WriteString("\n")
			return
		case "blockquote":
			var inner strings.Builder
			walkChildren(&inner, n, s)
			for _, line := range strings.Split(strings.TrimSpace(inner.String()), "\n") {
				b.WriteString("> " + line + "\n")
			}
			b.WriteString("\n")
			return
		case "hr":
			b.WriteString("\n---\n\n")
			return
		case "table":
			renderTable(b, n)
			return
		}
	}

	walkChildren(b, n, s)
}

func walkChildren(b *strings.Builder, n *html.Node, s state) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkNode(b, c, s)
	}
}

func attr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func renderTable(b *strings.Builder, table *html.Node) {
	var rows [][]string
	forEachTag(table, "tr", func(tr *html.Node) {
		var cells []string
		for c := tr.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && (c.Data == "td" || c.Data == "th") {
				var cb strings.Builder
				walkChildren(&cb, c, state{})
				cells = append(cells, strings.TrimSpace(cb.String()))
			}
		}
		if len(cells) > 0 {
			rows = append(rows, cells)
		}
	})
	if len(rows) == 0 {
		return
	}
	b.WriteString("\n")
	b.WriteString("| " + strings.Join(rows[0], " | ") + " |\n")
	sep := "|"
	for range rows[0] {
		sep += " --- |"
	}
	b.WriteString(sep + "\n")
	for _, row := range rows[1:] {
		b.WriteString("| " + strings.Join(row, " | ") + " |\n")
	}
	b.WriteString("\n")
}

func forEachTag(n *html.Node, tag string, fn func(*html.Node)) {
	if n.Type == html.ElementNode && n.Data == tag {
		fn(n)
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachTag(c, tag, fn)
	}
}
