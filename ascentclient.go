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
    value []rune
}

func (s *mutableString) Clone() specimens.Specimen {
    newValue := make([]rune, len(s.value))
    copy(newValue, s.value)
    return &mutableString{newValue}
}

func (s *mutableString) ToString() string {
    return string(s.value)
}

func randomLetterMutation(specimen specimens.Specimen) {
    s := specimen.(*mutableString)
    if (len(s.value) > 0) {
        index := rand.Intn(len(s.value))
        s.value = append(append(s.value[0:index], rune(rand.Intn(255-32)+32)), s.value[index+1:]...)
    }
}

func dropLetterMutation(specimen specimens.Specimen) {
    s := specimen.(*mutableString)
    if (len(s.value) > 0) {
        index := rand.Intn(len(s.value))
        s.value = append(s.value[0:index], s.value[index+1:]...)
    }
}

func addLetterMutation(specimen specimens.Specimen) {
    s := specimen.(*mutableString)
    index := rand.Intn(len(s.value)+1)
    newValue := make([]rune, index)
    copy(newValue, s.value[0:index])
    s.value = append(append(newValue, rune(rand.Intn(255-32)+32)), s.value[index:]...)
}

func main() {

    if len(os.Args) < 3 {
        println("Usage: ascentclient <start string> <goal string>\n")
        os.Exit(1)
    }

    rand.Seed(int64(time.Now().Nanosecond()))

    start, goal := os.Args[1], os.Args[2]
    goalRunes := []rune(goal)

    engine := ascent.New()
    engine.Mutations().Register(randomLetterMutation, 0.33)
    engine.Mutations().Register(dropLetterMutation, 0.33)
    engine.Mutations().Register(addLetterMutation, 0.33)

    generationCounter := 0
    specimenCounter := 0

    engine.SetGenerationCallback(func(winner specimens.Specimen) (bool){
        println(winner.(*mutableString).ToString())
        generationCounter++
        return winner.(*mutableString).ToString() != goal
    })

    println(start)

    var final specimens.Specimen;
    duration := timeFunction(func() {
        final = engine.Run(4, &mutableString{[]rune(start)}, func(specimen specimens.Specimen) (float32) {
            specimenCounter++
            value := specimen.(*mutableString).value

            score := -float32(math.Abs(float64(len(value) - len(goalRunes)))/4.0)
            for i := 0; i < len(value) && i < len(goalRunes); i++ {
                if value[i] == goalRunes[i] {
                    score += 1.0
                }
            }
            for i := 0; i < len(goalRunes); i++ {
                if strings.ContainsRune(string(value), goalRunes[i]) {
                    score += 1.0
                }
            }

            return score
        })
    })
    println(final.ToString())
    fmt.Printf("Done after %d generations, %d specimens, %f sec\n", generationCounter + 1, specimenCounter, duration.Seconds())
}

func timeFunction(function func()) time.Duration {
    start := time.Now()
    function()
    return time.Now().Sub(start)
}
