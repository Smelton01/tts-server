package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"
	pb "github.com/smelton01/tts-server/api"
)

func main() {
	backend := flag.String("b", "localhost:8080", "address of the read backend")
	output := flag.String("o", "output.wav", "wav file where the output will be written")
	input := flag.String("f", "", "input file to read")
	flag.Parse()

	message := flag.Arg(0)
	if *input != "" {
		b, err := ioutil.ReadFile(*input)
		if err != nil {
			log.Fatalf("could not read file %s: %v", *input, err)
		}
		message = string(b)
	}else if flag.NArg() < 1 {
		fmt.Printf("usage:\n\t%s \"text to speak\"\n\t%s -f inputfile.txt\n", os.Args[0], os.Args[0])
		os.Exit(1)
	}

	conn, err := grpc.Dial(*backend, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to %s: %v", *backend, err)
	}
	defer conn.Close()

	client := pb.NewTextToSpeechClient(conn)

	text := &pb.Text{Text: message}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.Read(ctx, text)
	if err != nil {
		log.Fatalf("could not read [%s]: %v", text.Text, err)
	}

	data := []byte{}
	for {
		res, err := stream.Recv()
		if err != nil {
			logrus.Fatal("could not receive data: ", err)
		}
		if res.Audio == nil {
			logrus.Printf("all data received!!!")
			break
		}
		data = append(data, res.Audio...)
	}

	if err := ioutil.WriteFile(*output, data, 0666); err != nil {
		log.Fatalf("could not write to %s: %v", *output, err)
	}
	cmd := exec.Command("afplay", *output)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("error playing the output file %s: %s", *output, err)
	}
}
