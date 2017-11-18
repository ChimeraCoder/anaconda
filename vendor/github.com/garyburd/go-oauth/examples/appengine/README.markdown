This example shows how to use the oauth package on App Engine.

The examples require a configuration file containing a consumer key and secret from Twitter:

1. Register an application at https://dev.twitter.com/apps/new
2. $ cp config.json.example config.json.
3. Edit config.json to include your Twitter consumer key and consumer secret from step 1.

To run the web example:

1. devapp\_server.py .
2. Go to http://127.0.0.1:8080/ in a browser to try the application.
