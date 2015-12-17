package markdownify

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// MarkdownifyReader takes a io.Reader with HTML and returns the text in Markdown
func MarkdownifyReader(r io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}
	selection := doc.Selection
	return strings.TrimSpace(markdownify(selection)), nil
}

func markdownify(s *goquery.Selection) string {
	var buf bytes.Buffer

	// Slightly optimized vs calling Each: no single selection object created
	for _, n := range s.Nodes {
		buf.WriteString(getNodeText(n))
	}
	return tooManyNLRetexp.ReplaceAllString(strings.TrimFunc(buf.String(), unicode.IsSpace), "\n\n")
}

// Get the specified node's text content.
// BUG: It doesn't respect <pre> tags
func getNodeText(node *html.Node) string {
	var buf bytes.Buffer
	// Clear redundant whitespace from text
	if node.Type == html.TextNode {
		text := normalizeWhitespace(node.Data)
		if node.NextSibling == nil || isBlock(node.NextSibling) {
			text = strings.TrimRightFunc(text, unicode.IsSpace)
		}
		if isBlock(node.NextSibling) {
			text = text + "\n\n"
		}
		if isBlock(node.PrevSibling) {
			text = strings.TrimLeftFunc(text, unicode.IsSpace)
		}
		return text
	}
	// change BRs to spaces unless it has two in which case we add extra
	if node.Data == "br" {
		return "\n"
	}
	if node.FirstChild == nil {
		return ""
	}
	if node.Data == "a" {
		href, exists := getAttributeValue("href", node)
		text := getNodeText(node.FirstChild)
		if !exists {
			return text
		}
		if strings.TrimSpace(text) == "" {
			return " "
		}
		return fmt.Sprintf("[%s](%s)", text, href)
	}
	//buf.WriteString("=> " + node.Data + "|")
	if isHeader(node) {
		buf.WriteString("# ")
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		buf.WriteString(getNodeText(c))
	}
	if isQuote(node) {
		str := strings.TrimSpace(buf.String())
		lines := strings.Split(str, "\n")
		for i, line := range lines {
			txt := strings.TrimFunc(line, unicode.IsSpace)
			if txt != "" {
				lines[i] = "> " + txt
			} else {
				lines[i] = ">" + txt
			}
		}
		return strings.Join(lines, "\n") + "\n\n"
	}
	if isBlock(node) {
		buf.WriteString("\n\n")
	}
	return buf.String()
}

func isBlock(node *html.Node) bool {
	return node != nil && (node.Data == "div" || isParagraph(node) || isHeader(node) || isQuote(node) || isList(node))
}

func isQuote(node *html.Node) bool {
	return node != nil && node.Data == "blockquote"
}

func isList(node *html.Node) bool {
	return node != nil && node.Data == "ul"
}

func isParagraph(node *html.Node) bool {
	return node != nil && node.Data == "p"
}

func isHeader(node *html.Node) bool {
	return node != nil && len(node.Data) == 2 && node.Data[0] == 'h' && node.Data[1] != 'r'
}

// Private function to get the specified attribute's value from a node.
func getAttributeValue(attrName string, n *html.Node) (val string, exists bool) {
	if n == nil {
		return
	}

	for _, a := range n.Attr {
		if a.Key == attrName {
			val = a.Val
			exists = true
			return
		}
	}
	return
}

var (
	spaceRegexp     = regexp.MustCompile("[[:space:]]+")
	tooManyNLRetexp = regexp.MustCompile("\n{3,}")
)

func normalizeWhitespace(str string) string {
	str = spaceRegexp.ReplaceAllString(str, " ")
	return str
}
