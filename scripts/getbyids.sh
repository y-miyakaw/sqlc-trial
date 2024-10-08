curl -X POST \
-H "Content-Type: application/json" \
-d '{
  "ids": ["1", "2", "3"],
  "color": "red"
}' \
http://localhost:8099/products/searchids