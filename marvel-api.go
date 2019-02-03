package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	apiKey = flag.String("pub", "9d2d890d05392417fa76b24057ef7d8f", "Public API key")
	secret = flag.String("priv", "9de5f9d35c5d9d7d72a10668f384cfbcd1693108", "Private API secret")
)

type CharacterDataWrapper struct {
	Code      int                    `json:"code,omitempty"`
	Status    string                 `json:"status,omitempty"`
	ETag      string                 `json:"etag,omitempty"`
	Copyright string                 `json:"copyright,omitempty"`
	Data      CharacterDataContainer `json:"data,omitempty"`
}

type CharacterDataContainer struct {
	Offset  int         `json:"offset,omitempty"`
	Limit   int         `json:"limit,omitempty"`
	Total   int         `json:"total,omitempty"`
	Count   int         `json:"count,omitempty"`
	Results []Character `json:"results,omitempty"`
}

type Character struct {
	ID          int          `json:"id,omitempty"`
	Name        string       `json:"name,omitempty"`
	Description string       `json:"description,omitempty"`
	Modified    Date         `json:"modified,omitempty"`
	ResourceURI string       `json:"resourceURI,omitempty"`
	URLS        []Url        `json:"urls,omitempty"`
	Comics      ComicList    `json:"comics,omitempty"`
	Stories     StoryList    `json:"stories,omitempty"`
	Events		EventList    `json:"events,omitempty"`
	Series		SeriesList   `json:"series,omitempty"`
}

type Url struct {
	Type string `json:"type,omitempty"`
	URL  string `json:"url,omitempty"`
}

type ComicList struct {
	Available     int            `json:"available,omitempty"`
	Returned      int            `json:"returned,omitempty"`
	CollectionURI string         `json:"collectionURI,omitempty"`
	Items         []ComicSummary `json:"items,omitempty"`
}

type ComicSummary struct {
	ResourceURI string `json:"resourceURI,omitempty"`
	Name        string `json:"name,omitempty"`
}

type StoryList struct {
	Available     int            `json:"available,omitempty"`
	Returned      int            `json:"returned,omitempty"`
	CollectionURI string         `json:"collectionURI,omitempty"`
	Items         []StorySummary `json:"items,omitempty"`
}

type StorySummary struct {
	ResourceURI string `json:"resourceURI,omitempty"`
	Name        string `json:"name,omitempty"`
	Type		string `json:"type,omitempty"`
}

type EventList struct {
	Available     int            `json:"available,omitempty"`
	Returned      int            `json:"returned,omitempty"`
	CollectionURI string         `json:"collectionURI,omitempty"`
	Items         []EventSummary `json:"items,omitempty"`
}

type EventSummary struct {
	ResourceURI string `json:"resourceURI,omitempty"`
	Name        string `json:"name,omitempty"`
}

type SeriesList struct {
	Available     int              `json:"available,omitempty"`
	Returned      int              `json:"returned,omitempty"`
	CollectionURI string           `json:"collectionURI,omitempty"`
	Items         []SeriesSummary  `json:"items,omitempty"`
}

type SeriesSummary struct {
	ResourceURI string `json:"resourceURI,omitempty"`
	Name        string `json:"name,omitempty"`
}

// -----------------------

type Date string

const dateLayout = "2006-01-02T15:04:05-0700"

func (d Date) Parse() time.Time {
	t, err := time.Parse(dateLayout, string(d))

	if err != nil {
		panic(err)
	}
	return t
}

// -----------------------

func makeTS() string {
	date := time.Now().UnixNano() / int64(time.Millisecond)
	return strconv.FormatInt(date, 10)
}

func makeHash(ts, privateKey, publicKey string) string {
	keys := []byte(ts + privateKey + publicKey)
	byteHash := md5.Sum(keys)
	hash := hex.EncodeToString(byteHash[:])

	return hash
}

func searchParameters(limit, name, order string) string {
	if limit != "" {
		limit = "&limit=" + limit
	}

	if name != "" {
		name = "&nameStartsWith=" + strings.Replace(name, " ", "+", -1)
	}

	if order != "" {
		order = "&orderBy=" + order
	}

	return limit + name + order
}

func getConnection(publicKey, privateKey, searchParams string) *http.Response {
	ts := makeTS()
	hash := makeHash(ts, privateKey, publicKey)

	URL := "http://gateway.marvel.com/v1/public/characters?ts=" + ts + "&apikey=" + publicKey + "&hash=" + hash + searchParams
	// fmt.Println(URL)

	// resp, err := http.Get("http://example.com/")
	response, err := http.Get(URL)

	if err != nil {
		fmt.Println("Error al establecer la conexión :", err)
	}
	/*
		if response.StatusCode == 200 {
			fmt.Println("Conexión establecida!")
		} */

	return response
}

