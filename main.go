package main

import (
	"bufio"
	"fmt"
	"github.com/projectdiscovery/cdncheck"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	fUtils "github.com/projectdiscovery/utils/file"
	"log"
	"net"
	"os"
	"sync"
)

type options struct {
	IP     string
	List   string
	Silent bool
}

var (
	Opt   = &options{}
	wg    sync.WaitGroup
	input []string
)

func main() {
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription("cdn-finder is a tool that find and separates the CDN IPs from your input")
	flagSet.StringVar(&Opt.IP, "ip", "", "IP for scan")
	flagSet.StringVar(&Opt.List, "list", "", "list of ip")
	flagSet.BoolVar(&Opt.Silent, "silent", false, "show silent output")
	if err := flagSet.Parse(); err != nil {
		log.Fatalf("Could not parse flags: %s\n", err)
	}

	if !Opt.Silent {
		banner()
	}

	// run get user input
	Input()

	if len(input) > 0 {
		Run(input)
	}
}

// Run get []string and exec cdn-finder
func Run(input []string) {
	client, err2 := cdncheck.NewWithCache()
	if err2 != nil {
		fmt.Println("Can't get request to cdnlist")
		return
	}

	for _, ip := range input {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			found, _, _ := client.Check(net.ParseIP(ip))

			if !found {
				fmt.Println(ip)
			}
		}(ip)
	}
	wg.Wait()
}

// Input get user input and append to input variable
func Input() {

	if fUtils.HasStdin() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input = append(input, scanner.Text())
		}
	}

	if Opt.IP != "" {
		input = append(input, Opt.IP)
	}

	if Opt.List != "" {
		lists, err := fUtils.ReadFile(Opt.List)
		if err != nil {
			fmt.Println(err)
			return
		}
		for list := range lists {
			input = append(input, list)
		}
		return
	}
}

func banner() {
	gologger.Print().Msgf(`
 ██████╗██████╗ ███╗   ██╗      ███████╗██╗███╗   ██╗██████╗ ███████╗██████╗
██╔════╝██╔══██╗████╗  ██║      ██╔════╝██║████╗  ██║██╔══██╗██╔════╝██╔══██╗
██║     ██║  ██║██╔██╗ ██║█████╗█████╗  ██║██╔██╗ ██║██║  ██║█████╗  ██████╔╝
██║     ██║  ██║██║╚██╗██║╚════╝██╔══╝  ██║██║╚██╗██║██║  ██║██╔══╝  ██╔══██╗
╚██████╗██████╔╝██║ ╚████║      ██║     ██║██║ ╚████║██████╔╝███████╗██║  ██║
 ╚═════╝╚═════╝ ╚═╝  ╚═══╝      ╚═╝     ╚═╝╚═╝  ╚═══╝╚═════╝ ╚══════╝╚═╝  ╚═╝
`)
	gologger.Print().Msgf("   Created by Sharo_k_h :)\n\n")
}
