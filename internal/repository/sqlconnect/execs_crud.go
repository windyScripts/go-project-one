package sqlconnect

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"restapi/internal/models"
	"restapi/pkg/utils"
	"strconv"

	"golang.org/x/crypto/argon2"
)

func GetExecByID(id int) (models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Error retrieving data.")
	}
	defer db.Close()

	var exec models.Exec
	err = db.QueryRow("SELECT id, first_name, last_name, email, username, user_created_at, inactive_status, role FROM execs WHERE id = ?", id).Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.UserCreatedAt, &exec.InactiveStatus, &exec.Role)
	if err == sql.ErrNoRows {
		return models.Exec{}, utils.ErrorHandler(err, "Exec not found.")
	} else if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Error retrieving data.")
	}
	return exec, nil
}

func GetExecsDbHandler(execs []models.Exec, r *http.Request) ([]models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error retrieving data.")
	}
	defer db.Close()

	query := "SELECT id, first_name, last_name, email, username, user_created_at, inactive_status, role FROM execs WHERE 1=1"
	var args []interface{}

	query, args = utils.AddFilters(r, query, args)

	query = utils.AddSorting(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		return nil, utils.ErrorHandler(err, "Error retrieving data.")
	}
	defer rows.Close()

	// execList := make([]models.Exec, 0)

	for rows.Next() {
		var exec models.Exec
		err := rows.Scan(&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.UserCreatedAt, &exec.InactiveStatus, &exec.Role)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error retrieving data.")
		}
		execs = append(execs, exec)
	}
	return execs, nil
}

func AddExecsDbHandler(newExecs []models.Exec) ([]models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error adding data.")
	}
	defer db.Close()

	// stmt, err := db.Prepare("INSERT INTO execs (first_name, last_name, email, username) VALUES (?,?,?,?,?)")
	stmt, err := db.Prepare(utils.GenerateInsertQuery("execs", models.Exec{}))
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error adding data.")
	}
	defer stmt.Close()

	addedExecs := make([]models.Exec, len(newExecs))
	for i, newExec := range newExecs {
		//res, err := stmt.Exec(newExec.FirstName, newExec.LastName, newExec.Email, newExec.Class)
		if newExec.Password == "" {
			return nil, utils.ErrorHandler(errors.New("Password is blank"), "please enter password")
		}
		salt := make([]byte, 16)
		_, err = rand.Read(salt)
		if err != nil {
			return nil, utils.ErrorHandler(errors.New("Error in generating salt"), "error adding data")
		}

		hash := argon2.IDKey([]byte(newExec.Password), salt, 1, 64*1024, 4, 32)
		saltBase64 := base64.StdEncoding.EncodeToString(salt)
		hashBase64 := base64.StdEncoding.EncodeToString(hash)

		encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)
		newExec.Password = encodedHash

		values := utils.GetStructValues(newExec)
		res, err := stmt.Exec(values...)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error adding data.")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error adding data.")
		}
		newExec.ID = int(lastID)
		addedExecs[i] = newExec
	}
	return addedExecs, nil
}

func PatchExecsDB(updates []map[string]interface{}) error {
	db, err := ConnectDb()
	if err != nil {
		log.Println(err)
		return utils.ErrorHandler(err, "Error updating data.")
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return utils.ErrorHandler(err, "Error updating data.")
	}

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			// http.Error(w, "Invalid Exec ID", http.StatusBadRequest)
			return utils.ErrorHandler(err, "Error updating data.")
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			// http.Error(w, "Error converting ID to Int", http.StatusBadRequest)
			return err
		}

		var execFromDb models.Exec
		err = db.QueryRow("SELECT id, first_name, last_name, email, username FROM execs WHERE id = ?", id).Scan(&execFromDb.ID, &execFromDb.FirstName, &execFromDb.LastName, &execFromDb.Email, &execFromDb.Username)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				// http.Error(w, "Exec not found", http.StatusNotFound)
				return utils.ErrorHandler(err, "Exec not found.")
			}
			// http.Error(w, "Error retrieving exec", http.StatusInternalServerError)
			return utils.ErrorHandler(err, "Error updating data.")
		}

		// apply updates using reflection

		execVal := reflect.ValueOf(&execFromDb).Elem()
		execType := execVal.Type()

		for k, v := range update {
			if k == "id" {
				continue // skip updating id
			}
			for i := 0; i < execVal.NumField(); i++ {
				field := execType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := execVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							tx.Rollback()
							log.Printf("Cannot convert %v to %v", val.Type(), fieldVal.Type())
							return err
						}
					}
					break
				}
			}
		}
		_, err = tx.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ? WHERE id = ?", execFromDb.FirstName, execFromDb.LastName, execFromDb.Email, execFromDb.Username, execFromDb.ID)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Error updating data.")
		}
	}
	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return utils.ErrorHandler(err, "Error updating data.")
	}
	return nil
}

func PatchExecDB(id int, updates map[string]interface{}) (models.Exec, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Error updating data.")
	}

	defer db.Close()

	var existingExec models.Exec
	err = db.QueryRow("SELECT id, first_name, last_name, email, username FROM execs where id = ?", id).Scan(&existingExec.ID, &existingExec.FirstName, &existingExec.LastName, &existingExec.Email, &existingExec.Username)
	if err == sql.ErrNoRows {
		return models.Exec{}, utils.ErrorHandler(err, "Exec not found.")
	}
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Error updating data.")
	}

	execVal := reflect.ValueOf(&existingExec).Elem() // elements of object
	execType := execVal.Type()

	for k, v := range updates {
		for i := 0; i < execVal.NumField(); i++ {
			field := execType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if execVal.Field(i).CanSet() {
					fieldVal := execVal.Field(i)
					fieldVal.Set(reflect.ValueOf(v).Convert(execVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ? WHERE id = ?", existingExec.FirstName, existingExec.LastName, existingExec.Email, existingExec.Username, existingExec.ID)
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Error updating data.")
	}
	return existingExec, nil
}

func DeleteExecDB(id int) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.ErrorHandler(err, "Error deleting data.")
	}

	defer db.Close()

	result, err := db.Exec("DELETE FROM execs WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err, "Error deleting data.")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "Error deleting data.")
	}
	if rowsAffected == 0 {
		return utils.ErrorHandler(err, "Error deleting data.")
	}
	return nil
}

/*
authorization examples:
role based access control (rbac)
attribute based access control (abac)
access control lists (acl)s
decides what the authenticated user is allowed to do.
*/

/*
cookies: Key value pairs.
	small pieces of data usually client side, sent with each request.
	session management
		session id, preferences, tracking info
		redis and mongo are used for in memory database.
*/

/*
userid, claims, expiration time in JWT

*/
