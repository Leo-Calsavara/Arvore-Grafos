package main

import (
	"fmt"
	"os"
	"strings"
)

//==========================================================
//			DECLARAÇÃO DAS ESTRUTURAS E TIPOS
//==========================================================

type Contato struct {
	Nome     string
	Endereco string
	Telefone string
	Apagado  bool
}

const t = 2

type Lista struct {
	Contatos []Contato
}

type DataType struct {
	nome  string
	index int
}

type BTreeNode struct {
	leaf     bool
	keys     []DataType
	children []*BTreeNode
}

type BTree struct {
	root *BTreeNode
}

var tree *BTree

//========================================================
//				FUNÇÕES SOBRE ARVORES
//========================================================

func Init() *BTree {

	return &BTree{
		root: InitNode(true),
	}
}

func InitNode(leaf bool) *BTreeNode {
	return &BTreeNode{
		leaf:     leaf,
		keys:     []DataType{},
		children: []*BTreeNode{},
	}
}

func (node *BTreeNode) Insert(key DataType) {

	if !node.leaf {
		i := len(node.keys) - 1
		for i >= 0 && key.nome < node.keys[i].nome {
			i--
		}

		if len(node.children[i+1].keys) == 2*t-1 {
			node.splitChild(int16(i) + 1)
			if key.nome > node.keys[i+1].nome {
				i++
			}
		}
		node.children[i+1].Insert(key)
	} else {
		empty := DataType{"", 0}
		i := len(node.keys) - 1
		node.keys = append(node.keys, empty)
		for i >= 0 && key.nome < node.keys[i].nome {
			node.keys[i+1] = node.keys[i]
			i--
		}
		node.keys[i+1] = key
	}
}


func (tree *BTree) Insert(key DataType) {
	root := tree.root
	if len(root.keys) == 2*t-1 {
		newRoot := InitNode(false)
		newRoot.children = append(newRoot.children, root)
		newRoot.splitChild(0)
		tree.root = newRoot
	}
	tree.root.Insert(key)
}

func (node *BTreeNode) splitChild(i int16) {
	child := node.children[i]
	newChild := InitNode(child.leaf)

	newChild.keys = append(newChild.keys, child.keys[t:]...)
	child.keys = child.keys[:t]
	if !child.leaf {
		newChild.children = append(newChild.children, child.children[t:]...)
		child.children = child.children[:t]
	}

	node.children = append(node.children, nil)
	copy(node.children[i+2:], node.children[i+1:])
	node.children[i+1] = newChild

	empty := DataType{"", 0}

	node.keys = append(node.keys, empty)
	copy(node.keys[i+1:], node.keys[i:])
	node.keys[i] = child.keys[t-1]
	child.keys = child.keys[:t-1]
}

func (tree *BTree) PercursoEmOrdem() {
	fmt.Println("Percurso em Ordem:")
	result := tree.root.percursoEmOrdem()
	fmt.Println(result)
}

func (node *BTreeNode) percursoEmOrdem() string {
	if node == nil {
		return ""
	}

	var result strings.Builder

	if !node.leaf {
		result.WriteString(node.children[0].percursoEmOrdem())
	}

	for i := 0; i < len(node.keys); i++ {
		result.WriteString(fmt.Sprintf("%v ", node.keys[i]))
		if !node.leaf {
			result.WriteString(node.children[i+1].percursoEmOrdem())
		}
	}

	return result.String()
}


