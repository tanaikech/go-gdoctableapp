// Package gdoctableapp (go-gdoctableapp.go) :
// This is a Golang library for managing tables in Google Document using Google Docs API.
package gdoctableapp

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	docs "google.golang.org/api/docs/v1"
)

// getTables : Retrieve all tables.
func (o *obj) getTables() {
	for i, table := range o.docTables {
		o.docTable = table
		o.parseTable()
		res := [][]string{}
		for _, e := range o.contents {
			temp1 := []string{}
			for _, f := range e {
				temp2 := []string{}
				for _, g := range f.tempColsContent {
					temp2 = append(temp2, strings.Replace(g.content, "\n", "", -1))
				}
				temp1 = append(temp1, strings.Join(temp2, "\n"))
			}
			res = append(res, temp1)
		}
		t := &Table{}
		t.Index = int64(i)
		t.Values = res
		t.TablePosition.StartIndex = table.StartIndex
		t.TablePosition.EndIndex = table.EndIndex
		o.result.Tables = append(o.result.Tables, *t)
	}
}

// getValues : Retrieve values from a table of Document.
func (o *obj) getValues() ([][]string, error) {
	o.parseTable()
	res := [][]string{}
	for _, e := range o.contents {
		temp1 := []string{}
		for _, f := range e {
			temp2 := []string{}
			for _, g := range f.tempColsContent {
				temp2 = append(temp2, strings.Replace(g.content, "\n", "", -1))
			}
			temp1 = append(temp1, strings.Join(temp2, "\n"))
		}
		res = append(res, temp1)
	}
	return res, nil
}

// createDeleteContentRangeRequest : Create DeleteContentRangeRequest.
func createDeleteContentRangeRequest(startIndex, endIndex int64) *docs.Request {
	r := &docs.DeleteContentRangeRequest{}
	r.Range = &docs.Range{
		StartIndex: startIndex,
		EndIndex:   endIndex,
	}
	req := &docs.Request{
		DeleteContentRange: r,
	}
	return req
}

// deleteTable : Delete table.
func (o *obj) deleteTable() error {
	o.parseTable()
	dr := createDeleteContentRangeRequest(o.docTable.StartIndex, o.docTable.EndIndex)
	br := &docs.BatchUpdateDocumentRequest{}
	br.Requests = append(br.Requests, dr)
	o.requestBody = br
	if err := o.documentbatchUpdate(); err != nil {
		return err
	}
	return nil
}

// deleteRowsColumns : Delete rows and columns of a table.
func (o *obj) deleteRowsColumns() error {
	if len(o.params.DeleteRowsColumnsRequest.Rows) == 0 && len(o.params.DeleteRowsColumnsRequest.Columns) == 0 {
		return fmt.Errorf("No parameters for using DeleteRowsAndColumns()")
	}
	rows := o.params.DeleteRowsColumnsRequest.Rows
	cols := o.params.DeleteRowsColumnsRequest.Columns
	sort.Slice(rows, func(i, j int) bool { return rows[i] > rows[j] })
	sort.Slice(cols, func(i, j int) bool { return cols[i] > cols[j] })
	maxDeleteRow := rows[0]
	maxDeleteCol := cols[0]
	table := o.docTable.Table
	if table.Rows < maxDeleteRow || table.Columns < maxDeleteCol {
		return fmt.Errorf("Rows and columns for deleting are outside of the table")
	}
	inputObj := o.params.DeleteRowsColumnsRequest
	l := &docs.Location{}
	l.Index = o.docTable.StartIndex
	br := &docs.BatchUpdateDocumentRequest{}
	if len(inputObj.Rows) > 0 {
		for _, e := range inputObj.Rows {
			tc := &docs.TableCellLocation{}
			tc.TableStartLocation = l
			tc.RowIndex = e
			r := &docs.DeleteTableRowRequest{}
			r.TableCellLocation = tc
			dr := &docs.Request{
				DeleteTableRow: r,
			}
			br.Requests = append(br.Requests, dr)
		}
	}
	if len(inputObj.Columns) > 0 {
		for _, e := range inputObj.Columns {
			tc := &docs.TableCellLocation{}
			tc.TableStartLocation = l
			tc.ColumnIndex = e
			r := &docs.DeleteTableColumnRequest{}
			r.TableCellLocation = tc
			dr := &docs.Request{
				DeleteTableColumn: r,
			}
			br.Requests = append(br.Requests, dr)
		}
	}
	o.requestBody = br
	if err := o.documentbatchUpdate(); err != nil {
		return err
	}
	return nil
}

