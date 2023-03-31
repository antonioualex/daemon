package domain

import (
	"time"
)

type URL struct {
	URLstring           string      `json:"url"`
	Counter             int         `json:"counter"`
	CreatedAt           []time.Time `json:"created_at"`
	FailInitDownload    bool        `json:"-"`
	DownloadTime        string      `json:"download_time"`
	SuccessfulDownloads int         `json:"successful_downloads"`
	FailedDownloads     int         `json:"failed_downloads"`
}

const (
	ASC     string = "ASC"
	DESC    string = "DESC"
	BIGGEST string = "BIGGEST"
	LOWEST  string = "LOWEST"
	TOP     string = "TOP"
)

type BatchInsights struct {
	TotalSuccessfulDownloads int
	TotalFailedDownloads     int
	ElapsedTime              time.Duration
}

type URLService interface {
	Add(URl string) error
	Get(latest int, filter string) ([]*URL, error)
	Download(url *URL, insights *BatchInsights)
	CheckUrls()
}

type URLRepository interface {
	Get(latest int, filter string) ([]*URL, error)
	Insert(URL string) error
	Exists(URL string, insight bool) (bool, *URL)
}
