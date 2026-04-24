#!/bin/bash
#
# Generate and verify Prometheus alert rules for specific namespace and job
#
# Usage:
#   ./scripts/generate_alerts.sh <namespace> <job> [output_file] [options]
#
# Options:
#   --no-verify                 Skip metric filter verification
#   --all-modules               Include all @box_module sections (incl. grpc_server, mongodb)
#   --modules <list>            Comma-separated module ids (overrides default selection)
#   --modules=a,b               Same as --modules a,b
#
# Module ids: http_server, http_client, grpc_server, db_client, redis, mongodb, schedule, go_runtime
# Aliases:    grpc, db, database, mongo, go, ...
#
# Default (when --modules / --all-modules not used): all modules except grpc_server and mongodb.
#
# Examples:
#   ./scripts/generate_alerts.sh prod api-service
#   ./scripts/generate_alerts.sh prod api-service alerts/prod_api.yaml
#   ./scripts/generate_alerts.sh prod api-service --no-verify
#   ./scripts/generate_alerts.sh prod api-service --all-modules
#   ./scripts/generate_alerts.sh prod api-service --modules http_server,http_client,db_client
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
EXTRACT_PY="${SCRIPT_DIR}/extract_alert_modules.py"

# Default: no gRPC, no MongoDB
DEFAULT_MODULES_CSV="http_server,http_client,db_client,redis,schedule,go_runtime"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

ERRORS=0
WARNINGS=0

# Check arguments
if [ $# -lt 2 ]; then
    echo -e "${RED}Error: Missing required arguments${NC}"
    echo ""
    echo "Usage: $0 <namespace> <job> [output_file] [options]"
    echo "  See script header for --modules / --all-modules."
    echo ""
    exit 1
fi

NAMESPACE="$1"
JOB="$2"
if [ -f "docs/prometheus_alerts_template.yaml" ]; then
    TEMPLATE="docs/prometheus_alerts_template.yaml"
else
    TEMPLATE="${SCRIPT_DIR}/../docs/prometheus_alerts_template.yaml"
fi
OUTPUT=""
SKIP_VERIFY=false
ALL_MODULES=false
MODULES_CSV_USER=""

# Parse arguments
shift 2
while [ $# -gt 0 ]; do
    case "$1" in
        --no-verify)
            SKIP_VERIFY=true
            shift
            ;;
        --all-modules)
            ALL_MODULES=true
            shift
            ;;
        --modules=*)
            MODULES_CSV_USER="${1#--modules=}"
            shift
            ;;
        --modules)
            if [ -z "$2" ]; then
                echo -e "${RED}Error: --modules requires a comma-separated list${NC}" >&2
                exit 1
            fi
            MODULES_CSV_USER="$2"
            shift 2
            ;;
        *)
            OUTPUT="$1"
            shift
            ;;
    esac
done

# Set default output if not specified
if [ -z "$OUTPUT" ]; then
    OUTPUT="docs/${NAMESPACE}_${JOB}_alerts.yaml"
fi

if [ "$ALL_MODULES" = true ] && [ -n "$MODULES_CSV_USER" ]; then
    echo -e "${RED}Error: use either --all-modules or --modules, not both${NC}" >&2
    exit 1
fi

if [ "$ALL_MODULES" = true ]; then
    MODARG="ALL"
else
    if [ -n "$MODULES_CSV_USER" ]; then
        MODARG="$MODULES_CSV_USER"
    else
        MODARG="$DEFAULT_MODULES_CSV"
    fi
fi

# Check if template exists
if [ ! -f "$TEMPLATE" ]; then
    echo -e "${RED}Error: Template file not found: $TEMPLATE${NC}"
    exit 1
fi
if [ ! -f "$EXTRACT_PY" ]; then
    echo -e "${RED}Error: extract module helper not found: $EXTRACT_PY${NC}"
    exit 1
