package main

import (
	"log"
	"net/http"
	"zoo/router"

	"github.com/gorilla/mux"
)

func initializeRouter() {
	r := mux.NewRouter()
	r.HandleFunc("/animal-types", router.GetAnimalTypes).Methods("GET")
	r.HandleFunc("/animal-types", router.CreateAnimalType).Methods("POST")
	r.HandleFunc("/animals", router.CreateAnimal).Methods("POST")
	r.HandleFunc("/food-types", router.CreateFoodType).Methods("POST")
	r.HandleFunc("/animal-feedings", router.FeedAnimal).Methods("POST")
	r.HandleFunc("/animals", router.GetAnimals).Methods("GET")
	r.HandleFunc("/animals/{id}", router.GetAnimal).Methods("GET")
	r.HandleFunc("/animals/{id}", router.DeleteAnimal).Methods("DELETE")
	r.HandleFunc("/save-animal-types-csv", router.SaveAnimalTypesCSVHandler).Methods("GET")
	r.HandleFunc("/save-animals-csv", router.SaveAnimalsCSVHandler).Methods("GET")
	r.HandleFunc("/save-food-types-csv", router.SaveFoodTypesCSVHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}
func main() {
	router.InitialMigration()
	initializeRouter()

}
