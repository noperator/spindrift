package main

import (
	// "encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	// "os/user"
	// "path/filepath"
	// "strings"
	// "time"

	"github.com/mxschmitt/playwright-go"
	"github.com/pelletier/go-toml"
)

func assertErrorToNilf(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

// Allows us to use the Page type in a function's receiver.
type Page struct {
	playwright.Page
}

type Config struct {
	Spot     string
	Location string
}

func (page *Page) weather(location string) error {
	// Load weather forecast page.
	if _, err := page.Goto(fmt.Sprintf("https://www.google.com/search?q=weather forecast %s", location)); err != nil {
		return fmt.Errorf("could not load forecast page: %w", err)
	}

	// Get width of day in forecast.
	// dayWidth, err := page.EvalOnSelector(".wob_df", "e => e.offsetWidth")
	// if err != nil {
	// 	return fmt.Errorf("could not get width: %w", err)
	// }

	// Resize forecast.
	if _, err := page.EvalOnSelector("#wob_dp", "e => e.setAttribute('style', 'width: calc(75px * 7); height: 72px')"); err != nil {
		return fmt.Errorf("could not resize forecast: %w", err)
	}

	// Remove margin from each day in forecast.
	if _, err := page.EvalOnSelectorAll(".wob_df", "days => days.map(day => day.setAttribute('style', 'margin: 0px'))"); err != nil {
		return fmt.Errorf("could not remove margin from forecast day: %w", err)
	}

	// Hide last day of the week (only show 7 days).
	if _, err := page.EvalOnSelector("div.wob_df:nth-child(8)", "e => e.setAttribute('style', 'display: none')"); err != nil {
		return fmt.Errorf("could not hide last forecast day: %w", err)
	}

	// Remove highlight from each highlighted day in forecast.
	if _, err := page.EvalOnSelectorAll(".wob_ds", "hils => hils.map(hil => hil.classList.remove('wob_ds'))"); err != nil {
		return fmt.Errorf("could not remove highlight from forecast day: %w", err)
	}

	// Remove day of week from each day in forecast.
	if _, err := page.EvalOnSelectorAll(".wob_df > div:nth-child(1)", "dows => dows.map(dow => dow.setAttribute('style', 'display: none'))"); err != nil {
		return fmt.Errorf("could not hide day of week: %w", err)
	}

	for element, selector := range map[string]string{
		"temp":    "#wob_gsp",
		"rain":    "#wob_gsp",
		"wind":    "#wob_gsp",
		"weather": "#wob_dp",
	} {
		log.Printf("taking a screenshit of %s\n", element)
		if element != "weather" {
			if err := page.Click(fmt.Sprintf("#wob_%s", element)); err != nil {
				return fmt.Errorf("could not click forecast element: %w", err)
			}
		}
		elem, err := page.QuerySelector(selector)
		if err != nil {
			return fmt.Errorf("could not query forecast element: %w", err)
		}
		if _, err = elem.Screenshot(playwright.ElementHandleScreenshotOptions{
			Path: playwright.String(fmt.Sprintf("img/%s.png", element)),
		}); err != nil {
			return fmt.Errorf("could not take screenshit: %w", err)
		}
	}

	return nil
}

func main() {

	headful := flag.Bool("f", false, "headful mode")
	flag.Parse()

	// Load config.
	configFile, err := ioutil.ReadFile("config.toml")
	assertErrorToNilf("could not read config: %v", err)
	config := Config{}
	err = toml.Unmarshal(configFile, &config)
	assertErrorToNilf("could not load config: %v", err)
	// fmt.Println(config)

	// Get user info for home directory later on.
	// usr, err := user.Current()
	// assertErrorToNilf("could not get current user: %v", err)

	// Load Playwright, browser, and page.
	log.Println("setting up")
	pw, err := playwright.Run()
	assertErrorToNilf("could not start playwright: %v", err)
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(!*headful),
		// DownloadsPath: playwright.String(filepath.Join(usr.HomeDir, "Downloads")),
	})
	assertErrorToNilf("could not launch browser: %v", err)
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		AcceptDownloads:   playwright.Bool(true),
		DeviceScaleFactor: playwright.Float(2),
		Viewport: &playwright.BrowserNewContextViewport{
			Width:  playwright.Int(1440),
			Height: playwright.Int(900),
		},
	})

	// Gracefully teardown upon interrupt.
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		log.Printf("Interrupted; tearing down")
		assertErrorToNilf("could not close browser: %v", browser.Close())
		assertErrorToNilf("could not stop Playwright: %v", pw.Stop())
		os.Exit(1)
	}()

	newPage, err := context.NewPage()
	assertErrorToNilf("could not create page: %v", err)
	page := Page{newPage}
	assertErrorToNilf("could not check weather: %v", page.weather(config.Location))

	// time.Sleep(10000 * time.Second)

	log.Println("tearing down")
	assertErrorToNilf("could not close browser: %v", browser.Close())
	assertErrorToNilf("could not stop Playwright: %v", pw.Stop())
}
