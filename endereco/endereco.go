package endereco

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/brenddonanjos/locacaolivro/banco"
)

type Endereco struct {
	Id         int    `json:"id"`
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Numero     string `json:"numero"`
	Bairro     string `json:"bairro"`
	Cidade     string `json:"cidade"`
	Uf         string `json:"uf"`
}

func EnderecoNovo(w http.ResponseWriter, r *http.Request) int64 {
	//instancia banco
	db := banco.Start()
	defer db.Close()

	//Pega os dados da requisição e salva nas struct
	reqBody, _ := ioutil.ReadAll(r.Body)
	var endereco Endereco
	err := json.Unmarshal(reqBody, &endereco)
	if err != nil {
		log.Fatal(err)
	}

	//Cria statement e salva no banco
	stmt, _ := db.Prepare("insert into enderecos (cep, logradouro, numero, bairro, cidade, uf) values (?, ?, ?, ?, ?, ?)")
	resp, err := stmt.Exec(endereco.Cep, endereco.Logradouro, endereco.Numero, endereco.Bairro, endereco.Cidade, endereco.Uf)
	if err != nil {
		log.Fatal(err)
	}
	_enderecoid, _ := resp.LastInsertId()

	return _enderecoid
}
