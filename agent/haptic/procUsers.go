package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	nan "nanocloud.com/zeroinstall/lib/libnan"

	"github.com/boltdb/bolt"
)

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
func ValidateCreateUserParams(accountParam AccountParams) *nan.Err {

	Log("Verifying parameters to create %s account", accountParam.Email)

	if !nan.ValidName(accountParam.FirstName) {
		return ErrFirstnameNonCompliant
	}

	if !nan.ValidName(accountParam.LastName) {
		return ErrLastnameNonCompliant
	}

	if !nan.ValidEmail(accountParam.Email) {
		return nan.ErrPbWithEmailFormat
	}

	if !nan.ValidPassword(accountParam.Password) {
		return nan.ErrPasswordNonCompliant
	}

	if nan.DryRun || nan.ModeRef {
		return nil
	}

	//TODO OPTionalize this

	// Refuse creation if problem with license file
	// NOTE: license may be valid but http license service backend may be down
	// resp, httpErr := http.Get(G_Account.License)
	// defer resp.Body.Close()
	// if httpErr != nil || resp.StatusCode != 200 {
	// 	return LogErrorCode(ErrInvalidLicenseFile)
	// }

	return nil
}

// ========================================================================================================================
// Procedure: RegisterUser
//
// Does:
// - Check Params
// - Register TAC user : insert record in db guacamode/talend_tac
// ========================================================================================================================
func RegisterUser(accountParam AccountParams) *nan.Err {

	if err := ValidateCreateUserParams(accountParam); err != nil {
		return LogErrorCode(err)
	}

	Log("STARTING registerUser for: %s", accountParam.Email)

	// Activation fails if no license specified
	// TODO Make license checking configurable, as a plugin or optional prerequisite step
	if accountParam.License == "" {
		accountParam.License = "n/a"
	}

	numReg, err := g_Db.CountRegisteredUsers()
	if err != nil {
		return LogErrorCode(ErrIssueWithAccountsDb)
	} else if numReg >= nan.Config().Proxy.MaxNumRegistrations {
		return LogErrorCode(ErrMaxNumAccountsRegistered)
	}

	bRegistered, err := g_Db.IsUserRegistered(accountParam.Email)
	if err != nil {
		return LogErrorCode(ErrIssueWithAccountsDb)
	} else if bRegistered {
		return LogErrorCode(ErrAccountExists)
	}

	var (
		user User = User{
			Activated:    false,
			Email:        accountParam.Email,
			Firstname:    accountParam.FirstName,
			Lastname:     accountParam.LastName,
			Password:     accountParam.Password,
			Sam:          "",
			CreationTime: "",
		}
		userJson []byte
		resp     string
	)

	t := time.Now()
	user.CreationTime = t.Format(time.RFC3339)

	params := fmt.Sprintf(`{ "UserEmail" : "%s" }`, user.Email)
	g_PluginLdap.Call("Ldap.ForceDisableAccount", params, &resp)
	if resp == "0" {
		LogError("Ldap.ForceDisableAccount failed")
	}
	Log("Configure Windows user profile")
	params = fmt.Sprintf(`{ "UserEmail" : "%s", "password" : "%s" }`, user.Email, user.Password)
	g_PluginLdap.Call("Ldap.AddUser", params, &user.Sam)

	userJson, e := json.Marshal(user)
	if e != nil {
		return LogErrorCode(nan.ErrSomethingWrong)
	}

	e = g_Db.Update(
		func(tx *bolt.Tx) error {
			bucket, e := tx.CreateBucketIfNotExists([]byte("users"))
			if e != nil {
				return e
			}

			return bucket.Put([]byte(accountParam.Email), userJson)
		})

	if e != nil {
		return LogErrorCode(nan.ErrSomethingWrong)
	} else {
		return OkAccountBeingCreated
	}
}

// ========================================================================================================================
// Procedure: ListUsers
//
// Does:
// - Return list of users
// ========================================================================================================================
func ListUsers() []User {
	var (
		user  User
		users []User
	)

	g_Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("users"))
		if bucket == nil {
			return errors.New("Bucket 'users' doesn't exist")
		}

		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			json.Unmarshal(value, &user)
			users = append(users, user)
		}

		return nil
	})

	return users
}

