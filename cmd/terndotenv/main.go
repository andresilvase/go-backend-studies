package main

import (
	"fmt"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("topics/transactions/database/migrations/.env"); err != nil {
		panic(err)
	}

	cmd := exec.Command(
		"tern",
		"migrate",
		"-m",
		"topics/transactions/database/migrations",
		"--config",
		"topics/transactions/database/migrations/tern.conf",
	)

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Command failed with err:%v\n", err)
		fmt.Printf("Output %s", string(output))
		panic(err)
	}

	fmt.Printf("Command executed successfully: %s ", string(output))
}
