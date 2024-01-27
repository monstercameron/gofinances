package structs

import (
	"fmt"
	"sort"
)

type RecurrinBill struct {
	Id         int
	Name       string
	Amount     float64
	DayOfMonth int
	Owner      string
	Notes      string
}

type RecurringBillList struct {
	Bills []RecurrinBill
}

func (m *RecurringBillList) GetByID(id int) *RecurrinBill {
	for i := range m.Bills {
		if m.Bills[i].Id == id {
			return &m.Bills[i]
		}
	}
	return nil
}

// get total costs
func (m *RecurringBillList) GetTotalCost() float64 {
	var total float64
	for i := range m.Bills {
		total += m.Bills[i].Amount
	}
	return total
}

// / SortBy sorts bills by name, amount, owner, or recurring date
func (m *RecurringBillList) SortBy(sortBy string) {
	switch sortBy {
	case "name":
		// Sort by Name
		sort.Slice(m.Bills, func(i, j int) bool {
			return m.Bills[i].Name < m.Bills[j].Name
		})
	case "amount":
		// Sort by Amount
		sort.Slice(m.Bills, func(i, j int) bool {
			return m.Bills[i].Amount < m.Bills[j].Amount
		})
	case "owner":
		// Sort by Owner
		sort.Slice(m.Bills, func(i, j int) bool {
			return m.Bills[i].Owner < m.Bills[j].Owner
		})
	case "day":
		// Sort by Recurring Date
		sort.Slice(m.Bills, func(i, j int) bool {
			return m.Bills[i].DayOfMonth < m.Bills[j].DayOfMonth
		})
	default:
		fmt.Println("RecurringBillList.SortBy(): Invalid sortBy value: ", sortBy)
	}
}

func PopulateRecurringBills() RecurringBillList {
	var bills RecurringBillList
	bills.Bills = append(bills.Bills, RecurrinBill{Id: 0, Name: "Rent", DayOfMonth: 17, Amount: 1000, Owner: "Cameron", Notes: "This is a note"})
	bills.Bills = append(bills.Bills, RecurrinBill{Id: 1, Name: "Car Payment", DayOfMonth: 2, Amount: 300, Owner: "Cameron", Notes: "This is a note"})
	bills.Bills = append(bills.Bills, RecurrinBill{Id: 2, Name: "Car Insurance", DayOfMonth: 3, Amount: 100, Owner: "Cameron", Notes: "This is a note"})
	bills.Bills = append(bills.Bills, RecurrinBill{Id: 3, Name: "Car Payment", Amount: 300, Owner: "Cameron", Notes: "This is a note"})
	bills.Bills = append(bills.Bills, RecurrinBill{Id: 4, Name: "Car Insurance", DayOfMonth: 20, Amount: 100, Owner: "Cameron", Notes: "This is a note"})
	return bills
}

var RecurringBills RecurringBillList

func init() {
	fmt.Println("RecurringBills.init(): Populating RecurringBills...")
	RecurringBills = PopulateRecurringBills()
}
