#!/bin/bash

echo "🚀 Running unit tests and checking code coverage before push..."

# Define paths
CONTROLLER_TEST_PATH="./test/ut/controller/..."
SERVICE_TEST_PATH="./test/ut/service/..."
REPOSITORY_TEST_PATH="./test/ut/repository/..."

# Run unit tests
echo "🔍 Running unit tests..."
go test -count=1 $CONTROLLER_TEST_PATH || exit 1
go test -count=1 $SERVICE_TEST_PATH || exit 1
go test -count=1 $REPOSITORY_TEST_PATH || exit 1

# Run coverage tests and ensure 100% coverage
echo "📊 Checking code coverage..."
COVERAGE_THRESHOLD=100.0  # Set required percentage

function check_coverage {
    PACKAGE=$1
    TEST_FILE=$2
    COVERAGE=$(go test -count=1 -cover -coverpkg=./$PACKAGE $TEST_FILE | grep 'coverage:' | awk '/coverage:/ {for (i=1; i<=NF; i++) if ($i == "coverage:") print $(i+1)}' | tr -d '%')

    if (( $(echo "$COVERAGE < $COVERAGE_THRESHOLD" | bc -l) )); then
        echo "❌ Coverage for $PACKAGE is below 100% ($COVERAGE%)"
        exit 1
    fi
    echo "✅ Coverage for $PACKAGE is 100% ($COVERAGE%)"
}

# Check coverage for all testable packages
check_coverage "controller" $CONTROLLER_TEST_PATH
check_coverage "service" $SERVICE_TEST_PATH
check_coverage "repository" $REPOSITORY_TEST_PATH

echo "✅ All tests passed, and 100% coverage met. Proceeding with push."