func (node *BTreeNode) Remove(key string) {
	i := 0
	for i < len(node.keys) && key > node.keys[i].nome {
		i++
	}

	if i < len(node.keys) && key == node.keys[i].nome {
		if node.leaf {
			copy(node.keys[i:], node.keys[i+1:])
			node.keys = node.keys[:len(node.keys)-1]
		} else {
			predecessor := node.children[i]
			for !predecessor.leaf {
				predecessor = predecessor.children[len(predecessor.children)-1]
			}
			node.keys[i] = predecessor.keys[len(predecessor.keys)-1]
			predecessor.Remove(node.keys[i].nome)
		}
	} else {
		if node.leaf {
			fmt.Println("A chave não está presente na árvore.")
			return
		}

		if len(node.children[i].keys) < t {
			child := node.children[i]
			if i > 0 && len(node.children[i-1].keys) >= t {
				leftSibling := node.children[i-1]
				child.keys = append([]DataType{node.keys[i-1]}, child.keys...)
				node.keys[i-1] = leftSibling.keys[len(leftSibling.keys)-1]
				leftSibling.keys = leftSibling.keys[:len(leftSibling.keys)-1]
				if !child.leaf {
					child.children = append([]*BTreeNode{leftSibling.children[len(leftSibling.children)-1]}, child.children...)
					leftSibling.children = leftSibling.children[:len(leftSibling.children)-1]
				}
			} else if i < len(node.children)-1 && len(node.children[i+1].keys) >= t {
				rightSibling := node.children[i+1]
				child.keys = append(child.keys, node.keys[i])
				node.keys[i] = rightSibling.keys[0]
				rightSibling.keys = rightSibling.keys[1:]
				if !child.leaf {
					child.children = append(child.children, rightSibling.children[0])
					rightSibling.children = rightSibling.children[1:]
				}
			} else {
				if i > 0 {
					leftSibling := node.children[i-1]
					leftSibling.keys = append(leftSibling.keys, node.keys[i-1])
					leftSibling.keys = append(leftSibling.keys, child.keys...)
					if !child.leaf {
						leftSibling.children = append(leftSibling.children, child.children...)
					}
					copy(node.keys[i-1:], node.keys[i:])
					copy(node.children[i:], node.children[i+1:])
				} else {
					rightSibling := node.children[i+1]
					child.keys = append(child.keys, node.keys[i])
					child.keys = append(child.keys, rightSibling.keys...)
					if !child.leaf {
						child.children = append(child.children, rightSibling.children...)
					}
					copy(node.keys[i:], node.keys[i+1:])
					copy(node.children[i+1:], node.children[i+2:])
				}
				node.keys = node.keys[:len(node.keys)-1]
				node.children = node.children[:len(node.children)-1]
			}
			child.Remove(key)
		} else {
			node.children[i].Remove(key)
		}
	}
}

func (tree *BTree) Remove(key string) {
	tree.root.Remove(key)
}

func (node *BTreeNode) Search(key string) (*DataType, int) {
	i := 0
	for i < len(node.keys) && key > node.keys[i].nome {
		i++
	}

	if i < len(node.keys) && key == node.keys[i].nome {
		return &node.keys[i], i
	} else if node.leaf {
		return nil, -1
	} else {
		return node.children[i].Search(key)
	}
}

func (tree *BTree) Search(key string) (*DataType, int) {
	return tree.root.Search(key)
}

func (node *BTreeNode) Printa(indent string, last bool) {
	fmt.Print(indent)
	if last {
		fmt.Print("└─ ")
		indent += "    "
	} else {
		fmt.Print("├─ ")
		indent += "|   "
	}
	keys := make([]string, len(node.keys))
	fmt.Print("[")
	for i, key := range node.keys {
		keys[i] = fmt.Sprintf("%v", key)
	}
	fmt.Println(strings.Join(keys, "|"), "]")
	childCount := len(node.children)
	for i, child := range node.children {
		child.Printa(indent, i == childCount-1)
	}
}

//==================================================================================

func Clear() {
	fmt.Print("\033[H\033[2J")
}

//===================================================================================
//							FUNÇÕES SOBRE LISTA
//===================================================================================

func (l *Lista) AchaIndex(nome string) int {
	var i int
	for _, Contato := range l.Contatos {
		i++
		if Contato.Nome == nome {
			return i
		}
	}
	return i
}

func (l *Lista) LerArquivo() {
	content, err := os.ReadFile("agenda.txt")
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}

	lines := strings.Split(string(content), "\n")

	var contato Contato

	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		campo := strings.TrimSpace(parts[0])
		valor := strings.TrimSpace(parts[1])

		switch campo {
		case "Nome":
			contato.Nome = valor
		case "Telefone":
			contato.Telefone = valor
		case "Endereço":
			contato.Endereco = valor
		case "Apagado":
			contato.Apagado = strings.ToLower(valor) == "true"
			l.Contatos = append(l.Contatos, contato)
			contato = Contato{}
		}
	}
	l.OrdenarContatos()
}

