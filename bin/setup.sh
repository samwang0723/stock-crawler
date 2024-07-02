#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Check if pre-commit hooks are already installed
if [ ! -f .git/hooks/pre-commit ]; then
    echo "Installing pre-commit hooks..."
    pre-commit install
else
    echo "Updating pre-commit hooks..."
    pre-commit autoupdate
fi

# Run make commands
make proto
make mock-gen
make lint

# Print success message
echo -e "\nSetup for the project has been completed successfully."
