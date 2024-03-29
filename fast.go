package main

import (
	"fmt"
	//"go/scanner"
	"io"
	"os"
	"regexp"
	"bufio"
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

var r = regexp.MustCompile("@")
var regexpAndroid = regexp.MustCompile("Android")
var regexpMSIE = regexp.MustCompile("MSIE")

type User struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"company"`
	Country  string   `json:"country"`
	Email    string   `json:"email"`
	Job      string   `json:"job"`
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
}

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson9e1087fdDecodeGithubComSample(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "company":
			out.Company = string(in.String())
		case "country":
			out.Country = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "job":
			out.Job = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "phone":
			out.Phone = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9e1087fdDecodeGithubComSample(&r, v)
	return r.Error()
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)

	if err != nil {
		panic(err)
	}

	defer func() {
		_ = file.Close()
	}()

	seenBrowsers := make([]string, 0, 10)
	uniqueBrowsers := 0
	user := User{}
	i := -1

	_, _ = fmt.Fprintln(out, "found users:")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i += 1
		// fmt.Printf("%v %v\n", err, line)
		err := user.UnmarshalJSON(scanner.Bytes())
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		browsers := user.Browsers

		for _, browser := range browsers {
			notSeenBefore := false
			if regexpAndroid.MatchString(browser) {
				isAndroid = true
				notSeenBefore = true
			}
			if regexpMSIE.MatchString(browser) {
				isMSIE = true
				notSeenBefore = true
			}
			if notSeenBefore {
				for _, item := range seenBrowsers {
					if item == browser {
						notSeenBefore = false
					}
				}
			}
			if notSeenBefore {
				// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
				seenBrowsers = append(seenBrowsers, browser)
				uniqueBrowsers++
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := r.ReplaceAllString(user.Email, " [at] ")
		_, _ = fmt.Fprintf(out, "[%d] %s <%s>\n", i, user.Name, email)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	_, _ = fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}
