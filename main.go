//pipe: use in go generate statements to pipe commands
//commands are executed left to right
//stdout from preceding is piped into stdin of next
//commands are separated by ::
//format is: pipe cmd0 arg0 arg1.. :: cmd1 arg0 arg1.. :: cmdn..
//use pipe -v (verbose) to print output from each command
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println(" valid usage is:")
		fmt.Println("  //go generate pipe cmd0 arg0 arg1.. :: cmd1 arg0 arg1.. :: cmdn..")
		fmt.Println("  //go generate pipe -v cmd0 arg0 arg1.. :: cmd1 arg0 arg1.. :: cmdn..")
		os.Exit(1)
	}
	//parse command string into nxm slice of strings
	verbose := false
	startindex := 0
	if os.Args[1] == "-v" { //verbose
		verbose = true
		startindex = 1
	}
	var cmda [][]string
	for i := startindex; i < len(os.Args); i++ {
		if (i == startindex) || (os.Args[i] == "::") {
			cmda = append(cmda, []string{})
		} else {
			cmda[len(cmda)-1] = append(cmda[len(cmda)-1], os.Args[i])
		}
	}

	var out []byte
	//execute all commands hooking up outputs to inputs
	for i := 0; i < len(cmda); i++ {
		fmt.Println("Command: ", cmda[i][0], cmda[i][1:])
		cmd := exec.Command(cmda[i][0], cmda[i][1:]...)
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Fatal(err)
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		if i > 0 {
			stdin, err := cmd.StdinPipe()
			if err != nil {
				log.Fatal(err)
			}

			go func() {
				defer stdin.Close()
				io.WriteString(stdin, string(out))
			}()
		}
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		newout, _ := ioutil.ReadAll(stdout)
		if verbose {
			fmt.Println(string(newout))
		}
		errtxt, _ := ioutil.ReadAll(stderr)
		if string(errtxt) != "" {
			fmt.Printf("%s\n", errtxt)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
		out = newout
	}
}
