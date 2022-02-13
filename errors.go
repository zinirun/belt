package belt

import (
	"log"
	"net/http"
)

func errorHandlerWithPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func errorHandlerWithFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkStatusCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status ", res.StatusCode)
	}
}
