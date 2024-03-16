package status

import (
	"io"
	"log"
	"net/http"

	"github.com/cbergoon/speedtest-go"
)

type NetworkStatus struct {
	Success           bool
	DownloadSpeedMbps float64
	UploadSpeedMbps   float64
	LatencyMs         float64
}

func GetNetworkInfo() NetworkStatus {
	user, _ := speedtest.FetchUserInfo()

	serverList, _ := speedtest.FetchServerList(user)
	targets, _ := serverList.FindServer([]int{})

	for _, s := range targets {
		s.PingTest()
		s.DownloadTest()
		s.UploadTest()
		return NetworkStatus{Success: true, LatencyMs: float64(s.Latency), DownloadSpeedMbps: s.DLSpeed, UploadSpeedMbps: s.ULSpeed}
	}
	return NetworkStatus{Success: false}
}
func GetLocationData() string {
	ipapiClient := http.Client{}

	req, err := http.NewRequest("GET", "https://ipapi.co/json/", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "ipapi.co/#go-v1.4.01")

	resp, err := ipapiClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}
