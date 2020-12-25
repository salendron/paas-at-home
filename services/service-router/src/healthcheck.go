package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/gommon/log"
)

type HealthCheckInterface interface {
	DoHealthCheck(address string) (int, error)
	DoHealthChecks(storage StorageInterface)
}

type HealthCheck struct{}

func (hc *HealthCheck) DoHealthCheck(address string) (int, error) {
	start := time.Now()
	resp, err := http.Get(fmt.Sprintf("http://%v/_health", address))
	latency := int(time.Since(start).Milliseconds())
	if err != nil {
		return latency, err
	}
	if resp.StatusCode != 200 {
		return latency, fmt.Errorf("Service responded with HTTP %v", resp.StatusCode)
	}

	return latency, nil
}

func (hc *HealthCheck) DoHealthChecks(storage StorageInterface) {
	checkFunc := func() {
		services := storage.GetAllServices()
		for i := 0; i < len(services); i++ {
			latency, err := hc.DoHealthCheck(*(services[i].Address))
			isHealthy := true

			if err != nil {
				log.Errorf("Healtcheck failed for service %v/%v at %v", *(services[i].Name), *(services[i].Version), *(services[i].Address))
				isHealthy = false
			}

			services[i].IsHealthy = newBool(isHealthy)
			services[i].Latency = newInt(latency)

			err = storage.SaveService(services[i])
			if err != nil {
				log.Errorf("Could not saveHealtcheck for service %v/%v at %v", services[i].Name, services[i].Version, services[i].Address)
			}
		}
	}

	run := func() {
		log.Info("Starting HealthCheck routine...")
		healthCheckTicker := time.NewTicker(120 * time.Second)

		for {
			select {
			case <-healthCheckTicker.C:
				log.Info("Starting HealthChecks...")
				checkFunc()
				log.Info("HealthChecks done")
			}
		}
	}

	go run()
}
