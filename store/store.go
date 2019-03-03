package store

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/bseishen/golang_doorcontrol/api"
	"github.com/bseishen/golang_doorcontrol/user"
	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	dbFile string
	db     *sql.DB
}

func New(file string) *Store {
	return &Store{
		dbFile: file,
	}
}

func (s *Store) Open() error {
	db, err := sql.Open("sqlite3", s.dbFile)
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) FindUser(key int, pw string) (*user.User, error) {
	err := s.Open()
	if err != nil {
		return nil, err
	}
	defer s.Close()

	var (
		u          *user.User
		hashbuffer string
		isActive   int
		ircName    string
	)

	rows, err := s.db.Query("SELECT hash, active, irc_name from users WHERE key = ?", key)
	if err != nil {
		return nil, err
	}

	//key should be unique, if not grab the last record.
	for rows.Next() {
		err = rows.Scan(&hashbuffer, &isActive, &ircName)
		if err != nil {
			return nil, err
		}
		u = &user.User{DBHash: hashbuffer, Active: isActive, IrcName: ircName, Key: key}
	}

	if u == nil {
		return nil, errors.New("Access Denied: User not found")
	}

	if u.Active == 0 {
		return nil, errors.New("Access Denied: User is not Active")
	}

	if !u.ValidatePass(pw) {
		return nil, errors.New("Access DENIED: Incorrect keycode")
	}

	return u, nil
}

//Deletes existing SQLite DB then Writes member data to a new SQLite database
func (s *Store) UpdateDatabase(memberdata api.Data) {
	os.Remove(s.dbFile)

	db, err := sql.Open("sqlite3", s.dbFile)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table users (id integer not null primary key, key interger, hash text, irc_name text, spoken_name text,added_by interger, date_created text, last_login text, admin interger, active interger,user_id interger, created_at string, updated_at string);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into users(id, key, hash, irc_name, spoken_name, added_by, date_created, last_login, admin, active, user_id, created_at, updated_at) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for l := range memberdata.Members {
		_, err = stmt.Exec(memberdata.Members[l].Id, memberdata.Members[l].Key, memberdata.Members[l].Hash, memberdata.Members[l].Irc_name, memberdata.Members[l].Spoken_name, memberdata.Members[l].Added_by, memberdata.Members[l].Date_created, memberdata.Members[l].Last_login, memberdata.Members[l].Admin, memberdata.Members[l].Active, memberdata.Members[l].User_id, memberdata.Members[l].Created_at, memberdata.Members[l].Updated_at)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()

	log.Println("Database updated")

}
