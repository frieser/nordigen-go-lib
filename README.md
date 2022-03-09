Nordigen Golang API client library
==================================

[Norgigen API documention](https://nordigen.com/en/account_information_documenation/api-documention/overview/)

How to use it([See a real example app using this library](https://github.com/frieser/openbanking-cli)):

```go
package main

import (
	"github.com/frieser/nordigen-go-lib/v2"
	"github.com/google/uuid"
	"strconv"
	"time"
	"log"
)

const redirectPort = ":3000"
	
func main() {
	c, err := nordigen.NewClient("secret_id", "secret_key")

	if err != nil {
        log.Fatal(err)
	}

	// supported banks in a country
	countryBanks, err := c.ListInstitutions(countryCode)
	
	// get authorization
	endUserId := uuid.NewString()
	// look and the function below
	r, err := GetAuthorization(c, bankId, endUserId)

	// get account metadata, details and balance
	mtdt, err := c.GetAccountMetadata(r.Accounts[0])
	
	dtls, err := c.GetAccountDetails(r.Accounts[0])

	blnc, err := c.GetAccountBalances(r.Accounts[0])
	
	//get accounts transactions
	txns, err := c.GetAccountTransactions(r.Accounts[0].Id, nil, nil)

	//get previous week accounts transactions
	from := &time.Now().Add(-14 * 24 * time.Hour)
	to := &time.Now().Add(-7 * 24 * time.Hour)
	txns, err := c.GetAccountTransactions(r.Accounts[0].Id, from, to)
}

func GetAuthorization(cli nordigen.Client, bankId string, endUserId string) (nordigen.Requisition, error) {
	requisition := nordigen.Requisition{
		Redirect:  "http://localhost" + redirectPort,
		Reference: strconv.Itoa(int(time.Now().Unix())),
		Agreement: "",
	}
	r, err := cli.CreateRequisition(requisition)

	if err != nil {
		return nordigen.Requisition{}, err
	}
	go internal.OpenBrowser(r.Redirect)

	ch := make(chan bool, 1)

	go internal.CatchRedirect(redirectPort, ch)

	<-ch

	for r.Status == "CR" {
		r, err = cli.GetRequisition(r.Id)

		if err != nil {

			return nordigen.Requisition{}, err
		}
		time.Sleep(1 * time.Second)
	}

	return r, nil
}

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func CatchRedirect(port string, ch chan bool) {
	handler := func(chan bool) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ch <- true
			w.Write([]byte("You can close this window now"))
		})
	}
	http.Handle("/", handler(ch))

	err := http.ListenAndServe(port, nil)

	if err != nil {
		panic(err)
	}
}

```