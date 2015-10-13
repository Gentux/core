package main

import (
	"fmt"
	"path/filepath"

	nan "nanocloud.com/zeroinstall/lib/libnan"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"io/ioutil"

	"os"
	"os/exec"
	"strings"
)

// ========================================================================================================================

// We wrap sql.DB in a user struct to which we can add our own methods
type Db struct {
	*sql.DB
}

// ========================================================================================================================
// TYPES
// =====

type AccountParams struct {
	FirstName, LastName, Email, Password, License string
}

// ========================================================================================================================
// ValidateCreateUserParams
// ========================================================================================================================

// Procedure [NOSIDEEFFECT] Check preconditions for valid/compliant account creation parameters
func ValidateCreateUserParams() {

	Log("Verifying parameters to create %s account", G_Account.Email)

	if !nan.ValidName(G_Account.FirstName) {
		ExitError(ErrFirstnameNonCompliant)
	}

	if !nan.ValidName(G_Account.LastName) {
		ExitError(ErrLastnameNonCompliant)
	}

	if !nan.ValidEmail(G_Account.Email) {
		ExitError(nan.ErrPbWithEmailFormat)
	}

	if !nan.ValidPassword(G_Account.Password) {
		ExitError(nan.ErrPasswordNonCompliant)
	}

	if nan.DryRun || nan.ModeRef {
		return
	}

	//TODO OPTionalize this

	// Refuse creation if problem with license file
	// NOTE: license may be valid but http license service backend may be down
	// resp, httpErr := http.Get(G_Account.License)
	// defer resp.Body.Close()
	// if httpErr != nil || resp.StatusCode != 200 {
	// 	ExitError(ErrInvalidLicenseFile)
	// }
}

// ========================================================================================================================
// Procedure: CreateAccount
//
// Does:
// - Check Params
// - Request pooled TAC VM from TAC host
// - Register TAC user : insert record in db guacamode/talend_tac
// - CreateUser (LDAP, AD)
// ========================================================================================================================
func CreateAccount(p AccountParams) {

	G_Account = p

	InitialiseDb()
	defer ShutdownDb()

	ValidateCreateUserParams()

	if !nan.DryRun && !nan.ModeRef {

		// bActive, err := g_Db.IsUserActivated(G_Account.Email)
		// if err != nil {
		// 	ExitError(ErrIssueWithAccountsDb)
		// } else if bActive {
		// 	ExitError(ErrAccountExists)
		// }

		// if maxNumAccounts := nan.Config().Proxy.MaxNumAccounts; maxNumAccounts > 0 {

		// 	if nAccounts, err := g_Db.CountActiveUsers(); err != nil {
		// 		ExitError(ErrIssueWithAccountsDb)
		// 	} else if nAccounts >= maxNumAccounts {
		// 		ExitError(ErrMaxNumAccountsReached)
		// 	}
		// }

	}

	// Email <=> TAC uuid
	//defer UndoIfFailed(G_ProcCreateTac)
	//[OPT] Create account resource such as TAC VM
	//G_ProcCreateTac.Do()

	//defer UndoIfFailed(G_ProcCreateWinUser)
	G_ProcCreateWinUser.Fqdn = "n/a" // [OPT] G_ProcCreateTac.Ans.TacUrl
	G_ProcCreateWinUser.Do()

	G_ProcRegisterProxyUser.sam = G_ProcCreateWinUser.out.sam
	G_ProcRegisterProxyUser.Do()

	ExitOk(OkAccountBeingCreated)
}

