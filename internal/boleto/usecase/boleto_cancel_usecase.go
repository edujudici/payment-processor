package usecase

import (
	"boleto-cancel/internal/boleto/adapters/outbound/service"
	"boleto-cancel/internal/boleto/ports"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type BoletoRequest struct {
	EndToEndId  string
	TxId        string
	Valor       string
	Horario     string
	InfoPagador string
	Chave       string
	DocCanal    string
}

type SendMessageRequest struct {
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
	CodBanco        *string `json:"codBanco"`
	NumeroConvenio  *string `json:"numeroConvenio"`
}

type Erro struct {
	CodigoRejeicao         string `json:"codigoRejeicao"`
	MensagemRetorno        string `json:"mensagemRetorno"`
	MensagemDetalheRetorno string `json:"mensagemDetalheRetorno"`
	Campo                  string `json:"campo"`
}

type BoletoCancelInterface interface {
	Execute(ctx context.Context, request BoletoRequest) error
}

type BoletoCancelUseCase struct {
	authService            ports.AuthInterface
	boletoCancelService    ports.BoletoCancelInterface
	sendMessageService     ports.SendMessageInterface
	boletoCancelRepository ports.BoletoCancelRepository
}

func NewBoletoCancelUseCase(
	authService ports.AuthInterface,
	boletoCancelService ports.BoletoCancelInterface,
	sendMessageService ports.SendMessageInterface,
	boletoCancelRepository ports.BoletoCancelRepository,
) *BoletoCancelUseCase {
	return &BoletoCancelUseCase{
		authService:            authService,
		boletoCancelService:    boletoCancelService,
		sendMessageService:     sendMessageService,
		boletoCancelRepository: boletoCancelRepository,
	}
}

func (uc *BoletoCancelUseCase) Execute(ctx context.Context, request BoletoRequest) error {
	log.Print("get access token to consume services")
	token, err := uc.authService.GetAccessToken()
	if err != nil {
		return err
	}

	v := strings.Split(request.DocCanal, "|")
	if len(v) < 5 {
		return fmt.Errorf("invalid text, does not contain all expected fields")
	}

	dueDate := v[0]
	codBanco := v[1]
	convenio := v[2]
	ourNumber := v[3]

	numConv, err := strconv.Atoi(convenio)
	if err != nil {
		return fmt.Errorf("error converting convenio to int: %v", err)
	}

	dateConv, err := convertDate(dueDate)
	if err != nil {
		return fmt.Errorf("error converting date: %v", err)
	}

	bank, err := uc.boletoCancelRepository.GetBankByAgreement(
		ctx,
		convenio,
	)
	if err != nil {
		return err
	}

	log.Print("get business data with id: ", bank.Company)
	business, err := uc.boletoCancelRepository.GetBusiness(
		ctx,
		bank.Company,
	)
	if err != nil {
		return err
	}

	input := &service.BoletoCancelInput{
		CodBanco:              bank.NameBank,
		CnpjClienteAccesstage: business.CNPJ,
		NumeroConvenio:        numConv,
		NossoNumero:           ourNumber,
		Operacao:              "BAIXAR",
		DataVencimentoTitulo:  dateConv,
		Beneficiario: service.Beneficiario{
			AgenciaBeneficiario:    bank.Agency,
			ContaBeneficiario:      bank.Account,
			DigitoVerificadorConta: bank.DigitAccount,
			Nome:                   business.Name,
			NumeroInscricao:        business.CNPJ,
			MensagemBeneficiario:   nil,
		},
	}

	log.Printf(
		"boleto cancellation request | txId=%s | nossoNumero=%s",
		request.TxId,
		input.NossoNumero,
	)

	_, response, err := uc.boletoCancelService.BoletoCancel(*token, input)
	if err != nil {
		return err
	}

	err = uc.Send(ctx, response, codBanco, convenio)
	if err != nil {
		return fmt.Errorf("error to send boleto cancel message: %v", err)
	}

	return nil
}

func (uc *BoletoCancelUseCase) Send(ctx context.Context, response *service.BoletoCancelOutput, codBanco, convenio string) error {

	sendRequest := SendMessageRequest{
		Registrado:      response.Registrado,
		DataHoraRetorno: response.DataHoraRetorno,
		SituacaoTitulo:  response.SituacaoTitulo,
		LinhaDigitavel:  response.LinhaDigitavel,
		CodigoBarras:    response.CodigoBarras,
		NumeroTitulo:    response.NumeroTitulo,
		NossoNumero:     response.NossoNumero,
		URL:             response.URL,
		QRCode:          response.QRCode,
		CodBanco:        &codBanco,
		NumeroConvenio:  &convenio,
	}

	sendRequest.Erros = make([]Erro, 0, len(response.Erros))

	for _, e := range response.Erros {
		sendRequest.Erros = append(sendRequest.Erros, Erro{
			CodigoRejeicao:         e.CodigoRejeicao,
			MensagemRetorno:        e.MensagemRetorno,
			MensagemDetalheRetorno: e.MensagemDetalheRetorno,
			Campo:                  e.Campo,
		})
	}

	messageBody, err := json.Marshal(sendRequest)
	if err != nil {
		return err
	}

	err = uc.sendMessageService.Send(ctx, string(messageBody))
	if err != nil {
		return err
	}

	return nil
}

func stringOrEmpty(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func convertDate(input string) (string, error) {
	const inputLayout = "20060102150405"
	const outputLayout = "02012006"

	t, err := time.Parse(inputLayout, input)
	if err != nil {
		return "", err
	}

	return t.Format(outputLayout), nil
}
