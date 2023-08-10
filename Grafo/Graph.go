package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Vertice struct {
	Label string
	X, Y  int
}

type Aresta struct {
	Label string
	Cost  int
}

type Graph struct {
	Vertices map[string]Vertice
	Arestas  map[string]map[string]Aresta
}

func (g *Graph) newGraph() {
	g.Vertices = make(map[string]Vertice)
	g.Arestas = make(map[string]map[string]Aresta)
}

func (g *Graph) addVertex(label string, x, y int) {
	if _, exists := g.Vertices[label]; exists {
		fmt.Println("Erro: Vértice já existe.")
		return
	}

	g.Vertices[label] = Vertice{Label: label, X: x, Y: y}

	g.Arestas[label] = make(map[string]Aresta)

	fmt.Printf("Vértice %s inserido com sucesso.\n", label)
}

func (g *Graph) addEdge(from, to, label string, cost int) {
	_, fromExists := g.Vertices[from]
	_, toExists := g.Vertices[to]

	if !fromExists || !toExists {
		fmt.Println("Erro: Vértices de origem e/ou destino não existem.")
		return
	}

	if _, exists := g.Arestas[from][to]; exists {
		fmt.Println("Erro: Aresta já existe.")
		return
	}

	g.Arestas[from][to] = Aresta{
		Label: label,
		Cost:  cost,
	}

	if _, exists := g.Arestas[to][from]; !exists {
		g.Arestas[to][from] = Aresta{
			Label: label,
			Cost:  cost,
		}
	}

	fmt.Printf("Aresta de %s para %s inserida com sucesso.\n", from, to)
}

func (g *Graph) removeVertex(label string) {

	if _, exists := g.Vertices[label]; !exists {
		fmt.Println("Erro: Vértice não existe.")
		return
	}

	for from := range g.Arestas[label] {
		delete(g.Arestas[label], from)

		if _, exists := g.Vertices[from]; exists {
			delete(g.Arestas[from], label)
		}
	}

	delete(g.Vertices, label)

	fmt.Printf("Vértice %s removido com sucesso.\n", label)
}

func (g *Graph) removeEdge(from, to string) {
	_, fromExists := g.Vertices[from]
	_, toExists := g.Vertices[to]

	if !fromExists || !toExists {
		fmt.Println("Erro: Vértices de origem e/ou destino não existem.")
		return
	}

	if _, exists := g.Arestas[from][to]; !exists {
		fmt.Println("Erro: Aresta não existe.")
		return
	}

	delete(g.Arestas[from], to)
	delete(g.Arestas[to], from)

	fmt.Printf("Aresta de %s para %s removida com sucesso.\n", from, to)
}

func (g *Graph) printEdges() {
	fmt.Println("Lista de arestas do grafo:")

	for from, edges := range g.Arestas {
		for to, edge := range edges {
			fmt.Printf("%s --(%s: %d)--> %s\n", from, edge.Label, edge.Cost, to)
		}
	}
}

// #########################################
// funções código do josué adaptadas
// func euleriano
func (g *Graph) Degree(v string) int {
	degree := 0

	// Verificar se o vértice existe no grafo
	if _, exists := g.Vertices[v]; !exists {
		fmt.Println("Erro: Vértice não existe.")
		return -1
	}

	// Percorrer todas as arestas do vértice v
	edges, exists := g.Arestas[v]

	if exists {
		degree = len(edges)
	}

	return degree
}

/***********************************
 * Verifica se o grafo é Euleriano *
 **********************************/
func (g *Graph) IsEulerian() bool {
	for _, vertex := range g.Vertices {
		if g.Degree(vertex.Label)%2 == 1 {
			return false
		}
	}
	return true
}

//##########################################

func (g *Graph) DepthFirstSearch(startVertex string) {
	visited := make(map[string]bool)
	g.depthFirstSearchHelper(startVertex, visited)
}

