package main

import (
	"encoding/json"
	"errors"

	nan "nanocloud.com/zeroinstall/lib/libnan"

	"github.com/boltdb/bolt"
)

var ()

type User struct {
	Activated    bool
	Email        string
	Firstname    string
	Lastname     string
	Password     string
	Sam          string
	CreationTime string
}

func InitialiseDb() {
	var err error

	if nan.DryRun || nan.ModeRef {
		return
	}

	g_Db.DB, err = bolt.Open(nan.Config().Database.ConnectionString, 0777, nil)
	if err != nil {
		ExitError(ErrIssueWithAccountsDb)
	}

	err = g_Db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("users"))
		return nil
	})
	if err != nil {
		ExitError(ErrIssueWithAccountsDb)
	}
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
func (p Db) GetUser(Email string) (User, error) {
	var user User

	e := g_Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("Bucket 'users' doesn't exist")
		}

		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			if string(key) == Email {
				json.Unmarshal(value, &user)
				break
			}
		}

		return nil
	})

	return user, e
}

func (p Db) GetUsers() ([]User, error) {
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

	return users, e
}

func (p Db) AddUser(user User) error {
	return g_Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("Bucket 'users' doesn't exist")
		}

		jsonUser, e := json.Marshal(user)
		bucket.Put([]byte(user.Email), jsonUser)

		return e
	})
}

func (p Db) DeleteUser(user User) error {
	return g_Db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("Bucket 'users' doesn't exist")
		}

		return bucket.Delete([]byte(user.Email))
	})
}

func (p Db) IsUserRegistered(Email string) (bool, error) {
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

	return result, e
}

func (p Db) GetSamFromEmail(Email string) (string, error) {
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

	return sam, e
}

func (p Db) CountRegisteredUsers() (int, error) {
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

	return count, e
}

func (p Db) CountActiveUsers() (int, error) {
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

	return count, e
}

func (p Db) IsUserActivated(Email string) (bool, error) {
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

	return result, e
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

func (p Db) GetRegisteredUsersInfo(pResults *[]RegisteredUserInfo) error {
	var user User

	*pResults = []RegisteredUserInfo{}

	return g_Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("Bucket 'users' doesn't exist")
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
}

func (p Db) GetActivatedUsersInfo(pResults *[]ActiveTacUserInfo) error {
	var user User

	*pResults = []ActiveTacUserInfo{}

	return g_Db.View(func(tx *bolt.Tx) error {
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
}
