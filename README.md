［## 使用拥有 UI 界面的全新 Bili Audio Downloader］(https://github.com/HIM049/BADownloaderUI)

---

# BiliAudioDownloader
一个用于批量下载B站收藏夹中视频音频的工具  
**如果项目对你有帮助，请给我一个 star!**  
项目写的比较杂乱，还请各位见谅。如果各位有不错的新功能提议或遇到了 BUG 欢迎提 issue 。  

## BADownloader 2.0
BAD 已经发布了 2.0 版本！以下是更新列表

- 重构了项目结构
- 重构了运行逻辑
- 大部分逻辑改用多线程
- 支持自动下载封面
- 支持自动写入元数据
- 支持批量转码为 MP3 （依赖 FFmpeg ）
- 程序的工作方式改为使用命令行参数

## 使用教程

### 生成任务列表
新版的 BADownloader 增加了一个下载列表。方便后续不同的步骤可以读取视频信息，无需重复请求内容。  
列表也可以手动修改。如删除个别不需要的分 P ，或是手动注释曲名等  
输入以下命令生成一个列表
```shell
./BADownloader.exe list bulid -id <你的收藏夹 ID> -c <需要下载的数量（为 0 则下载全部）> -p <是否下载分 P>
```
**示例**
```shell
./BADownloader.exe list bulid -id 742380048 -c 20 -p true
```
生成的文件会保存在 `./Cache/VideoList.json` 中。可以手动删改其中的内容以适应下载需求  

---  

### 下载列表内容
生成并修改完列表后，就可以开始下载列表中的内容了。  
输入以下命令自动下载列表内全部内容
```shell
./BADownloader.exe download -t <线程数>
```
**示例**
```shell
./BADownloader.exe download -t 5
```

---  


### 自动判断并写入歌曲元数据
目前 BADownload 2.0 只会写入元数据中的曲名和封面部分。  
自动识别的曲名为被标题中 书名号“《》” 包裹的内容。如标题中没有书名号，则填入视频标题本身。  
输入以下命令自动判断并写入歌曲元数据。
```shell
./BADownloader.exe tag add -t <线程数> -type <歌曲类型（未转码前即为 m4a，转码后即为 mp3 ）>
```
**示例**
```shell
./BADownloader.exe tag add -t 5 -type m4a
```

---  


### 批量转码为 MP3
目前 TAG 功能无法给 MP3 格式的音频添加封面，如需要自动添加封面请不要转码。  
输入以下命令将所有音频转码为 MP3  
```shell
./BADownloader.exe conver -t <线程数> -p <FFmpeg 路径>
```
**示例**
```shell
./BADownloader.exe conver -t 5 -p "C:\Program Files\ffmpeg\bin\ffmpeg.exe"
```

---  


### 导出音频
输入以下命令将音频输出到 `./Download/` 目录中并清除缓存目录。  
```shell
./BADownloader.exe out -t <线程数> -type <音频类型（未转码前即为 m4a，转码后即为 mp3 ）>
```
**示例**
```shell
./BADownloader.exe out -t 5 -type m4a
```

## TODO
- 完善部分功能
- 支持识别和写入更多元数据
- 添加 GUI


## 参考项目
感谢该项目提供的B站 API ，这对本项目有很大的帮助！  
https://github.com/SocialSisterYi/bilibili-API-collect
