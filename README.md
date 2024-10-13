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
        <a href="https://github.com/number571/micro-anon/pulse">
        	<img src="https://img.shields.io/github/commit-activity/m/number571/micro-anon" alt="Activity" />
		</a>
        <a href="https://github.com/number571/micro-anon/commits/master">
        	<img src="https://img.shields.io/github/last-commit/number571/micro-anon.svg" alt="Commits" />
		</a>
		<a href="https://github.com/number571/awesome-anonymity">
        	<img src="https://awesome.re/mentioned-badge.svg" alt="Awesome-Anonymity" />
		</a>
	</p>
	About project
</h2>

> [!WARNING]
> This anonymous network was written solely for the purpose of demonstrating the minimalism of the QB problem, and therefore it should be considered more as a template for modifications/editing than as a ready-made implementation. The implementation lacks an authentication mechanism, as well as mechanisms to counter DoS/DDoS attacks, spam and message repetition. You can read more about vulnerabilities and implementation-related issues at the end of this README.

The `Micro-Anonymous` network is based on a QB (queue-based) problem (also as [Hidden Lake](https://github.com/number571/hidden_lake)). The implementation uses only the standard library of the Go language. The goal of this network is to minimize the source code so that even a novice programmer can understand the entire mechanism of its functioning.

```bash
go run . [listen-address] [priv-key-file] [recv-key-file] [http-addr-1, http-addr-2, ...]
```

> More information about QB networks in research paper: [Анонимная сеть «Hidden Lake»](https://github.com/number571/go-peer/blob/master/docs/hidden_lake_anonymous_network.pdf)

## Example

```bash
# Terminal-1
$ go run . :7070 ./example/node2/priv.key ./example/node1/pub.key localhost:8080

# Terminal-2
$ go run . :8080 ./example/node1/priv.key ./example/node2/pub.key localhost:7070

# Terminal-1 <INPUT>
> hello

# Terminal-2 <OUTPUT>
> hello
```

## Advantages

1. <b>Simplicity</b>. The network is written without using third-party libraries, as well as without using hacks to pack it into 100 lines of code. As a result, even novice programmers can understand the logic of its operation, and even novice cryptographers can check for security.
2. <b>Anonymity</b>. This network really provides a good level of anonymity, protecting against all passive observations, including attacks by a global observer. Active observations are also impossible, because it requires the implementation of the composition of the conditions: 1) the attacker knows your public key, 2) you often talk to several subscribers at once, 3) the attacker must be in the list of subscribers with whom you are actively talking. But the implementation lacks the ability to communicate with multiple subscribers, and therefore the second condition will never be fulfilled.

## Vulnerabilities

1. <b>Lack of authentication</b>. It is unknown which particular participant sent you the message and there is no authenticated data to confirm that the interlocutor is exactly who he was introduced to at the beginning of the conversation.
2. <b>DoS/DDoS attacks</b>. An attacker can generate or collect many ciphertexts and send them to one node at a time, thereby overloading the processor power of the latter to perform decryption functions. Also, an attacker can generate large ciphertexts (or random bytes), thereby overloading the nodes' RAM due to the use of the io.ReadAll function.
3. <b>Spam</b>. Due to the lack of F2F or other trusted connection mechanisms, each node can communicate with any other node in the network if it knows the public key. As a result, an attacker can send many meaningless messages from different accounts without being able to block it. The only way to counteract it is to change your own private key.
4. <b>Repeat messages</b>. An attacker can copy the ciphertexts of the network and re-redirect them to a specific node due to the lack of verification of previously received messages. As a result, if the redirected ciphertext is true, the messages will be duplicated.

## License

Licensed under the MIT License. See [LICENSE](LICENSE) for the full license text.