func (g *Graph) depthFirstSearchHelper(currentVertex string, visited map[string]bool) {
	visited[currentVertex] = true

	fmt.Printf("%s ", currentVertex)

	edges := g.Arestas[currentVertex]

	for destinationVertex, _ := range edges {
		if !visited[destinationVertex] {
			g.depthFirstSearchHelper(destinationVertex, visited)
		}
	}
}

func Goodman(g *Graph) int {
	visited := make(map[string]bool)
	count := 0

	for vertex := range g.Vertices {
		if !visited[vertex] {
			g.depthFirstSearchHelper(vertex, visited)
			count++
		}
	}

	return count
}

func (g *Graph) BreadthFirstSearch(startVertex string) {
	if _, exists := g.Vertices[startVertex]; !exists {
		fmt.Println("Erro: O vértice informado não existe no grafo.")
		return
	}

	visited := make(map[string]bool)

	queue := []string{}

	visited[startVertex] = true
	queue = append(queue, startVertex)

	for len(queue) > 0 {
		currentVertex := queue[0]
		queue = queue[1:]

		fmt.Printf("%s ", currentVertex)

		edges := g.Arestas[currentVertex]

		for destinationVertex := range edges {
			if !visited[destinationVertex] {
				visited[destinationVertex] = true
				queue = append(queue, destinationVertex)
			}
		}
	}

	fmt.Println()
}

func (g *Graph) Dijkstra(source, destination string) {
	if _, sourceExists := g.Vertices[source]; !sourceExists {
		fmt.Println("Erro: Vértice de origem não existe no grafo.")
		return
	}
	if _, destinationExists := g.Vertices[destination]; !destinationExists {
		fmt.Println("Erro: Vértice de destino não existe no grafo.")
		return
	}

	distances := make(map[string]int)

	predecessors := make(map[string]string)

	visited := make(map[string]bool)

	for vertex := range g.Vertices {
		distances[vertex] = math.MaxInt32
	}
	distances[source] = 0

	for len(visited) < len(g.Vertices) {
		minDistance := math.MaxInt32
		var currentVertex string
		for vertex, distance := range distances {
			if !visited[vertex] && distance < minDistance {
				minDistance = distance
				currentVertex = vertex
			}
		}

		visited[currentVertex] = true

		edges := g.Arestas[currentVertex]
		for destinationVertex, edge := range edges {
			if !visited[destinationVertex] {
				newDistance := distances[currentVertex] + edge.Cost
				if newDistance < distances[destinationVertex] {
					distances[destinationVertex] = newDistance
					predecessors[destinationVertex] = currentVertex
				}
			}
		}
	}

	var path []string
	currentVertex := destination
	for currentVertex != source {
		path = append([]string{currentVertex}, path...)
		currentVertex = predecessors[currentVertex]
	}
	path = append([]string{source}, path...)

	fmt.Printf("Caminho de menor custo de %s para %s: %v\n", source, destination, path)
	fmt.Printf("Distância total: %d\n", distances[destination])
}

func (g *Graph) isBridge(u, v string) bool {
	edge := g.Arestas[u][v]
	delete(g.Arestas[u], v)
	delete(g.Arestas[v], u)

	visited := make(map[string]bool)
	g.depthFirstSearchHelper(u, visited)
	for vertex := range g.Vertices {
		if !visited[vertex] {
			g.Arestas[u][v] = edge
			g.Arestas[v][u] = Aresta{Label: edge.Label, Cost: edge.Cost}
			return false
		}
	}

	g.Arestas[u][v] = edge
	g.Arestas[v][u] = Aresta{Label: edge.Label, Cost: edge.Cost}

	return true
}

