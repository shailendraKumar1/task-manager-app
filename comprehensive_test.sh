#!/bin/bash

echo "=== Comprehensive Task Manager API Test ==="
echo

echo "1. Testing GET /tasks (list all tasks)"
curl -s http://localhost:8082/tasks | jq '.tasks | length' | xargs echo "Number of tasks:"
echo

echo "2. Testing POST /tasks (create new task)"
TASK_RESPONSE=$(curl -s -X POST http://localhost:8082/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Integration Test Task", "description": "Created via comprehensive test", "status": "Pending", "user_id": "test_user_456"}')
echo "$TASK_RESPONSE" | jq '.'
TASK_UUID=$(echo "$TASK_RESPONSE" | jq -r '.uuid')
echo "Created task UUID: $TASK_UUID"
echo

echo "3. Testing GET /tasks/:uuid (get specific task)"
curl -s "http://localhost:8082/tasks/$TASK_UUID" | jq '.'
echo

echo "4. Testing PUT /tasks/:uuid (update task)"
curl -s -X PUT "http://localhost:8082/tasks/$TASK_UUID" \
  -H "Content-Type: application/json" \
  -d '{"title": "Updated Integration Test Task", "description": "Updated via comprehensive test", "status": "InProgress"}' | jq '.'
echo

echo "5. Testing GET /tasks with status filter"
curl -s "http://localhost:8082/tasks?status=InProgress" | jq '.tasks | length' | xargs echo "InProgress tasks:"
echo

echo "6. Testing user validation with mock service"
curl -s "http://localhost:8080/api/users/test_user_456/validate" | jq '.'
echo

echo "7. Testing invalid user ID"
curl -s -w "\nStatus: %{http_code}\n" -X POST http://localhost:8082/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Invalid User Test", "description": "Testing invalid user", "status": "Pending", "user_id": "invalid_user"}'
echo

echo "8. Testing DELETE /tasks/:uuid"
curl -s -w "\nStatus: %{http_code}\n" -X DELETE "http://localhost:8082/tasks/$TASK_UUID"
echo

echo "=== Test Complete ==="
