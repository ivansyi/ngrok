package client

import (
	"fmt"
	"github.com/inconshreveable/mousetrap"
	"math/rand"
	"ngrok/log"
	"ngrok/util"
	"os"
	"runtime"
	"time"
	//"strconv"
	"crypto/md5"
	"encoding/hex"
)

func init() {
	if runtime.GOOS == "windows" {
		if mousetrap.StartedByExplorer() {
			fmt.Println("Don't double-click ngrok!")
			fmt.Println("You need to open cmd.exe and run it from the command line!")
			time.Sleep(5 * time.Second)
			os.Exit(1)
		}
	}
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Main() {
	// parse options
	opts, err := ParseArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	/*
	   if _, err := strconv.Atoi(opts.subdomain); err != nil {
	       fmt.Println("Illegal subdomain!")
	       os.Exit(1)
	   }

	   if GetMD5Hash(opts.subdomain + opts.protocol) != opts.signature {
	       fmt.Println("Illegal signature!")
	       os.Exit(1)
	   }
	*/

	// set up logging
	log.LogTo(opts.logto, opts.loglevel)

	// read configuration file
	config, err := LoadConfiguration(opts)
	fmt.Printf("Parsed configure:%s-%s", config.Ktvid, config.Dogname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
    //Check signature in command line: md5sum for ktvid+dogname.
    //if GetMD5Hash(config.Ktvid + "1652ec8ffe8245a2ad5b1f586f9baa94" + config.Dogname) != opts.signature {
    //    fmt.Println("Illegal signature!")
    //    os.Exit(1)
    //}

	// seed random number generator
	seed, err := util.RandomSeed()
	if err != nil {
		fmt.Printf("Couldn't securely seed the random number generator!")
		os.Exit(1)
	}
	rand.Seed(seed)

	NewController().Run(config)
}
