package models

import (
	"os"
	"path/filepath"
	"testing"
)

func setupExpenseTestFile(t *testing.T) {
	t.Helper()

	tempDir := t.TempDir()
	ExpensesCSVPath = filepath.Join(tempDir, "expenses.csv")

	if err := EnsureExpenseFile(); err != nil {
		t.Fatalf("EnsureExpenseFile() error = %v", err)
	}
}

func TestEnsureExpenseFile(t *testing.T) {
	tempDir := t.TempDir()
	ExpensesCSVPath = filepath.Join(tempDir, "expenses.csv")

	if err := EnsureExpenseFile(); err != nil {
		t.Fatalf("EnsureExpenseFile() error = %v", err)
	}

	if _, err := os.Stat(ExpensesCSVPath); err != nil {
		t.Fatalf("expected expenses CSV file to exist, got error = %v", err)
	}
}

func TestExpenseCRUD(t *testing.T) {
	setupExpenseTestFile(t)

	expense := &Expense{
		UserID:      1,
		Title:       "Lunch",
		Amount:      350.50,
		Category:    "Food",
		Note:        "Team lunch",
		ExpenseDate: "2025-06-10",
	}

	if err := CreateExpense(expense); err != nil {
		t.Fatalf("CreateExpense() error = %v", err)
	}

	got, err := GetExpenseByID(1, 1)
	if err != nil {
		t.Fatalf("GetExpenseByID() error = %v", err)
	}

	if got.Title != "Lunch" {
		t.Fatalf("Title = %s, want Lunch", got.Title)
	}

	got.Title = "Dinner"
	got.Amount = 500.00

	if err := UpdateExpense(got); err != nil {
		t.Fatalf("UpdateExpense() error = %v", err)
	}

	updated, err := GetExpenseByID(1, 1)
	if err != nil {
		t.Fatalf("GetExpenseByID() after update error = %v", err)
	}

	if updated.Title != "Dinner" {
		t.Fatalf("Title = %s, want Dinner", updated.Title)
	}

	if err := DeleteExpense(1, 1); err != nil {
		t.Fatalf("DeleteExpense() error = %v", err)
	}

	if _, err := GetExpenseByID(1, 1); err == nil {
		t.Fatal("expected deleted expense lookup to fail")
	}
}

func TestExpenseOwnershipIsolation(t *testing.T) {
	setupExpenseTestFile(t)

	expense := &Expense{
		UserID:      1,
		Title:       "Lunch",
		Amount:      350.50,
		Category:    "Food",
		ExpenseDate: "2025-06-10",
	}

	if err := CreateExpense(expense); err != nil {
		t.Fatalf("CreateExpense() error = %v", err)
	}

	if _, err := GetExpenseByID(1, 2); err == nil {
		t.Fatal("expected user 2 to be unable to access user 1 expense")
	}
}

