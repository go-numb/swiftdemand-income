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

	log.Infof("%+v", config)
}

func main() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// 初回起動
	err := Income()
	if err != nil {
		goto EXIT
	}

	for {
		select {
		case <-ticker.C:
			err = Income()
			if err != nil {
				log.Error(err)
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
				"--no-sandbox",
				// User-Agentがないとheadless modeでjavascriptを起動できない
				`--user-agent="Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.75 Safari/537.36"`,
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

	if err := page.FindByName("commit").Submit(); err != nil {
		if err := page.FindByXPath(`//*[@id="account"]/div[2]/div[1]/section/div/div[3]/input`).Click(); err != nil {
			return err
		}
	}
	time.Sleep(1 * time.Second)

	return nil
}
