package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand/v2"
	"time"
	"os/exec"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

func isPrintable(s string) bool {
	for _, r := range s {
		if r < 32 || r > 123 {
			return false
		}
	}
	return true
}

func bluetooth() {
	dev, err := linux.NewDevice()
	if err != nil {
		log.Fatalf("Failed to init BLE device: %v", err)
	}

	ble.SetDefaultDevice(dev)

	srvcUUID := ble.MustParse("fb7f3ba1-93ab-4eed-9e5f-6197aead8e07")
	charUUID := ble.MustParse("8c12eedc-9270-424e-9fc7-0dd09a9d13ec")

	service := ble.NewService(srvcUUID)

	// time stamp char
	tsChar := ble.NewCharacteristic(charUUID)
	tsChar.HandleWrite(ble.WriteHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		data := req.Data()

		log.Printf("Received %d bytes: %s", len(data), hex.EncodeToString(data))
		if s := string(data); isPrintable(s) {
			log.Printf("As string %q", s)
		}
	}))

	service.AddCharacteristic(tsChar)
	ble.AddService(service)

	ctx := ble.WithSigHandler(context.Background(), func() {
		log.Println("Signal received --- stopping advertisement")
		dev.Stop()
	})

	go func() {
		for {
			err := ble.AdvertiseNameAndServices(ctx, "TZeroHero", srvcUUID)
			if err != nil {
				log.Printf("Advertise error: %v", err)
			}
			time.Sleep(1 * time.Second)

		}
	}()

	<-ctx.Done()
	log.Println("exiting")
}

func playAudio(filePath string, expectedTimeStamp string) {
	cmd := exec.Command("aplay", filePath)
	if err := cmd.Run() ; err != nil {
		log.Fatalf("Failed to play %s: %v, at %s", filePath, err, expectedTimeStamp)
	}
}

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
	playAudio("./audios/onYourMarks.wav", onYourMarksString)
	//fmt.Printf("On your marks  %s, %s\n", onYourMarksString, time.Now().Format(time.RFC3339))

	waitGetSet := time.Until(getSet)
	time.Sleep(waitGetSet)
	playAudio("./audios/getSet.wav", getSetString)
	//fmt.Printf("Get Set %s, %s\n", getSetString, time.Now().Format(time.RFC3339))

	waitT0 := time.Until(t0)
	time.Sleep(waitT0)
	playAudio("./audios/gun.wav", t0String)
	//fmt.Printf("Gun!! %s, %s \n", t0String, time.Now().Format(time.RFC3339))
}

func main() {
	fmt.Printf("Empezando")
	setGun()
	bluetooth()
}
