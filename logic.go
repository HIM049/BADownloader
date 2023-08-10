package main

import (
	"fmt"
	"regexp"
	"strconv"
)

func Download(favlistID string, downloadCount int, downloadCompilation bool) error {
	// 先请求一次收藏夹基础信息，用于初始化循环
	favlist, err := GetFavListObj(favlistID, 1, 1)
	if err != nil {
		return err
	}
	fmt.Println("即将开始下载收藏夹“" + favlist.Data.Info.Title + "”")
	fmt.Println("共有 " + strconv.Itoa(favlist.Data.Info.Media_count) + " 个视频")
	// 如果用户输入的下载数量为 0 （全部下载）
	if downloadCount == 0 {
		downloadCount = favlist.Data.Info.Media_count
	}
	for i := 0; i < downloadCount; i++ {
		// 请求收藏夹信息，准备下载
		favlist, err := GetFavListObj(favlistID, 1, i+1)
		if err != nil {
			return err
		}
		// 将请求到的 BVID 传入，获取视频详情
		vInf, err := GetVideoPageInformationObj(favlist.Data.Medias[0].Bvid)
		if err != nil {
			return err
		}
		if vInf.Data.Videos <= 1 || !downloadCompilation {
			// 如果视频没有分 P / 用户不下载分 P
			fmt.Println("开始下载视频《" + vInf.Data.Title + "》")
			err = SimpleDownload(vInf.Data.Bvid, vInf.Data.Cid, CheckFileName(vInf.Data.Title))
			if err != nil {
				return err
			}
		} else {
			// 视频分 P 下载
			for i := 0; i < len(vInf.Data.Pages); i++ {
				fmt.Println("开始下载视频《" + vInf.Data.Title + "》的章节 " + strconv.Itoa(vInf.Data.Pages[i].Page))
				err = SimpleDownload(vInf.Data.Bvid, vInf.Data.Pages[i].Cid, CheckFileName(vInf.Data.Pages[i].Part))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func SimpleDownload(bvid string, cid int, fileName string) error {
	video, err := GetVideoObj(bvid, cid)
	if err != nil {
		return err
	}
	err = StreamingDownloader(video.Data.Dash.Audio[len(video.Data.Dash.Audio)-1].BaseUrl, fileName)
	if err != nil {
		return err
	}
	return nil
}

// 剔除文件名中的奇怪字符
func CheckFileName(SFileN string) string {
	re := regexp.MustCompile(`[/\$<>?:*|]`)
	newName := re.ReplaceAllString(SFileN, "")
	return newName
}

// func GetFavId(favListURL string) (string,error) {
// 	favURL, err := url.Parse(favListURL)
// 	if err != nil {
// 		return "", err
// 	}
// 	fid := GetStrNum(favURL.RawQuery)
// 	var result string
// 	result = fid[0]
// 	return result, nil
// }

// func GetStrNum(text string) []string {
// 	re := regexp.MustCompile(`\d+`)
// 	result := re.FindAllString(text, -1)
// 	return result
// }
