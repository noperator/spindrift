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

func (page *Page) screenshit(elem string, sel string) error {
	log.Printf("taking a screenshit: %s\n", elem)
	elemHand, err := page.QuerySelector(sel)
	if err != nil {
		return fmt.Errorf("could not query selector: %w", err)
	}
	if _, err = elemHand.Screenshot(playwright.ElementHandleScreenshotOptions{
		Path: playwright.String(fmt.Sprintf("img/%s.png", elem)),
	}); err != nil {
		return fmt.Errorf("could not take screenshit: %w", err)
	}
	return nil
}

func (page *Page) checkWeather(location string) error {
	// Load weather forecast page.
	if _, err := page.Goto(fmt.Sprintf("https://www.google.com/search?q=weather forecast %s", location)); err != nil {
		return fmt.Errorf("could not load weather forecast page: %w", err)
	}

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

	for elem, sel := range map[string]string{
		"temp":    "#wob_gsp",
		"rain":    "#wob_gsp",
		"wind":    "#wob_gsp",
		"weather": "#wob_dp",
	} {
		if elem != "weather" {
			if err := page.Click(fmt.Sprintf("#wob_%s", elem)); err != nil {
				return fmt.Errorf("could not click element: %w", err)
			}
		}
		if err := page.screenshit(elem, sel); err != nil {
			return fmt.Errorf("could not screenshit %s: %w", elem, err)
		}
	}

	return nil
}

func (page *Page) checkSurf(spot string) error {
	// Load surf forecast page.
	if _, err := page.Goto(fmt.Sprintf("https://magicseaweed.com/%s", spot), playwright.PageGotoOptions{
		WaitUntil: playwright.String("networkidle"),
	}); err != nil {
		return fmt.Errorf("could not load surf forecast page: %w", err)
	}

	// Remove padding from surf graph.
	if _, err := page.EvalOnSelectorAll(".scrubber-bars-container", "bars => bars.map(bar => bar.setAttribute('style', 'padding-top: 0px'))"); err != nil {
		return fmt.Errorf("could not remove surf padding: %w", err)
	}

	// Remove header from wind graph.
	if _, err := page.EvalOnSelectorAll(".scrubber-graph-header:below(span:text(\"Wind\"))", "heads => heads.map(head => head.setAttribute('style', 'display: none'))"); err != nil {
		return fmt.Errorf("could not remove wind headers: %w", err)
	}

	for elem, sel := range map[string]string{
		"current": "div:below(span:text(\"Current Surf Report\"))",
		"surf":    "#tab-7day .scrubber-forecast-graph-container",
	} {
		if err := page.screenshit(elem, sel); err != nil {
			return fmt.Errorf("could not screenshit %s: %w", elem, err)
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
		UserAgent:         playwright.String("User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36"),
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
	log.Println("checking weather")
	assertErrorToNilf("could not check weather: %v", page.checkWeather(config.Location))
	log.Println("checking surf")
	assertErrorToNilf("could not check surf: %v", page.checkSurf(config.Spot))

	// time.Sleep(10000 * time.Second)

	log.Println("tearing down")
	assertErrorToNilf("could not close browser: %v", browser.Close())
	assertErrorToNilf("could not stop Playwright: %v", pw.Stop())
}
