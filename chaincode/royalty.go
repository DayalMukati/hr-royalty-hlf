package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for music royalty accounting
type SmartContract struct {
	contractapi.Contract
}

// Rightsholder holds the accrued royalty balance for one party.
type Rightsholder struct {
	HolderID string `json:"HolderID"`
	Name     string `json:"Name"`
	Role     string `json:"Role"`
	Balance  int    `json:"Balance"`
	Status   string `json:"Status"`
}

// HistoryEntry represents one revision of a rightsholder from the ledger history.
type HistoryEntry struct {
	TxID      string        `json:"TxID"`
	Value     *Rightsholder `json:"Value"`
	Timestamp string        `json:"Timestamp"`
	IsDelete  bool          `json:"IsDelete"`
}

const (
	statusActive    = "ACTIVE"
	statusSuspended = "SUSPENDED"
)

// RegisterRightsholder creates a new rightsholder with a zero balance.
func (s *SmartContract) RegisterRightsholder(ctx contractapi.TransactionContextInterface, holderID string, name string, role string) error {
	existing, err := ctx.GetStub().GetState(holderID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("rightsholder %s already exists", holderID)
	}

	holder := Rightsholder{
		HolderID: holderID,
		Name:     name,
		Role:     role,
		Balance:  0,
		Status:   statusActive,
	}
	return putRightsholder(ctx, &holder)
}

// GetRightsholder returns the rightsholder identified by holderID.
func (s *SmartContract) GetRightsholder(ctx contractapi.TransactionContextInterface, holderID string) (*Rightsholder, error) {
	data, err := ctx.GetStub().GetState(holderID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if data == nil {
		return nil, fmt.Errorf("rightsholder %s does not exist", holderID)
	}

	var holder Rightsholder
	if err := json.Unmarshal(data, &holder); err != nil {
		return nil, err
	}
	return &holder, nil
}

// CreditRoyalty adds amount to the rightsholder's balance.
func (s *SmartContract) CreditRoyalty(ctx contractapi.TransactionContextInterface, holderID string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("credit amount must be positive")
	}

	holder, err := s.GetRightsholder(ctx, holderID)
	if err != nil {
		return err
	}
	if holder.Status != statusActive {
		return fmt.Errorf("rightsholder %s is not ACTIVE (current status: %s)", holderID, holder.Status)
	}

	holder.Balance += amount
	return putRightsholder(ctx, holder)
}

// DebitPayout subtracts amount from the rightsholder's balance.
func (s *SmartContract) DebitPayout(ctx contractapi.TransactionContextInterface, holderID string, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("payout amount must be positive")
	}

	holder, err := s.GetRightsholder(ctx, holderID)
	if err != nil {
		return err
	}
	if holder.Status != statusActive {
		return fmt.Errorf("rightsholder %s is not ACTIVE (current status: %s)", holderID, holder.Status)
	}
	if holder.Balance < amount {
		return fmt.Errorf("insufficient balance for %s: have %d, need %d", holderID, holder.Balance, amount)
	}

	holder.Balance -= amount
	return putRightsholder(ctx, holder)
}

// SuspendRightsholder sets the rightsholder's status to "SUSPENDED".
func (s *SmartContract) SuspendRightsholder(ctx contractapi.TransactionContextInterface, holderID string) error {
	holder, err := s.GetRightsholder(ctx, holderID)
	if err != nil {
		return err
	}
	if holder.Status == statusSuspended {
		return fmt.Errorf("rightsholder %s is already SUSPENDED", holderID)
	}

	holder.Status = statusSuspended
	return putRightsholder(ctx, holder)
}

// GetRightsholderHistory returns the full revision history, newest first.
func (s *SmartContract) GetRightsholderHistory(ctx contractapi.TransactionContextInterface, holderID string) ([]HistoryEntry, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(holderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for %s: %v", holderID, err)
	}
	defer resultsIterator.Close()

	var history []HistoryEntry
	for resultsIterator.HasNext() {
		modification, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		entry := HistoryEntry{
			TxID:      modification.TxId,
			Timestamp: time.Unix(modification.Timestamp.Seconds, int64(modification.Timestamp.Nanos)).UTC().Format(time.RFC3339),
			IsDelete:  modification.IsDelete,
		}
		if !modification.IsDelete {
			var holder Rightsholder
			if err := json.Unmarshal(modification.Value, &holder); err != nil {
				return nil, err
			}
			entry.Value = &holder
		}
		history = append(history, entry)
	}
	return history, nil
}

// --- helpers ---

func putRightsholder(ctx contractapi.TransactionContextInterface, holder *Rightsholder) error {
	bytes, err := json.Marshal(holder)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(holder.HolderID, bytes)
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
