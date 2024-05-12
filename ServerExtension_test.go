package main_test

//
//
//import (
//	"bpmntosd/custom"
//	"encoding/json"
//	"encoding/xml"
//	"fmt"
//	"gitlab.com/teec2/simplified/components/serverextender/client"
//	"gitlab.com/teec2/simplified/components/serverextender/simplifiedFunctions"
//	"gitlab.com/teec2/simplified/components/serverextender/simplifiedTypes"
//	"log"
//	"testing"
//	"time"
//)
//
//var host = "localhost:8088"
//var channel = "/ws"
//var user = "eruci18@gmail.com"
//var password = "Albania98!"
//
//func TestSingleModel(t *testing.T) {
//	var wsClient client.WebSocketClient
//	//wsClient.Logging = true
//
//	go custom.Start(&wsClient, "00000000-0000-0000-0000-000000000001", ".", "eruci18@gmail.com", "Albania98!", "acc.server.simplified.engineering", 443)
//	time.Sleep(4 * time.Second)
//	TstSingleModel(&wsClient, t)
//	wsClient.CloseWs()
//	log.Println("")
//}
//func TstSingleModel(wsClient *client.WebSocketClient, t *testing.T) {
//	const originalName = "TestModel BPMN"
//	const addName = "-to-SD"
//	if wsClient == nil || len(wsClient.Who) < 1 {
//		t.Error("Not correct connection to the server")
//		return
//	}
//	// get the models
//	mdls, errs := simplifiedFunctions.GetModel(wsClient, simplifiedTypes.ModelsByUser, wsClient.Who[0].UserId)
//	if errs != nil || mdls == nil {
//		t.Error("Nope", errs)
//		return
//	}
//
//	//Delete present translations
//	for _, x := range *mdls {
//		if x.Identification == originalName || x.Identification == originalName+addName {
//			fmt.Println(x.Identification + " " + x.Id)
//			deleted, errs := simplifiedFunctions.DeleteModel(wsClient, simplifiedTypes.Model, x)
//			if errs != nil || !deleted {
//				t.Error("Nope", errs)
//				return
//			}
//		}
//	}
//
//	//<editor-fold desc="Create Model">
//
//	// create an empty model without elements
//	emptyModel := simplifiedTypes.MessageModel{
//		Id:             simplifiedFunctions.NewId(),
//		Identification: originalName,
//		RepositoryId:   (*mdls)[0].RepositoryId,
//	}
//
//	demoModels, errs := simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.Model, emptyModel)
//	if errs != nil || len(*demoModels) < 1 {
//		t.Error("Nope", errs)
//		return
//	}
//	demoModel := (*demoModels)[0]
//	mfr1 := simplifiedTypes.MessageModelFolder{
//		XMLName:        xml.Name{},
//		Id:             simplifiedFunctions.NewId(),
//		Identification: "Elements",
//		ModelId:        demoModel.Id,
//	}
//	mfrs, _ := simplifiedFunctions.SaveModelFolder(wsClient, simplifiedTypes.ModelFolder, &mfr1)
//	mfr1 = (*mfrs)[0]
//
//	met1 := simplifiedTypes.MessageModelElement{
//		Id:             simplifiedFunctions.NewId(),
//		Identification: "AR01",
//		ModelId:        demoModel.Id,
//		ModelFolderId:  mfr1.Id,
//		ElementType:    "ElementaryActorRole1",
//		ElementKind:    simplifiedTypes.ElementKindElement,
//		Payload:        nil,
//	}
//	met1Server, errs := simplifiedFunctions.SaveModelElement(wsClient, simplifiedTypes.ModelElement, met1)
//	if errs != nil {
//		t.Error("Nope", errs)
//		return
//	}
//	//set local element to server element
//	met1 = (*met1Server)[0]
//
//	mea1 := simplifiedTypes.MessageModelAttribute{
//		Id:                  simplifiedFunctions.NewId(),
//		Identification:      "name",
//		ModelId:             met1.ModelId,
//		NotationAttributeId: "",
//		//FolderId:            met1.ModelFolderId,
//		ElementId:   met1.Id,
//		ElementType: met1.ElementType,
//		Value:       "ActorRole1Name",
//	}
//	mea1Server, errs := simplifiedFunctions.SaveModelAttribute(wsClient, simplifiedTypes.ModelAttribute, mea1)
//	if errs != nil {
//		t.Error("saving attribute failed", errs)
//		return
//	}
//	//set local attribute to server element
//	mea1 = (*mea1Server)[0]
//
//	data, err := json.Marshal(demoModel)
//	if err != nil {
//		t.Error("Marshal error")
//		return
//	}
//	msgMdl := &simplifiedTypes.Message{
//		MessageId:       simplifiedFunctions.NextMessageId(),
//		MessageType:     simplifiedTypes.MessageTypeGET,
//		ContentType:     simplifiedTypes.Model,
//		ProtocolVersion: simplifiedTypes.Version01,
//		Payload:         data,
//	}
//	//</editor-fold>
//
//	result, errs := custom.TransformOriginalToTarget(wsClient, msgMdl)
//	if errs != nil {
//		t.Error("Nope", errs)
//		return
//	}
//
//	//<editor-fold desc="Test transformation Model">
//
//	mdlRet := &simplifiedTypes.MessageModel{}
//	err = json.Unmarshal(result[0].Payload, mdlRet)
//	if err != nil {
//		t.Error("Nope", errs)
//		return
//	}
//
//	mdls3, errs := simplifiedFunctions.GetModel(wsClient, simplifiedTypes.Model, mdlRet.Id)
//	if errs != nil || mdls3 == nil {
//		t.Error("Nope", errs)
//		return
//	}
//	if (*mdls3)[0].Identification != originalName+addName {
//		t.Error("Nope")
//		return
//	}
//	mets, errs := simplifiedFunctions.GetModelElement(wsClient, simplifiedTypes.ModelElementsByModel, (*mdls3)[0].Id)
//	if errs != nil {
//		t.Error("Nope", errs)
//		return
//	}
//	if len(*mets) > 0 {
//		for _, met := range *mets {
//			fmt.Println(met.Identification + " " + met.ElementType + "" + met.ModelId)
//		}
//		return
//	}
//	//</editor-fold>
//
//}
