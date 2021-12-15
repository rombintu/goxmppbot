package bot

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/xuri/excelize/v2"
)

func GetXslxFromUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func OpenXslxFile(body []byte) (map[string][]string, error) {
	var buff bytes.Buffer
	buff.Write(body)
	f, err := excelize.OpenReader(&buff)
	if err != nil {
		return map[string][]string{}, err
	}
	rowsMap := make(map[string][]string)

	sheets := f.GetSheetList()
	rows, err := f.Rows(sheets[0])
	if err != nil {
		return map[string][]string{}, err
	}

	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return map[string][]string{}, err
		}
		var cells []string
		rowS := TrimS(strings.Join(row, " "))
		cellsDry := strings.Split(rowS, ":")
		for _, cell := range cellsDry {
			cells = append(cells, TrimS(cell))
		}
		if len(cells) > 1 {
			rowsMap[cells[0]] = cells[1:]
		}
	}
	return rowsMap, nil
}
