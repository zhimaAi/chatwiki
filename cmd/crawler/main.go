// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package main

import (
	"errors"
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	"github.com/playwright-community/playwright-go"
	"log"
	"net/http"
	netURL "net/url"
	"sync"
	"time"
)

var (
	cfg                   config
	browser               playwright.Browser
	browserActive         bool
	pw                    *playwright.Playwright
	browserMu             sync.Mutex
	concurrent            chan struct{}
	browserLastActiveTime time.Time
	idleTimeout           = 3 * time.Minute // close browser after 3 minutes of inactivity to release memory
)

var TooManyRequestsError = errors.New("too many requests, please try again later")

type config struct {
	IsProduction bool `env:"PRODUCTION" envDefault:"true"`
	Concurrent   int  `env:"CONCURRENT" envDefault:"5"`
}

func init() {
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	if cfg.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	concurrent = make(chan struct{}, cfg.Concurrent)

	// install playwright
	err := playwright.Install()
	if err != nil {
		panic(err)
	}
}

func main() {

	// start timer to check and restart browser
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			delayCloseBrowser()
		}
	}()

	// http server
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.POST("/content", handleContentRequest)
	s := &http.Server{
		Addr:           ":3800",
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Starting server on port 3800")
	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// handleContentRequest handles the content request
func handleContentRequest(c *gin.Context) {
	// check params
	var request struct {
		URL string `json:"url"`
	}
	if err := c.BindJSON(&request); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid params"})
		return
	}

	// parse url
	parsedURL, err := netURL.Parse(request.URL)
	if err != nil || parsedURL == nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	// open browser
	err = openBrowser()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// get page content
	pageInfo, err := fetchURLContent(parsedURL)
	if err != nil {
		log.Println(err)
		if errors.Is(err, TooManyRequestsError) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//c.String(http.StatusOK, pageInfo.RawHtml)
	c.JSON(http.StatusOK, pageInfo)
	return
}

// openBrowser opens a browser instance using Playwright
func openBrowser() error {
	browserMu.Lock()
	defer browserMu.Unlock()
	browserLastActiveTime = time.Now()

	if !browserActive {
		var err error
		pw, err = playwright.Run()
		if err != nil {
			return fmt.Errorf("could not start playwright: %v", err)
		}

		headless := true
		browser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless})
		if err != nil {
			return fmt.Errorf("could not launch browser: %v", err)
		}
		browserActive = true
	}
	return nil
}

// delayCloseBrowser closes the browser after idleTimeout
func delayCloseBrowser() {
	browserMu.Lock()
	defer browserMu.Unlock()

	if time.Since(browserLastActiveTime) > idleTimeout {
		if err := closeBrowser(); err != nil {
			panic(err)
		}
	}
}

// closeBrowser closes the browser and stops Playwright
func closeBrowser() error {
	if browserActive {
		if err := browser.Close(); err != nil {
			return fmt.Errorf("could not close browser: %v", err)
		}
		if err := pw.Stop(); err != nil {
			return fmt.Errorf("could not stop playwright: %v", err)
		}
		browserActive = false
	}

	return nil
}
