package nordigen

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type AccountMetadata struct {
	Id            string `json:"id,omitempty"`
	Created       string `json:"created,omitempty"`
	LastAccessed  string `json:"last_accessed,omitempty"`
	Iban          string `json:"iban,omitempty"`
	InstitutionId string `json:"institution_id,omitempty"`
	// There is an issue in the api, the status is still a string
	// like in v1
	Status string `json:"status,omitempty"`
	//Status        []string `json:"status"`
}

type AccountBalanceAmount struct {
	Amount   string `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type AccountBalance struct {
	BalanceAmount AccountBalanceAmount `json:"balanceAmount,omitempty"`
	BalanceType   string               `json:"balanceType,omitempty"`
}

type AccountBalances struct {
	Balances []AccountBalance `json:"balances,omitempty"`
}

type AccountDetails struct {
	Account struct {
		ResourceId string `json:"resourceId,omitempty"`
		Iban       string `json:"iban,omitempty"`
		Currency   string `json:"currency,omitempty"`
		OwnerName  string `json:"ownerName,omitempty"`
		Product    string `json:"product,omitempty,"`
		Status     string `json:"status,omitempty"`
	} `json:"account"`
}

type Transaction struct {
	TransactionId     string `json:"transactionId,omitempty"`
	EntryReference    string `json:"entryReference,omitempty"`
	BookingDate       string `json:"bookingDate,omitempty"`
	ValueDate         string `json:"valueDate,omitempty"`
	TransactionAmount struct {
		Amount   string `json:"amount,omitempty"`
		Currency string `json:"currency,omitempty"`
	} `json:"transactionAmount,omitempty"`
	CreditorName    string `json:"creditorName,omitempty"`
	CreditorAccount struct {
		Iban string `json:"iban,omitempty"`
	} `json:"creditorAccount,omitempty"`
	UltimateCreditor string `json:"ultimateCreditor,omitempty"`
	DebtorName       string `json:"debtorName,omitempty"`
	DebtorAccount    struct {
		Iban string `json:"iban,omitempty"`
	} `json:"debtorAccount,omitempty"`
	UltimateDebtor                         string   `json:"ultimateDebtor,omitempty"`
	RemittanceInformationUnstructured      string   `json:"remittanceInformationUnstructured"`
	RemittanceInformationUnstructuredArray []string `json:"RemittanceInformationUnstructuredArray"`
	BankTransactionCode                    string   `json:"bankTransactionCode,omitempty"`
	AdditionalInformation                  string   `json:"additionalInformation,omitempty"`
}

type AccountTransactions struct {
	Transactions struct {
		Booked  []Transaction `json:"booked,omitempty"`
		Pending []Transaction `json:"pending,omitempty"`
	} `json:"transactions,omitempty"`
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
		return AccountMetadata{}, &APIError{resp.StatusCode, string(body), err}
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
		return AccountBalances{}, &APIError{resp.StatusCode, string(body), err}
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
		return AccountDetails{}, &APIError{resp.StatusCode, string(body), err}
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
		return AccountTransactions{}, &APIError{resp.StatusCode, string(body), err}
	}
	accTxns := AccountTransactions{}
	err = json.Unmarshal(body, &accTxns)

	if err != nil {
		return AccountTransactions{}, err
	}

	return accTxns, nil
}
