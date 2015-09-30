package policies

import (
	"fmt"
	"math"
	"math/rand"
	"projeto_final/tetris"
	//"time"
)

var Politicas map[string]func([]*tetris.Playfield) (*tetris.Playfield, int)

func init() {
	Politicas = map[string]func([]*tetris.Playfield) (*tetris.Playfield, int){
		"aleatoria":    Random_policy,
		"menor_altura": Least_height_policy,
		"quadratica":   Least_height_quadratic,
		//"reinforcement_learning": Reinforcement_learning(),
	}
}

func TrainReinf() {
	Politicas["reinforcement_learning"] = Reinforcement_learning()
}

func Random_policy(tabuleiros []*tetris.Playfield) (*tetris.Playfield, int) {
	selected_play := tabuleiros[rand.Intn(len(tabuleiros))]
	lines_removed := selected_play.RemoveCompletedLines()
	return selected_play, lines_removed
}

func Least_height_policy(tabuleiros []*tetris.Playfield) (*tetris.Playfield, int) {
	selected_play := tabuleiros[0]
	linhas := selected_play.RemoveCompletedLines()
	height := selected_play.FreeHeight()
	for _, tabuleiro := range tabuleiros {
		tmp_linhas := tabuleiro.RemoveCompletedLines()
		tmp_height := tabuleiro.FreeHeight()
		if tmp_height > height || ((tmp_height == height) && (tmp_linhas > linhas)) {
			selected_play = tabuleiro
			linhas = tmp_linhas
			height = tmp_height
		}
	}
	return selected_play, linhas
}

func Least_height_quadratic(tabuleiros []*tetris.Playfield) (*tetris.Playfield, int) {
	selected_play := tabuleiros[0]
	linhas := selected_play.RemoveCompletedLines()
	quadratic_height := selected_play.QuadraticHeight()
	for _, tabuleiro := range tabuleiros {
		tmp_linhas := tabuleiro.RemoveCompletedLines()
		tmp_quadratic := tabuleiro.QuadraticHeight()
		if tmp_quadratic < quadratic_height || ((tmp_quadratic == quadratic_height) && (tmp_linhas > linhas)) {
			selected_play = tabuleiro
			linhas = tmp_linhas
			quadratic_height = tmp_quadratic
		}
	}
	return selected_play, linhas
}

// Primeiro aprende a jogar tetris, logo em seguida cria uma política
// de como jogar baseado no aprendizado
func Reinforcement_learning() func([]*tetris.Playfield) (*tetris.Playfield, int) {
	return Learn()
	// return func(tabuleiros []*tetris.Playfield) (*tetris.Playfield, int) {
	// 	policy += 1
	// 	return tabuleiros[0], 0
	// }
}

// Learn retorna um mapa de Playfield clusterizados para seu valor
// desse modos a funcao principal irá enquadrar cada tabuleiro nessa
// categoria e escolherá a que possui maior valor.
func Learn() func([]*tetris.Playfield) (*tetris.Playfield, int) {
	cluster_values := make(map[string]float64)
	cluster_values["L"] = -100.0
	explore_chance := 0.5
	taxa := 0.1
	taxa_decay := 0.50
	total_plays := 1000
	policy_from_cluster_learning := func(tabuleiros []*tetris.Playfield) (*tetris.Playfield, int) {
		var melhor_peso float64
		melhor_peso = math.Inf(-1)
		pontos := 0
		if rand.Float64() < explore_chance { //vou explorar
			tabuleiro_escolhido := tabuleiros[rand.Intn(len(tabuleiros))]
			pontos := tabuleiro_escolhido.RemoveCompletedLines()
			return tabuleiro_escolhido, pontos
		}
		var tabuleiro_escolhido *tetris.Playfield
		for _, tabuleiro := range tabuleiros {
			p := tabuleiro.RemoveCompletedLines()
			id := tabuleiro.GetClusterId()
			tabuleiro.ClusterId = id
			if cluster_values[id] > melhor_peso {
				melhor_peso = cluster_values[id]
				tabuleiro_escolhido = tabuleiro
				pontos = p
			}
		}
		return tabuleiro_escolhido, pontos
	}
	for i := 0; i < total_plays; i++ {
		//start := time.Now()
		if i%100 == 0 {
			println(i)
		}
		moves, _, plays := tetris.Play(-1, policy_from_cluster_learning)
		var value float64
		explore_chance *= 0.99
		if moves > 500 {
			value = 10
		} else {
			value = -10
		}
		taxa_var := taxa
		for play_count := len(plays) - 1; play_count >= 0; play_count-- {
			clusterId := plays[play_count].GetClusterId()
			peso_atual := cluster_values[clusterId]
			cluster_values[clusterId] = (peso_atual * (1 - taxa_var)) + (value * taxa_var)
			taxa_var *= taxa_decay
		}
		//fmt.Println("Um jogo tomou:", time.Since(start))
	}

	// _, _, plays := tetris.Play(-1, policy_from_cluster)
	// println("ClusterID:", plays[1].GetClusterId())
	// cluster_values["21-1"] = -50
	// cluster_values["22-2"] = -50
	// _, _, plays = tetris.Play(-1, policy_from_cluster)
	// println("ClusterID:", plays[1].GetClusterId())
	println("Aprendi")
	for key, value := range cluster_values {
		fmt.Println("ClusterId: ", key, " Valor do Id: ", value)
	}
	println(len(cluster_values), " clusters diferentes")

	policy_from_cluster := func(tabuleiros []*tetris.Playfield) (*tetris.Playfield, int) {
		var melhor_peso float64
		melhor_peso = math.Inf(-1)
		pontos := 0
		var tabuleiro_escolhido *tetris.Playfield
		for _, tabuleiro := range tabuleiros {
			p := tabuleiro.RemoveCompletedLines()
			id := tabuleiro.GetClusterId()
			if cluster_values[id] > melhor_peso {
				melhor_peso = cluster_values[id]
				tabuleiro_escolhido = tabuleiro
				pontos = p
			}
		}
		return tabuleiro_escolhido, pontos
	}
	return policy_from_cluster
}
