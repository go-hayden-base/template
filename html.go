package template

import "path"
import "errors"

import "github.com/go-hayden-base/fs"
import "os"
import "strings"
import "bytes"
import "path/filepath"

type Html struct {
	SourcePath      string
	OutputDirectory string

	valueMap map[string]string
}

func NewHtml(src, des string) *Html {
	aHtml := new(Html)
	aHtml.SourcePath = src
	aHtml.OutputDirectory = des
	aHtml.valueMap = make(map[string]string)
	return aHtml
}

func (s *Html) Output() []error {
	if err := s.check(); err != nil {
		return []error{err}
	}
	if errs := fs.CopyDirectory(s.SourcePath, s.OutputDirectory); errs != nil {
		return errs
	}
	return nil
}

func (s *Html) AddValue(key, val string) {
	if s.valueMap == nil {
		s.valueMap = make(map[string]string)
	}
	s.valueMap[key] = val
}

func (s *Html) RemoveValueForKey(key string) {
	if _, ok := s.valueMap[key]; ok {
		delete(s.valueMap, key)
	}
}

func (s *Html) ValueForKey(key string) (string, bool) {
	val, ok := s.valueMap[key]
	return val, ok
}

func (s *Html) renderIndex() error {
	indexPath := filepath.Join(s.OutputDirectory, "index.html")
	if !fs.FileExists(indexPath) {
		return errors.New("Can not find index html file in ")
	}
	buffer := new(bytes.Buffer)
	fs.ReadLine(indexPath, func(line string, finished bool, err error, stop *bool) {
		line = s.renderLine(line)
		buffer.WriteString(line + "\n")
	})
	if err := fs.WriteFile(indexPath, buffer.Bytes(), true, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (s *Html) renderLine(line string) string {
	for key, val := range s.valueMap {
		tag := "{{" + key + "}}"
		if strings.Index(line, tag) < 0 {
			continue
		}
		line = strings.Replace(line, tag, val, -1)
	}
	return line
}

func (s *Html) check() error {
	if s.OutputDirectory == "" || !filepath.IsAbs(s.OutputDirectory) {
		return errors.New("Please set output directory")
	}
	if s.SourcePath == "" || !filepath.IsAbs(s.SourcePath) {
		return errors.New("Please set source directory")
	}
	indexPath := path.Join(s.SourcePath, "index.html")
	if fs.FileExists(indexPath) {
		return errors.New("Source must contain index.html")
	}
	return nil
}
