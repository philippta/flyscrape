---
title: 'Retry'
weight: 6
---

The retry feature allows the scraper to automatically retry failed requests. This is particularly useful when dealing with unstable networks or servers that occasionally return error status codes.

The retry feature is automatically enabled and will retry requests that return the following HTTP status codes:

- 403 Forbidden
- 408 Request Timeout
- 425 Too Early
- 429 Too Many Requests
- 500 Internal Server Error
- 502 Bad Gateway
- 503 Service Unavailable
- 504 Gateway Timeout

### Retry Delays

After a failed request, the scraper will wait for a certain amount of time before retrying the request. The delay increases with each consecutive failed attempt, according to the following schedule:

- 1st retry: 1 second delay
- 2nd retry: 2 seconds delay
- 3rd retry: 5 seconds delay
- 4th retry: 10 seconds delay
