package cutc

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

func Run(input io.Reader, output io.Writer, fields []int, options Args) error {
	var delimiter, _ = utf8.DecodeRuneInString(options.Delimiter)

	// Stdin
	var is_header = true
	var csv_fields_len = 0
	var reader = csv.NewReader(input)
	reader.Comma = delimiter
	reader.ReuseRecord = true

	// Stdout
	var writer = csv.NewWriter(output)
	writer.Comma = delimiter
	defer writer.Flush()

	for {
		line, e := reader.Read()
		if e == io.EOF {
			break
		}
		if e != nil {
			continue
		}

		if options.SkipHeader && is_header {
			is_header = false
			continue
		}

		// A quick check if all the passed fields indexes to cut are in range
		if csv_fields_len == 0 {
			csv_fields_len = len(line)
			for _, f := range fields {
				if f > csv_fields_len {
					return errors.New(fmt.Sprintf("there is no field with index %v in the csv data", f))
				}
			}
		}

		var row []string
		for _, f := range fields {
			row = append(row, line[f-1])
		}

		err := writer.Write(row)
		if err != nil {
			return err
		}
	}

	return nil
}

// ParseColumnsIndexes cuts string with fields indexes into a slice
func ParseColumnsIndexes(s string) ([]int, error) {
	result := []int{}
	uniq := map[string]bool{}

	parts := strings.Split(strings.ReplaceAll(s, " ", ""), ",")
	for _, p := range parts {
		if p == "" {
			continue
		}
		if _, ok := uniq[p]; ok {
			continue
		}
		uniq[p] = false

		f, err := strconv.Atoi(p)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("\"%v\" is not a valid field index", p))
		}
		if f < 1 {
			return nil, errors.New("field index can't be less than 1")
		}
		result = append(result, f)
	}

	if len(result) == 0 {
		return nil, errors.New("no fields to cut")
	}

	return result, nil
}
