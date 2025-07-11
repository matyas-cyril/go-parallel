package parallel

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"sync"
	"time"

	"github.com/matyas-cyril/go-parallel/counter"
)

type Cmd struct {
	duration time.Duration
	err      error
	data     []byte
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (p *Parallel) command(cmd []string, cpt *counter.Counter, wg *sync.WaitGroup, cmdChan chan Cmd) {

	rtn := Cmd{}
	id := genRandID(10)
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.Timeout)*time.Second)

	defer func() {
		rtn.duration = time.Since(startTime)
		if p.Debug {
			log.Println(id, cmd, rtn.duration.Seconds())
		}
		cpt.Dec(1)
		wg.Done()
		cancel()
		cmdChan <- rtn
	}()

	args := cmd[1:]
	output, err := exec.CommandContext(ctx, cmd[0], args...).Output()
	if err != nil {
		rtn.err = fmt.Errorf("%s %s", id, err)
		return
	}

	rtn.data = output

}

func genRandID(length int) string {

	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
