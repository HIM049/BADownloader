package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/cheggaaa/pb/v3"
)

// 用于获取收藏夹基本信息的函数
// 传入收藏夹 ID ，ps 单页大小， pn 页码
// 获得如下结构体
type FavList struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 错误消息
	Data    struct {
		Info struct { // 收藏夹信息
			Title       string `json:"title"`       // 收藏夹标题
			Media_count int    `json:"media_count"` // 收藏夹数量
		}
		Medias []struct { // 收藏夹中的视频
			Id    int    `json:"id"`    // 稿件 avid
			Type  int    `json:"type"`  // 内容类型 （视频稿件2 音频12 合集21）
			Title string `json:"title"` // 标题
			Cover string `json:"cover"` // 封面 url
			Page  int    `json:"page"`  // 视频分P数
			Bvid  string `json:"bvid"`  // BV 号
		}
	}
}

func getFavList(favListId string, ps int, pn int) (string, error) {
	// 设置 URL 并发送 GET 请求
	params := url.Values{}
	Url, _ := url.Parse("https://api.bilibili.com/x/v3/fav/resource/list")
	params.Set("media_id", favListId)
	params.Set("ps", strconv.Itoa(ps))
	params.Set("platform", "web")
	params.Set("pn", strconv.Itoa(pn))
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, err := http.Get(urlPath)
	if err != nil {
		return "", err
	}

	// 将 body 转为字符串并返回
	body, _ := io.ReadAll(resp.Body)
	bodyString := string(body)
	defer resp.Body.Close()
	return bodyString, nil
}

func GetFavListObj(favListId string, ps int, pn int) (*FavList, error) {
	var obj FavList
	body, err := getFavList(favListId, ps, pn)
	if err != nil {
		return nil, err
	}
	err = decodeJson(body, &obj)
	if err != nil {
		return nil, err
	}
	// 错误检查
	if checkObj(obj.Code) {
		return nil, errors.New(obj.Message)
	}
	return &obj, nil
}

// 用于获取视频的详细信息
// 传入 BVID
// 获得如下结构体
type VideoInformation struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Bvid   string     `json:"bvid"`   // 稿件 BVID
		Videos int        `json:"videos"` // 稿件分 P 总数
		Pic    string     `json:"pic"`    // 稿件封面图片url
		Title  string     `json:"title"`  // 稿件标题
		Cid    int        `json:"cid"`    // 视频1P cid
		Pages  []struct { // 分 P 列表
			Cid  int    `json:"cid"`  // 分 P cid
			Page int    `json:"page"` // 分 P 序号
			Part string `json:"part"` // 分 P 标题
		}
	}
}

func getVideoPageInformation(bvid string) (string, error) {
	// 设置 URL 并发送 GET 请求
	params := url.Values{}
	Url, _ := url.Parse("https://api.bilibili.com/x/web-interface/view")
	params.Set("bvid", bvid)
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, err := http.Get(urlPath)
	if err != nil {
		return "", err
	}
	// 将 body 转为字符串并返回
	body, _ := io.ReadAll(resp.Body)
	bodyString := string(body)
	defer resp.Body.Close()
	return bodyString, nil
}

func GetVideoPageInformationObj(bvid string) (*VideoInformation, error) {
	var obj VideoInformation
	body, err := getVideoPageInformation(bvid)
	if err != nil {
		return nil, err
	}
	err = decodeJson(body, &obj)
	if err != nil {
		return nil, err
	}
	// 错误检查
	if checkObj(obj.Code) {
		return nil, errors.New(obj.Message)
	}
	return &obj, nil
}

// 用于获取视频流的详细信息
// 传入 BVID 和 CID
// 获得如下结构体
type Video struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Dash struct {
			Audio []struct {
				Id       int    `json:"id"`
				BaseUrl  string `json:"baseUrl"`
				MimeType string `json:"mimeType"`
			}
			Flac struct {
				Audio struct {
					Id       int    `json:"id"`
					BaseUrl  string `json:"baseUrl"`
					MimeType string `json:"mimeType"`
				}
			}
		}
	}
}

func getVideo(bvid string, cid int) (string, error) {
	// 设置 URL 并发送 GET 请求
	params := url.Values{}
	Url, _ := url.Parse("https://api.bilibili.com/x/player/playurl")
	params.Set("bvid", bvid)
	params.Set("cid", strconv.Itoa(cid))
	params.Set("fnval", "16")
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, err := http.Get(urlPath)
	if err != nil {
		return "", err
	}
	// 将 body 转为字符串并返回
	body, _ := io.ReadAll(resp.Body)
	bodyString := string(body)
	defer resp.Body.Close()
	return bodyString, nil
}
func GetVideoObj(bvid string, cid int) (*Video, error) {
	var obj Video
	body, err := getVideo(bvid, cid)
	if err != nil {
		return nil, err
	}
	err = decodeJson(body, &obj)
	if err != nil {
		return nil, err
	}
	if checkObj(obj.Code) {
		return nil, errors.New(obj.Message)
	}
	return &obj, nil
}

// 用于下载音频流的函数
// 传入流 URL 和文件名
func StreamingDownloader(audioURL string, filePathAndName string) error {
	// 音频流下载函数。接收音频url和文件名。
	client := &http.Client{}
	request, err := http.NewRequest("GET", audioURL, nil)
	if err != nil {
		return err
	}
	request.Header.Set("referer", "https://www.bilibili.com")
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	out, err := os.Create(filePathAndName)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return err
	}
	return nil
}

func SaveFromURL(url string, filePath string) error {
	// 发起 HTTP 请求获取图片内容
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// 创建目标文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将图片内容写入文件
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func ConcurrentSavePic(threads int, savePath string) error {
	sem := make(chan struct{}, threads)
	var wg sync.WaitGroup

	// 获取任务队列
	var list []VideoInformationList
	err := LoadJsonFile(VIDEO_LIST_PATH, &list)
	if err != nil {
		return err
	}
	// 设置进度条
	progressBar := pb.Full.Start(len(list))
	// 遍历下载队列
	for _, video := range list {

		go func(v VideoInformationList) {
			sem <- struct{}{} // 限制并发量
			wg.Add(1)         // 任务 +1

			err := SaveFromURL(v.Cover, savePath+strconv.Itoa(v.Cid)+".jpg")
			if err != nil {
				return
			}

			// 下载完成后
			defer func() {
				<-sem                   // 释放一个并发槽
				wg.Done()               // 发出任务完成通知
				progressBar.Increment() // 进度条增加
			}()
		}(video)
	}
	// 等待任务执行完成
	wg.Wait()
	progressBar.Finish()
	return nil
}

// 工具函数
// json解析函数
func decodeJson(jsonFile string, object any) error {
	err := json.Unmarshal([]byte([]byte(jsonFile)), object)
	if err != nil {
		return err
	}
	return nil
}

// 工具函数
// 检查结构体中的状态码
func checkObj(code int) bool {
	if code == 0 {
		return false
	} else {
		return true
	}
}
