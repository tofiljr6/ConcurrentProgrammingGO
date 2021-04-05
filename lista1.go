package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Node struct {
	id   int       // node identifier
	pack chan bool // pack channel is used for demonstrate if the package
	next *Node
}

func generateEdges(n int) [][]int {
	e := make([][]int, n-1)
	// generate edges in graph
	for i := 0; i < n-1; i++ {
		// create edge
		tmp := make([]int, 2)
		tmp[0] = i
		tmp[1] = i + 1

		// add edge to final array
		e[i] = tmp
	}
	return e
}

func generateVertices(n int) []int {
	v := make([]int, n)
	// generate vertices in graph
	for i := 0; i < n; i++ {
		v[i] = i
	}
	return v
}

func generateDigests(n, d int, e [][]int) [][]int {
	// generate digests in graph
	for i := 0; i < d; i += 0 {
		// create edge
		tmp := make([]int, 2)

		for {
			rand.Seed(time.Now().UnixNano())
			j := rand.Intn(n)
			k := rand.Intn(n)

			if k-j > 1 {
				tmp[0] = j
				tmp[1] = k

				g := 0 // check that draw by lot is not existed edge
				for l := 0; l < len(e); l++ {
					if e[l][0] == j && e[l][1] == k {
						g = -1 // (j, k) exists in edges array
						break
					}
				}

				if g != -1 {
					e = append(e, tmp)
					i++
					break
				}
			}
		}
	}
	return e
}

func getNexts(v int, e [][]int) []int {
	nexts := make([]int, 0)
	for i := 0; i < len(e); i++ {
		if v == e[i][0] {
			nexts = append(nexts, e[i][1])
		}
	}
	return nexts
}

func producer(link chan<- int) {
	pack := 0
	for x := 0; x < k; x++ {
		sec := rand.Intn(5) + 1
		fmt.Println("\tWysyłam paczkę za ", sec, "sec")

		pack++
		fmt.Println("\tWysłałem paczke", pack)
		link <- pack
	}
	close(link)
}

func consumer(link <-chan int, done chan<- bool, next []int) {
	for pack := range link {
		fmt.Println("Odebrałem paczkę", pack)

		sec := rand.Intn(5)
		fmt.Println("Analizuje paczke przez", sec, "sec")
		time.Sleep(time.Second * time.Duration(sec))

		fmt.Println("Przeanalizowałem pczacke", pack)
		//
		index := rand.Intn(len(next))
		x := next[index]
		fmt.Println("przesyłam paczkę do wierzchołka", x)
	}
	done <- true
}

const n = 5
const d = 4 // d <= n + 1
const k = 6

func main() {
	// create graph
	e := generateEdges(n)
	v := generateVertices(n)
	e = generateDigests(n, d, e)

	fmt.Println("E:", e)
	fmt.Println("V:", v)

	for i := 0; i < n; i++ {
		fmt.Println("Dla wierzcholka ", v[i], "następnikami są: ", getNexts(v[i], e))
	}

	link := make(chan int)
	done := make(chan bool)
	go producer(link)
	go consumer(link, done, getNexts(0, e))
	<-done
}
