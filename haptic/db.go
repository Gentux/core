package main

import (
	"encoding/json"
	"errors"

	nan "nanocloud.com/core/lib/libnan"

	"github.com/boltdb/bolt"
)

var ()

// We wrap the DB provider in a user struct to which we can add our own methods
type Db struct {
	*bolt.DB
}

type User struct {
	Activated    bool
	Email        string
	Firstname    string
	Lastname     string
	Password     string
	Sam          string
	CreationTime string
	Profile      string
}

func InitialiseDb() *nan.Err {

	var e error

	if nan.DryRun || nan.ModeRef {
		return nil
	}

	g_Db.DB, e = bolt.Open(nan.Config().Database.ConnectionString, 0777, nil)
	if e != nil {
		return LogErrorCode(ErrIssueWithAccountsDb)
	}

	e = g_Db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("users"))
		return nil
	})
	if e != nil {
		return LogErrorCode(ErrIssueWithAccountsDb)
	}

	var stats bolt.BucketStats
	e = g_Db.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("Bucket 'users' doesn't exist")
		}

		// If no user exists (new database) create first admin
		stats = bucket.Stats()
		return nil
	})
	if e != nil {
		return LogErrorCode(nan.ErrFrom(e))
	}

	if stats.KeyN == 0 {
		err := g_Db.AddUser(User{
			Activated:    true,
			Email:        nan.Config().AdminUser.Email,
			Firstname:    "admin",
			Lastname:     "admin",
			Password:     nan.Config().AdminUser.Password,
			Sam:          "",
			CreationTime: "",
			Profile:      "admin",
		})

		if err != nil {
			return LogErrorCode(err)
		}
	}

	return nil
}

func ShutdownDb() {
	if nan.DryRun || nan.ModeRef {
		return
	}

	defer g_Db.Close()
}

// ============================================================================================================
//
// DB utility functions
//
// ============================================================================================================

// Function:
func (p Db) GetUser(Email string) (User, *nan.Err) {
	var user User

	e := g_Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("Bucket 'users' doesn't exist")
		}

		userJson := bucket.Get([]byte(Email))
		json.Unmarshal(userJson, &user)

		return nil
	})

	return user, nan.ErrFrom(e)
}

func (p Db) GetUsers() ([]User, *nan.Err) {
	var (
		user  User
		users []User
	)

	e := g_Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("Bucket 'users' doesn't exist")
		}

		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			user = User{}
			json.Unmarshal(value, &user)
			users = append(users, user)
		}

		return nil
	})

	return users, nan.ErrFrom(e)
}

func (p Db) AddUser(user User) *nan.Err {
	e := g_Db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("Bucket 'users' doesn't exist")
		}

		jsonUser, e := json.Marshal(user)
		bucket.Put([]byte(user.Email), jsonUser)

		return e
	})

	return nan.ErrFrom(e)
}

func (p Db) DeleteUser(user User) *nan.Err {
	e := g_Db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("Bucket 'users' doesn't exist")
		}

		if user.Profile == "admin" {
			return errors.New("Admin users can't be deleted")
		}

		return bucket.Delete([]byte(user.Email))
	})

	return nan.ErrFrom(e)
}

func (p Db) IsUserRegistered(Email string) (bool, *nan.Err) {
	var result bool = false

	e := g_Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("Bucket 'users' doesn't exist")
		}

		cursor := bucket.Cursor()
		for key, _ := cursor.First(); key != nil; key, _ = cursor.Next() {
			if string(key) == Email {
				result = true
				break
			}
		}

		return nil
	})

	return result, nan.ErrFrom(e)
}

func (p Db) GetSamFromEmail(Email string) (string, *nan.Err) {
	var (
		user User
		sam  string
	)

	e := g_Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("Bucket 'users' doesn't exist")
		}

		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			if string(key) == Email {
				json.Unmarshal(value, &user)
				sam = user.Sam
				break
			}
		}

		return nil
	})

	return sam, nan.ErrFrom(e)
}

func (p Db) CountRegisteredUsers() (int, *nan.Err) {
	var count int = 0

	e := g_Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("Bucket 'users' doesn't exist")
		}

		cursor := bucket.Cursor()
		for key, _ := cursor.First(); key != nil; key, _ = cursor.Next() {
			count += 1
		}

		return nil
	})

	return count, nan.ErrFrom(e)
}

func (p Db) CountActiveUsers() (int, *nan.Err) {
	var (
		user  User
		count int = 0
	)

	e := g_Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("Bucket 'users' doesn't exist")
		}

		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			json.Unmarshal(value, &user)
			if user.Activated {
				count += 1
			}
		}

		return nil
	})

	return count, nan.ErrFrom(e)
}

func (p Db) IsUserActivated(Email string) (bool, *nan.Err) {
	var (
		user   User
		result bool = false
	)

	e := g_Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("Bucket 'users' doesn't exist")
		}

		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			if string(key) == Email {
				json.Unmarshal(value, &user)
				if user.Activated {
					result = true
				}
				break
			}
		}

		return nil
	})

	return result, nan.ErrFrom(e)
}

func (p Db) UpdateUserSamForEmail(Email, Sam string) bool {
	var (
		user   User
		result bool = false
	)

	g_Db.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("Bucket 'users' doesn't exist")
		}

		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			if string(key) == Email {
				json.Unmarshal(value, &user)
				user.Sam = Sam
				jsonUser, _ := json.Marshal(user)
				bucket.Put([]byte(user.Email), jsonUser)
				result = true
				break
			}
		}

		return nil
	})

	return result
}

func (p Db) GetRegisteredUsersInfo(pResults *[]RegisteredUserInfo) *nan.Err {
	var user User

	*pResults = []RegisteredUserInfo{}

	e := g_Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("DB bucket 'users' doesn't exist")
		}

		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			json.Unmarshal(value, &user)
			*pResults = append(*pResults, RegisteredUserInfo{
				Email:        user.Email,
				CreationTime: user.CreationTime,
				Activated:    user.Activated,
			})
		}

		return nil
	})

	return nan.ErrFrom(e)
}

func (p Db) GetActivatedUsersInfo(pResults *[]ActiveTacUserInfo) *nan.Err {
	var user User

	*pResults = []ActiveTacUserInfo{}

	e := g_Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("Bucket 'users' doesn't exist")
		}

		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			json.Unmarshal(value, &user)
			if user.Activated {
				*pResults = append(*pResults, ActiveTacUserInfo{
					TacId:        user.Email,
					TacUrl:       "TODO",
					CreationTime: user.CreationTime,
				})
			}
		}

		return nil
	})

	return nan.ErrFrom(e)
}

func (p Db) GetStats() ([]Stat, *nan.Err) {
	var (
		stat  Stat
		stats []Stat
	)

	e := g_Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("stats"))
		if bucket == nil {
			return errors.New("Bucket 'stats' doesn't exist")
		}

		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			stat = Stat{}
			json.Unmarshal(value, &stat)
			stats = append(stats, stat)
		}

		return nil
	})

	return stats, nan.ErrFrom(e)
}
