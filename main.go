package main // import "layeh.com/barnard/cmd/barnard"

import (
	"crypto/tls"
	"flag"
	"net"
	"fmt"
	"os"

	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleopenal"
	"layeh.com/gumble/gumbleutil"
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

	// Channel to keep main running until disconnect
	keepAlive := make(chan bool)

	// Attach Event Handlers
	b.Config.Attach(gumbleutil.Listener{
		TextMessage: func(e *gumble.TextMessageEvent) {
			fmt.Printf("Received text message: %s\n", e.Message)
		},

		Connect: func(e *gumble.ConnectEvent) {
			fmt.Printf("Connect.WelcomeMessage: %s\n", *e.WelcomeMessage)
		},

		Disconnect: func(e *gumble.DisconnectEvent) {
			keepAlive <- true
		},
	})

	b.Config.Attach(gumbleutil.AutoBitrate)

	// Connect to server
	var err error
	var client *gumble.Client
	client, err = gumble.DialWithDialer(new(net.Dialer), b.Address, b.Config, &b.TLSConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	// TODO: sensible printout of available channels?
	//fmt.Fprintf(os.Stderr, "%s\n", client.Channels)

	// Add Audio Stream
	if stream, err := gumbleopenal.New(client); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	} else {
		b.Stream = stream
		b.Stream.StartSource()
	}

	<-keepAlive
}