// ========================================================================================================================
// Procedure: RegisterUser
//
// Does:
// - Check Params
// - Register TAC user : insert record in db guacamode/talend_tac
// ========================================================================================================================
func RegisterUser(p AccountParams) {

	G_Account = p

	InitialiseDb()
	//[OPT] was done for Talend single executables
	//defer ShutdownDb()

	ValidateCreateUserParams()

	Log("STARTING registerUser for: %s", G_Account.Email)

	// Activation fails if no license specified
	// TODO Make license checking configurable, as a plugin or optional prerequisite step
	if G_Account.License == "" {
		G_Account.License = "n/a"
	}

	numReg, err := g_Db.CountRegisteredUsers()
	if err != nil {
		ExitError(ErrIssueWithAccountsDb)
	} else if numReg >= nan.Config().Proxy.MaxNumRegistrations {
		ExitError(ErrMaxNumAccountsRegistered)
	}

	bRegistered, err := g_Db.IsUserRegistered(G_Account.Email)
	if err != nil {
		ExitError(ErrIssueWithAccountsDb)
	} else if bRegistered {
		ExitError(ErrAccountExists)
	}

	G_ProcRegisterProxyUser.sam = "unactivated"
	G_ProcRegisterProxyUser.Do()

	// Store user parameters on disk so that they can be recovered easily by the next call to activateUser

	if CreateUserWorkspaceDir() {
		sUserAccountParams := fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
			G_Account.FirstName, G_Account.LastName, G_Account.Email, G_Account.Password, G_Account.License)

		sUserParamsFilePath := fmt.Sprintf("%s/studio/%s/account_params", nan.Config().CommonBaseDir, G_Account.Email)

		if err := ioutil.WriteFile(sUserParamsFilePath, []byte(sUserAccountParams), 0777); err != nil {
			LogError("Failed to save user account params into workspace file : %s", sUserParamsFilePath)
		}
	}

	//TODO ExitOk : make this behaviour configurable : exit or just print
	nan.PrintOk(OkAccountBeingCreated)
}

// Returns FirstName LastName Email Password License
func GetUserAccountParamsForActivation(_Email string) (string, string, string, string, string) {

	// Reload user account params from workspace
	sUserParamsFilePath := fmt.Sprintf("%s/studio/%s/account_params", nan.Config().CommonBaseDir, _Email)

	var err error
	var bytesRead []byte

	if bytesRead, err = ioutil.ReadFile(sUserParamsFilePath); err != nil {
		nan.ExitErrorf(0, "Failed to read user account params from workspace file : %s", sUserParamsFilePath)
	}

	var FirstName, LastName, tmpEmail, Password, License string

	nItemsParsed, err := fmt.Sscanf(string(bytesRead), "%s\n%s\n%s\n%s\n%s",
		&FirstName, &LastName, &tmpEmail, &Password, &License)

	if err != nil || nItemsParsed != 5 {
		nan.ExitErrorf(0, "Failed to parse user account params from workspace file : %s, read %d items",
			sUserParamsFilePath, nItemsParsed)
	}

	return FirstName, LastName, tmpEmail, Password, License
}

// ========================================================================================================================
// Procedure: ActivateUser
//
// Does:
// - Check Params
// - Register TAC user : insert record in db guacamode/talend_tac
// ========================================================================================================================
func ActivateUser(p AccountParams) {

	G_Account = p

	//[OPT] was previously needed for single executalbles
	//InitialiseDb()
	//defer ShutdownDb()

	if !nan.ValidEmail(G_Account.Email) {
		ExitError(nan.ErrPbWithEmailFormat)
	}

	Log("STARTING activateUser for: %s", G_Account.Email)

	// Reached maximum number of active users ?

	if maxNumAccounts := nan.Config().Proxy.MaxNumAccounts; maxNumAccounts > 0 {

		if nAccounts, err := g_Db.CountActiveUsers(); err != nil {
			ExitError(ErrIssueWithAccountsDb)
		} else if nAccounts >= maxNumAccounts {
			ExitError(ErrMaxNumAccountsReached)
		}
	}

	// User not registered yet ?

	if bRegistered, err := g_Db.IsUserRegistered(G_Account.Email); err != nil {
		ExitError(ErrIssueWithAccountsDb)
	} else if !bRegistered {
		ExitError(ErrAccountNotRegistered)
	}

	// If user already activated then do nothing
	if bValue, _ := g_Db.IsUserActivated(G_Account.Email); bValue {
		ExitOk(ErrAccountActivated)
	}

	tmpEmail := ""
	G_Account.FirstName, G_Account.LastName, tmpEmail, G_Account.Password, G_Account.License =
		GetUserAccountParamsForActivation(G_Account.Email)

	if G_Account.Email != tmpEmail {
		LogError("[INCONSISTENCY] Email passed as parameter (%s) doesn't match email loaded in account params: %s",
			G_Account.Email, tmpEmail)
	}

	//defer UndoIfFailed(G_ProcCreateTac)
	//[OPT] Create account resource such as TAC VM
	// G_ProcCreateTac.Do()

	//defer UndoIfFailed(G_ProcCreateWinUser)
	G_ProcCreateWinUser.Fqdn = "n/a" //[OPT] = G_ProcCreateTac.Ans.TacUrl
	G_ProcCreateWinUser.Do()

	if G_ProcCreateWinUser.out.sam == "" {
		ExitError(ErrIssueWithTacProvisioning)
	}

	g_Db.UpdateConnectionUserNameForEmail(G_Account.Email, G_ProcCreateWinUser.out.sam)

	ExitOk(OkAccountBeingActivated)
}

