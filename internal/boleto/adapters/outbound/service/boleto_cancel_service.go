package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type BoletoCancelInput struct {
	CodBanco              string       `json:"codBanco"`
	CnpjClienteAccesstage string       `json:"cnpjClienteAccesstage"`
	NumeroConvenio        int          `json:"numeroConvenio"`
	NossoNumero           string       `json:"nossoNumero"`
	Operacao              string       `json:"operacao"`
	DataVencimentoTitulo  string       `json:"dataVencimentoTitulo"`
	Beneficiario          Beneficiario `json:"beneficiario"`
}

type Beneficiario struct {
	AgenciaBeneficiario    int     `json:"agenciaBeneficiario"`
	ContaBeneficiario      int     `json:"contaBeneficiario"`
	DigitoVerificadorConta int     `json:"digitoVerificadorConta"`
	Nome                   string  `json:"nome"`
	NumeroInscricao        string  `json:"numeroInscricao"`
	MensagemBeneficiario   *string `json:"mensagemBeneficiario"`
}

type BoletoCancelOutput struct {
	Registrado      bool    `json:"registrado"`
	DataHoraRetorno string  `json:"dataHoraRetorno"`
	SituacaoTitulo  *string `json:"situacaoTitulo"`
	LinhaDigitavel  *string `json:"linhaDigitavel"`
	CodigoBarras    *string `json:"codigoBarras"`
	NumeroTitulo    *string `json:"numeroTitulo"`
	NossoNumero     *string `json:"nossoNumero"`
	URL             *string `json:"url"`
	QRCode          *string `json:"qrCode"`
	Erros           []Erro  `json:"erros"`
}

type Erro struct {
	CodigoRejeicao         string `json:"codigoRejeicao"`
	MensagemRetorno        string `json:"mensagemRetorno"`
	MensagemDetalheRetorno string `json:"mensagemDetalheRetorno"`
	Campo                  string `json:"campo"`
}

type Boleto struct {
	HTTPClient *http.Client
}

func NewBoletoCancelService(client *http.Client) *Boleto {
	if client == nil {
		client = &http.Client{}
	}
	return &Boleto{HTTPClient: client}
}

func (p *Boleto) BoletoCancel(token string, req *BoletoCancelInput) (int, *BoletoCancelOutput, error) {
	url := os.Getenv("CANCEL_BOLETO_URL")
	if url == "" {
		return http.StatusInternalServerError, nil, fmt.Errorf("required environment variables not found to cancel boleto")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	log.Printf(
		"cancel boleto request | nossoNumero=%s | request=%s",
		req.NossoNumero,
		body,
	)

	httpReq, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("error creating request: %v", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := p.HTTPClient.Do(httpReq)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("error executing request: %v", err)
	}
	defer httpResp.Body.Close()

	responseBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("error reading response body: %v", err)
	}

	log.Printf("cancel boleto response code: %d and body: %s", httpResp.StatusCode, string(responseBody))

	if httpResp.StatusCode == http.StatusOK || httpResp.StatusCode == http.StatusBadRequest {
		var boletoResponse BoletoCancelOutput
		if err := json.Unmarshal(responseBody, &boletoResponse); err != nil {
			return http.StatusInternalServerError, nil, fmt.Errorf("error parsing JSON: %v", err)
		}

		return httpResp.StatusCode, &boletoResponse, nil
	}

	return http.StatusInternalServerError, nil, fmt.Errorf("error boleto response. Code: %d and ResponseBody: %s", httpResp.StatusCode, string(responseBody))
}
