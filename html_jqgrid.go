package template

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"strconv"
)

type JQGrid struct {
	Data      string            `json:"datastr,omitempty" bson:"datastr,omitempty"`
	DataType  string            `json:"datatype,omitempty" bson:"datatype,omitempty"`
	ColNames  []string          `json:"colNames,omitempty" bson:"colNames,omitempty"`
	ColModel  []*JQGridColModel `json:"colModel,omitempty" bson:"colModel,omitempty"`
	RowNum    int               `json:"rowNum,omitempty" bson:"rowNum,omitempty"`
	SortName  string            `json:"sortname,omitempty" bson:"sortname,omitempty"`
	SortOrder string            `json:"sortorder,omitempty" bson:"sortorder,omitempty"`
	Caption   string            `json:"caption,omitempty" bson:"caption,omitempty"`

	rowData *JQGridRow
	col     int
	row     int
}

type JQGridColModel struct {
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Index    string `json:"index,omitempty" bson:"index,omitempty"`
	Width    int    `json:"width,omitempty" bson:"width,omitempty"`
	Sortable bool   `json:"sortable" bson:"sortable"`
	Align    string `json:"align,omitempty" bson:"align,omitempty"`
}

type JQGridRow struct {
	Rows []*JQGridRowData `json:"rows,omitempty" bson:"rows,omitempty"`
}

type JQGridRowData struct {
	Id   string   `json:"id,omitempty" bson:"id,omitempty"`
	Cell []string `json:"cell,omitempty" bson:"cell,omitempty"`
}

func NewJQGrid(caption, sortName, sortOrder string, col int) *JQGrid {
	if col < 1 {
		col = 1
	}
	aJQGrid := new(JQGrid)
	aJQGrid.col = col
	aJQGrid.Caption = caption
	aJQGrid.ColNames = make([]string, col, col)
	aJQGrid.ColModel = make([]*JQGridColModel, col, col)
	aJQGrid.DataType = "jsonstring"
	aJQGrid.RowNum = math.MaxInt32
	aJQGrid.SortName = sortName
	aJQGrid.SortOrder = sortOrder
	return aJQGrid
}

func (s *JQGrid) SetColCaption(caps ...string) {
	l := len(s.ColNames)
	if l == 0 {
		return
	}
	for idx, item := range caps {
		if idx < l {
			s.ColNames[idx] = item
		}
	}
}

func (s *JQGrid) SetColModel(col int, name, index, align string, width int, sortable bool) {
	l := len(s.ColModel)
	if l == 0 || col < 0 || col >= l {
		return
	}
	aJQGridColModel := new(JQGridColModel)
	if name == "" {
		aJQGridColModel.Name = "Col " + strconv.Itoa(col)
	} else {
		aJQGridColModel.Name = name
	}
	if index == "" {
		aJQGridColModel.Index = aJQGridColModel.Name
	} else {
		aJQGridColModel.Index = index
	}
	if align == "" {
		aJQGridColModel.Align = "left"
	} else {
		aJQGridColModel.Align = align
	}
	if width < 50 {
		aJQGridColModel.Width = 50
	} else {
		aJQGridColModel.Width = width
	}
	aJQGridColModel.Sortable = sortable
	s.ColModel[col] = aJQGridColModel
}

func (s *JQGrid) AddRowData(rowDatas ...string) {
	l := len(rowDatas)
	if l == 0 {
		return
	}
	if s.rowData == nil {
		s.rowData = new(JQGridRow)
	}
	if s.rowData.Rows == nil {
		s.rowData.Rows = make([]*JQGridRowData, 0, 10)
	}
	s.row++
	aRow := new(JQGridRowData)
	aRow.Id = strconv.Itoa(s.row)
	aRow.Cell = make([]string, s.col, s.col)
	for idx, itme := range rowDatas {
		if idx < s.col {
			aRow.Cell[idx] = itme
		}
	}
	s.rowData.Rows = append(s.rowData.Rows, aRow)
}

func (s *JQGrid) JSONString() (string, error) {
	if s.rowData == nil {
		return "", errors.New("no cell data")
	}
	str, err := s.rowData.JSONString()
	if err != nil {
		return "", err
	}
	s.Data = str

	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s *JQGrid) Output(src, des string) error {
	json, err := s.JSONString()
	if err != nil {
		return err
	}
	aHTML := NewHtml(src, des)
	aHTML.AddValue("Option", json)
	aHTML.AddValue("Title", s.Caption)
	errs := aHTML.Output()
	if errs == nil {
		return nil
	}
	buffer := new(bytes.Buffer)
	for _, aErr := range errs {
		buffer.WriteString(aErr.Error() + "\n")
	}
	return errors.New(buffer.String())
}

func (s *JQGridRow) JSONString() (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
