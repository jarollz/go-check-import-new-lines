package internal

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

type Parser struct {
	MaxNewLine             int32
	FilePath               string
	openerRegex            *regexp.Regexp
	closerRegex            *regexp.Regexp
	lineRegex              *regexp.Regexp
	readSourceCodeFromFile func() (string, error)
	countImportNewLines    func(src string) int32
}

type stringProcessor func(string) string

const (
	openerPattern = `^import \($`
	closerPattern = `^\)$`
	linePattern   = `^\s*$`
)

func New(maxNewLine int32, filePath string) (*Parser, error) {
	if maxNewLine < 0 {
		return nil, errors.New("maxNewLine must be a positive integer")
	}
	if filePath == "" {
		return nil, errors.New("filePath is required")
	}

	p := &Parser{
		MaxNewLine: maxNewLine,
		FilePath:   filePath,
	}

	p.openerRegex, _ = regexp.Compile(openerPattern)
	p.closerRegex, _ = regexp.Compile(closerPattern)
	p.lineRegex, _ = regexp.Compile(linePattern)
	p.readSourceCodeFromFile = p.ReadSourceCodeFromFile
	p.countImportNewLines = p.CountImportNewLines

	return p, nil
}

func (p *Parser) ValidateImportsNewLines() error {
	src, err := p.readSourceCodeFromFile()
	if err != nil {
		return err
	}
	srcImportsLineCount := p.countImportNewLines(src)
	if srcImportsLineCount > p.MaxNewLine {
		return errors.New(fmt.Sprintf("extra new lines in imports, got %d messy new lines (max %d)", srcImportsLineCount, p.MaxNewLine))
	}
	return nil
}

func (p *Parser) CountImportNewLines(sourceCode string) int32 {
	cleanedSourceCode := processString(sourceCode, []stringProcessor{
		func(in string) string {
			return strings.Replace(in, "\r\n", "\n", -1)
		},
		func(in string) string {
			return strings.Replace(in, "\r", "\n", -1)
		},
	})
	sourceCodeLines := strings.Split(cleanedSourceCode, "\n")

	lineCount := int32(0)
	importOpened := false
	for _, line := range sourceCodeLines {
		if importOpened {
			if p.lineRegex.MatchString(line) {
				lineCount++
				continue
			} else if p.closerRegex.MatchString(line) {
				importOpened = false
				continue
			}
		} else {
			if p.openerRegex.MatchString(line) {
				importOpened = true
				continue
			}
		}
	}

	return lineCount
}

func (p *Parser) ReadSourceCodeFromFile() (string, error) {
	srcBytes, err := ioutil.ReadFile(p.FilePath)
	if err != nil {
		return "", err
	}
	src := string(srcBytes)
	return src, nil
}

func processString(input string, processors []stringProcessor) string {
	output := input
	for _, p := range processors {
		output = p(output)
	}
	return output
}
