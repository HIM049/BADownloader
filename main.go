package main

import (
	"fmt"
	"os"
	"strconv"
)

const (
	CONFIG_PATH     = "./config.json"
	CACHE_PATH      = "./Cache/"
	VIDEO_LIST_PATH = "./Cache/VideoList.json"
	M4A_PATH        = "./Cache/m4a/"
	MP3_PATH        = "./Cache/mp3/"
	COVER_PATH      = "./Cache/jpg/"
	LOG_PATH        = "./Cache/log/"
	OUT_PATH        = "./Download/"
)

func init() {
	fmt.Println("仓库链接: (https://github.com/HIM049/BADownloader)")
	_ = os.MkdirAll("./Download", 0755)
	_ = os.MkdirAll("./Cache", 0755)
	_ = os.MkdirAll("./Cache/m4a", 0755)
	_ = os.MkdirAll("./Cache/jpg", 0755)
	_ = os.MkdirAll("./Cache/mp3", 0755)
	_ = os.MkdirAll("./Cache/log", 0755)

	if !fileExists(CONFIG_PATH) {
		fmt.Println("已创建默认配置文件")
		err := createDefaultConfig()
		fmt.Println("创建文件时错误", err)
	}
}

func main() {
	SetupCommands()

	// body, _ := MakeVideoList("742380048", 10, true)
	// SaveJsonFile(VIDEO_LIST_PATH, body)
	// err := DownloadList(2)
	// fmt.Println(err)

}

// func start() {
// 	cfg, err := configIndex()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	err = Download(cfg.FavlistID, cfg.DownloadCount, cfg.DownloadCompilation)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

func ConfigIndex() (*Config, error) {
	var cfg Config
	favlistID := getStrInput("请输入收藏夹 ID")
	downloadCount, _ := strconv.Atoi(getStrInput("请输入希望下载的视频数量：\n全部下载请输入【0】"))
	downloadCompilation, _ := strconv.Atoi(getStrInput("请选择是否下载合集：\n输入【0】下载全部视频分 P \n输入【1】不下载视频分 P \n如果需要单独设置可以在程序生成任务列表后单独修改"))
	threads, _ := strconv.Atoi(getStrInput("请输入希望的最大线程数：默认值【5】\n这项设置将会应用到全部多线程部分。提高上限可以提高运行速度，但会增加占用。"))
	ffmpegPath := getStrInput("请设置 FFmpeg 路径：（格式转换功能依赖 ffmpeg ，无需可跳过）\n")
	coverToMp3 := getBoolInput("请选择是否自动转换格式为 MP3：")
	writeMetaInf := getBoolInput("请选择是否自动填写歌曲元数据：\n（歌名、封面图等）")
	if downloadCount != 0 {
		favlist, err := GetFavListObj(favlistID, 1, 1)
		if err != nil {
			return nil, err
		}
		downloadCount = favlist.Data.Info.Media_count - downloadCount
	}
	cfg = Config{
		FavlistID:           favlistID,
		DownloadCount:       downloadCount,
		DownloadCompilation: downloadCompilation,
		DownloadThreads:     threads,
		FFmpegThreads:       threads,
		FFmpegPath:          ffmpegPath,
		CoverToMp3:          coverToMp3,
		WriteMetaInf:        writeMetaInf,
	}
	return &cfg, nil
}

// 配置文件结构
type Config struct {
	FavlistID           string // 收藏夹 ID
	DownloadCount       int    // 下载计数
	DownloadCompilation int    // 合集下载
	DownloadThreads     int    // 下载线程数
	CoverToMp3          bool   // 是否转换格式到 MP3
	FFmpegPath          string // ffmpeg 路径
	FFmpegThreads       int    // ffmpeg 线程数
	WriteMetaInf        bool   // 写入歌曲元数据
}

// 检查是否存在配置文件
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// 创建默认配置文件
func createDefaultConfig() error {
	defaultConfig := Config{
		DownloadCompilation: 0,
		DownloadThreads:     5,
		CoverToMp3:          false,
		FFmpegPath:          "ffmpeg",
		FFmpegThreads:       5,
	}

	err := SaveJsonFile(CONFIG_PATH, defaultConfig)
	if err != nil {
		return err
	}
	return nil
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