// ====================================================================================================
// Procedure: CreateUser
// => LDAP, Win user profile, security,
// ====================================================================================================

type ProcCreateWinUser struct {
	nan.ProcedureStruct

	Fqdn string // currently set to the TAC url as received in JSON from the TAC request

	out ProcCreateWinUserOutput
}

type ProcCreateWinUserOutput struct {
	sam string
}

func CreateUserWorkspaceDir() bool {

	sWorkspaceDirPath := fmt.Sprintf("%s/studio/%s", nan.Config().CommonBaseDir, G_Account.Email)

	if _, err := os.Stat(sWorkspaceDirPath); err != nil {
		if err = os.MkdirAll(sWorkspaceDirPath, os.ModePerm); err != nil {
			LogError("Failed to create directory: %s", err)
			return false
		}

		if err = os.Chmod(sWorkspaceDirPath, 0777); err != nil {
			LogError("Failed to set permissions on directory: %s", err)
			return false
		}
	}

	return true
}

func TestWinExe(sam string) {

	sWindowsServerSecurityFile := fmt.Sprintf(`\\winad.intra.nanocloud.com\NETLOGON\%s.setSecurity.bat`, sam)
	nan.Debug("Invoking WinExe on file:" + sWindowsServerSecurityFile)
	cmd := exec.Command(nan.Config().Proxy.WinExe,
		"-U", "intra.nanocloud.com/Administrator%3password",
		"--runas=intra.nanocloud.com/Administrator%3password", "//10.20.12.10", //TODOHARDCODED
		sWindowsServerSecurityFile)

	if out, err := cmd.Output(); err != nil {
		Log("Error <%s> returned by WinExe when running file <%s> with output: <%s>", err, sWindowsServerSecurityFile, out)
		ExitError(nan.ErrSomethingWrong)
	}
}

