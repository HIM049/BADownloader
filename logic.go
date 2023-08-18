package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"

	"github.com/cheggaaa/pb/v3"
)

func DownloadList(threads int) error {
	fmt.Println("开始下载任务列表")
	sem := make(chan struct{}, threads+1)
	var wg sync.WaitGroup
	var progressBar *pb.ProgressBar

	// 获取任务队列
	var list []VideoInformationList
	err := LoadJsonFile(VIDEO_LIST_PATH, &list)
	if err != nil {
		return err
	}
	// 设置进度条
	progressBar = pb.Full.Start(len(list))
	// 遍历下载队列
	for _, video := range list {

		go func(v VideoInformationList) {
			// fmt.Println("调用下载")
			// 下载完成后
			defer func() {
				progressBar.Increment()
				<-sem     // 释放一个并发槽
				wg.Done() // 发出任务完成通知
			}()

			sem <- struct{}{} // 给通道中
			wg.Add(1)         // 任务 +1

			err := SimpleDownload(v.Bvid, v.Cid, strconv.Itoa(v.Cid)+".m4a")
			if err != nil {
				fmt.Printf("SimpleDownload：%s", err)
			}

		}(video)
	}
	// 等待任务执行完成
	wg.Wait()
	progressBar.Finish()
	return nil
}

func SimpleDownload(bvid string, cid int, fileName string) error {
	video, err := GetVideoObj(bvid, cid)
	if err != nil {
		return err
	}
	// 下载媒体流
	err = StreamingDownloader(video.Data.Dash.Audio[0].BaseUrl, M4A_PATH+fileName)
	if err != nil {
		return err
	}
	return nil
}

func SaveJsonFile(filePath string, theData any) error {
	data, err := json.MarshalIndent(theData, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// func SaveListFile(slice []VideoInformationList) error {
// 	jsonFile, err := os.Create(VIDEO_LIST_PATH)
// 	if err != nil {
// 		return err
// 	}
// 	defer jsonFile.Close()

// 	// encoder := json.NewEncoder(jsonFile)
// 	// err = encoder.Encode(slice)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// // 使用 json.MarshalIndent 编码数据为格式化的 JSON 字符串
// 	// formattedJSON, err := json.MarshalIndent(slice, "", "    ")
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// // 写入格式化的 JSON 数据到文件
// 	// _, err = jsonFile.Write(formattedJSON)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	return nil
// }

func LoadJsonFile(filePath string, obj interface{}) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, obj)
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

func ExtractTitle(input string) (string, error) {
	// 定义书名号正则表达式
	re := regexp.MustCompile(`《(.*?)》`)

	// 查找匹配的字符串
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return "", errors.New("无法找到合适的书名号")
	}

	// 返回匹配的书名号内容
	return matches[1], nil
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
