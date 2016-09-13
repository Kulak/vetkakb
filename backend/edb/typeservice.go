package edb

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os/exec"

	"github.com/russross/blackfriday"
)

// RawToTextFunc provides raw data to text or html conversion.
type RawToTextFunc func(raw []byte) (string, error)

// TypeProvider provides a single point of reference
// for all type specific functionality.
type TypeProvider struct {
	TypeNum int
	Name    string
	ToHTML  RawToTextFunc
	ToPlain RawToTextFunc
}

// TypeService maps Entry table's RawType number into
// usable functionality based on that number.
type TypeService struct {
	byType map[int]*TypeProvider
}

// NewTypeService creates new service with no mappings.
func NewTypeService() *TypeService {
	return &TypeService{
		byType: make(map[int]*TypeProvider),
	}
}

// Initialize populates service with default type providers.
func (ts *TypeService) Initialize() {
	ts.AddProvider(plainTextProvider())
	ts.AddProvider(htmlProvider())
	ts.AddProvider(markdownProvider())
	ts.AddProvider(imageProvider())
	ts.AddProvider(binaryProvider())
}

// AddProvider registers new provider.
// If provider exists it does nothing (does not update).
func (ts *TypeService) AddProvider(tp *TypeProvider) {
	_, exists := ts.byType[tp.TypeNum]
	if !exists {
		ts.byType[tp.TypeNum] = tp
	} else {
		log.Fatalf("Type provider number %d already exists", tp.TypeNum)
	}
}

// List returns mapping.
func (ts TypeService) List() map[int]*TypeProvider {
	return ts.byType
}

// NameByNum returns type name or empty if not found.
func (ts TypeService) NameByNum(typeNum int) string {
	tp, exists := ts.byType[typeNum]
	if exists {
		return tp.Name
	}
	return ""
}

// NumByName returns type number from name or zero if not found.
func (ts TypeService) NumByName(name string) int {
	for k, v := range ts.byType {
		if v.Name == name {
			return k
		}
	}
	return 0
}

// ProviderByNum returns error if provider is not registered.
func (ts TypeService) ProviderByNum(typeNum int) (tp *TypeProvider, err error) {
	var exists bool
	tp, exists = ts.byType[typeNum]
	if exists {
		return tp, nil
	}
	return nil, fmt.Errorf("Cannot find type provider for type number %v", typeNum)
}

// ProviderByName returns provider by name.
func (ts TypeService) ProviderByName(name string) (*TypeProvider, error) {
	for _, v := range ts.byType {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, fmt.Errorf("Cannot find type provider for type name %v", name)
}

// plainTextProvider implements simple plain text provider.
// Input is sanitized.
func plainTextProvider() *TypeProvider {
	return &TypeProvider{
		TypeNum: 1,
		Name:    "Plain Text",
		ToHTML: func(raw []byte) (string, error) {
			t, err := template.New("foo").Parse(`{{define "T"}}{{.}}{{end}}`)
			if err != nil {
				return "", err
			}
			sanitized := &bytes.Buffer{}
			err = t.ExecuteTemplate(sanitized, "T", string(raw))
			return fmt.Sprintf("<pre>%s</pre>", string(sanitized.Bytes())), err
		},
		ToPlain: func(raw []byte) (string, error) {
			return string(raw), nil
		},
	}
}

// htmlProvider implements HTML format through pandoc.
// Input is NOT sanitized.
func htmlProvider() *TypeProvider {
	return &TypeProvider{
		TypeNum: 2,
		Name:    "HTML",
		ToHTML: func(raw []byte) (string, error) {
			return string(raw), nil
		},
		ToPlain: func(raw []byte) (string, error) {
			plain, err := pandoc("html", "plain", raw)
			plainStr := string(plain)
			return plainStr, err
		},
	}
}

func markdownProvider() *TypeProvider {
	return &TypeProvider{
		TypeNum: 3,
		Name:    "Markdown",
		ToHTML: func(raw []byte) (string, error) {
			output := blackfriday.MarkdownBasic(raw)
			return string(output), nil
		},
		ToPlain: func(raw []byte) (string, error) {
			return string(raw), nil
		},
	}
}

// provides binary dat maaping
func imageProvider() *TypeProvider {
	return &TypeProvider{
		TypeNum: 4,
		Name:    "Binary/Image",
		ToHTML: func(raw []byte) (string, error) {
			return "", nil
		},
		ToPlain: func(raw []byte) (string, error) {
			return "", nil
		},
	}
}

// provides binary dat maaping
func binaryProvider() *TypeProvider {
	return &TypeProvider{
		TypeNum: 5,
		Name:    "Binary/Custom",
		ToHTML: func(raw []byte) (string, error) {
			return "", nil
		},
		ToPlain: func(raw []byte) (string, error) {
			return "", nil
		},
	}
}

func pandoc(from, to string, stdin []byte) (stdout []byte, err error) {
	cmd := exec.Command("/opt/local/bin/pandoc", "-f", from, "-t", to)
	cmd.Stdin = bytes.NewBuffer(stdin)
	out := &bytes.Buffer{}
	cmd.Stdout = out
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("Run error: %v", err)
		return
	}
	stdout = out.Bytes()
	return
}
