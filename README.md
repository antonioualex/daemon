# Website URL Daemon
This is a Linux daemon built in Go that exposes a minimal HTTP server to handle incoming requests for adding and
retrieving website URLs. The daemon also periodically retrieves and logs information about the most frequently
requested URLs.

### Endpoints
The server provides the following endpoints:

- A **POST** Method with endpoint: **/url/add**, e.g., `localhost:8000/add/url`. <br>
  This endpoint adds a URL able to receive a request to add a website URL to an internal list of objects.
- A **GET** Method with endpoint **/url/retrieve** , being able to retrieve the latest 50 URLs sent through the previous method
  (sorted from newest to oldest or from the smallest to the biggest, upon user request) and a counter that would show
  how many times that specific URL have been submitted to the API since the program started. In order to get sorted results,
  you can set the query with key **sort** and value one of below cases.
  Possible query values for "sort":
```
"DESC": sort by date (newest to oldest)
"ASC": sort by date (oldest to newest)
"LOWEST": sort by frequency (smallest to largest)
"BIGGEST": sort by frequency (largest to smallest)
```
e.g., `localhost:8000/url/retrieve?sort=ASC`

### Other Features

Every 60 seconds the daemon gets the 10 most submitted/requested URLs from the ones that have been submitted and
measuring insights. All the download operations are happening in parallel with a concurrency factor of three .
All the downloaded times, successful downloads counter and failed downloads counter are logged on the stdout for each
batch.


### HowTo
To build the daemon service, you can either clone this repository and build the application yourself
```
go build ./cmd/main.go
export SERVE_ON_PORT=8000
./main
```
Either build the app with docker
```
docker build -t daemon .
docker run -p 8000:8000 daemon //running on 8000 port
```

###### Optional environment variables
```
SERVE_ON_PORT= Set service port, default value is 8000
```