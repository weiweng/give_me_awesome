package helper

import (
	"encoding/csv"
	"io"
	"os"

	"give_me_awesome/logs"
)

func ReadCSV(fileName string) [][]string {
	fs, err := os.Open(fileName)
	if err != nil {
		logs.Fatalf("can not open the file, err is %+v", err)
	}
	defer fs.Close()
	r := csv.NewReader(fs)
	m := make([][]string, 0, 300000)
	for {
		row, err := r.Read()
		if err != nil && err != io.EOF {
			logs.Fatalf("can not read, err is %+v", err)
		}
		if err == io.EOF {
			break
		}
		if len(row) == 0 {
			continue
		}
		m = append(m, row)
	}
	return m
}

func WriteCsv(fileName string, data [][]string) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// 写入UTF-8 BOM
	f.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(f)
	w.WriteAll(data)
	w.Flush()
}
