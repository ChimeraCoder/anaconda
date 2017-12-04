This example shows how to use the oauth package with Twitter from a command line application.

The examples require a configuration file containing a consumer key and secret:

1. Register an application at https://dev.twitter.com/apps/new (Note: create a callback url for the app)
2. $ cp config.json.example config.json.
3. Edit config.json to include your Twitter consumer key and consumer secret from step 1.

To run the command line example with OOB authorization:

1. $ go run main.go
