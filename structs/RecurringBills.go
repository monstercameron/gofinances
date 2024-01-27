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
	bills.Bills = append(bills.Bills, RecurrinBill{Id: 0, Name: "Rent", DayOfMonth: 17, Amount: 1000, Owner: "Cameron", Notes: "Monthly rent for apartment"})
	bills.Bills = append(bills.Bills, RecurrinBill{Id: 1, Name: "Car Payment", DayOfMonth: 2, Amount: 300, Owner: "Cameron", Notes: "Car loan payment"})
	bills.Bills = append(bills.Bills, RecurrinBill{Id: 2, Name: "Car Insurance", DayOfMonth: 3, Amount: 100, Owner: "Cameron", Notes: "Quarterly car insurance"})
	bills.Bills = append(bills.Bills, RecurrinBill{Id: 3, Name: "Internet Bill", DayOfMonth: 12, Amount: 60, Owner: "Alex", Notes: "Monthly internet service fee"})
	bills.Bills = append(bills.Bills, RecurrinBill{Id: 4, Name: "Electricity Bill", DayOfMonth: 20, Amount: 75, Owner: "Taylor", Notes: "Monthly electricity usage"})
	bills.Bills = append(bills.Bills, RecurrinBill{Id: 5, Name: "Gym Membership", DayOfMonth: 5, Amount: 40, Owner: "Jordan", Notes: "Annual gym subscription"})
	bills.Bills = append(bills.Bills, RecurrinBill{Id: 6, Name: "Streaming Service", DayOfMonth: 25, Amount: 15, Owner: "Morgan", Notes: "Monthly Netflix subscription"})
	bills.Bills = append(bills.Bills, RecurrinBill{Id: 7, Name: "Water Bill", DayOfMonth: 18, Amount: 50, Owner: "Casey", Notes: "Monthly water usage charge"})
	bills.Bills = append(bills.Bills, RecurrinBill{Id: 8, Name: "Grocery Delivery", DayOfMonth: 1, Amount: 120, Owner: "Jamie", Notes: "Monthly grocery delivery service"})
	bills.Bills = append(bills.Bills, RecurrinBill{Id: 9, Name: "Cell Phone", DayOfMonth: 11, Amount: 85, Owner: "Pat", Notes: "Monthly cell phone plan fee"})
	return bills
}

var RecurringBills RecurringBillList

func init() {
	fmt.Println("RecurringBills.init(): Populating RecurringBills...")
	RecurringBills = PopulateRecurringBills()
}
