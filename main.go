package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

func setGun() {
	// TODO: figure out the right math for this

	t0Delay := rand.IntN(31) + 40 // second from now to t0
	getSetDelay := t0Delay - rand.IntN(5) - 1
	onYourMarksDelay := getSetDelay - rand.IntN(10) - 15

	now := time.Now()
	t0 := now.Add(time.Duration(t0Delay) * time.Second)
	getSet := now.Add(time.Duration(getSetDelay) * time.Second)
	onYourMarks := now.Add(time.Duration(onYourMarksDelay) * time.Second)

	t0String := t0.Format(time.RFC3339)
	getSetString := getSet.Format(time.RFC3339)
	onYourMarksString := onYourMarks.Format(time.RFC3339)

	waitOnYourMarks := time.Until(onYourMarks)
	time.Sleep(waitOnYourMarks)
	fmt.Printf("On your marks  %s, %s\n", onYourMarksString, time.Now().Format(time.RFC3339))

	waitGetSet := time.Until(getSet)
	time.Sleep(waitGetSet)
	fmt.Printf("Get Set %s, %s\n", getSetString, time.Now().Format(time.RFC3339))

	waitT0 := time.Until(t0)
	time.Sleep(waitT0)
	fmt.Printf("Gun!! %s, %s \n", t0String, time.Now().Format(time.RFC3339))
}

func main() {
	fmt.Println("Hello World!\n")

	setGun()
}
