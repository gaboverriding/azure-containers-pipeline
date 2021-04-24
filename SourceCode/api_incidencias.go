package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// API REST - ENTIDAD Incidencia
type Incidencia struct {
	IdIncidencia string    `json:"id_incidencia"`
	Titulo       string    `json:"titulo"` 
	Descripcion  string    `json:"descripcion"`
	CreadaElDia  time.Time `json:"creada_el_dia"`
}

var datosIncidencias = make(map[string]Incidencia)
var id int

func main() {
	// ENRUTADOR DE GORILLA MUX
	gorilla_router := mux.NewRouter().StrictSlash(false)

	// MANEJADORES POR VERBOS HTTP
	gorilla_router.HandleFunc("/api/incidencias", GetNoteHandler).Methods("GET")
	gorilla_router.HandleFunc("/api/incidencias", PostNoteHandler).Methods("POST")
	gorilla_router.HandleFunc("/api/incidencias/{id}", PutNoteHandler).Methods("PUT")
	gorilla_router.HandleFunc("/api/incidencias/{id}", DeleteNoteHandler).Methods("DELETE")

	server := &http.Server{
		Addr:           ":8080",
		Handler:        gorilla_router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Go API Rest para la bitácora de incidencias escuchando en puerto  8080 ...")
	server.ListenAndServe()

}

// GetNoteHandler - GET - /api/incidencias
func GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Se invocó GET ...")
	// SLICE DE Incidencias
	var incidencias []Incidencia
	for _, valor := range datosIncidencias {
		incidencias = append(incidencias, valor)
	}
	// Set ESTABLECE CABECERAS HTTP
	w.Header().Set("Content-Type", "application/json")

	// PASAR LOS DATOS AL FORMATO JSON CON Marshall
	j, err := json.Marshal(incidencias)
	if err != nil {
		panic(err)
	}
	// ESCRIBIR LA RESPUESTA HTTP PARA EL USUARIO
	w.WriteHeader(http.StatusOK)

	// CONTENIDO  Y RESPUESTA EN FORMATO JSON
	w.Write(j)

}

// PostNoteHandler - POST - /api/incidencias
// { "title" : "Titulo que sea", "description" : "alguna descripcion"}
func PostNoteHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Se invocó POST ...")
	// Incidencia nueva
	var incidencia Incidencia
	// DECODIFICADOR PARA EL DATO DE ENTRADA Y VERIFICAR QUE EL JSON ENVIADO ES CORRECTO
	err := json.NewDecoder(r.Body).Decode(&incidencia)
	if err != nil {
		panic(err)
	}
	// SE AÑADE EL TIME STAMP A LA Incidencia
	incidencia.CreadaElDia = time.Now()
	id++
	k := strconv.Itoa(id)
	incidencia.IdIncidencia = k
	// SE AGREGA LA NUEVA INCIDENCIA AL ARREGLO DE INCIDENCIAS
	datosIncidencias[k] = incidencia

	// SE PREPARA LA RESPUESTA AL CLIENTE
	// SE ESTABLECE POR CABECERA EL TIPO DE RESPUESTA
	w.Header().Set("Content-Type", "application/json")

	// SE PASAN LOS DATOS AL FORMATO JSON CON Marshall
	j, err := json.Marshal(incidencia)
	if err != nil {
		panic(err)
	}
	// ESCRIBIR LA RESPUESTA HTTP PARA EL USUARIO
	w.WriteHeader(http.StatusCreated)

	// CONTENIDO  Y RESPUESTA EN FORMATO JSON
	w.Write(j)
}

// PutNoteHandler - PUT - /api/incidencias/id
func PutNoteHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Se invocó PUT ...")
	// SE RECUPERA EL PARAMETRO INDICADO, EN ESTE CASO "id"
	// DEVUELVE EN UN MAP DE STRING CUYO INDICE ES EL NOMBRE DEL PARAMETRO "id"
	vars := mux.Vars(r)

	// SE RECUPERA EL "id" INFORMADO
	k := vars["id"]

	// SE OBTIENEN LOS DATOS INFORMADOS EN EL PAYLOAD Y SE ASIGNAN A UNA ESTRUCTURA
	var incidenciaUpdate Incidencia
	err := json.NewDecoder(r.Body).Decode(&incidenciaUpdate)
	if err != nil {
		panic(err)
	}

	// SE REVISA SI EXISTE LA INCIDENCIA POR EL ID
	if incidencia, ok := datosIncidencias[k]; ok {
		// SE RECUPERA EL TIMESTAMP DE LA INCIDENCIA A ACTUALIZAR
		incidenciaUpdate.CreadaElDia = incidencia.CreadaElDia
		// SE RECUPERA EL ID DE LA INCIDENCIA A ACTUALIZAR
		incidenciaUpdate.IdIncidencia = incidencia.IdIncidencia
		// SE BORRA LA INCIDENCIA ACTUAL
		delete(datosIncidencias, k)
		// SE AGREGA UN NUEVO REGISTRO CON LOS DATOS NUEVOS
		datosIncidencias[k] = incidenciaUpdate
	} else {
		log.Printf("No se encontro la incidencia con el id: %d", k)
	}

	// SE MANDA LA RESPUESTA AL CLIENTE
	w.WriteHeader(http.StatusNoContent)

}

// DeleteNoteHandler - DELETE - /api/incidencias/id
func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Se invocó DELETE ...")
	// SE RECUPERAN LOS PARAMETROS, EN ESTE CASO EL "id"
	// DEVUELVE EN UN MAP DE STRING CUYO INDICE ES EL NOMBRE DEL PARAMETRO "id"
	vars := mux.Vars(r)
	k := vars["id"]

	// SE REVISA SI EXISTE LA INCIDENCIA
	if _, ok := datosIncidencias[k]; ok {
		delete(datosIncidencias, k)
	} else {
		log.Printf("No se encontro la incidencia con el id: %d", k)
	}

	// SE MANDA LA RESPUESTA AL CLIENTE
	w.WriteHeader(http.StatusNoContent)

}
