# Starting point for a web application wrapper

This is a very basic example of how to wrap an application so you can gradually remove traffic from an existing app, to pull the logic into a new app.

Have fun!

## How to use

You need to have your go path set up, then you can run the following commands:

```
cd httpd
go run main.go --baseURL="http://localhost:8888" --port=9000
```

Pass whatever you are trying to wrap in to the `baseURL` flag, and whatever port you want this new app to run on to the `port` flag.
