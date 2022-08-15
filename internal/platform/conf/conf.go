package conf

import (
	"encoding/json"
	"log"
	"os"

	"github.com/sherpaurgen/garagesale/internal/platform/database"
)

//	type DBconfig struct {
//		Host       string `json:"Host"`
//		Port       int    `json:"Port"`
//		User       string `json:"User"`
//		Pass       string `json:"Pass"`
//		DBname     string `json:"DBname"`
//		DisableTLS bool   `json:"DisableTLS,omitempty"`
//	}
type Webconfig struct {
	Addr         string `json:"Addr"`
	ReadTimeout  int    `json:"ReadTimeout"`
	WriteTimeout int    `json:"WriteTimeout"`
}

func GetWebConfig() Webconfig {
	var basePath string
	basePath = getUserHomeDir()
	webconfigFile := "/.triage/web.json"
	webjsonFilePath := basePath + webconfigFile
	webdata, err := os.ReadFile(webjsonFilePath)
	if err != nil {
		log.Fatal("Error reading config file:", err)
	}
	var webc Webconfig
	json.Unmarshal(webdata, &webc)
	return webc
}

func GetDbConfig() database.DBconfig {
	basePath := getUserHomeDir()

	dbconfigFile := "/.triage/db.json"

	dbjsonFilePath := basePath + dbconfigFile

	dbdata, err := os.ReadFile(dbjsonFilePath)

	if err != nil {
		log.Fatal("Error reading config file:", err)
	}

	var dbc database.DBconfig

	json.Unmarshal(dbdata, &dbc)
	// dbc = data.dbconfig
	// webc = data.webserverconfig
	log.Println(dbc)
	return dbc
}

func getUserHomeDir() string {
	userHomePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error parsing json files in userHomedir : ", err)
	}
	return userHomePath
}
