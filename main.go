package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"stats.k8s.devminer.xyz/internal/cache"
	"stats.k8s.devminer.xyz/internal/scraper"
	"stats.k8s.devminer.xyz/internal/web"
	"stats.k8s.devminer.xyz/pkg/prometheus"
	"time"
)

//go:embed "static"
var Static embed.FS

var (
	prometheusUrl = flag.String("prometheus", "http://localhost:9090", "Prometheus URL")
	listenAddr    = flag.String("listen-addr", ":3000", "Address to listen on")
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	flag.Parse()

	cacheObj := cache.New()

	promClient, err := prometheus.FromUrl(*prometheusUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create prometheus client")
	}

	scraperObj := scraper.New(promClient)

	app := fiber.New()
	app.Use(cors.New(cors.Config{AllowOrigins: "*"}))

	app.Get("/api/stats", func(c *fiber.Ctx) error {
		stats := cacheObj.Get()
		stats["cluster"] = calculateClusterStats(stats)

		return c.JSON(stats)
	})

	app.Get("/", func(c *fiber.Ctx) error {
		stats := cacheObj.Get()
		cluster := calculateClusterStats(stats)

		c.Set("Content-Type", fiber.MIMETextHTMLCharsetUTF8)
		return web.Index(cluster, stats).Render(c.UserContext(), c)
	})

	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(Static),
		PathPrefix:   "/static",
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		log.Info().Msg("starting HTTP server...")
		if err := app.Listen(*listenAddr); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start api")
		}
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Second)

		log.Info().Msg("starting scraper...")
		if err := cacheObj.Update(scraperObj.Scrape); err != nil {
			log.Error().Err(err).Msg("failed to update cache")
		}

		for {
			select {
			case <-ticker.C:
				if err := cacheObj.Update(scraperObj.Scrape); err != nil {
					log.Error().Err(err).Msg("failed to update cache")
				}

			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	cancel()

	log.Info().Msg("shutting down HTTP server...")
	if err := app.Shutdown(); err != nil {
		log.Fatal().Err(err).Msg("failed to shutdown HTTP server")
	}

	log.Info().Msg("exiting")
}

func calculateClusterStats(stats map[string]cache.Node) cache.Node {
	cluster := cache.Node{}

	for _, node := range stats {
		cluster.CPU = add(cluster.CPU, node.CPU)
		cluster.Memory = add(cluster.Memory, node.Memory)
		cluster.NetworkReceive = add(cluster.NetworkReceive, node.NetworkReceive)
		cluster.NetworkTransmit = add(cluster.NetworkTransmit, node.NetworkTransmit)
	}

	return cluster
}

func add(a cache.Value, b cache.Value) cache.Value {
	a.V1 += b.V1
	a.V2 += b.V2

	return a
}
