package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})))
	slog.Info(runtime.Version())

	// Setup graceful shutdown
	signals := make(chan os.Signal, 1)
	// Docker and Kubernetes use the SIGTERM signal to gracefully shut down a container.
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	withServer := len(os.Args) == 1
	if withServer {
		wg.Add(1)
		go serve(ctx, &wg)
	}

	wg.Add(1)
	go getPeriodically(ctx, &wg)

	// Wait for the termination signal
	v := <-signals
	slog.Info("Starting graceful shutdown", "Received termination signal", v)
	if withServer {
		err := server.Shutdown(ctx)
		if err != nil {
			slog.Error("server", "error", err)
		}
	}

	cancel()

	wg.Wait()
	slog.Info("App: graceful shutdown completed")
}

var interval = 1 * time.Second

func getPeriodically(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if interval <= 1*time.Second {
				interval = 5 * time.Second
				t.Reset(interval)
			}
			k, err := demo(ctx)
			if err != nil {
				slog.Error("get", "error", err.Error())
				break
			}
			slog.Info("get", "key", k)

		case <-ctx.Done():
			slog.Info("getPeriodically: graceful shutdown completed")
			return
		}
	}
}

var (
	address    = "http://127.0.0.1:8080/Acct"
	serverAddr = ":8080"
	client     = http.Client{
		Timeout: 10 * time.Second,
	}
	server *http.Server
)

func demo(ctx context.Context) (key string, err error) {
	get, err := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
	if err != nil {
		return
	}

	res, err := client.Do(get)
	if err != nil {
		return
	}
	defer res.Body.Close()

	var KeyId *Response
	err = json.NewDecoder(res.Body).Decode(&KeyId)
	if err != nil {
		return
	}

	key = KeyId.Key
	return
}

type Response struct {
	Key string
}

func serve(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	mux := http.NewServeMux()
	mux.HandleFunc("/Acct", home)

	// Create a custom HTTP server with read and write timeouts
	server = &http.Server{
		Addr:         serverAddr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second, // Set read timeout to 10 seconds
		WriteTimeout: 10 * time.Second, // Set write timeout to 10 seconds
	}

	slog.Info("Server", "Addr", server.Addr)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("Server", "error", err)
		panic(err)
	}
	slog.Info("Server: graceful shutdown completed")
}

func home(w http.ResponseWriter, r *http.Request) {
	userIP := r.Header.Get("X-Forwarded-For")
	if userIP == "" {
		userIP = r.RemoteAddr
	}
	slog.Info("server", "userIP", userIP)
	Data := []byte(`{ "Key":  "Key1234", "Account":  "Account1234" }`)
	w.Write(Data)
}
