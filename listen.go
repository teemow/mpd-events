package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fhs/gompd/mpd"
	"github.com/hoisie/web"
	"github.com/spf13/cobra"
)

type ListenFlags struct {
	Host    string
	Port    string
	OnPlay  string
	OnPause string
}

var (
	listenCmd = &cobra.Command{
		Use:   "listen",
		Short: "Listen to mpd state changes.",
		Long:  "Listen to mpd state changes.",
		Run:   listenRun,
	}

	listenFlags = &ListenFlags{}
)

func init() {
	listenCmd.PersistentFlags().StringVar(&listenFlags.Host, "host", "localhost", "mpd host")
	listenCmd.PersistentFlags().StringVar(&listenFlags.Port, "port", "6600", "mpd port")
	listenCmd.PersistentFlags().StringVar(&listenFlags.OnPlay, "onplay", "", "Execute if state changes to play")
	listenCmd.PersistentFlags().StringVar(&listenFlags.OnPause, "onpause", "", "Execute if state changes to stop/pause")
}

func execCmd(cmdString string) (string, error) {
	args := strings.Split(cmdString, " ")

	cmd := exec.Command(args[0], args[1:]...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s %s - %v", stdout.String(), stderr.String(), err)
	}

	return stdout.String(), nil
}

func onStateChange(state string) {
	if state == "play" && listenFlags.OnPlay != "" {
		out, err := execCmd(listenFlags.OnPlay)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Print(out)
	}
	if (state == "stop" || state == "pause") && listenFlags.OnPause != "" {
		out, err := execCmd(listenFlags.OnPause)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Print(out)
	}
}

func watchMPDState(mpc *mpd.Client) {
	// Loop printing the current status of MPD.
	oldstate := ""
	for {
		status, err := mpc.Status()
		if err != nil {
			log.Println(err)
		}
		if status["state"] != oldstate {
			onStateChange(status["state"])
			oldstate = status["state"]
		}
		time.Sleep(1e9)
	}
}

func manage(mpc *mpd.Client) func(string) string {
	return func(key string) string {
		status, err := mpc.Status()
		if err != nil {
			log.Println(err)
		}

		switch {
		case key == "left":
			mpc.Previous()
		case key == "right":
			mpc.Next()
		case key == "select":
			if status["state"] == "play" {
				mpc.Pause(true)
			} else {
				mpc.Pause(false)
			}
		case key == "up":
			execCmd("amixer set PCM 2dB+")
		case key == "down":
			execCmd("amixer set PCM 2dB-")
		}

		return key
	}
}

func listenRun(cmd *cobra.Command, args []string) {
	// Connect to MPD server
	mpc, err := mpd.Dial("tcp", fmt.Sprintf("%s:%s", listenFlags.Host, listenFlags.Port))
	if err != nil {
		log.Fatalln(err)
	}
	defer mpc.Close()

	go watchMPDState(mpc)

	web.Get("/(.*)", manage(mpc))
	go web.Run("0.0.0.0:80")

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	web.Close()
}
