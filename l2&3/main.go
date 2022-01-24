package main

import (
	"flag"
	"fmt"
	"math"
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

func printGraph(e [][]int) {
	for i := 0; i < len(e); i++ {
		if e[i][1] - e[i][0] == 1 {
			leftmargin := strings.Repeat("     ", e[i][0])
			fmt.Println(leftmargin, e[i][0], "- ", e[i][1])
		} else {
			leftmargin := strings.Repeat("     ", e[i][0])
			line := strings.Repeat("---", e[i][1]-e[i][0]+1)
			fmt.Println(leftmargin, e[i][0], line, e[i][1])
		}
	}
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
					if (e[l][0] == j && e[l][1] == k) || (e[l][0] == k && e[l][1] == j) {
						g = -1 // (j, k) or (k, j) exists in edges array
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

type para struct {
	j int
	rijcost int
}

type R struct {
	nexthop int
	cost int
	changed bool
}

type raport struct {
	i int
	ri map[int]R
}

func generateChannels(n int) map[int]chan para {
	mm := make(map[int]chan para)
	for i := 0; i < n; i++ {
		tmp := make(chan para, 1)
		//mm[i] = para {0,0}
		mm[i] = tmp
		fmt.Println(i, tmp)
	}
	return mm
}

func getNexts(v int, e [][]int) []int {
	nexts := make([]int, 0)
	for i := 0; i < len(e); i++ {
		if v == e[i][0] {
			nexts = append(nexts, e[i][1])
		} else if v == e[i][1] {
			nexts = append(nexts, e[i][0])
		}
	}
	return nexts
}

func getNextsChannels(next []int, m map[int] chan para) []chan para{
	result := make([]chan para, 0)
	for i := 0; i < len(next); i++ {
		tmp := make(chan para, 1)
		tmp = m[next[i]]
		result = append(result, tmp)
	}
	return result
}

func node(id int, nexts []int, basicedges [][]int, nextschannels []chan para, in chan para,  sp chan string, spri chan raport) {
	// creating routing table
	ri := make(map[int]R)
	for j := 0; j < len(basicedges) + 1; j++ {
		if id != j { // droga z wierzchołek do samego siebie
			for next := 0; next < len(nexts); next++ {
				if j == nexts[next] {
					ri[j] = R{nexthop: j, cost: 1, changed: true}
				} else {
					if id < j {
						ri[j] = R{nexthop: id + 1,
							cost:    int(math.Abs(float64(id - j))),
							changed: true,
						}
					} else {
						ri[j] = R{nexthop: id - 1,
							cost:    int(math.Abs(float64(id - j))),
							changed: true,
						}
					}
				}
			}
		}
	}
	// the started routing table is done
	fmt.Println(id, ri)

	// sender i
	go func() {
		for {
			for j := 0; j < len(ri) + 1; j++ {
				if ri[j].changed == true {
					for l := 0; l < len(nextschannels); l++ {
						mutex.Lock()
						p := para{
							j, ri[j].cost,
						}
						tmp := R {ri[j].nexthop, ri[j].cost, false}
						ri[j] = tmp
						mutex.Unlock()
						for s := 0; s < len(nextschannels); s++ { // wysyłanie oferty do sąsiadów
							spri <- raport{id, ri}
							sp <- fmt.Sprint(id, " wysyłam oferte ", p, " do ", nexts[s])
							nextschannels[s] <- p
						}
					}
				}
			}
			rand.Seed(time.Now().UnixNano())
			sec := rand.Float64() * 5
			time.Sleep(time.Second * time.Duration(sec))
		}
	}()


	// receiver i
	go func() {
		for {
			select {
			case paraIN := <-in :
				sp <- fmt.Sprint(id, " odebrałem ofertę ", paraIN)
				newcost := paraIN.rijcost + 1
				mutex.Lock()
				if newcost < ri[paraIN.j].cost {
					tmp := R {nexthop: paraIN.j, cost: newcost, changed: true}
					ri[paraIN.j] = tmp
					sp <- fmt.Sprint(id, " robie zmiane dla R[", id, "][", paraIN.j, "] = ", tmp)
				}
				mutex.Unlock()
			case <-time.After(500 * time.Millisecond):
			}
		}
	}()
}

var mutex = sync.Mutex{}

func main() {
	nPtr := flag.Int("n", 0, "an int") // G(n-1) 0..n-1
	dPtr := flag.Int("d", 0, "an int") // d <= n + 1
	rawPtr := flag.Int("raw", 0, "an int")
	flag.Parse()

	e := generateEdges(*nPtr)
	basicE := e
	v := generateVertices(*nPtr)
	e = generateShortcuts(*nPtr, *dPtr, e)
	m := generateChannels(*nPtr)
	fmt.Println(m)
	fmt.Println(e)
	fmt.Println(v)
	fmt.Println(basicE)
	printGraph(e)

	var serverPrinter = make(chan string)
	var serverPrinterRi = make(chan raport)
	var lastRaport = make(map[int]map[int]R)

	for i:= 0; i < len(v); i++ {
		node(i, getNexts(i, e), basicE, getNextsChannels(getNexts(i, e), m), m[i], serverPrinter, serverPrinterRi)
	}

	fmt.Println(" ")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			select {
			case sp := <- serverPrinter:
				fmt.Println(sp)
			case spri := <- serverPrinterRi:
				lastRaport[spri.i] = spri.ri
			case <-time.After(4 * time.Second):
				//fmt.Println("\t\t\t\tdawno juz nic nie drukowałem -> koniec programu")
				wg.Done()
			}
		}
	}()
	wg.Wait()

	fmt.Println("Koniec zmian w routing tables -> raporty")
	// drukowanie raportów
	for i := 0; i < len(lastRaport); i++ {
		if *rawPtr == 1 {
			fmt.Println(i, lastRaport[i])
		} else { // pretty printing
			fmt.Println("Wierzchołek", i, "do wierzchołka: ")
			for j := 0; j < len(lastRaport[i])+1; j++ {
				if i != j {
					fmt.Println("\t", j, "ma", lastRaport[i][j].cost, "skoki")
				}
			}
		}
	}
}