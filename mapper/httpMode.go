package mapper

import (
	"io/ioutil"
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"github.com/jeffail/tunny"
	"runtime"
	"bytes"
	"strconv"
	"github.com/dgmann/ma-shared"
	"time"
)

type RegisterRequest struct {
	Url string `json:"url"`
	NumWorker int `json:"numWorker"`
}

type HttpMode struct {
	StageName string
	pool *tunny.WorkPool
}

func NewHttpMode(stageName string, pool *tunny.WorkPool) *HttpMode {
	return &HttpMode{StageName:stageName, pool:pool}
}

func(mode *HttpMode) Listen(endpoint string, setResult func(message *shared.Message, result interface{})) {
	http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()

		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}

		message := shared.NewMessageFromJSON(bodyBytes)
		message.EnterStage(mode.StageName)

		result, err := mode.pool.SendWork(message)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)

		setResult(message, result)
		message.LeaveStage(mode.StageName)

		jsonResponse, err := message.ToJSON()
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jsonResponse))
	})

	log.Fatal(http.ListenAndServe(":80", nil))
}

func(mode *HttpMode) Register(url string, host string, numWorkerScale int) (error) {
	fmt.Printf("Register at dispatcher: %s\r\n", url)

	client := http.Client{}
	registerRequest := RegisterRequest{
		Url: "http://" + host + "/image",
		NumWorker: runtime.NumCPU() * numWorkerScale,
	}
	jsonBytes, err := json.Marshal(registerRequest)
	if err != nil {
		println("Error marshalling register request.")
		return err
	}

	buffer := bytes.NewBuffer(jsonBytes)
	req, err := http.NewRequest("POST", "http://" + url +"/register", buffer)
	if err != nil {
		println("Error creating register request.")
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(buffer.Len()))
	defer req.Body.Close()

	res := send(client, req, url)
	defer res.Body.Close()
	println("Registered")

	return nil
}

func send(client http.Client, req *http.Request, url string) *http.Response {
	for {
		res, err := client.Do(req)

		if err == nil {
			return res
		}

		log.Println(err)
		log.Printf("Trying to connect to dispatcher at %s\n", url)
		time.Sleep(500 * time.Millisecond)
	}
}

