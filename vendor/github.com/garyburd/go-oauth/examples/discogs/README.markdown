This example shows how to use the oauth package with [Discogs](http://www.discogs.com/developers/).

The examples require a configuration file containing a consumer key and secret:

1. [Create an application](https://www.discogs.com/settings/developers).
2. $ cp config.json.example config.json.
3. Edit config.json to include your consumer key and secret from step 1.


To run the web example:

1. $ go run main.go
2. Go to http://127.0.0.1:8080/ in a browser to try the application.
