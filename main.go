package main

import (
	"go-expense-tracker/src"
	"io"
	"log"
	"os"

	"github.com/akamensky/argparse"
)

func main() {
	// Ensure the logs directory exists
	src.CreateDir("logs")
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("Failed to open log file: %v", err)
	}
	defer file.Close()

	mw := io.MultiWriter(os.Stdout, file)

	parser := argparse.NewParser("go-expense-tracker", "Go Expense Tracker Terminal app")
	name := parser.String("n", "name", &argparse.Options{Required: true, Help: "Your Name"})
	addOperation := parser.NewCommand("add", "add an expense with a description and amount.")
	description := addOperation.String("d", "description", &argparse.Options{Required: true, Help: "Description of the expense"})
	amount := addOperation.Float("a", "amount", &argparse.Options{Required: true, Help: "Amount of the expense"})
	listOperation := parser.NewCommand("list", "list all expenses.")
	summaryOperation := parser.NewCommand("summary", "get a summary of expenses.")
	setBudgetOperation := parser.NewCommand("set-budget", "set a budget for expenses.")
	budget := setBudgetOperation.Float("b", "budget", &argparse.Options{Required: true, Help: "Budget amount"})
	deleteOperation := parser.NewCommand("delete", "delete an expense.")
	deleteID := deleteOperation.Int("i", "id", &argparse.Options{Required: true, Help: "ID of the expense to delete"})
	updateOperation := parser.NewCommand("update", "update an expense.")
	updateID := updateOperation.Int("i", "id", &argparse.Options{Required: true, Help: "ID of the expense to update"})
	updateDescription := updateOperation.String("d", "description", &argparse.Options{Required: true, Help: "New description of the expense"})
	updateAmount := updateOperation.Float("a", "amount", &argparse.Options{Required: true, Help: "New amount of the expense"})
	err = parser.Parse(os.Args)
	if err != nil {
		logger := log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		logger.Printf("Failed to parse arguments: %v", err)
		os.Exit(1)
	}

	logger := log.New(mw, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.SetPrefix("INFO: ")
	logger.Println("Welcome to Go Expense Tracker Application")
	if addOperation.Happened() {
		src.ExpensesWriter(*name, *description, *amount, file)
	} else if listOperation.Happened() {
		src.ListExpenses(*name, file)
	} else if summaryOperation.Happened() {
		src.SummarizeExpenses(*name, file)
	} else if setBudgetOperation.Happened() {
		src.SetBudget(*name, *budget, file)
	} else if deleteOperation.Happened() {
		src.DeleteExpense(*name, *deleteID, file)
	} else if updateOperation.Happened() {
		src.UpdateExpense(*name, *updateID, *updateDescription, *updateAmount, file)
	}
}
