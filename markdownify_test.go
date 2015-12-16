package markdownify

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
)

func TestMarkdownifyReader(t *testing.T) {
	file, err := os.Open("testdata/markdown.html")
	require.NoError(t, err)

	actual, err := MarkdownifyReader(file)
	require.NoError(t, err)

	md, err := ioutil.ReadFile("testdata/markdown_reader.md")
	require.NoError(t, err)
	expected := strings.TrimSpace(string(md))

	if actual != expected {
		t.Fatalf("Expected: %#v\n\tGot: %#v", expected, actual)
	}
}

func TestMarkdownifyEmptyString(t *testing.T) {
	actual, err := MarkdownifyReader(strings.NewReader(""))
	require.NoError(t, err)

	expected := ""

	if actual != expected {
		t.Fatalf("Expected: %#v\n\tGot: %#v", expected, actual)
	}
}

func TestMarkdownifyBrAsLastChild(t *testing.T) {
	str := "<span>content <br /></span>"

	actual, err := MarkdownifyReader(strings.NewReader(str))
	require.NoError(t, err)

	expected := "content"

	if actual != expected {
		t.Fatalf("Expected: %#v\n\tGot: %#v", expected, actual)
	}
}

func TestMarkdownConvert(t *testing.T) {
	file, err := os.Open("testdata/markdown.html")
	require.NoError(t, err)
	doc, err := goquery.NewDocumentFromReader(file)
	require.NoError(t, err)

	selection := doc.Find("#content")
	actual := strings.TrimSpace(markdownify(selection))

	md, err := ioutil.ReadFile("testdata/markdown.md")
	require.NoError(t, err)
	expected := strings.TrimSpace(string(md))

	if actual != expected {
		t.Fatalf("Expected: %#v\n\tGot: %#v", expected, actual)
	}
}
