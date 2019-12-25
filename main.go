package main

import (
	"flag"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/sclevine/agouti"
)

var config = &Config{}

type Config struct {
	ID, Password string
}

func init() {
	var id, pass string
	flag.StringVar(&id, "id", "null", "option -id is login ID")
	flag.StringVar(&pass, "pass", "null", "option -pass is login Password")
	flag.Parse()

	config.ID = id
	config.Password = pass
}

func main() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	var err error

	for {
		select {
		case <-ticker.C:
			err = Income()
			if err != nil {
				goto EXIT
			}
		}
	}

EXIT:

	log.Fatal(err)
}

// Income gets crypto by swiftdemand
func Income() error {
	d := agouti.ChromeDriver(
		agouti.ChromeOptions(
			"args", []string{
				"--headless", // headlessモードの指定
				"--disable-gpu",
				"no-sandbox",
				"--window-size=1280,800", // ウィンドウサイズの指定
			}),
		agouti.Debug,
	)
	d.Start()
	defer d.Stop()

	page, err := d.NewPage()
	if err != nil {
		return err
	}
	defer page.Destroy()

	if err := page.Navigate("https://www.swiftdemand.com/users/sign_in"); err != nil {
		return err
	}

	if err := page.FindByName("user[email]").Fill(config.ID); err != nil {
		return err
	}
	if err := page.FindByName("user[password]").Fill(config.Password); err != nil {
		return err
	}
	if err := page.FindByID("new_user").Submit(); err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	if err := page.FindByLink("Claim").Click(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)

	return nil
}
