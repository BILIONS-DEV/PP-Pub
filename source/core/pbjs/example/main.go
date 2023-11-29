package main

import (
	"fmt"
	"source/core/pbjs"
	"strconv"
	"time"
)

type Result struct {
	Out string
	Err error
}

func main() {
	now := time.Now()

	/*
	 * chạy lệnh cài gulp: npm install -g gulp
	 * run file init-Prebid.js.sh trên terminal (linux) hoặc Git Bash/ Linux shell (window) để lấy pbjsPath
	 * Yêu cầu: cài đặt git, node, npm
	 */
	pbjsPath := "/home/apd/go/selfserve/source/core/pbjs/Prebid.js"
	PBJSBundle := pbjs.NewPBJSBundle(pbjsPath)

	listModules := []string{"openxBidAdapter"}

	c := make(chan Result)
	numberThreads := 5
	for i := 0; i < numberThreads; i++ {
		go func(x chan Result, i int) {
			fileName := "prebid-" + strconv.Itoa(i+1) + ".js"
			out, err := PBJSBundle.Build(listModules, fileName)
			x <- Result{
				Out: out,
				Err: err,
			}
		}(c, i)
	}

	countResult := 0
	countError := 0
	for result := range c {
		countResult++
		fmt.Println(result.Out, result.Err)
		if result.Err != nil {
			countError++
		}
		if countResult == numberThreads {
			close(c)
		}
	}

	fmt.Println("Build", numberThreads, "file trong", time.Since(now))
	fmt.Println("Thành công:", numberThreads-countError, "Lỗi:", countError)

}
