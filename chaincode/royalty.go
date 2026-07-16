package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for music royalty accounting
type SmartContract struct {
	contractapi.Contract
}

// Rightsholder holds the accrued royalty balance for one party.
// Balance is in integer minor units (paise) to avoid floating-point errors.
type Rightsholder struct {
	HolderID string `json:"HolderID"` // unique id, e.g. "rh1"
	Name     string `json:"Name"`     // party name, e.g. "Writer"
	Role     string `json:"Role"`     // WRITER | PRODUCER | LABEL
	Balance  int    `json:"Balance"`  // accrued royalties in paise
	Status   string `json:"Status"`   // ACTIVE | SUSPENDED
}

// HistoryEntry represents one revision of a rightsholder from the ledger history.
type HistoryEntry struct {
	TxID      string        `json:"TxID"`
	Value     *Rightsholder `json:"Value"`
	Timestamp string        `json:"Timestamp"`
	IsDelete  bool          `json:"IsDelete"`
}

// RegisterRightsholder creates a new rightsholder with a zero balance and
// status "ACTIVE".
// It must fail if the rightsholder already exists.
func (s *SmartContract) RegisterRightsholder(ctx contractapi.TransactionContextInterface, holderID string, name string, role string) error {

	return nil
}

// GetRightsholder returns the rightsholder identified by holderID.
// It must fail if the rightsholder does not exist.
func (s *SmartContract) GetRightsholder(ctx contractapi.TransactionContextInterface, holderID string) (*Rightsholder, error) {

	return nil, nil
}

// CreditRoyalty adds amount to the rightsholder's balance.
// It must fail if the rightsholder does not exist, is SUSPENDED, or amount is
// not positive.
func (s *SmartContract) CreditRoyalty(ctx contractapi.TransactionContextInterface, holderID string, amount int) error {

	return nil
}

// DebitPayout subtracts amount from the rightsholder's balance (a payout).
// It must fail if the rightsholder does not exist, is SUSPENDED, amount is not
// positive, or the balance is insufficient.
func (s *SmartContract) DebitPayout(ctx contractapi.TransactionContextInterface, holderID string, amount int) error {

	return nil
}

// SuspendRightsholder sets the rightsholder's status to "SUSPENDED" (for
// example during a rights dispute). No credits or payouts are allowed while
// suspended.
// It must fail if the rightsholder does not exist or is already SUSPENDED.
func (s *SmartContract) SuspendRightsholder(ctx contractapi.TransactionContextInterface, holderID string) error {

	return nil
}

// GetRightsholderHistory returns the full revision history of a rightsholder,
// newest first, using GetHistoryForKey.
func (s *SmartContract) GetRightsholderHistory(ctx contractapi.TransactionContextInterface, holderID string) ([]HistoryEntry, error) {

	return nil, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		panic("Error creating royalty chaincode: " + err.Error())
	}

	if err := chaincode.Start(); err != nil {
		panic("Error starting royalty chaincode: " + err.Error())
	}
}
