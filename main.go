package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"time"
	"log"
	"net/http"
	"sync"
)

func addCORSHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Permite qualquer origem
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func isDeviceActive(ip string, timeout time.Duration) bool {
	cmd := exec.Command("ping", "-n", "1", "-w", fmt.Sprintf("%d", int(timeout.Seconds()*1000)), ip)
	err := cmd.Run()
	return err == nil
}

func getDeviceName(ip string) string {
	names, err := net.LookupAddr(ip)
	if err != nil || len(names) == 0 {
		return "Desconhecido"
	}
	return names[0]
}

func checkDevice(ip string, timeout time.Duration, wg *sync.WaitGroup, activeDevices chan<- map[string]string) {
	defer wg.Done() 

	if isDeviceActive(ip, timeout) {
		deviceName := getDeviceName(ip)
		activeDevices <- map[string]string{"ip": ip, "name": deviceName}
	}
}

func scanNetwork(w http.ResponseWriter, r *http.Request) {
	addCORSHeaders(w, r)

	ipBase := r.URL.Query().Get("ip")
	if ipBase == "" {
		http.Error(w, "Por favor, forneça um IP base para escanear", http.StatusBadRequest)
		return
	}

	timeout := 1 * time.Second

	var activeDevices []map[string]string

	activeDevicesChan := make(chan map[string]string)

	var wg sync.WaitGroup

	for i := 1; i <= 255; i++ {
		ip := fmt.Sprintf("%s.%d", ipBase, i)
		wg.Add(1) 
		go checkDevice(ip, timeout, &wg, activeDevicesChan) // Dispara uma goroutine para cada IP
	}

	go func() {
		wg.Wait()
		close(activeDevicesChan)
	}()

	for device := range activeDevicesChan {
		activeDevices = append(activeDevices, device)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activeDevices)
}

func main() {
	http.HandleFunc("/scan", scanNetwork)

	port := "8080"
	fmt.Printf("Iniciando servidor na porta %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
