package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	ID                 string
	Name               string
	LastName           string
	JobTitle           string
	SalaryByDay        int
	WorkedDays         int
	ExpectedWorkedDays int
	DiffWorked         int // If diff > 3, must be fired.
	MustFire           bool
}

// checkErr checks if there is an error and panics if it is the case.
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// GET all employees
func showEmployees(c *gin.Context) {
	var employees []Employee
	var employee Employee

	// MYSQL CONNECTION
	db, err := sql.Open("mysql", "root:Pepperonipizza123.@tcp(127.0.0.1:3306)/")
	checkErr(err)
	defer db.Close()

	_, err = db.Exec("USE employee_time_management")
	checkErr(err)

	// We execute a query iterating over all the rows.
	rows, err := db.Query("SELECT * FROM employees")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&employee.ID, &employee.Name, &employee.LastName, &employee.JobTitle, &employee.SalaryByDay, &employee.WorkedDays, &employee.ExpectedWorkedDays, &employee.DiffWorked, &employee.MustFire)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error. No employees found/bad request."})
		}

		employees = append(employees, employee)
	}

	// If everything goes well, we return Status code 200, and the list formatted in JSON as response.
	c.IndentedJSON(http.StatusOK, employees)
}

// GET a employee by his ID.
func employeeByID(c *gin.Context) {
	id := c.Param("id")

	var e Employee
	var es []Employee

	// MYSQL CONNECTION
	db, err := sql.Open("mysql", "root:Pepperonipizza123.@tcp(127.0.0.1:3306)/")
	checkErr(err)
	defer db.Close()

	_, err = db.Exec("USE employee_time_management")
	checkErr(err)

	// We execute a query iterating over all the rows.
	rows, err := db.Query("SELECT * FROM employees")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&e.ID, &e.Name, &e.LastName, &e.JobTitle, &e.SalaryByDay, &e.WorkedDays, &e.ExpectedWorkedDays, &e.DiffWorked, &e.MustFire)
		checkErr(err)
		es = append(es, e)
	}

	for _, employee := range es {
		if employee.ID == id {
			c.IndentedJSON(http.StatusOK, employee)
			return
		}
	}

	// If not found, returns a BadRequest Status code,
	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Employee not found."})
}

type userEmployee struct {
	Name               string
	LastName           string
	JobTitle           string
	SalaryByDay        int
	WorkedDays         int
	ExpectedWorkedDays int
}

// POST new employee in the company.
func newEmployee(c *gin.Context) {

	ue := userEmployee{}

	// MYSQL CONNECTION
	db, err := sql.Open("mysql", "root:Pepperonipizza123.@tcp(127.0.0.1:3306)/")
	checkErr(err)
	defer db.Close()

	_, err = db.Exec("USE employee_time_management")
	checkErr(err)

	// We send the reference to the variable and bind to JSON.
	if err := c.BindJSON(&ue); err != nil {
		return
	}

	var MustFire bool

	// Before adding to table, we calculate if the employee must be fired
	DiffWorked := ue.ExpectedWorkedDays - ue.WorkedDays
	if DiffWorked > 3 {
		MustFire = true
	} else {
		MustFire = false
	}

	// We prepare for the input in order to prevent SQL injections.
	prepareSentence, err := db.Prepare("INSERT INTO employees (name, last_name, job_title, salary_day, worked_days, expected_worked_days, diff_worked, must_fire) VALUES(?,?,?,?,?,?,?,?)")
	checkErr(err)

	defer prepareSentence.Close()

	// Execute sentence for every '?'
	_, err = prepareSentence.Exec(ue.Name, ue.LastName, ue.JobTitle, ue.SalaryByDay, ue.WorkedDays, ue.ExpectedWorkedDays, DiffWorked, MustFire)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error in request"})
	} else {
		c.IndentedJSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("New employee, %s %s, added to the database", ue.Name, ue.LastName)})
	}

}

type EmployeeMustBeFired struct {
	ID       string
	Name     string
	LastName string
	JobTitle string
}

// GET employees that must be fired by hours worked.
// If the difference between and expected worked hours are more than 3, must be fire.
func showMustBeFired(c *gin.Context) {
	var employeesToFire []EmployeeMustBeFired
	var employeeToFire EmployeeMustBeFired

	// MYSQL CONNECTION
	db, err := sql.Open("mysql", "root:Pepperonipizza123.@tcp(127.0.0.1:3306)/")
	checkErr(err)
	defer db.Close()

	_, err = db.Exec("USE employee_time_management")
	checkErr(err)

	// We execute a query iterating over all the rows.
	rows, err := db.Query("SELECT ID, name, last_name, job_title FROM employees WHERE must_fire = true;")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&employeeToFire.ID, &employeeToFire.Name, &employeeToFire.LastName, &employeeToFire.JobTitle)
		if err != nil {
			fmt.Println(err)
		}
		employeesToFire = append(employeesToFire, employeeToFire)
	}

	if len(employeesToFire) < 1 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error. No employees to fire found/bad request."})
	} else {
		// If everything goes well, we return Status code 200, and the list formatted in JSON as response.
		c.IndentedJSON(http.StatusOK, employeesToFire)
	}

}

// fireAll retrieves all the employees in condition and then deletes all from the database.
func fireAll(c *gin.Context) {
	var employeesToFire []EmployeeMustBeFired
	var employeeToFire EmployeeMustBeFired

	db, err := sql.Open("mysql", "root:Pepperonipizza123.@tcp(127.0.0.1:3306)/")
	checkErr(err)

	defer db.Close()

	_, err = db.Exec("USE employee_time_management")
	checkErr(err)

	// We execute a query iterating over all the rows.
	rows, err := db.Query("SELECT ID, name, last_name, job_title FROM employees WHERE must_fire = true;")
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&employeeToFire.ID, &employeeToFire.Name, &employeeToFire.LastName, &employeeToFire.JobTitle)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error. No employees to fire found/bad request."})
		}
		employeesToFire = append(employeesToFire, employeeToFire)
	}

	if len(employeesToFire) == 0 {
		c.IndentedJSON(http.StatusNoContent, gin.H{"Error": "No employees to fire."})
	}
	db.Exec("DELETE FROM employees WHERE must_fire = true")
	checkErr(err)

	// If everything goes well, the employees are deleted from the DB.
	c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%d employees were fired.", len(employeesToFire))})
}

func main() {
	// MYSQL DATABASE AND TABLE CREATION/CONNECTION.
	db, err := sql.Open("mysql", "root:Pepperonipizza123.@tcp(127.0.0.1:3306)/")
	checkErr(err)

	defer db.Close()

	_, err = db.Exec("CREATE DATABASE if not exists employee_time_management")
	checkErr(err)

	_, err = db.Exec("USE employee_time_management")
	checkErr(err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS employees (ID tinyint NOT NULL auto_increment, name tinytext NOT NULL, last_name tinytext NOT NULL, job_title tinytext NOT NULL, salary_day int NOT NULL, worked_days int not null, expected_worked_days int not null, diff_worked int not null, must_fire bool not null ,primary key (ID))")
	checkErr(err)

	r := gin.Default()

	// Search employees
	r.GET("/", showEmployees)
	r.GET("/employee/:id", employeeByID)

	// Create new employee
	r.POST("/newemployee", newEmployee)

	// GET and DELETE employees to be fired.
	r.GET("/mustbefired", showMustBeFired)
	r.DELETE("/mustbefired", fireAll)

	r.Run("localhost:8080")

}