func (l *Lista) IncluirContato(nome string, endereco string, telefone string) {

	Contato := Contato{
		Nome:     nome,
		Endereco: endereco,
		Telefone: telefone,
		Apagado:  false,
	}

	l.Contatos = append(l.Contatos, Contato)
	l.OrdenarContatos()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (l *Lista) AlterarContato(nome string, novoNome string, novoEndereco string, novoTelefone string) {

	for i, Contato := range l.Contatos {
		if Contato.Nome == nome && !Contato.Apagado {
			l.Contatos[i].Nome = novoNome
			l.Contatos[i].Endereco = novoEndereco
			l.Contatos[i].Telefone = novoTelefone
			l.OrdenarContatos()

			fmt.Println("Contato atualizado com sucesso!")
			return
		}
	}

	fmt.Println("Contato não encontrado ou está apagado.")
}

func (l *Lista) ApagarLogicamente(nome string) {
	for i, contato := range l.Contatos {
		if contato.Nome == nome && !contato.Apagado {
			l.Contatos[i].Apagado = true
			fmt.Println("Seu contato foi apagado!")
			return
		}
	}
	fmt.Println("Contato não encontrado ou já está apagado.")
}

func (l *Lista) RecuperarContato(nome string) {

	for i, contato := range l.Contatos {
		if contato.Nome == nome && contato.Apagado {
			l.Contatos[i].Apagado = false
			fmt.Println("Contato recuperado com sucesso!")
			return
		}
	}

	fmt.Println("Contato não encontrado ou não está apagado.")
}

func (l *Lista) EsvaziarLixeira() []string {
	var ContatosLimpos []Contato
	ret := []string{}

	for _, Contato := range l.Contatos {
		if !Contato.Apagado {
			ContatosLimpos = append(ContatosLimpos, Contato)
		}
		if Contato.Apagado{
			ret = append(ret, Contato.Nome)
		}
	}

	l.Contatos = ContatosLimpos
	fmt.Println("Sua lixeira foi esvaziada com sucesso!")
	return ret
}

func (l *Lista) Print() {
	i := 1
	for _, Contato := range l.Contatos {
		if Contato.Apagado == false {
			fmt.Println("Contato", i, ":")
			fmt.Println("Nome: ", Contato.Nome)
			fmt.Println("Telefone: ", Contato.Telefone)
			fmt.Println("Endereço: ", Contato.Endereco)
			fmt.Println("=================================")
			i++
		}
	}
}

func (l *Lista) OrdenarContatos() {
	n := len(l.Contatos)

	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if strings.Compare(l.Contatos[j].Nome, l.Contatos[j+1].Nome) > 0 {
				l.Contatos[j], l.Contatos[j+1] = l.Contatos[j+1], l.Contatos[j]
			}
		}
	}

	fmt.Println("Contatos ordenados com sucesso!")
}

func (l *Lista) gravarArquivo() {
	file, err := os.OpenFile("agenda.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	for _, contato := range l.Contatos {
		file.WriteString("Nome: " + contato.Nome + "\n")
		file.WriteString("Telefone: " + contato.Telefone + "\n")
		file.WriteString("Endereço: " + contato.Endereco + "\n")
		file.WriteString("Apagado: " + fmt.Sprintf("%v", contato.Apagado) + "\n\n")
	}
	fmt.Println("Contatos gravados no arquivo com sucesso!")
}


