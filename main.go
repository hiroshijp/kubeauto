package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"regexp"
)

type Data struct {
	MemUsage int `json:"mem_usage"`
	CpuUsage int `json:"cpu_usage"`
}

func main() {
	fs := http.FileServer(http.Dir("frontend"))
	http.Handle("/", fs)

	http.HandleFunc("/api", handler)

	http.ListenAndServe(":8080", nil)

}

func handler(w http.ResponseWriter, r *http.Request) {
	// ダミー
	data := &Data{
		MemUsage: 70,
		CpuUsage: 50,
	}
	log.Println(kubeTop())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func kubeTop() string {
	cmd := exec.Command("kubectl", "-n", "hcce", "top", "node")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(`(\d+)%`)
	matches := re.FindAllStringSubmatch(out.String(), -1)
	cpuPercent := matches[0][1]
	memoryPercent := matches[1][1]
	return cpuPercent + " " + memoryPercent
}
