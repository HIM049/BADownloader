package main

import (
	"os"
	"strconv"
	"sync"

	"github.com/cheggaaa/pb/v3"
	tag "github.com/gcottom/mp3-mp4-tag"
)

// 写入元数据和重命名文件
func ConcurrentChangeTagAndName(threads int, coverPath, audioType, aideoSourcePath, aideoDestPath string) error {
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
		// 转码
		go func(v VideoInformationList) {
			sem <- struct{}{} // 限制并发量
			wg.Add(1)         // 任务 +1

			// TODO: 将标签分析移动到列表

			// 处理音频标题
			NfileName := v.Title
			songName, err := ExtractTitle(v.Title)
			if err != nil {
				// 如果无法判断标题
				songName = v.Title
			}
			// 如果是分 P （以分 P 命名为主）
			if v.IsPage {
				NfileName = v.Title + "(" + v.PageTitle + ")"
				songName, err = ExtractTitle(v.PageTitle)
				if err != nil {
					// 如果无法判断标题
					songName = v.PageTitle
				}
			}
			// 预处理文件名+后缀
			fileName := strconv.Itoa(v.Cid) + audioType
			// 写入歌曲元数据
			err = ChangeTag(fileName, songName, coverPath+strconv.Itoa(v.Cid)+".jpg")
			if err != nil {
				return
			}
			// 重命名歌曲文件并移动位置
			err = RenameAndMoveFile(aideoSourcePath+fileName, aideoDestPath+NfileName+audioType)
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

// 修改 TAG
func ChangeTag(filename, songName, coverPath string) error {
	tags, err := tag.OpenTag("./Download/cache/" + filename)
	if err != nil {
		return err
	}

	tags.SetTitle(songName)                 // 歌曲名
	tags.SetAlbumArtFromFilePath(coverPath) // 封面路径

	// TODO: 将歌曲 tag 数据整理为结构体
	// TODO: 修改作词人，作曲人等，以及自动适配

	tags.Save() // 保存更改

	return nil
}

// 重命名和移动
func RenameAndMoveFile(sourcePath, destPath string) error {
	err := os.Rename(sourcePath, destPath)
	if err != nil {
		return err
	}
	return nil
}
