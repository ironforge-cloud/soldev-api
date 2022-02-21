package main

import (
	"api/internal/database"
	"api/internal/types"
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	db, err := database.GetConnection()
	defer db.Close()

	if err != nil {
		log.Println(err)
	}

	projects, err := database.GetProjects(db)

	var groupedProjects []types.GroupedProjects

	counter := 0
	for _, item := range projects {
		for _, e := range item {

      if groupedProjects == nil {
        groupedProjects = append(groupedProjects, struct{})
      }
			groupedProjects[counter].Projects = append(groupedProjects[counter].Projects, e.Project)
			groupedProjects[counter].CategoryName = e.CategoryName
			groupedProjects[counter].CategoryID = e.CategoryID
			counter++
		}

	}

	if err != nil {
		log.Println(err)
	}

	response, err := json.Marshal(&projects)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(response)
}
