package cutc

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

func Run(input io.Reader, output io.Writer, options Args) error {
	var (
		delimiter, _ = utf8.DecodeRuneInString(options.Delimiter)
		fields       []int
		err          error
	)

	// Stdin
	var (
		is_header = true
		reader    = csv.NewReader(input)
	)
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

		if is_header {
			// Parse fields when got first line of data
			fields, err = ParseFields(options.FieldsList, len(line))
			if err != nil {
				return err
			}

			is_header = false

			if options.SkipHeader {
				continue
			}
		}

		var row []string
		for _, f := range fields {
			field := line[f-1]

			if options.TrimFields {
				field = strings.TrimSpace(field)
			}

			row = append(row, field)
		}

		err = writer.Write(row)
		if err != nil {
			return err
		}
	}

	return nil
}

// ParseFields cuts string with regular fields indexes
func ParseFields(s string, data_length int) ([]int, error) {
	result := []int{}

	fields_pattern := regexp.MustCompile(`^\d+$`)
	range_fields_pattern := regexp.MustCompile(`^((\d)?-\d+|\d+-(\d+)?)$`)
	range_start_pattern := regexp.MustCompile(`^(\d+)-`)
	range_end_pattern := regexp.MustCompile(`-(\d+)$`)

	parts := strings.Split(strings.ReplaceAll(s, " ", ""), ",")
	for _, p := range parts {
		if p == "" {
			continue
		}

		// Parse regular field (1, 2, 3, ...)
		if fields_pattern.MatchString(p) {
			f, _ := strconv.Atoi(p)
			if f < 1 {
				return nil, errors.New("field index can't be less than 1")
			}
			if f > data_length {
				return nil, fmt.Errorf("there is no field with index %v in the csv data", f)
			}
			result = append(result, f)
		}

		// Parse range fields (1-150, 6-, -102)
		if range_fields_pattern.MatchString(p) {
			start_index := 0
			end_index := 0

			// Is valid range?
			if range_start_pattern.MatchString(p) && range_end_pattern.MatchString(p) {
				start_index, _ = strconv.Atoi(range_start_pattern.FindStringSubmatch(p)[1])
				end_index, _ = strconv.Atoi(range_end_pattern.FindStringSubmatch(p)[1])
			} else if range_start_pattern.MatchString(p) {
				// Only start index?
				start_index, _ = strconv.Atoi(range_start_pattern.FindStringSubmatch(p)[1])
				end_index = data_length
			} else {
				// Only end index?
				start_index = 1
				end_index, _ = strconv.Atoi(range_end_pattern.FindStringSubmatch(p)[1])
			}

			if start_index < 1 {
				return nil, fmt.Errorf("fields range start index can't be less than 1 - %v", p)
			}

			if end_index > data_length {
				return nil, fmt.Errorf("fields range end index can't be greater that the total fields length (%v) in the csv - %v", data_length, p)
			}

			if start_index > end_index {
				return nil, fmt.Errorf("fields range start index is greater that the end index - %v", p)
			}

			for i := start_index; i <= end_index; i++ {
				result = append(result, i)
			}
		}
	}

	if len(result) == 0 {
		return nil, errors.New("no fields to cut")
	}

	return result, nil
}
