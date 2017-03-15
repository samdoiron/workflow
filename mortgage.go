package workflow

import (
	"database/sql"
	"errors"
	"log"
)

// SubmitMortgageApplication validates the given application, and stores it if
// it is valid.
func SubmitMortgageApplication(app MortgageApplication) (int, error) {
	if !app.IsValid() {
		return 0, errors.New("invalid application")
	}

	return app.submit()
}

// GetMortgageApplication retrieves a mortgage application by name.
// If not present, second return is false
func GetMortgageApplication(name string) (MortgageApplication, bool) {
	db := getConnection()
	row := db.QueryRowx("SELECT * FROM mortgageco_application WHERE name = $1",
		name)
	var app MortgageApplication
	if err := row.StructScan(&app); err != nil {
		log.Println("[DEBUG] failed to scan application", err)
		return app, false
	}
	return app, true
}

// SetEmployerInfo sets the information from an application's employer in
// the database.
func SetEmployerInfo(id int64, info MortgageEmployerInfo) error {
	db := getConnection()
	_ = db.MustExec(`
    UPDATE mortgageco_application
        SET yearly_salary = $2,
            years_of_service = $3,
            position = $4
    WHERE id = $1
    `, id, info.YearlySalary, info.YearsOfService, info.Position)
	return nil
}

// MortgageApplication represents a single application for a mortgage.
// Data is not assumed to have been validated.
type MortgageApplication struct {
	ID                int
	Name              string
	Phone             string
	Address           string
	EmployerName      string        `db:"employer_name"`
	LifeInsuranceName string        `db:"life_insurance_name"`
	YearlySalary      sql.NullInt64 `db:"yearly_salary"`
	YearsOfService    sql.NullInt64 `db:"years_of_service"`
	Position          sql.NullString
}

// MortgageEmployerInfo contains all information that mortgage application's
// care about from the person's employer.
type MortgageEmployerInfo struct {
	Name           string
	YearlySalary   int `json:"yearly_salary"`
	YearsOfService int `json:"years_of_service"`
	Position       string
}

// IsValid checks if the application is valid.
func (m *MortgageApplication) IsValid() bool {
	return !hasPreviousApplication(m.Name) &&
		m.Name != "" &&
		m.Phone != "" &&
		m.Address != "" &&
		m.EmployerName != "" &&
		m.LifeInsuranceName != ""
}

func (m *MortgageApplication) submit() (int, error) {
	db := getConnection()
	stmt, err := db.Prepare(`
    INSERT INTO mortgageco_application
        (name, phone, address, employer_name, life_insurance_name)
        VALUES ($1, $2, $3, $4, $5)
    RETURNING id`)
	if err != nil {
		return 0, errors.New("could not prepare statement")
	}

	var id int
	stmt.
		QueryRow(m.Name, m.Phone, m.Address, m.EmployerName, m.LifeInsuranceName).
		Scan(&id)

	return id, nil
}

func hasPreviousApplication(name string) bool {
	db := getConnection()
	var exists bool
	res := db.QueryRow(`
    SElECT EXISTS(
        SELECT 1 FROM mortgageco_application WHERE name = $1
    )`, name)

	res.Scan(&exists)
	return exists
}

func resetMortgageTables() {
	db := getConnection()
	db.MustExec(`DROP TABLE IF EXISTS mortgageco_application`)
	db.MustExec(`
    CREATE TABLE mortgageco_application (
        id                  SERIAL PRIMARY KEY,
        name                TEXT NOT NULL,
        phone               TEXT NOT NULL,
        address             TEXT NOT NULL,
        employer_name       TEXT NOT NULL,
        life_insurance_name TEXT NOT NULL,
        yearly_salary       INT,
        years_of_service    INT,
        position            TEXT
    )`)
	db.MustExec(`CREATE INDEX ON mortgageco_application (name)`)
	db.MustExec(`CREATE INDEX ON mortgageco_application (employer_name)`)
	db.MustExec(`CREATE INDEX ON mortgageco_application (life_insurance_name)`)
}
