package policies

import (
	"fmt"
	"math"
	"math/rand"
	"os"
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
}

// Learn retorna um mapa de Playfield clusterizados para seu valor
// desse modos a funcao principal irá enquadrar cada tabuleiro nessa
// categoria e escolherá a que possui maior valor.
func Learn() func([]*tetris.Playfield) (*tetris.Playfield, int) {
	cluster_values := make(map[string]float64)
	// O cluster que representa um jogo perdido tem valor predefinido de -100
	cluster_values["L"] = -100.0

	// Chance de se explorar enquanto estiver aprendendo
	explore_chance := 0.75
	// Decaimento da taxa de exploracao
	var explore_decay float64
	explore_decay = 0.999

	// Porcentagem do valor de um único jogo que é incorporada ao valor
	// presente de um cluster
	taxa := 0.05

	// Total de jogos que serão jogados
	total_plays := 10000
	policy_from_cluster_learning :=
		func(tabuleiros []*tetris.Playfield) (*tetris.Playfield, int) {
			var melhor_peso float64
			melhor_peso = math.Inf(-1)
			pontos := 0
			// decide-se se vai explorar
			if rand.Float64() < explore_chance {
				// se sim, escolher um dos tabuleiros aleatoriamente
				tabuleiro_escolhido := tabuleiros[rand.Intn(len(tabuleiros))]
				pontos := tabuleiro_escolhido.RemoveCompletedLines()
				return tabuleiro_escolhido, pontos
			}
			// se não for explorar, passa por cada tabuleiro e
			// escolhe-se aquele que apresentar melhor valor de cluster
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

	results, err := os.Create("resultados.csv")
	if err != nil {
		println("Nao consegui abrir arquivo")
		return nil
	}
	var acc float64
	var avg float64
	acc = 0
	fmt.Sprintf("jogo,movimentos")
	for i := 1; i < total_plays+1; i++ {
		if i%10 == 0 {
			println(i)
			avg = acc / 10.0
			acc = 0
			fmt.Println("Explore:", explore_chance)
			fmt.Println("Media parcial:", avg)
			fmt.Fprintf(results, "%v,%v\n", i, avg)
		}
		moves, _, plays := tetris.Play(-1, policy_from_cluster_learning)
		acc += float64(moves)

		// Value representa o valor do jogo, que será adicionado ao cluster
		// que representa a última jogada e perpetuado para as jogadas anteriores
		//var value float64
		if explore_chance > 0.1 {
			explore_chance *= explore_decay
		}

		// O valor de um jogo é o seu número de pontos, dando um resultado negativo
		// para qualquer jogo cujo resultado foi menos que 50 movimentos
		value := (float64(moves) - avg/1.5)

		// Atualiza o valor de cada cluster que participou do jogo
		//println("Value:", value)
		taxa_temp := taxa
		for play_count := 0; play_count < len(plays)-1; play_count++ { //play_count := len(plays) - 1; play_count >= 0; play_count-- {
			clusterId := plays[play_count].GetClusterId()
			//fmt.Println("ClusterId:", clusterId)
			peso_atual := cluster_values[clusterId]
			//fmt.Println("PesoAtual:", peso_atual)
			cluster_values[clusterId] = (peso_atual * (1 - taxa_temp)) + (value * taxa_temp)
			//fmt.Println("Novo valor:", cluster_values[clusterId])
			value = cluster_values[clusterId]
			taxa_temp *= taxa
		}
	}

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
