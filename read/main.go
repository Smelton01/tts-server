package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"google.golang.org/grpc"

	pb "github.com/smelton01/tts-server/api"
)

func main() {
	backend := flag.String("b", "localhost:8080", "address of the say backend")
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
	res, err := client.Read(context.Background(), text)
	if err != nil {
		log.Fatalf("could not say %s: %v", text.Text, err)
	}

	if err := ioutil.WriteFile(*output, res.Audio, 0666); err != nil {
		log.Fatalf("could not write to %s: %v", *output, err)
	}
	cmd := exec.Command("afplay", *output)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("error playing the output file %s: %s", *output, err)
	}
}
