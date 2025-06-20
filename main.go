package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/exec"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"

	"github.com/go-audio/wav"
	"github.com/hajimehoshi/oto/v2"
)

type Cue struct {
	name       string
	pcm        []byte
	sampleRate int
	numChans   int
}

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

	play := func() error {
		fmt.Printf("Action expected time  %s, actual time %s\n", expectedTimeStamp, time.Now().Format(time.RFC3339Nano))
		cmd := exec.Command("aplay", filePath)
		return cmd.Run()
	}

	if err := play(); err != nil {
		log.Fatalf("Failed to play %s: %v, at %s", filePath, err, expectedTimeStamp)
	}
}

func loadWAV(filePath string) ([]byte, int, int, error) {
	f, err := os.Open(filePath)

	if err != nil {
		return nil, 0, 0, err
	}

	defer f.Close()

	dec := wav.NewDecoder(f)
	if !dec.IsValidFile() {
		return nil, 0, 0, fmt.Errorf("invalid WAV file: %s", filePath)
	}

	buf, err := dec.FullPCMBuffer()
	if err != nil {
		return nil, 0, 0, err
	}

	pcmBuf := &bytes.Buffer{}
	for _, sample := range buf.Data {
		if err := binary.Write(pcmBuf, binary.LittleEndian, int16(sample)); err != nil {
			return nil, 0, 0, err
		}
	}

	return pcmBuf.Bytes(), int(dec.SampleRate), int(dec.NumChans), nil
}

func setGun() {
	// TODO: figure out the right math for this

	t0Delay := rand.IntN(31) + 40 // second from now to t0
	getSetDelay := t0Delay - rand.IntN(5) - 2
	onYourMarksDelay := getSetDelay - rand.IntN(10) - 15

	now := time.Now()
	t0 := now.Add(time.Duration(t0Delay) * time.Second)
	getSet := now.Add(time.Duration(getSetDelay) * time.Second)
	onYourMarks := now.Add(time.Duration(onYourMarksDelay) * time.Second)
	warmUp := now.Add(3 * time.Second)

	files := []string{"onYourMarks.wav", "onYourMarks.wav", "getSet.wav", "gun.wav"}
	cues := make([]Cue, len(files))

	for i, fn := range files {
		pcm, sr, ch, err := loadWAV("./audios/" + fn)
		if err != nil {
			log.Fatalf("load %s: %v", fn, err)
		}
		cues[i] = Cue{name: fn, pcm: pcm, sampleRate: sr, numChans: ch}
	}

	ctx, ready, err := oto.NewContext(
		cues[0].sampleRate,
		cues[0].numChans,
		2,
	)

	if err != nil {
		log.Fatalf("auido init: %v", err)
	}
	<-ready

	schedule := []struct {
		cue  Cue
		when time.Time
	}{
		{cues[0], warmUp},
		{cues[1], onYourMarks},
		{cues[2], getSet},
		{cues[3], t0},
	}

	for _, item := range schedule {
		c := item.cue
		exp := item.when
		log.Printf("Scheduled %-15s at %s\n", c.name, exp.Format(time.RFC3339Nano))
		time.AfterFunc(time.Until(exp), func() {
			now := time.Now()
			log.Printf(
				"→ %-15s | Exp: %s | Act: %s | Δ: %v\n",
				c.name,
				exp.Format(time.RFC3339Nano),
				now.Format(time.RFC3339Nano),
				now.Sub(exp),
			)
			reader := bytes.NewReader(c.pcm)
			go func() {
				player := ctx.NewPlayer(reader)
				player.Play()
				time.Sleep(time.Duration(len(c.pcm)) * time.Second / time.Duration(c.sampleRate*c.numChans*2))
			}()
		})
	}

	runtime := schedule[len(schedule)-1].when.Sub(time.Now()) + time.Second
	time.Sleep(runtime)
}

func main() {
	setGun()
	bluetooth()
}
