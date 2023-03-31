package services

import (
	"daemon/domain"
	"daemon/persistence/fakes"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestURLService_Add(t *testing.T) {
	urlFakeRepo := &fakes.FakeURLRepository{}

	urlFakeRepo.InsertReturns(nil)

	urlService := NewURLService(urlFakeRepo)

	err := urlService.Add("http://example.com")
	if err != nil {
		t.Error(err)
	}
}

func TestURLService_Add_FAIL(t *testing.T) {
	urlFakeRepo := &fakes.FakeURLRepository{}

	urlFakeRepo.InsertReturns(errors.New("failed to add url"))

	urlService := NewURLService(urlFakeRepo)

	err := urlService.Add("http://example.com")
	if err == nil {
		t.Error(err)
	}
}

func TestURLService_Get(t *testing.T) {
	urlFakeRepo := &fakes.FakeURLRepository{}

	urlSlice := []*domain.URL{
		{
			URLstring:           "www.cat.com",
			Counter:             4,
			CreatedAt:           []time.Time{time.Date(2019, time.December, 10, 23, 0, 0, 0, time.UTC)},
			FailInitDownload:    false,
			DownloadTime:        "",
			SuccessfulDownloads: 2,
			FailedDownloads:     0,
		},
		{
			URLstring:           "www.dog.au",
			Counter:             7,
			CreatedAt:           []time.Time{time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)},
			FailInitDownload:    true,
			DownloadTime:        "",
			SuccessfulDownloads: 5,
			FailedDownloads:     0,
		},
		{
			URLstring:           "www.lion.au",
			Counter:             6,
			CreatedAt:           []time.Time{time.Date(2006, time.November, 10, 23, 0, 0, 0, time.UTC)},
			FailInitDownload:    true,
			DownloadTime:        "",
			SuccessfulDownloads: 1,
			FailedDownloads:     0,
		},
	}

	urlFakeRepo.GetReturns(urlSlice, nil)

	urlService := NewURLService(urlFakeRepo)

	urls, err := urlService.Get(5, domain.DESC)
	if err != nil {
		t.Error(err)
	}

	for key, url := range urls {
		if !reflect.DeepEqual(url, urlSlice[key]) {
			t.Error("urls do not match the urlSlice")
		}
	}
}

func TestURLService_Get_FAIL(t *testing.T) {
	urlFakeRepo := &fakes.FakeURLRepository{}

	urlFakeRepo.GetReturns(nil, errors.New("empty url list"))

	urlService := NewURLService(urlFakeRepo)

	urls, err := urlService.Get(5, domain.ASC)
	if err == nil {
		t.Error(err)
	}

	if len(urls) > 0 {
		t.Error("url list should be zero len")
	}
}
