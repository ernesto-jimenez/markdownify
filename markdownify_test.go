package markdownify

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/require"
)

func TestMarkdownifyReader(t *testing.T) {
	file, err := os.Open("testdata/markdown.html")
	require.NoError(t, err)

	actual, err := ConvertReader(file)
	require.NoError(t, err)

	md, err := ioutil.ReadFile("testdata/markdown_reader.md")
	require.NoError(t, err)
	expected := strings.TrimSpace(string(md))

	assertEqual(t, expected, actual)
}

func TestMarkdownifyEmptyString(t *testing.T) {
	actual, err := ConvertReader(strings.NewReader(""))
	require.NoError(t, err)

	expected := ""

	assertEqual(t, expected, actual)
}

func TestMarkdownifyBrAsLastChild(t *testing.T) {
	str := "<span>content <br /></span>"

	actual, err := ConvertReader(strings.NewReader(str))
	require.NoError(t, err)

	expected := "content"

	assertEqual(t, expected, actual)
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

	assertEqual(t, expected, actual)
}

func assertEqual(t *testing.T, expected, actual string) {
	if actual == expected {
		return
	}
	diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(expected),
		B:        difflib.SplitLines(actual),
		FromFile: "Expected",
		FromDate: "",
		ToFile:   "Actual",
		ToDate:   "",
		Context:  1,
	})
	t.Fatalf("\n\tExpected: %#v\n\tReceived: %#v\n%s", expected, actual, diff)
}
