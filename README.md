# Url Shortener

This is the API backend server for a simple URL shortener.

## Usage

To run the server, build the Docker image and run it:

```bash
docker build -t url-shortener .
docker run --rm -p 6666:6666 url-shortener
```

A request has the form:

```json
{
	"stub": "hello",
	"url": "https://some.long/url"
}
```

it contains the url and a preferred stub. the stub field can be omitted.
The Server responds with the same structure. When the stub is already taken or the stub was omitted, the server returns a randomly generated stub.
