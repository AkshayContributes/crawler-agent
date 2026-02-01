package extractor

import (
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ExtractText(body io.Reader) (string, error) {

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "", err
	}

	doc.Find("script, style, nav, footer, header, iframe").Remove()

	doc.Find("div, p, h1, h2, h3, h4, h5, h6, li, br").Each(func(i int, s *goquery.Selection) {
		s.AppendHtml(" ")
	})

	text := doc.Find("body").Text()

	return strings.Join(strings.Fields(text), " "), nil

}
