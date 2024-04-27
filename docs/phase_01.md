# Phase One Documentation

## Introduction

The goal of this phase is to implement a simple URL shortener with two endpoints using the Go programming language, Echo framework, and Logrus for logging. In this phase, we will not use any databases and use an in-memory solution to store data.

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
  "short_url": "http://short.url/abcdef"
  }

### Endpoint: Redirect

- **URL**: `/api/v1/{shortUrl}`
- **Method**: Get
- **Response**: Return longURL for HTTP redirection (301 status code)
  ```json
  {
    "url": "https://example.com/long/url/to/be/shortened"
  }

### Algorithm for Generating Short URLs
