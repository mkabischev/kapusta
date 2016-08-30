# Kapusta [![Build Status](https://travis-ci.org/mkabischev/kapusta.svg?branch=master)](https://travis-ci.org/mkabischev/kapusta) [![Coverage Status](https://coveralls.io/repos/github/mkabischev/kapusta/badge.svg)](https://coveralls.io/github/mkabischev/kapusta) [![Go Report Card](https://goreportcard.com/badge/github.com/mkabischev/kapusta?)](https://goreportcard.com/report/github.com/mkabischev/kapusta?) [![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hyperium/hyper/master/LICENSE)


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
	Do( *http.Request) (*http.Response, error)
}
```

`http.Client` supports it out of box!

## Middlewares:

```go
type MiddlewareFunc func(Client) Client
```

Kapusta provides some helpful middlewares for you:

- ```HeadersMiddleware(values map[string]string)``` Adds headers to requests
- ```HeaderMiddleware(name, value string)``` Like headers, but add only one header. 
- ```RecoverMiddleware()``` Converts all panics into errors
- ```BaseURLMiddleware(baseURL string)``` Replaces scheme and host to baseURL value.

## Usage

```go
decoratedClient := kapusta.Chain(
    http.DefaultClient,
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
type ClientFunc func(*http.Request) (*http.Response, error)
```

So the same example will be looks like:
```go
func AwesomeStuffDecorator(c kapusta.Client) kapusta.Client {
	return kapusta.ClientFunc(func(r *http.Request) (*http.Response, error) {
		// some stuff before call
        res, err := c.client.Do(r)
        // some stuff after call
        
        return res, err
	})
}
```

## Mock 

Also kapusta provides mock package for testing purposes. 



```go
client := mock.NewClient() // implements kapusta.Client
	
	client.
		Get("/path").
		WithBody(`{"foo": "bar"}`).
		WillReturn(200, `{"bar": "foo"}`)
		
// now you can inject client to your code.
```
