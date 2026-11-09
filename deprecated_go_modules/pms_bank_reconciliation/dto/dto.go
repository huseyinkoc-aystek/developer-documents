package dto

type IngestRequest struct {
	ExternalTxnID string  `json:"external_txn_id"`
	BankAccountID string  `json:"bank_account_id"`
	TxnDate       string  `json:"txn_date"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Direction     string  `json:"direction"`
	Description   string  `json:"description"`
	Counterparty  string  `json:"counterparty"`
}

type MatchRequest struct {
	UseAutoProposal bool `json:"use_auto_proposal"`
}

type ConfirmRequest struct {
	CandidateChargeID string `json:"candidate_charge_id,omitempty"`
	CreatePayment     bool   `json:"create_payment,omitempty"`
}

type RejectRequest struct {
	Reason string `json:"reason,omitempty"`
}
