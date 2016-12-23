package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"syscall"
	"time"

	logging "github.com/op/go-logging"
	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/gui"
	"github.com/skycoin/skycoin/src/util"
	"github.com/skycoin/skycoin/src/visor/blockdb"
)

//"github.com/wudaofan/mzcoin/src/cli"

//"github.com/wudaofan/mzcoin/src/wallet"

var (
	logger     = logging.MustGetLogger("main")
	logFormat  = "[mzcoin.%{module}:%{level}] %{message}"
	logModules = []string{
		"main",
		"daemon",
		"coin",
		"gui",
		"util",
		"visor",
		"wallet",
		"gnet",
		"pex",
		"webrpc",
	}

	//clear these after loading [????]
	GenesisSignatureStr = "ab58cd355f2e5b8c18ecfedba67d9410385c27588f0dc25f6cf18cc4fa7164456673e9d471e041fc25a74c70b33263654632a5ec058f532eb821f403976379ac01"
	GenesisAddressStr   = "ppu2zgS1H2aheeMNgVpUXjHeJJ7Uov3i4W"
	BlockchainPubkeyStr = "02e2016590cf0036a47482773316ec1d521425fcd214cd02adca556751fafb291e"
	BlockchainSeckeyStr = ""

	GenesisTimestamp  uint64 = 0
	GenesisCoinVolume uint64 = 300e12

	//use port 6001
	// DefaultServers = []string{
	// 	"40.74.80.119:6001",
	DefaultConnections = []string{
		"121.41.103.148:7000",
		"120.77.69.188:7000",
	}
)

// Command line interface arguments

type Config struct {
	// Disable peer exchange
	DisablePEX bool
	// Don't make any outgoing connections
	DisableOutgoingConnections bool
	// Don't allowing incoming connections
	DisableIncomingConnections bool
	// Disables networking altogether
	DisableNetworking bool
	// Only run on localhost and only connect to others on localhost
	LocalhostOnly bool
	// Which address to serve on. Leave blank to automatically assign to a
	// public interface
	Address string
	//gnet uses this for TCP incoming and outgoing
	Port int
	//max connections to maintain
	MaxConnections int
	// How often to make outgoing connections
	OutgoingConnectionsRate time.Duration
	// Wallet Address Version
	//AddressVersion string
	// Remote web interface
	WebInterface      bool
	WebInterfacePort  int
	WebInterfaceAddr  string
	WebInterfaceCert  string
	WebInterfaceKey   string
	WebInterfaceHTTPS bool

	// Launch System Default Browser after client startup
	LaunchBrowser bool

	// If true, print the configured client web interface address and exit
	PrintWebInterfaceAddress bool

	// Data directory holds app data -- defaults to ~/.mzcoin
	DataDirectory string
	// GUI directory contains assets for the html gui
	GUIDirectory string
	// Logging
	LogLevel logging.Level
	ColorLog bool
	// This is the value registered with flag, it is converted to LogLevel after parsing
	logLevel string

	// Wallets
	// Defaults to ${DataDirectory}/wallets/
	WalletDirectory string
	BlockchainFile  string
	BlockSigsFile   string

	// Centralized network configuration

	RunMaster bool

	GenesisSignature cipher.Sig
	GenesisTimestamp uint64
	GenesisAddress   cipher.Address

	BlockchainPubkey cipher.PubKey
	BlockchainSeckey cipher.SecKey

	/* Developer options */

	// Enable cpu profiling
	ProfileCPU bool
	// Where the file is written to
	ProfileCPUFile string
	// HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
	HTTPProf bool
	// Will force it to connect to this ip:port, instead of waiting for it
	// to show up as a peer
	ConnectTo string
}