// setValuesMain : Main method for setValues.
func (o *obj) setValuesMain() error {
	dupChk, err := o.checkDupValues()
	if err != nil {
		return err
	}
	if len(dupChk.dup) > 0 {
		return fmt.Errorf("Range of inputted values are duplicated")
	}
	o.parseInputValuesForSetValues(dupChk)
	o.addRowsAndColumnsForSetValues()
	if o.requestBody != nil {
		if err := o.documentbatchUpdate(); err != nil {
			return err
		}
		if err := o.getTable(); err != nil {
			return err
		}
	}
	o.parseTable()
	o.createSetValuesRequests()
	if err := o.documentbatchUpdate(); err != nil {
		return err
	}
	return nil

}

// setValues : Set values.
func (o *obj) setValues() error {
	if o.params.Works.DoValuesArray {
		vo := &ValueObject{}
		vo.Values = o.params.ValuesArray
		vo.Range.StartRowIndex = 0
		vo.Range.StartColumnIndex = 0
		o.params.ValuesObject = append(o.params.ValuesObject, *vo)
	}
	if err := o.setValuesMain(); err != nil {
		return err
	}
	return nil
}

// createInsertTableRowColumnRequestBody : Create requests body of insertTableRow and insertTableColumn.
func (o *obj) createInsertTableRowColumnRequestBody(maxRow, maxCol int64) *docs.BatchUpdateDocumentRequest {
	startIndex := o.docTable.StartIndex
	tableRow := o.docTable.Table.Rows
	tableCol := o.docTable.Table.Columns
	addRows := maxRow - tableRow
	addColumns := maxCol - tableCol
	br := &docs.BatchUpdateDocumentRequest{}
	if addRows > 0 {
		for i := int64(0); i < addRows; i++ {
			l := &docs.Location{}
			l.Index = startIndex
			tc := &docs.TableCellLocation{}
			tc.RowIndex = tableRow - 1 + i
			tc.TableStartLocation = l
			t := &docs.InsertTableRowRequest{}
			t.InsertBelow = true
			t.TableCellLocation = tc
			dr := &docs.Request{
				InsertTableRow: t,
			}
			br.Requests = append(br.Requests, dr)
		}
	}
	if addColumns > 0 {
		for i := int64(0); i < addColumns; i++ {
			l := &docs.Location{}
			l.Index = startIndex
			tc := &docs.TableCellLocation{}
			tc.ColumnIndex = tableCol - 1 + i
			tc.TableStartLocation = l
			t := &docs.InsertTableColumnRequest{}
			t.InsertRight = true
			t.TableCellLocation = tc
			dr := &docs.Request{
				InsertTableColumn: t,
			}
			br.Requests = append(br.Requests, dr)
		}
	}
	return br
}

// createInsertTableRowColumnRequest : Create requests of insertTableRow and insertTableColumn.
func (o *obj) createInsertTableRowColumnRequest(maxRow, maxCol int64) {
	br := o.createInsertTableRowColumnRequestBody(maxRow, maxCol)
	if len(br.Requests) > 0 {
		o.requestBody = br
	}
}

// createSetValuesRequests : Create the requests for putting values.
func (o *obj) createSetValuesRequests() {
	br := &docs.BatchUpdateDocumentRequest{}
	for i := int64(len(o.parsedValues)) - 1; i >= 0; i-- {
		r := o.parsedValues[i].row
		c := o.parsedValues[i].col
		v := o.parsedValues[i].content
		delReq := o.delCell[r][c]
		if delReq.DeleteContentRange.Range.StartIndex != delReq.DeleteContentRange.Range.EndIndex {
			br.Requests = append(br.Requests, delReq)
		}
		if v != "" {
			t := &docs.InsertTextRequest{}
			location := &docs.Location{}
			location.Index = delReq.DeleteContentRange.Range.StartIndex
			t.Location = location
			t.Text = v
			dr := &docs.Request{
				InsertText: t,
			}
			br.Requests = append(br.Requests, dr)
		}
	}
	o.requestBody = br
}

