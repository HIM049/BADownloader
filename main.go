package main

import (
	"fmt"
	"os"
)

func main() {
	start()
}

func start() {
	_ = os.MkdirAll("./Download", 0755)
	fmt.Println("请输入要下载的 收藏夹URL / 收藏夹ID / 视频BV号")
}
