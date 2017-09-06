# ggm
A tool created to receive Gmail with Raspberry Pi 3 using  golang Gmail SDK.

## Description
I made it for Raspberry Pi 3 to work.

Check Gmail using golang Gmail SDK and execute the specified command when there is a mail that meets the conditions.

Speech synthesis(text-to-speech) and USB power management commands must be installed in Raspberry Pi 3 beforehand as specified commands.

## Demo

<blockquote class="embedly-card"><h4><a href="https://twitter.com/uchimanajet7/status/905570296039481344">Manabu Uchida on Twitter</a></h4><p>音声合成できたー楽しい https://t.co/ZRDM7dXPoG</p></blockquote>
<script async src="//cdn.embedly.com/widgets/platform.js" charset="UTF-8"></script>

## Features
- It is made by golang so it supports multi form.
- If you get Gamil certification and describe the setting in the setting file, it works.

## Requirement
- Go 1.9+
- Speech synthesis(text-to-speech) tool 
	- AquesTalk Pi - Raspberry Pi用の音声合成
		- https://www.a-quest.com/products/aquestalkpi.html
- USB power management tool
	- codazoda/hub-ctrl.c: Control USB power on a port by port basis on some USB hubs. 
		- https://github.com/codazoda/hub-ctrl.c
- Gamil certification
	- Go Quickstart  |  Gmail API  |  Google Developers 
		- https://developers.google.com/gmail/api/quickstart/go


## Usage
Just run the only one command.

```	sh
$ ./ggm
```

However, setting is necessary to execute.

### Setting Example

1. In the same place as the binary file create `.ggm` dir.

1. Get Gamil certification.
	- Go Quickstart  |  Gmail API  |  Google Developers 
		- https://developers.google.com/gmail/api/quickstart/go
1. Save Gamil certification `.ggm/client_secret.json`

```json
{
	"installed": {
		"client_id": "123456789012-abcdefghijklnmopqrstuvwxyz123456.apps.googleusercontent.com",
		"project_id": "ex-gmail-test",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
        "token_uri": "https://accounts.google.com/o/oauth2/token",
        "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
        "client_secret": "AbCdEfGhijklNmopqr135790",
        "redirect_uris": [
            "urn:ietf:wg:oauth:2.0:oob",
            "http://localhost"
        ]
	}
}
```
	
1. Run command `./ggm`
1. Copy the URL displayed on the terminal and access it with the browser.
1. Perform an authorization operation in the browser, display the code, and input it to the terminal.
1. `.ggm/client_token.json` file is created.

```json
{
	"access_token": "1234.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"token_type": "Bearer",
	"refresh_token": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"expiry": "2017-09-04T20:25:28.251300122+09:00"
}
```

1. Execution settings are done with `.ggm/user_config.json` file.

```sh
{
	"LastDate": 1504243346130,
	"LastTotal": 749,
	"UserEmail": "test@example.com",
	"SpeakCommands": [["/home/pi/aquestalkpi/AquesTalkPi", "-s", "120", "%s"],["aplay"]],
	"UsbCommands": [["sudo", "hub-ctrl", "-h", "0", "-P", "2", "-p", "%d"]],
	"Filters": [
		{
			"From": "user1@example.com",
			"Subjects": ["test","user1"]
		},
		{
			"From": "user2@example.com",
			"Subjects": null
		}
	]
}
```


## Installation

```	sh
$ go get github.com/uchimanajet7/ggm
```

## Author
[uchimanajet7](https://github.com/uchimanajet7)


## Licence
[MIT License](https://github.com/uchimanajet7/ggm/blob/master/LICENSE)
