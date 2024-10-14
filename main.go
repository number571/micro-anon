package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	privKeyBytes, _ := os.ReadFile(os.Args[2])
	privKey, err := x509.ParsePKCS1PrivateKey(privKeyBytes)
	doif(err != nil, func() { panic(err) })

	pubKeyBytes, _ := os.ReadFile(os.Args[3])
	pubKey, err := x509.ParsePKCS1PublicKey(pubKeyBytes)
	doif(err != nil, func() { panic(err) })

	ctx := context.TODO()
	go func() { _ = runQBProblem(ctx, pubKey, os.Args[4:]) }()
	fmt.Println(runMessageHandler(ctx, privKey, os.Args[1]))
}

func runMessageHandler(ctx context.Context, privateKey *rsa.PrivateKey, addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		encBytes, _ := io.ReadAll(r.Body)
		decBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encBytes, nil)
		doif(err == nil, func() { fmt.Println(string(decBytes)) })
	})
	server := &http.Server{Addr: addr, Handler: mux}
	go func() {
		<-ctx.Done()
		server.Close()
	}()
	return server.ListenAndServe()
}

func runQBProblem(ctx context.Context, receiverKey *rsa.PublicKey, hosts []string) error {
	queue := make(chan []byte, 256)
	go func() {
		pr, err := rsa.GenerateKey(rand.Reader, receiverKey.N.BitLen())
		doif(err != nil, func() { panic(err) })
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if len(queue) == 0 {
					encBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &pr.PublicKey, []byte("_"), nil)
					doif(err == nil, func() { queue <- encBytes })
				}
			}
		}
	}()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				input, _, _ := bufio.NewReader(os.Stdin).ReadLine()
				encBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, receiverKey, input, nil)
				doif(err == nil, func() { queue <- encBytes })
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
			encBytes := <-queue
			client := &http.Client{Timeout: time.Second}
			for _, host := range hosts {
				_, _ = client.Post(fmt.Sprintf("http://%s/push", host), "text/plain", bytes.NewBuffer(encBytes))
			}
		}
	}
}

func doif(isTrue bool, do func()) {
	if isTrue {
		do()
	}
}
