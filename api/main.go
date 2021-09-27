package api

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
	pb "github.com/smelton01/tts-server/internal/protofiles"
	"google.golang.org/grpc"
)

// chunk size to stream 20kb/s
const chunkSize = 20_000

func Serve() {
	port := flag.Int("p", 8080, "port to listen to")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logrus.Fatalf("could not listen to port %d: %v", *port, err)
	}
	logrus.Infof("listening to port %d", *port)

	s := grpc.NewServer()
	pb.RegisterTextToSpeechServer(s, server{})
	err = s.Serve(lis)
	if err != nil {
		logrus.Fatalf("could not serve: %v", err)
	}
}

type server struct {
	pb.UnimplementedTextToSpeechServer
}

// Read method uses the gtts-cli package to convert the input
// text to audio and streams the result back to the client
func (server) Read(text *pb.Text, stream pb.TextToSpeech_ReadServer) error {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return fmt.Errorf("could not create tmp file: %v", err)

	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("could not close %s: %v", f.Name(), err)
	}

	cmd := exec.Command("gtts-cli", text.Text, "-o", f.Name())
	if data, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("gTTS failed: %s", data)
	}

	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return fmt.Errorf("could not read tmp file: %v", err)
	}
	// Stream audio file in chunks
	var wg sync.WaitGroup
	for index, count := 0, 0; index < len(data); index, count = index+chunkSize, count+1 {
		end := index + chunkSize
		if end >= len(data) {
			end = len(data) - 1
		}
		wg.Add(1)
		go func(count, end, index int) {
			defer wg.Done()
			if err := sendChunk(wg, stream, data[index:end], int32(count)); err != nil {
				panic(err)
			}
		}(count, end, index)
	}
	wg.Wait()
	stream.Send(&pb.Speech{Audio: []byte{}, Index: -1})
	return nil
}

func sendChunk(wg sync.WaitGroup, stream pb.TextToSpeech_ReadServer, data []byte, index int32) error {
	res := pb.Speech{Audio: data, Index: index}
	err := stream.Send(&res)
	if err != nil {
		logrus.Print(err)
		return err
	}
	return nil
}
