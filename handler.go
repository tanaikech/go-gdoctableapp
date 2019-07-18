// Package gdoctableapp (methods.go) :
// This is a Golang library for managing tables in Google Document using Google Docs API.
// This file includes handler method.
package gdoctableapp

func (o *obj) handler() (*Result, error) {
	if err := o.optionChecker(); err != nil {
		return nil, err
	}

	o.result.LibraryVersion = version
	if !o.params.Works.DoCreateTable {
		if o.params.Works.DoGetTables {
			if err := o.getAllTables(); err != nil {
				return nil, err
			}
		} else {
			if err := o.getTable(); err != nil {
				return nil, err
			}
		}
	}

	// getTables
	if o.params.Works.DoGetTables {
		o.getTables()
		if !o.params.ShowAPIResponseFlag {
			o.result.ResponseFromAPIs = nil
		}
		return &o.result, nil
	}

	// getValues
	if o.params.Works.DoGetValues {
		values, err := o.getValues()
		if err != nil {
			return nil, err
		}
		o.result.Values = values
		if !o.params.ShowAPIResponseFlag {
			o.result.ResponseFromAPIs = nil
		}
		return &o.result, nil
	}

	// setValues
	if o.params.Works.DoValuesArray || o.params.Works.DoValuesObject {
		if err := o.setValues(); err != nil {
			return nil, err
		}
		if !o.params.ShowAPIResponseFlag {
			o.result.ResponseFromAPIs = nil
		}
		return &o.result, nil
	}

	// deleteTable
	if o.params.Works.DoDeleteTable {
		if err := o.deleteTable(); err != nil {
			return nil, err
		}
		if !o.params.ShowAPIResponseFlag {
			o.result.ResponseFromAPIs = nil
		}
		return &o.result, nil
	}

	// deleteRowsColumns
	if o.params.Works.DoDeleteRowsColumns {
		if err := o.deleteRowsColumns(); err != nil {
			return nil, err
		}
		if !o.params.ShowAPIResponseFlag {
			o.result.ResponseFromAPIs = nil
		}
		return &o.result, nil
	}

	// createTable
	if o.params.Works.DoCreateTable {
		if err := o.crateTable(); err != nil {
			return nil, err
		}
		if !o.params.ShowAPIResponseFlag {
			o.result.ResponseFromAPIs = nil
		}
		return &o.result, nil
	}

	// appendRow
	if o.params.Works.DoAppendRow {
		if err := o.appendRow(); err != nil {
			return nil, err
		}
		if !o.params.ShowAPIResponseFlag {
			o.result.ResponseFromAPIs = nil
		}
		return &o.result, nil
	}

	return nil, nil
}
