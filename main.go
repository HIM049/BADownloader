package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	fmt.Println("仓库链接: (https://github.com/HIM049/BADownloader)")
	start()
}

func start() {
	_ = os.MkdirAll("./Download", 0755)
	idType, ID := inputID()
	config := inputOption(idType)
	downloadCount := config.DownloadCount

	if idType == "fid" {
		favTitle, favCount := GetFavlistInf(ID)
		if downloadCount == 0 {
			downloadCount = favCount
		}
		fmt.Println("即将开始下载收藏夹 《" + favTitle + "》 ")
		fmt.Println("视频总数 " + strconv.Itoa(favCount) + " 个，将下载其中的 " + strconv.Itoa(downloadCount) + " 个")
		for i := 1; i <= downloadCount; i++ {
			Favobj := FavListObj(ID, i)
			fmt.Println("开始下载视频 " + strconv.Itoa(i) + " " + Favobj.Data.Medias[0].Title)
			BvDownload(Favobj.Data.Medias[0].Bvid, Favobj.Data.Medias[0].Title, config.DownloadCompilation)
		}

	} else {
		obj := videoPageInformationObj(ID)
		fmt.Println("开始下载 " + " 《" + obj.Data.Pages[0].Part + "》")
		BvDownload(ID, obj.Data.Pages[0].Part, config.DownloadCompilation)
	}
}

func inputID() (string, string) {
	for {
		fmt.Println("请输入要下载的 收藏夹URL / 收藏夹ID / 视频BV号")
		var favListURL string = ""
		fmt.Scanln(&favListURL)
		switch checkInput(favListURL) {
		case "https":
			fid := ToFavId(favListURL)
			return "fid", fid

		case "bv":
			return "bv", favListURL

		case "id":
			return "fid", favListURL

		default:
			fmt.Println("输入了错误的数据，请重新输入")
		}
	}
}

func inputOption(mode string) Config {
	var obj Config
	if mode == "fid" {
		fmt.Println("请输入希望下载的数量（按照收藏夹顺序从前往后。输入 “0” 完整下载收藏夹）")
		fmt.Scanln(&obj.DownloadCount)
	}
	fmt.Println("是否下载视频合集: (y/n) (默认值: y)")
	var downloadCompilation string = ""
	fmt.Scanln(&downloadCompilation)
	if downloadCompilation == "n" {
		obj.DownloadCompilation = false
	} else {
		obj.DownloadCompilation = true
	}
	return obj
}

func checkInput(text string) string {
	if regexp.MustCompile(`(?i)^https`).MatchString(text) {
		return "https"
	}
	if regexp.MustCompile(`(?i)^BV`).MatchString(text) {
		return "bv"
	}
	if regexp.MustCompile(`^\d`).MatchString(text) {
		return "id"
	}
	return "error"
}
