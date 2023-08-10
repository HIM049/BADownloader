package main

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	FavlistID           string
	DownloadCount       int
	DownloadCompilation bool
}

func main() {
	fmt.Println("仓库链接: (https://github.com/HIM049/BADownloader)")
	_ = os.MkdirAll("./Download", 0755)
	start()
}

func start() {
	cfg, err := userConfig()
	if err != nil {
		fmt.Println(err)
	}
	err = Download(cfg.FavlistID, cfg.DownloadCount, cfg.DownloadCompilation)
	if err != nil {
		fmt.Println(err)
	}
}

func userConfig() (*Config, error) {
	var cfg Config
	favlistID := getStrInput("请输入收藏夹 ID")
	downloadCount, err := strconv.Atoi(getStrInput("请输入下载数量：\n全部下载请输入【0】"))
	if err != nil {
		return nil, err
	}
	downloadCompilation := getBoolInput("是否下载合集？")
	cfg = Config{
		FavlistID:           favlistID,
		DownloadCount:       downloadCount,
		DownloadCompilation: downloadCompilation,
	}
	return &cfg, nil
}

func getStrInput(text string) string {
	// 打印提示词
	fmt.Println(text)
	// 等待输入并返回
	var input string
	fmt.Scan(&input)
	return input
}

func getBoolInput(text string) bool {
	// 打印提示词
	fmt.Println(text)
	for {
		// 调用 getStrInput 获取输入
		input := getStrInput("【输入 [T] 表示是，输入 [F] 表示否】")
		choices := map[string]bool{
			"T": true,
			"t": true,
			"F": false,
			"f": false,
		}

		result, ok := choices[input]
		if !ok {
			// 如果无法判断，再次输入
			fmt.Println("无效的输入")
			continue
		}
		// 判断成功，返回结果
		return result
	}

}

// func inputID() (string, string) {
// 	for {
// 		fmt.Println("请输入要下载的 收藏夹URL / 收藏夹ID / 视频BV号")
// 		var favListURL string = ""
// 		fmt.Scanln(&favListURL)
// 		switch checkInput(favListURL) {
// 		case "https":
// 			fid := ToFavId(favListURL)
// 			return "fid", fid

// 		case "bv":
// 			return "bv", favListURL

// 		case "id":
// 			return "fid", favListURL

// 		default:
// 			fmt.Println("输入了错误的数据，请重新输入")
// 		}
// 	}
// }

// func inputOption(mode string) Config {
// 	var obj Config
// 	if mode == "fid" {
// 		fmt.Println("请输入希望下载的数量（按照收藏夹顺序从前往后。输入 “0” 完整下载收藏夹）")
// 		fmt.Scanln(&obj.DownloadCount)
// 	}
// 	fmt.Println("是否下载视频合集: (y/n) (默认值: y)")
// 	var downloadCompilation string = ""
// 	fmt.Scanln(&downloadCompilation)
// 	if downloadCompilation == "n" {
// 		obj.DownloadCompilation = false
// 	} else {
// 		obj.DownloadCompilation = true
// 	}
// 	return obj
// }

// func checkInput(text string) string {
// 	if regexp.MustCompile(`(?i)^https`).MatchString(text) {
// 		return "https"
// 	}
// 	if regexp.MustCompile(`(?i)^BV`).MatchString(text) {
// 		return "bv"
// 	}
// 	if regexp.MustCompile(`^\d`).MatchString(text) {
// 		return "id"
// 	}
// 	return "error"
// }