func (p *ProcCreateWinUser) Do() {

	var err error

	if nan.DryRun || nan.ModeRef {
		Log("Creating Windows user + LDAP declaration")
		return
	}

	// TODO: Remove these check, redundant with earlier check ?
	if !nan.ValidEmail(G_Account.Email) {
		ExitError(nan.ErrPbWithEmailFormat)
	}

	if !nan.ValidPassword(G_Account.Password) {
		ExitError(nan.ErrPasswordNonCompliant)
	}

	if !G_TwoStageActivation {
		bRegistered, err := g_Db.IsUserRegistered(G_Account.Email)
		if err != nil {
			ExitError(ErrIssueWithAccountsDb)
		} else if bRegistered {
			ExitError(ErrAccountExists)
		}
	}

	// If this account is still enabled on the LDAP, deactivate it

	resp := ""

	params := fmt.Sprintf(`{ "userid" : "%s" }`, G_Account.Email)

	g_PluginLdap.Call("Ldap.ForceDisableAccount", params, &resp)

	if resp == "0" {
		LogError("Ldap.ForceDisableAccount failed")
	}

	// Active Directory user
	// =====================

	Log("Configure Windows user profile")

	// Add LDAP user

	params = fmt.Sprintf(`{ "userid" : "%s", "password" : "%s" }`, G_Account.Email, G_Account.Password)

	for idx := 0; idx < 3; idx++ {
		g_PluginLdap.Call("Ldap.AddUser", params, &resp)

		if resp[0] != '$' {
			LogError("Failed to add LDAP user, got output: <%s>. Retrying for user <%s> and password <%s>", resp, G_Account.Email, G_Account.Password)
		} else {
			break
		}
	}

	p.out.sam = strings.Trim(resp, " ")

	// sAddUserPhpScript := fmt.Sprintf("%s/add_LDAP_user.php", nan.Config().CommonBaseDir)

	// for idx := 0; idx < 3; idx++ {
	// 	cmd := exec.Command("/usr/bin/php", "-f", sAddUserPhpScript, "--", G_Account.Email, G_Account.Password,
	// 		/* necessary ? */ "2>/dev/null")

	// 	out, err := cmd.Output()
	// 	if err != nil {
	// 		LogError("Failed to run script add_LDAP_user.php for email <%s> and password <%s>, error: %s", G_Account.Email, G_Account.Password, err)
	// 		ExitError(nan.ErrSomethingWrong)
	// 	}

	// 	p.out.sam = string(out)

	// 	if p.out.sam[0] != '$' {
	// 		LogError("Failed to add LDAP user, got output: <%s>. Retrying for user <%s> and password <%s>", p.out.sam, G_Account.Email, G_Account.Password)
	// 	} else {
	// 		break
	// 	}
	//}

	// Setup and finally upload the Workspace script on the Windows server
	// ===================================================================

	CreateUserWorkspaceDir()

	// sWorkspaceDirPath := fmt.Sprintf("%s/studio/%s", Config().CommonBaseDir, G_Account.Email)
	// if err = os.MkdirAll(sWorkspaceDirPath, os.ModePerm); err != nil {
	// 	LogError("Failed to create directory: %s", err)
	// 	ExitError(ErrFilesystemError)
	// }

	// if err = os.Chmod(sWorkspaceDirPath, 0777); err != nil {
	// 	LogError("Failed to set permissions on directory: %s", err)
	// 	ExitError(ErrSystemError)
	// }

	// TODO make this configurable in the case of Talend

	Log("Configuring workspace for %s", G_Account.Email)

	workspaceScriptFilePath := fmt.Sprintf("%s/studio/workspace.cmd", nan.Config().CommonBaseDir)
	samFilePath := fmt.Sprintf("%s/studio/%s/%s.config.bat", nan.Config().CommonBaseDir, G_Account.Email, p.out.sam)

	if err := nan.CopyFile(workspaceScriptFilePath, samFilePath); err != nil {
		LogError(nan.ErrFilesystemError.Message)
	}

	if err := nan.ReplaceInFile(samFilePath, "USEREMAIL", G_Account.Email); err != nil {
		LogError(nan.ErrFilesystemError.Message)
	}

	if err := nan.ReplaceInFile(samFilePath, "USERPASSWORD", G_Account.Password); err != nil {
		LogError(nan.ErrFilesystemError.Message)
	}

	if err := nan.ReplaceInFile(samFilePath, "USERFQDN", p.Fqdn); err != nil {
		LogError(nan.ErrFilesystemError.Message)
	}

	if err := exec.Command("/usr/bin/unix2dos", samFilePath).Run(); err != nil {
		Log("Error when attempting to use unix2dos on file : %s, error: %s", samFilePath, err)
	}

	// TODO add timeout + retry on this call + message: did not respond in a timely manner

	if err = exec.Command("/usr/bin/scp", samFilePath, "Administrator@10.20.12.20:/cygdrive/c/Windows/SYSVOL/domain/scripts/").Run(); err != nil {
		Log("Error when attempting to scp file: %s on server", samFilePath)
		ExitError(nan.ErrSomethingWrong)
	}

	// Setup, upload and execute the security script on Windows Server
	// ===============================================================

	sSecurityScriptPath := fmt.Sprintf("%s/studio/setSecurity.cmd", nan.Config().CommonBaseDir)
	sUserSecurityScriptPath := fmt.Sprintf("%s/studio/%s/%s.setSecurity.bat", nan.Config().CommonBaseDir, G_Account.Email, p.out.sam)

	if err := nan.CopyFile(sSecurityScriptPath, sUserSecurityScriptPath); err != nil {
		LogError(nan.ErrFilesystemError.Message)
	}

	if err := nan.ReplaceInFile(sUserSecurityScriptPath, "SAMUSER", p.out.sam); err != nil {
		LogError(nan.ErrFilesystemError.Message)
	}

	if err = exec.Command("/usr/bin/unix2dos", sUserSecurityScriptPath).Run(); err != nil {
		Log("Error when attempting to use unix2dos on file : %s", sUserSecurityScriptPath)
	}

	// TODO add timeout + retry on this call + message: did not respond in a timely manner

	if err = exec.Command("/usr/bin/scp", sUserSecurityScriptPath,
		"Administrator@10.20.12.20:/cygdrive/c/Windows/SYSVOL/domain/scripts/").Run(); err != nil {
		Log("Error when attempting to scp file: %s on server", sUserSecurityScriptPath)
		ExitError(nan.ErrSomethingWrong)
	}

	// Execute then delete security setup script on Windows Server
	// ===========================================================

	bExecSecurityScript := true

	// TODO add timeout on this call + did not respond in a timely manner

	if bExecSecurityScript {
		sWindowsServerSecurityFile := fmt.Sprintf(`\\winad.intra.nanocloud.com\NETLOGON\%s.setSecurity.bat`, p.out.sam)
		nan.Debug("Invoking WinExe on file:" + sWindowsServerSecurityFile)
		cmd := exec.Command(nan.Config().Proxy.WinExe,
			"-U", "intra.nanocloud.com/Administrator%3password",
			"--runas=intra.nanocloud.com/Administrator%3password", "//10.20.12.10", //TODO HARDCODED
			sWindowsServerSecurityFile)
		if out, err := cmd.Output(); err != nil {
			Log("Error <%s> returned by WinExe when running file <%s> with output: <%s>", err, sWindowsServerSecurityFile, out)
			ExitError(nan.ErrSomethingWrong)
		}

		// TODO add timeout + retry on this call + message: did not respond in a timely manner

		cmd = exec.Command(nan.Config().Proxy.WinExe,
			"-U", "intra.nanocloud.com/Administrator%3password",
			"--runas=intra.nanocloud.com/Administrator%3password", "//10.20.12.10", //TODO HARDCODED
			"cmd.exe /C DEL "+sWindowsServerSecurityFile)

		if out, err := cmd.Output(); err != nil {
			Log("Error code %s returned by WinExe when running file : %s with output: %s", err, sWindowsServerSecurityFile, out)
			ExitError(nan.ErrSomethingWrong)
		}
	}

	p.Result = nil
}