func getBody(response *http.Response) []byte {
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Printf("Error al obtener la data: %v\n", err)
		os.Exit(1)
	}

	return body
}

func getCharacters(body []byte) []Character {
	characterDataWrapper := CharacterDataWrapper{}
	err := json.Unmarshal(body, &characterDataWrapper)

	if err != nil {
		fmt.Println("Error al obtener los datos...")
		duration := time.Duration(3) * time.Second
		time.Sleep(duration)
		os.Exit(1)
	}

	return characterDataWrapper.Data.Results
}

func printCharacters(response *http.Response) {
	body := getBody(response)
	characters := getCharacters(body)

	counter := 0
	for k := range characters {
		counter++

		fmt.Println("NÚMERO:        ", counter)
		fmt.Println("ID:            ", characters[k].ID)
		fmt.Println("NOMBRE:        ", characters[k].Name)
		fmt.Println("DESCRIPCIÓN:   ", characters[k].Description)
		fmt.Println("MODIFICADO:    ", characters[k].Modified)
		fmt.Println("RESOURCE_URI:  ", characters[k].ResourceURI)
		fmt.Println("---------------------------------------------")
		fmt.Println("                   URLS")
		fmt.Println("---------------------------------------------")
		urldata := characters[k].URLS
		for url := range urldata {
			fmt.Println("··· ··· ···")
			fmt.Println("type: ", urldata[url].Type)
			fmt.Println("url: ", urldata[url].URL)
		}

		fmt.Println("---------------------------------------------")
		fmt.Println("                  COMICS")
		fmt.Println("---------------------------------------------")
		fmt.Println("Disponible: ", characters[k].Comics.Available)
		fmt.Println("Returned: ", characters[k].Comics.Returned)
		fmt.Println("CollectionURI: ", characters[k].Comics.CollectionURI)
		comics := characters[k].Comics.Items
		for comic := range comics {
			fmt.Println("··· ··· ···")
			fmt.Println("Nombre: ", comics[comic].Name)
			fmt.Println("ResourceURI: ", comics[comic].ResourceURI)
		}

		fmt.Println("---------------------------------------------")
		fmt.Println("                  HISTORIA")
		fmt.Println("---------------------------------------------")
		fmt.Println("Disponible: ", characters[k].Stories.Available)
		fmt.Println("Returned: ", characters[k].Stories.Returned)
		fmt.Println("CollectionURI: ", characters[k].Stories.CollectionURI)
		stories := characters[k].Stories.Items
		for story := range stories {
			fmt.Println("··· ··· ···")
			fmt.Println("Nombre: ", stories[story].Name)
			fmt.Println("Tipo: ", stories[story].Type)
			fmt.Println("ResourceURI: ", stories[story].ResourceURI)
		}

		fmt.Println("---------------------------------------------")
		fmt.Println("                  EVENTOS")
		fmt.Println("---------------------------------------------")
		fmt.Println("Disponible: ", characters[k].Events.Available)
		fmt.Println("Returned: ", characters[k].Events.Returned)
		fmt.Println("CollectionURI: ", characters[k].Events.CollectionURI)
		events := characters[k].Events.Items
		for event := range events {
			fmt.Println("··· ··· ···")
			fmt.Println("ResourceURI: ", events[event].ResourceURI)
			fmt.Println("Nombre: ", events[event].Name)
		}

		fmt.Println("---------------------------------------------")
		fmt.Println("                  SERIES")
		fmt.Println("---------------------------------------------")
		fmt.Println("Disponible: ", characters[k].Series.Available)
		fmt.Println("Returned: ", characters[k].Series.Returned)
		fmt.Println("CollectionURI: ", characters[k].Series.CollectionURI)
		theseries := characters[k].Series.Items
		for series := range theseries {
			fmt.Println("··· ··· ···")
			fmt.Println("ResourceURI: ", theseries[series].ResourceURI)
			fmt.Println("Nombre: ", theseries[series].Name)
		}

		fmt.Println("---------------------------------------------")
		fmt.Println("INFORMACIÓN: Este es el resultado nro ", counter)
		fmt.Println("---------------------------------------------")


		fmt.Println()
		fmt.Println("******************************************************************************************")
		fmt.Println("//////////////////////////////////////////////////////////////////////////////////////////")
		fmt.Println("******************************************************************************************")
		fmt.Println()
	}
}

