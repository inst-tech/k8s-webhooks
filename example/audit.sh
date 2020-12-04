#!/bin/bash
curl -d @audit.json -X POST http://localhost:8000/audit
curl http://localhost:8000/metrics
