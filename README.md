# Introduction

This is a really simple go server to take a `.env` file and serve it up in different formats for different use cases.

Currently, this supports the following formats:

- env
- json

## Usage

To use this, create a standard `.env` file and load in any variables you want to use. Then, run the server and make a request to the `/{format}` endpoint with the desired format.

For example, if you have a `.env` file that looks like this:

```
HELLO=world
```

You can make a request to `http://localhost/json` and get the following response:

```json
{
  "HELLO": "world"
}
```

## Running the server

This is set up to either run as a standalone project, or run via the included Dockerfile. To run as a standalone project, you can run the following:

```bash
go run main.go
```

To run via Docker with an `.env` file, you can run the following:

```bash
# Buidl the image
docker build -t config-server .

# Run the image, exposing port 8080
docker run -p 8080:80 --env-file .env config-server
```

