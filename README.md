# go-gdoctableapp

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENCE)

<a name="top"></a>

# Overview

This is a Golang library for managing tables on Google Document using Google Docs API.

# Description

Google Docs API has been released. When I used this API, I found that it is very difficult for me to manage the tables on Google Document using Google Docs API. Although I checked [the official document](https://developers.google.com/docs/api/how-tos/tables), unfortunately, I thought that it's very difficult for me. So in order to easily manage the tables on Google Document, I created this library.

## Features

- All values can be retrieved from the table on Google Document.
- Values can be put to the table.
- Delete table, rows and columns of the table.
- New table can be created by including values.
- Append rows to the table by including values.
- Replace texts with images.
	- The image data can be retrieved from URL.
	- The image data can be uploaded from the local PC.

## Languages

I manages the tables on Google Document using several languages. So I created the libraries for 4 languages which are golang, node.js and python. Google Apps Script has Class DocumentApp. So I has never created the GAS library yet.

- [go-gdoctableapp](https://github.com/tanaikech/go-gdoctableapp)
	- This library uses [google-api-go-client](https://github.com/googleapis/google-api-go-client).
- [node-gdoctableapp](https://github.com/tanaikech/node-gdoctableapp)
- [gdoctableapppy](https://github.com/tanaikech/gdoctableapppy)
- [google-docs-table-factory](https://github.com/gumatias/google-docs-table-factory) by gumatias

# Method

| Method                                                                       | Explanation                                       |
| :--------------------------------------------------------------------------- | :------------------------------------------------ |
| [`GetTables()`](#gettables)                                                  | Get all tables from Document.                     |
| [`GetValues()`](#getvalues)                                                  | Get values from a table from Document.            |
| [`SetValuesBy2DArray(values [][]interface{})`](#setvaluesby2darray)          | Set values to a table with 2 dimensional array.   |
| [`SetValuesByObject(values []ValueObject)`](#setbaluesbyobject)              | Set values to a table with an object.             |
| [`DeleteTable()`](#deletetable)                                              | Delete a table.                                   |
| [`DeleteRowsAndColumns(d *DeleteRowsColumnsRequest)`](#deleterowsandcolumns) | Delete rows and columns of a table.               |
| [`CreateTable(c *CreateTableRequest)`](#createtable)                         | Create new table including sell values.           |
| [`AppendRow(c *AppendRowRequest)`](#appendrow)                               | Append row to a table by including values.        |
| [`ReplaceTextsToImagesByURL(from, to string)`](#replacetexts)                | Replace texts with images from URL.               |
| [`ReplaceTextsToImagesByFile(from, to string)`](#replacetexts)               | Replace texts with images from files on local PC. |

This library uses [google-api-go-client](https://github.com/googleapis/google-api-go-client).

## Responses

The structure of response from this library is as follows.

```golang
Result struct {
	Tables           []Table       `json:"tables,omitempty"`
	Values           [][]string    `json:"values,omitempty"`
	ResponseFromAPIs []interface{} `json:"responseFromAPIs,omitempty"`
	LibraryVersion   string        `json:"libraryVersion"`
}
```

- When `GetTables()` is used, you can see the values with `Tables`.
- When `GetValues()` is used, you can see the values with `Values`.
- When other methods are used and the option of `ShowAPIResponse` is `true`, you can see the responses from APIs which were used for the method. And also, you can know the number of APIs, which were used for the method, by the length of array of `ResponseFromAPIs`.

# Usage

About the authorization, please check the section of [Authorization](#authorization). In order to use this library, it is required to confirm that [the Quickstart](https://developers.google.com/docs/api/quickstart/go) works fine.

Please import this libray as follows.

```go
gdoctableapp "github.com/tanaikech/go-gdoctableapp"
```

And, please use as follows.

```go
client := ### // Please use your client.
g := gdoctableapp.New()

res, err := g.Docs(documentID).TableIndex(tableIndex).GetValues().Do(client)
```

- You can see the sample script including the authorization script at [here](#authorization)

## Scope

In this library, using the scope of `https://www.googleapis.com/auth/documents` is recommended. When the method of `ReplaceTextsToImagesByFile` is used, also please add `https://www.googleapis.com/auth/drive`.

<a name="gettables"></a>

## 1. GetTables

Get all tables from Document. All values, table index and table position are retrieved.

### Sample script

This sample script retrieves all tables from the Google Document of document ID.

```golang
documentID := "###"
g := gdoctableapp.New()
res, err := g.Docs(documentID).GetTables().Do(client)
if err != nil {
	fmt.Println(err)
	os.Exit(1)
}
fmt.Println(res.Tables) // You can see the retrieved values like this.
```

the structure of `res.Tables` is as follows.

```golang
Table struct {
	Index         int64      `json:"index"` // TableIdx
	Values        [][]string `json:"values"`
	TablePosition struct {
		StartIndex int64 `json:"startIndex"`
		EndIndex   int64 `json:"endIndex"`
	}
}
```

When the option of `ShowAPIResponse` is used, the responses from Docs API can be retrieved. **This option can be used for all methods.**

```golang
documentID := "###"
g := gdoctableapp.New()
res, err := g.Docs(documentID).GetTables().ShowAPIResponse(true).Do(client)
if err != nil {
	fmt.Println(err)
	os.Exit(1)
}
fmt.Println(res.Tables)           // You can see the retrieved values like this.
fmt.Println(res.ResponseFromAPIs) // You can see the responses from Docs API like this.
```

<a name="getvalues"></a>

## 2. GetValues

Get values from the table. All values are retrieved.

### Sample script

This sample script retrieves the values from 1st table in Google Document. You can see the retrieved values as `[][]string`. Because when the values are retrieved by Docs API, all values are automatically converted to the string data.

```golang
documentID := "###"
tableIndex := 0
g := gdoctableapp.New()
res, err := g.Docs(documentID).TableIndex(tableIndex).GetValues().Do(client)
if err != nil {
	fmt.Println(err)
	os.Exit(1)
}
fmt.Println(res.Values) // You can see the retrieved values like this.
```

- `documentID`: Document ID.
- `tableIndex`: Table index. If you want to use the 3rd table in Google Document. It's 2. The start number of index is 0.
- `client`: `*Client` for using Docs API. Please check the section of [Authorization](#authorization).

<a name="setvaluesby2darray"></a>

## 3. SetValuesBy2DArray

Set values to the table with 2 dimensional array. When the rows and columns of values which are put are over those of the table, this method can automatically expand the rows and columns.

### Sample script

This sample script puts the values to the first table in Google Document.

```golang
documentID := "###"
tableIndex := 0
g := gdoctableapp.New()
valuesBy2DArray := [][]interface{}{{"a1", "b1"}, {"a2", "b2"}, {"a3", "b3", "c3"}}
res, err := g.Docs(documentID).TableIndex(tableIndex).SetValuesBy2DArray(valuesBy2DArray).Do(client)
if err != nil {
	fmt.Println(err)
	os.Exit(1)
}
fmt.Println(res)
```

- `documentID`: Document ID.
- `tableIndex`: Table index. If you want to use the 3rd table in Google Document. It's 2. The start number of index is 0.
- `client`: `*Client` for using Docs API. Please check the section of [Authorization](#authorization).
- `valuesBy2DArray`: `[][]interface{}`

### Result

When above script is run, the following result is obtained.

#### From:

![](images/fig1.png)

#### To:

![](images/fig2.png)

<a name="setbaluesbyobject"></a>

## 4. SetValuesByObject

Set values to a table with an object. In this method, you can set the values using the range. When the rows and columns of values which are put are over those of the table, this method can automatically expand the rows and columns.

### Sample script

This script puts the values with the range to the first table in Google Document.

```golang
documentID := "###"
tableIndex := 0
g := gdoctableapp.New()

valuesByObject := []gdoctableapp.ValueObject{}

vo1 := &gdoctableapp.ValueObject{}
vo1.Range.StartRowIndex = 0
vo1.Range.StartColumnIndex = 0
vo1.Values = [][]interface{}{{"A1"}, {"A2", "B2", "c2", "d2"}, {"A3"}}
valuesByObject = append(valuesByObject, *vo1)

vo2 := &gdoctableapp.ValueObject{}
vo2.Range.StartRowIndex = 0
vo2.Range.StartColumnIndex = 1
vo2.Values = [][]interface{}{{"B1", "C1"}}
valuesByObject = append(valuesByObject, *vo2)

res, err := g.Docs(documentID).TableIndex(tableIndex).SetValuesByObject(valuesByObject).Do(client)
if err != nil {
	fmt.Println(err)
	os.Exit(1)
}
fmt.Println(res)
```

- `documentID`: Document ID.
- `tableIndex`: Table index. If you want to use the 3rd table in Google Document. It's 2. The start number of index is 0.
- `client`: `*Client` for using Docs API. Please check the section of [Authorization](#authorization).
- `Range.StartRowIndex` of `valuesByObject`: Row index of `values[0][0]`.
- `Range.StartColumnIndex` of `valuesByObject`: Column index of `values[0][0]`.
- `Values` of `valuesByObject`: Values you want to put.

For example, when the row, column indexes and values are 1, 2 and "value", respectively, "value" is put to "C3".

### Result

When above script is run, the following result is obtained.

#### From:

![](images/fig1.png)

#### To:

![](images/fig3.png)

<a name="deleteuable"></a>

## 5. DeleteTable

### Sample script

This script deletes the first table in Google Document.

```golang
documentID := "###"
tableIndex := 0
g := gdoctableapp.New()

res, err := g.Docs(documentID).TableIndex(tableIndex).DeleteTable().Do(client)
if err != nil {
	fmt.Println(err)
	os.Exit(1)
}
fmt.Println(res)
```

- `documentID`: Document ID.
- `tableIndex`: Table index. If you want to use the 3rd table in Google Document. It's 2. The start number of index is 0.
- `client`: `*Client` for using Docs API. Please check the section of [Authorization](#authorization).

<a name="deleterowsandcolumns"></a>

## 6. DeleteRowsAndColumns

### Sample script

This script deletes rows of indexes of 3, 1 and 2 of the first table in Google Document. And also this script deletes columns of indexes of 2, 1 and 3.

```golang
documentID := "###"
tableIndex := 0
g := gdoctableapp.New()

obj := &gdoctableapp.DeleteRowsColumnsRequest{
	Rows:    []int64{3, 1, 2}, // Start index is 0.
	Columns: []int64{2, 1, 3}, // Start index is 0.
}
res, err := g.Docs(documentID).TableIndex(tableIndex).DeleteRowsAndColumns(obj).Do(client)
if err != nil {
	fmt.Println(err)
	os.Exit(1)
}
fmt.Println(res)
```

- `documentID`: Document ID.
- `tableIndex`: Table index. If you want to use the 3rd table in Google Document. It's 2. The start number of index is 0.
- `client`: `*Client` for using Docs API. Please check the section of [Authorization](#authorization).
- `Rows` of `obj`: Indexes of rows you want to delete.
- `Columns` of `obj`: Indexes of columns you want to delete.

<a name="createtable"></a>

## 7. CreateTable

### Sample script

This script creates new table to the top of Google Document, and the cells of the table have values.

```golang
documentID := "###"
g := gdoctableapp.New()

obj := &gdoctableapp.CreateTableRequest{
	Rows:    3,
	Columns: 5,
	Index:   1,
	// Append:  true, // When this is used instead of "Index", new table is created to the end of Document.
	Values: [][]interface{}{{"a1", "b1"}, {"a2", "b2"}, {"a3", "b3", "c3"}},
}
res, err := g.Docs(documentID).CreateTable(obj).Do(client)
if err != nil {
	fmt.Println(err)
	os.Exit(1)
}
fmt.Println(res)
```

- `documentID`: Document ID.
- `client`: `*Client` for using Docs API. Please check the section of [Authorization](#authorization).
- `Rows` of `obj`: Number of rows of new table.
- `Columns` of `obj`: Number of columns of new table.
- `Index` of `obj`: Index of Document for putting new table. For example, `1` is the top of Document.
- `Append` of `obj`: When `Append` is `true` instead of `Index`, the new table is created to the end of Google Document.
- `Values` of `obj`: If you want to put the values when new table is created, please use this.

### Result

When above script is run, the following result is obtained. In this case, the new table is created to the top of Google Document.

![](images/fig4.png)

<a name="appendrow"></a>

## 8. AppendRow

### Sample script

This sample script appends the values to the first table of Google Document.

```golang
documentID := "###"
tableIndex := 0
g := gdoctableapp.New()

obj := &gdoctableapp.AppendRowRequest{
	Values: [][]interface{}{{"a1", "b1", "c1", 1, "", 2}, {"a2", "b2", "c2", 1, "", 2}},
}
res, err := g.Docs(documentID).TableIndex(tableIndex).AppendRow(obj).Do(client)
if err != nil {
	fmt.Println(err)
	os.Exit(1)
}
fmt.Println(res)
```

- `documentID`: Document ID.
- `tableIndex`: Table index. If you want to use the 3rd table in Google Document. It's 2. The start number of index is 0.
- `client`: `*Client` for using Docs API. Please check the section of [Authorization](#authorization).
- `Values` of `obj`: Values you want to append to the existing table.

### Result

When above script is run, the following result is obtained. In this case, the values are put to the last row. And you can see that 3 columns are automatically added when the script is run.

#### From:

![](images/fig5.png)

#### To:

![](images/fig6.png)

<a name="replacetexts"></a>

## 9. ReplaceTextsToImagesByURL and ReplaceTextsToImagesByFile

### Sample script 1

In this sample, the texts `{{sample}}` in all tables are replaced with the image retrieved by **the URL of `https://###/sample.png`**.

```golang
documentID := "###"
searchText := "{{sample}}"
tableOnly := true
replaceImageURL := "https://###/sample.png"
g := gdoctableapp.New()

res, err := g.Docs(documentID).ReplaceTextsToImagesByURL(searchText, replaceImageURL).TableOnly(tableOnly).Do(client)
if err != nil {
	fmt.Println(err)
	os.Exit(1)
}
fmt.Println(res)
```

### Sample script 2

In this sample, the texts `{{sample}}` in all tables are replaced with the image retrieved by **the file of `./sample.png` on your local PC**.

```golang
documentID := "###"
searchText := "{{sample}}"
tableOnly := true // default is false
replaceImageFilePath := "./sample.png"
g := gdoctableapp.New()

res, err := g.Docs(documentID).ReplaceTextsToImagesByFile(searchText, replaceImageFilePath).TableOnly(tableOnly).Do(client)
if err != nil {
	fmt.Println(err)
	os.Exit(1)
}
fmt.Println(res)
```

- `documentID`: Document ID.
- `client`: `*Client` for using Docs API. Please check the section of [Authorization](#authorization).
- `searchText`: Search text. This text is replaced with image.
- `tableOnly`: When this is `true`, only texts in the table are replaced with image. When this is `false`, the texts in the body are replaced.
- `replaceImageURL`: URL of the image.
- `replaceImageFilePath`: File path of the image.

If you want to change the width and height of the image, please use the method of `SetImageSize(width, height float64)` like below.

```golang
res, err := g.Docs(documentID).SetImageSize(100, 100).ReplaceTextsToImagesByFile(searchText, replaceImageFilePath).TableOnly(tableOnly).Do(client)
```

### Note

- The flow for replacing the text with the image on the local PC.

  1. Upload the image from local PC to Google Drive.
  2. Publicly share the image file. - The time for sharing is several seconds. The file is delete after the image is put.
  3. Put the image using the URL of the publicly shared file.
  4. Delete the image. - Even when the image is delete from Google Drive, the put image on Google Document is not deleted.

- About `SetImageSize`
  > [**objectSize**](https://developers.google.com/docs/api/reference/rest/v1/documents/request#insertinlineimagerequest): The size that the image should appear as in the document. This property is optional and the final size of the image in the document is determined by the following rules: _ If neither width nor height is specified, then a default size of the image is calculated based on its resolution. _ If one dimension is specified then the other dimension is calculated to preserve the aspect ratio of the image. \* If both width and height are specified, the image is scaled to fit within the provided dimensions while maintaining its aspect ratio.

### Result

When above script is run, the following result is obtained.

#### From:

![](images/fig7.png)

#### To:

![](images/fig8.png)

The image of `https://cdn.sstatic.net/Sites/stackoverflow/company/img/logos/so/so-logo.png` was used as the sample image.

When `tableOnly` is `false`, the following result is retrieved.

![](images/fig9.png)

<a name="authorization"></a>

# Authorization

There are 2 patterns for using this library.

## 1. Use OAuth2

Document of OAuth2 is [here](https://developers.google.com/identity/protocols/OAuth2).

### Sample script

In this sample script, the authorization process uses [the Quickstart for Go](https://developers.google.com/docs/api/quickstart/go). You can see the detail information at there.

```golang
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	gdoctableapp "github.com/tanaikech/go-gdoctableapp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
)

// Retrieves a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Requests a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache OAuth token: %v", err)
	}
	json.NewEncoder(f).Encode(token)
}

// OAuth2 : Use OAuth2
func OAuth2() *http.Client {
	credentialFile := "credentials.json"
	b, err := os.ReadFile(credentialFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, docs.DocumentsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)
	return client
}

func main() {
	documentID := "###" // Please set here
	tableIndex := 0     // Please set here

	client := OAuth2()
	g := gdoctableapp.New()

	res, err := g.Docs(documentID).TableIndex(tableIndex).GetValues().Do(client)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(res.Values)
}
```

## 2. Use Service account

Document of Service account is [here](https://developers.google.com/identity/protocols/OAuth2ServiceAccount). When you use Service account, please share Google Document with the email of Service account.

### Sample script

```golang
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	gdoctableapp "github.com/tanaikech/go-gdoctableapp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	docs "google.golang.org/api/docs/v1"
)

// ServiceAccount : Use Service account
func ServiceAccount(credentialFile string) *http.Client {
	b, err := ioutil.ReadFile(credentialFile)
	if err != nil {
		log.Fatal(err)
	}
	var c = struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
	}{}
	json.Unmarshal(b, &c)
	config := &jwt.Config{
		Email:      c.Email,
		PrivateKey: []byte(c.PrivateKey),
		Scopes: []string{
			docs.DocumentsScope,
		},
		TokenURL: google.JWTTokenURL,
	}
	client := config.Client(oauth2.NoContext)
	return client
}

func main() {
	documentID := "###" // Please set here
	tableIndex := 0     // Please set here

	client := ServiceAccount("credential.json") // Please set here
	g := gdoctableapp.New()

	res, err := g.Docs(documentID).TableIndex(tableIndex).GetValues().Do(client)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(res.Values)
}
```

# Sample scripts

- [Creating a Table to Google Document by Retrieving Values from Google Spreadsheet for Golang](https://gist.github.com/tanaikech/0589a673cae9569181def8ccd10793cf)

# Limitations

- In the current stage, unfortunately, `tableCellStyle` cannot be modified by Google Docs API. By this, the formats of cells cannot be modified. About this, I have posted as [Feature Request](https://issuetracker.google.com/issues/135136221).

# References:

- Official document: [Inserting or deleting table rows](https://developers.google.com/docs/api/how-tos/tables#inserting_or_deleting_table_rows)
- If you want to know the relationship between the index and startIndex of each cell, you can see it at [here](https://stackoverflow.com/a/56944149).

---

<a name="licence"></a>

# Licence

[MIT](LICENCE)

<a name="author"></a>

# Author

[Tanaike](https://tanaikech.github.io/about/)

If you have any questions and commissions for me, feel free to tell me.

<a name="updatehistory"></a>

# Update History

- v1.0.0 (July 18, 2019)

  1. Initial release.

- v1.0.5 (January 21, 2020)

  1. When the inline objects and tables are put in the table. An error occurred. This bug was removed by this update.

- v1.1.0 (January 22, 2020)

  1. [2 new methods were added.](#replacetexts) From this version, the texts can be replaced by images. The direct link and local file can be used as the image.

- v1.1.0 (April 18, 2023)

  1. Confirmed all methods and updated the sample scripts. All methods works fine.

[TOP](#top)
