# JSana

*An attempt at static javascript analysis*

### If you appreciate what I do please help

https://commerce.coinbase.com/checkout/1f3b4d8f-afb9-47f4-8a18-1e6f4502ce82

### Installation

1) `git clone https://github.com/iamSm9l/JSana.git`
2) `cd JSana`
3) `chmod +x install.sh`
4) `sudo ./install.sh`

### Usage:

```
JSana -u <FILE>
-u <FILE> : A file of js url's one on each line, (eg output from my tool 'wriggle')
-v : verbose mode, not advisiable unless you love spam
-h : Display this help page
```

### Features

This ***Does not*** spit out money, it is a tool which needs *manual* review after

- Looks for keywords amazon buckets "s3." (for US) and "s3-" for everwhere else
- API keys through "apikey", "api_key", "api key" 
- xss Issues ".innerHTML" ".html(" "eval(" ".dangerouslySetInnerHTML"
- serverside issues "eval("

### Example:

`JSana -u urlFile`

- where `urlFile` is a file with a new JS file on each line 

An example of both files can be seen in the examples folder