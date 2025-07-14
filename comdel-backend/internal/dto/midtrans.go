package dto

type PremiumPlan int;

const (
	Newbie PremiumPlan = iota
	Creator
)


type TransactionStatus struct {
	OrderID		string;
	StatusCode	string;
	Status 		string;
}