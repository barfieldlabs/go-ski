# Go-Ski 🎿

<img src="https://github.com/barfieldlabs/go-ski/assets/73442540/54fc642f-63a1-4f81-afa7-1754822b436b" width="200" height="200">

---

## Overview 📖

Go-Ski is a web scraping library built in Go. It aims to provide a simple yet powerful way to perform various web scraping tasks. The library is built on top of the Chrome DevTools Protocol, using the `chromedp` package for the heavy lifting.

## Features 👀

- Clicking on elements
- Submitting forms
- Navigating through iframes
- Switching tabs
- Extensible for more complex tasks

## Installation ⚡️

To install Go-Ski, run the following command:

```bash
go get github.com/barfieldlabs/go-ski
```

## Usage 🚀

Here's a simple example that demonstrates how to perform a click and form submission:

```go
package main

import (
"context"
"log"

    "github.com/barfieldlabs/go-ski/core"
    "github.com/chromedp/cdproto/target"

)

func main() {
proc := core.NewProcedures()
ctx := context.Background()
var initialTargets []\*target.Info

    proc.Actions = []core.Action{
    	{
    		Type:  core.Click,
    		XPath: "some_xpath",
    	},
    	{
    		Type:       core.FormSubmit,
    		XPath:      "form_xpath",
    		FormFields: map[string]string{"username": "user", "password": "pass"},
    		SubmitBtn:  "submit_btn_xpath",
    	},
    }

    err := proc.Execute(ctx, initialTargets)
    if err != nil {
    	log.Fatalf("Failed to perform actions: %v", err)
    }

    log.Println("Successfully completed web scraping.")

}
```

## Contributing 🫱🏾‍🫲🏽

Feel free to open issues or submit pull requests. All contributions are welcome!

## License

This project is open-source and available under the MIT License.
