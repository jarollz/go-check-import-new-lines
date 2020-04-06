package internal

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	openerRegex, err := regexp.Compile(openerPattern)
	if err != nil {
		t.Fatal(err)
	}
	closerRegex, err := regexp.Compile(closerPattern)
	if err != nil {
		t.Fatal(err)
	}
	lineRegex, err := regexp.Compile(linePattern)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		maxNewLine int32
		filePath   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				maxNewLine: 1,
				filePath:   "/file.go",
			},
			wantErr: false,
		},
		{
			name: "should err with negative < 0 number",
			args: args{
				maxNewLine: -2,
				filePath:   "/file.go",
			},
			wantErr: true,
		},
		{
			name: "should err with empty file path",
			args: args{
				maxNewLine: 1,
				filePath:   "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.maxNewLine, tt.args.filePath)
			if err == nil {
				if tt.wantErr {
					assert.NotNil(t, err)
					return
				}
				assert.Equal(t, tt.args.maxNewLine, got.MaxNewLine)
				assert.Equal(t, tt.args.filePath, got.FilePath)
				assert.Equal(t, openerRegex, got.openerRegex)
				assert.Equal(t, closerRegex, got.closerRegex)
				assert.Equal(t, lineRegex, got.lineRegex)
				assert.NotNil(t, got.readSourceCodeFromFile)
				assert.NotNil(t, got.countImportNewLines)
			} else {
				if !tt.wantErr {
					assert.Nil(t, err)
					return
				}
			}
		})
	}
}

func TestParser_ValidateImportsNewLines(t *testing.T) {
	type fields struct {
		maxNewLine             int32
		readSourceCodeFromFile func() (string, error)
		countImportNewLines    func(src string) int32
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				maxNewLine: 1,
				readSourceCodeFromFile: func() (string, error) {
					return "some source code", nil
				},
				countImportNewLines: func(string) int32 {
					return 1
				},
			},
			wantErr: false,
		},
		{
			name: "read src err",
			fields: fields{
				maxNewLine: 1,
				readSourceCodeFromFile: func() (string, error) {
					return "", errors.New("some io error")
				},
				countImportNewLines: func(string) int32 {
					return 1
				},
			},
			wantErr: true,
		},
		{
			name: "too many new lines",
			fields: fields{
				maxNewLine: 3,
				readSourceCodeFromFile: func() (string, error) {
					return "some code", nil
				},
				countImportNewLines: func(string) int32 {
					return 5
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				MaxNewLine:             tt.fields.maxNewLine,
				readSourceCodeFromFile: tt.fields.readSourceCodeFromFile,
				countImportNewLines:    tt.fields.countImportNewLines,
			}
			if err := p.ValidateImportsNewLines(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateImportsNewLines() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParser_CountImportNewLines(t *testing.T) {
	type args struct {
		sourceCode string
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{
			name: "simple",
			args: args{
				`package main
import (
	"errors"
	"fmt"
)`,
			},
			want: 0,
		},
		{
			name: "goimported",
			args: args{
				`package main
import (
	"errors"
	"fmt"

	"github.com/somebody/package-a"
	"github.com/somebody/package-b"
)`,
			},
			want: 1,
		},
		{
			name: "goimported messy",
			args: args{
				`package main
import (
	"errors"
	"fmt"

	"github.com/somebody/package-a"
	
	"github.com/somebody/package-b"
)`,
			},
			want: 2,
		},
		{
			name: "goimported messy with other import style",
			args: args{
				`package main

import "strconv"
import "strings"
import (
	"errors"
	"fmt"

	"github.com/somebody/package-a"
	
	"github.com/somebody/package-b"
	"github.com/somebody/package-c"

	"github.com/somebody/package-d"
)`,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := New(1, "/file.go")
			if err != nil {
				t.Fatal(err)
			}
			got := p.CountImportNewLines(tt.args.sourceCode)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParser_ReadSourceCodeFromFile(t *testing.T) {
	tests := []struct {
		name           string
		fileName       string
		wantStartsWith string
		wantErr        bool
	}{
		{
			name:           "readme",
			fileName:       "README.md",
			wantStartsWith: "# Check New Line in Imports",
			wantErr:        false,
		},
		{
			name:           "gitignore",
			fileName:       ".gitignore",
			wantStartsWith: "# Binaries for programs and plugins",
			wantErr:        false,
		},
		{
			name:           "no file",
			fileName:       "__________.___",
			wantStartsWith: "",
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goPath := os.Getenv("GOPATH")
			basePath := path.Join(goPath, "src/github.com/jarollz/go-check-import-new-lines/")
			filePath := basePath + "/" + tt.fileName
			p := &Parser{
				FilePath: filePath,
			}
			got, err := p.ReadSourceCodeFromFile()
			if tt.wantErr {
				assert.Equal(t, "", got)
				assert.NotNil(t, err)
				return
			} else {
				assert.Nil(t, err)
			}
			if strings.Index(got, tt.wantStartsWith) != 0 {
				t.Errorf("ReadSourceCodeFromFile, should start with '%s', but got started with '%s'", tt.wantStartsWith, got[:len(tt.wantStartsWith)])
			}
		})
	}
}

func Test_processString(t *testing.T) {
	testProcessor1 := func(in string) string {
		return "12" + in + "3"
	}
	testProcessor2 := func(in string) string {
		return "45" + in + "6"
	}

	type args struct {
		input      string
		processors []stringProcessor
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				input:      "abc",
				processors: []stringProcessor{testProcessor1, testProcessor2},
			},
			want: "4512abc36",
		},
		{
			name: "case 2",
			args: args{
				input:      "abc",
				processors: []stringProcessor{testProcessor1, testProcessor1},
			},
			want: "1212abc33",
		},
		{
			name: "case 3",
			args: args{
				input:      "abc",
				processors: []stringProcessor{testProcessor2, testProcessor2},
			},
			want: "4545abc66",
		},
		{
			name: "case 4",
			args: args{
				input:      "abc",
				processors: []stringProcessor{testProcessor2, testProcessor1},
			},
			want: "1245abc63",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processString(tt.args.input, tt.args.processors); got != tt.want {
				t.Errorf("processString() = %v, want %v", got, tt.want)
			}
		})
	}
}
