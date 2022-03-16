# Least Recently Used (LRU) Cache

Thread-safe Go implementation of a [Least Recently Used cache](https://en.wikipedia.org/wiki/Cache_replacement_policies#Least_recently_used_(LRU)).

## Usage

The library uses generic Key and Values [recently introduced in Go 1.18](https://tip.golang.org/doc/go1.18).


### String example

```golang
package main

import (
  lru "github.com/adamjq/lru-cache"
)

func main() {
  cache, err := lru.New[string, string](2) // init cache of capacity 2
  if err != nil {
    panic(err)
  }

  cache.Put("key-1", "value-1")
  cache.Put("key-2", "value-2")
  cache.Put("key-3", "value-3") // evicts key-1 from cache

  v := cache.Get("key-1") // nil
  v = cache.Get("key-2") // "value-2" and updates key-2 as most recently used
}
```

### Custom struct as value

```golang
package main

import (
  lru "github.com/adamjq/lru-cache"
)

type user struct {
  userID string
  name string
}

func main() {
  cache, err := lru.New[string, user](2) // instantiate cache with custom types
  if err != nil {
    panic(err)
  }

  userId := "22bc77a3-1456-470f-bdb0-0c893b8778a8"

  key := userId
  value := user{
    userID: userId,
    name:   "Adam",
  }

  cache.Put(key, value)
  v := cache.Get(key) // user
}
```

## Development

```bash
go test -race ./...
```
