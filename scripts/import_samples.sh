#!/usr/bin/env bash
#
# Round-trip test: import workflows from actions/starter-workflows,
# rebuild to YAML, and validate with actionlint.
#
# Usage:
#   ./scripts/import_samples.sh              # Test all starter-workflows
#   ./scripts/import_samples.sh --quick      # Test first 10 only
#   ./scripts/import_samples.sh <workflow>   # Test specific workflow
#
# Prerequisites:
#   - wetwire-github CLI must be built (run ./scripts/ci.sh first)
#   - actionlint must be installed (go install github.com/rhysd/actionlint/cmd/actionlint@latest)
#

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PACKAGE_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$PACKAGE_ROOT"

QUICK=""
SPECIFIC=""
for arg in "$@"; do
    case $arg in
        --quick)
            QUICK="1"
            ;;
        *)
            SPECIFIC="$arg"
            ;;
    esac
done

# Check prerequisites
if [ ! -f "./wetwire-github" ] && ! command -v wetwire-github &> /dev/null; then
    echo "Error: wetwire-github CLI not found"
    echo "Run ./scripts/ci.sh first to build the CLI"
    exit 1
fi

CLI="./wetwire-github"
if [ ! -f "$CLI" ]; then
    CLI="wetwire-github"
fi

# Create temp directories
SAMPLES_DIR=$(mktemp -d)
OUTPUT_DIR=$(mktemp -d)
trap "rm -rf $SAMPLES_DIR $OUTPUT_DIR" EXIT

echo "=== Round-trip testing: starter-workflows ==="
echo ""
echo "Samples dir: $SAMPLES_DIR"
echo "Output dir:  $OUTPUT_DIR"
echo ""

# Clone starter-workflows (shallow)
echo ">>> Fetching actions/starter-workflows..."
git clone --depth 1 --filter=blob:none --sparse \
    https://github.com/actions/starter-workflows.git \
    "$SAMPLES_DIR" 2>/dev/null
cd "$SAMPLES_DIR"
git sparse-checkout set ci automation code-scanning deployments pages 2>/dev/null
cd "$PACKAGE_ROOT"
echo ""

# Find all workflow files
WORKFLOWS=$(find "$SAMPLES_DIR" -name "*.yml" -o -name "*.yaml" | sort)

if [ -n "$SPECIFIC" ]; then
    WORKFLOWS=$(echo "$WORKFLOWS" | grep -i "$SPECIFIC" || true)
fi

if [ -n "$QUICK" ]; then
    WORKFLOWS=$(echo "$WORKFLOWS" | head -10)
fi

TOTAL=$(echo "$WORKFLOWS" | wc -l | tr -d ' ')
PASSED=0
FAILED=0

echo ">>> Testing $TOTAL workflows..."
echo ""

for workflow in $WORKFLOWS; do
    name=$(basename "$workflow")

    # Import workflow to Go
    go_dir="$OUTPUT_DIR/$(basename "$workflow" .yml)"
    mkdir -p "$go_dir"

    if $CLI import "$workflow" -o "$go_dir" 2>/dev/null; then
        # Build back to YAML
        if $CLI build "$go_dir" -o "$go_dir/.github/workflows" 2>/dev/null; then
            # Validate with actionlint (if available)
            if command -v actionlint &> /dev/null; then
                if actionlint "$go_dir/.github/workflows/"*.yml 2>/dev/null; then
                    echo "✓ $name"
                    ((PASSED++))
                else
                    echo "✗ $name (actionlint failed)"
                    ((FAILED++))
                fi
            else
                echo "✓ $name (no actionlint)"
                ((PASSED++))
            fi
        else
            echo "✗ $name (build failed)"
            ((FAILED++))
        fi
    else
        echo "✗ $name (import failed)"
        ((FAILED++))
    fi
done

echo ""
echo "=== Results ==="
echo "Passed: $PASSED/$TOTAL"
echo "Failed: $FAILED/$TOTAL"

if [ "$FAILED" -gt 0 ]; then
    exit 1
fi

echo ""
echo "=== All round-trip tests passed! ==="
