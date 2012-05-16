package main

import (
    "ascent"
    "ascent/specimens"
    "math/rand"
    "math"
    "os"
    "time"
    "strings"
    "fmt"
)

type mutableString struct {
    value string
}

func (s *mutableString) Clone() specimens.Specimen {
    clone := *s
    return &clone
}

func (s *mutableString) ToString() string {
    return s.value
}

func randomLetterMutation(specimen specimens.Specimen) {
    s := specimen.(*mutableString)
    if (len(s.value) > 0) {
        index := rand.Intn(len(s.value))
        s.value = s.value[0:index] + string(rand.Intn(126-32)+32) + s.value[index+1:]
    }
}

func dropLetterMutation(specimen specimens.Specimen) {
    s := specimen.(*mutableString)
    if (len(s.value) > 0) {
        index := rand.Intn(len(s.value))
        s.value = s.value[0:index] + s.value[index+1:]
    }
}

func addLetterMutation(specimen specimens.Specimen) {
    s := specimen.(*mutableString)
    index := rand.Intn(len(s.value)+1)
    s.value = s.value[0:index] + string(rand.Intn(126-32)+32) + s.value[index:]
}

func main() {

    if len(os.Args) < 3 {
        println("Usage: ascentclient <start string> <goal string>\n")
        os.Exit(1)
    }

    rand.Seed(int64(time.Now().Nanosecond()))

    start, goal := os.Args[1], os.Args[2]

    engine := ascent.New()
    engine.Mutations().Register(randomLetterMutation, 0.33)
    engine.Mutations().Register(dropLetterMutation, 0.33)
    engine.Mutations().Register(addLetterMutation, 0.33)
    engine.Mutations().SetIdentityChance(0.01)

    generationCounter := 0

    engine.SetGenerationCallback(func(winner specimens.Specimen) {
        println(winner.(*mutableString).value)
        generationCounter++
    })

    println(start)

    final := engine.Run(&mutableString{start}, func(specimen specimens.Specimen) (float32, bool) {
        value := specimen.ToString()

        if value == goal {
            return 0, true
        }

        score := float32(-math.Abs(float64(len(value) - len(goal)))/4.0)
        for i := 0; i < len(value) && i < len(goal); i++ {
            if value[i] == goal[i] {
                score += 1.0
            }
        }
        for i := 0; i < len(goal); i++ {
            if strings.ContainsRune(value, rune(goal[i])) {
                score += 1.0
            }
        }

        return score, false
    })

    println(final.ToString())
    fmt.Printf("Done after %d generations\n", generationCounter + 1)
}
