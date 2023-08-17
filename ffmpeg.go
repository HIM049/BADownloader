package main

import (
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/cheggaaa/pb/v3"
)

// 多线程转码
func ConcurrentToMp3(threads int, ffmpegPath, inputPath, outputPath, logPath string) error {
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
			sem <- struct{}{} // 限制并发量
			wg.Add(1)         // 任务 +1

			err := ConvertM4aToMp3(ffmpegPath, inputPath+strconv.Itoa(v.Cid), outputPath+strconv.Itoa(v.Cid), logPath)
			if err != nil {
				return
			}

			// 下载完成后
			defer func() {
				<-sem     // 释放一个并发槽
				wg.Done() // 发出任务完成通知
			}()
		}(video)
	}
	// 等待任务执行完成
	wg.Wait()
	progressBar.Finish()
	return nil
}

// 通过命令行直接调用 ffmpeg
func ConvertM4aToMp3(ffmpegPath, inputPath, outputPath, logPath string) error {
	cmd := exec.Command(ffmpegPath, "-y", "-i", inputPath+".m4a", "-vn", "-acodec", "libmp3lame", "-ab", "192k", "-ar", "44100", outputPath+".mp3")
	stderrFile, err := os.Create(logPath + "ffmpeg_output.txt")
	if err != nil {
		return err
	}
	defer stderrFile.Close()

	cmd.Stderr = stderrFile
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
