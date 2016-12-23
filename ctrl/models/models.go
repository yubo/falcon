package models

import (
	"errors"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	DB_PREFIX = ""
	PAGE_PER  = 20
)

const (
	CTL_M_HOST = iota
	CTL_M_ROLE
	CTL_M_SYSTEM
	CTL_M_TAG
	CTL_M_USER
	CTL_M_SIZE
)

const (
	CTL_A_ADD = iota
	CTL_A_DEL
	CTL_A_SET
	CTL_A_GET
	CTL_A_SIZE
)

var (
	cacheModule [CTL_M_SIZE]cache

	moduleName [CTL_M_SIZE]string = [CTL_M_SIZE]string{
		"host", "role", "system", "tag", "user",
	}

	actionName [CTL_A_SIZE]string = [CTL_A_SIZE]string{
		"add", "del", "set", "get",
	}

	ErrExist      = errors.New("object exists")
	ErrNoExits    = errors.New("object not exists")
	ErrAuthFailed = errors.New("auth failed")
	ErrNoUsr      = errors.New("user not exists")
	ErrNoHost     = errors.New("host not exists")
	ErrNoTag      = errors.New("tag not exists")
	ErrNoRole     = errors.New("role not exists")
	ErrNoSystem   = errors.New("system not exists")
	ErrNoScope    = errors.New("scope not exists")
	ErrNoLogged   = errors.New("not Logged")

	ErrUnsupported = errors.New("unsupported")
	ErrParam       = errors.New("param error")
	ErrEmpty       = errors.New("empty items")
	EPERM          = errors.New("Operation not permitted")
	ENOENT         = errors.New("No such file or directory")
	ESRCH          = errors.New("No such process")
	EINTR          = errors.New("Interrupted system call")
	EIO            = errors.New("I/O error")
	ENXIO          = errors.New("No such device or address")
	E2BIG          = errors.New("Argument list too long")
	ENOEXEC        = errors.New("Exec format error")
	EBADF          = errors.New("Bad file number")
	ECHILD         = errors.New("No child processes")
	EAGAIN         = errors.New("Try again")
	ENOMEM         = errors.New("Out of memory")
	EACCES         = errors.New("Permission denied")
	EFAULT         = errors.New("Bad address")
	ENOTBLK        = errors.New("Block device required")
	EBUSY          = errors.New("Device or resource busy")
	EEXIST         = errors.New("File exists")
	EXDEV          = errors.New("Cross-device link")
	ENODEV         = errors.New("No such device")
	ENOTDIR        = errors.New("Not a directory")
	EISDIR         = errors.New("Is a directory")
	EINVAL         = errors.New("Invalid argument")
	ENFILE         = errors.New("File table overflow")
	EMFILE         = errors.New("Too many open files")
	ENOTTY         = errors.New("Not a typewriter")
	ETXTBSY        = errors.New("Text file busy")
	EFBIG          = errors.New("File too large")
	ENOSPC         = errors.New("No space left on device")
	ESPIPE         = errors.New("Illegal seek")
	EROFS          = errors.New("Read-only file system")
	EMLINK         = errors.New("Too many links")
	EPIPE          = errors.New("Broken pipe")
	EDOM           = errors.New("Math argument out of domain of func")
	ERANGE         = errors.New("Math result not representable")
	EFMT           = errors.New("Invalid format") // custom
	EALLOC         = errors.New("Allocation Failure")
)

func init() {
	orm.RegisterModelWithPrefix(DB_PREFIX, new(User), new(Host),
		new(Tag), new(Role), new(System), new(Log))

	AuthMap = make(map[string]AuthInterface)
	Auths = make([]AuthInterface, 0)

	beego.AddAPPStartHook(start)
}

/* The hookfuncs will run in beego.Run() */
func start() error {
	for _, auth := range strings.Split(beego.AppConfig.String("authmodule"), ",") {
		beego.Debug(auth)
		if auth, ok := AuthMap[auth]; ok {
			Auths = append(Auths, auth)
		}
	}

	for _, module := range strings.Split(beego.AppConfig.String("cachemodule"), ",") {
		for k, v := range moduleName {
			if v == module {
				cacheModule[k] = cache{
					enable: true,
					data:   make(map[int]interface{}),
				}
				break
			}
		}
	}
	return nil
}
