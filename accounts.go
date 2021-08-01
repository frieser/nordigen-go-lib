package nordigen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type AccountMetadata struct {
	Id              string `json:"id"`
	Created         string `json:"created"`
	LastAccessed    string `json:"last_accessed"`
	Iban            string `json:"iban"`
	AspspIdentifier string `json:"aspsp_identifier"`
	Status          string `json:"status"`
}

type AccountBalanceAmount struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type AccountBalance struct {
	BalanceAmount AccountBalanceAmount `json:"balanceAmount"`
	BalanceType   string               `json:"balanceType"`
}

type AccountBalances struct {
	Balances []AccountBalance `json:"balances"`
}

type AccountDetails struct {
	Account struct {
		ResourceId string `json:"resourceId"`
		Iban       string `json:"iban"`
		Currency   string `json:"currency"`
		OwnerName  string `json:"ownerName"`
		Product    string `json:"product"`
		Status     string `json:"status"`
	} `json:"account"`
}

type AccountTransactions struct {
	Transactions struct {
		Booked []struct {
			TransactionId     string `json:"transactionId"`
			EntryReference    string `json:"entryReference"`
			BookingDate       string `json:"bookingDate"`
			ValueDate         string `json:"valueDate"`
			TransactionAmount struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
			} `json:"transactionAmount"`
			CreditorName    string `json:"creditorName,omitempty"`
			CreditorAccount struct {
				Iban string `json:"iban"`
			} `json:"creditorAccount"`
			UltimateCreditor string `json:"ultimateCreditor,omitempty"`
			DebtorName       string `json:"debtorName,omitempty"`
			DebtorAccount    struct {
				Iban string `json:"iban"`
			} `json:"debtorAccount,omitempty"`
			UltimateDebtor                    string `json:"ultimateDebtor,omitempty"`
			RemittanceInformationUnstructured string `json:"remittanceInformationUnstructured"`
			BankTransactionCode               string `json:"bankTransactionCode,omitempty"`
		} `json:"booked"`
		Pending []struct {
			TransactionAmount struct {
				Amount   string `json:"amount"`
				Currency string `json:"currency"`
			} `json:"transactionAmount"`
			ValueDate                         string `json:"valueDate"`
			RemittanceInformationUnstructured string `json:"remittanceInformationUnstructured"`
		} `json:"pending"`
	} `json:"transactions"`
}

const accountPath = "accounts"
const balancesPath = "balances"
const detailsPath = "details"
const transactionsPath = "transactions"

func (c Client) GetAccountMetadata(id string) (AccountMetadata, error) {
	req := http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: strings.Join([]string{accountPath, id, ""}, "/"),
		},
	}
	resp, err := c.c.Do(&req)

	if err != nil {
		return AccountMetadata{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return AccountMetadata{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return AccountMetadata{}, fmt.Errorf("expected %d status code: got %d", http.StatusOK, resp.StatusCode)
	}
	accMtdt := AccountMetadata{}
	err = json.Unmarshal(body, &accMtdt)

	if err != nil {
		return AccountMetadata{}, err
	}

	return accMtdt, nil
}

func (c Client) GetAccountBalances(id string) (AccountBalances, error) {
	req := http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: strings.Join([]string{accountPath, id, balancesPath, ""}, "/"),
		},
	}
	resp, err := c.c.Do(&req)

	if err != nil {
		return AccountBalances{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return AccountBalances{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return AccountBalances{}, fmt.Errorf("expected %d status code: got %d", http.StatusOK, resp.StatusCode)
	}
	accBlnc := AccountBalances{}
	err = json.Unmarshal(body, &accBlnc)

	if err != nil {
		return AccountBalances{}, err
	}

	return accBlnc, nil
}

func (c Client) GetAccountDetails(id string) (AccountDetails, error) {
	req := http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: strings.Join([]string{accountPath, id, detailsPath, ""}, "/"),
		},
	}
	resp, err := c.c.Do(&req)

	if err != nil {
		return AccountDetails{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return AccountDetails{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return AccountDetails{}, fmt.Errorf("expected %d status code: got %d", http.StatusOK, resp.StatusCode)
	}
	accDtl := AccountDetails{}
	err = json.Unmarshal(body, &accDtl)

	if err != nil {
		return AccountDetails{}, err
	}

	return accDtl, nil
}

func (c Client) GetAccountTransactions(id string) (AccountTransactions, error) {
	req := http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: strings.Join([]string{accountPath, id, transactionsPath, ""}, "/"),
		},
	}
	resp, err := c.c.Do(&req)

	if err != nil {
		return AccountTransactions{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return AccountTransactions{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return AccountTransactions{}, fmt.Errorf("expected %d status code: got %d", http.StatusOK, resp.StatusCode)
	}
	accTxns := AccountTransactions{}
	err = json.Unmarshal(body, &accTxns)

	if err != nil {
		return AccountTransactions{}, err
	}

	return accTxns, nil
}
