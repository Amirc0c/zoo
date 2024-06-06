package router

import (
	"encoding/json"
	"log"
	"net/http"

	"encoding/csv"
	"fmt"

	"os"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AnimalType struct {
	ID          int    `json:"id"`
	Type        string `json:"name"` // name -> type
	Description string `json:"description"`
	// AnimalID    Animal `json:"animalID"`
}

type Animal struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	TypeID   int        `json:"type_id"`
	Type     AnimalType `json:"type" gorm:"foreignKey:TypeID"`
	Age      int        `json:"age"`
	Gender   string     `json:"gender"`
	Health   int        `json:"health"`
	FoodType []FoodType `json:"food_type" gorm:"foreignKey:AnimalID"`
}

type FoodType struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Amount       int    `json:"amount"`
	HealthPoints int    `json:"health_points"`
	AnimalID     uint   `json:"animalID"`
}

type AnimalFeedingRequest struct {
	AnimalID         int `json:"animal_id"`
	FoodHealthPoints int `json:"food_health_points"`
}

type CacheData struct {
	Value      interface{}
	Expiration time.Time
}

type Cache struct {
	storage map[string]CacheData
}

func NewCache() *Cache {
	return &Cache{
		storage: make(map[string]CacheData),
	}
}

var db *gorm.DB
var err error

var CacheInstance *Cache

func InitialMigration() {
	dsn := "host=localhost user=postgres password=Amir2009 dbname=postgres port=5432 sslmode=disable"

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	err = db.AutoMigrate(&Animal{}, &AnimalType{}, &FoodType{}, &AnimalFeedingRequest{})
	if err != nil {
		log.Fatalf("Ошибка при автомиграции таблиц: %v", err)
	}
}

func GetAnimalTypes(w http.ResponseWriter, r *http.Request) {

	var animalTypes []AnimalType
	if err := db.Find(&animalTypes).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func CreateAnimalType(w http.ResponseWriter, r *http.Request) {

	var animalType AnimalType
	if err := json.NewDecoder(r.Body).Decode(&animalType); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.Create(&animalType)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Println("Created animal type:", animalType)
}

func CreateAnimal(w http.ResponseWriter, r *http.Request) {
	var animal Animal
	if err := json.NewDecoder(r.Body).Decode(&animal); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.Create(&animal)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

}

func CreateFoodType(w http.ResponseWriter, r *http.Request) {
	var foodType FoodType
	if err := json.NewDecoder(r.Body).Decode(&foodType); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.Create(&foodType)
	fmt.Println(result)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println(result.RowsAffected)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(foodType)
	fmt.Println("create food type")
}

func FeedAnimal(w http.ResponseWriter, r *http.Request) {
	var feedingRequest AnimalFeedingRequest
	if err := json.NewDecoder(r.Body).Decode(&feedingRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var animal Animal
	if err := db.First(&animal, feedingRequest.AnimalID).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	animal.Health += feedingRequest.FoodHealthPoints
	if animal.Health > 100 {
		animal.Health = 100
	}
	if err := db.Save(&animal).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println("feedanimal")
}

func GetAnimals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var animals []Animal
	db.Find(&animals)
	json.NewEncoder(w).Encode(animals)
	fmt.Println("getanimals")
}

func GetAnimal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var animal Animal
	db.First(&animal, params["id"])
	json.NewEncoder(w).Encode(animal)
	fmt.Println("get animal")
}

func DeleteAnimal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var animal Animal
	db.Delete(&animal, params["id"])
	json.NewEncoder(w).Encode(" успешно")
	fmt.Println("delete")
}
func SaveAnimalTypesToCSV(filename string, animalTypes []AnimalType) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"ID", "Name", "Description"}
	writer.Write(headers)

	for _, animalType := range animalTypes {
		record := []string{
			fmt.Sprintf("%d", animalType.ID),
			animalType.Type,
			animalType.Description,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
func SaveAnimalTypesCSVHandler(w http.ResponseWriter, r *http.Request) {
	var animalTypes []AnimalType
	if err := db.Find(&animalTypes).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := SaveAnimalTypesToCSV("animal_types.csv", animalTypes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println("Saved AnimalType to CSV")
}

func SaveAnimalsToCSV(filename string, animals []Animal) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"ID", "Name", "TypeID", "Age", "Gender", "Health"}
	writer.Write(headers)

	for _, animal := range animals {
		record := []string{
			fmt.Sprintf("%d", animal.ID),
			animal.Name,
			fmt.Sprintf("%d", animal.TypeID),
			fmt.Sprintf("%d", animal.Age),
			animal.Gender,
			fmt.Sprintf("%d", animal.Health),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
func SaveAnimalsCSVHandler(w http.ResponseWriter, r *http.Request) {
	var animals []Animal
	if err := db.Find(&animals).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := SaveAnimalsToCSV("animals.csv", animals); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println("Saved Animals to CSV")
}

func SaveFoodTypesToCSV(filename string, foodTypes []FoodType) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"ID", "Name", "Description", "Amount", "HealthPoints", "AnimalID"}
	writer.Write(headers)

	for _, foodType := range foodTypes {
		record := []string{
			fmt.Sprintf("%d", foodType.ID),
			foodType.Name,
			foodType.Description,
			fmt.Sprintf("%d", foodType.Amount),
			fmt.Sprintf("%d", foodType.HealthPoints),
			fmt.Sprintf("%d", foodType.AnimalID),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
func SaveFoodTypesCSVHandler(w http.ResponseWriter, r *http.Request) {
	var foodTypes []FoodType
	if err := db.Find(&foodTypes).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := SaveFoodTypesToCSV("food_types.csv", foodTypes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println("Saved FoodTypes to CSV")
}
