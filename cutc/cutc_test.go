package cutc

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Run(t *testing.T) {
	var mock_reader *strings.Reader
	var mock_writer *bufio.Writer
	var writer_buffer bytes.Buffer

	options := Args{FieldsList: "1,4", Delimiter: ",", SkipHeader: false}
	mock_reader = strings.NewReader("column1,column2,column3")
	err := Run(mock_reader, nil, options)
	require.Error(t, err)

	options = Args{FieldsList: "4-1", Delimiter: ",", SkipHeader: false}
	mock_reader = strings.NewReader("column1,column2,column3")
	err = Run(mock_reader, nil, options)
	require.Error(t, err)

	options = Args{FieldsList: "1-8", Delimiter: ",", SkipHeader: false}
	mock_reader = strings.NewReader("column1,column2,column3")
	err = Run(mock_reader, nil, options)
	require.Error(t, err)

	tests := []struct {
		args Args
		data string
		want string
	}{
		{
			args: Args{FieldsList: "1", Delimiter: ",", SkipHeader: false},
			data: "column1,column2,column3",
			want: "column1",
		},
		{
			args: Args{FieldsList: "3,  1", Delimiter: ",", SkipHeader: false},
			data: "column1,column2,column3",
			want: "column3,column1",
		},
		{
			args: Args{FieldsList: "1-3", Delimiter: ",", SkipHeader: false},
			data: "column1,column2,column3",
			want: "column1,column2,column3",
		},
		{
			args: Args{FieldsList: "-3", Delimiter: ",", SkipHeader: false},
			data: "column1,column2,column3",
			want: "column1,column2,column3",
		},
		{
			args: Args{FieldsList: "-1", Delimiter: ",", SkipHeader: false},
			data: "column1,column2,column3",
			want: "column1",
		},
		{
			args: Args{FieldsList: "-2,1,3", Delimiter: ",", SkipHeader: false},
			data: "column1,column2,column3",
			want: "column1,column2,column1,column3",
		},
	}

	for _, tt := range tests {
		mock_reader = strings.NewReader(tt.data)
		mock_writer = bufio.NewWriter(&writer_buffer)
		err = Run(mock_reader, mock_writer, tt.args)
		mock_writer.Flush()
		require.NoError(t, err)
		require.Equal(t, tt.want+"\n", writer_buffer.String())
		writer_buffer.Reset()
	}
}

func Test_ParseFields(t *testing.T) {
	_, err := ParseFields("", 10)
	require.Error(t, err)

	_, err = ParseFields("0", 10)
	require.Error(t, err)

	_, err = ParseFields("1", 0)
	require.Error(t, err)

	_, err = ParseFields("a,b,c", 10)
	require.Error(t, err)

	got, err := ParseFields("1", 10)
	require.NoError(t, err)
	require.Equal(t, []int{1}, got)

	got, err = ParseFields("1,1,1", 10)
	require.NoError(t, err)
	require.Equal(t, []int{1, 1, 1}, got)

	got, err = ParseFields("1,2,3", 10)
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3}, got)

	got, err = ParseFields("1,2,3,3", 10)
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3, 3}, got)

	got, err = ParseFields("9,1,5", 10)
	require.NoError(t, err)
	require.Equal(t, []int{9, 1, 5}, got)

	got, err = ParseFields("1,b,c", 10)
	require.NoError(t, err)
	require.Equal(t, []int{1}, got)

	// Ranges ---------------------
	_, err = ParseFields("6-1", 0)
	require.Error(t, err)

	_, err = ParseFields("0-3", 0)
	require.Error(t, err)

	_, err = ParseFields("1-3", 1)
	require.Error(t, err)

	got, err = ParseFields("1-3", 3)
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3}, got)

	got, err = ParseFields("62-64", 100)
	require.NoError(t, err)
	require.Equal(t, []int{62, 63, 64}, got)

	got, err = ParseFields("2-2", 3)
	require.NoError(t, err)
	require.Equal(t, []int{2}, got)

	got, err = ParseFields("1-", 3)
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3}, got)

	got, err = ParseFields("-3", 3)
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3}, got)

	got, err = ParseFields("1-", 1)
	require.NoError(t, err)
	require.Equal(t, []int{1}, got)

	got, err = ParseFields("-1", 1)
	require.NoError(t, err)
	require.Equal(t, []int{1}, got)

	// Mixed ---------------------
	got, err = ParseFields("1, 2, 3-3, 4", 100)
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3, 4}, got)

	got, err = ParseFields("1, 2, 3, 62-64, -5, 99-, 95", 100)
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3, 62, 63, 64, 1, 2, 3, 4, 5, 99, 100, 95}, got)

}
