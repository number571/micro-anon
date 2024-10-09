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
	ctx := context.Background()
	go func() { _ = runQBProblem(ctx, getReceiverKey(os.Args[3]), os.Args[4:]) }()
	runHandler(ctx, getPrivateKey(os.Args[2]), os.Args[1])
}

func runHandler(_ context.Context, privateKey *rsa.PrivateKey, addr string) {
	http.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		encBytes, _ := io.ReadAll(r.Body)
		decBytes, _ := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encBytes, nil)
		if len(decBytes) == 0 {
			return
		}
		fmt.Println(string(decBytes))
	})
	http.ListenAndServe(addr, nil)
}

func runQBProblem(ctx context.Context, receiverKey *rsa.PublicKey, hosts []string) error {
	queue := make(chan []byte, 256)
	go func() {
		pr, _ := rsa.GenerateKey(rand.Reader, receiverKey.N.BitLen())
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if len(queue) == 0 {
					encBytes, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, &pr.PublicKey, []byte("_"), nil)
					push(queue, encBytes)
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
				encBytes, _ := rsa.EncryptOAEP(sha256.New(), rand.Reader, receiverKey, input, nil)
				push(queue, encBytes)
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
			for _, host := range hosts {
				client := &http.Client{Timeout: time.Second}
				_, _ = client.Post(fmt.Sprintf("http://%s/push", host), "text/plain", bytes.NewBuffer(<-queue))
			}
		}
	}
}

func push(queue chan<- []byte, bytes []byte) {
	if len(bytes) == 0 {
		return
	}
	queue <- bytes
}

func getPrivateKey(privateKeyFile string) *rsa.PrivateKey {
	privKeyBytes, _ := os.ReadFile(privateKeyFile)
	priv, err := x509.ParsePKCS1PrivateKey(privKeyBytes)
	if err != nil {
		panic(err)
	}
	return priv
}

func getReceiverKey(receiverKeyFile string) *rsa.PublicKey {
	pubKeyBytes, _ := os.ReadFile(receiverKeyFile)
	pub, err := x509.ParsePKCS1PublicKey(pubKeyBytes)
	if err != nil {
		panic(err)
	}
	return pub
}
