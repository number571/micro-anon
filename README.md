<img src="images/ma_logo.png" alt="ma_logo.png"/>

<h2>
	<p align="center">
    	<strong>
        	Anonymity network in 100 lines of code
   		</strong>
	</p>
	<p align="center">
        <a href="https://github.com/topics/golang">
        	<img src="https://img.shields.io/github/go-mod/go-version/number571/micro-anon" alt="Go" />
		</a>
        <a href="https://github.com/number571/micro-anon/blob/master/LICENSE">
        	<img src="https://img.shields.io/github/license/number571/micro-anon.svg" alt="License" />
		</a>
	</p>
	About project
</h2>

The `Micro-Anonymous` network is based on a QB (queue-based) problem (also as [Hidden Lake](https://github.com/number571/go-peer/tree/master/cmd/hidden_lake)). The implementation uses only the standard library of the Go language. The goal of this network is to minimize the source code so that even a novice programmer can understand the entire mechanism of its functioning.

```bash
usage: 
    go run . [listen-address] [private-key-file] [receiver-key-file] [http-address-1, http-address-2, ...]
```

> More information about QB networks in research paper: [Анонимная сеть «Hidden Lake»](https://github.com/number571/go-peer/blob/master/docs/hidden_lake_anonymous_network.pdf)

## Example

```bash
# Terminal-1
$ go run . :7070 ./example/node2/priv.key ./example/node1/pub.key localhost:8080
> hello <INPUT>

# Terminal-2
$ go run . :8080 ./example/node1/priv.key ./example/node2/pub.key localhost:7070
> hello <OUTPUT>
```
