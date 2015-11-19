package tetris

import (
	"fmt"
	"github.com/mohae/deepcopy"
	"math"
	"math/rand"
)

type Coordinate struct{ x, y int }
type TetPosition [4]Coordinate
type Tetromino struct {
	position []TetPosition
}

type Playfield struct {
	// Largura do tabuleiro
	width int
	// Altura do tabuleiro
	height int
	// Vetor representando o tabuleiro
	field []int
	// Cluster no qual o tabuleiro se encaixa
	ClusterId string
}

const StandardWidth int = 10
const StandardHeight int = 20

/************************************
Square Tetromino

1 1  |
1 1  |
*************************************/
var tSquare = Tetromino{
	[]TetPosition{
		{{1, 1}, {0, 0}, {0, 1}, {1, 0}},
	}}

/************************************
Z Tetromino

1 1    |    1  |
  1 1  |  1 1  |
       |  1    |
*************************************/
var tJog1 = Tetromino{
	[]TetPosition{
		{{0, 0}, {1, 0}, {1, 1}, {2, 1}},
		{{1, 0}, {0, 1}, {1, 1}, {0, 2}},
	}}

/************************************
S Tetromino

  1 1  |  1    |
1 1    |  1 1  |
       |    1  |
*************************************/
var tJog2 = Tetromino{
	[]TetPosition{
		{{1, 0}, {2, 0}, {0, 1}, {1, 1}},
		{{0, 0}, {0, 1}, {1, 1}, {1, 2}},
	}}

/************************************
T Tetromino

  1    |    1    |         |    1
1 1 1  |    1 1  |  1 1 1  |  1 1
       |    1    |    1    |    1
*************************************/
var tTee = Tetromino{
	[]TetPosition{
		{{1, 1}, {0, 1}, {2, 1}, {1, 0}},
		{{1, 1}, {1, 2}, {2, 1}, {1, 0}},
		{{1, 1}, {0, 1}, {2, 1}, {1, 2}},
		{{1, 1}, {0, 1}, {1, 2}, {1, 0}},
	}}

/************************************
J Tetromino

  1 1  |         |    1  |  1
  1    |  1 1 1  |    1  |  1 1 1
  1    |      1  |  1 1  |
*************************************/
var tEl1 = Tetromino{
	[]TetPosition{
		{{1, 1}, {1, 0}, {2, 0}, {1, 2}},
		{{1, 1}, {0, 1}, {2, 1}, {2, 2}},
		{{1, 1}, {1, 0}, {0, 2}, {1, 2}},
		{{1, 1}, {0, 0}, {0, 1}, {2, 1}},
	}}

/************************************
L Tetromino

1 1  |      1  |    1    |
  1  |  1 1 1  |    1    |  1 1 1
  1  |         |    1 1  |  1
*************************************/
var tEl2 = Tetromino{
	[]TetPosition{
		{{1, 1}, {0, 0}, {1, 0}, {1, 2}},
		{{1, 1}, {2, 0}, {0, 1}, {2, 1}},
		{{1, 1}, {1, 0}, {1, 2}, {2, 2}},
		{{1, 1}, {0, 1}, {0, 2}, {2, 1}},
	}}

