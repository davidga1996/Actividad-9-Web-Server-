package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func form(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(res, cargarHtml("index.html"))
}

func cargarHtml(a string) string {
	html, _ := ioutil.ReadFile(a)
	return string(html)
}

//Aux ...
type Aux struct {
	Alumno       string
	Materia      string
	Calificacion float64
}

//Materias ...
var Materias = make(map[string]map[string]float64)

//Alumno ...
var Alumno = make(map[string]map[string]float64)

func agregar(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		n, _ := strconv.ParseFloat(req.FormValue("calificacion"), 64)
		alumno := Aux{Alumno: req.FormValue("alumno"), Calificacion: n, Materia: req.FormValue("materia")}
		AgregarCalificacion(alumno)
		imprimir()
		res.Header().Set(
			"Conten-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("agregar.html"),
		)
	case "GET":
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("general.html"),
			PromedioGeneral(),
		)
	}
}

//AgregarCalificacion ...
func AgregarCalificacion(alumno Aux) {
	_, flagMateria := Materias[alumno.Materia]               //Bool si la materia existe
	_, flagAlumno := Materias[alumno.Materia][alumno.Alumno] //bool si el alumno tiene calificación en la meteria
	_, flagAlumnoClases := Alumno[alumno.Alumno]             //si el alumno existe
	if flagMateria {
		if flagAlumno {
			fmt.Println("El alumno ya tiene calificación para la materia")
		} else {
			Materias[alumno.Materia][alumno.Alumno] = alumno.Calificacion
		}
		if flagAlumnoClases {
			Alumno[alumno.Alumno][alumno.Materia] = alumno.Calificacion
		} else {
			var materiaAux = make(map[string]float64)
			materiaAux[alumno.Materia] = alumno.Calificacion
			Alumno[alumno.Alumno] = materiaAux
		}
	} else {
		//Crear materia
		var alumnos = make(map[string]float64)
		alumnos[alumno.Alumno] = alumno.Calificacion
		Materias[alumno.Materia] = alumnos
		if flagAlumnoClases {
			Alumno[alumno.Alumno][alumno.Materia] = alumno.Calificacion
		} else {
			var materiaAux = make(map[string]float64)
			materiaAux[alumno.Materia] = alumno.Calificacion
			Alumno[alumno.Alumno] = materiaAux
		}
	}
}

func promedioAlumno(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		promedio := ObtenerPromedioAlumno(req.FormValue("alumno"))
		res.Header().Set(
			"Conten-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("promedio.html"),
			req.FormValue("alumno"),
			promedio,
		)
	}
}

//ObtenerPromedioAlumno ...
func ObtenerPromedioAlumno(nombre string) float64 {
	var promedio float64
	promedio = 0
	numMaterias := 0.0
	_, flagAlumno := Alumno[nombre]
	if flagAlumno {
		for _, calificación := range Alumno[nombre] {
			promedio += calificación
			numMaterias++
		}
		promedio = promedio / numMaterias
		return promedio
	}
	fmt.Println("El alumno no existe")
	return 0
}

func promedioMateria(res http.ResponseWriter, req *http.Request) {
	fmt.Println(req.Method)
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		promedio := PromedioMateria(req.FormValue("materia"))
		res.Header().Set(
			"Conten-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("promedioMateria.html"),
			req.FormValue("materia"),
			promedio,
		)
	}
}

// PromedioMateria ...
func PromedioMateria(materia string) float64 {
	_, ok := Materias[materia]
	if ok {
		promedio := 0.0
		numAlumnos := 0.0
		for _, cal := range Materias[materia] {
			promedio += cal
			numAlumnos++
		}
		return promedio / numAlumnos
	}
	return 0
}

//PromedioGeneral ...
func PromedioGeneral() float64 {
	alumnos := len(Alumno)
	var PromedioGeneral float64
	var numMaterias float64
	if alumnos > 0 {
		for materia := range Materias {
			numMaterias++
			promedioMateria := 0.0
			numAlumnos := 0.0
			for _, cal := range Materias[materia] {
				promedioMateria += cal
				numAlumnos++
			}
			PromedioGeneral += promedioMateria / numAlumnos
		}
		return PromedioGeneral / numMaterias
	}
	return 0
}

func imprimir() {
	fmt.Println(Alumno)
	fmt.Println(Materias)
}

func main() {
	http.HandleFunc("/index", form)
	http.HandleFunc("/agregar", agregar)
	http.HandleFunc("/promedioAlumno", promedioAlumno)
	http.HandleFunc("/promedioMateria", promedioMateria)
	fmt.Println("Corriendo servirdor de Calificaciones...")
	http.ListenAndServe(":9000", nil)
}
