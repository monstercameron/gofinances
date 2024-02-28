#!/bin/bash

# Base directory for .templ files
BASE_DIR="./features"

# Current time in seconds since the epoch
current_time=$(date +%s)

# Find all .templ files and iterate over them
find "$BASE_DIR" -name "*.templ" | while read FILE_PATH; do
    echo "Processing file: $FILE_PATH"

    # Get the last modification time of the file in seconds since the epoch
    file_mod_time=$(date -r "$FILE_PATH" +%s)

    # Calculate the difference in time from now
    time_diff=$((current_time - file_mod_time))

    # Check if the file was modified in the last 10 seconds
    if [ $time_diff -le 10 ]; then
        echo "File $FILE_PATH has been updated in the last 10 seconds."

        # Extract the directory from the file path
        DIR_PATH=$(dirname "$FILE_PATH")

        # Navigate to the directory where the .templ file exists
        cd "$DIR_PATH" || exit 1

        # Run templ generate for the specific directory
        /Users/earlcameron/go/bin/templ generate

        # Output to indicate which directory was processed
        echo "Processed directory: $PWD"

        # Return to the base directory if needed (optional, depending on your workflow)
        cd - > /dev/null
    # else
        # echo "File $FILE_PATH was not modified in the last 10 seconds."
    fi
done
