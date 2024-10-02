package storage

import (
	"database/sql"
	"fmt"

	"math/rand"

	_ "github.com/go-sql-driver/mysql"
)

const ID_LENGTH = 6

type MySQLDB struct {
	rdatabase *sql.DB
}

func NewMySQLDB(usename string, password string, address string, dbname string) (*MySQLDB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", usename, password, address, dbname))
	if err != nil {
		return nil, fmt.Errorf("couldn't open database: %w", err)
	}

	return &MySQLDB{
		rdatabase: db,
	}, nil
}

func (mdb MySQLDB) Save(requestEntry PicRequestEntry) error {
	query := "INSERT INTO Requests(RequestID,Email,RequestStatus,ImageCaption,NewImageURL) VALUES (?,?,?,?,?)"
	_, err := mdb.rdatabase.Exec(query,
		requestEntry.ReqId,
		requestEntry.Email,
		requestEntry.ReqStatus,
		requestEntry.ImageCaption,
		requestEntry.NewImageURL)
	if err != nil {
		return fmt.Errorf("couldn't save the entry in database: %w", err)
	}
	return nil
}

func (mdb MySQLDB) Get(requestId string) (PicRequestEntry, error) {
	query := "SELECT RequestID,Email,RequestStatus,ImageCaption,NewImageURL FROM Requests WHERE RequestID = ?"
	row := mdb.rdatabase.QueryRow(query, requestId)
	var newEntry PicRequestEntry

	err := row.Scan(
		&newEntry.ReqId,
		&newEntry.Email,
		&newEntry.ReqStatus,
		&newEntry.ImageCaption,
		&newEntry.NewImageURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return newEntry, RequestNotFoundError{ReqId: requestId}
		}
		return newEntry, err
	}

	return newEntry, nil
}

func (mdb MySQLDB) GenerateUniqueID() (string, error) {
	query := "SELECT RequestID FROM Requests WHERE RequestID = ?"
	for i := 0; i < 10; i++ {
		randID := generateAlphanumericSequence(ID_LENGTH)
		var empID int
		err := mdb.rdatabase.QueryRow(query, randID).Scan(&empID)
		if err != nil {
			if err == sql.ErrNoRows {
				return randID, nil
			}
			return "", fmt.Errorf("unable to query for generated ID: %w", err)
		}
	}
	return "", fmt.Errorf("no unique ID found after 10 tries")
}

func generateAlphanumericSequence(length int) string {
	charset := []byte("0123456789abcdefghijklmnopqrstuvwxyz")
	randseq := make([]byte, length)
	for i := range length {
		randseq[i] = charset[rand.Intn(len(charset))]
	}
	return string(randseq)
}
