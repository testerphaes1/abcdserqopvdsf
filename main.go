package main

import "test-manager/cmd"

//go:generate sqlboiler --wipe --no-tests --add-soft-deletes psql -o usecase_models/boiler

func main() {
	cmd.Execute()
}
