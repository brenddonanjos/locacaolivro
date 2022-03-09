package livro

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/brenddonanjos/locacaolivro/banco"
)

type Livro struct {
	Id      int    `json:"id"`
	Titulo  string `json:"titulo"`
	Autor   string `json:"autor"`
	Ano     int    `json:"ano"`
	Edicao  string `json:"edicao"`
	Editora string `json:"editora"`
}

func LivroHandler(w http.ResponseWriter, r *http.Request) {
	sid := strings.TrimPrefix(r.URL.Path, "/livro/") //isola o id e salva

	id := 0
	if sid != "" {
		id, _ = strconv.Atoi(sid)
	}

	switch {
	case r.Method == "GET" && id != 0:
		livro := LivroPorId(id)
		json, _ := json.Marshal(livro)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(json))
	case r.Method == "GET":
		livroTodos(w)
	case r.Method == "POST":
		livroNovo(w, r)
	case r.Method == "PUT":
		livroEditar(w, r)
	case r.Method == "DELETE":
		livroExcluir(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, r)
	}
}

func LivroPorId(id int) Livro {
	db := banco.Start()
	defer db.Close()

	var livro Livro
	row := db.QueryRow("SELECT id, titulo, autor, ano, edicao, editora FROM livros WHERE id = ?", id)
	row.Scan(&livro.Id, &livro.Titulo, &livro.Autor, &livro.Ano, &livro.Edicao, &livro.Editora)

	return livro
}

func livroTodos(w http.ResponseWriter) {
	db := banco.Start()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM livros ORDER BY id DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var lista []Livro
	//percorre row e adiciona na lista
	for rows.Next() {
		var livro Livro
		rows.Scan(&livro.Id, &livro.Titulo, &livro.Autor, &livro.Ano, &livro.Edicao, &livro.Editora)
		lista = append(lista, livro)
	}

	json, _ := json.Marshal(lista)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}

func livroNovo(w http.ResponseWriter, r *http.Request) {
	//instancia banco
	db := banco.Start()
	defer db.Close()

	//Pega os dados da requisição e salva nas struct
	reqBody, _ := ioutil.ReadAll(r.Body)
	var livro Livro
	err1 := json.Unmarshal(reqBody, &livro)
	if err1 != nil {
		log.Fatal(err1)
	}

	//Cria statement e salva no banco
	stmt, _ := db.Prepare("INSERT INTO livros (titulo, autor, ano, edicao, editora) values (?, ?, ?, ?, ?)")
	_, err := stmt.Exec(livro.Titulo, livro.Autor, livro.Ano, livro.Edicao, livro.Editora)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "Livro cadastrado com suceso !")

}

func livroEditar(w http.ResponseWriter, r *http.Request) {
	db := banco.Start()
	defer db.Close()

	//Pega os dados da requisição e salva nas struct
	reqBody, _ := ioutil.ReadAll(r.Body)
	var livro Livro
	err1 := json.Unmarshal(reqBody, &livro)
	if err1 != nil {
		log.Fatal(err1)
	}
	stmt, _ := db.Prepare("UPDATE livros SET titulo = ?, autor = ?, ano = ?, edicao = ?, editora = ? WHERE id = ?")
	_, err := stmt.Exec(livro.Titulo, livro.Autor, livro.Ano, livro.Edicao, livro.Editora, livro.Id)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "Cliente "+livro.Titulo+" alterado com suceso !")
}

func livroExcluir(w http.ResponseWriter, r *http.Request) {
	//instancia banco
	db := banco.Start()
	defer db.Close()

	//Pega os dados da requisição e salva nas struct
	reqBody, _ := ioutil.ReadAll(r.Body)
	var livro Livro
	err1 := json.Unmarshal(reqBody, &livro)
	if err1 != nil {
		log.Fatal(err1)
	}

	//Cria statement e salva no banco
	stmt, _ := db.Prepare("DELETE FROM livros WHERE id = ?")
	_, err := stmt.Exec(livro.Id)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "Livro Apagado !")

}
