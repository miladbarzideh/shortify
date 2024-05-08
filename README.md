## Introduction

The goal of this project is to implement a simple URL shortener with two endpoints using the Go programming language, Echo framework, Gorm, Cobra and Logrus for logging. [WIP]

### Running the Application

1. Build the project
```sh
make build
```
You can run the application using the `serve` command. If the port is not specified, it will be read from the configuration file.
```sh
./shortify serve -p <port>
```

### Migrating the Database

To migrate the database and create all tables if they do not exist, use the migrate command.
```sh
./shortify migrate
```

## Implementation Details

### Endpoint: Create Short URL

- **URL**: `/api/v1/urls/shorten`
- **Method**: POST
- **Request Body**: JSON object with the following structure:
  ```json
  {
    "url": "https://example.com/long/url"
  }

- **Response Body**: Return short url:
  ```json
  {
  "url": "http://short.url/abcdef"
  }

### Endpoint: Redirect

- **URL**: `/api/v1/urls/{shortUrl}`
- **Method**: Get
- **Response**: Return longURL for HTTP redirection (301 status code)

### Algorithm for Generating Short URLs

Short codes are randomly generated Base62 strings, composed of alphanumeric characters. Short code length is configurable. In case of collisions, a retry mechanism generates new codes.

### Configuration
Shortify uses a configuration file (config.yaml) to specify settings such as database connection details. An example configuration file is provided (config.example.yaml).

