package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var server = log.New(os.Stdout, "", 0)

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

func getNextChannels(next []int, m map[int]chan int) []chan int {
	result := make([]chan int, 0)
	for i := 0; i < len(next); i++ {
		tmp := make(chan int, 1)
		tmp = m[next[i]]
		result = append(result, tmp)
	}
	return result
}

func generateChannels(n int) map[int]chan int {
	mm := make(map[int]chan int)
	for i := 0; i < n; i++ {
		tmp := make(chan int, 1)
		mm[i] = tmp
	}
	return mm
}

func generateArrayHistoryVertices(n int) map[int][]int {
	mm := make(map[int][]int)
	for i := 0; i < n; i++ {
		tmp := make([]int, 0)
		mm[i] = tmp
	}
	return mm
}

func generateArrayHistoryPackages(k int) map[int][]int {
	mm := make(map[int][]int)
	pack := 1000
	for i := 1; i < k+1; i++ {
		tmp := make([]int, 0)
		mm[i*pack] = tmp
	}
	return mm
}

func producer(nc []chan int, vp map[int][]int, pv map[int][]int) {
	// nc - nexts channels from current vertices
	pack := 1000
	for q := 1; q < k+1; q++ {
		rand.Seed(time.Now().UnixNano())
		sec := rand.Intn(2)
		time.Sleep(time.Second * time.Duration(sec))

		randomChannel := rand.Intn(len(nc))
		server.Println("Pakiet", q*pack, "jest w wierchołku 0")
		vp[0] = append(vp[0], pack*q)
		pv[pack*q] = append(pv[pack*q], 0) // producers id is 0
		nc[randomChannel] <- pack * q
	}
}

func node(id int, in <-chan int, nc []chan int, pv map[int][]int, vp map[int][]int) {
	for {
		p := <-in
		pv[id] = append(pv[id], p)
		vp[p] = append(vp[p], id)

		leftmargin := strings.Repeat("-", id)
		server.Println(leftmargin, "Pakiet", p, "jest w wierzchołku", id)

		rand.Seed(time.Now().UnixNano())
		sec := rand.Intn(2)
		time.Sleep(time.Second * time.Duration(sec))

		randomChannel := rand.Intn(len(nc))
		//server.Println("Wysyłam pakiet", p, "do wierzchołka", nc[randomChannel])
		nc[randomChannel] <- p
	}
}

func consumer(id int, in <-chan int, d chan<- bool, pv map[int][]int, vp map[int][]int) {
	for l := 0; l < k; l++ {
		p := <-in
		pv[id] = append(pv[id], p)
		vp[p] = append(vp[p], id)

		leftmargin := strings.Repeat("-", id)
		server.Println(leftmargin, "Pakiet", p, "został odebrany")
	}
	d <- true
}

const n = 4 // G(n-1) 0..n-1
const d = 2 // d <= n + 1
const k = 4 // k - number of packages

func main() {
	// create graph
	e := generateEdges(n)
	v := generateVertices(n)
	e = generateDigests(n, d, e)
	m := generateChannels(n)
	vp := generateArrayHistoryVertices(n) // get history of packages in i-vertices
	pv := generateArrayHistoryPackages(k)

	server.Println("E:", e)
	server.Println("V:", v)
	//server.Println(m)

	var done = make(chan bool)

	go producer(getNextChannels(getNexts(v[0], e), m), vp, pv)

	for i := 1; i < n-1; i++ {
		go node(i, m[i], getNextChannels(getNexts(v[i], e), m), vp, pv)
	}

	go consumer(n-1, m[n-1], done, vp, pv)

	<-done

	// Reports
	fmt.Println("\nWierzchołek -> pakiet")
	for i := 0; i < n; i++ {
		fmt.Println(i, "->", vp[i])
	}

	fmt.Println("Pakiet -> wierzchołek")
	for i := 0; i < n; i++ {
		fmt.Println((i+1)*1000, "->", pv[(i+1)*1000])
	}
}
