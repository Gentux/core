package main

import (
	"encoding/json"
	"fmt"

	nan "nanocloud.com/zeroinstall/lib/libnan"

	"github.com/dullgiulio/pingo"
	"os/exec"
)

var (
	g_Hostname string
)

type AccountParams struct {
	UserId   string
	Password string
}

type Ldap struct{}

func (p *Ldap) Configure(_hostname string, _outMsg *string) error {
	g_Hostname = _hostname
	return nil
}

func (p *Ldap) AddUser(jsonParams string, _outMsg *string) error {
	*_outMsg = "0" // return code meaning failure of operation

	var params AccountParams

	if err := json.Unmarshal([]byte(jsonParams), &params); err != nil {
		r := nan.NewExitCode(0, "ERROR: failed to unmarshal Ldap.AccountParams : "+err.Error())
		*_outMsg = r.ToJson()
		return nil
	}

	sAddUserPhpScript := "add_LDAP_user.php"

	cmd := exec.Command("/usr/bin/php", "-f", sAddUserPhpScript, "--", params.UserId, params.Password,
		/* necessary ? */ "2>/dev/null")

	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to run script add_LDAP_user.php for email <%s> and password <%s>, error: %s, output: %s\n",
			params.UserId, params.Password, err, string(out))

	} else {
		*_outMsg = string(out)

	}

	return nil
}

func (p *Ldap) ForceDisableAccount(jsonParams string, _outMsg *string) error {
	*_outMsg = "0" // return code meaning failure of operation

	var params AccountParams

	if err := json.Unmarshal([]byte(jsonParams), &params); err != nil {
		r := nan.NewExitCode(0, "ERROR: failed to unmarshal Ldap.AccountParams : "+err.Error())
		*_outMsg = r.ToJson()
		return nil
	}

	sForceDisableUserPhpScript := "force_disable_LDAP_user.php"

	//fmt.Println("Running force disable for:", params.UserId)

	cmd := exec.Command("/usr/bin/php", "-f", sForceDisableUserPhpScript, "--", params.UserId,
		/* necessary ? */ "2>/dev/null")

	//fmt.Println("Php done")

	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to run script force_disable_LDAP_user.php for email <%s>, error: %s, output: %s\n",
			params.UserId, err, string(out))

	} else {
		*_outMsg = "1" //success
		//*_outMsg = string(out)

	}

	// Log("LDAP Check... %s account(s) disabled", string(out))

	return nil
}

func (p *Ldap) DisableAccount(jsonParams string, _outMsg *string) error {
	*_outMsg = "0" // return code meaning failure of operation

	var params AccountParams

	if err := json.Unmarshal([]byte(jsonParams), &params); err != nil {
		r := nan.NewExitCode(0, "ERROR: failed to unmarshal Ldap.AccountParams : "+err.Error())
		*_outMsg = r.ToJson()
		return nil
	}

	// *_outMsg = ImpLdapDisableAccount(params.UserId)

	sDisableUserPhpScript := "disable_LDAP_user.php"

	cmd := exec.Command("/usr/bin/php", "-f", sDisableUserPhpScript, "--", params.UserId,
		/* necessary ? */ "2>/dev/null")

	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to run script disable_LDAP_user.php for email <%s>, error: %s, output: %s\n",
			params.UserId, err, string(out))
	} else {
		// *_outMsg = string(out)
		*_outMsg = "1" // success
	}

	// sPhpScript := fmt.Sprintf("%s/disable_LDAP_user.php ", nan.Config().CommonBaseDir)
	// cmd := exec.Command("/usr/bin/php", "-f", sPhpScript, "--", _Sam)

	// _, err := cmd.Output()
	// if err != nil {
	// 	LogError("Error returned by script disable_LDAP_user.php, error: %s", err)
	// 	ExitError(nan.ErrSomethingWrong)
	// }

	return nil
}

func main() {

	plugin := &Ldap{}

	pingo.Register(plugin)

	pingo.Run()
}
