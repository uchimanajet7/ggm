package main

import (
	"fmt"
	"log"
	"time"

	"github.com/uchimanajet7/ggm/utils"
)

func main() {
	// Get user config
	config, err := utils.GetUserConfig()
	if err != nil {
		log.Fatalf("+%v\n", err)
	}

	// usb power off
	err = config.RunUsbCommand(false)
	if err != nil {
		log.Fatalf("+%v\n", err)
	}

	nowEpoch := utils.GetNowEpoch()
	messages, err := utils.GetGmailData(config.LastDate)
	if err != nil {
		log.Fatalf("+%v\n", err)
	}

	usbPower := false
	for _, v := range messages {
		// Is target data?
		if !config.IsTargetData(v) {
			continue
		}

		if !usbPower {
			// usb power on
			err = config.RunUsbCommand(true)
			if err != nil {
				log.Fatalf("+%v\n", err)
			}
			usbPower = true
			time.Sleep(1 * time.Second)
		}

		// get speak text
		text := v.GetSpeakText()
		fmt.Printf("%+v\n\n", text)

		// speak text
		err := config.RunSpeakCommand(text)
		if err != nil {
			log.Fatalf("+%v\n", err)
		}
	}

	// update user config
	config.UpdateUserConfig(nowEpoch)

	fmt.Print("\nAll execution completed normally.\n\n")
}