func UpdateUserPassword(Email string, Password string) bool {
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

				user.Password = Password
				jsonUser, _ := json.Marshal(user)
				bucket.Put([]byte(user.Email), jsonUser)
				result = true
				break
			}
		}

		return nil
	})

	Log("TODO : Call LDAP plugin to update passwords")

	return result
}

// Returns FirstName LastName Email Password License
func GetUserAccountParamsForActivation(_Email string) (string, string, string, string, string, *nan.Err) {

	// Reload user account params from workspace
	sUserParamsFilePath := fmt.Sprintf("%s/studio/%s/account_params", nan.Config().CommonBaseDir, _Email)

	var err error
	var bytesRead []byte

	if bytesRead, err = ioutil.ReadFile(sUserParamsFilePath); err != nil {
		return "", "", "", "", "",
			nan.LogError("Failed to read user account params from workspace file : %s", sUserParamsFilePath)
	}

	var FirstName, LastName, tmpEmail, Password, License string

	nItemsParsed, err := fmt.Sscanf(string(bytesRead), "%s\n%s\n%s\n%s\n%s",
		&FirstName, &LastName, &tmpEmail, &Password, &License)

	if err != nil || nItemsParsed != 5 {
		return "", "", "", "", "",
			LogError("Failed to parse user account params from workspace file : %s, read %d items",
				sUserParamsFilePath, nItemsParsed)
	}

	return FirstName, LastName, tmpEmail, Password, License, nil
}

// ========================================================================================================================
// Procedure: ActivateUser
//
// Does:
// - Check Params
// - Register TAC user : insert record in db guacamode/talend_tac
// ========================================================================================================================
func ActivateUser(accountParams AccountParams) *nan.Err {

	if !nan.ValidEmail(accountParams.Email) {
		return nan.ErrPbWithEmailFormat
	}

	Log("STARTING activateUser for: %s", G_Account.Email)

	// Reached maximum number of active users ?

	if maxNumAccounts := nan.Config().Proxy.MaxNumAccounts; maxNumAccounts > 0 {

		if nAccounts, err := g_Db.CountActiveUsers(); err != nil {
			return LogErrorCode(ErrIssueWithAccountsDb)
		} else if nAccounts >= maxNumAccounts {
			return LogErrorCode(ErrMaxNumAccountsReached)
		}
	}

	// User not registered yet ?

	if bRegistered, err := g_Db.IsUserRegistered(G_Account.Email); err != nil {
		return LogErrorCode(ErrIssueWithAccountsDb)
	} else if !bRegistered {
		return LogErrorCode(ErrAccountNotRegistered)
	}

	// If user already activated then do nothing
	if bValue, _ := g_Db.IsUserActivated(G_Account.Email); bValue {
		return LogErrorCode(ErrAccountActivated)
	}

	tmpEmail := ""
	err := nan.NewErr()
	G_Account.FirstName, G_Account.LastName, tmpEmail, G_Account.Password, G_Account.License, err =
		GetUserAccountParamsForActivation(G_Account.Email)
	if err != nil {
		return err
	}

	if G_Account.Email != tmpEmail {
		nan.LogError("[INCONSISTENCY] Email passed as parameter (%s) doesn't match email loaded in account params: %s",
			G_Account.Email, tmpEmail)
	}

	//defer UndoIfFailed(G_ProcCreateTac)
	//[OPT] Create account resource such as TAC VM
	// G_ProcCreateTac.Do()

	//defer UndoIfFailed(G_ProcCreateWinUser)
	G_ProcCreateWinUser.Fqdn = "n/a" //[OPT] = G_ProcCreateTac.Ans.TacUrl
	G_ProcCreateWinUser.Do()

	if G_ProcCreateWinUser.out.sam == "" {
		return LogErrorCode(ErrIssueWithTacProvisioning)
	}

	g_Db.UpdateUserSamForEmail(G_Account.Email, G_ProcCreateWinUser.out.sam)

	return OkAccountBeingActivated
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
		"-U", "intra.nanocloud.com/Administrator%3nexbAie2050",
		"--runas=intra.nanocloud.com/Administrator%3nexbAie2050", "//10.20.12.10", //TODOHARDCODED
		sWindowsServerSecurityFile)

	if out, err := cmd.Output(); err != nil {
		LogError("<%s> returned by WinExe when running file <%s> with output: <%s>", err, sWindowsServerSecurityFile, out)
		return
	}
}

