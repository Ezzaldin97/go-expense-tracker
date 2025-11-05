# Go Expense Tracker

Welcome to the Go Expense Tracker! This is a simple and easy-to-use command-line application to help you keep track of your daily expenses.

## Objectives

*   Provide a straightforward way to manage personal finances from the command line.
*   Allow users to add, view, update, and delete their expenses.
*   Offer a summary of expenses to help with budgeting.
*   Enable users to set and manage a budget.

## Features

*   **Add Expense**: Add a new expense with a description and amount.
*   **List Expenses**: View a list of all your recorded expenses.
*   **Summarize Expenses**: Get a summary of your total expenses. You can also filter the summary by month.
*   **Set Budget**: Set a monthly budget to keep your spending in check.
*   **Update Expense**: Modify the details of an existing expense.
*   **Delete Expense**: Remove an expense that you no longer need.

## Installation

1.  **Prerequisites**: Make sure you have Go installed on your system. You can download it from [https://golang.org/dl/](https://golang.org/dl/).

2.  **Clone the repository**:
    ```bash
    git clone https://github.com/your-username/go-expense-tracker.git
    cd go-expense-tracker
    ```

3.  **Build the application**:
    ```bash
    go build -o go-expense-tracker .
    ```

## How to Use

The application is controlled via command-line arguments. Here’s how to use it on different operating systems:

### Windows

```powershell
.\go-expense-tracker -n "YourName" <command> [options]
```

### macOS and Linux

```bash
./go-expense-tracker -n "YourName" <command> [options]
```

### Commands

*   **`add`**: Add a new expense.
    *   `--description` or `-d`: A description of the expense (required).
    *   `--amount` or `-a`: The amount of the expense (required).
    *   **Example**:
        ```bash
        ./go-expense-tracker -n "Mohamed" add -d "Lunch" -a 15.50
        ```

*   **`list`**: List all your expenses.
    *   **Example**:
        ```bash
        ./go-expense-tracker -n "Mohamed" list
        ```

*   **`summary`**: Get a summary of your expenses.
    *   `--month` or `-m`: The month (1-12) to summarize expenses for (optional).
    *   **Example**:
        ```bash
        # Summary of all expenses
        ./go-expense-tracker -n "Mohamed" summary

        # Summary for a specific month (e.g., October)
        ./go-expense-tracker -n "Mohamed" summary -m 10
        ```

*   **`set-budget`**: Set a budget.
    *   `--budget` or `-b`: The budget amount (required).
    *   **Example**:
        ```bash
        ./go-expense-tracker -n "Mohamed" set-budget -b 500.00
        ```

*   **`update`**: Update an existing expense.
    *   `--id` or `-i`: The ID of the expense to update (required).
    *   `--description` or `-d`: The new description (optional).
    *   `--amount` or `-a`: The new amount (optional).
    *   **Example**:
        ```bash
        ./go-expense-tracker -n "Mohamed" update -i 1 -d "Coffee" -a 3.00
        ```

*   **`delete`**: Delete an expense.
    *   `--id` or `-i`: The ID of the expense to delete (required).
    *   **Example**:
        ```bash
        ./go-expense-tracker -n "Mohamed" delete -i 1
        ```

We hope you find the Go Expense Tracker useful!
