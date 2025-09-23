package main

const Help = `cutc (v%v) | https://tpkn.me

Yes, it's "cut" for csv data.

Usage:
  cutc [ -options ] < <file.csv>

Options:
  -d, --delimiter    Fields delimiter
  -f, --fields       Fields indexes to cut (starting from 1, order matters)
  -h, --header       Skip csv header
  --help             Help
  --version          Version

Examples:
  # Cut columns 1 and 4
  cutc -f 1,4 < <file.csv>
  
  # Cut columns 1, 4 and 7, but print them in a specific order - 4,1,7
  cutc -f 4,1,7 < <file.csv>
`
