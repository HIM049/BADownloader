package main

import (
	"encoding/json"
	"fmt"
	"regexp"
)

func CheckErr(err error) {
	// 错误检查函数
	if err != nil {
		fmt.Println(err)
	}
}

func checkFileName(SFileN string) string {
	re := regexp.MustCompile(`[/\$<>?:*|]`)
	NFileN := re.ReplaceAllString(SFileN, "")
	return NFileN
}

func decodeJson(jsonFile string, object any) {
	// json解析函数。无返回
	err := json.Unmarshal([]byte([]byte(jsonFile)), object)
	CheckErr(err)
}

// 获取json数据并返回处理后的对象
func FavListObj(favListId string, pn int) FavList {
	var obj FavList
	decodeJson(GetFavList(favListId, pn), &obj)
	return obj
}
func videoPageInformationObj(bvid string) VideoInformation {
	var obj VideoInformation
	decodeJson(GetVideoPageInformation(bvid), &obj)
	return obj
}
func videoObj(bvid string, cid int) Video {
	var obj Video
	decodeJson(GetVideoJson(bvid, cid), &obj)
	return obj
}
func videoPageListObj(bvid string) VideoPageList {
	var obj VideoPageList
	decodeJson(GetVideoPageList(bvid), &obj)
	return obj
}

func GetAudioURL(bvid string, cid int) string {
	// 提取音频流URL
	audio := videoObj(bvid, cid)
	return audio.Data.Dash.Audio[0].BaseUrl
}

func AnalysisVideoList(bvid string) (int, int) {
	videoPageList := videoPageListObj(bvid)
	page := len(videoPageList.Data)  // 返回分P数量
	cid := videoPageList.Data[0].Cid // 返回CID
	return cid, page
}

func GetVideoInf(bvid string) (string, string) {
	// 使用BV获取视频分P标题和封面
	videoInf := videoPageInformationObj(bvid)
	title := checkFileName(videoInf.Data.Pages[0].Part)
	cover := videoInf.Data.Pic
	return title, cover
}
