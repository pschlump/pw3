pw: parse words
===============
 
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/pschlump/Go-FTL/master/LICENSE)


This is for parsing a string into a  set of words.  For Example:

	import (
		pw "github.com/pschlump/pw3"
	)

	/* ... */

	func ParseLineIntoWords(line string) []string {
		Pw := pw.NewParseWords()
		Pw.SetOptions("C", false, false)
		Pw.SetLine(line)
		rv := Pw.GetWords()
		return rv
	}

Will take a line like this one and return an array of words.

Examples using the above code.

```
	s := `set aa "{\"ab\":123}"`
	words := ParseLineIntoWords(s)
	fmt.Printf ( "%s\n", words )
```

Outputs in (length 3)

```
[set aa {"ab":123}]
```




