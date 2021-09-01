package main

import (
	"encoding/json"
	"flag"
	"github.com/OpenDiablo2/AbyssEngine/providers/renderprovider/raylibrenderprovider"
	"io/ioutil"
	"os"
	"path"

	"github.com/OpenDiablo2/AbyssEngine/engine"
	"github.com/pkg/profile"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var runPath string
var doProfile bool

func initFlags() {
	flag.StringVar(&runPath, "path", "", "path to the engine runtime files")
	flag.BoolVar(&doProfile, "profile", false, "profile the engine")
	flag.Parse()

	if runPath == "" {
		runPath, _ = os.Getwd()
	}
}

func initLogging() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func main() {
	initFlags()
	initLogging()

	if doProfile {
		defer profile.Start(profile.ProfilePath(".")).Stop()
	}

	renderProvider := raylibrenderprovider.New()

	log.Info().Msg("Abyss Engine")
	log.Debug().Msgf("Runtime Path: %s", runPath)

	renderProvider.SetLoggerCallback(func(logLevel int, s string) {
		[]func() *zerolog.Event{
			log.Trace,
			log.Debug,
			log.Info,
			log.Warn,
			log.Error,
			log.Fatal,
		}[logLevel-1]().Msg(s)
	})
	renderProvider.SetLoggerLevel(0)

	engineConfig := engine.Configuration{
		RootPath: runPath,
	}

	jsonFile, err := os.Open(path.Join(runPath, "config.json"))
	if err == nil {
		bytes, _ := ioutil.ReadAll(jsonFile)
		_ = json.Unmarshal(bytes, &engineConfig)
		_ = jsonFile.Close()
	}

	coreEngine := engine.New(engineConfig, renderProvider)

	coreEngine.Run()
	coreEngine.Destroy()

}
