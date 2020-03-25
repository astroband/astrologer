package commands

import (
	"bufio"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/astroband/astrologer/es"
)

type FastReplayCommandConfig struct {
	UpTo  int
	Count int
}

type FastReplayCommand struct {
	ES     es.Adapter
	Config FastReplayCommandConfig
}

// Execute starts the export process
func (cmd *FastReplayCommand) Execute() {
	log.Printf("Fast replay %d ledgers up to %d\n", cmd.Config.Count, cmd.Config.UpTo)

	pipeFile := "astrologer_meta_stream"

	os.Remove(pipeFile)
	err := syscall.Mkfifo(pipeFile, 0666)

	if err != nil {
		log.Fatal("Make named pipe file error:", err)
	}

	stellarCoreInstance := exec.Command(
		"stellar-core",
		"catchup",
		fmt.Sprintf("%d/%d", cmd.Config.UpTo, cmd.Config.Count),
		"--replay-in-memory",
		"--conf",
		"./stellar-core.cfg",
	)

	// stdout, err := stellarCoreInstance.StdoutPipe()

	// if err != nil {
	// 	log.Fatal(err)
	// }

	go func() {
		if err := stellarCoreInstance.Start(); err != nil {
			log.Fatal(err)
		}
		// scanner := bufio.NewScanner(stdout)

		// for scanner.Scan() {
		// 	log.Println(scanner.Text())
		// }

		// if err := scanner.Err(); err != nil {
		// 	log.Fatal(err)
		// }
	}()

	file, err := os.OpenFile(pipeFile, os.O_RDONLY, os.ModeNamedPipe)

	if err != nil {
		log.Fatal("Open named pipe file error:", err)
	}

	reader := bufio.NewReader(file)
	sizeBytes := make([]byte, 4)

	for {
		bytesRead, err := io.ReadFull(reader, sizeBytes)

		if err != nil {
			log.Fatal(err)
		}

		if bytesRead == 0 {
			break
		}

		sizeBytes[0] &= 0x7f
		size := binary.BigEndian.Uint32(sizeBytes)
		log.Printf("Gonna read %d bytes\n", size)

		data := make([]byte, size)

		_, err = io.ReadFull(reader, data)

		if err == nil {
			sEnc := base64.StdEncoding.EncodeToString(data)
			log.Printf("META: %s\n\n", sEnc)
		} else {
			log.Fatal("Error reading from pipe:", err)
		}
	}
}
