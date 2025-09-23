package cutc

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Run(t *testing.T) {
	var options = Args{Delimiter: ",", SkipHeader: false}
	var mock_reader *strings.Reader
	var mock_writer *bufio.Writer
	var writer_buffer bytes.Buffer

	mock_reader = strings.NewReader("column1,column2,column3")
	err := Run(mock_reader, nil, []int{1, 4}, options)
	require.Error(t, err)

	mock_reader = strings.NewReader("column1,column2,column3")
	mock_writer = bufio.NewWriter(&writer_buffer)
	err = Run(mock_reader, mock_writer, []int{1}, options)
	mock_writer.Flush()
	require.NoError(t, err)
	require.Equal(t, "column1\n", writer_buffer.String())

	writer_buffer.Reset()

	mock_reader = strings.NewReader("column1,column2,column3")
	mock_writer = bufio.NewWriter(&writer_buffer)
	err = Run(mock_reader, mock_writer, []int{3, 1}, options)
	mock_writer.Flush()
	require.NoError(t, err)
	require.Equal(t, "column3,column1\n", writer_buffer.String())
}

func Test_ParseColumnsIndexes(t *testing.T) {
	_, err := ParseColumnsIndexes("")
	require.Error(t, err)

	_, err = ParseColumnsIndexes("0")
	require.Error(t, err)

	_, err = ParseColumnsIndexes("a,b,c")
	require.Error(t, err)

	_, err = ParseColumnsIndexes("1,b,c")
	require.Error(t, err)

	// -----------------------

	got, err := ParseColumnsIndexes("1")
	require.NoError(t, err)
	require.Equal(t, []int{1}, got)

	got, err = ParseColumnsIndexes("1,1,1")
	require.NoError(t, err)
	require.Equal(t, []int{1}, got)

	got, err = ParseColumnsIndexes("1,2,3,3")
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3}, got)

	got, err = ParseColumnsIndexes("1,2,3")
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3}, got)

	got, err = ParseColumnsIndexes("9,1,5")
	require.NoError(t, err)
	require.Equal(t, []int{9, 1, 5}, got)
}
