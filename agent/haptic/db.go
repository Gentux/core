package main

import (
	"fmt"

	"database/sql"

	nan "nanocloud.com/zeroinstall/lib/libnan"
)

var ()

func InitialiseDb() {

	if nan.DryRun || nan.ModeRef {
		return
	}

	pDb, err := sql.Open(nan.Config().Database.Type, nan.Config().Database.ConnectionString)
	g_Db = &Db{pDb}

	if err != nil || g_Db == nil {
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
func (p Db) IsUserRegistered(_sEmail string) (bool, error) {
	var err error = nil
	count := 0

	sRequest := fmt.Sprintf("SELECT COUNT(*) FROM guacamole_user WHERE username='%s'", _sEmail)

	rows, err := p.Query(sRequest)
	if err != nil {
		return false, err
	}

	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return false, err
		}
	}

	return count > 0, err
}

// func SAM=$(echo -e "$(eval "echo -e \"`<$MYPATH/sql/sam.sql`\"")" |
// mysql -u root --password='password' --database='guacamole' --skip-column-names)
func (p Db) GetSamFromEmail(_Email string) (string, error) {
	var err error = nil
	var sam string

	sRequest := fmt.Sprintf(`SELECT parameter_value FROM guacamole_connection_parameter 
INNER JOIN guacamole_connection_permission ON guacamole_connection_parameter.connection_id = guacamole_connection_permission.connection_id 
INNER JOIN guacamole_user ON guacamole_connection_permission.user_id = guacamole_user.user_id 
WHERE guacamole_connection_parameter.parameter_name='username'
AND guacamole_user.username='%s'
LIMIT 1;`, _Email)

	rows, err := p.Query(sRequest)
	if err != nil {
		return "", err
	}

	if rows.Next() == false {
		return "", nil
	}

	err = rows.Scan(&sam)

	return sam, err
}

func (p Db) CountRecordsInTable(_tableName string) (int, error) {

	var err error = nil
	count := 0

	sQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", _tableName)

	rows, err := p.Query(sQuery)
	if err != nil {
		return 0, err
	}

	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}

	return count, err
}

func (p Db) CountRegisteredUsers() (int, error) {
	if numRegs, e := p.CountRecordsInTable("guacamole_user"); e != nil {
		return 0, e
	} else {
		if numRegs > 0 {
			// Minus the guacamole admin account
			numRegs--
		}
		return numRegs, e
	}
}

func (p Db) CountActiveUsers() (int, error) {
	return p.CountRecordsInTable("talend_tac")
}

func (p Db) GetConnectionIdForEmail(Email string) string {
	connId := ""

	sQuery := fmt.Sprintf(`select connection_id from guacamole_connection where guacamole_connection.parent_id IN 
		(select connection_group_id from guacamole_connection_group where connection_group_name='%s')`, Email)

	rows, err := g_Db.Query(sQuery)
	if err != nil {
		LogError("Failed to select connection_id for email: %s, with error: %s", G_Account.Email, err)
		return ""
	}

	if rows.Next() == false {
		return ""
	}

	if err = rows.Scan(&connId); err != nil {
		LogError("Failed to parse result of query to get connection_id for email: %s, error: %s", G_Account.Email, err)
		return ""
	}

	return connId
}

// func (p Db) IsUserActivated(_sEmail string) (bool, error) {
// 	var err error = nil
// 	count := 0

// 	sRequest := fmt.Sprintf("SELECT COUNT(*) FROM talend_tac WHERE user_id='%s'", _sEmail)

// 	rows, err := p.Query(sRequest)
// 	if err != nil {
// 		return false, err
// 	}

// 	for rows.Next() {
// 		err = rows.Scan(&count)
// 		if err != nil {
// 			return false, err
// 		}
// 	}

// 	// !!! TODO LogError if count > 0 !!!

// 	return count > 0, err
// }

func (p Db) IsUserActivated(Email string) (bool, error) {
	sQuery := fmt.Sprintf(`SELECT parameter_value from guacamole_connection_parameter where parameter_name='username' AND guacamole_connection_parameter.connection_id 
		IN (select connection_id from guacamole_connection where guacamole_connection.parent_id 
			IN (select connection_group_id from guacamole_connection_group where connection_group_name='%s'))`, G_Account.Email)

	accountUserName := ""
	rows, err := g_Db.Query(sQuery)

	if err != nil {
		LogError("SQL error when selecting parameter_value for account email <%s> in SQL request, error: %s", G_Account.Email, err)
		return false, err
	} else if rows.Next() == false {
		LogError("Could not select a connection_parameter user_name associated to email : <%s>", G_Account.Email)
		return false, nil
	} else if err = rows.Scan(&accountUserName); err != nil {
		LogError("Failed to parse result of select on connection_parameter for <%s>", G_Account.Email)
		return false, err
	}

	return (accountUserName != "unactivated"), nil
}

