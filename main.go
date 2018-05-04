package main // import "layeh.com/barnard/cmd/barnard"

import (
	"crypto/tls"
	"flag"
	"net"
	"fmt"
	"os"

	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleopenal"
	_ "layeh.com/gumble/opus"
)

type Barney struct {
	Config *gumble.Config
	Client *gumble.Client

	Address   string
	TLSConfig tls.Config

	Stream *gumbleopenal.Stream
}


func main() {
	// Command line flags
	server := flag.String("server", "localhost:64738", "the server to connect to")
	username := flag.String("username", "barney", "the username of the client")
	password := flag.String("password", "", "the password of the server")
	//insecure := flag.Bool("insecure", false, "skip server certificate verification")
	insecure := flag.Bool("insecure", true, "skip server certificate verification")
	certificate := flag.String("certificate", "", "PEM encoded certificate and private key")

	flag.Parse()

	// Initialize
	b := Barney{
		Config: gumble.NewConfig(),
		Address: *server,
	}

	b.Config.Username = *username
	b.Config.Password = *password

	// b.Config.Attach(MyConsolePrintListeners)
		// How to act on few meaningful events? reconnect/exit

	if *insecure {
		b.TLSConfig.InsecureSkipVerify = true
	}
	if *certificate != "" {
		cert, err := tls.LoadX509KeyPair(*certificate, *certificate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		b.TLSConfig.Certificates = append(b.TLSConfig.Certificates, cert)
	}

	var err error
	var client *gumble.Client
	client, err = gumble.DialWithDialer(new(net.Dialer), b.Address, b.Config, &b.TLSConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "%s\n", client.Channels)
	//fmt.Fprintf(os.Stderr, "%s\n", client.AudioEncoder)

	if stream, err := gumbleopenal.New(b.Client); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	} else {
		b.Stream = stream
	}

	for true {}
}
