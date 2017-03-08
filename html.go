package template

import "path"
import "errors"
import "io/ioutil"
import "github.com/go-hayden-base/fs"
import "os"
import "strings"
import "bytes"

const (
	_TAG_CSS = `{{CSS}}`
	_TAG_JS  = `{{JS}}`
)

type Html struct {
	Css             []string
	Js              []string
	Index           string
	OutputDirectory string

	valueMap map[string]string
}

func (s *Html) Output() error {
	if err := s.outputCss(); err != nil {
		return err
	}
	if err := s.outputJs(); err != nil {
		return err
	}
	if err := s.outputIndex(); err != nil {
		return err
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

func (s *Html) outputIndex() error {
	if !fs.FileExists(s.Index) {
		return errors.New("Can not find index html file in ")
	}
	buffer := new(bytes.Buffer)
	fs.ReadLine(s.Index, func(line string, finished bool, err error, stop *bool) {
		line = s.renderLine(line)
		buffer.WriteString(line + "\n")
	})
	p := path.Join(s.OutputDirectory, "index.html")
	if err := fs.WriteFile(p, buffer.Bytes(), true, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (s *Html) renderLine(line string) string {
	if strings.Index(line, _TAG_CSS) > -1 {
		tags := s.htmlTags(_TAG_CSS)
		return strings.Replace(line, _TAG_CSS, strings.Join(tags, ""), -1)
	}
	if strings.Index(line, _TAG_JS) > -1 {
		tags := s.htmlTags(_TAG_JS)
		return strings.Replace(line, _TAG_JS, strings.Join(tags, ""), -1)
	}
	for key, val := range s.valueMap {
		tag := "{{" + key + "}}"
		if tag == _TAG_CSS || tag == _TAG_JS {
			continue
		}
		if strings.Index(line, tag) < 0 {
			continue
		}
		line = strings.Replace(line, tag, val, -1)
	}
	return line
}

func (s *Html) htmlTags(t string) []string {
	var arr []string
	var dir, prefix, suffix string
	if t == _TAG_CSS {
		arr, dir, prefix, suffix = s.Css, "css", `<link rel="stylesheet" href="`, `" />`
	} else if t == _TAG_JS {
		arr, dir, prefix, suffix = s.Js, "js", `<script type="text/javascript" src="`, `"></script>`
	}
	res := make([]string, 0, len(arr))
	for _, p := range arr {
		name := path.Base(p)
		res = append(res, prefix+path.Join(dir, name)+suffix)
	}
	return res
}

func (s *Html) outputCss() error {
	if len(s.Css) == 0 {
		return nil
	}
	if err := s.check(); err != nil {
		return err
	}
	for _, cssFilePath := range s.Css {
		if b, err := ioutil.ReadFile(cssFilePath); err != nil {
			return err
		} else {
			newCssFilePath := path.Join(s.OutputDirectory, "css", path.Base(cssFilePath))
			if err := fs.WriteFile(newCssFilePath, b, true, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Html) outputJs() error {
	if len(s.Js) == 0 {
		return nil
	}
	if err := s.check(); err != nil {
		return err
	}
	for _, jsFilePath := range s.Js {
		if b, err := ioutil.ReadFile(jsFilePath); err != nil {
			return err
		} else {
			newJsFilePath := path.Join(s.OutputDirectory, "css", path.Base(jsFilePath))
			if err := fs.WriteFile(newJsFilePath, b, true, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Html) check() error {
	if s.OutputDirectory == "" || !path.IsAbs(s.OutputDirectory) {
		return errors.New("Please set output directory")
	}
	return nil
}
