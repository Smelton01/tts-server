/*
Copyright © 2021 Simon Mduduzi Juba scimail09@gmail.com
*/

package tts

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"time"

	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"
	pb "github.com/smelton01/tts-server/internal/protofiles"
)

const output = "output.wav"
const timeout = 200

func Read(message, api string) ([]byte, error) {
	conn, err := grpc.Dial(api, grpc.WithInsecure())

	if err != nil {
		return nil, fmt.Errorf("could not connect to %s: %v", api, err)
	}
	defer conn.Close()

	client := pb.NewTextToSpeechClient(conn)

	text := &pb.Text{Text: message}
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	stream, err := client.Read(ctx, text)

	if err != nil {
		return nil, fmt.Errorf("could not read [%s]: %v", text.Text, err)
	}

	data := map[int][]byte{}
	for {
		res, err := stream.Recv()
		if err != nil {
			return nil, fmt.Errorf("could not receive data: %v", err)
		}
		if res.Index == -1 {
			logrus.Printf("all data received!!!")
			break
		}
		data[int(res.Index)] = res.Audio
	}

	audio := []byte{}
	for i := 0; i < len(data); i++ {
		audio = append(audio, data[i]...)
	}
	return audio, nil
}

func PlayAudio(audio []byte) error {

	if err := ioutil.WriteFile(output, audio, 0666); err != nil {
		return fmt.Errorf("could not write to %s: %v", output, err)
	}
	cmd := exec.Command("afplay", output)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error playing the output file %s: %s", output, err)
	}
	return nil
}
