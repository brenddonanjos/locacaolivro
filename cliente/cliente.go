package cliente

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

type Cliente struct {
	Id         int    `json:"id"`
	Nome       string `json:"nome"`
	Email      string `json:"email"`
	Telefone   string `json:"telefone"`
	Cpf        string `json:"cpf"`
	Ativo      int    `json:"ativo"`
	EnderecoId string `json:"endereco_id"`
}

//ClienteHandler analisa o request e invoca a função adequada
func ClienteHandler(w http.ResponseWriter, r *http.Request) {
	sid := strings.TrimPrefix(r.URL.Path, "/cliente/") //isola o id e salva

	id := 0
	if sid != "" {
		id, _ = strconv.Atoi(sid)
	}

	switch {
	case r.Method == "GET" && id != 0:
		cli := ClientePorId(id)
		json, _ := json.Marshal(cli)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(json))
	case r.Method == "GET":
		clienteTodos(w, r)
	case r.Method == "POST":
		clienteNovo(w, r)
	case r.Method == "PUT":
		clienteEditar(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, r)
	}
}

func clienteTodos(w http.ResponseWriter, r *http.Request) {
	db := banco.Start()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM clientes WHERE ativo = ? ORDER BY id DESC", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var lista []Cliente
	//percorre row e adiciona na lista
	for rows.Next() {
		var cli Cliente
		rows.Scan(&cli.Id, &cli.Nome, &cli.Email, &cli.Telefone, &cli.Cpf, &cli.Ativo, &cli.EnderecoId)
		lista = append(lista, cli)
	}

	json, _ := json.Marshal(lista)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}

func ClientePorId(id int) Cliente {
	db := banco.Start()
	defer db.Close()

	var cli Cliente
	row := db.QueryRow("SELECT * FROM clientes WHERE id = ? and ativo = ?", id, 1)
	row.Scan(&cli.Id, &cli.Nome, &cli.Email, &cli.Telefone, &cli.Cpf, &cli.Ativo, &cli.EnderecoId)

	return cli
}

func clienteNovo(w http.ResponseWriter, r *http.Request) {
	//instancia banco
	db := banco.Start()
	defer db.Close()

	//Pega os dados da requisição e salva nas struct
	reqBody, _ := ioutil.ReadAll(r.Body)
	var cliente Cliente
	err1 := json.Unmarshal(reqBody, &cliente)
	if err1 != nil {
		log.Fatal(err1)
	}

	//Cria statement e salva no banco
	stmt, _ := db.Prepare("INSERT INTO clientes (nome, cpf, email, telefone, ativo) values (?, ?, ?, ?, ?)")
	_, err := stmt.Exec(cliente.Nome, cliente.Cpf, cliente.Email, cliente.Telefone, true)
	if err != nil {
		log.Fatal(err)
	}
	//_clienteid, _ := resp.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "Cliente salvo com suceso !")

}

func clienteEditar(w http.ResponseWriter, r *http.Request) {
	db := banco.Start()
	defer db.Close()

	//Pega os dados da requisição e salva nas struct
	reqBody, _ := ioutil.ReadAll(r.Body)
	var cliente Cliente
	err1 := json.Unmarshal(reqBody, &cliente)
	if err1 != nil {
		log.Fatal(err1)
	}
	stmt, _ := db.Prepare("UPDATE clientes SET nome = ?, email = ?, telefone = ?, cpf = ?, ativo = ? WHERE id = ?")
	_, err := stmt.Exec(cliente.Nome, cliente.Email, cliente.Telefone, cliente.Cpf, cliente.Ativo, cliente.Id)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "Cliente "+cliente.Nome+" alterado com suceso !")
}