fi

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  Prometheus Alert Rules Generator & Verifier              ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${YELLOW}📋 Configuration:${NC}"
echo "  Namespace: $NAMESPACE"
echo "  Job: $JOB"
echo "  Template: $TEMPLATE"
echo "  Modules:  $MODARG"
echo "  Output: $OUTPUT"
echo "  Verify: $([ "$SKIP_VERIFY" = true ] && echo "Disabled" || echo "Enabled")"
echo ""

# ============================================================
# STEP 1: Generate Alert Rules
# ============================================================
echo -e "${YELLOW}🔨 Step 1: Generating alert rules...${NC}"

# Create output directory if it doesn't exist
mkdir -p "$(dirname "$OUTPUT")"

# Generate header + top-level group wrapper (rule fragments from template live under rules:)
cat > "$OUTPUT" << EOF
# Prometheus Alert Rules
# Generated for namespace: ${NAMESPACE}, job: ${JOB}
#
# This file is auto-generated. Do not edit manually.
# To regenerate, run:
#   ./scripts/generate_alerts.sh ${NAMESPACE} ${JOB} [output] [--all-modules | --modules a,b]
# Default modules omit grpc_server and mongodb.

groups:
  - name: box-alerts
    rules:

EOF

FRAG=$(mktemp)
trap 'rm -f "$FRAG"' EXIT

python3 "$EXTRACT_PY" "$TEMPLATE" "$MODARG" > "$FRAG"

