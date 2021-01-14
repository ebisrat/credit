package repo

import (
	"bloom/read"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Defines a credit record for a Person
type Tag struct {
	CreditTag string `db:"credit_tag"`
	UserID    int    `db:"errouser_id"`
}

func SetupDB(credits []read.Credit) {
	// first create database called credit
	db, err := sql.Open("mysql", "root:rootpass@tcp(127.0.0.1:3306)/credit")
	if err != nil {
		panic(err.Error())
	}
	query := `DROP TABLE IF EXISTS credit;`
	_, err = db.Exec(query)
	if err != nil {
		panic(err.Error())
	}
	query = `
	CREATE TABLE credit (
		id int unsigned not null primary key auto_increment,
		uuid VARCHAR(191) NOT NULL,
		name varchar(64) not null,
		social_security int not null);
`
	_, err = db.Exec(query)

	if err != nil {
		panic(err.Error())
	}
	query = `DROP TABLE IF EXISTS tag;`
	_, err = db.Exec(query)
	if err != nil {
		panic(err.Error())
	}
	query = `
		CREATE TABLE tag (
			user_id int unsigned not null,
			credit_tag int not null);
		`
	_, err = db.Exec(query)

	if err != nil {
		panic(err.Error())
	}
	// takes a long time to run
	go func(db *sql.DB, credits []read.Credit) {
		MigrateDB(db, credits)
		defer db.Close()
	}(db, credits)
}

// migrate all data from credit file into db
func MigrateDB(db *sql.DB, credits []read.Credit) {
	query := ""
	for i, credit := range credits {
		query = `insert into credit (uuid, name, social_security)
			VALUES(UUID(), ?, ?)`
		u64, err := strconv.ParseUint(strings.TrimSpace(credit.SocialSecurity), 10, 32)
		if err != nil {
			u64 = 0
			fmt.Println(err)
		}
		ssnUint := uint(u64)
		res, err := db.Exec(query, strings.TrimSpace(credit.Name), ssnUint)
		if err != nil {
			panic(err.Error())
		}
		userID, err := res.LastInsertId()

		fmt.Println("Adding record: ", i+1)
		// credits := strings.Split(credit.CreditTag, " ")
		tagsCount := len(credit.CreditTag) / 9

		// for _, elem := range credit.CreditTag {
		offset := 0
		for i := 0; i < tagsCount; i++ {
			tag := ""
			if i == 0 {
				// first tag may be either 8 or 9 digits long
				if credit.CreditTag[0] == '-' {
					tag = credit.CreditTag[0 : (i+1)*8]
					offset = 1
				} else {
					tag = credit.CreditTag[0 : (i+1)*8]
				}
			} else {
				tag = credit.CreditTag[(i*8)+offset : ((i+1)*8)+1]
				offset = 1
			}
			u64, err := strconv.ParseUint(strings.TrimSpace(tag), 10, 32)
			if err != nil {
				u64 = 0
				fmt.Println(err)
			}
			tagUint := uint(u64)
			query = `insert into tag (credit_tag, user_id)
				VALUES(?, ?)`
			_, err = db.Exec(query, tagUint, userID)
			if err != nil {
				panic(err.Error())
			}
		}
		// }
	}

}

func GetUserTagByID(userID uint) []map[string]interface{} {
	db, err := sql.Open("mysql", "root:rootpass@tcp(127.0.0.1:3306)/credit")
	if err != nil {
		panic(err.Error())
		return nil
	}
	query := `select tag.* from tag where user_id = ?`
	// c := new(Tag)
	rows, err := db.Query(query, userID)
	if err != nil {
		panic(err.Error())
		return nil
	}
	columns, err := rows.Columns()

	// for each database row / record, a map with the column names and row values is added to the allMaps slice
	var allMaps []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i, _ := range values {
			pointers[i] = &values[i]
		}
		err := rows.Scan(pointers...)
		if err != nil {
			panic(err.Error())
			return nil
		}
		resultMap := make(map[string]interface{})
		for i, val := range values {
			fmt.Printf("Adding key=%s val=%v\n", columns[i], val)
			resultMap[columns[i]] = val
		}
		allMaps = append(allMaps, resultMap)
	}
	// defer results.Close()
	// colNames, err := results.Columns()
	// if err != nil {
	// 	panic(err)
	// }
	// tags := make([]Tag, len(colNames))
	// // for results.Next() {
	// c := new(Tag)
	// err = results.Scan(tags...)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil
	// }

	// tags = append(tags, *c)
	// }

	return allMaps

}