/************************************
Straight Tetromino

1  |
1  |  1 1 1 1
1  |
1  |
*************************************/
var tLong = Tetromino{
	[]TetPosition{
		{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
		{{0, 1}, {1, 1}, {2, 1}, {3, 1}},
	}}

const (
	square = iota
	jog1
	jog2
	tee
	el1
	el2
	long
)

var Tetrominos = []Tetromino{
	tSquare,
	tJog1,
	tJog2,
	tTee,
	tEl1,
	tEl2,
	tLong,
}

// Cria um campo novo
func NewPlayfield(width, height int) Playfield {
	var p Playfield
	p.width = width
	p.height = height + 4
	p.field = make([]int, p.width*p.height)
	return p
}

// Copia um campo
func (p *Playfield) deepCopy() Playfield {
	var pf Playfield
	pf.width = p.width
	pf.height = p.height
	pf.field = deepcopy.IntSlice(p.field)
	return pf
}

// Retorna o quadrado na posicao (x,y) , -1 caso posicao invalida
func (p *Playfield) at(x, y int) int {
	if x < 0 || x >= p.width || y < 0 || y >= p.height {
		return -1
	}
	return p.field[x+y*p.width]
}

// Checa se uma linha está completa
func (p *Playfield) isLineComplete(y int) bool {
	for i := 0; i < p.width; i++ {
		if p.at(i, y) == 0 {
			return false
		}
	}
	return true
}

// Checa se uma linha está ocupada por alguma peça
func (p *Playfield) isLineOccupied(y int) bool {
	for i := 0; i < p.width; i++ {
		if p.at(i, y) != 0 {
			return true
		}
	}
	return false
}

// Retorna a primeira linha ocupada
func (p *Playfield) FreeHeight() int {
	var i int
	for i = 0; i < p.height; i++ {
		if p.isLineOccupied(i) {
			break
		}
	}
	return i
}

// Conta o número de buracos em um tabuleiro
func (p *Playfield) Holes() int {
	total_holes := 0
	for y := 1; y < p.height; y++ {
		for x := 0; x < p.width; x++ {
			if p.at(x, y) == 0 && p.at(x, y-1) != 0 {
				total_holes += 1
			}
		}
	}
	return total_holes
}

// Checa se o jogo foi perdido
func (p *Playfield) lost() bool {
	if p.FreeHeight() < 4 {
		return true
	}
	return false
}

// Remove uma linha, criando uma nova linha vazia no topo, e fazendo cada linha acima da linha removida
// descer uma posicao no tabuleiro
func (p *Playfield) removeLine(y int) {
	if y > 0 {
		for i := y; i > 0; i-- {
			for j := 0; j < p.width; j++ {
				p.set(j, i, p.at(j, i-1))
			}
		}
	}
	for j := 0; j < p.width; j++ {
		p.field[j] = 0
	}
}

// Passa pelo tabuleiro removendo linhas completas
func (p *Playfield) RemoveCompletedLines() int {
	count := 0
	for i := 0; i < p.height; i++ {
		if p.isLineComplete(i) {
			p.removeLine(i)
			count++
		}
	}
	return count
}

// Marca uma posicao do tabuleiro como ocupada
func (p *Playfield) set(x, y, v int) bool {
	if x < 0 || x >= p.width || y < 0 || y >= p.height {
		return false
	}
	p.field[x+y*p.width] = v
	return true
}

// Imprime o tabuleiro
func (p *Playfield) Print() {
	for i := 0; i < p.height; i++ {
		fmt.Print("|")
		for j := 0; j < p.width; j++ {
			v := p.at(j, i)
			if v == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print(p.at(j, i))
			}
		}
		fmt.Println("|")
	}
	fmt.Println("ClusterId:", p.GetClusterId())
}

// Coloca uma peca no tabuleiro
func (p *Playfield) place(piece *Piece) bool {
	for i := 0; i < 4; i++ {
		r := piece.rot
		placed := p.set(
			piece.pos.x+piece.tet.position[r][i].x,
			piece.pos.y+piece.tet.position[r][i].y, 1)
		if placed == false {
			return false
		}
	}
	return true
}

// Uma peca possui uma coordenada, rotacao e o tetormino que ela representa
type Piece struct {
	tet      *Tetromino
	tet_code int
	pos      Coordinate
	rot      int
}

func (p *Piece) deepCopy() *Piece {
	var pi Piece
	pi.tet = p.tet
	pi.tet_code = p.tet_code
	pi.pos.x = p.pos.x
	pi.pos.y = p.pos.y
	pi.rot = p.rot
	return &pi
}

func (p *Piece) Print() {
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			for _, pos := range p.tet.position[p.rot] {
				if pos.x == x && pos.y == y {
					print(1)
				}
			}
		}
		print("\n")
	}
}

func (p *Piece) to_s() string {
	return fmt.Sprintf("x:%v, y:%v, rot:%v", p.pos.x, p.pos.y, p.rot)
}

// Move uma peca no tabuleiro, retornando true se ela pode ser movida e false se ela nao pode ser movida
func (p *Piece) move(field *Playfield, xi, yi int) bool {
	if p == nil {
		return false
	}
	tx := p.pos.x + xi
	ty := p.pos.y + yi
	for i := 0; i < 4; i++ {
		x := tx + p.tet.position[p.rot][i].x
		y := ty + p.tet.position[p.rot][i].y
		if field.at(x, y) != 0 {
			return false
		}
	}
	p.pos.x = tx
	p.pos.y = ty
	return true
}

// Roda uma peca no tabuleiro, retornando true se ela pode ser rodada e false se ela nao pode ser rodada
func (p *Piece) rotate(field *Playfield, ri int) bool {
	if p == nil {
		return false
	}
	var tr int
	if p.tet_code == square {
		return true
	} else if p.tet_code == long || p.tet_code == jog1 || p.tet_code == jog2 {
		tr = int(uint(p.rot+ri) % 2)
	} else {
		tr = int(uint(p.rot+ri) % 4)
	}
	for i := 0; i < 4; i++ {
		x := p.pos.x + p.tet.position[tr][i].x
		y := p.pos.y + p.tet.position[tr][i].y
		if field.at(x, y) != 0 {
			return false
		}
	}
	p.rot = tr
	return true
}

