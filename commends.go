package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

func SetupCommands() {
	app := cli.NewApp()
	app.Name = "BADownloader"
	app.Usage = "一个用于下载 BiliBili 视频音频的工具"

	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			fmt.Println("您没有输入任何参数。 请输入 '-h' 获取帮助。")
			return nil
		}
		return nil
	}

	app.Commands = []cli.Command{
		// {
		// 	Name:    "config",
		// 	Aliases: []string{"cfg"},
		// 	Usage:   "设置相关",
		// 	Subcommands: []cli.Command{
		// 		{
		// 			Name:  "index",
		// 			Usage: "设置引导",
		// 			Action: func() {
		// 				cfg, err := ConfigIndex()
		// 				if err != nil {
		// 					fmt.Println("ConfigIndex: ", err)
		// 				}
		// 				err = SaveJsonFile(CONFIG_PATH, cfg)
		// 				if err != nil {
		// 					fmt.Println("SaveJsonFile: ", err)
		// 				}
		// 			},
		// 		},
		// 	},
		// },
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "关于任务列表",
			Subcommands: []cli.Command{
				{
					Name:    "bulid",
					Aliases: []string{"b"},
					Usage:   "添加指定收藏夹视频到列表",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "id",
							Usage: "收藏夹 ID",
						},
						&cli.IntFlag{
							Name:  "c",
							Usage: "下载数量设置",
						},
						&cli.BoolFlag{
							Name:  "p",
							Usage: "下载分 P 设置",
						},
					},
					Action: func(c *cli.Context) error {
						if c.IsSet("id") && c.IsSet("c") && c.IsSet("p") {
							// 参数齐全
							id := c.String("id")
							count := c.Int("c")
							downloadPage := c.Bool("p")

							fmt.Printf("正在创建任务列表：%s 视频数量：%d 下载分 P：%t\n", id, count, downloadPage)
							err := MakeAndSaveList(VIDEO_LIST_PATH, id, count, downloadPage)
							if err != nil {
								fmt.Println("MakeAndSaveList: ", err)
							}
						} else {
							fmt.Println("参数不足")
						}
						return nil
					},
				},
			},
		},
		{
			Name:    "download",
			Aliases: []string{"d"},
			Usage:   "执行下载任务",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "t",
					Usage: "下载线程数",
				},
			},
			Action: func(c *cli.Context) {
				if c.IsSet("t") {
					fmt.Printf("正在下载音频：线程数：%d\n", c.Int("t"))
					err := DownloadList(c.Int("t"))
					if err != nil {
						fmt.Println("DownloadList: ", err)
					}
					fmt.Printf("正在下载封面图：线程数：%d\n", c.Int("t"))
					err = ConcurrentSavePic(c.Int("t"), COVER_PATH)
					if err != nil {
						fmt.Println("ConcurrentSavePic: ", err)
					}

				} else {
					fmt.Println("参数不足")
				}
			},
		},
		{
			Name:  "tag",
			Usage: "为歌曲添加元数据",
			Subcommands: []cli.Command{
				{
					Name:    "add",
					Aliases: []string{"a"},
					Usage:   "添加 TAG",
					Flags: []cli.Flag{
						&cli.IntFlag{
							Name:  "t",
							Usage: "下载线程数",
						},
						&cli.StringFlag{
							Name:  "type",
							Usage: "音频类型",
						},
					},
					Action: func(c *cli.Context) {
						if c.IsSet("t") && c.IsSet("type") {
							path := MP3_PATH
							if c.String("type") == "m4a" {
								path = M4A_PATH
							}
							fmt.Println("正在添加元数据")
							err := ConcurrentChangeTag(c.Int("t"), path, COVER_PATH, "."+c.String("type"))
							if err != nil {
								fmt.Println("ConcurrentChangeTag: ", err)
							}
						} else {
							fmt.Println("参数不足")
						}
					},
				},
				// {
				// 	Name:    "remove",
				// 	Aliases: []string{"r"},
				// 	Usage:   "移除 TAG",
				// },
			},
		},
		{
			Name:    "conver",
			Aliases: []string{"con"},
			Usage:   "批量转码歌曲",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "t",
					Usage: "下载线程数",
				},
				&cli.StringFlag{
					Name:  "p",
					Usage: "ffmpeg 路径",
				},
			},
			Action: func(c *cli.Context) {
				if c.IsSet("t") && c.IsSet("p") {
					fmt.Println("正在转码音频")
					err := ConcurrentToMp3(c.Int("t"), c.String("p"), M4A_PATH, MP3_PATH, LOG_PATH)
					if err != nil {
						fmt.Println("ConcurrentToMp3", err)
					}
				} else {
					fmt.Println("参数不足")
				}
			},
		},
		{
			Name:    "out",
			Aliases: []string{"o"},
			Usage:   "设置为完成状态，输出音频并清理文件",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "t",
					Usage: "下载线程数",
				},
				&cli.StringFlag{
					Name:  "type",
					Usage: "音频类型",
				},
			},
			Action: func(c *cli.Context) {
				if c.IsSet("t") && c.IsSet("type") {
					path := MP3_PATH
					if c.String("type") == "m4a" {
						path = M4A_PATH
					}
					fmt.Println("正在导出音频")
					err := ConcurrentChangeName(c.Int("t"), "."+c.String("type"), path, OUT_PATH)
					if err != nil {
						fmt.Println("ConcurrentToMp3", err)
					}
					err = os.RemoveAll(CACHE_PATH)
					if err != nil {
						fmt.Println("清理缓存数据时发生错误：", err)
						fmt.Println("您可以手动删除目录下的 “ Cache ” 文件夹来清除缓存文件")
					}
				} else {
					fmt.Println("参数不足")
				}
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