func (c *Config) register() {
	flag.BoolVar(&c.DisablePEX, "disable-pex", c.DisablePEX,
		"disable PEX peer discovery")
	flag.BoolVar(&c.DisableOutgoingConnections, "disable-outgoing",
		c.DisableOutgoingConnections, "Don't make outgoing connections")
	flag.BoolVar(&c.DisableIncomingConnections, "disable-incoming",
		c.DisableIncomingConnections, "Don't make incoming connections")
	flag.BoolVar(&c.DisableNetworking, "disable-networking",
		c.DisableNetworking, "Disable all network activity")
	flag.StringVar(&c.Address, "address", c.Address,
		"IP Address to run application on. Leave empty to default to a public interface")
	flag.IntVar(&c.Port, "port", c.Port, "Port to run application on")
	flag.BoolVar(&c.WebInterface, "web-interface", c.WebInterface,
		"enable the web interface")
	flag.IntVar(&c.WebInterfacePort, "web-interface-port",
		c.WebInterfacePort, "port to serve web interface on")
	flag.StringVar(&c.WebInterfaceAddr, "web-interface-addr",
		c.WebInterfaceAddr, "addr to serve web interface on")
	flag.StringVar(&c.WebInterfaceCert, "web-interface-cert",
		c.WebInterfaceCert, "cert.pem file for web interface HTTPS. "+
			"If not provided, will use cert.pem in -data-directory")
	flag.StringVar(&c.WebInterfaceKey, "web-interface-key",
		c.WebInterfaceKey, "key.pem file for web interface HTTPS. "+
			"If not provided, will use key.pem in -data-directory")
	flag.BoolVar(&c.WebInterfaceHTTPS, "web-interface-https",
		c.WebInterfaceHTTPS, "enable HTTPS for web interface")
	flag.BoolVar(&c.LaunchBrowser, "launch-browser", c.LaunchBrowser,
		"launch system default webbrowser at client startup")
	flag.BoolVar(&c.PrintWebInterfaceAddress, "print-web-interface-address",
		c.PrintWebInterfaceAddress, "print configured web interface address and exit")
	flag.StringVar(&c.DataDirectory, "data-dir", c.DataDirectory,
		"directory to store app data (defaults to ~/.mzcoin)")
	flag.StringVar(&c.ConnectTo, "connect-to", c.ConnectTo,
		"connect to this ip only")
	flag.BoolVar(&c.ProfileCPU, "profile-cpu", c.ProfileCPU,
		"enable cpu profiling")
	flag.StringVar(&c.ProfileCPUFile, "profile-cpu-file",
		c.ProfileCPUFile, "where to write the cpu profile file")
	flag.BoolVar(&c.HTTPProf, "http-prof", c.HTTPProf,
		"Run the http profiling interface")
	flag.StringVar(&c.logLevel, "log-level", c.logLevel,
		"Choices are: debug, info, notice, warning, error, critical")
	flag.BoolVar(&c.ColorLog, "color-log", c.ColorLog,
		"Add terminal colors to log output")
	flag.StringVar(&c.GUIDirectory, "gui-dir", c.GUIDirectory,
		"static content directory for the html gui")

	//Key Configuration Data
	flag.BoolVar(&c.RunMaster, "master", c.RunMaster,
		"run the daemon as blockchain master server")

	flag.StringVar(&BlockchainPubkeyStr, "master-public-key", BlockchainPubkeyStr,
		"public key of the master chain")
	flag.StringVar(&BlockchainSeckeyStr, "master-secret-key", BlockchainSeckeyStr,
		"secret key, set for master")

	flag.StringVar(&GenesisAddressStr, "genesis-address", GenesisAddressStr,
		"genesis address")
	flag.StringVar(&GenesisSignatureStr, "genesis-signature", GenesisSignatureStr,
		"genesis block signature")
	flag.Uint64Var(&c.GenesisTimestamp, "genesis-timestamp", c.GenesisTimestamp,
		"genesis block timestamp")

	flag.StringVar(&c.WalletDirectory, "wallet-dir", c.WalletDirectory,
		"location of the wallet files. Defaults to ~/.mzcoin/wallet/")

	flag.StringVar(&c.BlockchainFile, "blockchain-file", c.BlockchainFile,
		"location of the blockchain file. Default to ~/.mzcoin/blockchain.bin")
	flag.StringVar(&c.BlockSigsFile, "blocksigs-file", c.BlockSigsFile,
		"location of the block signatures file. Default to ~/.mzcoin/blockchain.sigs")

	flag.DurationVar(&c.OutgoingConnectionsRate, "connection-rate",
		c.OutgoingConnectionsRate, "How often to make an outgoing connection")
	flag.BoolVar(&c.LocalhostOnly, "localhost-only", c.LocalhostOnly,
		"Run on localhost and only connect to localhost peers")
	//flag.StringVar(&c.AddressVersion, "address-version", c.AddressVersion,
	//	"Wallet address version. Options are 'test' and 'main'")
}

func (c *Config) Parse() {
	c.register()
	flag.Parse()
	c.postProcess()
}

