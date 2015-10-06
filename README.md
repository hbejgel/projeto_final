# Projeto Final - 2015.2
## Aprendizado por Reforço para Jogar Tetris

O objetivo desse projeto é criar um jogador inteligente de Tetris.

Para isso é utilizado o conceito de **aprendizado por reforço**,
onde a  máquina aprende algo sem nenhuma assisteência externa, ou seja,
dada apenas um *valor* para a derrota e outro para a vitória,
que no caso é chegar a 200 movimentos sem perder, o algoritmo
aprende a jogar por si só, escolhendo novas jogadas tanto para explorar
novas oportunidades quanto para aproveitar oportunidades já identificadas
como positivas. Ao final desse processo de treinamento, temos uma polītica
que consegue escolher um tabuleiro dentre uma lista de opções.

Fontes utilizadas:

https://webdocs.cs.ualberta.ca/~sutton/book/the-book.html

http://www.cs.cmu.edu/afs/cs/project/ACRL/www/TetrisReports/Breelyn_Eric_Don_Project.pdf
