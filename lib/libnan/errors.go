package libnan

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const ()

// ===================================================================================================
// TYPES
// ===================================================================================================

type Error struct {
	Code    int
	Message string
	Details string

	jsonBytes []byte `json:"-"`
}

func NewExitCode(_code int, _msg string) *Error {
	p := &Error{Code: _code, Message: _msg}

	jsonBytes, err := json.Marshal(p)
	if err != nil {
		log.Printf("Error when json marshalling : { %d, %s }", _code, _msg)
		os.Exit(-1)
	}

	p.jsonBytes = jsonBytes

	return p
}

func (o Error) Ok() bool {
	return o.Code == 1
}

func (o Error) Failed() bool {
	return o.Code != 1
}

func (p *Error) ToJson() string {
	return string(p.jsonBytes)
}

func (p *Error) ToString() string {
	return p.Message
}

func (p *Error) Unmarshal(s string) bool {
	if err := json.Unmarshal([]byte(s), p); err != nil {
		LogError("Failed to unmarshal exit code from string: %s", s)
		return false
	}
	return true
}

func PrintErrorJson(_pError *Error) {
	fmt.Println(_pError.ToJson())
}

func PrintOk(_pExitCode *Error) {
	Log(_pExitCode.Message)
	fmt.Println(_pExitCode.ToJson())
}

func ExitOk(_pExitCode *Error) {
	Log(_pExitCode.Message)
	fmt.Println(_pExitCode.ToJson())

	os.Exit(0)
}

func ExitError(_pExitCode *Error) {
	LogError(_pExitCode.Message)
	fmt.Println(_pExitCode.ToJson())

	os.Exit(-1)
}

func Errorf(_msg string, _args ...interface{}) *Error {
	msg := fmt.Sprintf(_msg, _args...)

	return &Error{Code: 0, Message: msg}
}

func ExitErrorf(_code int, _msg string, _args ...interface{}) {
	msg := fmt.Sprintf(_msg, _args...)

	LogError(msg)

	p := &Error{Code: _code, Message: msg}

	jsonBytes, err := json.Marshal(p)
	if err != nil {
		log.Printf("Error when json marshalling : { %d, %s }", _code, _msg)
		return
	}

	fmt.Println(string(jsonBytes))
	os.Exit(_code)
}

// ===================================================================================================

var (
	ErrUnset                = NewExitCode(0, "Error not set")
	ErrSomethingWrong       = NewExitCode(0, "Something went wrong")
	ErrOk                   = NewExitCode(1, "Operation succeeded")
	OkSuccess               = NewExitCode(1, "Operation succeeded")
	ErrUnknownUuid          = NewExitCode(3, "Unknown resource UUID")
	ErrOpFailed             = NewExitCode(4, "Operation failed: resource is not in the state required to perform the operation")
	ErrPbWithEmailFormat    = NewExitCode(2, "Problem with email format")
	ErrPasswordNonCompliant = NewExitCode(3, "The password does not respect the security policy")
	ErrFilesystemError      = NewExitCode(16, "Filesystem error : failed to create/delete file/directory")
	ErrSystemError          = NewExitCode(17, "System error")
	ErrConfigError          = NewExitCode(18, "Error in configuration file")
	ErrPluginError          = NewExitCode(19, "Plugin communication error")
	ErrConsulNotFound       = NewExitCode(100, "Could not access Consul executable")
	ErrCouldNotPingVm       = NewExitCode(101, "Could not ping VM")
	ErrJsonParsingError     = NewExitCode(103, "Error when parsing JSON, see TAC log")
	ErrSshConnFailureonVm   = NewExitCode(104, "Failed to initiate SSH root on vm")

	ErrErrorWithExternalExe = NewExitCode(2, "Error returned by external executable, see TAC log")

	// VM creation
	ErrDuringVmCreation        = NewExitCode(0, "nc_create_vm did not return a valid VM uuid")
	ErrFailedToLocateVmProcess = NewExitCode(1, "Failed to locate VM pid in process list")

	// "message" : "Corrupt state, pool said to be non empty but has no VMs listed" }`)

	// Inside VM
	ErrCommandFailedInVm = NewExitCode(105, "Failed to run command on VM via SSH")
)