func (c *Config) postProcess() {
	var err error
	if GenesisSignatureStr != "" {
		c.GenesisSignature, err = cipher.SigFromHex(GenesisSignatureStr)
		panicIfError(err, "Invalid Signature")
	}
	if GenesisAddressStr != "" {
		c.GenesisAddress, err = cipher.DecodeBase58Address(GenesisAddressStr)
		panicIfError(err, "Invalid Address")
	}
	if BlockchainPubkeyStr != "" {
		c.BlockchainPubkey, err = cipher.PubKeyFromHex(BlockchainPubkeyStr)
		panicIfError(err, "Invalid Pubkey")
	}
	if BlockchainSeckeyStr != "" {
		c.BlockchainSeckey, err = cipher.SecKeyFromHex(BlockchainSeckeyStr)
		panicIfError(err, "Invalid Seckey")
		BlockchainSeckeyStr = ""
	}
	if BlockchainSeckeyStr != "" {
		c.BlockchainSeckey = cipher.SecKey{}
	}

	c.DataDirectory = util.InitDataDir(c.DataDirectory)
	if c.WebInterfaceCert == "" {
		c.WebInterfaceCert = filepath.Join(c.DataDirectory, "cert.pem")
	}
	if c.WebInterfaceKey == "" {
		c.WebInterfaceKey = filepath.Join(c.DataDirectory, "key.pem")
	}

	if c.BlockchainFile == "" {
		c.BlockchainFile = filepath.Join(c.DataDirectory, "blockchain.bin")
	}
	if c.BlockSigsFile == "" {
		c.BlockSigsFile = filepath.Join(c.DataDirectory, "blockchain.sigs")
	}
	if c.WalletDirectory == "" {
		c.WalletDirectory = filepath.Join(c.DataDirectory, "wallets/")
	}

	ll, err := logging.LogLevel(c.logLevel)
	panicIfError(err, "Invalid -log-level %s", c.logLevel)
	c.LogLevel = ll

}

func panicIfError(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Panicf(msg+": %v", append(args, err)...)
	}
}

func printProgramStatus() {
	fn := "goroutine.prof"
	logger.Debug("Writing goroutine profile to %s", fn)
	p := pprof.Lookup("goroutine")
	f, err := os.Create(fn)
	defer f.Close()
	if err != nil {
		logger.Error("%v", err)
		return
	}
	err = p.WriteTo(f, 2)
	if err != nil {
		logger.Error("%v", err)
		return
	}
}

func catchInterrupt(quit chan<- int) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	signal.Stop(sigchan)
	quit <- 1
}

// Catches SIGUSR1 and prints internal program state
func catchDebug() {
	sigchan := make(chan os.Signal, 1)
	//signal.Notify(sigchan, syscall.SIGUSR1)
	signal.Notify(sigchan, syscall.Signal(0xa)) // SIGUSR1 = Signal(0xa)
	for {
		select {
		case <-sigchan:
			printProgramStatus()
		}
	}
}

func initLogging(level logging.Level, color bool) {
	format := logging.MustStringFormatter(logFormat)
	logging.SetFormatter(format)
	for _, s := range logModules {
		logging.SetLevel(level, s)
	}
	stdout := logging.NewLogBackend(os.Stdout, "", 0)
	stdout.Color = color
	logging.SetBackend(stdout)
}

