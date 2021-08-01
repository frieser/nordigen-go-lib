Nordigen Golang API client library
==================================

[Norgigen API documention](https://nordigen.com/en/account_information_documenation/api-documention/overview/)

How to use it:

```go
package main

import (
	"github.com/frieser/nordigen-go-lib"
	"github.com/google/uuid"
	"strconv"
	"time"
)

const redirectPort = ":3000"
	
func main() {
	token := "your_token"
	c := nordigen.NewClient(token)

	// supported banks in a country
	countryBanks, err := c.ListAspsps(countryCode)
	
	
	// get authorization
	endUserId := uuid.NewString()
	// look and the function below
	r, err := GetAuthorization(c, bankId, endUserId)

	// get account metadata, details and balance
	mtdt, err := c.GetAccountMetadata(r.Accounts[0])
	
	dtls, err := c.GetAccountDetails(r.Accounts[0])

	blnc, err := c.GetAccountBalances(r.Accounts[0])
	
	//get accounts transactions
	txns, err := c.GetAccountTransactions(r.Accounts[0].Id)

}

func GetAuthorization(cli nordigen.Client, bankId string, endUserId string) (nordigen.Requisition, error) {
	requisition := nordigen.Requisition{
		Redirect:  "http://localhost" + redirectPort,
		Reference: strconv.Itoa(int(time.Now().Unix())),
		EnduserId: endUserId,
		Agreements: []string{

		},
	}
	r, err := cli.CreateRequisition(requisition)

	if err != nil {
		return nordigen.Requisition{}, err
	}
	rr, err := cli.CreateRequisitionLink(r.Id, nordigen.RequisitionLinkRequest{
		AspspsId: bankId})

	if err != nil {
		return nordigen.Requisition{}, err
	}
	go internal.OpenBrowser(rr.Initiate)

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