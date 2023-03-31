package services

import (
	"daemon/domain"
	"log"
	"sync"
	"time"
)

type URLService struct {
	URLRepository domain.URLRepository
	semaphore     chan string
}

func (s *URLService) Add(URl string) error {
	return s.URLRepository.Insert(URl)
}

func (s *URLService) Get(latest int, filter string) ([]*domain.URL, error) {
	return s.URLRepository.Get(latest, filter)
}

func (s *URLService) Download(url *domain.URL, bi *domain.BatchInsights) {

	select {
	case s.semaphore <- "locked":
		// Acquired a semaphore slot - execute the function logic
		defer func() {
			// Release the semaphore slot after the function has finished executing
			<-s.semaphore
		}()

		start := time.Now()
		s.URLRepository.Exists(url.URLstring, true)
		url.SuccessfulDownloads += 1
		elapsed := time.Since(start)
		url.DownloadTime = elapsed.String()
		bi.ElapsedTime += elapsed
		bi.TotalSuccessfulDownloads += 1
		return
	default:
		if !url.FailInitDownload {
			url.FailedDownloads += 1
			url.FailInitDownload = true
		}
		url.FailedDownloads += 1
		bi.TotalFailedDownloads += 1
		return
	}
}

func (s *URLService) CheckUrls() {
	urls, err := s.URLRepository.Get(10, domain.TOP)
	if err != nil {
		log.Println(err)
	}

	var wg sync.WaitGroup // create a wait group to synchronize goroutines
	bi := domain.BatchInsights{}
	for _, url := range urls {
		wg.Add(1) // increment wait group counter

		go func(url *domain.URL) {
			defer func() {
				wg.Done() // decrement wait group counter
			}()
			s.Download(url, &bi)
		}(url)
	}

	wg.Wait() // wait for all goroutines to finish
	log.Printf("Total Successfull Downloads: %v, Total Failed Downloads: %v, Total time: %v", bi.TotalSuccessfulDownloads, bi.TotalFailedDownloads, bi.ElapsedTime.String())
}

func NewURLService(URLRepository domain.URLRepository) *URLService {
	semaphore := make(chan string, 3)
	return &URLService{URLRepository, semaphore}
}
