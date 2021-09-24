package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"

	"github.com/sirupsen/logrus"
	pb "github.com/smelton01/tts-server/api"
	"google.golang.org/grpc"
)

// chunk size to stream 20kb/s
const chunkSize = 20_000

func main() {
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

type server struct{
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
	for index := 0; index < len(data); index += chunkSize{
		if index+chunkSize >= len(data) {
			res := pb.Speech{Audio: data[index:len(data)-1] }
			if err := stream.Send(&res); err != nil {
				return err
			}
			stream.Send(&pb.Speech{})
			return nil 
		}
		res := pb.Speech{Audio: data[index:index+chunkSize]}
		if err := stream.Send(&res); err != nil {
			return err
		}
	}
	return nil
}
