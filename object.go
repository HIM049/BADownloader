package main

type Config struct {
	FavListLink string
	FileName    bool
}

type FavList struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Info struct {
			Title       string `json:"title"`
			Media_count int    `json:"media_count"`
		}
		Medias []struct {
			Id    int    `json:"id"`
			Type  int    `json:"type"`
			Title string `json:"title"`
			Cover string `json:"cover"`
			Page  int    `json:"page"`
			Bvid  string `json:"bvid"`
		}
	}
}

type VideoPageList struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		Cid int `json:"cid"`
	}
}

type VideoInformation struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Pic   string `json:"pic"`
		Pages []struct {
			Cid  int    `json:"cid"`
			Part string `json:"part"`
		}
	}
}

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
