package pw

import (
	"bytes"
	"fmt"
	"regexp"
	"unicode"
)

const (
	Version = "Version: 1.0.0"
)

type ParseWords struct {
	buf            string // Current string
	bufx           []rune //
	pos            int    // where we are in string
	db             bool   // Debuging Flag
	qf             string // "C"		== '" with \,   "SQL", "...???..."  - style of quoting.		Not really implemented yet.
	keep_quote     bool   // Keep ' and " in output
	keep_backslash bool   // Keep \\ in output
}

func NewParseWords() (pw *ParseWords) {
	return &ParseWords{qf: "C", db: false, keep_quote: false, keep_backslash: false}
}

// SetOptions configures from the default options.
//
// kq - keep qutoes in string.
// kb - keep backslashes
func (this *ParseWords) SetOptions(qf string, kq, kb bool) {
	this.qf = qf
	this.keep_quote = kq
	this.keep_backslash = kb
}

func (this *ParseWords) SetDebug(b bool) {
	this.db = b
}

func (this *ParseWords) AppendLine(s string) {
	this.buf += s
}

func (this *ParseWords) SetLine(s string) {
	this.buf = s
}

/*

From: http://blog.golang.org/strings

	const nihongo = "日本語"
    for i, w := 0, 0; i < len(nihongo); i += w {
        runeValue, width := utf8.DecodeRuneInString(nihongo[i:])
        fmt.Printf("%#U starts at byte position %d\n", runeValue, i)
        w = width
    }
*/

var noneRe *regexp.Regexp = regexp.MustCompile("[ \t\f]+")

func (this *ParseWords) GetWords() []string {
	if this.qf == "none" {
		return noneRe.Split(this.buf[this.pos:], -1)
	}

	i := this.pos
	this.bufx = []rune(this.buf)
	l := len(this.bufx)
	rv := make([]string, 0, 10)
	cs := &bytes.Buffer{} // var cs bytes.Buffer
	var c rune
	wf := false

	var getC = func(this_st int) {
		if i < l {
			// c = this.buf[i : i+1]
			c = this.bufx[i]
			if this.db {
				fmt.Printf("top st=%d c->%c<-\n", this_st, c)
			}
		} else {
			c = rune(0)
		}
	}

x0: // scan across to blank
	getC(0)
	if i >= l {
		goto x99
	} else if c == '"' {
		// this.st = 1
		cs = &bytes.Buffer{}
		if this.keep_quote {
			cs.WriteRune(c)
		}
		wf = true
		i++
		goto x1
	} else if c == '\'' {
		// this.st = 11
		cs = &bytes.Buffer{}
		if this.keep_quote {
			cs.WriteRune(c)
		}
		wf = true
		i++
		goto x11
	} else if unicode.IsSpace(c) {
		// this.st = 2
		if wf {
			rv = append(rv, cs.String())
			cs = &bytes.Buffer{}
			wf = false
		}
		i++
		goto x2
	} else {
		wf = true
		cs.WriteRune(c)
		i++
		goto x0
	}
x1: // Start of "
	getC(1)
	if i >= l {
		goto x99
	} else if c == '\\' {
		// this.st = 3
		if this.keep_backslash {
			cs.WriteRune(c)
		}
		i++
		goto x3
	} else if c == '"' {
		// this.st = 0
		if this.keep_quote {
			cs.WriteRune(c)
		}
		if wf {
			rv = append(rv, cs.String())
			cs = &bytes.Buffer{}
			wf = false
		}
		i++
		goto x0
	} else {
		cs.WriteRune(c)
		i++
		goto x1
	}
x11: // Start of "
	getC(11)
	if i >= l {
		goto x99
	} else if c == '\\' {
		// this.st = 13
		if this.keep_backslash {
			cs.WriteRune(c)
		}
		i++
		goto x13
	} else if c == '\'' {
		// this.st = 0
		if this.keep_quote {
			cs.WriteRune(c)
		}
		if wf {
			rv = append(rv, cs.String())
			cs = &bytes.Buffer{}
			wf = false
		}
		i++
		goto x0
	} else {
		cs.WriteRune(c)
		i++
		goto x11
	}
x2: // Found blank
	// Scan across blanks until non-blank
	getC(2)
	if i >= l {
		goto x99
	} else if unicode.IsSpace(c) {
		i++
		goto x2
	} else if c == '"' {
		wf = true
		// this.st = 4
		wf = true
		cs = &bytes.Buffer{}
		if this.keep_quote {
			cs.WriteRune(c)
		}
		i++
		goto x4
	} else if c == '\'' {
		wf = true
		// this.st = 14
		wf = true
		cs = &bytes.Buffer{}
		if this.keep_quote {
			cs.WriteRune(c)
		}
		i++
		goto x14
	} else {
		wf = true
		// this.st = 0
		cs.WriteRune(c)
		i++
		goto x0
	}
x3: // \" processing
	getC(3)
	if i >= l {
		goto x99
	} else {
		// this.st = 1
		cs.WriteRune(c)
		i++
		goto x1
	}
x13: // \' processing
	getC(13)
	if i >= l {
		goto x99
	} else {
		// this.st = 11
		cs.WriteRune(c)
		i++
		goto x11
	}
x4: // scan across to blank
	getC(4)
	if i >= l {
		goto x99
	} else if c == '"' {
		// this.st = 0
		if this.keep_quote {
			cs.WriteRune(c)
		}
		rv = append(rv, cs.String())
		cs = &bytes.Buffer{}
		wf = false
		i++
		goto x0
	} else if c == '\\' {
		// this.st = 5
		if this.keep_backslash {
			cs.WriteRune(c)
		}
		i++
		goto x5
	} else {
		wf = true
		cs.WriteRune(c)
		i++
		goto x4
	}
x14: // scan across to blank
	getC(14)
	if i >= l {
		goto x99
	} else if c == '\'' {
		// this.st = 0
		if this.keep_quote {
			cs.WriteRune(c)
		}
		rv = append(rv, cs.String())
		cs = &bytes.Buffer{}
		wf = false
		i++
		goto x0
	} else if c == '\\' {
		// this.st = 15
		if this.keep_backslash {
			cs.WriteRune(c)
		}
		i++
		goto x15
	} else {
		wf = true
		cs.WriteRune(c)
		i++
		goto x14
	}
x5: // \" processing
	getC(5)
	if i >= l {
		goto x99
	} else {
		// this.st = 4
		cs.WriteRune(c)
		i++
		goto x4
	}
x15: // \' processing
	getC(15)
	if i >= l {
		goto x99
	} else {
		// this.st = 14
		cs.WriteRune(c)
		i++
		goto x14
	}

x99:
	if wf {
		rv = append(rv, cs.String())
	}
	return rv
}
