package main

import (
	"fmt"
	"strings"
)

const t = 2 

type DataType uint16


type BTreeNode struct {
	leaf     bool           
	keys     []DataType     
	children []*BTreeNode   
	del: bool
}


func InitNode(leaf bool) *BTreeNode {
	return &BTreeNode{
		leaf:     leaf,
		keys:     []DataType{},
		children: []*BTreeNode{},
	}
}


type BTree struct {
	root *BTreeNode
}


func Init() *BTree {
	return &BTree{
		root: InitNode(true),
	}
}

func (node *BTreeNode) Print(indent string, last bool) {
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
		child.Print(indent, i == childCount-1)
	}
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

	node.keys = append(node.keys, 0)
	copy(node.keys[i+1:], node.keys[i:])
	node.keys[i] = child.keys[t-1]
	child.keys = child.keys[:t-1]
}


func (node *BTreeNode) Insert(key DataType) {
	if !node.leaf {

		i := len(node.keys) - 1
		for i >= 0 && key < node.keys[i] {
			i--
		}

		if len(node.children[i+1].keys) == 2*t-1 {
			node.splitChild(int16(i) + 1)
			if key > node.keys[i+1] {
				i++
			}
		}
		node.children[i+1].Insert(key)
	} else {
		i := len(node.keys) - 1
		node.keys = append(node.keys, 0)
		for i >= 0 && key < node.keys[i] {
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


func (node *BTreeNode) Search(key DataType) *DataType {
	i := 0
	for i < len(node.keys) && key > node.keys[i] {
		i++
	}

	if i < len(node.keys) && key == node.keys[i] {
		return &node.keys[i]
	} else if node.leaf {
		return nil
	} else {
		return node.children[i].Search(key)
	}
}

func (tree *BTree) Search(key DataType) *DataType {
	return tree.root.Search(key)
}


func Clear(){
   fmt.Print("\033[H\033[2J") 
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




func main() {
   var (
      op DataType
      tree *BTree
      x DataType
   )
	
   for {
      Clear() 
      fmt.Println("ÁRVORES-B (B-Trees)")
      fmt.Println("\t Menu Principal")
      fmt.Println("[ 0] Sair")
      fmt.Println("[ 1] Criar árvore")
      fmt.Println("[ 2] Imprimir árvore")
      fmt.Println("[ 3] Inserir elemento")
      fmt.Println("[ 4] Remover elemento")
      fmt.Println("[ 5] Procurar elemento")
	  fmt.Println("[ 5] Percurso")
      
      
      fmt.Print("\nQual a sua opção? >> ")
      fmt.Scan(&op)
      switch op {
      case 0: {
         Clear()
         fmt.Println("Programa Encerrado!\nTecle [ENTER]")
         fmt.Scanln(&op)
         return
      }
	   case 1: {
         Clear()
         fmt.Println("Criando árvore...")
         tree = Init()
         fmt.Println("Árvore criada!\n Tecle [ENTER]")
         fmt.Scanln()
      }
      case 2: {
         Clear()
         fmt.Println("Árvore armazenada:")
         tree.root.Print("", true)
         fmt.Println("\nTecle [ENTER]")
         fmt.Scanln()
      }
      case 3: {
         Clear()
         fmt.Println("Inserindo elemento")
         fmt.Print("Número a inserir >> ")
         fmt.Scanln(&x)
         tree.Insert(x)
         fmt.Println("Elemento inserido. Tecle [ENTER]")
         fmt.Scanln()
      }
    case 6: {
         Clear()
         fmt.Println("Percurso \n")
         fmt.Print("Número a remover >> ")
         tree.PercursoEmOrdem()
         fmt.Println("Elemento Removido. Tecle [ENTER]")
         fmt.Scanln()
      }
      /*case 5: {
         Clear()
         fmt.Println("Procurar elemento\n")
         fmt.Print("Número a procurar >> ")
         fmt.Scanln(&x)
         a := binTree.Find(x)
         if a == nil {
            fmt.Println("Elemento não encontrado. Tecle [ENTER]")
            fmt.Scanln()
         } else {
            fmt.Println(x, "está na árvore:", a)
            fmt.Scanln()
         } 
      }
*/	 default: {
	    fmt.Println("Opção Inválida!\nTecle [ENTER]")
	    fmt.Scanln(&op)
	 }
    }
   }
}
