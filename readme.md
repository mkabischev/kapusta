# Kapusta

[![Build Status](https://travis-ci.org/mkabischev/kapusta.svg)](https://travis-ci.org/mkabischev/kapusta)

It`s middleware approach for using **http.Client**,  inspired by [Embrace the Interface talk](https://github.com/gophercon/2015-talks/tree/master/Tom%C3%A1s%20Senart%20-%20Embrace%20the%20Interface). You can wrap your client with different functionality: 

 - log every request
 - append auth headers
 - use http cache
 - use etcd for service discovery
 - and whatever you want!
 
**Just like a cabbage!**

![](http://2.bp.blogspot.com/-LtmW_ktxtXU/Um28ElCtQnI/AAAAAAAAB04/aaXWbmTdbnE/s1600/cabbage.png)

## Client interface

Internal http package doesn`t have any interface for http clients, so Kapusta provides very simple client interface:
```go
type Client interface {
	Do(*http.Request) (*http.Response, error)
}
```
`http.Client` supports it out of box.

## Decorators

Decorator is like a middleware:
```go
type DecoratorFunc func(Client) Client
```

Kapusta provides some helpful decorators for you:

- ```HeadersDecorator(values map[string]string)``` Adds headers to requests
- ```HeaderDecorator(name, value string)``` Like headers, but add only one header. 
- ```RecoverDecorator()``` Converts all panics into errors
- ```BaseURLDecorator(baseURL string)``` Replaces scheme and host to baseURL value.

## Usage

```go
client := http.DefaultClient

decoratedClient := kapusta.Decorate(
    client,
    kapusta.HeaderDecorator("X-Auth", "123"),
    kapusta.RecoverDecorator(), // better to place it last to recover panics from decorators too
)
```

## Create your own decorator

There are two ways of creating new decorators.

You can create some new struct:
```go
struct AwesomeStuffClient {
    client kapusta.Client
}

func(c *AwesomeStuffClient) Do(r *http.Request) (*http.Response, error) {
    // some stuff before call
    res, err := c.client.Do(r)
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

Sometimes it`s required to pass some params in decorator, for details see Headers decorator.
