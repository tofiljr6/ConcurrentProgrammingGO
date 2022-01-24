package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

func generateEdges(n int) [][]int {
	e := make([][]int, n-1)
	// generate edges in graph
	for i := 0; i < n-2; i++ {
		// create edge
		tmp := make([]int, 2)
		tmp[0] = i
		tmp[1] = i + 1

		// add edge to final array
		e[i] = tmp
	}

	// the last one - one receiver
	tmp := make([]int, 2)
	tmp[0] = n - 2
	tmp[1] = n - 1
	e[n-2] = tmp
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

func generateShortcuts(n, d int, e [][]int) [][]int {
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

				if k == n-1 || j == 0 { // the last is receiver nth node
					break
				}

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

func generateShortcutsBack(n, b int, e [][]int) [][]int {
	for i := 0; i < b; i += 0 {
		// create edge
		tmp := make([]int, 2)

		for {
			rand.Seed(time.Now().UnixNano())
			j := rand.Intn(n)
			k := rand.Intn(n)

			if j-k >= 1 {
				tmp[0] = j
				tmp[1] = k
				//fmt.Println(tmp)

				if j == n-1 || k == 0 {
					// the last is receiver nth node and producer is only one
					// after producer make all packages => dead
					break
				}

				g := 0
				for l := 0; l < len(e); l++ {
					if e[l][0] == j && e[l][1] == k {
						g = -1
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
		} else if v == e[i][0] && e[i][0] > e[i][1] {
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

func generateArrayOfPackageLife(k, h int) map[int]int {
	mm := make(map[int]int)
	for i := 1; i < k+1; i++ {
		mm[i] = h
	}
	return mm
}

func producer(k int, nc []chan int, vp map[int][]int, pv map[int][]int, sp chan string) {
	// nc - nexts channels from current vertices
	pack := 1
	for q := 1; q < k+1; q++ {
		rand.Seed(time.Now().UnixNano())
		sec := rand.Float64() * 2
		time.Sleep(time.Second * time.Duration(sec))

		randomChannel := rand.Intn(len(nc))
		sp <- fmt.Sprint("Pakiet ", q*pack, " jest w wierchołku 0")
		vp[0] = append(vp[0], pack*q)
		pv[pack*q] = append(pv[pack*q], 0) // producers id is 0
		nc[randomChannel] <- pack * q
	}
}

func node(id int, in <-chan int, nc []chan int, pv map[int][]int, vp map[int][]int, sp chan string, h map[int]int) {
	for {
		p := <-in
		pv[id] = append(pv[id], p)
		vp[p] = append(vp[p], id)

		leftmargin := strings.Repeat("-", id)
		sp <- fmt.Sprint(leftmargin, "Pakiet ", p, " jest w wierzchołku ", id)

		rand.Seed(time.Now().UnixNano())
		sec := rand.Float64() * 2
		time.Sleep(time.Second * time.Duration(sec))

		var l = sync.Mutex{}
		l.Lock()
		h[p]--
		plife := h[p]
		l.Unlock()

		if plife > 0 {
			randomChannel := rand.Intn(len(nc))
			nc[randomChannel] <- p
		} else {
			sp <- fmt.Sprint("\t\t\t\tŚmierć pakietu ", p)
			l.Lock()
			lifepack--
			l.Unlock()
		}
	}
}

func consumer(id int, in <-chan int, d chan<- bool, pv map[int][]int, vp map[int][]int, sp chan string) {
	for {
		select {
		case p := <-in:
			pv[id] = append(pv[id], p)
			vp[p] = append(vp[p], id)
			leftmargin := strings.Repeat("~", id)
			sp <- fmt.Sprint(leftmargin, "\t\t\t\tPakiet ", p, " został odebrany")
			lifepack--
		case <-time.After(1 * time.Second):
			// timeout - next loop
			//fmt.Println("timeout", lifepack)
		}

		if lifepack == 0 {
			break
		}
	}
	d <- true
}

func printGraph(e [][]int) {
	for i := 0; i < len(e); i++ {
		if e[i][0] < e[i][1] {
			leftmargin := strings.Repeat("     ", e[i][0])
			fmt.Println(leftmargin, e[i][0], "->", e[i][1])
		} else {
			leftmargin := strings.Repeat("     ", e[i][0]-1)
			fmt.Println(leftmargin, e[i][1], "<-", e[i][0])
		}
	}
}

var lifepack int = 1

func main() {
	// parse params from command line
	nPtr := flag.Int("n", 0, "an int") // G(n-1) 0..n-1
	dPtr := flag.Int("d", 0, "an int") // d <= n + 1
	kPtr := flag.Int("k", 0, "an int") // k - number of packages
	bPtr := flag.Int("b", 0, "an int") // back shortcuts
	hPtr := flag.Int("h", 0, "an int") // life of package
	flag.Parse()

	*nPtr = *nPtr + 2 // fast solution with sender and receiver - i will write better solution in the future

	if (*nPtr > 0 && *dPtr > 0 && *kPtr > 0 && *bPtr > 0 && *hPtr > 0) && (*dPtr <= *nPtr-4) && (*bPtr <= *nPtr-2) {
		/*
			conditions
			1* all of params must be positive
			2* jump to forward must be d<n+1,
				but i increment before if nPtr so its - 2 and first and the last one node
				is only receiver and sender so -4 finally
			3* jump to backward (shortcuts) must be b < n+1
		*/

		// create graph
		e := generateEdges(*nPtr)
		v := generateVertices(*nPtr)
		e = generateShortcuts(*nPtr, *dPtr, e)
		e = generateShortcutsBack(*nPtr, *bPtr, e)
		m := generateChannels(*nPtr)
		vp := generateArrayHistoryVertices(*nPtr) // get history of packages in i-vertices
		pv := generateArrayHistoryPackages(*kPtr)
		h := generateArrayOfPackageLife(*kPtr, *hPtr)
		fmt.Println(*hPtr)
		fmt.Println(h)
		lifepack = *kPtr

		// printing graph
		fmt.Println("GRAPH:", *bPtr)
		printGraph(e)

		fmt.Println("E:", e)
		fmt.Println("V:", v)
		//server.Println(m)

		fmt.Println("NEXT:")
		for i := 0; i < len(v); i++ {
			fmt.Println("i=", i, getNexts(v[i], e))
		}

		var done = make(chan bool)
		var serverPrinter = make(chan string)

		go producer(*kPtr, getNextChannels(getNexts(v[0], e), m), vp, pv, serverPrinter)

		for i := 1; i < *nPtr-1; i++ {
			go node(i, m[i], getNextChannels(getNexts(v[i], e), m), vp, pv, serverPrinter, h)
		}

		go consumer(*nPtr-1, m[*nPtr-1], done, vp, pv, serverPrinter)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			for {
				select {
				case sp := <-serverPrinter:
					fmt.Println(sp)
				case <-done:
					wg.Done()
				}
			}
		}()
		wg.Wait()

		// Reports
		fmt.Println("\nWierzchołek -> pakiet")
		for i := 0; i < *nPtr; i++ {
			fmt.Println(i, "->", vp[i])
		}

		fmt.Println("Pakiet -> wierzchołek")
		for i := 0; i < *kPtr; i++ {
			fmt.Println((i + 1), "->", pv[(i+1)])
		}
	} else {
		fmt.Println("Nie poprawne parametry")
	}
}
