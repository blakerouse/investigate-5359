# Analyze Duplicate Entries

This quick tool verifies that filebeat will not duplicate logs when there are a 1000 log files each with 1000 lines of content.

## How to use?

First download the version of filebeat you want to test and then place just the binary into this directory next to the `filebeat.yml` file.

Then run the following commands and as long as you get a clean exit code of 0 and no error messages printed then everything worked correctly and no duplicate was created.

```shell
go run main.go generate
./filebeat -e  # ctrl-c once you see a events-{{datetime}}.ndjson file created
go run main.go analyze
./filebeat -e  # run it again ensure it reads registry on start (wait about 30 seconds, then ctrl-c)
go run main.go analyze
```
