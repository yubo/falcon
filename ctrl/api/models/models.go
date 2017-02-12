package models

import (
	"errors"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	DB_PREFIX      = ""
	PAGE_PER       = 10
	SYS_TAG_SCHEMA = "cop,owt,pdl;servicegroup;service,jobgroup;job,sbs;mod;srv;grp;cluster;"
	SYS_W_SCOPE    = "falcon_write"
	SYS_R_SCOPE    = "falcon_read"
	SYS_B_SCOPE    = "falcon_bind"
	SYS_O_SCOPE    = "falcon_operate"
	SYS_A_SCOPE    = "falcon_admin"
)

var (
	dbTables = []string{
		"action",
		"expression",
		"host",
		"kv",
		"log",
		"role",
		"rule",
		"session",
		"strategy",
		"tag",
		"tag_host",
		"tag_rel",
		"tag_tpl",
		"team",
		"team_user",
		"template",
		"token",
		"tpl_rel",
		"trigger",
		"user",
	}
)

// ctl meta name
const (
	CTL_M_HOST = iota
	CTL_M_ROLE
	CTL_M_SYSTEM
	CTL_M_TAG
	CTL_M_USER
	CTL_M_TOKEN
	CTL_M_TPL
	CTL_M_RULE
	CTL_M_TEMPLATE
	CTL_M_TRIGGER
	CTL_M_EXPRESSION
	CTL_M_TEAM
	CTL_M_SIZE
)

// ctl method name
const (
	CTL_A_ADD = iota
	CTL_A_DEL
	CTL_A_SET
	CTL_A_GET
	CTL_A_SIZE
)

// db.tpl_rel.type_id
const (
	TPL_REL_T_ACL_USER = iota
	TPL_REL_T_ACL_TOKEN
	TPL_REL_T_RULE_TRIGGER
)

// db.kv.type_id
const (
	KV_T_CONFIG = iota
	KV_T_CACHE
)

type Id struct {
	Id int64 `json:"id"`
}

type Total struct {
	Total int64 `json:"total"`
}

var (
	cacheModule  [CTL_M_SIZE]cache
	sysTagSchema *TagSchema

	moduleName [CTL_M_SIZE]string = [CTL_M_SIZE]string{
		"host", "role", "system", "tag", "user", "token",
		"template", "rule", "trigger", "expression", "team",
	}

	actionName [CTL_A_SIZE]string = [CTL_A_SIZE]string{
		"add", "del", "set", "get",
	}

	ErrExist        = errors.New("object exists")
	ErrLogged       = errors.New("already logged in")
	ErrNoExits      = errors.New("object not exists")
	ErrAuthFailed   = errors.New("auth failed")
	ErrNoUsr        = errors.New("user not exists")
	ErrNoHost       = errors.New("host not exists")
	ErrNoTag        = errors.New("tag not exists")
	ErrNoRole       = errors.New("role not exists")
	ErrNoToken      = errors.New("token not exists")
	ErrNoModule     = errors.New("module not exists")
	ErrNoRel        = errors.New("relation not exists")
	ErrNoRule       = errors.New("rule not exists")
	ErrNoTemplate   = errors.New("template not exists")
	ErrNoTrigger    = errors.New("trigger not exists")
	ErrNoExpression = errors.New("expression not exists")
	ErrNoTeam       = errors.New("team not exists")
	ErrNoStrategy   = errors.New("strategy not exists")
	ErrNoLogged     = errors.New("not logged in")
	ErrRePreStart   = errors.New("multiple times PreStart")
	ErrUnsupported  = errors.New("unsupported")
	ErrDelDefault   = errors.New("You cannot delete this basic data")
	ErrDelInUse     = errors.New("Still in use, cannot remove")
	ErrParam        = errors.New("param error")
	ErrEmpty        = errors.New("empty items")
	EPERM           = errors.New("Operation not permitted")
	ENOENT          = errors.New("No such file or directory")
	ESRCH           = errors.New("No such process")
	EINTR           = errors.New("Interrupted system call")
	EIO             = errors.New("I/O error")
	ENXIO           = errors.New("No such device or address")
	E2BIG           = errors.New("Argument list too long")
	ENOEXEC         = errors.New("Exec format error")
	EBADF           = errors.New("Bad file number")
	ECHILD          = errors.New("No child processes")
	EAGAIN          = errors.New("Try again")
	ENOMEM          = errors.New("Out of memory")
	EACCES          = errors.New("Permission denied")
	EFAULT          = errors.New("Bad address")
	ENOTBLK         = errors.New("Block device required")
	EBUSY           = errors.New("Device or resource busy")
	EEXIST          = errors.New("File exists")
	EXDEV           = errors.New("Cross-device link")
	ENODEV          = errors.New("No such device")
	ENOTDIR         = errors.New("Not a directory")
	EISDIR          = errors.New("Is a directory")
	EINVAL          = errors.New("Invalid argument")
	ENFILE          = errors.New("File table overflow")
	EMFILE          = errors.New("Too many open files")
	ENOTTY          = errors.New("Not a typewriter")
	ETXTBSY         = errors.New("Text file busy")
	EFBIG           = errors.New("File too large")
	ENOSPC          = errors.New("No space left on device")
	ESPIPE          = errors.New("Illegal seek")
	EROFS           = errors.New("Read-only file system")
	EMLINK          = errors.New("Too many links")
	EPIPE           = errors.New("Broken pipe")
	EDOM            = errors.New("Math argument out of domain of func")
	ERANGE          = errors.New("Math result not representable")
	EFMT            = errors.New("Invalid format") // custom
	EALLOC          = errors.New("Allocation Failure")
)

func init() {
	orm.RegisterModelWithPrefix(DB_PREFIX,
		new(User), new(Host), new(Tag),
		new(Role), new(Token), new(Log),
		new(Tag_rel), new(Tpl_rel), new(Team),
		new(Template), new(Trigger), new(Expression),
		new(Action), new(Strategy))

	// tag
	sysTagSchema, _ = NewTagSchema(SYS_TAG_SCHEMA)

	// auth
	AuthMap = make(map[string]AuthInterface)
	Auths = make([]AuthInterface, 0)

	// The hookfuncs will run in beego.Run()
	beego.AddAPPStartHook(start)
}

func start() (err error) {
	for _, auth := range strings.Split(beego.AppConfig.String("authmodule"), ",") {
		if auth, ok := AuthMap[auth]; ok {
			if auth.PreStart() == nil {
				Auths = append(Auths, auth)
			}
		}
	}

	CacheInit()

	PluginStart()
	// err = ConfigStart()

	metricInit(beego.AppConfig.String("metricfile"))

	return
}

func CacheInit() {
	for _, module := range strings.Split(beego.AppConfig.String("cachemodule"), ",") {
		for k, v := range moduleName {
			if v == module {
				cacheModule[k] = cache{
					enable: true,
					data:   make(map[int64]interface{}),
				}
				break
			}
		}
	}
}