// Função para encontrar um ciclo euleriano em um grafo euleriano a partir de um vértice informado
func (g *Graph) Fleury(startVertex string) {
	// Verificar se o vértice informado existe no grafo
	if _, exists := g.Vertices[startVertex]; !exists {
		fmt.Println("Erro: O vértice informado não existe no grafo.")
		return
	}

	// Verificar se o grafo é euleriano
	if !g.IsEulerian() {
		fmt.Println("Erro: O grafo não é euleriano, portanto não possui um ciclo euleriano.")
		return
	}

	// Criar uma cópia do grafo para manipulação durante o algoritmo
	copyGraph := g.copyGraph()

	// Inicializar uma lista para armazenar o ciclo euleriano
	eulerianCycle := []string{}

	// Chamar a função auxiliar para encontrar o ciclo euleriano a partir do vértice informado
	g.fleuryHelper(startVertex, copyGraph, &eulerianCycle)

	// Imprimir o ciclo euleriano encontrado
	fmt.Println("Ciclo Euleriano:")
	for _, vertex := range eulerianCycle {
		fmt.Printf("%s ", vertex)
	}
	fmt.Println()
}

// Função auxiliar (recursiva) para encontrar o ciclo euleriano a partir de um vértice informado
func (g *Graph) fleuryHelper(currentVertex string, copyGraph *Graph, eulerianCycle *[]string) {
	// Adicionar o vértice atual ao ciclo euleriano
	*eulerianCycle = append(*eulerianCycle, currentVertex)

	// Enquanto houverem arestas no vértice atual
	for destinationVertex := range copyGraph.Arestas[currentVertex] {
		// Remover a aresta do vértice atual para o vértice de destino da cópia do grafo
		delete(copyGraph.Arestas[currentVertex], destinationVertex)

		// Remover a aresta do vértice de destino para o vértice atual da cópia do grafo
		delete(copyGraph.Arestas[destinationVertex], currentVertex)

		// Se o vértice de destino não ficar isolado após a remoção da aresta, continuar a busca a partir dele
		if copyGraph.Degree(destinationVertex) > 0 {
			g.fleuryHelper(destinationVertex, copyGraph, eulerianCycle)
		}
	}

	// Se não houverem mais arestas saindo do vértice atual, verificar se ainda existem arestas no grafo
	// (caso contrário, o ciclo está completo e a recursão termina)
	if copyGraph.Degree(currentVertex) == 0 {
		for vertex, edges := range copyGraph.Arestas {
			if len(edges) > 0 {
				g.fleuryHelper(vertex, copyGraph, eulerianCycle)
				break
			}
		}
	}
}

// Função para criar uma cópia do grafo
func (g *Graph) copyGraph() *Graph {
	copyGraph := &Graph{
		Vertices: make(map[string]Vertice),
		Arestas:  make(map[string]map[string]Aresta),
	}

	// Copiar os vértices
	for label, vertex := range g.Vertices {
		copyGraph.Vertices[label] = vertex
	}

	// Copiar as arestas
	for from, edges := range g.Arestas {
		copyGraph.Arestas[from] = make(map[string]Aresta)
		for to, edge := range edges {
			copyGraph.Arestas[from][to] = edge
		}
	}

	return copyGraph
}

func SaveGraphToFile(g *Graph, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	visitedEdges := make(map[string]bool)
	visitedEdgeLabels := make(map[string]bool)

	for _, vertex := range g.Vertices {
		line := fmt.Sprintf("Vertice: %s %d %d\n", vertex.Label, vertex.X, vertex.Y)
		writer.WriteString(line)
	}

	for _, vertex := range g.Vertices {
		for to, edge := range g.Arestas[vertex.Label] {
			edgeKey := fmt.Sprintf("%s-%s-%d", vertex.Label, to, edge.Cost)
			if !visitedEdges[edgeKey] && !visitedEdgeLabels[edge.Label] {
				line := fmt.Sprintf("Aresta: %s %s %s %d\n", vertex.Label, to, edge.Label, edge.Cost)
				writer.WriteString(line)
				visitedEdges[edgeKey] = true
				visitedEdgeLabels[edge.Label] = true
			}
		}
	}

	return writer.Flush()
}

