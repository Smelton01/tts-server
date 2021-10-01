package api_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/smelton01/tts-server/api"

	pb "github.com/smelton01/tts-server/internal/protofiles"
	"github.com/smelton01/tts-server/internal/tts"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterTextToSpeechServer(s, &api.Server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDailer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestRead(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDailer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	testText := "testing 1 2 3"
	client := pb.NewTextToSpeechClient(conn)
	stream, err := client.Read(ctx, &pb.Text{Text: testText})
	if err != nil {
		t.Fatalf("failed to read [%v]: %v", testText, err)
	}

	local, err := testAudio(testText)
	if err != nil {
		t.Fatalf("failed to get local audio: %v", err)
	}

	data := map[int][]byte{}
	for {
		res, err := stream.Recv()
		if err != nil {
			t.Fatalf("could not receive data: %v", err)
		}
		if res.Index == -1 {
			break
		}
		data[int(res.Index)] = res.Audio
	}

	audio := []byte{}
	for i := 0; i < len(data); i++ {
		audio = append(audio, data[i]...)
	}

	tts.PlayAudio(audio)
	if res := bytes.Compare(audio, local[:len(local)-1]); res != 0 {
		t.Errorf("Wrong audio response data:\n\t expected [%v]\n\tgot [%v]", len(local), len(audio))
	}
}

func testAudio(text string) ([]byte, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, fmt.Errorf("could not create tmp file: %v", err)

	}
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("could not close %s: %v", f.Name(), err)
	}

	cmd := exec.Command("gtts-cli", text, "-o", f.Name())
	if data, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("gTTS failed: %s", data)
	}

	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, fmt.Errorf("could not read tmp file: %v", err)
	}
	return data, nil
}
