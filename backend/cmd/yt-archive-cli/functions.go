package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/google/uuid"
)

func showErroredTasks() {
	rows, err := db.Query("select id, description from tasks where status=4")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	fmt.Println("Errored tasks:")
	for rows.Next() {
		var id uuid.UUID
		var description string
		err = rows.Scan(&id, &description)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(id, description)
	}
}

func resetAllErroredTasks() {
	rowsAffected, err := execRowsAffected("update tasks set status=0 where status=4")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d tasks updated.\n", rowsAffected)
}

func showFinishedTasks() {
	rows, err := db.Query("select id, description from tasks where status=3")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	fmt.Println("Finished tasks:")

	count := 0
	const limit = 10

	for rows.Next() {
		var id uuid.UUID
		var description string
		err = rows.Scan(&id, &description)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(id, description)

		count++
		if count > limit {
			fmt.Println("... and more. Stopping here.")
			return
		}
	}
}

func positiveIntValidator(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("enter positive integer")
	}

	if n < 0 {
		return fmt.Errorf("enter positive integer")
	}

	return nil
}

func deleteFinishedTasks() {
	strPreserveN := "10"

	err := huh.NewInput().
		Title("Delete finished tasks").
		Description("How many finished tasks do you want to preserve?").
		Validate(positiveIntValidator).
		Value(&strPreserveN).
		Run()

	if err != nil {
		log.Fatal(err)
	}

	preserveN, err := strconv.Atoi(strPreserveN)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := execRowsAffected("delete from tasks where status=3 and id not in (select id from tasks where status=3 order by id desc limit ?)", preserveN)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d tasks deleted.\n", rowsAffected)
}