# Add namespace and job filters to all metric queries in the extracted fragment.
# go_*: only rewrite bare "metric >" in expr (not annotation text) — keep annotations free of go_ names.
sed -E \
    -e "s/http_server_requests_total\{/http_server_requests_total{namespace=\"${NAMESPACE}\",job=\"${JOB}\",/g" \
    -e "s/http_server_requests_total\[/http_server_requests_total{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/http_server_request_duration_seconds_bucket\{/http_server_request_duration_seconds_bucket{namespace=\"${NAMESPACE}\",job=\"${JOB}\",/g" \
    -e "s/http_server_request_duration_seconds_bucket\[/http_server_request_duration_seconds_bucket{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/http_server_request_duration_seconds_count\{/http_server_request_duration_seconds_count{namespace=\"${NAMESPACE}\",job=\"${JOB}\",/g" \
    -e "s/http_server_request_duration_seconds_count\[/http_server_request_duration_seconds_count{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/http_client_requests_total\{/http_client_requests_total{namespace=\"${NAMESPACE}\",job=\"${JOB}\",/g" \
    -e "s/http_client_requests_total\[/http_client_requests_total{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/http_client_requests_inflight([^{])/http_client_requests_inflight{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}\1/g" \
    -e "s/http_client_request_duration_seconds_bucket\{/http_client_request_duration_seconds_bucket{namespace=\"${NAMESPACE}\",job=\"${JOB}\",/g" \
    -e "s/http_client_request_duration_seconds_bucket\[/http_client_request_duration_seconds_bucket{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/grpc_server_requests_total\{/grpc_server_requests_total{namespace=\"${NAMESPACE}\",job=\"${JOB}\",/g" \
    -e "s/grpc_server_requests_total\[/grpc_server_requests_total{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/grpc_server_panics_total\[/grpc_server_panics_total{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/db_client_requests_total\{/db_client_requests_total{namespace=\"${NAMESPACE}\",job=\"${JOB}\",/g" \
    -e "s/db_client_requests_total\[/db_client_requests_total{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/db_client_request_duration_seconds_count\{/db_client_request_duration_seconds_count{namespace=\"${NAMESPACE}\",job=\"${JOB}\",/g" \
    -e "s/db_client_request_duration_seconds_count\[/db_client_request_duration_seconds_count{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/db_client_request_duration_seconds_bucket\{/db_client_request_duration_seconds_bucket{namespace=\"${NAMESPACE}\",job=\"${JOB}\",/g" \
    -e "s/db_client_request_duration_seconds_bucket\[/db_client_request_duration_seconds_bucket{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/db_client_connections_in_use([^{])/db_client_connections_in_use{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}\1/g" \
    -e "s/db_client_connections_max_open([^{])/db_client_connections_max_open{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}\1/g" \
    -e "s/redis_client_requests_total\{/redis_client_requests_total{namespace=\"${NAMESPACE}\",job=\"${JOB}\",/g" \
    -e "s/redis_client_requests_total\[/redis_client_requests_total{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/redis_client_request_duration_seconds_bucket\{/redis_client_request_duration_seconds_bucket{namespace=\"${NAMESPACE}\",job=\"${JOB}\",/g" \
    -e "s/redis_client_request_duration_seconds_bucket\[/redis_client_request_duration_seconds_bucket{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/mongo_client_requests_total\{/mongo_client_requests_total{namespace=\"${NAMESPACE}\",job=\"${JOB}\",/g" \
    -e "s/mongo_client_requests_total\[/mongo_client_requests_total{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/mongo_client_request_duration_seconds_bucket\{/mongo_client_request_duration_seconds_bucket{namespace=\"${NAMESPACE}\",job=\"${JOB}\",/g" \
    -e "s/mongo_client_request_duration_seconds_bucket\[/mongo_client_request_duration_seconds_bucket{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/schedule_jobs_total\{/schedule_jobs_total{namespace=\"${NAMESPACE}\",job=\"${JOB}\",/g" \
    -e "s/schedule_jobs_total\[/schedule_jobs_total{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/schedule_job_duration_seconds_bucket\{/schedule_job_duration_seconds_bucket{namespace=\"${NAMESPACE}\",job=\"${JOB}\",/g" \
    -e "s/schedule_job_duration_seconds_bucket\[/schedule_job_duration_seconds_bucket{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/go_goroutines\{/go_goroutines{namespace=\"${NAMESPACE}\",job=\"${JOB}\",/g" \
    -e "s/go_goroutines\[/go_goroutines{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/go_goroutines[[:space:]]+>/go_goroutines{namespace=\"${NAMESPACE}\",job=\"${JOB}\"} >/g" \
    -e "s/go_threads\{/go_threads{namespace=\"${NAMESPACE}\",job=\"${JOB}\",/g" \
    -e "s/go_threads[[:space:]]+>/go_threads{namespace=\"${NAMESPACE}\",job=\"${JOB}\"} >/g" \
    -e "s/go_memstats_sys_bytes[[:space:]]+>/go_memstats_sys_bytes{namespace=\"${NAMESPACE}\",job=\"${JOB}\"} >/g" \
    -e "s/go_memstats_heap_alloc_bytes\[/go_memstats_heap_alloc_bytes{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/go_gc_duration_seconds\{/go_gc_duration_seconds{namespace=\"${NAMESPACE}\",job=\"${JOB}\",/g" \
    -e "s/go_gc_duration_seconds_count\[/go_gc_duration_seconds_count{namespace=\"${NAMESPACE}\",job=\"${JOB}\"}[/g" \
    -e "s/go_memstats_gc_cpu_fraction[[:space:]]+>/go_memstats_gc_cpu_fraction{namespace=\"${NAMESPACE}\",job=\"${JOB}\"} >/g" \
    "$FRAG" >> "$OUTPUT"

echo -e "${GREEN}✓ Alert rules generated successfully!${NC}"
echo ""

# ============================================================
# STEP 2: Verify Alert Rules (if not skipped)
# ============================================================
if [ "$SKIP_VERIFY" = true ]; then
    echo -e "${YELLOW}⚠ Verification skipped (--no-verify flag)${NC}"
    echo ""
