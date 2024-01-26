#!/bin/bash

# Get the path of the changed file
FILE_PATH=$1

# Navigate to the correct directory based on the file path
if [[ "$FILE_PATH" == views/pages/*.templ ]]; then
    cd views/pages
elif [[ "$FILE_PATH" == views/components/*.templ ]]; then
    cd views/components
fi

# Run templ generate
/Users/earlcameron/go/bin/templ generate

# Navigate back to the original directory if needed
cd - 