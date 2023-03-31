package presentation

import (
	"daemon/app/fakes"
	"daemon/domain"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestURLHandler_Add(t *testing.T) {
	fakeURLService := &fakes.FakeURLService{}

	h := URLHandler{fakeURLService}

	reqBody := `{"url": "https://www.example.com"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/url/add", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	h.Add(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	expectedBody := `{"status":"URL added successfully"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %q, but got %q", expectedBody, w.Body.String())
	}
}

func TestURLHandler_Add_FAIL(t *testing.T) {
	fakeURLService := &fakes.FakeURLService{}
	fakeURLService.AddReturns(errors.New(`{"error":"Invalid url format"}`))

	h := URLHandler{fakeURLService}
	reqBody := `{"url": "https://example.com"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/url/add", strings.NewReader(reqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	h.Add(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, w.Code)
	}

	expectedBody := `{"error":"Invalid url format"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %q, but got %q", expectedBody, w.Body.String())
	}
}

func TestURLHandler_Retrieve(t *testing.T) {
	fakeURLService := &fakes.FakeURLService{}
	h := URLHandler{fakeURLService}

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
	}
	fakeURLService.GetReturns(urlSlice, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/url/retrieve", nil)
	h.Retrieve(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}
}

func TestURLHandler_Retrieve_FAIL(t *testing.T) {
	fakeURLService := &fakes.FakeURLService{}
	h := URLHandler{fakeURLService}

	fakeURLService.GetReturns(nil, errors.New("there are no urls available"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/url/retrieve", nil)
	h.Retrieve(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}
	expectedBody := `"there are no urls available"`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %q, but got %q", expectedBody, w.Body.String())
	}

}
