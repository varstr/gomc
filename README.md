#gomc#

golang binding for libmemcached

- Check Godoc[http://godoc.org/github.com/ianoshen/gomc] for more details.
- Check libmemcached[http://docs.libmemcached.org/] for even more details. 

##Examples##

```go
package main

import (
    "github.com/ianoshen/gomc"
    "fmt"
)

var (
    servers = []string {
        "host1:11211",
        "host2:11211",
        "host3:11211",
    }
)

func main() {
    cli, _ := gomc.NewClient(servers, 1, gomc.ENCODING_DEFAULT)
    cli.SetBehavior(gomc.BEHAVIOR_TCP_KEEPALIVE, 1)
    cli.SetBehavior(gomc.BEHAVIOR_HASH, uint64(gomc.HASH_MD5))

    var val string
    cli.Get("test-key", &val)
    fmt.Println(val)

    keys := []string{"key1", "key2", "key3"}
    res, _ := cli.GetMulti(keys)
    for _, key := keys {
        var val string
        res.Get("key1", &val)
        fmt.Println(val)
    }
}
```

##Encoding##

gomc will handle some encode/decode stuff between Go types and raw bytes stored in memcached. 

###Base Types###
- For base types, just use strconv no matter what encoding flag you choose when initializing the client.
- For string and byte silce, do nothing.

###Complex Types###
- Two options are available for complex structs and types: ENCODING_GOB/ENCODING_JSON.
- JSON is faster, but unable to dump some types like map[int]string.
- Decoding gob is relatively slow, but works for almost anything. In the worse case, implements GobEncoder and GobDecoder by yourself.

###Benchmark Detail###
```
BenchmarkEncodeDefault  10000000               292 ns/op
BenchmarkDecodeDefault   5000000               393 ns/op
BenchmarkEncodeGob        100000             23157 ns/op
BenchmarkDecodeGob         10000            227261 ns/op
BenchmarkEncodeJSON       100000             20631 ns/op
BenchmarkDecodeJSON        50000             34432 ns/op
```
