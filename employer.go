package workflow

import (
	"bytes"
	"fmt"
	"log"

	"golang.org/x/crypto/scrypt"
)

const (
	// PasswordSalt is the salt used in password creation / decryption
	PasswordSalt = "bobloblaw"
)

func encryptPassword(password string) []byte {
	e, _ := scrypt.Key([]byte(password), []byte(PasswordSalt), 16384, 8, 1, 32)
	return e
}

// AuthenticateEmployee checks to see if the given employee exists
func AuthenticateEmployee(id, password string) bool {
	if employee, ok := GetEmployee(id); ok {
		fmt.Printf("%v == %v\n", employee.Password, encryptPassword(password))
		return bytes.Equal(employee.Password, encryptPassword(password))
	}
	return false
}

// Employee is an employee of the employer (?!)
type Employee struct {
	ID             string
	Name           string
	YearlySalary   int    `db:"yearly_salary"`
	YearsOfService int    `db:"years_of_service"`
	Password       []byte `db:"password"`
	Position       string
}

// GetEmployee tries to retrieve an employee by id.
// Second argument is true if found.
func GetEmployee(id string) (Employee, bool) {
	db := getConnection()
	var employee Employee

	res := db.QueryRowx(`SELECT * FROM loblaw_employee WHERE id = $1`, id)
	if err := res.StructScan(&employee); err != nil {
		log.Println("[DEBUG] failed to get employee:", err)
		return employee, false
	}

	return employee, true
}

func resetEmployeeTables() {
	db := getConnection()
	db.MustExec(`DROP TABLE IF EXISTS loblaw_employee`)
	db.MustExec(`
    CREATE TABLE loblaw_employee (
        id               VARCHAR(255) PRIMARY KEY,
        name             TEXT NOT NULL,
		password         BYTEA NOT NULL,
        yearly_salary    INT NOT NULL,
        years_of_service INT NOT NULL,
        position         TEXT NOT NULL
    )`)

	pass := encryptPassword("bob")
	db.MustExec(`
    INSERT INTO loblaw_employee (id, name, yearly_salary,
                                 years_of_service, position, password)
                VALUES('bob', 'Bob Loblaw', 40000, 22, 'Founder and CEO', $1)
    `, pass)
}
