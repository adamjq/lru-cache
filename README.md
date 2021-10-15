# lru-cache

Thread-safe Golang implementation of a Least Recently Used cache.

## Usage

```golang
package main

import (
  lru "github.com/adamjq/lru-cache"
)

func main() {
  cache, err := lru.New(2) // init cache of capacity 2
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

## Development

```bash
make test
```