// createTable : Create new table with values.
func (o *obj) crateTable() error {
	br := &docs.BatchUpdateDocumentRequest{}
	table := &docs.InsertTableRequest{}
	rows := o.params.CreateTableRequest.Rows
	columns := o.params.CreateTableRequest.Columns
	if rows == 0 || columns == 0 {
		return fmt.Errorf("Values of Rows and/or Columns are not found")
	}
	table.Rows = rows
	table.Columns = columns
	var idx int64
	if o.params.CreateTableRequest.Append {
		el := &docs.EndOfSegmentLocation{
			SegmentId: "",
		}
		table.EndOfSegmentLocation = el
		dr1 := &docs.Request{
			InsertTable: table,
		}
		br.Requests = append(br.Requests, dr1)
		o.requestBody = br
		if err := o.documentbatchUpdate(); err != nil {
			return err
		}
		br = &docs.BatchUpdateDocumentRequest{} // Reset
		contents, err := o.getDocument()
		if err != nil {
			return err
		}
		for i := len(contents) - 1; i >= 0; i-- {
			if table := contents[i].Table; table != nil {
				o.docTable = contents[i]
				break
			}
		}
		idx = o.docTable.StartIndex - 1
	} else if o.params.CreateTableRequest.Index != 0 {
		loc := &docs.Location{}
		loc.Index = o.params.CreateTableRequest.Index
		table.Location = loc
		dr1 := &docs.Request{
			InsertTable: table,
		}
		br.Requests = append(br.Requests, dr1)
		idx = o.params.CreateTableRequest.Index
	} else {
		return fmt.Errorf("Please set Index (> 0) or Append")
	}

	if len(o.params.CreateTableRequest.Values) > 0 {
		val, err := parseInputValues(
			o.params.CreateTableRequest.Values,
			idx,
			o.params.CreateTableRequest.Rows,
			o.params.CreateTableRequest.Columns,
		)
		if err != nil {
			return err
		}
		for i := int64(len(val)) - 1; i >= 0; i-- {
			v := val[i].content
			if v != "" {
				t := &docs.InsertTextRequest{
					Location: &docs.Location{
						Index: val[i].index,
					},
					Text: v,
				}
				dr2 := &docs.Request{
					InsertText: t,
				}
				br.Requests = append(br.Requests, dr2)
			}
		}
	}
	if o.params.CreateTableRequest.Index > 0 || len(o.params.CreateTableRequest.Values) > 0 {
		o.requestBody = br
		if err := o.documentbatchUpdate(); err != nil {
			return err
		}
	}
	return nil
}

// appendRow : Append row with values.
func (o *obj) appendRow() error {
	if len(o.params.AppendRowRequest.Values) == 0 {
		return fmt.Errorf("Values for putting are not set")
	}
	vo := &ValueObject{}
	vo.Values = o.params.AppendRowRequest.Values
	vo.Range.StartRowIndex = o.docTable.Table.Rows
	vo.Range.StartColumnIndex = 0
	o.params.ValuesObject = append(o.params.ValuesObject, *vo)
	if err := o.setValuesMain(); err != nil {
		return err
	}
	return nil
}

// convertStr : Convert value to string.
func convertStr(v interface{}) (string, error) {
	switch value := v.(type) {
	case int:
		if g, ok := v.(int); ok {
			return fmt.Sprint(g), nil
		}
		return "", fmt.Errorf("error")
	case int64:
		if g, ok := v.(int64); ok {
			return fmt.Sprint(g), nil
		}
		return "", fmt.Errorf("error")
	case float64:
		if g, ok := v.(float64); ok {
			return fmt.Sprint(g), nil
		}
		return "", fmt.Errorf("error")
	case string:
		if g, ok := v.(string); ok {
			return g, nil
		}
		return "", fmt.Errorf("error")
	default:
		return "", fmt.Errorf("error: Unknown value: %+v, %+v", v, value)
	}
}