func (p *ProcCreateWinUser) Do() *nan.Err {

	var err error

	if nan.DryRun || nan.ModeRef {
		Log("Creating Windows user + LDAP declaration")
		return nil
	}

	// TODO: Remove these check, redundant with earlier check ?
	if !nan.ValidEmail(G_Account.Email) {
		return LogErrorCode(nan.ErrPbWithEmailFormat)
	}

	if !nan.ValidPassword(G_Account.Password) {
		return LogErrorCode(nan.ErrPasswordNonCompliant)
	}

	if !G_TwoStageActivation {
		bRegistered, err := g_Db.IsUserRegistered(G_Account.Email)
		if err != nil {
			return LogErrorCode(ErrIssueWithAccountsDb)
		} else if bRegistered {
			return LogErrorCode(ErrAccountExists)
		}
	}

	// If this account is still enabled on the LDAP, deactivate it

	resp := ""

	params := fmt.Sprintf(`{ "UserEmail" : "%s" }`, G_Account.Email)

	g_PluginLdap.Call("Ldap.ForceDisableAccount", params, &resp)

	if resp == "0" {
		LogError("Ldap.ForceDisableAccount failed")
	}

	// Active Directory user
	// =====================

	Log("Configure Windows user profile")

	// Add LDAP user

	params = fmt.Sprintf(`{ "UserEmail" : "%s", "password" : "%s" }`, G_Account.Email, G_Account.Password)

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
	// 		return LogError("Failed to run script add_LDAP_user.php for email <%s> and password <%s>, error: %s", G_Account.Email, G_Account.Password, err)
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
	// 	return LogError("ErrFilesystemError : Failed to create directory: %s", err)
	// }

	// if err = os.Chmod(sWorkspaceDirPath, 0777); err != nil {
	// 	LogError("ErrSystemError : Failed to set permissions on directory: %s", err)
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
		LogError("when attempting to use unix2dos on file : %s, error: %s", samFilePath, err)
	}

	// TODO add timeout + retry on this call + message: did not respond in a timely manner

	if err = exec.Command("/usr/bin/scp", samFilePath, "Administrator@10.20.12.20:/cygdrive/c/Windows/SYSVOL/domain/scripts/").Run(); err != nil {
		return LogError("when attempting to scp file: %s on server", samFilePath)
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
		return LogError("when attempting to scp file: %s on server", sUserSecurityScriptPath)
	}

	// Execute then delete security setup script on Windows Server
	// ===========================================================

	bExecSecurityScript := true

	// TODO add timeout on this call + did not respond in a timely manner

	if bExecSecurityScript {
		sWindowsServerSecurityFile := fmt.Sprintf(`\\winad.intra.nanocloud.com\NETLOGON\%s.setSecurity.bat`, p.out.sam)
		nan.Debug("Invoking WinExe on file:" + sWindowsServerSecurityFile)
		cmd := exec.Command(nan.Config().Proxy.WinExe,
			"-U", "intra.nanocloud.com/Administrator%3nexbAie2050",
			"--runas=intra.nanocloud.com/Administrator%3nexbAie2050", "//10.20.12.10", //TODO HARDCODED
			sWindowsServerSecurityFile)
		if out, err := cmd.Output(); err != nil {
			return LogError("<%s> returned by WinExe when running file <%s> with output: <%s>", err, sWindowsServerSecurityFile, out)
		}

		// TODO add timeout + retry on this call + message: did not respond in a timely manner

		cmd = exec.Command(nan.Config().Proxy.WinExe,
			"-U", "intra.nanocloud.com/Administrator%3nexbAie2050",
			"--runas=intra.nanocloud.com/Administrator%3nexbAie2050", "//10.20.12.10", //TODO HARDCODED
			"cmd.exe /C DEL "+sWindowsServerSecurityFile)

		if out, err := cmd.Output(); err != nil {
			return LogError("code %s returned by WinExe when running file : %s with output: %s", err, sWindowsServerSecurityFile, out)
		}
	}

	p.Result = nil
	return p.Result
}

