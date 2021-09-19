package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//func errgroupUse() {
//	var g  errgroup.Group
//	var urls = []string{
//		"https://www.baidu.com",
//		"https://www.bilibili.com",
//		//"wwwa.bilibili.com",
//	}
//
//	for _, url := range urls {
//		g.Go(func() error {
//
//			url := url
//			response, err := http.Get(url)
//
//			if err == nil {
//				err = response.Body.Close()
//			}
//			return err
//		})
//	}
//
//	if err := g.Wait(); err == nil {
//		println("Successful fetched all URLs")
//	} else {
//		fmt.Printf("meet %+v\n", err)
//	}
//}

var (
	Exit       = errors.New("exit")
	Reload     = errors.New("reload")
	HotUpgrade = errors.New("hot upgrade")
)

func helloworld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello,World!"))
}

func startHttpServer(addr *http.Server) error {
	http.HandleFunc("/hello", helloworld)
	println("start http server")
	err := addr.ListenAndServe()
	return err
}

// homework
// 基于 errgroup 实现一个 http server 的启动和关闭 ，
// 以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。
func homework() {

	g, ctx := errgroup.WithContext(context.Background())

	// register signal
	g.Go(func() error {
		registerSignals := []os.Signal{
			syscall.SIGINT,
			syscall.SIGTERM,
		}
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, registerSignals...)
		select {
		case s := <-signals:
			switch s {
			case syscall.SIGINT:
				return Exit
			default:
				return nil
			}
		}
	})

	// start http server
	server := &http.Server{Addr: ":7878"}
	g.Go(func() error {

		err := startHttpServer(server)
		return err
	})
	
	g.Go(func() error {
		var err error
		select {
		case <-ctx.Done():
			err = server.Close()
		}
		return err
	})

	//g.Go(func() error {
	//	for {
	//
	//		println("goroutine 1")
	//		time.Sleep(1 * time.Second)
	//
	//		if ctx.Err() != nil {
	//			return ctx.Err()
	//		}
	//
	//	}
	//	return nil
	//})
	//
	//g.Go(func() error {
	//	for {
	//		time.Sleep(500 * time.Millisecond)
	//		println("goroutine 2")
	//
	//		if ctx.Err() != nil {
	//			return ctx.Err()
	//		}
	//	}
	//})

	if err := g.Wait(); err == nil {
		println("successful done")
	} else {
		fmt.Printf("with err : %+v", err)
	}
}

func main() {
	//errgroupUse()
	homework()
}
