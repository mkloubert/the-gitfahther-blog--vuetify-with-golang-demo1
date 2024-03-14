// MIT License
//
// Copyright (c) 2024 Marcel Joachim Kloubert (https://marcel.coffee)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/georgecookeIW/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// storing the data of the POST request
type TextToSpeechRequestBody struct {
	Text  string `json:"text"`
	Voice string `json:"voice"`
}

func main() {
	OPENAI_API_KEY := os.Getenv("OPENAI_API_KEY")
	const port = "4000"

	router := fasthttprouter.New()

	// these settings are OK for this demo,
	// BUT NOT FOR PRODUCTION!
	//
	// learn more at: https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
	router.HandleOPTIONS = true
	router.HandleCORS.Handle = true
	router.HandleCORS.AllowOrigin = "*"
	router.HandleCORS.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

	// POST endpoint "/"
	router.POST("/", func(ctx *fasthttp.RequestCtx) {
		sendError := func(err error) {
			ctx.SetStatusCode(500)
			ctx.WriteString(err.Error())
		}

		// read data from request
		postBodyData := ctx.PostBody()

		var postBody TextToSpeechRequestBody
		err := json.Unmarshal(postBodyData, &postBody)
		if err != nil {
			sendError(err) // parsing input as JSON failed
			return
		}

		// setup request for OpenAI
		requestBody := map[string]interface{}{
			"model":           "tts-1",
			"input":           postBody.Text,
			"response_format": "opus",
			"voice":           postBody.Voice,
		}

		requestBodyData, err := json.Marshal(requestBody)
		if err != nil {
			sendError(err) // could not create JSON string from `requestBody`
			return
		}

		// start the request
		request, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/speech", bytes.NewBuffer(requestBodyData))
		if err != nil {
			sendError(err) // failed
			return
		}

		// setup headers like API key
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", OPENAI_API_KEY))

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			sendError(err) // could not doing the request
			return
		}

		responseBodyData, err := io.ReadAll(response.Body)
		if err != nil {
			sendError(err) // could not read the response
			return
		}

		if response.StatusCode != 200 {
			sendError(fmt.Errorf("unexpected response: %v", response.StatusCode))
			return
		}

		// now prepare response data
		// for returning as data URI
		base64ResponseBodyData := base64.StdEncoding.EncodeToString(responseBodyData)
		responseBodyDataURI := fmt.Sprintf("data:audio/ogg;base64,%s", base64ResponseBodyData)
		responseBodyDataURIData := []byte(responseBodyDataURI)

		ctx.SetStatusCode(200)
		ctx.Response.Header.Add("Content-Length", fmt.Sprint(len(responseBodyDataURIData)))
		ctx.Response.Header.Add("Content-Type", "text/plain; charset=UTF-8")
		ctx.Write(responseBodyDataURIData)
	})

	fmt.Println("Server now listening on port", port, "...")
	log.Fatal(fasthttp.ListenAndServe(":"+port, router.Handler))
}
