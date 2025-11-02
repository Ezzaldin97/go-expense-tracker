package src

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
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

func idSetter(name string) int {
	file, err := os.OpenFile(fmt.Sprintf("data/%s/config.yaml", name), os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		os.Exit(1)
	}

	var allExpenses []ExpenseIds
	if len(data) > 0 {
		err = yaml.Unmarshal(data, &allExpenses)
		if err != nil {
			os.Exit(1)
		}
	}

	var newId int

	if len(allExpenses) == 0 {
		// it works like a constructor if no config file for the name...
		newId = 1
		newCategory := ExpenseIds{
			Ids:    []int{newId},
			Budget: -1,
		}
		allExpenses = append(allExpenses, newCategory)
	} else {
		// override the existing ids
		newId = len(allExpenses[0].Ids) + 1
		allExpenses[0].Ids = append(allExpenses[0].Ids, newId)
	}

	newData, err := yaml.Marshal(&allExpenses)
	if err != nil {
		os.Exit(1)
	}

	err = file.Truncate(0)
	if err != nil {
		os.Exit(1)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		os.Exit(1)
	}

	_, err = file.Write(newData)
	if err != nil {
		os.Exit(1)
	}

	return newId
}

func ExpensesWriter(name string, description string, amount float64, logsFile *os.File) {
	CreateDir("data")
	CreateDir(fmt.Sprintf("data/%s", name))
	nextId := idSetter(name)
	filename := fmt.Sprintf("data/%v/%s.json", name, strconv.Itoa(nextId))
	expenseFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("Failed to open expense tracker file: %v", err)
		os.Exit(1)
	}
	defer expenseFile.Close()
	encoder := json.NewEncoder(expenseFile)
	expense := Expense{
		ID:          nextId,
		Amount:      amount,
		Category:    "General",
		Description: description,
		Date:        time.Now(),
	}
	mw := io.MultiWriter(os.Stdout, logsFile)
	err = encoder.Encode(expense)
	if err != nil {
		logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("Failed to write to expense tracker file: %v", err)
	}
	logger := log.New(mw, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.SetPrefix("INFO: ")
	logger.Printf("Expense added successfully (ID: %d)", nextId)
}

func ListExpenses(name string, logsFile *os.File) {
	mw := io.MultiWriter(os.Stdout, logsFile)
	file, err := os.OpenFile(fmt.Sprintf("data/%s/config.yaml", name), os.O_RDWR, 0666)
	if err != nil {
		logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("Failed to Get Expenses: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("Failed to Get Expenses: %v", err)
		os.Exit(1)
	}

	var allExpenses []ExpenseIds
	if len(data) > 0 {
		err = yaml.Unmarshal(data, &allExpenses)
		if err != nil {
			logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
			logger.Printf("Failed to Get Expenses: %v", err)
			os.Exit(1)
		}
	} else {
		logger := log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("No expenses found.")
	}

	for _, id := range allExpenses[0].Ids {
		filename := fmt.Sprintf("data/%v/%s.json", name, strconv.Itoa(id))
		expenseFile, err := os.OpenFile(filename, os.O_RDONLY, 0666)
		if err != nil {
			logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
			logger.Printf("Failed to Get Expenses: %v", err)
			continue
		}
		defer expenseFile.Close()

		decoder := json.NewDecoder(expenseFile)
		var expense Expense
		err = decoder.Decode(&expense)
		if err != nil {
			logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
			logger.Printf("Failed to Get Expenses: %v", err)
			continue
		}
		logger := log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("ID: %d | Amount: %.2f | Category: %s | Description: %s | Date: %s", expense.ID, expense.Amount, expense.Category, expense.Description, expense.Date)
	}
}

func SummarizeExpenses(name string, logsFile *os.File) {
	mw := io.MultiWriter(os.Stdout, logsFile)
	file, err := os.OpenFile(fmt.Sprintf("data/%s/config.yaml", name), os.O_RDWR, 0666)
	if err != nil {
		logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("Failed to Get Expenses: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("Failed to Get Expenses: %v", err)
		os.Exit(1)
	}
	var allExpenses []ExpenseIds
	if len(data) > 0 {
		err = yaml.Unmarshal(data, &allExpenses)
		if err != nil {
			logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
			logger.Printf("Failed to Get Expenses: %v", err)
			os.Exit(1)
		}
	} else {
		logger := log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("Total Expenses: $0")
	}

	totalAmount := 0.0

	for _, id := range allExpenses[0].Ids {
		filename := fmt.Sprintf("data/%v/%s.json", name, strconv.Itoa(id))
		expenseFile, err := os.OpenFile(filename, os.O_RDONLY, 0666)
		if err != nil {
			logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
			logger.Printf("Failed to Get Expenses: %v", err)
			continue
		}
		defer expenseFile.Close()
		decoder := json.NewDecoder(expenseFile)
		var expense Expense
		err = decoder.Decode(&expense)
		if err != nil {
			logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
			logger.Printf("Failed to Get Expenses: %v", err)
			continue
		}

		totalAmount += expense.Amount
	}
	logger := log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Printf("Total Amount: $%v", totalAmount)
}

func SetBudget(name string, budget float64, logsFile *os.File) {
	CreateDir("data")
	CreateDir(fmt.Sprintf("data/%s", name))
	mw := io.MultiWriter(os.Stdout, logsFile)
	file, err := os.OpenFile(fmt.Sprintf("data/%s/config.yaml", name), os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("Failed to set budget: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("Failed to set budget: %v", err)
		os.Exit(1)
	}

	var allExpenses []ExpenseIds
	if len(data) > 0 {
		err = yaml.Unmarshal(data, &allExpenses)
		if err != nil {
			logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
			logger.Printf("Failed to set budget: %v", err)
			os.Exit(1)
		}
	}

	if len(allExpenses) == 0 {
		// it works like a constructor if no config file for the name...
		newCategory := ExpenseIds{
			Ids:    []int{},
			Budget: budget,
		}
		allExpenses = append(allExpenses, newCategory)
	} else {
		// override the existing ids
		newExpenses := ExpenseIds{
			Ids:    allExpenses[0].Ids,
			Budget: budget,
		}
		allExpenses[0] = newExpenses
	}
	encoder := yaml.NewEncoder(file)
	err = encoder.Encode(allExpenses)
	if err != nil {
		logger := log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("Failed to set budget: %v", err)
		os.Exit(1)
	}
	logger := log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Printf("Budget set successfully to $%.2f", budget)
}
