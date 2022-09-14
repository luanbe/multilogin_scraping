package util

import (
	"encoding/csv"
	"io"
	"mime/multipart"
)

func ReadCsvFile(file multipart.File) ([]map[string]string, error) {
	var rows []map[string]string
	isFirstRow := true
	headerMap := make(map[string]int)

	// read csv values using csv.Reader
	csvReader := csv.NewReader(file)
	for {
		data := make(map[string]string)
		// read row
		record, err := csvReader.Read()

		// stop at EOF
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			return nil, err
		}

		// Handle first row case
		if isFirstRow {
			isFirstRow = false
			// Add mapping: Column/property name --> record index
			for i, v := range record {
				headerMap[v] = i
			}

			// skip next code
			continue
		}

		for k, i := range headerMap {
			data[k] = record[i]
		}
		rows = append(rows, data)
	}
	return rows, nil
}