func (p *ProcCreateWinUser) Undo() {

	p.Result = nil

	// Refuse deletion if user account doesn't exist
	bRegistered, err := g_Db.IsUserRegistered(G_Account.Email)
	if err != nil {
		ExitError(ErrIssueWithAccountsDb)
	} else if !bRegistered {
		Log("Email address not listed in accounts database")
	}

	// TODO LDAP Plugin
	params := fmt.Sprintf(`{ "userid" : "%s" }`, G_Account.Email)
	resp := ""

	g_PluginLdap.Call("Ldap.ForceDisableAccount", params, &resp)
	if resp == "0" {
		LogError("Ldap.ForceDisableAccount failed")
	}
	sam, e := g_Db.GetSamFromEmail(G_Account.Email)
	if e != nil {
		ExitError(ErrIssueWithAccountsDb)
	}

	sam = strings.Trim(sam, " ")

	if sam == "" || sam == "unactivated" {
		LogError("Email ok but found no matching SAM user")
		return
	}

	// Delete Talend Studio instance : logoff user and remove user profile

	removeProfileSourcePath := fmt.Sprintf("%s/studio/removeProfile.cmd", nan.Config().CommonBaseDir)
	removeProfileDestPath := fmt.Sprintf("%s/studio/%s/%s.removeProfile.bat", nan.Config().CommonBaseDir, G_Account.Email, sam)

	if nan.CopyFile(removeProfileSourcePath, removeProfileDestPath) != nil {
		LogError("Failed to copy removeProfile for preparation")
	} else if err := nan.ReplaceInFile(removeProfileDestPath, "SAMUSER", sam); err != nil {
		LogError("Failed to edit removeProfile script")
	} else if err := exec.Command("/usr/bin/unix2dos", removeProfileDestPath).Run(); err != nil {
		LogError("Error when attempting to use unix2dos on file : %s", removeProfileDestPath)
	} else if err = exec.Command("/usr/bin/scp", removeProfileDestPath, "Administrator@10.20.12.20:/cygdrive/c/Windows/SYSVOL/domain/scripts/").Run(); err != nil {
		LogError("Error when attempting to scp file: %s on server", removeProfileDestPath)
	} else {

		// Execute and then delete the remove profile script

		AdRemoveProfileUncPath := fmt.Sprintf(`\\winad.intra.nanocloud.com\NETLOGON\%s.removeProfile.bat`, sam)

		// Run the logoff script and profile removal process

		Log("Requested remote exec: %s", AdRemoveProfileUncPath)

		cmd := exec.Command(nan.Config().Proxy.WinExe,
			"-U", "intra.nanocloud.com/Administrator%3password",
			"--runas=intra.nanocloud.com/Administrator%3password", "//10.20.12.10", //TODO HARDCODED
			AdRemoveProfileUncPath)
		if out, err := cmd.Output(); err != nil {
			Log("Error returned by WinExe when running user removeProfile.bat: %s, outpout: %s", err, string(out))
		}

		// Delete the removeProfile File on Windows Server

		cmd = exec.Command(nan.Config().Proxy.WinExe,
			"-U", "intra.nanocloud.com/Administrator%3password",
			"--runas=intra.nanocloud.com/Administrator%3password", "//10.20.12.10", //TODO HARDCODED
			"cmd.exe /C DEL "+AdRemoveProfileUncPath)
		if out, err := cmd.Output(); err != nil {
			Log("Error code %s returned by WinExe after attempting to delete removeProfile script: %s", err, string(out))
			// ? ExitError(ErrSomethingWrong)
		}

		// Delete the config File on Windows Server

		AdConfigScriptUncPath := fmt.Sprintf(`\\winad.intra.nanocloud.com\NETLOGON\%s.config.bat`, sam)

		cmd = exec.Command(nan.Config().Proxy.WinExe,
			"-U", "intra.nanocloud.com/Administrator%3password",
			"--runas=intra.nanocloud.com/Administrator%3password", "//10.20.12.10", //TODO HARDCODED
			"cmd.exe /C DEL "+AdConfigScriptUncPath)
		if out, err := cmd.Output(); err != nil {
			Log("Error %s returned by WinExe when attempting to delete : %s with output: %s", err, AdConfigScriptUncPath, out)
			// ? ExitError(ErrSomethingWrong)
		}
	}

	sFilePattern := fmt.Sprintf("%s/studio/%s/*", nan.Config().CommonBaseDir, G_Account.Email)

	workspaceFilenames, _ := filepath.Glob(sFilePattern)
	for _, fileName := range workspaceFilenames {
		if err := os.Remove(fileName); err != nil {
			LogError("Error when deleting file: %s, err: %s", fileName, err)
		}
	}

	workspaceDir := fmt.Sprintf("%s/studio/%s", nan.Config().CommonBaseDir, G_Account.Email)
	if err := os.Remove(workspaceDir); err != nil {
		LogError("Error when deleting directory: %s, err: %s", workspaceDir, err)
	}
}

