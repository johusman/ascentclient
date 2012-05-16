package main

import (
    "ascent"
    "ascent/specimens"
    "math/rand"
    "os"
    "time"
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
    index := rand.Intn(len(s.value))
    s.value = s.value[0:index] + string(rand.Intn(126-32)+32) + s.value[index+1:]
}

func main() {
    rand.Seed(int64(time.Now().Nanosecond()))

    start, end := os.Args[1], os.Args[2]

    engine := ascent.New()
    engine.Mutations().Register(randomLetterMutation, 0.99)
    engine.Mutations().SetIdentityChance(0.01)
    engine.SetGenerationCallback(func(winner specimens.Specimen) {
        println(winner.(*mutableString).value)
    })

    println(start)

    final := engine.Run(&mutableString{start}, func(specimen specimens.Specimen) (float32, bool) {
        value := specimen.ToString()

        if value == end {
            return 0, true
        }

        var score float32 = 0.0
        for i := 0; i < len(value) && i < len(end); i++ {
            if value[i] == end[i] {
                score += 1.0
            }
        }

        return score, false
    })

    println(final.ToString())
}
