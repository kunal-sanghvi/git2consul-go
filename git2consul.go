package main

import (
	"github.com/kunal-sanghvi/git2consul-go/actor"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	start := time.Now()
	var configPath string
	var privateKeyPath string

	homeDir, _ := os.UserHomeDir()
	defaultPemFile := homeDir + "/.ssh/id_rsa"
	defaultConfigPath, _ := filepath.Abs("config.json")

	app := cli.NewApp()
	app.Name = "git2consul-go"
	app.Usage = "Populate config from git to consul... Faster!!"
	app.Description = "Fast and Furious Git2Consul - Because it's written in Go :p"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "config, c",
			Destination: &configPath,
			Usage: "Specify config `FILE`, defaults to config.json in current directory",
			Value: defaultConfigPath,
		},
		cli.StringFlag{
			Name: "pemfile, p",
			Destination: &privateKeyPath,
			Usage: "Specify pem `FILE`, defaults to ~/.ssh/id_rsa in home directory",
			Value: defaultPemFile,
		},
	}
	app.Action = func(c *cli.Context) error {
		actor.Git2Consul(configPath, privateKeyPath)
		elapsed := time.Since(start)
		log.Printf("git2Consul-go took %s", elapsed)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
