package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

// 通过命令行直接调用 ffmpeg
func ConvertM4aToMp3(inputPath, outputPath, logPath string) error {
	cmd := exec.Command("ffmpeg", "-y", "-i", inputPath+".m4a", "-vn", "-acodec", "libmp3lame", "-ab", "192k", "-ar", "44100", outputPath+".mp3")
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

// 多线程转码
func ConcurrentToMp3(threads int, inputPath, outputPath, logPath string) error {
	sem := make(chan struct{}, threads)
	var wg sync.WaitGroup
	// 转码
	conver := func(vInf *VideoInfLite, args ...interface{}) {
		sem <- struct{}{} // 限制并发量
		wg.Add(1)         // 任务 +1

		if len(args) >= 2 {
			outputPath, ok1 := args[0].(string)
			logPath, ok2 := args[1].(string)

			if ok1 && ok2 {
				err := ConvertM4aToMp3(inputPath+strconv.Itoa(vInf.Cid), outputPath+strconv.Itoa(vInf.Cid), logPath)
				if err != nil {
					return
				}
			} else {
				fmt.Println(("参数类型错误"))
				return
			}
		} else {
			fmt.Println("参数数量不足")
			return
		}

		// 下载完成后
		defer func() {
			<-sem     // 释放一个并发槽
			wg.Done() // 发出任务完成通知
		}()
	}

	// 分配任务
	err := VideoListLoop(conver, outputPath, logPath)
	if err != nil {
		return err
	}

	// 等待任务执行完成
	wg.Wait()
	return nil
}
