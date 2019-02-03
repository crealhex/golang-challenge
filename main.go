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
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	// Constante para quitar el fin de línea que trae por defecto tu consola
	inputrefactor = "\r\n"
)

var (
	// Inserta aquí tus propias credenciales de la API
	apiKey = flag.String("pub", "9d2d890d05392417fa76b24057ef7d8f", "Public API key")
	secret = flag.String("priv", "9de5f9d35c5d9d7d72a10668f384cfbcd1693108", "Private API secret")
)

// CharacterDataWrapper provee los estados e información de la conexión
type CharacterDataWrapper struct {
	Code      int                    `json:"code,omitempty"`
	Status    string                 `json:"status,omitempty"`
	ETag      string                 `json:"etag,omitempty"`
	Copyright string                 `json:"copyright,omitempty"`
	Data      CharacterDataContainer `json:"data,omitempty"`
}

// CharacterDataContainer provee estados de la consulta
type CharacterDataContainer struct {
	Offset  int         `json:"offset,omitempty"`
	Limit   int         `json:"limit,omitempty"`
	Total   int         `json:"total,omitempty"`
	Count   int         `json:"count,omitempty"`
	Results []Character `json:"results,omitempty"`
}

// Character provee información detallada de cada personaje
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

// Url provee los datos de todas las URLs de un personaje
type Url struct {
	Type string `json:"type,omitempty"`
	URL  string `json:"url,omitempty"`
}

// ComicLists provee los comics en un array e info de la consulta
type ComicList struct {
	Available     int            `json:"available,omitempty"`
	Returned      int            `json:"returned,omitempty"`
	CollectionURI string         `json:"collectionURI,omitempty"`
	Items         []ComicSummary `json:"items,omitempty"`
}

// ComicSummary provee los datos de cada comic de un personaje
type ComicSummary struct {
	ResourceURI string `json:"resourceURI,omitempty"`
	Name        string `json:"name,omitempty"`
}

// StoryList provee las historias en un array e info de la consulta
type StoryList struct {
	Available     int            `json:"available,omitempty"`
	Returned      int            `json:"returned,omitempty"`
	CollectionURI string         `json:"collectionURI,omitempty"`
	Items         []StorySummary `json:"items,omitempty"`
}

// StorySummary provee los datos de cada historia de un personaje
type StorySummary struct {
	ResourceURI string `json:"resourceURI,omitempty"`
	Name        string `json:"name,omitempty"`
	Type		string `json:"type,omitempty"`
}

// EventList provee los eventos en un array e info de la consulta
type EventList struct {
	Available     int            `json:"available,omitempty"`
	Returned      int            `json:"returned,omitempty"`
	CollectionURI string         `json:"collectionURI,omitempty"`
	Items         []EventSummary `json:"items,omitempty"`
}

// EventSummary provee los datos de cada evento de un personaje
type EventSummary struct {
	ResourceURI string `json:"resourceURI,omitempty"`
	Name        string `json:"name,omitempty"`
}

// SeriesList provee las series en un array e info de la consulta
type SeriesList struct {
	Available     int              `json:"available,omitempty"`
	Returned      int              `json:"returned,omitempty"`
	CollectionURI string           `json:"collectionURI,omitempty"`
	Items         []SeriesSummary  `json:"items,omitempty"`
}

// SeriesSummary provee los datos de cada serie de un personaje
type SeriesSummary struct {
	ResourceURI string `json:"resourceURI,omitempty"`
	Name        string `json:"name,omitempty"`
}

/*
 * Esto establece el formato de tiempo que usa la API
 * Nos sirve para ser usado en la estructura Character
 * para compatiblizar el valor de Modified
 */
type Date string
const dateLayout = "2006-01-02T15:04:05-0700"
func (d Date) Parse() time.Time {
	t, err := time.Parse(dateLayout, string(d))

	if err != nil {
		panic(err)
	}
	return t
}


/*
 * Las siguientes funciones nos ayudan
 * a generar una URL correcta para poder establecer
 * un 200 StatusCode en la conexión
 */
func makeTimestamp() string {
	date := time.Now().UnixNano() / int64(time.Millisecond)
	return strconv.FormatInt(date, 10)
}

func makeHash(ts, privateKey, publicKey string) string {
	keys := []byte(ts + privateKey + publicKey)
	byteHash := md5.Sum(keys)
	hash := hex.EncodeToString(byteHash[:])

	return hash
}

