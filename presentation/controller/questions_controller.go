package controller

import (
	"encoding/xml"
	"fmt"
	"goserver/utils/gsrender"
	"net/http"
)

type AnswerOption struct {
	XMLName xml.Name `xml:"answerOptions"`
	OptionId string `xml:"optionId"`
	Text string `xml:"text"`
}

type Question struct {
	XMLName xml.Name `xml:"questions"`
	QuestionId string `xml:"questionId"`
	Order string `xml:"orden"`
	Text string `xml:"text"`
	AnswerOptions []*AnswerOption
}

type QuestionsResponse struct {
	XMLName xml.Name `xml:"ns=http://webservices.idvalidator.veraz.com ns:obtenerPreguntasResponse"`
	Document string `xml:"return>integrantes>documento"`
	Birthdate string `xml:"return>integrantes>fecha_nac"`
	Name string `xml:"return>integrantes>nombre"`
	Sex string `xml:"return>integrantes>sexo"`
	Lot string `xml:"return>lote"`
	Questionary string `xml:"return>questionary"`
	Questions []*Question `xml:"return>questions"`
}

type GetQuestionsResponse struct {
	XMLName xml.Name `xml:"soapenv=http://schemas.xmlsoap.org/soap/envelope/ soapenv:Envelope"`
	QuestionsResponse *QuestionsResponse `xml:"soapenv:Body>ns:obtenerPreguntasResponse"`
}

type QuestionsRequest struct {
	XMLName xml.Name `xml:"obtenerPreguntas"`
	User string `xml:"user"`
	Password string `xml:"password"`
	Sucursal string `xml:"sucursal"`
	DocumentNumber string `xml:"documentNumber"`
	Gender string `xml:"gender"`
	Matrix string `xml:"matrix"`
}

type GetQuestionsRequest struct {
	XMLName xml.Name `xml:"Envelope"`
	Body struct {
		XMLName xml.Name `xml:"Body"`
		QuestionsRequest
	}
}

type Plant struct {
    XMLName xml.Name `xml:"plant"`
    Id      int      `xml:"id,attr"`
    Name    string   `xml:"name"`
    Origin  []string `xml:"origin"`
}

func GetQuestions(w http.ResponseWriter, r *http.Request) {

	var requestBody GetQuestionsRequest

	err := xml.NewDecoder(r.Body).Decode(&requestBody)
	if (err != nil) {
		fmt.Println(err)
		return
	}

	answerOptionA := &AnswerOption{
		OptionId: "1",
		Text: "25/06/1997",
	}
	answerOptionB := &AnswerOption{
		OptionId: "2",
		Text: "18/02/1996",
	}

	questionA := &Question{
		QuestionId: "1",
		Order: "0",
		Text: "Su fecha de nacimiento es:",
		AnswerOptions: []*AnswerOption{answerOptionA, answerOptionB},
	}

	questionsResponse := &QuestionsResponse{
		Document: "39148499",
		Birthdate: "1995-10-18",
		Name: "BIAGINI, MARTIN",
		Sex: "M",
		Lot: "87918595-20391484997",
		Questionary: "12",
		Questions: []*Question{questionA},
	}

	response := &GetQuestionsResponse{}
	response.QuestionsResponse = questionsResponse

	gsrender.WriteXML(w, 200, response)
}