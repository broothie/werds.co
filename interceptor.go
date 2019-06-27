package main

import "net/http"

type Interceptor struct {
	w             http.ResponseWriter
	interceptFunc InterceptFunc
	statusCode    int
	header        http.Header
	body          []byte
}

type InterceptFunc func(*Interceptor)

func newInterceptor(w http.ResponseWriter, interceptFunc InterceptFunc) *Interceptor {
	return &Interceptor{
		w:             w,
		interceptFunc: interceptFunc,
		statusCode:    http.StatusOK,
		header:        make(map[string][]string),
	}
}

func (i *Interceptor) WriteHeader(statusCode int) {
	i.statusCode = statusCode
}

func (i *Interceptor) Header() http.Header {
	return i.header
}

func (i *Interceptor) Write(b []byte) (int, error) {
	i.body = b
	if i.interceptFunc != nil {
		i.interceptFunc(i)
	}
	return i.write()
}

func (i *Interceptor) write() (int, error) {
	for key, values := range i.header {
		for _, value := range values {
			i.w.Header().Add(key, value)
		}
	}

	if i.statusCode != http.StatusOK {
		i.w.WriteHeader(i.statusCode)
	}
	return i.w.Write(i.body)
}