func TestFilterExpensesByDate(t *testing.T) {
	expenses := []Expense{
		{ID: 1, ExpenseDate: "2025-06-01"},
		{ID: 2, ExpenseDate: "2025-06-15"},
		{ID: 3, ExpenseDate: "2025-07-01"},
	}

	got, err := FilterExpensesByDate(expenses, "2025-06-01", "2025-06-30")
	if err != nil {
		t.Fatalf("FilterExpensesByDate() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
}

func TestFilterExpensesByDateInvalidDate(t *testing.T) {
	expenses := []Expense{
		{ID: 1, ExpenseDate: "2025-06-01"},
	}

	if _, err := FilterExpensesByDate(expenses, "invalid", "2025-06-30"); err == nil {
		t.Fatal("expected invalid date_from error")
	}
}

func TestSortExpenses(t *testing.T) {
	tests := []struct {
		name      string
		sortBy    string
		sortOrder string
		wantFirst int
	}{
		{
			name:      "amount desc",
			sortBy:    "amount",
			sortOrder: "desc",
			wantFirst: 2,
		},
		{
			name:      "amount asc",
			sortBy:    "amount",
			sortOrder: "asc",
			wantFirst: 1,
		},
		{
			name:      "date desc default",
			sortBy:    "expense_date",
			sortOrder: "",
			wantFirst: 3,
		},
		{
			name:      "date asc",
			sortBy:    "expense_date",
			sortOrder: "asc",
			wantFirst: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expenses := []Expense{
				{ID: 1, Amount: 100, ExpenseDate: "2025-06-01"},
				{ID: 2, Amount: 300, ExpenseDate: "2025-06-10"},
				{ID: 3, Amount: 200, ExpenseDate: "2025-07-01"},
			}

			if err := SortExpenses(expenses, tt.sortBy, tt.sortOrder); err != nil {
				t.Fatalf("SortExpenses() error = %v", err)
			}

			if expenses[0].ID != tt.wantFirst {
				t.Fatalf("first ID = %d, want %d", expenses[0].ID, tt.wantFirst)
			}
		})
	}
}

func TestSortExpensesValidationFailures(t *testing.T) {
	expenses := []Expense{{ID: 1, Amount: 100, ExpenseDate: "2025-06-01"}}

	tests := []struct {
		name      string
		sortBy    string
		sortOrder string
	}{
		{name: "sort order without sort by", sortBy: "", sortOrder: "asc"},
		{name: "invalid sort by", sortBy: "category", sortOrder: "asc"},
		{name: "invalid sort order", sortBy: "amount", sortOrder: "highest"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SortExpenses(expenses, tt.sortBy, tt.sortOrder); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestBuildExpenseSummary(t *testing.T) {
	expenses := []Expense{
		{ID: 1, Amount: 100, Category: "Food"},
		{ID: 2, Amount: 200, Category: "Food"},
		{ID: 3, Amount: 50, Category: "Transport"},
	}

	summary := BuildExpenseSummary(expenses, "2025-06-01", "2025-06-30")

	if summary.TotalAmount != 350 {
		t.Fatalf("TotalAmount = %.2f, want 350.00", summary.TotalAmount)
	}

	if summary.TotalCount != 3 {
		t.Fatalf("TotalCount = %d, want 3", summary.TotalCount)
	}

	if len(summary.ByCategory) != 2 {
		t.Fatalf("len(ByCategory) = %d, want 2", len(summary.ByCategory))
	}
}

func TestGetExpensesMalformedCSV(t *testing.T) {
	tempDir := t.TempDir()
	ExpensesCSVPath = filepath.Join(tempDir, "expenses.csv")

	content := []byte("id,user_id,title,amount,category,note,expense_date,created_at\n1,1,Lunch\n")

	if err := os.WriteFile(ExpensesCSVPath, content, 0644); err != nil {
		t.Fatalf("failed to write malformed expenses CSV: %v", err)
	}

	if _, err := GetExpensesByUserID(1); err == nil {
		t.Fatal("expected malformed expenses CSV error, got nil")
	}
}

func TestFilterExpensesByDateSingleBound(t *testing.T) {
	expenses := []Expense{
		{ID: 1, ExpenseDate: "2025-06-01"},
		{ID: 2, ExpenseDate: "2025-06-15"},
		{ID: 3, ExpenseDate: "2025-07-01"},
	}

	tests := []struct {
		name     string
		dateFrom string
		dateTo   string
		wantLen  int
	}{
		{
			name:     "date_from only",
			dateFrom: "2025-06-15",
			dateTo:   "",
			wantLen:  2,
		},
		{
			name:     "date_to only",
			dateFrom: "",
			dateTo:   "2025-06-15",
			wantLen:  2,
		},
		{
			name:     "no bounds",
			dateFrom: "",
			dateTo:   "",
			wantLen:  3,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := FilterExpensesByDate(expenses, testCase.dateFrom, testCase.dateTo)
			if err != nil {
				t.Fatalf("FilterExpensesByDate() error = %v", err)
			}

			if len(got) != testCase.wantLen {
				t.Fatalf("len(got) = %d, want %d", len(got), testCase.wantLen)
			}
		})
	}
}
