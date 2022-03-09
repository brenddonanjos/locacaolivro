package locacao

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/brenddonanjos/locacaolivro/banco"
	"github.com/brenddonanjos/locacaolivro/cliente"
	"github.com/brenddonanjos/locacaolivro/livro"
)

type Locacao struct {
	Id          int             `json:"id"`
	Status      string          `json:"status"`
	DataLocacao string          `json:"data_locacao"`
	PrazoDias   int             `json:"prazo_dias"`
	ClienteId   int             `json:"cliente_id"`
	LivroId     int             `json:"livro_id"`
	Cliente     cliente.Cliente `json:"cliente"`
	Livro       livro.Livro     `json:"livro"`
}

func LocacaoHandler(w http.ResponseWriter, r *http.Request) {
	sid := strings.TrimPrefix(r.URL.Path, "/cliente/") //isola o id e salva

	id := 0
	if sid != "" {
		id, _ = strconv.Atoi(sid)
	}

	switch {
	case r.Method == "GET" && id != 0:
		locacao := locacaoPorId(id)
		json, _ := json.Marshal(locacao)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(json))
	case r.Method == "GET":
		locacaoTodos(w)
	case r.Method == "POST":
		locacaoNovo(w, r)
	case r.Method == "PUT":
		locacaoEditar(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, r)
	}
}

func locacaoPorId(id int) Locacao {
	db := banco.Start()
	defer db.Close()

	var locacao Locacao
	row := db.QueryRow("SELECT * FROM locacoes WHERE id = ? ", id)
	row.Scan(&locacao.Status, &locacao.DataLocacao, &locacao.PrazoDias, &locacao.ClienteId, &locacao.LivroId)
	locacao.Cliente = cliente.ClientePorId(locacao.ClienteId)
	locacao.Livro = livro.LivroPorId(locacao.LivroId)

	return locacao
}

func locacaoTodos(w http.ResponseWriter) {
	db := banco.Start()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM locacoes ORDER BY id DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var lista []Locacao
	//percorre row e adiciona na lista
	for rows.Next() {
		var locacao Locacao
		rows.Scan(&locacao.Id, &locacao.Status, &locacao.DataLocacao, &locacao.PrazoDias, &locacao.ClienteId, &locacao.LivroId)
		locacao.Cliente = cliente.ClientePorId(locacao.ClienteId)
		locacao.Livro = livro.LivroPorId(locacao.LivroId)
		lista = append(lista, locacao)
	}

	json, _ := json.Marshal(lista)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}

func locacaoNovo(w http.ResponseWriter, r *http.Request) {
	//instancia banco
	db := banco.Start()
	defer db.Close()

	//Pega os dados da requisição e salva nas struct
	reqBody, _ := ioutil.ReadAll(r.Body)
	var locacao Locacao
	err1 := json.Unmarshal(reqBody, &locacao)
	if err1 != nil {
		log.Fatal(err1)
	}

	//Cria statement e salva no banco
	stmt, _ := db.Prepare("INSERT INTO locacoes (status, data_locacao, prazo_dias, cliente_id, livro_id) values (?, ?, ?, ?, ?)")
	_, err := stmt.Exec(locacao.Status, locacao.DataLocacao, locacao.PrazoDias, locacao.ClienteId, locacao.LivroId)
	if err != nil {
		log.Fatal(err)
	}
	//_clienteid, _ := resp.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "Locação registrada com suceso !")

}
func locacaoEditar(w http.ResponseWriter, r *http.Request) {
	db := banco.Start()
	defer db.Close()

	//Pega os dados da requisição e salva nas struct
	reqBody, _ := ioutil.ReadAll(r.Body)
	var locacao Locacao

	err1 := json.Unmarshal(reqBody, &locacao)
	if err1 != nil {
		log.Fatal(err1)
	}

	stmt, _ := db.Prepare("UPDATE locacoes SET status = ?, prazo_dias = ? WHERE id = ?")
	_, err := stmt.Exec(locacao.Status, locacao.PrazoDias, locacao.Id)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "Locação atualizada com suceso !")
}
