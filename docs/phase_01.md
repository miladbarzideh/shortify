# Phase One Documentation

## Introduction

In this phase, we will not use any databases and use an in-memory solution to store data.

## Implementation Details

### Endpoint: Create Short URL

- **URL**: `/api/v1/shorten`
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

- **URL**: `/api/v1/{shortUrl}`
- **Method**: Get
- **Response**: Return longURL for HTTP redirection (301 status code)

### Algorithm for Generating Short URLs

To generate short URLs, we use a unique ID generator to create an ID for each long URL. This ID is then encoded using base64 URL encoding. Two maps are utilized to store long and short URLs, with short URL keys in one map and long URL keys in the other.