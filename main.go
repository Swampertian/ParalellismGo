package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type INMET struct {
	Codigo string `json:"CD_OSCAR"`
}

type DADOS struct {
	Temp string `json:"TEMP_MIN"`
}

func fetch_endpoint1(ch chan<- []INMET) {
	start := time.Now()
	url := "https://apitempo.inmet.gov.br/estacoes/T"

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var estacoes []INMET
	json.Unmarshal(body, &estacoes)

	elapsed := time.Since(start)
	fmt.Printf("Tempo de execução de fetch_endpoint1: %v\n", elapsed)

	ch <- estacoes
}

func fetch_endpoint2(ch chan<- []DADOS) {
	start := time.Now()
	token := os.Getenv("INMET_TOKEN")
	if token == "" {
		panic("INMET_TOKEN environment variable not set")
	}
	url := "https://apitempo.inmet.gov.br/token/estacao/diaria/2022-11-01/2022-11-01/A001/" + token

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var dados []DADOS
	json.Unmarshal(body, &dados)

	elapsed := time.Since(start)
	fmt.Printf("Tempo de execução de fetch_endpoint2: %v\n", elapsed)

	ch <- dados
}

func main() {
	start := time.Now()
	ch1 := make(chan []INMET)
	ch2 := make(chan []DADOS)

	go fetch_endpoint1(ch1)
	go fetch_endpoint2(ch2)

	estacoes := <-ch1
	dados := <-ch2

	elapsed := time.Since(start)
	fmt.Printf("Tempo total de execução: %v\n", elapsed)

	fmt.Println("=== ESTACOES ===")
	for _, e := range estacoes {
		fmt.Println("Código:", e.Codigo)
	}
	fmt.Println("=== DADOS ===")
	for _, d := range dados {
		fmt.Println("Temp:", d.Temp)
	}
}