func initProfiling(httpProf, profileCPU bool, profileCPUFile string) {
	if profileCPU {
		f, err := os.Create(profileCPUFile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if httpProf {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
}

var devConfig Config = Config{
	// Disable peer exchange
	DisablePEX: true,
	// Don't make any outgoing connections
	DisableOutgoingConnections: false,
	// Don't allowing incoming connections
	DisableIncomingConnections: false,
	// Disables networking altogether
	DisableNetworking: false,
	// Only run on localhost and only connect to others on localhost
	LocalhostOnly: false,
	// Which address to serve on. Leave blank to automatically assign to a
	// public interface
	Address: "",
	//gnet uses this for TCP incoming and outgoing
	Port: 7000,

	MaxConnections: 16,
	// How often to make outgoing connections, in seconds
	OutgoingConnectionsRate: time.Second * 5,
	// Wallet Address Version
	//AddressVersion: "test",
	// Remote web interface
	WebInterface:             true,
	WebInterfacePort:         7420,
	WebInterfaceAddr:         "127.0.0.1",
	WebInterfaceCert:         "",
	WebInterfaceKey:          "",
	WebInterfaceHTTPS:        false,
	PrintWebInterfaceAddress: false,
	LaunchBrowser:            true,
	// Data directory holds app data -- defaults to ~/.mzcoin
	DataDirectory: ".mzcoin",
	// Web GUI static resources
	GUIDirectory: "./src/gui/static/",
	// Logging
	LogLevel: logging.DEBUG,
	ColorLog: true,
	logLevel: "DEBUG",

	// Wallets
	WalletDirectory: "",
	BlockchainFile:  "",
	BlockSigsFile:   "",

	// Centralized network configuration
	RunMaster:        false,
	BlockchainPubkey: cipher.PubKey{},
	BlockchainSeckey: cipher.SecKey{},

	GenesisAddress:   cipher.Address{},
	GenesisTimestamp: GenesisTimestamp,
	GenesisSignature: cipher.Sig{},

	/* Developer options */

	// Enable cpu profiling
	ProfileCPU: false,
	// Where the file is written to
	ProfileCPUFile: "mzcoin.prof",
	// HTTP profiling interface (see http://golang.org/pkg/net/http/pprof/)
	HTTPProf: false,
	// Will force it to connect to this ip:port, instead of waiting for it
	// to show up as a peer
	ConnectTo: "",
}

func configureDaemon(c *Config) daemon.Config {
	//cipher.SetAddressVersion(c.AddressVersion)

	dc := daemon.NewConfig()
	dc.Peers.DataDirectory = c.DataDirectory
	dc.Peers.Disabled = c.DisablePEX
	dc.Daemon.DisableOutgoingConnections = c.DisableOutgoingConnections
	dc.Daemon.DisableIncomingConnections = c.DisableIncomingConnections
	dc.Daemon.DisableNetworking = c.DisableNetworking
	dc.Daemon.Port = c.Port
	dc.Daemon.Address = c.Address
	dc.Daemon.LocalhostOnly = c.LocalhostOnly
	dc.Daemon.OutgoingMax = c.MaxConnections

	daemon.DefaultConnections = DefaultConnections

	if c.OutgoingConnectionsRate == 0 {
		c.OutgoingConnectionsRate = time.Millisecond
	}
	dc.Daemon.OutgoingRate = c.OutgoingConnectionsRate

	dc.Visor.Config.BlockchainFile = c.BlockchainFile
	dc.Visor.Config.BlockSigsFile = c.BlockSigsFile

	dc.Visor.Config.IsMaster = c.RunMaster

	dc.Visor.Config.BlockchainPubkey = c.BlockchainPubkey
	dc.Visor.Config.BlockchainSeckey = c.BlockchainSeckey

	dc.Visor.Config.GenesisAddress = c.GenesisAddress
	dc.Visor.Config.GenesisSignature = c.GenesisSignature
	dc.Visor.Config.GenesisTimestamp = c.GenesisTimestamp
	dc.Visor.Config.GenesisCoinVolume = GenesisCoinVolume
	return dc
}

func Run(c *Config) {

	c.GUIDirectory = util.ResolveResourceDirectory(c.GUIDirectory)

	scheme := "http"
	if c.WebInterfaceHTTPS {
		scheme = "https"
	}
	host := fmt.Sprintf("%s:%d", c.WebInterfaceAddr, c.WebInterfacePort)
	fullAddress := fmt.Sprintf("%s://%s", scheme, host)
	logger.Critical("Full address: %s", fullAddress)

	if c.PrintWebInterfaceAddress {
		fmt.Println(fullAddress)
		return
	}

	initProfiling(c.HTTPProf, c.ProfileCPU, c.ProfileCPUFile)
	initLogging(c.LogLevel, c.ColorLog)

	// start the block db.
	blockdb.Start()
	defer blockdb.Stop()

	// start the transaction db.
	// transactiondb.Start()
	// defer transactiondb.Stop()

	// If the user Ctrl-C's, shutdown properly
	quit := make(chan int)
	go catchInterrupt(quit)
	// Watch for SIGUSR1
	go catchDebug()

	gui.InitWalletRPC(c.WalletDirectory)

	dconf := configureDaemon(c)
	d := daemon.NewDaemon(dconf)

	stopDaemon := make(chan int)
	go d.Start(stopDaemon)

	// start the webrpc
	closingC := make(chan struct{})
	go webrpc.Start("0.0.0.0:7430",
		webrpc.ChanBuffSize(1000),
		webrpc.Gateway(d.Gateway),
		webrpc.ThreadNum(1000),
		webrpc.Quit(closingC))

	// Debug only - forces connection on start.  Violates thread safety.
	if c.ConnectTo != "" {
		_, err := d.Pool.Pool.Connect(c.ConnectTo)
		if err != nil {
			log.Panic(err)
		}
	}

	if c.WebInterface {
		var err error
		if c.WebInterfaceHTTPS {
			// Verify cert/key parameters, and if neither exist, create them
			errs := util.CreateCertIfNotExists(host, c.WebInterfaceCert, c.WebInterfaceKey, "Skycoind")
			if len(errs) != 0 {
				for _, err := range errs {
					logger.Error(err.Error())
				}
				logger.Error("gui.CreateCertIfNotExists failure")
				os.Exit(1)
			}

			err = gui.LaunchWebInterfaceHTTPS(host, c.GUIDirectory, d, c.WebInterfaceCert, c.WebInterfaceKey)
		} else {
			err = gui.LaunchWebInterface(host, c.GUIDirectory, d)
		}

		if err != nil {
			logger.Error(err.Error())
			logger.Error("Failed to start web GUI")
			os.Exit(1)
		}

		if c.LaunchBrowser {
			go func() {
				// Wait a moment just to make sure the http interface is up
				time.Sleep(time.Millisecond * 100)

				logger.Info("Launching System Browser with %s", fullAddress)
				if err := util.OpenBrowser(fullAddress); err != nil {
					logger.Error(err.Error())
				}
			}()
		}
	}
	/*
		time.Sleep(5)
		tx := InitTransaction()
		_ = tx
		err, _ = d.Visor.Visor.InjectTxn(tx)
		if err != nil {
			log.Panic(err)
		}
	*/

	//first transaction
	if c.RunMaster == true {
		//log.Printf("BLOCK SEQ= %d \n", d.Visor.Visor.Blockchain.Head().Seq())
		go func() {
			for d.Visor.Visor.Blockchain.Head().Seq() < 1 {
				time.Sleep(5 * time.Second)
				tx := InitTransaction()
				err, _ := d.Visor.Visor.InjectTxn(tx)
				if err != nil {
					log.Printf("%s\n", err)
				}
			}
		}()
	}

	<-quit
	stopDaemon <- 1
	close(closingC)

	logger.Info("Shutting down")
	d.Shutdown()
	logger.Info("Goodbye")
}

func main() {

	/*
		mzcoin.Run(&cli.DaemonArgs)
	*/

	/*
	   mzcoin.Run(&cli.ClientArgs)
	   stop := make(chan int)
	   <-stop
	*/

	//mzcoin.Run(&cli.DevArgs)
	devConfig.Parse()
	Run(&devConfig)
}

//addresses for storage of coins
var AddrList []string = []string{
	"bhfN2SXFoJdgfd6k2MXbNQXZix3pQTEfJb",
	"xa4F9twZsHjyBMvrDuQtimh34Wn4xAWRb5",
	"2LkUBbUf2LvCnsPjitSn3PSiP3vSep2L6Pd",
	"2W5Wb6nEe9KByv486xujBK3uQKYbpssYby1",
	"wwszxeD6hWNveQLvqYmrz6CcQYr7H9Z5W2",
	"2PYexNtTQEk1V8gjwAqRunxwQUgdxygWE4p",
	"2k6UHuzFSLZgaz9XPxPV8rCrLJRtJ67oQBe",
	"2fymCfVQLaU1EUaNm5fyzuLHPAhYYbBFhjM",
	"RBrHhZt2WoJEN8Tg9WkxaaRyjkV476Efvj",
	"2ay6NmBcSTAMscm8cuXTvu258KNS352kteq",
	"WeiHBmqS4ob2xuAU15sVpQwrspxKFzUuLR",
	"2irzCEdSikB4U161xwvhhSD4HYdKybWfGYw",
	"23PApg2FomMbZ4pS2DjyWL6PPrM7219Tun8",
	"2Rk8BMLwreSHTNBXXLPXJZVfupxEwUyaLtT",
	"wf76eyZvVybrAqCwV1A4BwvYWNPFktieTy",
	"873xgFM5UpzQnnTdM2REtKHYi2FM1Zt2R2",
	"27rYHhBSfzVSCjnT3bRn2xLXuM9x53UBYEZ",
	"2DWHBRuJZTG9BBnjsDNArJJ7XDnDvye7PMF",
	"QrTVJgz8knNyMHWaaHEfU1sWA3TjzGxNGA",
	"GYcigvRmUcNi4q2yqYjbcKHc5WVVNUddQF",
	"yicKja6yRyRwPS9Hugjg7nztAKXtJwh8xm",
	"2XmZsewt88nxnRJPTWufknBrsChD6jzS343",
	"NBBnK5fLDLMvAX6VFtt4qwnbHeGkpqDxoC",
	"iAVptbD8Np9E1tX4vc14oNRTa5LCRfaibv",
	"7SjN9L6rUs8mPeh51gXAfHofY3D3kpssMy",
	"2g39jwf8h2Q9NXjjtJnZfqpafKZNKs3QaY",
	"21hEAszZL1VkLQaxCDGLQ3FqaaqDNjgHrK1",
	"WJ29xWNNFsVzhsB8cqh1ekDg9s7wHPLLZX",
	"2DbYEiUm7EpHdjVfd7g2coWAjB8WyXPpHWk",
	"hLeQurePKTvGTsaVxsp5wWbLbqnzLzqsDQ",
	"2CWoQDwm6YNKDcDrXZg1FSFCXFPnhDWcrot",
	"7t2dsWpxpqVtww7q7YNxCrschWSJYdJnpZ",
	"2c8csipycrUFkFFh3i6Rzkf2QvYhVUru5Rq",
	"ENt1iuQs2prq9cDevvc7xAGV9kMHsXirGV",
	"K2EQAZMm7rjzXnBhb1CvpbidZron2YPcss",
	"2eH1DQxS83RB7rgGbx5JUNMhVZTLdiPSyuH",
	"wFgv2JLrTYhkzExAfGcvgb8CZhKhQbz7TY",
	"PHCDfMBjd4uWfQ6EqVKPU571kHnrrE55N8",
	"CSULm5QWcx3E7Qgyi7doiVLg2JkJmBx5HH",
	"2drH12moEXeR4Uj3jBpFSiXXg53mmHpop3U",
	"DWycLBPX88qythtZFUjAUe9nW1rTL2137M",
	"5nVPC9kAWcyyCA27RxT9SvS4QDqefERovY",
	"2StQ2tFaLgS6Zv6ZKh9vKqRUTEag3v53gAX",
	"2LfvLvg9MRcFZ1J3KSQZSfU66gnBKYPoMvC",
	"pggCTaNQ2YGtLHMd2K8bzKGomwd4JZzwHZ",
	"2d1JMT7s5B3rEPSBXAnh7JXGomgbWJ9xsf9",
	"2mCSXePhCPrCJyok6kmCAiMmcfd2WACurDK",
	"4s15J311npW9gE35WvzPcdntPv3sp83DUw",
	"2eVtPJbKkDJkcNnEV3YKnwv2MtaTvLUbc3t",
	"3WgJz3K8HYxsQr7qBM7CT9k7SSw5h8NFHH",
	"gLw7TQHJ69zk8vWiViW8n3VumaCFgcZZsC",
	"rfqBRPWmFudi1PEKZToA8hX82a392uLn2q",
	"24gFre2Wyt5E2L45qW8jULq38MFxtJ8d41R",
	"2Z1V2UpEWpBbMuUbbGGs2h9Lnt7dL1pVdow",
	"2dcrSobzPtpjf4U36Btf1CESEjnSATZpkre",
	"bH4jqy8qsWxfKV21Zc76qGACFioGVrNPMN",
	"FapGCWGob6zxiayKJsHX1WTmRUKStu8gAG",
	"UF3Zs4YrEZ1qW8TXZqmeniLbB91UDtsZCs",
	"2EGnRrze66vzBosPLz7yzMSZ7L6rCgNgAyS",
	"LxT1Qy4GhxyeuU6o5AXhNyUyVmxhDxrWBa",
	"2etQcyH8Ex9wyDHQ15KXPZkn9jEQS1VdDQJ",
	"PY9XsZxAnygh2wPzuPcEh3qPyBWTvHfGVB",
	"JHLXWvVaJrwQkQXmU4pMBXbkp53BEUDBEJ",
	"FUyrJL2nGW1L5Exbv7B2MzxRPcrcji8TRp",
	"2P3E1jhgFoZgtC2YZnJTjFXZNqvYL3N7ue5",
	"2771xbXanzeydayszcEzPyYF7HD72h4BggQ",
	"MaSGs1B4sGivEYKcPh4LfWRz3fVjbvVMmh",
	"ruDmeGYSW62gB9rQj4Yy974HLwbMtgog1j",
	"2c7PNBwwTYP9newxcxdwqGRnSXBS2ZDD31X",
	"2GKZnD2vfWR5X7ZUSsafL27xSme8LmJ9yH6",
	"2ZEXuo6n68XetWZqmgaPaFEPmyVatwVvuDz",
	"2GFdxbv3C7gnFoJ4UBWqyai1LDtroappku",
	"BDLgoo4jb8gMsHexTNzzAwmc1FEoWVUrAC",
	"nCeqRZZa6r2e52yo5rC5zTvLWFoxUGCV69",
	"2NzHpU4ogdH9F91YaZ2RkGAs3tCK29Fm6V9",
	"qgnm7MhxME5oAj5AXpfzse9sREQ6qshf2u",
	"2G1X2Atfyc3AQNyp3UszU6uHNLs2uSTL8XF",
	"kdgzfJAP5A7A77R8m4QEpi2gNjSQzRxZZP",
	"2NhZJaCErM4mqUtpcu9aDcUgyVA6EzCM9pm",
	"EEhqTbQP8fpPBZij5WdBRhzRZBzWnhCQsZ",
	"W6UhJheYvExzwYNFeyRT6Sm2fsZxfNSCFF",
	"SDhfRZ8twZ4Ybtctuaf3Qr9yE5A49a1zJp",
	"NTTCre1H4SznbfKUsuvEtEctoqcsxdzjuJ",
	"E5P4iMWjvEUNQn9uYinzp3fFDJHtC7kPLq",
	"rDV5K5q8j9VEW4JRu34c5erT2PMEFtwTH9",
	"8qGtRTpYQbAzfJc4FcaWbuwaku1LnkPPbD",
	"NcKq9ipiePzgZBKbaXvRmbc9cJEvKhVqb5",
	"ucWmmA91iaevnd3WkHgFExm9qzTV4mtLpg",
	"2BuEM1ruAFp3nYVcwviT5NDvY8UdHUgZhWG",
	"2kk8BAJtTTU6aNvT1Lw31rmbCbMU3Xnmp4h",
	"puY8tYhFtU8yJU9yJ7s8cJDpzhKgfoN2mg",
	"R43teMKjkYmc433SYeo1btLmwt7DZ4H8Hi",
	"2H9tdF3uF1mAbRREvVHjjZSN1v2hN54t384",
	"2Z6Ty4GfgQviuzrJMmNZu3nb2TVDMRKsXgb",
	"sneUfLKiWZQReydqbYp4EFTxggCbQZk9gi",
	"zGh2xdLfEKLxPYkRdSoQckJZjkNGxMPL39",
	"HZXxJyc4iR1UV22Tm4JLBRfS5gFexzMJfD",
	"2U2o4ncEgyBw321gGL7aeExQ9vmVNdPanPL",
	"2CjUzn3qb89Z14VWwH7GqstWJ9bW7GfLm1v",
	"2fU9fYdTidLDzK8brNT3tKxKmMEpcXsyY6v",
}

func InitTransaction() coin.Transaction {

	genesis_output := "33dc904e0e697509b216d13cafbae49cc4b3da7073bab19e85a6d428645f429e"
	genesis_sig := "9244e5a6bb7db77a796951324110439d04f1770a28ae22c92f22c9fb23c522fb0498a384e907524764fa7670981588e5994edf810d6db8feda487ff30183cad401" //sig for spending genesis output

	var tx coin.Transaction

	output := cipher.MustSHA256FromHex(genesis_output)
	tx.PushInput(output)

	for i := 0; i < 100; i++ {
		addr := cipher.MustDecodeBase58Address(AddrList[i])
		tx.PushOutput(addr, 3e12, 1) // 10e6*10e6
	}

	txs := make([]cipher.Sig, 1)
	sig := genesis_sig
	txs[0] = cipher.MustSigFromHex(sig)
	tx.Sigs = txs

	tx.UpdateHeader()

	err := tx.Verify()

	if err != nil {
		log.Panic(err)
	}

	//log.Printf("signature= %s", tx.Sigs[0].Hex())
	return tx
}
