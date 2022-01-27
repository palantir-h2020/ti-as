package backend

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"PalantirIpAnonymizationService/backend/db"
	"PalantirIpAnonymizationService/cryptopan"
	"PalantirIpAnonymizationService/logger"
	"PalantirIpAnonymizationService/utils"
)

type App struct {
	Router *mux.Router
}

type RequestBody struct {
	IpAddr string
}

// Initialize logger
var log = logger.New(false)
var cpan = initializeCryptoPan()
var dbController *db.RedisDb

func (a *App) Initialize(config *utils.Config) {
	dbController = initializeDbConnection(config)

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Info("Palanir Anonymization Service Running on " + addr + " (Press CTRL+C to quit)")
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", a.homeRoute).Methods("GET")
	a.Router.HandleFunc("/anonymize", a.anonymizeRoute).Methods("POST")
	a.Router.HandleFunc("/deanonymize", a.deanonymizeRoute).Methods("POST")
}

func (a *App) homeRoute(w http.ResponseWriter, r *http.Request) {
	log.Info("Endpoint Hit: Index (/)")

	response := make(map[string]string)
	response["result"] = "success"
	response["info"] = "Palantir IP Anonymization Service"
	response["status"] = "up"

	respondWithJSON(w, http.StatusOK, response)
}

func (a *App) anonymizeRoute(w http.ResponseWriter, r *http.Request) {
	log.Info("Endpoint Hit: Anonymize (/anonymize)")

	startTime := time.Now()

	// Get HTTP POST JSON body
	var reqBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Cannot anonymize IP address, IpAddr param is missing.")
		return
	}

	originalIp := reqBody.IpAddr
	if len(originalIp) == 0 {
		respondWithError(w, http.StatusBadRequest, "Cannot anonymize IP address, IpAddr param is missing.")
		return
	}

	obfIp := cpan.Anonymize(net.ParseIP(originalIp))

	log.Info("Anonymization: " + originalIp + " --> " + obfIp.String())

	res, err := dbController.Insert(obfIp.String(), originalIp)
	if err != nil {
		log.Error("Cannot insert IP pair(originalIp=" + originalIp + ", obfuscatedIp=" + obfIp.String() + ") to Redis database.")
		respondWithError(w, http.StatusBadRequest, "Cannot store IP pair. Connection with database cannot be established.")
		return
	}

	if !res {
		respondWithError(w, http.StatusBadRequest, "Cannot store IP pair. Result is false.")
		return
	}

	elapsed, _ := time.ParseDuration(time.Since(startTime).String())

	response := make(map[string]string)
	response["result"] = "success"
	response["originalIp"] = originalIp
	response["obfuscatedIp"] = obfIp.String()
	response["execution_time"] = fmt.Sprintf("%f", elapsed.Seconds())

	respondWithJSON(w, http.StatusOK, response)
}

func (a *App) deanonymizeRoute(w http.ResponseWriter, r *http.Request) {
	log.Info("Endpoint Hit: Deanonymize (/deanonymize)")

	startTime := time.Now()

	// Get HTTP POST JSON body
	var reqBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Cannot deanonymize IP address, IpAddr param is missing.")
		return
	}

	obfuscatedIp := reqBody.IpAddr
	if len(obfuscatedIp) == 0 {
		respondWithError(w, http.StatusBadRequest, "Cannot deanonymize IP address, IpAddr param is missing.")
		return
	}

	originalIp, err := dbController.Get(obfuscatedIp)
	if err != nil {
		log.Error("Cannot deanonymize IP " + obfuscatedIp + ".")
		respondWithError(w, http.StatusBadRequest, "Cannot deanonymize IP "+obfuscatedIp+". Connection with database cannot be established.")
		return
	}

	if len(originalIp) == 0 {
		log.Error("Cannot deanonymize IP " + obfuscatedIp + ".")
		respondWithError(w, http.StatusBadRequest, "Cannot deanonymize IP "+obfuscatedIp+". Cannot find obfuscated IP.")
		return
	}

	elapsed, _ := time.ParseDuration(time.Since(startTime).String())

	response := make(map[string]string)
	response["result"] = "success"
	response["originalIp"] = originalIp
	response["obfuscatedIp"] = obfuscatedIp
	response["execution_time"] = fmt.Sprintf("%f", elapsed.Seconds())

	respondWithJSON(w, http.StatusOK, response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func initializeCryptoPan() *cryptopan.Cryptopan {
	var encrKey = []byte{21, 34, 23, 141, 51, 164, 207, 128, 19, 10, 91, 22, 73, 144, 125, 16, 216, 152, 143, 131, 121, 121, 101, 39, 98, 87, 76, 45, 42, 132, 34, 2}
	cpan, err := cryptopan.New(encrKey)

	if err != nil {
		// Cannot initalize CryptoPan object. Print error
		log.Error("Error (Cannot initialize CryptoPan library): \n")
		log.PrintError(err)

		// Exit
		os.Exit(-1)
	}

	return cpan
}

func initializeDbConnection(config *utils.Config) *db.RedisDb {
	host := config.Database.Host
	port := config.Database.Port
	username := config.Database.Username
	password := config.Database.Password
	database := config.Database.Name

	// TODO: Read from configuration file. Also, read db type and load proper controller
	redisDb, err := db.New(host, port, username, password, database)
	if err != nil {
		log.Fatal(err)
	}

	return redisDb
}
