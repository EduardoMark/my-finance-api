package main

import (
	"fmt"
	"os/exec"

	"github.com/EduardoMark/my-finance-api/pkg/config"
)

func main() {
	_, err := config.LoadEnv()
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(
		"tern",
		"migrate",
		"--migrations",
		"./internal/store/pgstore/migrations",
		"config",
		"--config",
		"./internal/store/pgstore/migrations/tern.conf",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Command execution failed: ", err)
		fmt.Println("Output", string(output))
		panic(err)
	}

	fmt.Println("Command executed successfuly: ", string(output))
}
