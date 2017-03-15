package workflow

// ResetTables deletes all data and recreates the database schema
func ResetTables() {
	resetMortgageTables()
	resetEmployeeTables()
}
