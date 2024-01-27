package structs

import (
	"fmt"
	"sort"
)

// RecurringBill represents a single recurring bill.
// It is used to store information about a bill that is paid on a regular basis.
type RecurringBill struct {
	Id         int
	Name       string
	Amount     float64
	DayOfMonth int
	Owner      string
	Notes      string
}

// RecurringBillList represents a list of RecurringBill.
type RecurringBillList struct {
	Bills []RecurringBill
}

// GetByID returns a pointer to the RecurringBill with the given ID.
// Returns nil if no bill is found.
func (m *RecurringBillList) GetByID(id int) *RecurringBill {
	for i := range m.Bills {
		if m.Bills[i].Id == id {
			return &m.Bills[i]
		}
	}
	return nil
}

// GetTotalCost calculates and returns the total cost of all bills.
// It is used to calculate the total cost of all bills in the list.
func (m *RecurringBillList) GetTotalCost() float64 {
	var total float64
	for _, bill := range m.Bills {
		total += bill.Amount
	}
	return total
}

// RemoveByID removes the bill with the specified ID from the list.
func (m *RecurringBillList) RemoveByID(id int) {
	for i, bill := range m.Bills {
		if bill.Id == id {
			m.Bills = append(m.Bills[:i], m.Bills[i+1:]...)
			return
		}
	}
}

// SortBy sorts the bills based on the provided criterion: name, amount, owner, or day of the month.
// It is used to sort the bills in the list by a specific criterion.
func (m *RecurringBillList) SortBy(sortBy string) {
	switch sortBy {
	case "name":
		sort.Slice(m.Bills, func(i, j int) bool { return m.Bills[i].Name < m.Bills[j].Name })
	case "amount":
		sort.Slice(m.Bills, func(i, j int) bool { return m.Bills[i].Amount < m.Bills[j].Amount })
	case "owner":
		sort.Slice(m.Bills, func(i, j int) bool { return m.Bills[i].Owner < m.Bills[j].Owner })
	case "day":
		sort.Slice(m.Bills, func(i, j int) bool { return m.Bills[i].DayOfMonth < m.Bills[j].DayOfMonth })
	default:
		fmt.Println("Invalid sortBy value:", sortBy)
	}
}

// GetLastID returns the ID of the last bill in the list.
func (m *RecurringBillList) GetLastID() int {
	if len(m.Bills) > 0 {
		return m.Bills[len(m.Bills)-1].Id
	}
	return 0 // Return 0 if the list is empty
}

// PopulateRecurringBills creates and returns a RecurringBillList with initial data.
func PopulateRecurringBills() RecurringBillList {
	var bills RecurringBillList
	bills.Bills = append(bills.Bills, RecurringBill{Id: 0, Name: "Rent", DayOfMonth: 17, Amount: 1000, Owner: "Cameron", Notes: "Monthly rent for apartment"})
	bills.Bills = append(bills.Bills, RecurringBill{Id: 1, Name: "Car Payment", DayOfMonth: 2, Amount: 300, Owner: "Cameron", Notes: "Car loan payment"})
	bills.Bills = append(bills.Bills, RecurringBill{Id: 2, Name: "Car Insurance", DayOfMonth: 3, Amount: 100, Owner: "Cameron", Notes: "Quarterly car insurance"})
	bills.Bills = append(bills.Bills, RecurringBill{Id: 3, Name: "Internet Bill", DayOfMonth: 12, Amount: 60, Owner: "Alex", Notes: "Monthly internet service fee"})
	bills.Bills = append(bills.Bills, RecurringBill{Id: 4, Name: "Electricity Bill", DayOfMonth: 20, Amount: 75, Owner: "Taylor", Notes: "Monthly electricity usage"})
	bills.Bills = append(bills.Bills, RecurringBill{Id: 5, Name: "Gym Membership", DayOfMonth: 5, Amount: 40, Owner: "Jordan", Notes: "Annual gym subscription"})
	bills.Bills = append(bills.Bills, RecurringBill{Id: 6, Name: "Streaming Service", DayOfMonth: 25, Amount: 15, Owner: "Morgan", Notes: "Monthly Netflix subscription"})
	bills.Bills = append(bills.Bills, RecurringBill{Id: 7, Name: "Water Bill", DayOfMonth: 18, Amount: 50, Owner: "Casey", Notes: "Monthly water usage charge"})
	bills.Bills = append(bills.Bills, RecurringBill{Id: 8, Name: "Grocery Delivery", DayOfMonth: 1, Amount: 120, Owner: "Jamie", Notes: "Monthly grocery delivery service"})
	bills.Bills = append(bills.Bills, RecurringBill{Id: 9, Name: "Cell Phone", DayOfMonth: 11, Amount: 85, Owner: "Pat", Notes: "Monthly cell phone plan fee"})
	return bills
}

// RecurringBills is a global instance of RecurringBillList.
var RecurringBills RecurringBillList

// init populates RecurringBills with initial data.
func init() {
	fmt.Println("Populating RecurringBills...")
	RecurringBills = PopulateRecurringBills()
}
