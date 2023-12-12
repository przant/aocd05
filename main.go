package main

import (
    "bufio"
    "fmt"
    "log"
    "math"
    "os"
    "strconv"
    "strings"
    "sync"
)

const (
    digits = "0123456789"
)

type Map struct {
    Src int
    Dst int
    Rng int
}

type Category struct {
    Name  string
    Maps  []Map
    Pairs map[int]int
}

func main() {
    pf, err := os.Open("example.txt")
    if err != nil {
        log.Fatalf("while opening file %q: %s", pf.Name(), err)
    }
    defer pf.Close()

    scnr := bufio.NewScanner(pf)

    lines := make([]string, 0)
    for scnr.Scan() {
        lines = append(lines, scnr.Text())
    }

    cs := make([]Category, 0)
    seeds := parseNumbers(lines[0])

    for i := 1; i < len(lines); {
        if strings.HasSuffix(lines[i], "map:") {
            cat := Category{
                Name:  getCategoryName(lines[i]),
                Maps:  make([]Map, 0),
                Pairs: make(map[int]int),
            }
            j := i + 1
            for ; j < len(lines) && strings.ContainsAny(lines[j], digits); j++ {
                dsr := parseNumbers(lines[j])
                m := Map{
                    Src: dsr[1],
                    Dst: dsr[0],
                    Rng: dsr[2],
                }
                cat.Maps = append(cat.Maps, m)
            }
            cs = append(cs, cat)
            i = j
            continue
        }
        i++
    }

    wg := &sync.WaitGroup{}

    wg.Add(len(cs))
    for c := 0; c < len(cs); c++ {
        go pairs(c, cs, wg)
    }
    wg.Wait()

    ml := math.MaxInt
    for _, s := range seeds {
        d := s
        for _, c := range cs {
            if n, ok := c.Pairs[d]; ok {
                d = n
            }
        }
        if d < ml {
            ml = d
        }
    }

    fmt.Println(ml)
}

func parseNumbers(line string) []int {
    var nums []string
    numbers := make([]int, 0)
    switch {
    case strings.Contains(line, ":"):
        nums = strings.Split(strings.TrimSpace(strings.Split(line, ":")[1]), " ")
    default:
        nums = strings.Split(line, " ")
    }
    for _, num := range nums {
        seed, _ := strconv.Atoi(strings.TrimSpace(num))
        numbers = append(numbers, seed)
    }

    return numbers
}

func getCategoryName(line string) string {
    str := strings.Split(line, " ")[0]
    str = strings.ReplaceAll(str, "-", " ")
    name := strings.ToTitle(str)

    return name
}

func pairs(c int, cs []Category, wg *sync.WaitGroup) {
    defer wg.Done()

    for m := 0; m < len(cs[c].Maps); m++ {
        src := cs[c].Maps[m].Src
        dst := cs[c].Maps[m].Dst
        rng := cs[c].Maps[m].Rng
        for n := 0; n < rng; n++ {
            cs[c].Pairs[src+n] = dst + n
        }
    }
}
