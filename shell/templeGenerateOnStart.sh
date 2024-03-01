#!/bin/bash

echo "Starting template generation process..."

# Base directory for .tmpl files
BASE_DIR="./features"

# Counter for processed files
count=0

# Find all .tmpl files and iterate over them
find "$BASE_DIR" -name "*.templ" | while read -r FILE_PATH; do
    echo "Processing file: $FILE_PATH"

    # Extract the directory from the file path
    DIR_PATH=$(dirname "$FILE_PATH")
    echo "Changing directory to $DIR_PATH"

    # Navigate to the directory where the .tmpl file exists
    cd "$DIR_PATH" || { echo "Failed to change directory to $DIR_PATH"; exit 1; }

    # Run templ generate for the specific directory
    echo "Running templ generate in $PWD"
    /Users/earlcameron/go/bin/templ generate

    # Output to indicate which directory was processed
    echo "Successfully processed directory: $PWD"
    count=$((count+1))

    # Return to the base directory if needed (optional, depending on your workflow)
    cd - > /dev/null
done

echo "Template generation process completed. $count files were processed."

# Print out other npm scripts as suggestions
echo "You might want to run one of the following npm scripts next:"
echo "npm run setup - Install dependencies and run initial setup."
echo "npm run tailwinds - Generate tailwind CSS."
echo "npm run dev - Start the development server with nodemon."
echo "npm run templ - Watch for changes in .templ files and regenerate them."
