package main

import (
	"context"
	"encoding/json"
	"fmt" // пакет для форматированного ввода вывода
	"regexp"
	"server/config"
	"server/pkg/postgres"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"

	//"github.com/pkg/errors"
	"io/ioutil"
	"log"      // пакет для логирования
	"net/http" // пакет для поддержки HTTP протокола
	"strings"  // пакет для работы с  UTF-8 строками
)

var pool *pgxpool.Pool
var reg *regexp.Regexp

func Error(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	_, _ = fmt.Fprintln(w, fmt.Sprintf(`{"message":"%s"}`, message))
}

func HomeRouterHandler(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json; charset=utf-8") // ну понятно, кодировка
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	res.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
	res.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
	res.Header().Set("X-Content-Type-Options", "nosniff")

	if req.Method == "OPTIONS" {
		return
	}

	if strings.Contains(req.URL.Path, ".") {
		Error(res, "Method not found", http.StatusNotFound)
		return
	}

	path :=
		reg.ReplaceAllString(
			strings.Replace(
				strings.TrimPrefix(req.URL.Path, "/"), // убираем первый слеш
				"/", "_", -1),                         // остальные заменяем на _ (-1 == реплейс на все вхождения)
			"") // убираем всё кроме цифр и букв и _

	splittedPath := strings.Split(path, "_")

	schemaName := fmt.Sprintf(splittedPath[0])
	functionName := fmt.Sprintf(strings.Join(splittedPath[1:], "_"))

	if functionName == "" {
		Error(res, "Function name is missing", http.StatusNotFound)
		return
	}
	//log.Println("schema name", schemaName)
	//log.Println("function name", functionName)

	auth_token := req.Header.Get("Authorization")
	var token string
	var tokenRef *string

	if strings.HasPrefix(auth_token, "Token") {
		token = auth_token[6:]
		tokenRef = &token
	}

	err := req.ParseForm() //анализ аргументов,
	if err != nil {
		panic(err)
	}

	queryParams := req.URL.Query()

	fixedParams := make(map[string]interface{}) // здесь будут храниться все параметры

	for k, v := range queryParams {
		if len(v) == 1 {
			fixedParams[k] = v[0] // если значение одно, вытаскиваем из массива
		} else {
			fixedParams[k] = v // если несколько — оставляем массивом
		}
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var postParams map[string]interface{}
	// пробуем получить json из body
	if len(body) > 0 {
		if err = json.Unmarshal(body, &postParams); err != nil {
			log.Printf("BODY DECODE FAILED: %s\n%s", string(body), err)
			Error(res, "Body decode failed. Expected json OBJECT", http.StatusBadRequest)
			return
		} else {
			for k, v := range postParams {
				fixedParams[k] = v
			}
		}
	}

	jsonString, err := json.Marshal(fixedParams)

	query := fmt.Sprintf("select %s.%s($1::json, $2::uuid);", schemaName, functionName) // формируем строку запроса
	var data string                                                                     // переменная для результата

	conn, err := pool.Acquire(req.Context())
	if err != nil {
		fmt.Println("Error acquiring connection from pool:", err)
		return
	}
	defer conn.Release()

	err = conn.QueryRow(context.Background(), query, string(jsonString), tokenRef).Scan(&data)

	log.Printf("query: %s\n", strings.Replace(strings.Replace(query, "$1", "'"+string(jsonString)+"'", -1), "$2", fmt.Sprintf("'%v'", token), -1))
	if err != nil {
		log.Printf("%s", err)
		if dberr, ok := err.(pgx.PgError); ok {
			status := http.StatusBadRequest
			var message string
			switch dberr.Code {
			case "42883":
				status = http.StatusNotFound
				message = "MethodNotFound"
			case "3F000":
				status = http.StatusNotFound
				message = fmt.Sprintf(`/%s/ api does not exists`, schemaName)
			case "ER401":
				status = http.StatusUnauthorized
				message = dberr.Message
			case "ER403":
				status = http.StatusForbidden
				message = dberr.Message
			case "P0001":
				message = dberr.Message
			case "XX000":
				message = dberr.Message
			default:
				log.Printf("%s: %s\n", dberr.Code, dberr.Message)
				message = "Unhandled Exception"
			}
			Error(res, message, status)
		} else {
			log.Println(err.Error())
			Error(res, "System Error", http.StatusInternalServerError)
		}

		return
	}
	_, _ = fmt.Fprintf(res, data) // отдаём данные в поток writer'а
}

func main() {
	ctx := context.Background()

	config := config.ParseConfigurationFile()
	pool = postgres.OpenPoolConnection(ctx, config)
	defer pool.Close()

	var err error
	if err != nil {
		panic(err)
	}
	reg, err = regexp.Compile("[^a-zA-Z0-9_]+") // компилируем regexp для замены всего опасного в имени функции

	http.HandleFunc("/", HomeRouterHandler) // установим роутер
	fmt.Println("Server started on 9999")
	err = http.ListenAndServe(":9999", nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
