package main

import "fmt"

type Employee struct {
	ID         int
	Name       string
	Department string
	Salary     float64
}

func (e Employee) Describe() string {
	return fmt.Sprintf("Employee #%d: %s (%s) — PKR %.0f/month", e.ID, e.Name, e.Department, e.Salary)
}

func (e Employee) AnnualSalary() float64 {
	return e.Salary * 12
}

type Manager struct {
	Employee
	Reports []*Employee
}

func (m Manager) Describe() string {
	// baseDescription := m.Employee.Describe()
	return fmt.Sprintf("Manager #%d: %s (%s) — PKR %.0f/month — manages %d reports", m.ID, m.Name, m.Department, m.Salary, len(m.Reports))
}

func (m *Manager) AddReport(e *Employee) {
	m.Reports = append(m.Reports, e)
}

func (m Manager) TeamAnnualCost() float64 {
	totalCost := m.AnnualSalary()
	for _, report := range m.Reports {
		totalCost += report.AnnualSalary()
	}
	return totalCost
}

func main() {
	employee1 := Employee{
		ID:         1,
		Name:       "Sara Ahmed",
		Department: "Engineering",
		Salary:     150000,
	}
	employee2 := Employee{
		ID:         2,
		Name:       "Faisal Nisar",
		Department: "Engineering",
		Salary:     150000,
	}
	employee3 := Employee{
		ID:         3,
		Name:       "Saqib Nisar",
		Department: "Engineering",
		Salary:     150000,
	}

	fmt.Println(employee1.Describe())
	fmt.Printf("Annual salary: %.f\n\n", employee1.AnnualSalary())

	employee4 := Employee{
		ID:         2,
		Name:       "Ali Khan",
		Department: "Engineering",
		Salary:     220000,
	}
	manager := Manager{
		Employee: employee4,
		Reports:  []*Employee{},
	}
	manager.AddReport(&employee1)
	manager.AddReport(&employee2)
	manager.AddReport(&employee3)
	fmt.Println(manager.Describe())
	fmt.Printf("Team annual cost: %.0f", manager.TeamAnnualCost())
}
