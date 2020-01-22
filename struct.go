// Package gdoctableapp (methods.go) :
// This is a Golang library for managing tables in Google Document using Google Docs API.
// This file includes struct.
package gdoctableapp

import (
	"net/http"

	docs "google.golang.org/api/docs/v1"
	drive "google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

const (
	version       = "1.1.0"
	defaultFields = "body(content(endIndex,startIndex,table))"
)

type (
	// obj : Main object
	obj struct {
		params Params // Input values
		result Result // Output values

		cell1stIndex int64
		contents     [][]*tempColsContents
		delCell      [][]*docs.Request
		docTable     *docs.StructuralElement
		docTables    []*docs.StructuralElement
		parsedValues []tempCheckDupValues
		requestBody  *docs.BatchUpdateDocumentRequest
		srv          *docs.Service
		srvDrive     *drive.Service
		fields       googleapi.Field
	}

	// Result : Result from gdoctableapp
	Result struct {
		Tables           []Table       `json:"tables,omitempty"`
		Values           [][]string    `json:"values,omitempty"`
		ResponseFromAPIs []interface{} `json:"responseFromAPIs,omitempty"`
		LibraryVersion   string        `json:"libraryVersion"`
		Message          string        `json:"message,omitempty"`
	}

	// Params : Parameters inputted by users.
	Params struct {
		AppendRowRequest         *AppendRowRequest
		Client                   *http.Client `json:"client"`
		CreateTableRequest       *CreateTableRequest
		DeleteRowsColumnsRequest *DeleteRowsColumnsRequest
		DocumentID               string          `json:"documentID"`
		ShowAPIResponseFlag      bool            `json:"showAPIResponseFlag"`
		TableIdx                 int             `json:"tableIdx"`
		ValuesArray              [][]interface{} `json:"valuesArray"`
		ValuesObject             []ValueObject   `json:"valuesObject"`
		ReplaceTextsToImagesP    struct {
			FileID           string  `json:"fileID"`
			ReplaceFromText  string  `json:"replaceFromText"`
			ReplaceToImage   string  `json:"replaceToImage"`
			ReplaceTableOnly bool    `json:"replaceTableOnly"`
			Width            float64 `json:"width"`
			Height           float64 `json:"height"`
		}
		Works struct {
			DoAppendRow                  bool `json:"doAppendRow"`
			DoCreateTable                bool `json:"doCreateTable"`
			DoDeleteTable                bool `json:"doDeleteTable"`
			DoDeleteRowsColumns          bool `json:"doDeleteRowsColumns"`
			DoGetValues                  bool `json:"doGetValues"`
			DoGetTables                  bool `json:"doGetTables"`
			DoValuesArray                bool `json:"doValuesArray"`
			DoValuesObject               bool `json:"doValuesObject"`
			DoReplaceTextsToImagesByURL  bool `json:"doReplaceTextsToImagesByURL"`
			DoReplaceTextsToImagesByFile bool `json:"doReplaceTextsToImagesByFile"`
		}
	}

	// AppendRowRequest : Object for appending row and values to existing table.
	AppendRowRequest struct {
		Values [][]interface{} `json:"values"`
	}

	// CreateTableRequest : Object for creating new table with values.
	CreateTableRequest struct {
		Rows    int64           `json:"rows"`
		Columns int64           `json:"columns"`
		Append  bool            `json:"append"`
		Index   int64           `json:"index"`
		Values  [][]interface{} `json:"values"`
	}

	// DeleteRowsColumnsRequest : Object for deleting rows and columns of a table.
	DeleteRowsColumnsRequest struct {
		Rows    []int64 `json:"deleteRows"`
		Columns []int64 `json:"deleteColumns"`
	}

	// ValueObject : Object for putting values.
	ValueObject struct {
		Range struct {
			StartRowIndex    int64 `json:"startRowIndex"`
			StartColumnIndex int64 `json:"startColumnIndex"`
		} `json:"range"`
		Values [][]interface{} `json:"values"`
	}

	// Table : Retrieved table.
	Table struct {
		Index         int64      `json:"index"`
		Values        [][]string `json:"values"`
		TablePosition struct {
			StartIndex int64 `json:"startIndex"`
			EndIndex   int64 `json:"endIndex"`
		}
	}

	// dupCheck : For cheking duplicated values.
	dupCheck struct {
		dup   []tempCheckDupValues
		noDup []tempCheckDupValues
	}

	// for temporal
	tempColsContents struct {
		tempColsContent []tempColsContent
	}

	// for temporal
	tempColsContent struct {
		startIndex int64
		endIndex   int64
		content    string
	}

	// for temporal
	tempCheckDupValues struct {
		row     int64
		col     int64
		content string
		index   int64
	}
)
