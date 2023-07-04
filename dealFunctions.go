package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
)

func BvDownload(bvid string, fileName string, downCompilation bool) {
	cid, pageNum := analysisVideoList(bvid)
	if pageNum > 1 && downCompilation {
		fmt.Println("检测到该视频包含分集")
		for i := 0; i < pageNum; i++ {
			videoInf := videoPageInformationObj(bvid)
			fmt.Println("正在下载视频《" + fileName + "》中的第 " + strconv.Itoa(i+1) + " 个视频" + videoInf.Data.Pages[i].Part)
			videURL := getAudioURL(bvid, videoInf.Data.Pages[i].Cid)
			Downloader(videURL, videoInf.Data.Pages[i].Part)
			// videoInf.Data.Pic
		}
	} else {
		videURL := getAudioURL(bvid, cid)
		Downloader(videURL, fileName)
	}
}

func ToFavId(favListURL string) string {
	favURL, err := url.Parse(favListURL)
	CheckErr(err)
	fid := getStrNum(favURL.RawQuery)
	var result string
	result = fid[0]
	return result
}
func getStrNum(text string) []string {
	re := regexp.MustCompile(`\d+`)
	result := re.FindAllString(text, -1)
	return result
}

func CheckErr(err error) {
	// 错误检查函数
	if err != nil {
		fmt.Println(err)
	}
}

func checkObj(code int, message string) {
	if code == 0 {
		return
	} else {
		fmt.Println(message)
	}
}

func CheckFileName(SFileN string) string {
	re := regexp.MustCompile(`[/\$<>?:*|]`)
	NFileN := re.ReplaceAllString(SFileN, "")
	return NFileN
}

func decodeJson(jsonFile string, object any) {
	// json解析函数
	err := json.Unmarshal([]byte([]byte(jsonFile)), object)
	CheckErr(err)
}

func FavListObj(favListId string, pn int) FavList {
	// 获取json数据并返回处理后的对象
	var obj FavList
	decodeJson(GetFavList(favListId, pn), &obj)
	checkObj(obj.Code, obj.Message)
	return obj
}

func videoPageInformationObj(bvid string) VideoInformation {
	// 获取json数据并返回处理后的对象
	var obj VideoInformation
	decodeJson(GetVideoPageInformation(bvid), &obj)
	checkObj(obj.Code, obj.Message)
	return obj
}

func videoObj(bvid string, cid int) Video {
	// 获取json数据并返回处理后的对象
	var obj Video
	decodeJson(GetVideoJson(bvid, cid), &obj)
	checkObj(obj.Code, obj.Message)
	return obj
}

func videoPageListObj(bvid string) VideoPageList {
	// 获取json数据并返回处理后的对象
	var obj VideoPageList
	decodeJson(GetVideoPageList(bvid), &obj)
	checkObj(obj.Code, obj.Message)
	return obj
}

func getAudioURL(bvid string, cid int) string {
	// 提取音频流URL
	audio := videoObj(bvid, cid)
	return audio.Data.Dash.Audio[0].BaseUrl
}

func analysisVideoList(bvid string) (int, int) {
	videoPageList := videoPageListObj(bvid)
	page := len(videoPageList.Data)  // 返回分P数量
	cid := videoPageList.Data[0].Cid // 返回CID
	return cid, page
}

func GetFavlistInf(favid string) (string, int) {
	// 获取收藏夹视频数量及标题
	obj := FavListObj(favid, 1)
	count := obj.Data.Info.Media_count
	title := obj.Data.Info.Title
	return title, count
}

func VideoDownloadInf(favlistid string, pn int) (string, string) {
	obj := FavListObj(favlistid, pn)
	title := obj.Data.Medias[0].Title
	bvid := obj.Data.Medias[0].Bvid
	return bvid, title
}
