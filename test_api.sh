#!/bin/bash

echo "=== Testing Task Manager API ==="
echo

echo "1. Testing GET /tasks (should return empty list initially)"
curl -s http://localhost:8082/tasks | jq '.' 2>/dev/null || curl -s http://localhost:8082/tasks
echo -e "\n"

echo "2. Testing mock user service"
curl -s http://localhost:8080/users/user123 | jq '.' 2>/dev/null || curl -s http://localhost:8080/users/user123
echo -e "\n"

echo "3. Testing POST /tasks with valid user"
curl -s -X POST http://localhost:8082/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Test Task", "description": "This is a test task", "status": "Pending", "user_id": "user123"}' \
  | jq '.' 2>/dev/null || curl -s -X POST http://localhost:8082/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Test Task", "description": "This is a test task", "status": "Pending", "user_id": "user123"}'
echo -e "\n"

echo "4. Testing GET /tasks again (should show created task)"
curl -s http://localhost:8082/tasks | jq '.' 2>/dev/null || curl -s http://localhost:8082/tasks
echo -e "\n"
