package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

func main() {

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error { return myServe(ctx) })
	g.Go(func() error { return mySingal(ctx) })
	g.Wait()
}

func myServe(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, rep *http.Request) {
		fmt.Fprintln(resp, "Hello,Qcon")
	})

	s := http.Server{
		Addr:    ":22222",
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		fmt.Println("sShutdown")
		s.Shutdown(context.Background())

	}()
	return s.ListenAndServe()
}

func mySingal(ctx context.Context) error {
	sigs := make(chan os.Signal)
	signal.Notify(sigs)

	var err error = nil
	select {
	case <-sigs:
		fmt.Println("mySingal reveive signal")
		err = taskOne("reveive signal")
	case <-ctx.Done():
		fmt.Println("mySingal reveive done")
	}
	signal.Stop(sigs)
	close(sigs)
	fmt.Println("mySingal reveive end")
	return err
}

func taskOne(err string) error {
	return errors.New(err)
}
