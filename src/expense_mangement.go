package src

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Expense struct {
	ID          int       `json:"ID"`
	Amount      float64   `json:"Amount"`
	Category    string    `json:"Category"`
	Description string    `json:"Description"`
	Date        time.Time `json:"Date"`
}

type ExpenseIds struct {
	Ids    []int   `yaml:"ids"`
	Budget float64 `yaml:"budget"`
}

// handleError logs an error and exits the program.
func handleError(err error, mw io.Writer, message string) {
	if err != nil {
		logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("%s: %v", message, err)
		os.Exit(1)
	}
}

// getConfigPath returns the path to the config file.
func getConfigPath(name string) string {
	return filepath.Join("data", name, "config.yaml")
}

// getExpensePath returns the path to an expense file.
func getExpensePath(name string, id int) string {
	return filepath.Join("data", name, fmt.Sprintf("%d.json", id))
}

// readConfig reads the config.yaml file and returns the unmarshalled data.
func readConfig(name string, mw io.Writer) ([]ExpenseIds, *os.File) {
	configPath := getConfigPath(name)
	file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE, 0666)
	handleError(err, mw, "Failed to open config file")

	data, err := io.ReadAll(file)
	handleError(err, mw, "Failed to read config file")

	var allExpenses []ExpenseIds
	if len(data) > 0 {
		err = yaml.Unmarshal(data, &allExpenses)
		handleError(err, mw, "Failed to unmarshal config data")
	}
	return allExpenses, file
}

// writeConfig writes the config data back to the config.yaml file.
func writeConfig(file *os.File, allExpenses []ExpenseIds, mw io.Writer) {
	newData, err := yaml.Marshal(&allExpenses)
	handleError(err, mw, "Failed to marshal config data")

	err = file.Truncate(0)
	handleError(err, mw, "Failed to truncate config file")

	_, err = file.Seek(0, 0)
	handleError(err, mw, "Failed to seek in config file")

	_, err = file.Write(newData)
	handleError(err, mw, "Failed to write to config file")
}

// idSetter gets the next available ID for a new expense.
func idSetter(name string, mw io.Writer) int {
	allExpenses, file := readConfig(name, mw)
	defer file.Close()

	var newId int
	if len(allExpenses) == 0 {
		newId = 1
		newCategory := ExpenseIds{
			Ids:    []int{newId},
			Budget: -1,
		}
		allExpenses = append(allExpenses, newCategory)
	} else {
		if len(allExpenses[0].Ids) == 0 {
			newId = 1
		} else {
			// Find the max ID and add 1 to it for the new ID.
			maxId := 0
			for _, id := range allExpenses[0].Ids {
				if id > maxId {
					maxId = id
				}
			}
			newId = maxId + 1
		}
		allExpenses[0].Ids = append(allExpenses[0].Ids, newId)
	}

	writeConfig(file, allExpenses, mw)
	return newId
}

func ExpensesWriter(name string, description string, amount float64, logsFile *os.File) {
	mw := io.MultiWriter(os.Stdout, logsFile)
	CreateDir("data")
	CreateDir(filepath.Join("data", name))

	nextId := idSetter(name, mw)
	filename := getExpensePath(name, nextId)

	expenseFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	handleError(err, mw, "Failed to open expense tracker file")
	defer expenseFile.Close()

	encoder := json.NewEncoder(expenseFile)
	expense := Expense{
		ID:          nextId,
		Amount:      amount,
		Category:    "General",
		Description: description,
		Date:        time.Now(),
	}

	err = encoder.Encode(expense)
	if err != nil {
		logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("Failed to write to expense tracker file: %v", err)
	} else {
		logger := log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("Expense added successfully (ID: %d)", nextId)
	}
}

