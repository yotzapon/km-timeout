package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	FAIL_DELAY           = 5000
	SUCESS_DELAY         = 500
	SUCCESS_DELAY_URL    = "https://httpstat.us/200?sleep=1000"
	FAIL_DELAY_URL       = "https://httpstat.us/504?sleep=6000"
	SUCCESS_NO_DELAY_URL = "https://httpstat.us/200?sleep=300000" // max delay 5minutes
)

type errorResponseModel struct {
	ErrorCode   string
	ErrorString string
	TimeElapse  interface{}
}

type responseModel struct {
	HttpResponse *http.Response
	Error        error
	TimeElapse   time.Duration
}

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(5)

	// the request has no timeout
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		rs := noTimeout()
		if rs.Error != nil {
			log.Printf("something went wrong: %+v", errorResponseModel{
				ErrorCode: rs.HttpResponse.Status,
				ErrorString: rs.Error.Error(),
				TimeElapse: rs.TimeElapse.Seconds(),
			})
		} else {
			log.Printf("success: noTimeout, took time(sec): %v", rs.TimeElapse.Seconds())
		}
	}(wg)

	// the request has configure timeout but the response from server within http client configured.
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		rs := timeoutWithHttpClientSuccess()
		if rs.Error != nil {
			log.Printf("something went wrong: %+v", errorResponseModel{
				ErrorString: rs.Error.Error(),
				TimeElapse: rs.TimeElapse.Seconds(),
			})
		} else {
			log.Printf("success: timeoutWithHttpClientSuccess, took time(sec): %v", rs.TimeElapse.Seconds())
		}
	}(wg)

	// the http client timeout before server response
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		rs := timeoutWithHttpClientFail()
		if rs.Error != nil {
			log.Printf("something went wrong: %+v", errorResponseModel{
				ErrorString: rs.Error.Error(),
				TimeElapse:  rs.TimeElapse.Seconds(),
			})
		} else {
			log.Printf("success: timeoutWithHttpClientFail, took time(sec): %v", rs.TimeElapse.Seconds())
		}
	}(wg)

	// the request from http client handle timeout by using context
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		rs := timeoutWithContext()
		if rs.Error != nil {
			log.Printf("something went wrong: %+v", errorResponseModel{
				ErrorString: rs.Error.Error(),
				TimeElapse: rs.TimeElapse.Seconds(),
			})
		} else {
			log.Printf("success: timeoutWithContext, took time(sec): %v", rs.TimeElapse.Seconds())
		}
	}(wg)

	// the request timeout with http transport
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		rs := timeoutWithHttpTransport()
		if rs.Error != nil {
			log.Printf("something went wrong: %+v", errorResponseModel{
				ErrorString: rs.Error.Error(),
				TimeElapse: rs.TimeElapse,
			})
		} else {
			log.Printf("success: timeoutWithHttpTransport, took time(sec): %v", rs.TimeElapse.Seconds())
		}
	}(wg)

	wg.Wait()
}

func noTimeout() *responseModel {
	log.Print("begin request by noTimeout method")
	start := time.Now()
	httpClient := &http.Client{}
	resp, err := httpClient.Get(SUCCESS_NO_DELAY_URL)
	end := time.Since(start)

	return &responseModel{
		HttpResponse: resp,
		Error:        err,
		TimeElapse:   end,
	}
}

func timeoutWithHttpClientFail() *responseModel {
	log.Print("begin request by timeoutWithHttpClientFail method")
	httpClient := &http.Client{
		Timeout: time.Duration(FAIL_DELAY) * time.Millisecond,
	}

	req, err := http.NewRequest(http.MethodGet, FAIL_DELAY_URL, nil)
	if err != nil {
		return nil
	}
	start := time.Now()
	resp, err := httpClient.Do(req)
	end := time.Since(start)
	return &responseModel{
		Error:        err,
		HttpResponse: resp,
		TimeElapse:   end,
	}
}

func timeoutWithHttpClientSuccess() *responseModel {
	log.Print("begin request by timeoutWithHttpClientSuccess method")
	httpClient := &http.Client{
		Timeout: time.Duration(FAIL_DELAY) * time.Millisecond,
	}

	req, err := http.NewRequest(http.MethodGet, SUCCESS_DELAY_URL, nil)
	if err != nil {
		return nil
	}
	start := time.Now()
	resp, err := httpClient.Do(req)
	end := time.Since(start)
	return &responseModel{
		Error:        err,
		HttpResponse: resp,
		TimeElapse:   end,
	}
}

func timeoutWithContext() *responseModel {
	log.Print("begin request by timeoutWithContext method")
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * FAIL_DELAY)
	defer cancel()

	httpClient := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, FAIL_DELAY_URL, nil)
	if err != nil {
		return nil
	}

	start := time.Now()
	resp, err := httpClient.Do(req.WithContext(ctx))
	end := time.Since(start)
	return &responseModel{
		Error: err,
		HttpResponse: resp,
		TimeElapse: end,
	}
}

func timeoutWithHttpTransport() *responseModel {
	log.Print("begin request by timeoutWithHttpTransport method")
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: FAIL_DELAY * time.Millisecond,
		}).DialContext,
	}

	httpClient := &http.Client{Transport: transport}
	req, err := http.NewRequest(http.MethodGet, FAIL_DELAY_URL, nil)
	if err != nil {
		return nil
	}

	start := time.Now()
	resp, err := httpClient.Do(req)
	end := time.Since(start)

	return &responseModel{
		Error: err,
		HttpResponse: resp,
		TimeElapse: end,
	}
}


// output
//2020/09/12 00:01:38 begin request by timeoutWithHttpClientSuccess method
//2020/09/12 00:01:38 begin request by noTimeout method
//2020/09/12 00:01:38 begin request by timeoutWithHttpClientFail method
//2020/09/12 00:01:38 begin request by timeoutWithHttpTransport method
//2020/09/12 00:01:38 begin request by timeoutWithContext method
//2020/09/12 00:01:40 success: timeoutWithHttpClientSuccess, took time(sec): 1.601113784
//2020/09/12 00:01:43 something went wrong: {ErrorCode: ErrorString:Get "https://httpstat.us/504?sleep=6000": context deadline exceeded (Client.Timeout exceeded while awaiting headers) TimeElapse:5.003266729}
//2020/09/12 00:01:43 something went wrong: {ErrorCode: ErrorString:Get "https://httpstat.us/504?sleep=6000": context deadline exceeded TimeElapse:5.003306626}
//2020/09/12 00:01:45 success: timeoutWithHttpTransport, took time(sec): 6.414960478
//2020/09/12 00:03:19 success: noTimeout, took time(sec): 100.415444659