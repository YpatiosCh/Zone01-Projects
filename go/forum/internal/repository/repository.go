package repository

import (
	"database/sql"
	"fmt"
	uuid "forum/pkg/UUID"
	"strings"
	"time"
)

type Manager struct {
	Db *sql.DB
}

func NewManager(db *sql.DB) *Manager {
	return &Manager{Db: db}
}

// CreateRecord inserts a new record into the specified table
// It generates a UUID for the ID if not provided and sets the creation time
// It returns the ID of the newly created record
// or an error if the operation fails
func (m *Manager) CreateRecord(table string, data map[string]interface{}) (string, error) {

	// Validate input
	if table == "" {
		return "", fmt.Errorf("table name cannot be empty")
	}
	if len(data) == 0 {
		return "", fmt.Errorf("no data provided for insertion")
	}

	// Start a transaction
	tx, err := m.Db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // This will rollback if not committed

	// Generate UUID for the ID if not provided
	id, ok := data["id"]
	if !ok || id == "" {
		id = uuid.GenerateUUID()
		data["id"] = id
	}

	// Set creation time if not provided and the table has a created_at column
	_, hasCreatedAt := data["created_at"]
	if !hasCreatedAt {
		// Use transaction to check if the table has a created_at column
		checkQuery := fmt.Sprintf("PRAGMA table_info(%s)", strings.ToLower(table))
		rows, err := tx.Query(checkQuery)
		if err != nil {
			return "", fmt.Errorf("failed to check table structure: %w", err)
		}

		hasCreatedAtColumn := false
		for rows.Next() {
			var cid, notnull, pk int
			var name, typeName string
			var dfltValue interface{}
			err = rows.Scan(&cid, &name, &typeName, &notnull, &dfltValue, &pk)
			if err != nil {
				rows.Close()
				return "", fmt.Errorf("error scanning column info: %w", err)
			}
			if strings.ToLower(name) == "created_at" {
				hasCreatedAtColumn = true
			}
		}
		rows.Close()

		if err = rows.Err(); err != nil {
			return "", fmt.Errorf("error iterating columns: %w", err)
		}

		if hasCreatedAtColumn {
			data["created_at"] = time.Now()
		}
	}

	// Build column names and placeholders
	var columns []string
	var placeholders []string
	var values []interface{}

	for col, val := range data {
		columns = append(columns, strings.ToLower(col))
		placeholders = append(placeholders, "?")
		values = append(values, val)
	}

	// Build and execute query using the transaction
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		strings.ToLower(table),
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	// Use the transaction for execution
	_, err = tx.Exec(query, values...)
	if err != nil {
		return "", fmt.Errorf("failed to insert into %s: %w", table, err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id.(string), nil
}

// UpdateRecord updates an existing record in the specified table
func (m *Manager) UpdateRecord(table string, id string, data map[string]interface{}) error {
	// Validate input
	if table == "" {
		return fmt.Errorf("table name cannot be empty")
	}
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}
	if len(data) == 0 {
		return fmt.Errorf("no data provided for update")
	}

	// Start a transaction
	tx, err := m.Db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // This will rollback if not committed

	// Remove id from data if present (we don't want to update the primary key)
	delete(data, "id")

	// Build SET clause
	var setClauses []string
	var values []interface{}

	for col, val := range data {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", strings.ToLower(col)))
		values = append(values, val)
	}

	// Add ID to values for the WHERE clause
	values = append(values, id)

	// Build and execute query
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?",
		strings.ToLower(table),
		strings.Join(setClauses, ", "))

	// Use transaction for execution
	result, err := tx.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to update %s: %w", table, err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no record found with id %s in table %s", id, table)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteRecord deletes a specific record from the specified table
func (m *Manager) DeleteRecord(table string, id string) error {
	// Validate input
	if table == "" {
		return fmt.Errorf("table name cannot be empty")
	}
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	// Start a transaction
	tx, err := m.Db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // This will rollback if not committed

	// Build and execute query
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", strings.ToLower(table))

	// Use transaction for execution
	result, err := tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete from %s: %w", table, err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no record found with id %s in table %s", id, table)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetAll retrieves all rows from the specified table
// It returns a slice of maps, where each map represents a row
// GetAll retrieves all rows from the specified table
// It returns a slice of maps, where each map represents a row
func (m *Manager) GetAll(table string) ([]map[string]interface{}, error) {
	// Start a transaction
	tx, err := m.Db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // This will rollback if not committed

	query := fmt.Sprintf("SELECT * FROM %s", strings.ToLower(table))

	// Execute the query using the transaction
	rows, err := tx.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query table %s: %w", table, err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns for table %s: %w", table, err)
	}

	// Create slice to hold all rows
	var result []map[string]interface{}

	// Create a slice of interface{} to hold each row
	values := make([]interface{}, len(columns))
	// Create a slice of pointers to those values
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Process each row
	for rows.Next() {
		// Scan the row into the values slice
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Create a map for this row
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]

			// Convert []byte to string for text values
			b, ok := val.([]byte)
			if ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}

		// Add this row to the result
		result = append(result, rowMap)
	}

	// Check for errors after iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}

// Get retrieves rows from a table where a specific column matches a value
// Get retrieves rows from a table where a specific column matches a value
func (m *Manager) Get(table string, column string, value interface{}) ([]map[string]interface{}, error) {
	// Start a transaction
	tx, err := m.Db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // This will rollback if not committed

	// Build the query
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?",
		strings.ToLower(table),
		strings.ToLower(column))

	// Execute the query using the transaction
	rows, err := tx.Query(query, value)
	if err != nil {
		return nil, fmt.Errorf("failed to query table %s: %w", table, err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns for table %s: %w", table, err)
	}

	// Create slice to hold all rows
	var result []map[string]interface{}

	// Create slices for scanning row values
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Process each row
	for rows.Next() {
		// Scan the row into the values slice
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Create a map for this row
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]

			// Convert []byte to string for text values
			b, ok := val.([]byte)
			if ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}

		// Add this row to the result
		result = append(result, rowMap)
	}

	// Check for errors after iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return nil if no results found to be consistent
	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

// CreateJunction creates a record in a junction table with a composite primary key
func (m *Manager) CreateJunction(table string, data map[string]interface{}) error {
	// Validate input
	if table == "" {
		return fmt.Errorf("table name cannot be empty")
	}
	if len(data) == 0 {
		return fmt.Errorf("no data provided for insertion")
	}

	// Start a transaction
	tx, err := m.Db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // This will rollback if not committed

	// Build column names and placeholders
	var columns []string
	var placeholders []string
	var values []interface{}

	for col, val := range data {
		columns = append(columns, strings.ToLower(col))
		placeholders = append(placeholders, "?")
		values = append(values, val)
	}

	// Build and execute query using the transaction
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		strings.ToLower(table),
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	// Use the transaction for execution
	_, err = tx.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to insert into %s: %w", table, err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
