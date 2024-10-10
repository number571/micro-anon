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
	ctx := context.TODO()
	go func() { _ = runQBProblem(ctx, getReceiverKey(os.Args[3]), os.Args[4:]) }()
	_ = runMessageHandler(ctx, getPrivateKey(os.Args[2]), os.Args[1])
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
			for _, host := range hosts {
				client := &http.Client{Timeout: time.Second}
				_, _ = client.Post(fmt.Sprintf("http://%s/push", host), "text/plain", bytes.NewBuffer(encBytes))
			}
		}
	}
}

func getPrivateKey(privateKeyFile string) *rsa.PrivateKey {
	privKeyBytes, _ := os.ReadFile(privateKeyFile)
	priv, err := x509.ParsePKCS1PrivateKey(privKeyBytes)
	doif(err != nil, func() { panic(err) })
	return priv
}

func getReceiverKey(receiverKeyFile string) *rsa.PublicKey {
	pubKeyBytes, _ := os.ReadFile(receiverKeyFile)
	pub, err := x509.ParsePKCS1PublicKey(pubKeyBytes)
	doif(err != nil, func() { panic(err) })
	return pub
}

func doif(isTrue bool, do func()) {
	if isTrue {
		do()
	}
}
