package main

import (
	"fmt"
	"os"
	"projeto_final/policies"
	"projeto_final/tetris"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Arquivo de entrada não encontrado")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Falha de leitura")
		os.Exit(1)
	}

	var input string
	_, err = fmt.Fscanln(file, &input)
	if err != nil {
		fmt.Println("Arquivo Vazio")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	start := time.Now()
	for true {
		read, err := fmt.Fscanln(file, &input)
		values := strings.Split(input, ",")
		if read == 0 || err != nil {
			break
		}

		if len(values) != 4 {
			fmt.Println("Inputs errados")
			fmt.Println(values)
			os.Exit(1)
		}
		politica, ok := policies.Politicas[values[1]]
		if !ok {
			if values[1] == "reinforcement_learning" {
				policies.TrainReinf()
			} else {
				fmt.Println("Política Invalida")
				os.Exit(1)
			}

		}
		games, err := strconv.Atoi(values[0])
		pieces, err := strconv.Atoi(values[2])
		seed, err := strconv.Atoi(values[3])

		fmt.Println(tetris.Play_series(int64(seed), pieces, games, politica))
	}
	elapsed := time.Since(start)
	fmt.Println("Everything took", elapsed)
	os.Exit(0)
}