func getSysmode() string {
	scanner := bufio.NewScanner(os.Stdin)
	answer := ""
	sysmode := ""

	state := 0
	for state < 1 {
		fmt.Print(`
¿En qué consola estás corriendo esto?
  1. Consola de sistemas windows
  2. Consola de sistemas UNIX

  `)
		for scanner.Scan() {
			answer = scanner.Text()
			break
		}

		if answer == "1" {
			sysmode = "\r\n"
			state = 1

		} else if answer == "2" {
			sysmode = "\n"
			state = 1

		} else {
			fmt.Println("Por favor ingresa una respuesta válida")
			fmt.Println()
		}
	}


	//switch system := runtime.GOOS; system {
	//case "darwin":
	//	sysmode = "\n"
	//	fmt.Println("OS X...")
	//case "linux":
	//	sysmode = "\n"
	//	fmt.Println("Linux...")
	//default:
	//	sysmode = "\r\n"
	//	fmt.Println(system + "...")
	//}

	return sysmode
}

func getKeys(sysmode string) [2]string {
	reader := bufio.NewReader(os.Stdin)

	flag.Parse()

	publicKey := ""
	privateKey := ""

	state := 0
	for state < 1 {
		fmt.Println("")
		fmt.Print("¿Deseas ingresar tus propias apikeys?(y/N): ")
		preference, _ := reader.ReadString('\n')
		preference = strings.ToLower(strings.Replace(preference, sysmode, "", -1))

		if preference == "y" {
			fmt.Print("Ingresa tu llave privada: ")
			privateKey, _ = reader.ReadString('\n')
			privateKey = strings.Replace(privateKey, sysmode, "", -1)

			fmt.Print("Ingresa tu llave pública: ")
			publicKey, _ = reader.ReadString('\n')
			publicKey = strings.Replace(publicKey, sysmode, "", -1)

			state = 1

		} else if preference == "n" || preference == "" {

			if *apiKey != "" || *secret != "" {
				publicKey = *apiKey
				privateKey = *secret
				state = 1

			} else {
				fmt.Println("Error: Revisa tus credenciales en el código")
				fmt.Println()
			}

		} else {
			fmt.Println("Ingresa una respuesta válida")
			fmt.Println()
		}
	}
	var keys [2]string
	keys[0] = publicKey
	keys[1] = privateKey

	return keys
}

func getParamsExtra(sysmode string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Println("1. Buscar por nombre")
	fmt.Println("2. Listar")
	fmt.Println()
	option, _ := reader.ReadString('\n')
	option = strings.Replace(option, sysmode, "", -1)

	params := searchParameters("", "", "")

	if option == "1" {
		fmt.Print("Escribe el nombre de tu personaje favorito: ")
		character, _ := reader.ReadString('\n')
		character = strings.Replace(character, sysmode, "", -1)

		fmt.Println()
		fmt.Println("******************************************************************************************")
		fmt.Println("//////////////////////////////////////////////////////////////////////////////////////////")
		fmt.Println("******************************************************************************************")
		fmt.Println()

		params = searchParameters("1", character, "")

	} else if option == "2" {
		fmt.Println("Listando los primeros 20 personajes ordenados por nombre...")

		fmt.Println()
		fmt.Println("******************************************************************************************")
		fmt.Println("//////////////////////////////////////////////////////////////////////////////////////////")
		fmt.Println("******************************************************************************************")
		fmt.Println()

		params = searchParameters("20", "", "name")
	} else {
		fmt.Println("Iniciando búsqueda por defecto...")

		fmt.Println()
		fmt.Println("******************************************************************************************")
		fmt.Println("//////////////////////////////////////////////////////////////////////////////////////////")
		fmt.Println("******************************************************************************************")
		fmt.Println()
		fmt.Println()
	}

	return params
}

func main() {

	sysmode := getSysmode()
	keys := getKeys(sysmode)

	params := getParamsExtra(sysmode)
	response := getConnection(keys[0], keys[1], params)
	printCharacters(response)

	reader := bufio.NewReader(os.Stdin)

	state := 0
	for state < 2 {
		fmt.Println()
		fmt.Print("¿Deseas realizar otra búsqueda?(Y/n)")
		fmt.Println()
		answer, _ := reader.ReadString('\n')
		answer = strings.ToLower(strings.Replace(answer, sysmode, "", -1))

		if answer == "y" || answer == "" {
			params := getParamsExtra(sysmode)
			response := getConnection(keys[0], keys[1], params)
			printCharacters(response)

		} else if answer == "n" {
			fmt.Println("Gracias por probar...")
			duration := time.Duration(3) * time.Second
			time.Sleep(duration)
			state = 2

		} else {
			fmt.Println("Ingresa una respuesta válida")
			fmt.Println()
		}
	}
}
