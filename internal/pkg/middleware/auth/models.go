package auth

import (
	"encoding/json"
	"log"
	"net/http"
)

type GraphQLRequestBody struct {
	OperationName *string     `json:"operationName"`
	Query         *string     `json:"query"`
	Variables     interface{} `json:"variables"`
}

type CapturingResponseWriter struct {
	http.ResponseWriter
	CapturedBody []byte
	StatusCode   int
}

func (c *CapturingResponseWriter) WriteHeader(statusCode int) {
	c.ResponseWriter.WriteHeader(statusCode)
	c.StatusCode = statusCode
}

func (c *CapturingResponseWriter) Write(b []byte) (int, error) {
	type tokenResponse struct {
		Data struct {
			Login struct {
				Token string `json:"token"`
			} `json:"login"`
		} `json:"data"`
	}
	tokenResp := &tokenResponse{}

	err := json.Unmarshal(b, &tokenResp)
	if err != nil {
		log.Println(err)
		return c.ResponseWriter.Write(b)
	}

	http.SetCookie(c.ResponseWriter, &http.Cookie{Name: "Authorization", Value: tokenResp.Data.Login.Token})

	return c.ResponseWriter.Write(b)
}