func (p Db) UpdateConnectionUserNameForEmail(Email, UserName string) bool {

	sRequest := fmt.Sprintf(`update guacamole_connection_parameter set parameter_value='%s' 
		where parameter_name='username' AND guacamole_connection_parameter.connection_id 
			IN (select connection_id from guacamole_connection where guacamole_connection.parent_id 
				IN (select connection_group_id from guacamole_connection_group where connection_group_name='%s'));`, UserName, Email)

	if _, err := g_Db.Exec(sRequest); err != nil {
		LogError("Failed to update username in guacamole_connection_parameter for email: %s, with error: %s ", Email, err)
		return false
	}

	return true
}

func (p Db) GetRegisteredUsersInfo(pResults *[]RegisteredUserInfo) {
	var registeredUserInfo RegisteredUserInfo

	*pResults = []RegisteredUserInfo{}

	// 1) List registered users that are not activated

	sQuery := fmt.Sprintf(`SELECT connection_group_name FROM guacamole_connection_group WHERE guacamole_connection_group.connection_group_id 
							IN (SELECT parent_id FROM guacamole_connection WHERE connection_id 
								IN (SELECT connection_id FROM guacamole_connection_parameter WHERE parameter_name='username' 
									AND parameter_value='unactivated'));`)

	rows, err := g_Db.Query(sQuery)
	if err != nil {
		LogError("GetActiveUsersInfo: failed to perform select, error: %s", err)
		ExitError(ErrIssueWithAccountsDb)
	}

	for rows.Next() != false {

		if err = rows.Scan(&registeredUserInfo.Email); err != nil {
			LogError("GetRegisteredUsersInfo E1: failed to parse row with error: %s", err)
			ExitError(ErrIssueWithAccountsDb)
		}

		registeredUserInfo.CreationTime = GetUserAccountRegistrationTime(registeredUserInfo.Email)
		registeredUserInfo.Activated = false

		*pResults = append(*pResults, registeredUserInfo)
	}

	// 2) List registered users that are activated

	sQuery = fmt.Sprintf(`SELECT connection_group_name FROM guacamole_connection_group WHERE guacamole_connection_group.connection_group_id 
							IN (SELECT parent_id FROM guacamole_connection WHERE connection_id 
								IN (SELECT connection_id FROM guacamole_connection_parameter WHERE parameter_name='username' 
									AND parameter_value!='unactivated'));`)

	rows, err = g_Db.Query(sQuery)
	if err != nil {
		LogError("GetActiveUsersInfo E2: failed to perform select, error: %s", err)
		ExitError(ErrIssueWithAccountsDb)
	}

	for rows.Next() != false {

		if err = rows.Scan(&registeredUserInfo.Email); err != nil {
			LogError("GetRegisteredUsersInfo: failed to parse row with error: %s", err)
			ExitError(ErrIssueWithAccountsDb)
		}

		registeredUserInfo.CreationTime = GetUserAccountRegistrationTime(registeredUserInfo.Email)
		registeredUserInfo.Activated = true

		*pResults = append(*pResults, registeredUserInfo)
	}

}

func (p Db) GetActivatedUsersInfo(pResults *[]ActiveTacUserInfo) {
	var tacUserInfo ActiveTacUserInfo

	*pResults = []ActiveTacUserInfo{}

	sQuery := fmt.Sprintf("SELECT tac_id, tac_url, created_date FROM talend_tac")

	rows, err := g_Db.Query(sQuery)
	if err != nil {
		LogError("GetActiveUsersInfo: failed to perform select, error: %s", err)
		ExitError(ErrIssueWithAccountsDb)
	}

	for rows.Next() != false {

		if err = rows.Scan(&tacUserInfo.TacId, &tacUserInfo.TacUrl, &tacUserInfo.CreationTime); err != nil {
			LogError("GetActivatedUsersInfo: failed to parse result of query on table talend_tac, error: %s", err)
			ExitError(ErrIssueWithAccountsDb)
		}

		*pResults = append(*pResults, tacUserInfo)
	}
}