func main() {
	var op int
	lista := Lista{}
	var tree * BTree
	tree = Init()
	lista.LerArquivo()
	for  _, ctt := range lista.Contatos{
		var a DataType 
		a.nome = ctt.Nome
		a.index = lista.AchaIndex(ctt.Nome)
		tree.Insert(a)
	}
	tree.root.Printa("", true)

	for {
		Clear()
		fmt.Println("Agenda (B-Trees)")
		fmt.Println("\t Menu Principal")
		fmt.Println("[ 0] Sair")
		fmt.Println("[ 1] Adicionar Contato")
		fmt.Println("[ 2] Editar Contato")
		fmt.Println("[ 3] Apagar Lógicamente Contato")
		fmt.Println("[ 4] Recuperar Contato")
		fmt.Println("[ 5] Esvaziar lixeira")
		fmt.Println("[ 6] Pesquisar contatos")
		fmt.Println("[ 7] Listar todos os Contatos")

		fmt.Print("\nQual a sua opção? >> ")
		fmt.Scan(&op)
		switch op {
		case 0:
			{
				Clear()
				fmt.Println("Programa Encerrado!\nTecle [ENTER]")
				lista.gravarArquivo()
				lista.OrdenarContatos()
				fmt.Scanln(&op)
				return
			}
		case 1:
			{
				Clear()
				var nome, endereco, telefone string

				fmt.Println("Digite o nome (máximo 30 caracteres):")
				fmt.Scanln(&nome)
				nome = strings.TrimSpace(nome)[:min(len(strings.TrimSpace(nome)), 30)]

				fmt.Println("Digite o endereço (máximo 50 caracteres):")
				fmt.Scanln(&endereco)
				endereco = strings.TrimSpace(endereco)[:min(len(strings.TrimSpace(endereco)), 50)]

				fmt.Println("Digite o telefone (máximo 15 caracteres):")
				fmt.Scanln(&telefone)
				fmt.Print("\n\n\n\n")
				telefone = strings.TrimSpace(telefone)[:min(len(strings.TrimSpace(telefone)), 15)]
				lista.IncluirContato(nome, endereco, telefone)
				lista.OrdenarContatos()
				var a DataType 
				a.nome = nome
				a.index = lista.AchaIndex(nome)
				tree.Insert(a)

				fmt.Println("\nTecle [ENTER]")
				fmt.Scanln()
			}
		case 2:
			{

				Clear()
				var nome, novoNome, novoEndereco, novoTelefone string
				fmt.Println("Digite o nome do contato que deseja atualizar:")
				fmt.Scanln(&nome)

				fmt.Println("Digite o novo nome (máximo 30 caracteres):")
				fmt.Scanln(&novoNome)
				novoNome = strings.TrimSpace(novoNome)[:min(len(strings.TrimSpace(novoNome)), 30)]

				fmt.Println("Digite o novo endereço (máximo 50 caracteres):")
				fmt.Scanln(&novoEndereco)
				novoEndereco = strings.TrimSpace(novoEndereco)[:min(len(strings.TrimSpace(novoEndereco)), 50)]

				fmt.Println("Digite o novo telefone (máximo 15 caracteres):")
				fmt.Scanln(&novoTelefone)
				novoTelefone = strings.TrimSpace(novoTelefone)[:min(len(strings.TrimSpace(novoTelefone)), 15)]
				lista.AlterarContato(nome, novoNome, novoEndereco, novoTelefone)
				lista.OrdenarContatos()

				fmt.Println("\nTecle [ENTER]")
				fmt.Scanln()

			}
		case 3:
			{
				Clear()
				fmt.Println("Qual contato você deseja apagar?")
				var nome string
				fmt.Scanln(&nome)
				lista.ApagarLogicamente(nome)
				lista.OrdenarContatos()
				fmt.Println("Tecle [ENTER]")
				fmt.Scanln()
			}
		case 4:
			{
				Clear()
				fmt.Println("Qual contato deseja recuperar?")
				var nome string
				fmt.Scan(&nome)
				lista.RecuperarContato(nome)
				lista.OrdenarContatos()
				fmt.Println("Tecle [ENTER]")
				fmt.Scanln()
			}
		case 5:
			{
				apagados := []string{}
				Clear()
				apagados = lista.EsvaziarLixeira()
				lista.OrdenarContatos()
				for _, b := range apagados{
				tree.Remove(b)
				}
				fmt.Println("Tecle [ENTER]")
				fmt.Scanln()
			}
		case 6:
			{
				Clear()
				var nome string
				var b int
				fmt.Scan(&nome)

				_, b = tree.Search(nome)
				if lista.Contatos[b].Apagado == true {
					fmt.Print("Contato não existe ou esta apagado!\n")
				}else{
					fmt.Print(lista.Contatos[b])
				}
				lista.OrdenarContatos()
				fmt.Println("Tecle [ENTER]")
				fmt.Scanln()
			}
		case 7:
			{
				Clear()
				lista.OrdenarContatos()
				tree.root.Printa("", true)
				lista.Print()
				fmt.Println("Tecle [ENTER]")
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
