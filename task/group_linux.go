package task

import "time"

func createGroup(g Group) error {
	_, err := ExecCmd(time.Second*10, "groupadd", g.Name)
	return err
}

func removeGroup(g Group) error {
	_, err := ExecCmd(time.Second*10, "groupdel", g.Name)
	return err
}
