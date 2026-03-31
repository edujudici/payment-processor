package dto

import "boleto-cancel/internal/boleto/usecase"

type Request struct {
	EndToEndId  string `json:"endToEndId"`
	TxId        string `json:"txid"`
	Valor       string `json:"valor"`
	Horario     string `json:"horario"`
	InfoPagador string `json:"infoPagador"`
	Chave       string `json:"chave"`
	DocCanal    string `json:"docCanal"`
}

func (r Request) ToUsecaseInput() usecase.BoletoRequest {
	return usecase.BoletoRequest{
		EndToEndId:  r.EndToEndId,
		TxId:        r.TxId,
		Valor:       r.Valor,
		Horario:     r.Horario,
		InfoPagador: r.InfoPagador,
		Chave:       r.Chave,
		DocCanal:    r.DocCanal,
	}
}
