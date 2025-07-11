package parallel

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/matyas-cyril/go-parallel/counter"
)

type Parallel struct {
	Cmd       string
	File      string
	Separator string
	Occ       int
	CptMax    int
	Timeout   int
	NoWait    bool
	Debug     bool
}

type Stat struct {
	Total uint64 `json:"total"`
	Skip  uint64 `json:"skip"`
	OK    uint64 `json:"succes"`
	KO    uint64 `json:"fail"`
}

func (p *Parallel) ParseFile(file string) (*Stat, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		f.Close()
	}()

	var wg sync.WaitGroup // instanciation de notre structure WaitGroup

	total := counter.NewCounter()
	skip := counter.NewCounter()

	ok := counter.NewCounter()
	ko := counter.NewCounter()

	cpt := counter.NewCounter()

	separator := p.Separator
	occ := p.Occ

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())
		elts := strings.Split(line, separator)

		total.Add(1)

		// Cela ne match pas avec le nombre de %[N]s
		if len(elts) != occ {
			skip.Add(1)
			if p.Debug {
				log.Printf("Skip line %d: %s\n", total.Value(), line)
			}
			continue
		}

		if cpt.Value() < uint64(p.CptMax) {
			wg.Add(1)
			cpt.Add(1)

			// Convertir []string en []any - Compatilibité Sprintf à partir de 1.18
			var anyElts []any
			for _, v := range elts {
				anyElts = append(anyElts, v)
			}

			go func() {
				rtnCommand := make(chan Cmd)
				go p.command(strings.Fields(fmt.Sprintf(p.Cmd, anyElts...)), cpt, &wg, rtnCommand)
				resultCommand := <-rtnCommand
				if resultCommand.err != nil {
					log.Println(resultCommand.err)
					ko.Add(1)
					return
				}
				ok.Add(1)
			}()

			continue
		}

		if cpt.Value() >= uint64(p.CptMax) {
			// On attend...
			if !p.NoWait {
				wg.Wait()
			}
			continue
		}

	}

	if cpt.Value() > 0 {
		wg.Wait()
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	stat := Stat{}
	stat.Total = total.Value()
	stat.Skip = skip.Value()
	stat.OK = ok.Value()
	stat.KO = ko.Value()

	return &stat, nil
}