// parseInputValues : Parse input values for 2 dimensional array.
func parseInputValues(values [][]interface{}, index, rows, cols int64) ([]tempCheckDupValues, error) {
	index += 4
	v := []tempCheckDupValues{}
	var maxCol int64
	maxRow := int64(len(values))
	for row := int64(0); row < rows; row++ {
		if maxRow > row {
			maxCol = int64(len(values[row]))
		} else {
			maxCol = cols
		}
		for col := int64(0); col < cols; col++ {
			if maxRow > row && maxCol > col && values[row][col] != "" {
				colVal, err := convertStr(values[row][col])
				if err != nil {
					return nil, err
				}
				temp := &tempCheckDupValues{
					row:     int64(row),
					col:     int64(col),
					content: colVal,
					index:   index,
				}
				v = append(v, *temp)
			}
			index += 2
		}
		index++
	}
	return v, nil
}

// addRowsAndColumnsForSetValues : Create requests for adding rows and columns for the inputted values.
func (o *obj) addRowsAndColumnsForSetValues() {
	values := o.params.ValuesObject
	var maxRow, maxCol int64
	for _, e := range values {
		tMaxRow := int64(len(e.Values)) + e.Range.StartRowIndex
		tMaxCol := func() int64 {
			var n int64
			for _, f := range e.Values {
				if n < int64(len(f)) {
					n = int64(len(f))
				}
			}
			return n
		}() + e.Range.StartColumnIndex
		if maxRow < tMaxRow {
			maxRow = tMaxRow
		}
		if maxCol < tMaxCol {
			maxCol = tMaxCol
		}
	}
	o.createInsertTableRowColumnRequest(maxRow, maxCol)
}

// parseInputValuesForSetValues : Sort the inputted values.
func (o *obj) parseInputValuesForSetValues(dupChk *dupCheck) {
	sort.Slice(dupChk.noDup, func(i, j int) bool {
		if dupChk.noDup[i].row < dupChk.noDup[j].row {
			return true
		} else if (dupChk.noDup[i].col < dupChk.noDup[j].col) && (dupChk.noDup[i].row == dupChk.noDup[j].row) {
			return true
		}
		return false
	})
	o.parsedValues = dupChk.noDup
}

// checkDupValues : Check the duplication of values.
func (o *obj) checkDupValues() (*dupCheck, error) {
	values := o.params.ValuesObject
	temp := []tempCheckDupValues{}
	for _, e := range values {
		rowOffset := e.Range.StartRowIndex
		colOffset := e.Range.StartColumnIndex
		for i, row := range e.Values {
			temp2 := []tempCheckDupValues{}
			for j, col := range row {
				colVal, err := convertStr(col)
				if err != nil {
					return nil, err
				}
				t := &tempCheckDupValues{
					row:     int64(i) + rowOffset,
					col:     int64(j) + colOffset,
					content: colVal,
				}
				temp2 = append(temp2, *t)
			}
			temp = append(temp, temp2...)
		}
	}
	dc := &dupCheck{}
	for _, e := range temp {
		c := false
		for _, f := range dc.noDup {
			if f.row == e.row && f.col == e.col {
				c = true
			}
		}
		if c {
			dc.dup = append(dc.dup, e)
		} else {
			dc.noDup = append(dc.noDup, e)
		}
	}
	return dc, nil
}

