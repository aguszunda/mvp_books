#!/bin/bash

# Create some books
echo "Creating books..."
curl -X POST http://localhost:8080/books -d '{"title": "The Go Programming Language", "author": "Alan A. A. Donovan"}'
echo ""
curl -X POST http://localhost:8080/books -d '{"title": "Clean Code", "author": "Robert C. Martin"}'
echo ""
curl -X POST http://localhost:8080/books -d '{"title": "The Pragmatic Programmer", "author": "Andrew Hunt"}'
echo ""

# Get all books
echo "Listing books..."
curl http://localhost:8080/books
echo ""

# Get a specific book
echo "Getting book 1..."
curl http://localhost:8080/books/1
echo ""

# Get a non-existent book (to generate error trace)
echo "Getting book 999..."
curl http://localhost:8080/books/999
echo ""
