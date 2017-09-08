package task

import "time"

func createUser(u User) error {
	x := []string{}
	if u.Group != "" {
		x = append(x, "-g", u.Group)
	}
	if u.Home != "" {
		x = append(x, "-d", u.Home)
	}
	x = append(x, u.Name)
	if _, err := ExecCmd(time.Second*10, "useradd", x...); err != nil {
		return err
	}
	return nil
}

func removeUser(u User) error {
	_, err := ExecCmd(time.Second*10, "userdel", u.Name)
	return err
}
