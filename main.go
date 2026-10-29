package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

type Pagina struct {
	Titulo  string
	Mensaje string
	Imagen  string
}

var plantilla = template.Must(
	template.ParseFiles("templates/index.html"),
)

func inicio(w http.ResponseWriter, r *http.Request) {
	datos := Pagina{
		Titulo:  "Producto 1 - Entorno de desarrollo",
		Mensaje: "Soy alumno de la UOC",
		Imagen:  "/static/img/uoc.png",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := plantilla.Execute(w, datos); err != nil {
		log.Printf("Error al generar la página: %v", err)
		http.Error(
			w,
			"Error al generar la página",
			http.StatusInternalServerError,
		)
	}
}

func main() {
	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "8080"
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle(
		"/static/",
		http.StripPrefix("/static/", fs),
	)

	http.HandleFunc("/", inicio)

	log.Printf(
		"Servidor escuchando en http://localhost:%s",
		puerto,
	)

	log.Fatal(http.ListenAndServe(":"+puerto, nil))
}