func LoadGraphFromFile(filename string) (*Graph, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	graph := &Graph{}
	graph.newGraph()

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		switch parts[0] {
		case "Vertice:":
			if len(parts) != 4 {
				continue
			}
			label := parts[1]
			x, _ := strconv.Atoi(parts[2])
			y, _ := strconv.Atoi(parts[3])
			graph.addVertex(label, x, y)
		case "Aresta:":
			if len(parts) != 5 {
				continue
			}
			from := parts[1]
			to := parts[2]
			label := parts[3]
			cost, _ := strconv.Atoi(parts[4])
			graph.addEdge(from, to, label, cost)
		}
	}

	return graph, nil
}

func Clear() {
	fmt.Print("\033[H\033[2J")
}

func main() {
	var (
		x, y, op, custo                                   int
		name_vertice, vertice_c1, vertice_c2, name_aresta string
		graph                                             Graph
	)

	for {
		Clear()
		fmt.Println("manipulação e visualização de Grafos\n")
		fmt.Println("\t Menu Principal\n")
		fmt.Println("[ 0] Sair")
		fmt.Println("[ 1] Criar grafo")
		fmt.Println("[ 2] Imprimir arestas do grafo")
		fmt.Println("[ 3] Inserir novo vértice")
		fmt.Println("[ 4] Inserir nova aresta")
		fmt.Println("[ 5] Remover vértice")
		fmt.Println("[ 6] Remover aresta")
		fmt.Println("[ 7] Algoritmo de Goodman")
		fmt.Println("[ 8] Verificar grafo Euleriano") //se sim algoritmo de fleury
		fmt.Println("[ 9] Busca em Profundidade")
		fmt.Println("[10] Busca em Largura")
		fmt.Println("[11] Dijkstra ")
		fmt.Println("[12] Salvar em um arquivo")
		fmt.Println("[13] Ler de arquivo")

		fmt.Print("\nQual a sua opção? >> ")
		fmt.Scan(&op)
		switch op {
		case 0:
			{
				Clear()
				fmt.Println("Programa Encerrado!\nTecle [ENTER]")
				fmt.Scanln(&op)
				return
			}
		case 1:
			{
				Clear()
				fmt.Println("Criando Grafo...")
				graph.newGraph()
				fmt.Println("Grafo criada!\n Tecle [ENTER]")
				fmt.Scanln()
			}
		case 2:
			{
				Clear()
				fmt.Println("Grafo armazenada:\n")
				graph.printEdges()
				fmt.Println("\nTecle [ENTER]")
				fmt.Scanln()
			}
		case 3:
			{
				Clear()
				fmt.Println("Inserindo vértice\n")
				fmt.Print("Qual nome deseja dar ao vértice(único) >> ")
				fmt.Scanln(&name_vertice)
				fmt.Print("Quais as coordenadas(x,y) deseja dar ao vértice >> ")
				fmt.Scanln(&x, &y)
				graph.addVertex(name_vertice, x, y)
				fmt.Println("Vértice inserido.")
				graph.printEdges()
				fmt.Scanln()
			}
		case 4:
			{
				Clear()
				fmt.Println("Inserindo aresta\n")
				fmt.Print("Informe o nome de dois vértices ligados por esta aresta(existentes)>> ")
				fmt.Scanln(&vertice_c1, &vertice_c2)
				fmt.Print("Dê um nome para esta aresta>> ")
				fmt.Scanln(&name_aresta)
				fmt.Print("Qual o custo desta aresta>> ")
				fmt.Scanln(&custo)
				graph.addEdge(vertice_c1, vertice_c2, name_aresta, custo)
				fmt.Println("Aresta criada.")
				graph.printEdges()
				fmt.Scanln()
			}
		case 5:
			{
				Clear()
				fmt.Println("Removendo vértice\n")
				fmt.Print("Qual vértice deseja remover >> ")
				fmt.Scanln(&name_vertice)
				graph.removeVertex(name_vertice)
				fmt.Println("Vértice removido.")
				graph.printEdges()
				fmt.Scanln()
			}
		case 6:
			{
				Clear()
				fmt.Println("Removendo aresta\n")
				fmt.Print("Informe o  vértice de saída e do chegada da aresta que deseja remover >> ")
				fmt.Scanln(&vertice_c1, &vertice_c2)
				graph.removeEdge(vertice_c1, vertice_c2)
				fmt.Println("Aresta removida.")
				graph.printEdges()
				fmt.Scanln()
			}
		case 7:
			{
				Clear()
				fmt.Println("Realizando algoritmo de Goodman\n")
				fmt.Printf("A quantidade de fusões é: %d", Goodman(&graph))
				fmt.Println("\nTecle [ENTER]")
				fmt.Scanln()
			}
		case 8:
			{
				Clear()
				fmt.Println("Verificando se o Grafo é Euleriano\n")
				if graph.IsEulerian() {
					fmt.Println("O grafo inserido é Euleriano")
					fmt.Println("Realizando Fleury")
					fmt.Print("Digite um vértice para iniciar o algoritmo >> ")
					fmt.Scanln(&name_vertice)
					graph.Fleury(name_vertice)
				} else {
					fmt.Println("O grafo inserido não é Euleriano")
				}
				fmt.Println("\nTecle [ENTER]")
				fmt.Scanln()
			}
		case 9:
			{
				Clear()
				fmt.Println("Iniciando busca em Profundidade\n")
				fmt.Println("Digite um vértice para iniciar a busca >> \n")
				fmt.Scanln(&name_vertice)
				graph.DepthFirstSearch(name_vertice)
				graph.printEdges()
				fmt.Println("\nTecle [ENTER]")
				fmt.Scanln()
			}
		case 10:
			{
				Clear()
				fmt.Println("Iniciando busca em Largura\n")
				fmt.Println("Digite um vértice para iniciar a busca >> \n")
				fmt.Scanln(&name_vertice)
				graph.BreadthFirstSearch(name_vertice)
				graph.printEdges()
				fmt.Println("\nTecle [ENTER]")
				fmt.Scanln()
			}
		case 11:
			{
				Clear()
				fmt.Println("Iniciando Dijkstra \n")
				fmt.Println("Digite um vértice de inicio >> ")
				fmt.Scanln(&vertice_c1)
				fmt.Println("Digite um vértice final >> ")
				fmt.Scanln(&vertice_c2)
				graph.Dijkstra(vertice_c1, vertice_c2)
				fmt.Println("\nTecle [ENTER]")
				fmt.Scanln()
			}
		case 12:
			{
				Clear()
				fmt.Println("Salvando o grafo em um arquivo...")
				fmt.Print("Digite o nome do arquivo para salvar o grafo(com extensão): ")
				var filename string
				fmt.Scanln(&filename)
				err := SaveGraphToFile(&graph, filename)
				if err != nil {
					fmt.Printf("Erro ao salvar o arquivo: %s\n", err)
				} else {
					fmt.Println("Grafo salvo com sucesso!")
				}
				fmt.Println("\nTecle [ENTER]")
				fmt.Scanln()
			}
		case 13:
			{
				Clear()
				fmt.Println("Carregando o grafo de um arquivo...")
				fmt.Print("Digite o nome do arquivo para carregar o grafo(com extensão): ")
				var filename string
				fmt.Scanln(&filename)
				newGraph, err := LoadGraphFromFile(filename)
				if err != nil {
					fmt.Printf("Erro ao carregar o arquivo: %s\n", err)
				} else {
					graph = *newGraph
					fmt.Println("Grafo carregado com sucesso!")
				}
				fmt.Println("\nTecle [ENTER]")
				fmt.Scanln()
			}
		default:
			{
				fmt.Println("Opção Inválida!\nTecle [ENTER]")
				fmt.Scanln(&op)
			}
		}
	}
}
