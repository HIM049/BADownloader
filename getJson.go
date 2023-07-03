package main

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func GetFavList(favListId string, pn int) string {
	// 获取收藏夹视频列表。接收收藏夹ID和页码，返回json。
	params := url.Values{}
	Url, _ := url.Parse("https://api.bilibili.com/x/v3/fav/resource/list")
	params.Set("media_id", favListId)
	params.Set("ps", "1")
	params.Set("platform", "web")
	params.Set("pn", strconv.Itoa(pn))
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, err := http.Get(urlPath)
	CheckErr(err)
	body, _ := io.ReadAll(resp.Body)
	bodyString := string(body)
	defer resp.Body.Close()
	return bodyString
}

func GetVideoJson(bvid string, cid int) string {
	// 获取视频流URL。接收BV号和CID，返回视频元数据json。
	params := url.Values{}
	Url, _ := url.Parse("https://api.bilibili.com/x/player/playurl")
	params.Set("bvid", bvid)
	params.Set("cid", strconv.Itoa(cid))
	params.Set("fnval", "16")
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, err := http.Get(urlPath)
	CheckErr(err)
	body, _ := io.ReadAll(resp.Body)
	bodyString := string(body)
	defer resp.Body.Close()
	return bodyString
}

func GetVideoPageInformation(bvid string) string {
	// 获取视频流详细信息
	params := url.Values{}
	Url, _ := url.Parse("https://api.bilibili.com/x/web-interface/view")
	params.Set("bvid", bvid)
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, err := http.Get(urlPath)
	CheckErr(err)
	body, _ := io.ReadAll(resp.Body)
	bodyString := string(body)
	defer resp.Body.Close()
	return bodyString
}

func GetVideoPageList(bvid string) string {
	// 获取视频分P列表。接受BV号，返回视频分P json。
	params := url.Values{}
	Url, _ := url.Parse("https://api.bilibili.com/x/player/pagelist")
	params.Set("bvid", bvid)
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, err := http.Get(urlPath)
	CheckErr(err)
	body, _ := io.ReadAll(resp.Body)
	bodyString := string(body)
	defer resp.Body.Close()
	return bodyString
}