// ====================================================================================================
// Procedure: RegisterProxyUser
// =>
// ====================================================================================================
type ProcRegisterProxyUser struct {
	nan.ProcedureStruct

	sam string
}

func (p *ProcRegisterProxyUser) Do() {

	if nan.DryRun || nan.ModeRef {
		Log("DRYRUN: Registring user in guacamole database")
		return
	}

	// uid := uuid.New()
	// saltBytes := sha256.Sum256([]byte(uid))
	// saltString := string(saltBytes[:])
	// saltHexString := hex.EncodeToString(saltBytes[:])

	// saltedPasswordBytes := sha256.Sum256([]byte(G_Account.Password + saltHexString))
	// saltedPassword := string(saltedPasswordBytes[:])

	//!!!!!!!!!!!!!!!!!!!!!!!!
	//TODO : Got a syntax error once with next request. Could be scrambled chars that affect SQL ?
	//!!!!!!!!!!!!!!!!!!!!!!!!
	//Log("Email: <%s>, saltString <%s>, saltedPassword <%s>", G_Account.Email, saltString, saltedPassword)

	// guacamole_user
	// sRequest := fmt.Sprintf("INSERT INTO guacamole_user(username, password_salt, password_hash) VALUES ('%s', '%s', '%s');",
	// 	G_Account.Email, saltString, saltedPassword)

	// if _, err := g_Db.Exec(sRequest); err != nil {
	// 	LogError("Failed to insert user data into guacamole_user, error: %s", err)
	// 	ExitError(ErrSomethingWrong)
	// }

	if _, err := g_Db.Exec("SET @salt = UNHEX(SHA2(UUID(), 256))"); err != nil {
		LogError("Failed to compute salt")
		ExitError(nan.ErrSomethingWrong)
	}

	sInsertUserPassword := fmt.Sprintf(`INSERT INTO guacamole_user(username, password_salt, password_hash) 
		VALUES ('%s', @salt, UNHEX(SHA2(CONCAT('%s', HEX(@salt)), 256)) );`, G_Account.Email, G_Account.Password)

	if _, err := g_Db.Exec(sInsertUserPassword); err != nil {
		LogError("Failed to insert user data into guacamole_user, error: %s", err)
		ExitError(nan.ErrSomethingWrong)
	}

	// guacamole_connection_group
	sRequest := fmt.Sprintf("INSERT INTO guacamole_connection_group (connection_group_name, type) VALUES ('%s', '%s')",
		G_Account.Email, "ORGANIZATIONAL")

	if _, err := g_Db.Exec(sRequest); err != nil {
		LogError("Failed to insert user data into guacamole_connection_group, error: %s", err)
		ExitError(nan.ErrSomethingWrong)
	}

	// Retrieve group ID
	sQuery := fmt.Sprintf("select connection_group_id from guacamole_connection_group where connection_group_name='%s' LIMIT 1",
		G_Account.Email)

	rows, err := g_Db.Query(sQuery)
	if err != nil {
		LogError("Failed to select connection_group_id for user %s, error: %s", G_Account.Email, err)
		ExitError(nan.ErrSomethingWrong)
	}

	groupId := ""

	for rows.Next() {
		if err = rows.Scan(&groupId); err != nil {
			LogError("Failed to parse result of query on connection_group_id for user %s, error: %s", G_Account.Email, err)
			ExitError(nan.ErrSomethingWrong)
		}
	}

	// retrieve user ID
	userId := ""

	sQuery = fmt.Sprintf("select user_id from guacamole_user where username='%s' LIMIT 1", G_Account.Email)
	rows, err = g_Db.Query(sQuery)
	if err != nil {
		LogError("Failed to select user_id for user %s, error: %s", G_Account.Email, err)
		ExitError(nan.ErrSomethingWrong)
	}

	for rows.Next() {
		if err = rows.Scan(&userId); err != nil {
			LogError("Failed to parse result of query on connection_group_id for user %s, error: %s", G_Account.Email, err)
			ExitError(nan.ErrSomethingWrong)
		}
	}

	// guacamole_connection_group_permission

	// Add group to user

	sRequest = fmt.Sprintf("INSERT INTO guacamole_connection_group_permission (user_id, connection_group_id, permission) VALUES (%s, %s, 'READ')",
		userId, groupId)

	if _, err = g_Db.Exec(sRequest); err != nil {
		LogError("Failed to insert record into guacamole_connection_group_permission, error : %s", err)
		ExitError(nan.ErrSomethingWrong)
	}

	//TODO OPTIONAL (was done for Talend)
	connectionName := "Visual"           //ESI 	//"Talend Studio"
	RemoteAppName := "VisualEnvironment" //ESI 	//"TalendStudiowinx86_64"

	// guacamole_connection
	sRequest = fmt.Sprintf("INSERT INTO guacamole_connection (connection_name, parent_id,  protocol) VALUES ('%s', '%s', 'rdp')",
		connectionName, groupId)

	var connId int64

	if result, err := g_Db.Exec(sRequest); err != nil {
		LogError("Failed to insert record into guacamole_connection, error : %s", err)
		ExitError(nan.ErrSomethingWrong)
	} else {
		var err error
		if connId, err = result.LastInsertId(); err != nil {
			LogError("Failed to get LastInsert from query on guacamole_connection, error : %s", err)
			ExitError(nan.ErrSomethingWrong)
		}
	}

	insertRequests := []string{

		fmt.Sprintf("INSERT INTO guacamole_connection_parameter VALUES (%d, 'hostname', '10.20.12.10')", connId), //TODO HARDCODED
		fmt.Sprintf("INSERT INTO guacamole_connection_parameter VALUES (%d, 'port', '3389')", connId),
		fmt.Sprintf("INSERT INTO guacamole_connection_parameter VALUES (%d, 'username', '%s')", connId, p.sam),
		fmt.Sprintf("INSERT INTO guacamole_connection_parameter VALUES (%d, 'password', '%s')", connId, G_Account.Password),
		fmt.Sprintf("INSERT INTO guacamole_connection_parameter VALUES (%d, 'domain', 'intra.nanocloud.com')", connId),
		fmt.Sprintf("INSERT INTO guacamole_connection_parameter VALUES (%d, 'color-depth', '24')", connId),
		fmt.Sprintf("INSERT INTO guacamole_connection_parameter VALUES (%d, 'disable-audio', 'true')", connId),
		fmt.Sprintf("INSERT INTO guacamole_connection_parameter VALUES (%d, 'remote-app', '||%s')", connId, RemoteAppName),
	}

	for _, requestStr := range insertRequests {
		if _, err = g_Db.Exec(requestStr); err != nil {
			LogError("Failed to execute query: [%s], error: %s", requestStr, err)
			nan.ExitError(nan.ErrSomethingWrong)
		}
	}

	// guacamole_connection_parameter

	// -- Add connection to user

	sRequest = fmt.Sprintf("INSERT INTO guacamole_connection_permission (user_id, connection_id, permission) VALUES ('%s', '%d', 'READ')",
		userId, connId)

	if _, err := g_Db.Exec(sRequest); err != nil {
		LogError("Failed to insert record into guacamole_connection_permission, error : %s", err)
		ExitError(nan.ErrSomethingWrong)
	}

	p.Result = nil
}

