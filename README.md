# dd
Due diligence or dumb decisions?

## DESCRIPTION
* Pull webpage content from news sources on different tickers
* Find something stupid to do with all this data

## TRY IT OUT
```bash
$ cd ~ && git clone https://github.com/jonwho/dd.git
$ cd dd
$ ./scripts/build.sh
$ ./scripts/run.sh
$ go run .
$ cd tmp/ && ls
```

Behold bootleg webscraping is done! By default it will look for 10 sources to write to file.
Get more files by supplying some command line args. When you run `./scripts/run.sh` it'll
login you into bash on the docker container so you can then run `go run . -n 100` to get
100 scraped pages instead.
