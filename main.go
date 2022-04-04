package main

import (
	"fmt"
	"net"
	"sort"
)

func main() {
	// //单线程端口扫描器
	// start := time.Now()
	// for i := 21; i < 120; i++ {
	// 	address := fmt.Sprintf("185.209.85.183:%d", i)
	// 	conn, err := net.Dial("tcp", address)
	// 	if err != nil {
	// 		fmt.Printf("%s closed\n", address)
	// 		continue
	// 	}
	// 	conn.Close()
	// 	fmt.Printf("%s opened\n", address)
	// }
	// // //计时
	// elapsed := time.Since(start) / 1e9
	// fmt.Printf("\n\n %d seconds", elapsed)

	// 多线程端口扫描器(goroutine)
	// 使用waitgroup等待里面的进程执行结束再整体退出
	// var waitgroup sync.WaitGroup
	// start := time.Now()
	// for i := 21; i < 65535; i++ {
	// 	waitgroup.Add(1)
	// 	go func(j int) {
	// 		defer waitgroup.Done()
	// 		address := fmt.Sprintf("185.209.85.183:%d", j)
	// 		conn, err := net.Dial("tcp", address)
	// 		if err != nil {
	// 			fmt.Printf("%s closed\n", address)
	// 			return
	// 		}
	// 		conn.Close()
	// 		fmt.Printf("%s opened\n", address)
	// 	}(i)

	// }
	// waitgroup.Wait()
	// //计时
	// elapsed := time.Since(start) / 1e9
	// fmt.Printf("\n\n %d seconds", elapsed)

	//多线程（并发）tcp端口扫描器 worker池（goroutine池）
	ports := make(chan int, 100)
	results := make(chan int)
	var openports []int
	var closeports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := 1; i < 1024; i++ {
			ports <- i
		}
	}()

	for i := 1; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		} else {
			closeports = append(closeports, port)
		}
	}
	close(ports)
	close(results)
	sort.Ints(openports)
	sort.Ints(closeports)

	for _, port := range closeports {
		fmt.Printf("%d Closed\n", port)
	}

	for _, port := range openports {
		fmt.Printf("%d Opened\n", port)
	}

}

func worker(ports chan int, result chan int) {
	for p := range ports {
		address := fmt.Sprintf("185.209.85.183:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			result <- 0
			continue
		}
		conn.Close()
		result <- p
	}

}
