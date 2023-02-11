package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/playwright-community/playwright-go"
)

func numBetween(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func main() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}

	userDataDir := getEnv("MS_REWARDS_BOT_USER_DIR", os.TempDir())

	browser, err := pw.Chromium.LaunchPersistentContext(
		userDataDir,
		playwright.BrowserTypeLaunchPersistentContextOptions{
			Headless: playwright.Bool(false),
			Channel:  playwright.String("msedge"),
		},
	)
	if err != nil {
		log.Fatalf("Couldn't launch browser")
	}

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("couldn't launch page")
	}

	if _, err = page.Goto("https://rewards.bing.com"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	time.Sleep(3 * time.Second)

	fmt.Println("Earning activity points now...")
	activities, err := page.QuerySelectorAll(".mee-icon-AddMedium")
	if err != nil {
		log.Fatalf("could not get activities: %v", err)
	}
	for _, activity := range activities {
		time.Sleep(3 * time.Second)
		activity.Click()
	}

	fmt.Println("Earning search points now...")
	if err = page.BringToFront(); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	if _, err = page.Goto("https://bing.com"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	time.Sleep(3 * time.Second)

	for i := 0; i < 33; i++ {
		locator, err := page.Locator("#sb_form_q")
		if err != nil {
			log.Fatal()
		}
		sentence := gofakeit.Sentence(numBetween(12, 44))
		locator.Fill(sentence)
		fmt.Printf("Search #%d: %s\n", i, sentence)
		if err = page.Keyboard().Press("Enter"); err != nil {
			log.Fatalf("could not press key: %v", err)
		}
		time.Sleep(time.Duration(numBetween(3000, 6000)) * time.Millisecond)
	}

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
