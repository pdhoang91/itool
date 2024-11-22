#!/bin/bash

# Base directory
BASE_DIR="ExpenseTracker"

# List of files to create
FILES=(
    # App
    "App/ExpenseTrackerApp.swift"

    # Models
    "Models/User.swift"
    "Models/Transaction.swift"
    "Models/SpendingLimit.swift"
    "Models/Report.swift"

    # Views - Authentication
    "Views/Authentication/LoginView.swift"
    "Views/Authentication/RegisterView.swift"

    # Views - Main
    "Views/Main/MainTabView.swift"
    "Views/Main/HomeView.swift"
    "Views/Main/TransactionListView.swift"
    "Views/Main/AddTransactionView.swift"
    "Views/Main/ReportView.swift"
    "Views/Main/SettingsView.swift"

    # Views - Components
    "Views/Components/TransactionRowView.swift"
    "Views/Components/ChartView.swift"

    # ViewModels
    "ViewModels/AuthViewModel.swift"
    "ViewModels/TransactionViewModel.swift"
    "ViewModels/ReportViewModel.swift"
    "ViewModels/SettingsViewModel.swift"

    # Services
    "Services/APIService.swift"
    "Services/AuthService.swift"
    "Services/KeychainService.swift"

    # Utilities
    "Utilities/Constants.swift"
    "Utilities/Extensions/Color+Extension.swift"
    "Utilities/Extensions/Date+Extension.swift"
)

# Function to create files and directories
create_structure() {
    for filepath in "${FILES[@]}"; do
        # Get the full path
        full_path="$BASE_DIR/$filepath"

        # Get the directory of the file
        dir=$(dirname "$full_path")

        # Create the directory if it doesn't exist
        mkdir -p "$dir"

        # Create the file if it doesn't exist
        touch "$full_path"
    done
}

# Main script execution
echo "Creating project structure..."
create_structure
echo "Project structure created successfully in $BASE_DIR."
