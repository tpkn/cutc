<p align="center">
   <img width="300" src="icon.png" alt="" align="center">
</p>
<p align="center">
   It <code>cuts</code>, <code>duplicates</code> and <code>rearranges</code> fields in a csv file.
</p>

## Usage

```
cutc [ -options ] < <file.csv>
```

## Options

```
-d, --delimiter    Fields delimiter (default: comma)
-f, --fields       Fields indexes to cut (starting from 1; could be a ranges; order matters)
-h, --header       Skip csv header
-t, --trim         Remove spaces, tabs and new line chars from the beginnig and end of the each field
--help             Help
--version          Version
```

## Examples

Cut columns 1 and 4

```shell
cutc -f 1,4 < input.csv
```

Cut columns 1, 2 and 3 and skip header row

```shell
cutc -f 1,2,3 -h < input.csv
```

Cut columns 1, 4 and 7, but print them in a specific order - `4,1,7`

```shell
cutc -f 4,1,7 < input.csv
```

Cut and gzip

```shell
cutc -f 1,4,7 < input.csv | gzip -c > output.csv.gz
```

Duplicate fields 1 and 7 multiple times

```shell
cutc -f 4,1,1,7,7 < input.csv
```

Slice and dice... to get fields: `1,2,3,62,63,64,1,2,3,4,5,99,100,95`

```shell
cutc -f 1,2,3,62-64,-5,99-,95 < input.csv
```

