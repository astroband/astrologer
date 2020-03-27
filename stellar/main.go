package stellar

import (
	"bufio"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/astroband/astrologer/support"
	"github.com/guregu/null"
	log "github.com/sirupsen/logrus"
	"github.com/stellar/go/xdr"
)

func StreamLedgers(firstLedger, lastLedger int) chan xdr.LedgerCloseMeta {
	const pipeFile = "astrologer_meta_stream"
	os.Remove(pipeFile)
	ch := make(chan xdr.LedgerCloseMeta)

	err := syscall.Mkfifo(pipeFile, 0666)

	if err != nil {
		log.Fatal("Make named pipe file error:", err)
	}

	stellarCoreInstance := exec.Command(
		"stellar-core",
		"catchup",
		fmt.Sprintf("%d/%d", lastLedger, lastLedger-firstLedger+1),
		"--replay-in-memory",
		"--conf",
		"./stellar-core/pubnet.cfg",
	)
	stellarCoreLogger := log.WithFields(log.Fields{"process": "stellar-core"})
	stellarCoreStdOut, err := stellarCoreInstance.StdoutPipe()

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := stellarCoreInstance.Start(); err != nil {
			stellarCoreLogger.Fatal(err)
		}

		scanner := bufio.NewScanner(stellarCoreStdOut)

		for scanner.Scan() {
			stellarCoreLogger.Println(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			stellarCoreLogger.Fatal(err)
		}
	}()

	pipe, err := os.OpenFile(pipeFile, os.O_RDONLY, os.ModeNamedPipe)

	if err != nil {
		log.Fatal("Open named pipe file error:", err)
	}

	reader := bufio.NewReader(pipe)
	sizeBytes := make([]byte, 4)

	readerLogger := log.WithFields(log.Fields{"process": "reader"})

	go func() {
		defer close(ch)

		for {
			bytesRead, err := io.ReadFull(reader, sizeBytes)

			// No more to read
			if bytesRead == 0 {
				readerLogger.Info("Reached the end of the stream")
				break
			}

			if err != nil {
				readerLogger.Fatal("Failed to read from pipe", err)
			}

			sizeBytes[0] &= 0x7f
			size := binary.BigEndian.Uint32(sizeBytes)

			data := make([]byte, size)

			_, err = io.ReadFull(reader, data)

			if err == nil {
				var meta xdr.LedgerCloseMeta
				meta.UnmarshalBinary(data)

				readerLogger.Info("Writing to the channel...")
				ch <- meta
			} else {
				readerLogger.Fatal("Error reading from pipe:", err)
			}
		}

		err = pipe.Close()

		if err != nil {
			readerLogger.Println(err)
		}
	}()

	return ch
}

func MemoValue(memo xdr.Memo) null.String {
	var (
		value string
		valid bool
	)
	switch memo.Type {
	case xdr.MemoTypeMemoNone:
		value, valid = "", false
	case xdr.MemoTypeMemoText:
		scrubbed := support.Utf8Scrub(memo.MustText())
		notnull := strings.Join(strings.Split(scrubbed, "\x00"), "")
		value, valid = notnull, true
	case xdr.MemoTypeMemoId:
		value, valid = fmt.Sprintf("%d", memo.MustId()), true
	case xdr.MemoTypeMemoHash:
		hash := memo.MustHash()
		value, valid =
			base64.StdEncoding.EncodeToString(hash[:]),
			true
	case xdr.MemoTypeMemoReturn:
		hash := memo.MustRetHash()
		value, valid =
			base64.StdEncoding.EncodeToString(hash[:]),
			true
	default:
		panic(fmt.Errorf("invalid memo type: %v", memo.Type))
	}

	return null.NewString(value, valid)
}

func OperationMeta(txMeta xdr.TransactionMeta, opIndex int) (result *xdr.OperationMeta) {
	if v1, ok := txMeta.GetV1(); ok {
		ops := v1.Operations
		return &ops[opIndex]
	}

	ops, ok := txMeta.GetOperations()
	if !ok {
		return nil
	}

	return &ops[opIndex]
}
