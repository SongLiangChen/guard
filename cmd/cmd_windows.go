// +build windows

package cmd

import (
	"guard/model"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func RunCommand(cmd string) (string, error) {
	result, err := exec.Command("cmd", "/c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result)), err
}

func GetProcessStat(processName string) *model.Stat {
	a := `ps aux | awk '/` + processName + `/ && !/awk/ {print $2,$3,$4}'`
	println(a)
	val, err := RunCommand(a)
	if err != nil {
		println(err.Error())
		return nil
	}

	println(val)

	ss := strings.Split(val, " ")
	if len(ss) != 3 {
		return &model.Stat{}
	}

	stat := &model.Stat{
		Pid: ss[0],
	}
	stat.CPU, _ = strconv.ParseFloat(ss[1], 64)
	stat.MEM, _ = strconv.ParseFloat(ss[2], 64)

	return stat
}

func RestartProcess(restartCmd string) error {
	cmd := exec.Command("/bin/sh", "-c", restartCmd)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_UNICODE_ENVIRONMENT,
	}
	return cmd.Run()
}