/*
 * Esta particular función ayuda a agregar
 * parámetros de búsqueda en la URL
 */
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

// Crea la conexión con la Marvel API
func getConnection(publicKey, privateKey, searchParams string) *http.Response {
	ts := makeTimestamp()
	hash := makeHash(ts, privateKey, publicKey)
	URL := "http://gateway.marvel.com/v1/public/characters?ts=" + ts + "&apikey=" + publicKey + "&hash=" + hash + searchParams

	response, err := http.Get(URL)

	if err != nil {
		fmt.Println("Error al establecer la conexión :", err)
	}

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

//////
// Characters
/////
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

/*
 * Imprime todos los datos consultados en consola
 * Son la gran mayoría de datos que ofrece la API
 */
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
		fmt.Println("***************************************************************************")
		fmt.Println("///////////////////////////////////////////////////////////////////////////")
		fmt.Println("***************************************************************************")
		fmt.Println()
	}
}

/////
// JUST FOR FUN
func getSysmode() {
	fmt.Print("Corriendo en ")

	switch system := runtime.GOOS; system {
	case "darwin":
		fmt.Println("OS X...")
	case "linux":
		fmt.Println("Linux...")
	default:
		fmt.Println(system + "...")
	}
}

/*
 * Servible para saber la preferencia del usuario
 * ingresar tokens por consola o usar las tokens's flag
 */
func getKeys() [2]string {
	reader := bufio.NewReader(os.Stdin)

	flag.Parse()

	publicKey := ""
	privateKey := ""

	state := 0
	for state < 1 {
		fmt.Println("")
		fmt.Print("¿Deseas ingresar tus propias apikeys?[y/N]: ")
		preference, _ := reader.ReadString('\n')
		preference = strings.TrimRight(preference, inputrefactor)

		if preference == "y" {
			fmt.Print("Ingresa tu llave privada: ")
			privateKey, _ = reader.ReadString('\n')
			privateKey = strings.TrimRight(privateKey, inputrefactor)

			fmt.Print("Ingresa tu llave pública: ")
			publicKey, _ = reader.ReadString('\n')
			publicKey = strings.TrimRight(publicKey, inputrefactor)

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

/*
 * Usado para generar las opciones de búsqueda,
 * estas establecen qué parámetros llevará
 * searchParameters()
 */
func getParamsExtra() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Println("1. Buscar por nombre")
	fmt.Println("2. Listar")
	fmt.Println()

	fmt.Print("Digita una opción: ")
	option, _ := reader.ReadString('\n')
	option = strings.TrimRight(option, inputrefactor)

	params := searchParameters("", "", "")

	if option == "1" {
		fmt.Print("Escribe el nombre de tu personaje favorito: ")
		character, _ := reader.ReadString('\n')
		character = strings.TrimRight(character, inputrefactor)

		fmt.Println()
		fmt.Println("***************************************************************************")
		fmt.Println("///////////////////////////////////////////////////////////////////////////")
		fmt.Println("***************************************************************************")
		fmt.Println()

		params = searchParameters("1", character, "")

	} else if option == "2" {
		fmt.Println("Listando los primeros 20 personajes ordenados por nombre...")

		fmt.Println()
		fmt.Println("***************************************************************************")
		fmt.Println("///////////////////////////////////////////////////////////////////////////")
		fmt.Println("***************************************************************************")
		fmt.Println()

		params = searchParameters("20", "", "name")
	} else {
		fmt.Println("Iniciando búsqueda por defecto...")

		fmt.Println()
		fmt.Println("***************************************************************************")
		fmt.Println("///////////////////////////////////////////////////////////////////////////")
		fmt.Println("***************************************************************************")
		fmt.Println()
		fmt.Println()
	}

	return params
}

/*
 * Aquí se incluye también un loop
 * para saber en qué momento terminar la ejecución
 * por decisión del usuario
 */
func main() {
	getSysmode()
	keys := getKeys()

	params := getParamsExtra()
	response := getConnection(keys[0], keys[1], params)
	printCharacters(response)

	reader := bufio.NewReader(os.Stdin)

	state := 0
	for state < 2 {
		fmt.Println()
		fmt.Print("¿Deseas realizar otra búsqueda?[Y/n]")
		fmt.Println()
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimRight(answer, inputrefactor)

		if answer == "y" || answer == "" {
			params := getParamsExtra()
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