// parseTable : Parse the retrieved table.
func (o *obj) parseTable() {
	docContent := o.docTable
	tableRows := docContent.Table.TableRows
	var rowsDelCell [][]*docs.Request
	var rowsContents [][]*tempColsContents
	for i := 0; i < len(tableRows); i++ {
		tableCells := tableRows[i].TableCells
		var tRowsDelCell []*docs.Request
		var tRowsContents []*tempColsContents
		for j := 0; j < len(tableCells); j++ {
			tColsContents := &tempColsContents{}
			contents := tableCells[j].Content
			var si int64
			var ei int64
			for k := 0; k < len(contents); k++ {
				if contents[k].Paragraph != nil {
					elements := contents[k].Paragraph.Elements
					for l := 0; l < len(elements); l++ {
						if k == 0 && l == 0 {
							si = elements[l].StartIndex
						}
						if k == len(contents)-1 && l == len(elements)-1 {
							ei = elements[l].EndIndex - 1
						}
						cellContent := ""
						if elements[l].TextRun != nil {
							cellContent = elements[l].TextRun.Content
						} else if elements[l].InlineObjectElement != nil {
							cellContent = "[INLINE OBJECT]"
						} else {
							cellContent = "[UNSUPPORTED CONTENT]"
						}
						tColsContent := &tempColsContent{
							startIndex: elements[l].StartIndex,
							endIndex:   elements[l].EndIndex,
							content:    cellContent, // At Docs API, content is automatically converted to string.
						}
						tColsContents.tempColsContent = append(tColsContents.tempColsContent, *tColsContent)
					}
				} else if contents[k].Table != nil {
					tColsContent := &tempColsContent{
						startIndex: contents[k].StartIndex,
						endIndex:   contents[k].EndIndex,
						content:    "[TABLE]",
					}
					tColsContents.tempColsContent = append(tColsContents.tempColsContent, *tColsContent)
				} else {
					tColsContent := &tempColsContent{
						startIndex: contents[k].StartIndex,
						endIndex:   contents[k].EndIndex,
						content:    "[UNSUPPORTED CONTENT]",
					}
					tColsContents.tempColsContent = append(tColsContents.tempColsContent, *tColsContent)
				}
			}
			tRowsDelCell = append(tRowsDelCell, createDeleteContentRangeRequest(si, ei))
			tRowsContents = append(tRowsContents, tColsContents)
		}
		rowsDelCell = append(rowsDelCell, tRowsDelCell)
		rowsContents = append(rowsContents, tRowsContents)
	}
	o.delCell = rowsDelCell
	o.contents = rowsContents
	o.cell1stIndex = rowsContents[0][0].tempColsContent[0].startIndex
}

// getAllTables : Retrieve all tables from Google Document.
func (o *obj) getAllTables() error {
	contents, err := o.getDocument()
	if err != nil {
		return err
	}
	for _, e := range contents {
		if table := e.Table; table != nil {
			o.docTables = append(o.docTables, e)
		}
	}
	return nil
}

// getTable : Retrieve table from Google Document.
func (o *obj) getTable() error {
	contents, err := o.getDocument()
	if err != nil {
		return err
	}
	c := 0
	for _, e := range contents {
		if table := e.Table; table != nil {
			if o.params.TableIdx == c {
				o.docTable = e
				break
			}
			c++
		}
	}
	if o.docTable == nil {
		return fmt.Errorf("Table of index of %d was not found", o.params.TableIdx)
	}
	return nil
}

// documentbatchUpdate : Request the method of batchUpdate for Google Document.
func (o *obj) documentbatchUpdate() error {
	if o.requestBody != nil {
		doc, err := o.srv.Documents.BatchUpdate(o.params.DocumentID, o.requestBody).Do()
		if err != nil {
			return err
		}
		o.result.ResponseFromAPIs = append(o.result.ResponseFromAPIs, doc)
		o.requestBody = nil
	}
	return nil
}

// getDocument : Retrieve Document object from Google Document.
func (o *obj) getDocument() ([]*docs.StructuralElement, error) {
	doc, err := o.srv.Documents.Get(o.params.DocumentID).Fields("body(content(endIndex,startIndex,table))").Do()
	if err != nil {
		return nil, err
	}
	o.result.ResponseFromAPIs = append(o.result.ResponseFromAPIs, doc)
	return doc.Body.Content, nil
}

// optionChecker : Check inputted options.
func (o *obj) optionChecker() error {
	r := reflect.ValueOf(&o.params.Works).Elem()
	rt := r.Type()
	fl := 0
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		rv := reflect.ValueOf(&o.params.Works)
		value := reflect.Indirect(rv).FieldByName(field.Name)
		if value.Bool() {
			fl++
		}
	}
	if fl == 1 {
		return nil
	}
	return fmt.Errorf("There are many options for methods. Please use one method for one call")
}
