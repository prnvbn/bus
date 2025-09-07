# bus

A simple CLI to check London bus arrivals from the terminal.

![Screenshot](images/image.png)

### Install

- Using Go: `go install github.com/prnvbn/bus@latest`
- From source: `go build -o bq` then run `./bq ...`

### Notes

- Data sourced from [TfL Unified API](https://api.tfl.gov.uk/), rate limited at 50 requests per minute
