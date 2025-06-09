package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

func setGun() {
	// TODO: figure out the right math for this
	t0Delay := rand.IntN(31) + 40
	getSetDelay := t0Delay - rand.IntN(5) - 1
	onYourMarksDelay := getSetDelay - rand.IntN(10) - 15

	now := time.Now()
	t0 := now.Add(time.Duration(t0Delay) * time.Second)
	getSet := now.Add(time.Duration(getSetDelay) * time.Second)
	onYourMarks := now.Add(time.Duration(onYourMarksDelay) * time.Second)

	t0String := t0.Format(time.RFC3339)
	getSetString := getSet.Format(time.RFC3339)
	onYourMarksString := onYourMarks.Format(time.RFC3339)

	fmt.Printf("On your marks  %s\n", onYourMarksString)
	fmt.Printf("Set %s\n", getSetString)
	fmt.Printf("Gun!! %s\n", t0String)
}

func main() {
	now := time.Now()

	nowFormated := now.Format(time.RFC3339)

	fmt.Println("Hello World!\n")
	fmt.Printf("Time now\n", nowFormated)
	setGun()
}
