package lib

// JournalJSON holds mako structured json
type JournalJSON struct {
	Cursor                  string `json:"__CURSOR,omitempty"`
	RealtimeTimestamp       string `json:"__REALTIME_TIMESTAMP,omitempty"`
	MonotonicTimestamp      string `json:"__MONOTONIC_TIMESTAMP,omitempty"`
	BootID                  string `json:"_BOOT_ID,omitempty"`
	Priority                string `json:"PRIORITY,omitempty"`
	Message                 string `json:"MESSAGE,omitempty"`
	ContainerID             string `json:"CONTAINER_ID,omitempty"`
	ContainerIDFull         string `json:"CONTAINER_ID_FULL,omitempty"`
	ContainerName           string `json:"CONTAINER_NAME,omitempty"`
	Transport               string `json:"_TRANSPORT,omitempty"`
	PID                     string `json:"_PID,omitempty"`
	UID                     string `json:"_UID,omitempty"`
	GID                     string `json:"_GID,omitempty"`
	Comm                    string `json:"_COMM,omitempty"`
	Exe                     string `json:"_EXE,omitempty"`
	CmdLine                 string `json:"_CMDLINE,omitempty"`
	CapEffective            string `json:"_CAP_EFFECTIVE,omitempty"`
	SystemdCGroup           string `json:"_SYSTEMD_CGROUP,omitempty"`
	SystemdUnit             string `json:"_SYSTEMD_UNIT,omitempty"`
	SystemdSlice            string `json:"_SYSTEMD_SLICE,omitempty"`
	SeLinuxContext          string `json:"_SELINUX_CONTEXT,omitempty"`
	SourceRealtimeTimestamp string `json:"_SOURCE_REALTIME_TIMESTAMP,omitempty"`
	MachineID               string `json:"_MACHINE_ID,omitempty"`
	HostName                string `json:"_HOSTNAME"`
}
