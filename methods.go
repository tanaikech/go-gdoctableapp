// Package gdoctableapp (methods.go) :
// This is a Golang library for managing tables in Google Document using Google Docs API.
// This file includes all public methods.
package gdoctableapp

import (
	"net/http"

	docs "google.golang.org/api/docs/v1"
)

///
/// Methods
///

// AppendRow : Append rows and values to existing table.
func (p *Params) AppendRow(c *AppendRowRequest) *Params {
	p.Works.DoAppendRow = true
	p.AppendRowRequest = c
	return p
}

// CreateTable : Create new table with values.
func (p *Params) CreateTable(c *CreateTableRequest) *Params {
	p.Works.DoCreateTable = true
	p.CreateTableRequest = c
	return p
}

// DeleteRowsAndColumns : Delete rows and columns of a table.
func (p *Params) DeleteRowsAndColumns(d *DeleteRowsColumnsRequest) *Params {
	p.Works.DoDeleteRowsColumns = true
	p.DeleteRowsColumnsRequest = d
	return p
}

// DeleteTable : Delete table.
func (p *Params) DeleteTable() *Params {
	p.Works.DoDeleteTable = true
	return p
}

// SetValuesByObject : Put values using object.
func (p *Params) SetValuesByObject(values []ValueObject) *Params {
	p.Works.DoValuesObject = true
	p.ValuesObject = values
	return p
}

// SetValuesBy2DArray : Put values using 2 dimensional array.
func (p *Params) SetValuesBy2DArray(values [][]interface{}) *Params {
	p.Works.DoValuesArray = true
	p.ValuesArray = values
	return p
}

// GetValues : Retrieve values from a table of Google Document.
func (p *Params) GetValues() *Params {
	p.Works.DoGetValues = true
	return p
}

// GetTables : Retrieve all tables from Google Document.
func (p *Params) GetTables() *Params {
	p.Works.DoGetTables = true
	return p
}

///
/// Required parameters
///

// TableIndex : Set table index. If there are 5 tables in Document, tableIndex of 3rd table is 2.
func (p *Params) TableIndex(tableIndex int) *Params {
	p.TableIdx = tableIndex
	return p
}

// Docs : Set Document ID
func (p *Params) Docs(documentID string) *Params {
	p.DocumentID = documentID
	return p
}

// ShowAPIResponse : Show responses from Docs API.
func (p *Params) ShowAPIResponse(f bool) *Params {
	p.ShowAPIResponseFlag = true
	return p
}

// init : Initialize
func (o *obj) init() error {
	srv, err := docs.New(o.params.Client)
	if err != nil {
		return err
	}
	o.srv = srv
	return nil
}

// Do : Retrieve all file list and folder tree under root.
func (p *Params) Do(client *http.Client) (*Result, error) {
	o := &obj{
		params: *p,
	}
	o.params.Client = client
	if err := o.init(); err != nil {
		return nil, err
	}
	res, err := o.handler()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// New : Create an object for using gdoctableapp
func New() *Params {
	p := &Params{}
	return p
}
