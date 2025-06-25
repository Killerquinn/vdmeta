package models

import "sync"

type InstagramContent struct {
	mu     sync.Mutex
	IgUser string
	Type   string
	ID     string
}

type TikTokContent struct {
}

type PinterestContent struct {
}

type VimeoContent struct {
}

type YoutubeContent struct {
}
