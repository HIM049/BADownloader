# BiliAudioDownloader
一个用于批量下载B站收藏夹中视频音频的工具

---
常在B站音乐区混，攒了许多歌曲收藏。前段时间有制作本地曲库的想法，没有找到可以批量下载收藏夹的工具。于是便有了这个项目。
如果项目对你有帮助，请给我一个star!
项目写的比较杂乱，还请各位见谅。如果各位有不错的新功能或遇到了 BUG 欢迎提 issue 。

目前只提供了 Windows 的预编译版本。其他平台的用户可自行下载源码后使用 go 编译或直接运行。

### 已实现的功能：
- 收藏夹URL、收藏夹ID、BV号 下载
- 视频合集下载（可开关）
- 可自定义下载数量

### TODO:
- [x] 重构项目
- [ ] 自动写入歌曲信息到文件标签
- [ ] 多线程下载
- [ ] 可选择的音质

### 参考项目
感谢该项目提供的B站 API ，这对本项目有很大的帮助！  
https://github.com/SocialSisterYi/bilibili-API-collect
