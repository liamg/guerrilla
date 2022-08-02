package output

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/liamg/guerrilla/pkg/guerrilla"
	"github.com/liamg/tml"
	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

const (
	borderTopLeft     = '╭'
	borderTopRight    = '╮'
	borderBottomLeft  = '╰'
	borderBottomRight = '╯'
	borderVertical    = '│'
	borderHorizontal  = '─'
	borderLeftT       = '├'
	borderRightT      = '┤'
)

type Printer interface {
	PrintSummary(address string)
	PrintEmail(email guerrilla.Email)
}

type printer struct {
	w     io.Writer
	width int
}

func New(w io.Writer) Printer {

	width, _, err := term.GetSize(0)
	if err != nil {
		width = 80
	}

	return &printer{
		w:     w,
		width: width,
	}
}

var _ Printer = (*printer)(nil)

func (p *printer) printf(format string, args ...interface{}) {
	_ = tml.Fprintf(p.w, format, args...)
}

func (p *printer) printHeader(heading string) {
	p.printf(
		"\r\n<dim>%c%c%c</dim> <bold>%s</bold> <dim>%c%s%c</dim>\n",
		borderTopLeft,
		borderHorizontal,
		borderRightT,
		heading,
		borderLeftT,
		safeRepeat(string(borderHorizontal), p.width-7-runewidth.StringWidth(heading)),
		borderTopRight,
	)
	p.printBlank()
}

func (p *printer) printDivider(heading string) {
	p.printBlank()
	p.printf(
		"<dim>%c%c%c</dim> <bold>%s</bold> <dim>%c%s%c</dim>\n",
		borderLeftT,
		borderHorizontal,
		borderRightT,
		heading,
		borderLeftT,
		safeRepeat(string(borderHorizontal), p.width-7-runewidth.StringWidth(heading)),
		borderRightT,
	)
	p.printBlank()
}

func safeRepeat(input string, repeat int) string {
	if repeat <= 0 {
		return ""
	}
	return strings.Repeat(input, repeat)
}

func (p *printer) printIn(indent int, strip bool, format string, args ...interface{}) {

	var lines []string
	if strip {
		lines = p.limitSizeWithStrip(fmt.Sprintf(format, args...), p.width-indent-4)
	} else {
		lines = p.limitSize(fmt.Sprintf(format, args...), p.width-indent-4)
	}

	for _, line := range lines {
		realStr := line
		if strip {
			realStr = stripTags(line)
		}
		repeat := p.width - 4 - indent - runewidth.StringWidth(realStr)
		coloured := line
		if strip {
			coloured = tml.Sprintf(line)
		}
		padded := coloured + safeRepeat(" ", repeat)
		p.printf("<dim>%c</dim> %s", borderVertical, safeRepeat(" ", indent))
		p.printf("%s", padded)
		p.printf(" <dim>%c</dim>\n", borderVertical)
	}
}

func (p *printer) limitSizeWithStrip(input string, size int) []string {
	var word string
	var words []string
	var inTag bool
	for _, r := range []rune(input) {
		if inTag {
			word += string(r)
			if r == '>' {
				inTag = false
				if word != "" {
					words = append(words, word)
				}
				word = ""
			}
		} else {
			if r == '<' {
				if word != "" {
					words = append(words, word)
				}
				word = "<"
				inTag = true
			} else if r == ' ' {
				if word != "" {
					words = append(words, word)
					word = ""
				} else {
					words[len(words)-1] += " "
				}
			} else if r == '\n' {
				if word != "" {
					words = append(words, word)
					word = ""
				}
				words = append(words, "\n")
			} else {
				word += string(r)
			}
		}
	}
	if word != "" {
		words = append(words, word)
	}

	var line string
	var currentSize int
	var lines []string
	var hasContent bool

	for _, word := range words {
		if word == "\n" {
			lines = append(lines, line)
			line = ""
			continue
		}
		if word[0] == '<' {
			line += word
			continue
		}
		if currentSize+runewidth.StringWidth(word)+1 > size {
			lines = append(lines, line)
			line = word
			hasContent = true
			currentSize = runewidth.StringWidth(word)
		} else {
			if line != "" && hasContent {
				line += " "
				currentSize++
			}
			line += word
			hasContent = true
			currentSize += runewidth.StringWidth(word)
		}
	}

	if line != "" {
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		return []string{""}
	}

	return lines
}

func (p *printer) limitSize(input string, max int) []string {
	lines := strings.Split(input, "\n")
	var output []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		for runewidth.StringWidth(line) > max {
			output = append(output, line[:max])
			line = line[max:]
		}
		until := max
		if until > len(line) {
			until = len(line)
		}
		output = append(output, line[:until])
	}
	return output
}

func (p *printer) printBlank() {
	p.printIn(0, true, "")
}

func (p *printer) printFooter() {
	p.printBlank()
	p.printf("<dim>%c%s%c</dim>\n", borderBottomLeft, safeRepeat(string(borderHorizontal), p.width-2), borderBottomRight)
}

var rgxHtml = regexp.MustCompile(`<[^>]+>`)

func stripTags(input string) string {
	return rgxHtml.ReplaceAllString(input, "")
}

func (p *printer) PrintSummary(address string) {
	p.printHeader("New Address Created")
	p.printIn(0, true, "Your disposable email address is:")
	p.printIn(0, true, "")
	p.printIn(4, true, "<blue>%s</blue>", address)
	p.printIn(0, true, "")
	p.printIn(0, true, "Emails will appear below as they are received.")
	p.printFooter()
}

func (p *printer) PrintEmail(email guerrilla.Email) {
	p.printHeader("Email #" + email.ID)
	p.printIn(0, true, "Subject:   <blue>%s", email.Subject)
	p.printIn(0, true, "From:      <blue>%s", email.From)
	p.printIn(0, true, "Time:      <blue>%s", email.Timestamp.Format(time.RFC1123))
	p.printDivider("Body")
	p.printIn(0, false, "%s", email.Body)
	p.printFooter()
}