// Cria uma peca nova
func makeNewTet(code int) *Piece {
	if code < 0 || code > len(Tetrominos)-1 {
		code = rand.Intn(len(Tetrominos))
	}
	tet := &Tetrominos[code]
	return &Piece{tet, code, Coordinate{0, 0}, 0}
}

// bfs, dado um campo e uma peca, retorna uma lista de campos finais com aquela peca
func (pf *Playfield) bfs(piece *Piece) []*Playfield {
	frontier := []*Piece{piece}
	visited := map[string]bool{}
	next := []*Piece{}
	solutions := []*Playfield{}
	visited[piece.to_s()] = true
	for len(frontier) > 0 {
		next = []*Piece{}
		for _, p := range frontier {
			tmp := p.deepCopy()
			if !tmp.move(pf, 0, 1) && tmp.move(pf, 0, 0) {
				field_solution := pf.deepCopy()
				field_solution.place(tmp)
				solutions = append(solutions, &field_solution)
			}
			for _, neighbour := range pf.bfs_frontier(p, visited) {
				next = append(next, neighbour)
				visited[neighbour.to_s()] = true
			}
		}
		frontier = next
	}

	return solutions
}

// Dado uma peca, um tabuleiro e a lista de estados visitados, retorna as proximas posicoes da peca no tabuleiro
func (pf *Playfield) bfs_frontier(piece *Piece, visited map[string]bool) []*Piece {
	next := []*Piece{}
	p_left := piece.deepCopy()
	if p_left.move(pf, -1, 0) && !visited[p_left.to_s()] {
		next = append(next, p_left)
	}
	p_right := piece.deepCopy()
	if p_right.move(pf, 1, 0) && !visited[p_right.to_s()] {
		next = append(next, p_right)
	}
	p_down := piece.deepCopy()
	if p_down.move(pf, 0, 1) && !visited[p_down.to_s()] {
		next = append(next, p_down)
	}
	p_rotate := piece.deepCopy()
	if p_rotate.rotate(pf, -1) && !visited[p_rotate.to_s()] {
		next = append(next, p_rotate)
	}
	return next
}

// Joga uma partida de tetris dado uma semente aleatoria, as peças que participam da partida e uma política para escolha de peças
func Play(piece_code int, policy func(tabuleiros []*Playfield) (*Playfield, int)) (int, int, []Playfield) {
	playfield := NewPlayfield(StandardWidth, StandardHeight)
	moves := 0
	points := 0
	var plays []Playfield
	plays = append(plays, playfield)
	for true {
		piece := makeNewTet(piece_code)
		piece.pos.x = (playfield.width / 2) - 2
		outcomes := playfield.bfs(piece)
		if len(outcomes) < 1 {
			break
		}

		p, score := policy(outcomes)
		playfield = *p
		plays = append(plays, playfield)
		points += score
		moves += 1
		if p.lost() {
			break
		}

	}
	return moves, points, plays
}

func Play_series(random_seed int64, piece_code, games int, policy func(tabuleiros []*Playfield) (*Playfield, int)) (float64, float64) {
	var total_moves, total_points float64
	rand.Seed(random_seed)
	for i := 0; i < games; i++ {
		moves, points, _ := Play(piece_code, policy)
		total_moves += float64(moves)
		total_points += float64(points)
	}

	return total_moves / float64(games), total_points / float64(games)
}

func (p *Playfield) QuadraticHeight() int {
	var total float64
	for i := 0; i < p.width; i++ {
		for j := 0; j < p.height; j++ {
			if p.at(i, j) != 0 {
				total += math.Pow(float64(p.height-j), 2)
			}
		}
	}
	return int(total)
}

func (p *Playfield) quadraticHeight() int {
	var total float64
	for i := 0; i < p.width; i++ {
		for j := 0; j < p.height; j++ {
			if p.at(i, j) != 0 {
				total += math.Pow(float64(p.height-j), 2)
			}
		}
	}
	return int(total)
}

func (p *Playfield) GetClusterId() string {
	if p.ClusterId != "" {
		return p.ClusterId
	}
	freeHeight := p.FreeHeight()
	if freeHeight < 4 {
		return "L"
	}
	return fmt.Sprintf("%v-%v", freeHeight, p.Holes())
}

// Converte uma quantidade de linhas em pontos
func points_from_lines(lines int) int {
	switch lines {
	case 1:
		return 40
	case 2:
		return 100
	case 3:
		return 300
	case 4:
		return 1200
	}
	fmt.Println("Scoring Error")
	return 0
}
