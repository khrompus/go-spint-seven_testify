package main

import (
	"strings"

	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Проверка запроса и тело ответа
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

// Проверка на правильно указанный город
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

// Проверка при запросе большего числа городов
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req, err := http.NewRequest("GET", "/cafe?city=moscow&count=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	//Возможно здесь я сделал не правильно(костылем), в фидбекэ отпишите пожалуйста, правильно сделал или нет.
	expectedResponseLen := strings.Split(responseRecorder.Body.String(), ",")
	assert.Equal(t, totalCount, len(expectedResponseLen), "Expected total count to be %d", totalCount)

	//Такая же проверка на длину, если в запросе count > 4
	//expectedResponse := "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент"
	//assert.Equal(t, expectedResponse, responseRecorder.Body.String(), "Response body mismatch")
}
