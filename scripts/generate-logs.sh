#!/bin/bash

# Generate unstructured logs for the legacy application
echo "Generating unstructured logs for legacy application..."
for i in {1..10}
do
  echo "Legacy Log Entry $i: $(date) - This is an unstructured log message."
done

# Generate structured logs for the modern application
echo "Generating structured logs for modern application..."
for i in {1..10}
do
  echo "{\"log_entry\": \"Modern Log Entry $i\", \"timestamp\": \"$(date)\", \"level\": \"info\"}"
done

echo "Log generation completed."