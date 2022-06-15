package main

import (
	"fmt"
	"math/rand"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 10 * time.Second
	targeter := NewCustomTargeter("GET")
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}

func NewCustomTargeter(option string) vegeta.Targeter {
	if option == "POST" {
		return func(tgt *vegeta.Target) error {
			if tgt == nil {
				return vegeta.ErrNilTarget
			}
			tgt.Method = "POST"
			tgt.URL = "http://localhost:8000/route"

			name := fmt.Sprintf("Name teste %v ", rand.Int63())
			cpf := rand.Intn(99999999999)
			password := rand.Float64()
			date := time.Now()
			payload := fmt.Sprintf(`{
										"name": "%v",
										"cpf":"%v",
										"password":"%v",
										"date":"%v"}`, name, cpf, password, date)
			tgt.Body = []byte(payload)
			tgt.Header = map[string][]string{
				"Authorization": []string{"Bearer ---Token---"},
				"Accept":        []string{"application/json"},
				"Content-Type":  []string{"application/json"},
			}
			return nil
		}
	} else if option == "GET" {
		return func(tgt *vegeta.Target) error {
			if tgt == nil {
				return vegeta.ErrNilTarget
			}
			tgt.Method = "GET"
			tgt.URL = "http://localhost:8000/route"

			tgt.Header = map[string][]string{
				"Authorization": []string{"Bearer ---Token---"},
			}
			return nil
		}
	} else {
		return nil
	}
}
