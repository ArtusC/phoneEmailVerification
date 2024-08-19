#!/bin/bash

cmd="go test -timeout 800s -tags=integration -v"

# Function to display usage information
display_usage() {
    echo "Usage: $0 testfile=<PTAH_TO_TEST_FILE> [testname=<NAME_OF_THE_TEST>]"
}

if [ -z "$testfile" ] && [ -z "$testname" ] && [ -z "$race" ]; then
    echo "Running all tests"
    $cmd ./...
    exit
fi

# To run a specific test, the -testfile flag must be provided
if [ -z "$testfile" ]; then
    echo "Error: testfile parameter is required."
    display_usage
    exit 1
fi

# Iterate over parameters
for param in "$@"; do
    case $param in
        testfile=*)
            # Extract testfile value
            testfile_value="${param#*=}"
            ;;
        testname=*)
            # Extract testname value
            testname_value="${param#*=}"
            ;;
        *)
            # Unknown parameter
            echo "Error: Unknown parameter $param"
            display_usage
            exit 1
            ;;
    esac
done

# Check conditions for testname parameter
if [ -n "$testfile_value" ] && [ -n "$testname_value" ]; then
    # Both testfile and testname are provided
    echo "Running specific test: testfile=$testfile_value, testname=$testname_value"
    cmd="${cmd} ${testfile} -run ^${testname}$"
elif [ -n "$testfile_value" ]; then
    # Only testfile is provided
    echo "Running test file: testfile=$testfile_value"
    cmd="${cmd} ${testfile}"
else
    echo "Error: testfile parameter is required."
    display_usage
    exit 1
fi

echo "CMD: ${cmd}"
$cmd