// 	Remove Guacamole user
func (p *ProcRegisterProxyUser) Undo() {

	ok := true

	// Archive the user connection logs
	sRequest := fmt.Sprintf(`INSERT IGNORE INTO guacamole_deleted_connection_history
SELECT history_id, guacamole_connection_history.user_id, username, start_date, end_date
from guacamole_connection_history
INNER JOIN guacamole_user ON guacamole_connection_history.user_id = guacamole_user.user_id
WHERE username = '%s'`, G_Account.Email)

	/* Also delete from ? :
	guacamole_connection_group_permission
	guacamole_connection_parameter
	guacamole_connection_permission
	*/

	if _, err := g_Db.Exec(sRequest); err != nil {
		LogError("Failed to archive user connection info in table guacamole_deleted_connection_history : %s", err)
		ok = false
	}

	// Delete group
	sRequest = fmt.Sprintf("DELETE FROM guacamole_connection_group WHERE connection_group_name='%s'", G_Account.Email)
	if _, err := g_Db.Exec(sRequest); err != nil {
		LogError("Failed to delete record from table guacamole_connection_group, error: ", err)
		ok = false
	}

	// Delete user
	sRequest = fmt.Sprintf("DELETE FROM guacamole_user WHERE username = '%s'", G_Account.Email)
	if _, err := g_Db.Exec(sRequest); err != nil {
		LogError("Failed to delete record from table guacamole_user, error: ", err)
		ok = false
	}

	if !ok {
		ExitError(ErrIssueWithAccountsDb)
	}
}
