package persistence

import (
	domain "daemon/domain"
	"errors"
	"sort"
	"sync"
	"time"
)

type URLRepository struct {
	URLs []*domain.URL
	mu   *sync.Mutex
}

func (r *URLRepository) Exists(URL string, insight bool) (bool, *domain.URL) {
	if insight {
		r.mu.Lock()
		defer r.mu.Unlock()
	}
	if len(r.URLs) <= 0 {
		return false, &domain.URL{}
	}

	for index, v := range r.URLs {
		if v.URLstring == URL {
			if !insight {
				r.URLs[index].Counter += 1
				r.URLs[index].CreatedAt = append(r.URLs[index].CreatedAt, time.Now())
			}
			return true, r.URLs[index]
		}
	}

	return false, &domain.URL{}

}

func (r *URLRepository) Insert(url string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	exists, _ := r.Exists(url, false)

	if exists {
		return nil
	}

	r.URLs = append(r.URLs, &domain.URL{
		URLstring: url,
		Counter:   1,
		CreatedAt: []time.Time{
			time.Now(),
		},
	})

	return nil
}

func (r *URLRepository) Get(num int, filter string) ([]*domain.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.URLs) <= 0 {
		return nil, errors.New("no URLs available")
	}

	if len(r.URLs) < num {
		num = len(r.URLs)
	}

	switch filter {
	case domain.ASC: // ASC by Time
		sort.Slice(r.URLs, func(i, j int) bool {
			return r.URLs[i].CreatedAt[0].Before(r.URLs[j].CreatedAt[0])
		})
		return r.URLs[:num], nil
	case domain.DESC:
		sort.Slice(r.URLs, func(i, j int) bool {
			return r.URLs[i].CreatedAt[len(r.URLs[i].CreatedAt)-1].After(r.URLs[j].CreatedAt[len(r.URLs[j].CreatedAt)-1])
		})
		return r.URLs[:num], nil
	case domain.BIGGEST:
		sort.Slice(r.URLs, func(i, j int) bool {
			return len(r.URLs[i].URLstring) > len(r.URLs[j].URLstring)
		})
		return r.URLs[:num], nil
	case domain.LOWEST:
		sort.Slice(r.URLs, func(i, j int) bool {
			return len(r.URLs[i].URLstring) < len(r.URLs[j].URLstring)
		})
		return r.URLs[:num], nil
	case domain.TOP:
		sort.Slice(r.URLs, func(i, j int) bool {
			return r.URLs[i].Counter > r.URLs[j].Counter
		})
		return r.URLs[:num], nil
	default:
		return r.URLs[:num], nil
	}
}

func NewURLRepository() *URLRepository {
	URLs := []*domain.URL{}
	return &URLRepository{
		URLs: URLs,
		mu:   &sync.Mutex{},
	}
}
