package lbdl

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/anacrolix/torrent"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type xfer struct {
	torrent *torrent.Torrent
	gauge   *widgets.Gauge
}

var (
	client      *torrent.Client // torrent client
	downloadDir string          // directory to store downloads
	magnetList  string          // filename to read magnet links from
	pwd         string          // working directory of application
	sigc        chan os.Signal  // channel for listening to syscalls
	torrentDir  string          // directory to read torrent files from
	transfers   []*xfer         // list of transfers
)

func init() {
	sigc = make(chan os.Signal, 1)
	transfers = make([]*xfer, 0)
}

// Start the lbdl process
func Start() (err error) {
	var pwd string
	pwd, err = getCmdDir()
	if nil != err {
		return
	}

	defaultDownloadDir := fmt.Sprintf("%s/downloads", pwd)
	defaultMagnetList := fmt.Sprintf("%s/magnet.list", pwd)
	defaultTorrentDir := fmt.Sprintf("%s/torrents", pwd)

	flag.StringVar(&downloadDir, "d", defaultDownloadDir, "torrent download directory")
	flag.StringVar(&magnetList, "m", defaultMagnetList, "magnet list file")
	flag.StringVar(&torrentDir, "t", defaultTorrentDir, "torrent file directory")
	flag.Parse()

	log.Print("lbdl starting")

	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = downloadDir
	client, err = torrent.NewClient(cfg)
	if nil != err {
		return
	}
	defer client.Close()

	// read from the ./torrents/ dir
	var tfs []string
	tfs, err = getTorrentFiles()
	if nil != err {
		return
	}
	for _, tf := range tfs {
		t, err := client.AddTorrentFromFile(tf)
		if nil != err {
			log.Print(err.Error())
		} else {
			<-t.GotInfo()
			t.DownloadAll()
			log.Printf("downloading torrent from file %s", tf)
			x := &xfer{torrent: t}
			transfers = append(transfers, x)
		}
	}
	// read from ./magnet.list
	ms, err := getMagnetLinks()
	if err != nil {
		return err
	}
	for _, m := range ms {
		t, err := client.AddMagnet(m)
		if nil != err {
			log.Print(err.Error())
		} else {
			<-t.GotInfo()
			t.DownloadAll()
			log.Printf("downloading torrent from magnet link %s", m)
			x := &xfer{torrent: t}
			transfers = append(transfers, x)
		}
	}

	// panic if no torrents found
	if 0 == len(transfers) {
		return errors.New("no torrent files or magnet links found")
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p := widgets.NewParagraph()
	p.Text = "Welcome to lbdl - press Q to quit"
	p.SetRect(0, 0, 80, 3)
	ui.Render(p)

	// listen for exit signals
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		var close bool
		s := <-sigc
		switch s {
		case syscall.SIGHUP:
			close = true
		case syscall.SIGINT:
			close = true
		case syscall.SIGTERM:
			close = true
		case syscall.SIGQUIT:
			close = true
		}
		if close {
			log.Print("shutting down")
			client.Close()
			ui.Close()
			os.Exit(0)
		}
	}()

	for i, t := range transfers {
		g := widgets.NewGauge()
		g.Title = fmt.Sprintf(" %s ", t.torrent.Name())
		g.SetRect(0, 3+3*i, 80, 6+3*i)
		g.Percent = 0
		g.Label = fmt.Sprintf("%v%% complete", g.Percent)
		g.BarColor = ui.ColorYellow
		g.BorderStyle.Fg = ui.ColorWhite
		transfers[i].gauge = g
		ui.Render(g)
	}

	go func() {
		if ok := client.WaitAll(); ok {
			log.Print("all torrents complete")
		} else {
			log.Print("client stopped")
		}
	}()

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				sigc <- syscall.SIGQUIT
			}
		case <-ticker:
			updateGauges()
		}
	}
}

func getCmdDir() (dir string, err error) {
	ex, err := os.Executable()
	if err != nil {
		return
	}
	dir = filepath.Dir(ex)
	return
}

func getMagnetLinks() (links []string, err error) {
	links = make([]string, 0)

	file, err := os.Open(magnetList)
	if err != nil {
		return
	}
	defer file.Close()

	log.Printf("reading magnet links from %s", file.Name())

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		links = append(links, scanner.Text())
	}

	err = scanner.Err()

	return
}

func getTorrentFiles() (tf []string, err error) {
	tf = make([]string, 0)

	log.Printf("reading torrent files from %s", torrentDir)

	files, err := ioutil.ReadDir(torrentDir)
	if err != nil {
		return
	}

	for _, file := range files {
		n := file.Name()
		if 8 >= len(n) {
			continue
		}
		if `.torrent` != n[len(n)-8:] {
			continue
		}
		tf = append(tf, fmt.Sprintf("%s/%s", torrentDir, n))
	}

	return
}

func updateGauges() {
	for _, t := range transfers {
		done := t.torrent.BytesCompleted()
		total := done + t.torrent.BytesMissing()
		pct := int(float64(done) / float64(total) * 100)
		t.gauge.Percent = pct
		t.gauge.Label = fmt.Sprintf("%v%% complete", pct)
		if 100 == pct {
			t.gauge.BarColor = ui.ColorGreen
		}
		ui.Render(t.gauge)
	}
}
