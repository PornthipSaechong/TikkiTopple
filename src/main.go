package main

import (
	"fmt"
	"net/http"
	"math/rand"
	"time"
)

type Server struct{}

// method associating with Server struct
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	if r.Method == "GET" {
		game := Game{}
		game = game.Start(2)
		fmt.Println(game)
		fmt.Fprint(w, "<h1>Test</h1>") // Fprint write to a buffer output
		player1 := game.Players[0]
		// player2 := game.Players[1]
		board := game.Board
		fmt.Println("Before")
		fmt.Println(board)
		fmt.Println(player1.Cards)
		cardIndex := 7 - 1 // input as array instead of card
		tikkiIndex := 1 - 1
		player1, board = player1.Move(cardIndex, tikkiIndex, board)
		fmt.Println("After")
		fmt.Println(board)
		fmt.Println(player1.Cards)
		if game.CheckEnd() {
			fmt.Println("Game has ended")
		} else {
			fmt.Println("Game is still ongoing")
		}

	}
}

func main (){
	var s Server
	err := http.ListenAndServe("localhost:4000", s)
	checkError(err)
}

func checkError (err error) {
	if err != nil {
		panic(err) // panic will stop excution
	}
}

type Game struct {
	Board Board
	Players []Player
	Tikis []Tiki
	Discard []Card
}

type Tiki struct{
	Color string
	Type string
	Name string
}

type Board struct{
	Tiki []Tiki
}

type SecretCard struct{
	First Tiki
	Second Tiki
	Third Tiki
}

type Card struct {
	Name string
	Action string
	Count int
}

type Player struct {
	Cards []Card
	Secret SecretCard
}

func (g Game) Start(playerNumber int) Game {
	t1 := Tiki{"Red","1","r1"}
	t2 := Tiki{"Orange","1","o1"}
	t3 := Tiki{"Yellow","1","y1"}
	t4 := Tiki{"Green","2","g2"}
	t5 := Tiki{"Blue","2","b2"}
	t6 := Tiki{"Purple","2","p2"}
	t7 := Tiki{"Grey","3","g3"}
	t8 := Tiki{"Pink","3","p3"}
	t9 := Tiki{"White","3","w3"}
	g1 := []Tiki{t1,t2,t3}
	g2 := []Tiki{t4,t5,t6}
	g3 := []Tiki{t7,t8,t9}

	r := rand.New(rand.NewSource(time.Now().Unix()))

	b := Board{}
	g.Board = b.New(g1,g2,g3)

	p := Player{}
	p = p.New()
	
	for i := 0; i<playerNumber; i++ {
		s := SecretCard{}
		s.First = g1[r.Intn(len(g1))]
		s.Second = g2[r.Intn(len(g2))]
		s.Third = g3[r.Intn(len(g3))]
		p.Secret = s 
		g.Players = append(g.Players, p)
	}
	// fmt.Printf("%p,%p", &g.Players[0], &g.Players[1])
	return g
}


func (g Game) CheckEnd() bool {
	if len(g.Board.Tiki) == 3 {
		return true
	}
	for p := 0; p<len(g.Players); p++ {
		if len(g.Players[p].Cards) == 0 {
			return true
		}
	}
	
	return false
}

func (b Board) New(g1 , g2, g3 []Tiki) Board {
	g1 = Shuffle(g1)
	g2 = Shuffle(g2)
	g3 = Shuffle(g3)

	r := rand.New(rand.NewSource(time.Now().Unix()))
	perm := r.Perm(3)

	gA := [][]Tiki {g1,g2,g3}
	b.Tiki = append(b.Tiki,gA[perm[0]]...)
	b.Tiki = append(b.Tiki,gA[perm[1]]...)
	b.Tiki = append(b.Tiki,gA[perm[2]]...)
	return b	
}

func (p Player) New() Player {
	c1 := Card{"Move1","UP",1}
	c2 := Card{"Move2","UP",2}
	c3 := Card{"Move3","UP",3}
	cx := Card{"X","X",1}
	cdrop := Card{"Drop","DROP",0}
	p.Cards = make([]Card, 0, 7)
	p.Cards = append(p.Cards, c1)
	p.Cards = append(p.Cards, c1)
	p.Cards = append(p.Cards, c2)
	p.Cards = append(p.Cards, c3)
	p.Cards = append(p.Cards, cx)
	p.Cards = append(p.Cards, cx)
	p.Cards = append(p.Cards, cdrop)
	return p
}

func Shuffle(t []Tiki) []Tiki {
	// fmt.Println(t)
	// fmt.Println(time.Now().Unix())
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n:=len(t); n>0; n-- {
		rInt:= r.Intn(len(t))
		// fmt.Println("random int ", rInt)
		t[n-1], t[rInt] = t[rInt], t[n-1]
	}
	return t
} 

func (p Player) Move(cardIndex int, posIndex int, b Board) (Player, Board) {
	c := p.Cards[cardIndex] // get card from index input
	fmt.Println(c.Action)
	t := b.Tiki[posIndex] // get tiki to move
	switch c.Action {
		case "UP":
			fmt.Printf("Move tiki: %v by %v\n", t.Name, c.Count)
			if posIndex-c.Count<0 {
				fmt.Println("Invalid move")
			} else {
				b = b.Move(posIndex,c.Count,"up") // move the tiki on the board
			}				
		case "X":
			fmt.Printf("Destroy tiki: %v\n", t.Name)
			if (len(b.Tiki) <= 3) {
				fmt.Println("Invalid move")
			} else {
				b = b.Remove()
			}
		case "DROP":
			fmt.Printf("Drop tiki: %v\n", t.Name)
			if posIndex == len(b.Tiki) {
				fmt.Println("Invalid move")
			} else {
				b = b.Move(posIndex,len(b.Tiki)-1-(posIndex),"down")
			}
		default:
			fmt.Println("Invalid input")
	}
	p.Cards = append(p.Cards[:cardIndex], p.Cards[cardIndex+1:]...) // remove used card
	return p,b
} 

func (b Board) Move(index, count int, direction string) Board {
	var newPos int
	if (direction == "up") {
		newPos = index-count
	} else if (direction == "down") {
		newPos = index+count
	} else {
		fmt.Println("Invalid direction")
		return b
	}
	fmt.Println("new position", newPos)
	// remove the tiki from the chosen position
	chosen := b.Tiki[index]
	fmt.Println("chosen", chosen)
	b.Tiki = append(b.Tiki[:index], b.Tiki[index+1:]...)
	fmt.Println("removed chosen tiki array", b.Tiki)
	// insert tiki into position
	b.Tiki = append(b.Tiki,Tiki{})
	copy(b.Tiki[newPos+1:],b.Tiki[newPos:])
	b.Tiki[newPos] = chosen
	fmt.Println("inserted right array", b.Tiki)
	fmt.Printf("Moving tiki %v by %v\n", chosen.Name, count)
	return b
}

func (b Board) Remove() Board {
	lastIndex := len(b.Tiki)-1
	b.Tiki = append(b.Tiki[:lastIndex], b.Tiki[lastIndex+1:]...)
	return b
}
