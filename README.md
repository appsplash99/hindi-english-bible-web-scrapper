# hindi-english-bible-web-scrapper

a headless golang web scrapper for bible-data

# Step 1

Generate JSON files for each bible book by scraping data by running _main.go_

# Step 2

The remaining Scripts are to be run on the json files

### Scripts

1. _main.scrap-json.go_

   - populate json files from scrapping bible data
   - Note - this script(in _main.scrap-json.go_) needs to be executed multiple times,
     as there's a limit to the goRoutine or something;
     until the directory _./json-data_ does not contain all 66 bible book json files

2. _main.combine-json.go.txt_ - combines all 66 json files(for individual bible books - genesis/exodus/...) into one file

### docs: https://github.com/go-rod/rod

1. `$ go run . -rod=show` - The _show_ option means "show the browser UI on the foreground".

2. `$ go run . -rod=show,devtools` - also shows devTools

3. `$ go run . -rod="show,slow=1s,trace"` - slow-motion and visual trace