else
    echo -e "${YELLOW}🔍 Step 2: Verifying alert rules...${NC}"
    echo ""

    # List of metrics that should have filters
    METRICS=(
        "http_server_requests_total"
        "http_server_request_duration_seconds_bucket"
        "http_server_request_duration_seconds_count"
        "http_client_requests_total"
        "http_client_requests_inflight"
        "http_client_request_duration_seconds_bucket"
        "grpc_server_requests_total"
        "grpc_server_panics_total"
        "db_client_requests_total"
        "db_client_request_duration_seconds_count"
        "db_client_request_duration_seconds_bucket"
        "db_client_connections_in_use"
        "db_client_connections_max_open"
        "redis_client_requests_total"
        "redis_client_request_duration_seconds_bucket"
        "mongo_client_requests_total"
        "mongo_client_request_duration_seconds_bucket"
        "schedule_jobs_total"
        "schedule_job_duration_seconds_bucket"
        "go_goroutines"
        "go_threads"
        "go_memstats_sys_bytes"
        "go_memstats_heap_alloc_bytes"
        "go_gc_duration_seconds"
        "go_gc_duration_seconds_count"
        "go_memstats_gc_cpu_fraction"
    )

    # Check each metric
    ERRORS=0
    WARNINGS=0

    for metric in "${METRICS[@]}"; do
        # Find lines with this metric
        if grep -q "$metric" "$OUTPUT"; then
            # Check if all occurrences have the correct filter
            unfiltered=$(grep "$metric" "$OUTPUT" | grep -v "{namespace=\"${NAMESPACE}\",job=\"${JOB}\"" || true)

            if [ -n "$unfiltered" ]; then
                echo -e "${RED}✗ Found unfiltered metric: $metric${NC}"
                echo "$unfiltered" | head -3
                echo ""
                ((ERRORS++))
            fi
        fi
    done

    # Check for any metrics that might have been missed
    missed=$(grep -E '(http_|grpc_|db_|redis_|mongo_|schedule_|go_)' "$OUTPUT" | \
             grep -v "^#" | \
             grep -v "{namespace=\"${NAMESPACE}\",job=\"${JOB}\"" | \
             grep -v "namespace:" | \
             grep -v "job:" | \
             grep -v "sum by" | \
             grep -v "rate by" | \
             grep -v "humanize" || true)

    if [ -n "$missed" ]; then
        echo -e "${YELLOW}⚠ Potentially missed metrics (may be false positives):${NC}"
        echo "$missed" | head -5
        echo ""
        ((WARNINGS++))
    fi

    # Verification summary
    if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
        echo -e "${GREEN}✓ All metrics are correctly filtered!${NC}"
    elif [ $ERRORS -eq 0 ]; then
        echo -e "${YELLOW}⚠ Verification completed with warnings${NC}"
        echo -e "${YELLOW}  (Warnings may be false positives)${NC}"
    else
        echo -e "${RED}✗ Verification failed!${NC}"
        echo -e "${RED}  Found $ERRORS unfiltered metrics${NC}"
        exit 1
    fi
    echo ""
fi

# ============================================================
# Final Summary
# ============================================================
echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  Summary                                                   ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${GREEN}✓ Generation completed${NC}"
if [ "$SKIP_VERIFY" = false ]; then
    if [ $ERRORS -eq 0 ]; then
        echo -e "${GREEN}✓ Verification passed${NC}"
    fi
fi
echo ""
echo -e "${YELLOW}📄 Output file:${NC} $OUTPUT"
echo ""
echo -e "${YELLOW}📝 Next steps:${NC}"
echo ""
echo "  1. Review the generated file:"
echo "     cat $OUTPUT"
echo ""
echo "  2. (Optional) Validate with promtool:"
echo "     promtool check rules $OUTPUT"
echo ""
echo "  3. Deploy to Kubernetes:"
echo "     kubectl create configmap ${NAMESPACE}-${JOB}-alerts \\"
echo "       --from-file=$OUTPUT \\"
echo "       -n monitoring"
echo ""
echo "  4. Or add to prometheus.yml:"
echo "     rule_files:"
echo "       - '$OUTPUT'"
echo ""

exit 0
