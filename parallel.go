package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/matyas-cyril/go-parallel/parallel"
)

func checkArgs() (*parallel.Parallel, error) {

	var (
		argCmd     string
		argFile    string
		argSep     string
		argCptMax  int
		argNoWait  bool
		argDebug   bool
		argTimeout int
	)

	const defSep string = ","

	flag.StringVar(&argCmd, "cmd", "", "commande - OBLIGATOIRE")
	flag.StringVar(&argFile, "file", "", "fichier contenant")
	flag.StringVar(&argSep, "sep", defSep, fmt.Sprintf("separateur - défaut '%s'", defSep))
	flag.IntVar(&argCptMax, "nbr", 4, "Nombre max de commandes en //. Max 256")
	flag.BoolVar(&argNoWait, "noWait", false, "Activer le mode buffer")
	flag.IntVar(&argTimeout, "timeout", 10, "timeout (sec) d'exécution d'une commande")
	flag.BoolVar(&argDebug, "debug", false, "Activer le mode debug")

	flag.Parse()

	argSep = strings.TrimSpace(argSep)
	if argSep == "" {
		return nil, fmt.Errorf("arg '-sep' must be not empty")
	}

	argCmd = strings.TrimSpace(argCmd)
	if argCmd == "" {
		return nil, fmt.Errorf("arg '-cmd' must defined")
	}

	if argCptMax < 1 || argCptMax > 256 {
		return nil, fmt.Errorf("arg '-nbr' must be in [1-256]")
	}

	if argTimeout < 1 || argTimeout > 3600 {
		return nil, fmt.Errorf("arg '-timeout' must be in [1-3600]")
	}

	data := parallel.Parallel{}
	data.Cmd = argCmd
	data.Separator = argSep
	data.CptMax = argCptMax
	data.Debug = argDebug
	data.Timeout = argTimeout
	data.NoWait = argNoWait

	argFile = strings.TrimSpace(argFile)
	if argFile == "" {
		return nil, fmt.Errorf("arg '-file' must defined")
	}

	if !isFileExist(argFile) {
		return nil, fmt.Errorf("file '%s' not exist", argFile)
	}

	data.File = argFile

	// Déterminer le nombre d'occurences %[N]s dans la commande
	pattern := `\%\[\d+\]s`
	re := regexp.MustCompile(pattern)
	data.Occ = len(re.FindAllString(argCmd, -1))

	return &data, nil
}

func isFileExist(file string) bool {
	if _, err := os.Stat(file); err == nil || os.IsExist(err) {
		return true
	}
	return false
}

func main() {

	data, err := checkArgs()
	if err != nil {
		log.Fatalln(err)
	}

	rslt, err := data.ParseFile(data.File)
	if err != nil {
		log.Fatalln(err)
	}

	jsonData, err := json.Marshal(rslt)
	if err != nil {
		log.Fatalf("error converting to JSON: %s", err)
	}

	os.Stdout.Write(jsonData)
	os.Exit(0)
}
