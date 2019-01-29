package main

import (
	"fmt"
	"net/http"
)

type Server struct{}

// method associating with Server struct
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	if r.Method == "GET" {
		board := Board{}
		board = board.Start()
		fmt.Println(board)
		fmt.Fprint(w, "<h1>Test</h1>") // Fprint write to a buffer output
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

type Tiki struct{
	Color string
	Type string
	Name string
}

type Board struct{
	Position1 Tiki
	Position2 Tiki
	Position3 Tiki
	Position4 Tiki
	Position5 Tiki
	Position6 Tiki
	Position7 Tiki
	Position8 Tiki
	Position9 Tiki
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

func (b Board) Start() Board {
	t1 := Tiki{"Red","1","r1"}
	t2 := Tiki{"Orange","1","o1"}
	t3 := Tiki{"Yellow","1","y1"}
	t4 := Tiki{"Green","2","g2"}
	t5 := Tiki{"Blue","2","b2"}
	t6 := Tiki{"Purple","2","p2"}
	t7 := Tiki{"Grey","3","g3"}
	t8 := Tiki{"Pink","3","p3"}
	t9 := Tiki{"White","3","w3"}
	startBoard := Board{t1,t2,t3,t4,t5,t6,t7,t8,t9}
	return startBoard
}