func ListExpenses(name string, logsFile *os.File) {
	mw := io.MultiWriter(os.Stdout, logsFile)
	allExpenses, file := readConfig(name, mw)
	defer file.Close()

	if len(allExpenses) == 0 || len(allExpenses[0].Ids) == 0 {
		log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile).Printf("No expenses found.")
		return
	}

	for _, id := range allExpenses[0].Ids {
		expense, err := readExpense(name, id, mw)
		if err != nil {
			continue
		}
		log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile).Printf("ID: %d | Amount: %.2f | Category: %s | Description: %s | Date: %s", expense.ID, expense.Amount, expense.Category, expense.Description, expense.Date)
	}
}

func readExpense(name string, id int, mw io.Writer) (*Expense, error) {
	filename := getExpensePath(name, id)
	expenseFile, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile).Printf("Failed to open expense file: %v", err)
		return nil, err
	}
	defer expenseFile.Close()

	decoder := json.NewDecoder(expenseFile)
	var expense Expense
	err = decoder.Decode(&expense)
	if err != nil {
		log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile).Printf("Failed to decode expense file: %v", err)
		return nil, err
	}
	return &expense, nil
}

func SummarizeExpenses(name string, logsFile *os.File) {
	mw := io.MultiWriter(os.Stdout, logsFile)
	allExpenses, file := readConfig(name, mw)
	defer file.Close()

	if len(allExpenses) == 0 || len(allExpenses[0].Ids) == 0 {
		log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile).Printf("Total Expenses: $0")
		return
	}

	totalAmount := 0.0
	for _, id := range allExpenses[0].Ids {
		expense, err := readExpense(name, id, mw)
		if err != nil {
			continue
		}
		totalAmount += expense.Amount
	}
	log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile).Printf("Total Amount: $%v", totalAmount)
}

func SetBudget(name string, budget float64, logsFile *os.File) {
	mw := io.MultiWriter(os.Stdout, logsFile)
	CreateDir("data")
	CreateDir(filepath.Join("data", name))

	allExpenses, file := readConfig(name, mw)
	defer file.Close()

	if len(allExpenses) == 0 {
		newCategory := ExpenseIds{
			Ids:    []int{},
			Budget: budget,
		}
		allExpenses = append(allExpenses, newCategory)
	} else {
		allExpenses[0].Budget = budget
	}

	writeConfig(file, allExpenses, mw)
	log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile).Printf("Budget set successfully to $%.2f", budget)
}

func DeleteExpense(name string, id int, logsFile *os.File) {
	mw := io.MultiWriter(os.Stdout, logsFile)
	allExpenses, file := readConfig(name, mw)
	defer file.Close()

	if len(allExpenses) == 0 || len(allExpenses[0].Ids) == 0 {
		log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile).Printf("No expenses found.")
		return
	}

	indexToRemove := -1
	for i, expenseID := range allExpenses[0].Ids {
		if expenseID == id {
			indexToRemove = i
			break
		}
	}

	if indexToRemove == -1 {
		log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile).Printf("Expense with ID %d not found.", id)
		return
	}

	allExpenses[0].Ids = remove(allExpenses[0].Ids, indexToRemove)
	filePath := getExpensePath(name, id)
	err := os.Remove(filePath)
	if err != nil {
		log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile).Printf("Error removing file: %v\n", err)
	}

	writeConfig(file, allExpenses, mw)
	log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile).Printf("Expense with ID %d deleted successfully.", id)
}

func UpdateExpense(name string, id int, logsFile *os.File, description string, amount float64) {
	mw := io.MultiWriter(os.Stdout, logsFile)
	expense, err := readExpense(name, id, mw)
	if err != nil {
		log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile).Printf("Failed to read expense for update: %v", err)
		return
	}

	if description != "" {
		expense.Description = description
	}

	if amount != -1.0 {
		expense.Amount = amount
	}
	expense.Date = time.Now()

	filePath := getExpensePath(name, id)
	expenseFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("Failed to open expense file for update: %v", err)
		return
	}
	defer expenseFile.Close()

	encoder := json.NewEncoder(expenseFile)
	err = encoder.Encode(expense)
	if err != nil {
		logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("Failed to encode updated expense: %v", err)
	} else {
		logger := log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("Expense with ID %d updated successfully.", id)
	}
}
