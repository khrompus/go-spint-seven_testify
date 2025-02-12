package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func TestMainHandlerWhenRequestOk(t *testing.T) {
	req, err := http.NewRequest("GET", "/cafe?city=moscow&count=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status code 200")
	assert.NotNil(t, responseRecorder.Body)
	assert.NotEqual(t, "", responseRecorder.Body.String(), "expected body to be not empty")
}

func TestMainHandlerWhenCityNotEqual(t *testing.T) {
	req, err := http.NewRequest("GET", "/cafe?city=mos&count=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Expected status code 400")

	expectedBodyValue := "wrong city value"

	assert.Equal(t, expectedBodyValue, responseRecorder.Body.String(), "expected body value: wrong city value")
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req, err := http.NewRequest("GET", "/cafe?city=moscow&count=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	//Возможно здесь я сделал не правильно(костылем), в фидбек отпишите пожалуйста, правильно сделал или нет.
	expectedResponseLen := strings.Split(responseRecorder.Body.String(), ",")
	assert.Equal(t, totalCount, len(expectedResponseLen), "Expected total count to be %d", totalCount)

	//Такая же проверка на длину, если в запросе count > 4
	//expectedResponse := "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент"
	//assert.Equal(t, expectedResponse, responseRecorder.Body.String(), "Response body mismatch")
}
