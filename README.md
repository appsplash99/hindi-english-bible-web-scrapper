# hindi-english-bible-web-scrapper
a headless golang web scrapper for bible-data


### FYI

the script in _main.go_ needs to be executed multiple times,
as there's a limit to the goRoutine or something;
until the directory _./json-data_ does not contain all 66 bible book json files

### docs: https://github.com/go-rod/rod

1. `$ go run . -rod=show` - The _show_ option means "show the browser UI on the foreground".

2. `$ go run . -rod=show,devtools` - also shows devTools

3. `$ go run . -rod="show,slow=1s,trace"` - slow-motion and visual trace

4.