func (p *ProcCreateWinUser) Undo(accountParams AccountParams) *nan.Err {

	p.Result = nil

	// Refuse deletion if user account doesn't exist
	bRegistered, err := g_Db.IsUserRegistered(accountParams.Email)
	if err != nil {
		return LogErrorCode(ErrIssueWithAccountsDb)
	} else if !bRegistered {
		Log("Email address not listed in accounts database")
	}

	params := fmt.Sprintf(`{ "UserEmail" : "%s" }`, accountParams.Email)
	resp := ""

	g_PluginLdap.Call("Ldap.ForceDisableAccount", params, &resp)
	if resp == "0" {
		LogError("Ldap.ForceDisableAccount failed")
	}
	sam, e := g_Db.GetSamFromEmail(accountParams.Email)
	if e != nil {
		return LogErrorCode(ErrIssueWithAccountsDb)
	}

	sam = strings.Trim(sam, " ")

	if sam == "" || sam == "unactivated" {
		LogError("Email ok but found no matching SAM user")
		return nil
	}

	// Delete Talend Studio instance : logoff user and remove user profile

	removeProfileSourcePath := fmt.Sprintf("%s/studio/removeProfile.cmd", nan.Config().CommonBaseDir)
	removeProfileDestPath := fmt.Sprintf("%s/studio/%s/%s.removeProfile.bat", nan.Config().CommonBaseDir, accountParams.Email, sam)

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
			"-U", "intra.nanocloud.com/Administrator%3nexbAie2050",
			"--runas=intra.nanocloud.com/Administrator%3nexbAie2050", "//10.20.12.10", //TODO HARDCODED
			AdRemoveProfileUncPath)
		if out, err := cmd.Output(); err != nil {
			Log("Error returned by WinExe when running user removeProfile.bat: %s, outpout: %s", err, string(out))
		}

		// Delete the removeProfile File on Windows Server

		cmd = exec.Command(nan.Config().Proxy.WinExe,
			"-U", "intra.nanocloud.com/Administrator%3nexbAie2050",
			"--runas=intra.nanocloud.com/Administrator%3nexbAie2050", "//10.20.12.10", //TODO HARDCODED
			"cmd.exe /C DEL "+AdRemoveProfileUncPath)
		if out, err := cmd.Output(); err != nil {
			LogError("code %s returned by WinExe after attempting to delete removeProfile script: %s", err, string(out))
		}

		// Delete the config File on Windows Server

		AdConfigScriptUncPath := fmt.Sprintf(`\\winad.intra.nanocloud.com\NETLOGON\%s.config.bat`, sam)

		cmd = exec.Command(nan.Config().Proxy.WinExe,
			"-U", "intra.nanocloud.com/Administrator%3nexbAie2050",
			"--runas=intra.nanocloud.com/Administrator%3nexbAie2050", "//10.20.12.10", //TODO HARDCODED
			"cmd.exe /C DEL "+AdConfigScriptUncPath)
		if out, err := cmd.Output(); err != nil {
			LogError("%s returned by WinExe when attempting to delete : %s with output: %s", err, AdConfigScriptUncPath, out)
		}
	}

	sFilePattern := fmt.Sprintf("%s/studio/%s/*", nan.Config().CommonBaseDir, accountParams.Email)

	workspaceFilenames, _ := filepath.Glob(sFilePattern)
	for _, fileName := range workspaceFilenames {
		if err := os.Remove(fileName); err != nil {
			LogError("Error when deleting file: %s, err: %s", fileName, err)
		}
	}

	workspaceDir := fmt.Sprintf("%s/studio/%s", nan.Config().CommonBaseDir, accountParams.Email)
	if err := os.Remove(workspaceDir); err != nil {
		return LogError("Error when deleting directory: %s, err: %s", workspaceDir, err)
	}

	return nil
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

	g_Db.AddUser(G_User)
	// TODO Handle Users connections
}

// 	Remove Guacamole user
func (p *ProcRegisterProxyUser) Undo() {
	g_Db.DeleteUser(G_User)
	// TODO Handle Users connections
}
