package dag

import (
	// nolint // gosec
	"crypto/md5"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

// DAG contains all information about a workflow.
type DAG struct {
	// Location is the absolute path to the DAG file.
	Location string `json:"Location"`
	// Group is the group name of the DAG. This is optional.
	Group string `json:"Group"`
	// Name is the name of the DAG. The default is the filename without the
	// extension.
	Name string `json:"Name"`
	// Tags contains the list of tags for the DAG. optional.
	Tags []string `json:"Tags"`
	// Description is the description of the DAG. optional.
	Description string `json:"Description"`

	// Schedule configuration.
	// This is used by the scheduler to start / stop / restart the DAG.
	Schedule        []Schedule `json:"Schedule"`
	StopSchedule    []Schedule `json:"StopSchedule"`
	RestartSchedule []Schedule `json:"RestartSchedule"`

	// Env contains a list of environment variables to be set before running
	// the DAG.
	Env []string `json:"Env"`

	// LogDir is the directory where the logs are stored.
	// The actual log directory is LogDir + Name (with invalid characters
	// replaced with '_').
	LogDir string `json:"LogDir"`

	// Parameters configuration.
	// The DAG definition contains only DefaultParams. Params are automatically
	// set by the DAG loader.
	// DefaultParams contains the default parameters to be passed to the DAG.
	DefaultParams string `json:"DefaultParams"`
	// Params contains the list of parameters to be passed to the DAG.
	Params []string `json:"Params"`

	// Commands configuration to be executed in the DAG.
	// Steps contains the list of steps in the DAG.
	Steps []Step `json:"Steps"`
	// HandlerOn contains the steps to be executed on different events.
	HandlerOn HandlerOn `json:"HandlerOn"`

	// Preconditions contains the conditions to be met before running the DAG.
	// If the conditions are not met, the whole DAG is skipped.
	Preconditions []Condition `json:"Preconditions"`

	// Mail notification configuration.
	// MailOn contains the conditions to send mail.
	// SMTP contains the SMTP configuration.
	// If you don't want to repeat the SMTP configuration for each DAG, you can
	// set it in the base configuration.
	SMTP *SMTPConfig `json:"Smtp"`
	// ErrorMail contains the mail configuration for error.
	ErrorMail *MailConfig `json:"ErrorMail"`
	// InfoMail contains the mail configuration for info.
	InfoMail *MailConfig `json:"InfoMail"`
	// MailOn contains the conditions to send mail.
	MailOn *MailOn `json:"MailOn"`

	// Timeout is a field to specify the maximum execution time of the DAG task
	Timeout time.Duration `json:"Timeout"`
	// Misc configuration for DAG execution.
	// Delay is the delay before starting the DAG.
	Delay time.Duration `json:"Delay"`
	// RestartWait is the time to wait before restarting the DAG.
	RestartWait time.Duration `json:"RestartWait"`
	// MaxActiveRuns specifies the maximum concurrent steps to run in an
	// execution.
	MaxActiveRuns int `json:"MaxActiveRuns"`
	// MaxCleanUpTime is the maximum time to wait for cleanup when the DAG is
	// stopped.
	MaxCleanUpTime time.Duration `json:"MaxCleanUpTime"`
	// HistRetentionDays is the number of days to keep the history.
	HistRetentionDays int `json:"HistRetentionDays"`
}

// Schedule contains the cron expression and the parsed cron schedule.
type Schedule struct {
	// Expression is the cron expression.
	Expression string `json:"Expression"`
	// Parsed is the parsed cron schedule.
	Parsed cron.Schedule `json:"-"`
}

// HandlerOn contains the steps to be executed on different events in the DAG.
type HandlerOn struct {
	Failure *Step `json:"Failure"`
	Success *Step `json:"Success"`
	Cancel  *Step `json:"Cancel"`
	Exit    *Step `json:"Exit"`
}

// MailOn contains the conditions to send mail.
type MailOn struct {
	Failure bool `json:"Failure"`
	Success bool `json:"Success"`
}

// SMTPConfig contains the SMTP configuration.
type SMTPConfig struct {
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}

// MailConfig contains the mail configuration.
type MailConfig struct {
	From string `json:"From"`
	To   string `json:"To"`
	// Prefix is the prefix for the subject of the mail.
	Prefix string `json:"Prefix"`
	// AttachLogs is the flag to attach the logs in the mail.
	AttachLogs bool `json:"AttachLogs"`
}

// HandlerType is the type of the handler.
type HandlerType string

const (
	HandlerOnSuccess HandlerType = "onSuccess"
	HandlerOnFailure HandlerType = "onFailure"
	HandlerOnCancel  HandlerType = "onCancel"
	HandlerOnExit    HandlerType = "onExit"
)

func (e HandlerType) String() string {
	return string(e)
}

// ParseHandlerType converts a string to a HandlerType.
func ParseHandlerType(s string) HandlerType {
	return nameToHandlerType[s]
}

var (
	nameToHandlerType = map[string]HandlerType{
		"onSuccess": HandlerOnSuccess,
		"onFailure": HandlerOnFailure,
		"onCancel":  HandlerOnCancel,
		"onExit":    HandlerOnExit,
	}
)

var (
	defaultHistoryRetentionDays = 30
	defaultMaxCleanUpTime       = time.Second * 60
)

// setup sets the default values for the DAG.
func (d *DAG) setup() {
	// The default history retention days is 30 days.
	// It is the number of days to keep the history.
	// The older history is deleted when the DAG is executed.
	if d.HistRetentionDays == 0 {
		d.HistRetentionDays = defaultHistoryRetentionDays
	}

	// The default max cleanup time is 60 seconds.
	// It is the maximum time to wait for cleanup when the DAG gets a stop
	// signal. If the cleanup takes more than this time, the process of the DAG
	// is killed.
	if d.MaxCleanUpTime == 0 {
		d.MaxCleanUpTime = defaultMaxCleanUpTime
	}

	// set the default working directory for the steps if not set
	dir := filepath.Dir(d.Location)
	for i := range d.Steps {
		d.Steps[i].setup(dir)
	}

	// set the default working directory for the handler steps if not set
	if d.HandlerOn.Exit != nil {
		d.HandlerOn.Exit.setup(dir)
	}
	if d.HandlerOn.Success != nil {
		d.HandlerOn.Success.setup(dir)
	}
	if d.HandlerOn.Failure != nil {
		d.HandlerOn.Failure.setup(dir)
	}
	if d.HandlerOn.Cancel != nil {
		d.HandlerOn.Cancel.setup(dir)
	}
}

// HasTag checks if the DAG has the given tag.
func (d *DAG) HasTag(tag string) bool {
	for _, t := range d.Tags {
		if t == tag {
			return true
		}
	}

	return false
}

// SockAddr returns the unix socket address for the DAG.
// The address is used to communicate with the agent process.
// TODO: It needs to be unique for each process so that multiple processes can
// run in parallel.
func (d *DAG) SockAddr() string {
	s := strings.ReplaceAll(d.Location, " ", "_")
	name := strings.Replace(filepath.Base(s), filepath.Ext(filepath.Base(s)), "", 1)
	// nolint // gosec
	h := md5.New()
	_, _ = h.Write([]byte(s))
	bs := h.Sum(nil)
	// Socket name length must be shorter than 108 characters,
	// so we truncate the name.
	// 108 - 16 (length of the hash) - 34 (length remaining non-name) - 8 padding = 50
	lengthLimit := 50
	if len(name) > lengthLimit {
		name = name[:lengthLimit-1]
	}
	return filepath.Join("/tmp", fmt.Sprintf("@blackdagger-%s-%x.sock", name, bs))
}

// String implements the Stringer interface.
// It returns the string representation of the DAG.
// TODO: Remove if not needed.
func (d *DAG) String() string {
	ret := "{\n"
	ret = fmt.Sprintf("%s\tName: %s\n", ret, d.Name)
	ret = fmt.Sprintf(
		"%s\tDescription: %s\n", ret, strings.TrimSpace(d.Description),
	)
	ret = fmt.Sprintf("%s\tEnv: %v\n", ret, strings.Join(d.Env, ", "))
	ret = fmt.Sprintf("%s\tLogDir: %v\n", ret, d.LogDir)
	for i, s := range d.Steps {
		ret = fmt.Sprintf("%s\tStep%d: %s\n", ret, i, s.String())
	}
	ret = fmt.Sprintf("%s}\n", ret)
	return ret
}
