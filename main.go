package main

import (
    "bufio"
    "fmt"
    "log"
    "math"
    "os"
    "strconv"
    "strings"
)

const (
    digits = "0123456789"
)

type Map struct {
    Src uint64
    Dst uint64
    Rng uint64
}

type Category struct {
    Name string
    Maps []Map
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
                Name: getCategoryName(lines[i]),
                Maps: make([]Map, 0),
            }
            j := i + 1
            for ; j < len(lines) && strings.ContainsAny(lines[j], digits); j++ {
                dsr := parseNumbers(lines[j])
                m := Map{
                    Src: uint64(dsr[1]),
                    Dst: uint64(dsr[0]),
                    Rng: uint64(dsr[2]),
                }
                cat.Maps = append(cat.Maps, m)
            }
            cs = append(cs, cat)
            i = j
            continue
        }
        i++
    }

    min := uint64(math.MaxUint64)
    for _, seed := range seeds {
        loc := seed
        for c := 0; c < len(cs); c++ {
            for m := 0; m < len(cs[c].Maps); m++ {
                if cs[c].Maps[m].Src <= loc && loc <= cs[c].Maps[m].Src+cs[c].Maps[m].Rng {
                    loc = loc - cs[c].Maps[m].Src + cs[c].Maps[m].Dst
                    break
                }
            }
        }
        if loc < min {
            min = loc
        }
    }
    fmt.Println(min)
}

func parseNumbers(line string) []uint64 {
    var nums []string
    numbers := make([]uint64, 0)
    switch {
    case strings.Contains(line, ":"):
        nums = strings.Split(strings.TrimSpace(strings.Split(line, ":")[1]), " ")
    default:
        nums = strings.Split(line, " ")
    }
    for _, num := range nums {
        seed, _ := strconv.Atoi(strings.TrimSpace(num))
        numbers = append(numbers, uint64(seed))
    }

    return numbers
}

func getCategoryName(line string) string {
    str := strings.Split(line, " ")[0]
    str = strings.ReplaceAll(str, "-", " ")
    name := strings.ToTitle(str)

    return name
}
