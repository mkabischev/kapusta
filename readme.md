# Kapusta

[![Build Status](https://travis-ci.org/mkabischev/kapusta.svg?branch=master)](https://travis-ci.org/mkabischev/kapusta)

It`s middleware approach for using http.Client. You can wrap your client with different functionality: 

 - log every request
 - append auth headers
 - use http cache
 - use etcd/consul for service discovery
 - and whatever you want!
 
**Just like a cabbage!**

![](http://2.bp.blogspot.com/-LtmW_ktxtXU/Um28ElCtQnI/AAAAAAAAB04/aaXWbmTdbnE/s1600/cabbage.png)

## Client interface

Internal http package doesn`t have any interface for http clients, so Kapusta provides very simple client interface:
```go
type Client interface {
	Do(ctx context.Context, *http.Request) (*http.Response, error)
}
```

## Middlewares:

```go
type MiddlewareFunc func(IClient) IClient
```

Kapusta provides some helpful middlewares for you:

- ```HeadersMiddleware(values map[string]string)``` Adds headers to requests
- ```HeaderMiddleware(name, value string)``` Like headers, but add only one header. 
- ```RecoverMiddleware()``` Converts all panics into errors
- ```BaseURLMiddleware(baseURL string)``` Replaces scheme and host to baseURL value.

## Usage

```go
client := http.DefaultClient

decoratedClient := kapusta.Chain(
    client,
    middleware.HeaderMiddleware("X-Auth", "123"),
    middleware.RecoverMiddleware(), // better to place it last to recover panics from middlewares too
)
```

## Create your own middleware

There are two ways of creating new middleware.

You can create some new struct:
```go
struct AwesomeStuffClient {
    client kapusta.Client
}

func(c *AwesomeStuffClient) Do(ctx context.Context, r *http.Request) (*http.Response, error) {
    // some stuff before call
    res, err := c.client.Do(ctx, r)
    // some stuff after call
    
    return res, err
}

func AwesomeStuffDecorator(c kapusta.Client) kapusta.Client {
    return &AwesomeStuffClient{client: c}
}
```

Or you can create just a function with type:
```go 
type ClientFunc func(ctx context.Context, *http.Request) (*http.Response, error)
```

So the same example will be looks like:
```go
func AwesomeStuffDecorator(c kapusta.Client) kapusta.Client {
	return kapusta.ClientFunc(func(ctx context.Context, r *http.Request) (*http.Response, error) {
		// some stuff before call
        res, err := c.client.Do(ctx, r)
        // some stuff after call
        
        return res, err
	})
}
```