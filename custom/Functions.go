package custom

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/go-errors/errors"
	"gitlab.com/teec2/simplified/components/serverextender/client"
	"gitlab.com/teec2/simplified/components/serverextender/simplifiedFunctions"
	"gitlab.com/teec2/simplified/components/serverextender/simplifiedTypes"
	"log"
	"strings"
)

func GetOriginalFolder(foldername string, sourcefolders *[]simplifiedTypes.MessageModelFolder) (result *simplifiedTypes.MessageModelFolder) {
	for _, folder := range *sourcefolders {
		if folder.Identification == foldername {
			result = &folder
			break
		}
	}
	return
}
func GetOriginalDiagram(folderid string, elements *[]simplifiedTypes.MessageModelElement) (result *simplifiedTypes.MessageModelElement) {
	for _, element := range *elements {
		if element.ElementKind == simplifiedTypes.ElementKindDiagram {
			result = &element
			break
		}
	}
	return
}

func CreateStockConn(wsClient *client.WebSocketClient, Identification, modelid string, SFDdiagram string, sourcelementid, targetElemntid, path string) (*simplifiedTypes.MessageModelConnection, *simplifiedTypes.MessageModelVisualConnection) {
	conn1 := simplifiedTypes.MessageModelConnection{
		Id:                   simplifiedFunctions.NewId(),
		Identification:       Identification,
		ModelId:              modelid,
		ConnectionType:       Flow,
		ModelElementIdSource: sourcelementid,
		ModelElementIdTarget: targetElemntid,
		ForeignId:            "",
		ForeignName:          "",
		Payload:              nil,
	}
	cons, errs := simplifiedFunctions.SaveModelConnection(wsClient, simplifiedTypes.ModelConnection, conn1)
	if errs != nil {
		fmt.Println("Error creating Flow connection in SD diagram", errs)
		return nil, nil
	}

	results := &simplifiedTypes.MessageModelConnection{}
	for _, c := range *cons {
		fmt.Println("Conn", c)
		if c.Identification == Identification {
			results = &c
			break
		}
	}

	visualcon := CreateSFDStockVisualConnection(wsClient, results, SFDdiagram, modelid, Identification, path)

	return results, visualcon
}

func CreateSFDStockVisualConnection(wsClient *client.WebSocketClient, SDconn *simplifiedTypes.MessageModelConnection, ModelElementDiagramId, modelid, identification string, path string) *simplifiedTypes.MessageModelVisualConnection {
	con1 := simplifiedTypes.MessageModelVisualConnection{
		Id:                         simplifiedFunctions.NewId(),
		Identification:             &identification,
		ModelId:                    modelid,
		ConnectionId:               SDconn.Id,
		ModelElementDiagramId:      ModelElementDiagramId,
		ModelVisualElementIdSource: SDconn.ModelElementIdSource,
		ModelVisualElementIdTarget: SDconn.ModelElementIdTarget,
		ConnectionType:             SDconn.ConnectionType,
		Visible:                    true,
		Path:                       nil,
		Payload:                    nil,
	}

	cons, errs := simplifiedFunctions.SaveModelVisualConnection(wsClient, simplifiedTypes.ModelVisualConnection, con1)
	if errs != nil {
		fmt.Println("Error creating Source element", errs)
		return nil
	}

	results := &simplifiedTypes.MessageModelVisualConnection{}
	for _, c := range *cons {
		if c.ConnectionId == SDconn.Id {
			results = &c
		}
	}
	fmt.Printf("0000000000000000>>>>>>>>>>>>>>>>>%s\n", *results.Path)
	return results
}

func CreateSFDElement(wsClient *client.WebSocketClient, modelid, folderid, elementname, elementtype, diagramid, sfddiagram string, positionX, positionY int) (sfdElement *simplifiedTypes.MessageModelElement, visualsfdElement *simplifiedTypes.MessageModelVisualElement) {
	met1 := simplifiedTypes.MessageModelElement{
		Id:             simplifiedFunctions.NewId(),
		Identification: elementname,
		ModelId:        modelid,
		ModelFolderId:  folderid,
		ElementType:    elementtype,
		ElementKind:    simplifiedTypes.ElementKindElement,
		Payload:        nil,
	}
	elementss, errs := simplifiedFunctions.SaveModelElement(wsClient, simplifiedTypes.ModelElement, met1)
	if errs != nil {
		fmt.Println("Error creating Source element", errs)
		return nil, nil
	}
	results := &simplifiedTypes.MessageModelElement{}
	for _, ele := range *elementss {
		if ele.Identification == elementname {
			results = &ele
		}
	}
	visualCons, errs := simplifiedFunctions.GetModelVisualElement(wsClient, simplifiedTypes.ModelVisualConnectionsByDiagram, diagramid)
	if errs != nil {
		fmt.Println("Error getting Model Visual Connections", errs)
		//return
	}
	bpmvisual := &simplifiedTypes.MessageModelVisualElement{}
	if len(*visualCons) > 0 {
		for _, vel := range *visualCons {
			if &vel.ElementType == &elementtype {
				bpmvisual = &vel
			}
		}
	}
	vc := CreateSFDVisualElement(wsClient, results, bpmvisual, sfddiagram, modelid, elementtype, positionX, positionY)

	return results, vc
}
func CreateSFDVisualElement(wsClient *client.WebSocketClient, sfdCreatedElemnt *simplifiedTypes.MessageModelElement, bpmvisual *simplifiedTypes.MessageModelVisualElement, ModelElementDiagramId, modelid, elementtype string, positionX, positionY int) (visualelement *simplifiedTypes.MessageModelVisualElement) {
	met1 := simplifiedTypes.MessageModelVisualElement{
		Id:                    simplifiedFunctions.NewId(),
		ModelId:               sfdCreatedElemnt.ModelId,
		ElementType:           sfdCreatedElemnt.ElementType,
		PositionX:             positionX,
		PositionY:             positionY,
		Width:                 50,
		Height:                55,
		ScaleX:                1,
		ScaleY:                1,
		Rotation:              0,
		ModelElementId:        sfdCreatedElemnt.Id,
		ModelVisualElementId:  bpmvisual.ModelVisualElementId,
		ModelElementDiagramId: ModelElementDiagramId,
		ZOrder:                bpmvisual.ZOrder,
		ShapeVariables:        bpmvisual.ShapeVariables,
		Payload:               nil,
	}
	fmt.Println("----------------------------befroecreating visual element")
	vels, errs := simplifiedFunctions.SaveModelVisualElement(wsClient, simplifiedTypes.ModelVisualElement, met1)
	if errs != nil {
		fmt.Println("Error creating SFD Diagram", errs)
		return nil
	}

	results := &simplifiedTypes.MessageModelVisualElement{}
	for _, el := range *vels {
		if el.ModelElementId == met1.ModelElementId {
			results = &el
		}
	}
	fmt.Println("visual element ---------id", results)
	return results
}

func AdjustVisualElementPosition(positionX, positionY int, adjusterX, adjusterY int) (newPositionX, newPositionY int) {
	return positionX + adjusterX, positionY + adjusterY
}
func CreateSFDVisualConnection(wsClient *client.WebSocketClient, SDconn *simplifiedTypes.MessageModelConnection, ModelElementDiagramId, modelid string, visualHeadEl, visualTailEl *simplifiedTypes.MessageModelVisualElement, vcpath string) *simplifiedTypes.MessageModelVisualConnection {
	con1 := simplifiedTypes.MessageModelVisualConnection{
		Id:                         simplifiedTypes.NewId(),
		Identification:             &SDconn.Identification,
		ModelId:                    modelid,
		ConnectionId:               SDconn.Id,
		ModelVisualElementIdSource: visualHeadEl.Id,
		ModelVisualElementIdTarget: visualTailEl.Id,
		ConnectionType:             Flow,
		Path:                       &vcpath,
		ModelElementDiagramId:      ModelElementDiagramId,
		Visible:                    true,
		Payload:                    nil,
	}

	con1.Id = con1.Id
	cons, errs := simplifiedFunctions.SaveModelVisualConnection(wsClient, simplifiedTypes.ModelVisualConnection, con1)
	if errs != nil {
		fmt.Println("Error Save Model Visual Connection", errs, cons)
		return nil
	}

	results := &simplifiedTypes.MessageModelVisualConnection{}
	for _, c := range *cons {
		if c.ConnectionId == SDconn.Id {
			results = &c
		}
	}
	return results
	//fmt.Println("create conn step 7 ===============================>", *con1.Path)
	//results := &simplifiedTypes.MessageModelVisualConnection{}
	//for _, c := range *cons {
	//	results = &c
	//}
	//return results
}

func GetModelFolder(wsClient *client.WebSocketClient, modelId, folderName string) {

}
func CreateSDFolder(wsClient *client.WebSocketClient, model simplifiedTypes.MessageModel, folderName string) *simplifiedTypes.MessageModelFolder {
	mfr1 := simplifiedTypes.MessageModelFolder{
		XMLName:         xml.Name{},
		Id:              simplifiedFunctions.NewId(),
		Identification:  folderName,
		ModelId:         model.Id,
		VersionObjectId: model.RepositoryId,
	}
	_, errs := simplifiedFunctions.SaveModelFolder(wsClient, simplifiedTypes.ModelFolder, &mfr1)
	if errs != nil {
		fmt.Println("Saving transformed folder failed", errs)
	}
	return &mfr1
}

func CreateTargetDiagram(wsClient *client.WebSocketClient, model simplifiedTypes.MessageModel, folderid, elementname string) *simplifiedTypes.MessageModelElement {
	met1 := simplifiedTypes.MessageModelElement{
		Id:              simplifiedFunctions.NewId(),
		Identification:  elementname,
		ModelId:         model.Id,
		ModelFolderId:   folderid,
		ElementType:     SFD,
		ElementKind:     simplifiedTypes.ElementKindDiagram,
		Payload:         nil,
		VersionObjectId: model.RepositoryId,
	}
	vds, errs := simplifiedFunctions.SaveModelElement(wsClient, simplifiedTypes.ModelElement, met1)
	if errs != nil {
		fmt.Println("Error creating SFD Diagram", errs)
		return nil
	}
	results := &simplifiedTypes.MessageModelElement{}
	for _, dv := range *vds {
		if dv.Identification == elementname {
			results = &dv
		}
	}
	return results
}

func CreateVisualTargetEl(wsClient *client.WebSocketClient, modelElementId, modelId, diagramId, nvts, eType string,
	BPMNvisual simplifiedTypes.MessageModelVisualElement) *simplifiedTypes.MessageModelVisualElement {
	met1 := simplifiedTypes.MessageModelVisualElement{
		Id:                      simplifiedFunctions.NewId(),
		ModelElementId:          modelElementId,
		ModelVisualElementId:    "",
		ModelId:                 modelId,
		ModelElementDiagramId:   diagramId,
		NotationVisualElementId: nvts,
		ElementType:             eType,
		PositionX:               BPMNvisual.PositionX,
		PositionY:               BPMNvisual.PositionY,
		Width:                   50,
		Height:                  50,
		Payload:                 nil,
		ScaleX:                  1,
		ScaleY:                  1,
		Rotation:                0,
		ZOrder:                  0,
	}
	vels, errs := simplifiedFunctions.SaveModelVisualElement(wsClient, simplifiedTypes.ModelVisualElement, met1)
	if errs != nil {
		fmt.Println("Error creating Visual Target element", eType, errs)
		return nil
	}
	results := &simplifiedTypes.MessageModelVisualElement{}
	for _, el := range *vels {
		if el.ModelElementId == met1.ModelElementId {
			results = &el
		}
	}
	return results
}
func CreateVisualActiveStock(wsClient *client.WebSocketClient, modelElementId, modelId, diagramId, nvts, eType string, BPMNvisual simplifiedTypes.MessageModelVisualElement) *simplifiedTypes.MessageModelVisualElement {
	met1 := simplifiedTypes.MessageModelVisualElement{
		Id:                      simplifiedFunctions.NewId(),
		ModelElementId:          modelElementId,
		ModelVisualElementId:    "",
		ModelId:                 modelId,
		ModelElementDiagramId:   diagramId,
		NotationVisualElementId: nvts,
		ElementType:             eType,
		PositionX:               BPMNvisual.PositionX - 50,
		PositionY:               BPMNvisual.PositionY,
		Width:                   50,
		Height:                  50,
		Payload:                 nil,
		ScaleX:                  1,
		ScaleY:                  1,
		Rotation:                0,
		ZOrder:                  0,
	}
	vels, errs := simplifiedFunctions.SaveModelVisualElement(wsClient, simplifiedTypes.ModelVisualElement, met1)
	if errs != nil {
		fmt.Println("Error creating Visual Target element", eType, errs)
		return nil
	}
	results := &simplifiedTypes.MessageModelVisualElement{}
	for _, el := range *vels {
		if el.ModelElementId == met1.ModelElementId {
			results = &el
		}
	}
	return results
}
func CreateVisualFinishedStock(wsClient *client.WebSocketClient, modelElementId, modelId, diagramId, nvts, eType string, BPMNvisual simplifiedTypes.MessageModelVisualElement) *simplifiedTypes.MessageModelVisualElement {
	met1 := simplifiedTypes.MessageModelVisualElement{
		Id:                      simplifiedFunctions.NewId(),
		ModelElementId:          modelElementId,
		ModelVisualElementId:    "",
		ModelId:                 modelId,
		ModelElementDiagramId:   diagramId,
		NotationVisualElementId: nvts,
		ElementType:             eType,
		PositionX:               BPMNvisual.PositionX + 60,
		PositionY:               BPMNvisual.PositionY,
		Width:                   50,
		Height:                  50,
		Payload:                 nil,
		ScaleX:                  1,
		ScaleY:                  1,
		Rotation:                0,
		ZOrder:                  0,
	}
	vels, errs := simplifiedFunctions.SaveModelVisualElement(wsClient, simplifiedTypes.ModelVisualElement, met1)
	if errs != nil {
		fmt.Println("Error creating Visual Target element", eType, errs)
		return nil
	}
	results := &simplifiedTypes.MessageModelVisualElement{}
	for _, el := range *vels {
		if el.ModelElementId == met1.ModelElementId {
			results = &el
		}
	}
	return results
}

func CreateTargetVisualElement(wsClient *client.WebSocketClient, bpmvisual *simplifiedTypes.MessageModelVisualElement, modelId, modelElementDiagramId, modelElementId, elementType string) *simplifiedTypes.MessageModelVisualElement {
	met1 := simplifiedTypes.MessageModelVisualElement{
		Id:                      simplifiedFunctions.NewId(),
		ModelElementId:          modelElementId,
		ModelId:                 modelId,
		ModelVisualElementId:    bpmvisual.ModelVisualElementId,
		ModelElementDiagramId:   modelElementDiagramId,
		NotationVisualElementId: bpmvisual.NotationVisualElementId,
		ElementType:             elementType,
		PositionX:               bpmvisual.PositionX,
		PositionY:               bpmvisual.PositionY,
		Width:                   bpmvisual.Width,
		Height:                  bpmvisual.Height,
		ScaleX:                  bpmvisual.ScaleX,
		ScaleY:                  bpmvisual.ScaleY,
		Rotation:                bpmvisual.Rotation,
		ZOrder:                  bpmvisual.ZOrder,
		ShapeVariables:          bpmvisual.ShapeVariables,
		Payload:                 nil,
	}
	vels, errs := simplifiedFunctions.SaveModelVisualElement(wsClient, simplifiedTypes.ModelVisualElement, met1)
	if errs != nil {
		fmt.Println("Error creating Source element", errs)
		return nil
	}
	results := &simplifiedTypes.MessageModelVisualElement{}
	for _, el := range *vels {
		if el.ModelElementId == met1.ModelElementId {
			results = &el
		}
	}
	return results
}

func CreateVisualFlowTargetConn(wsClient *client.WebSocketClient, BPMNconn *simplifiedTypes.MessageModelVisualConnection, SDconn *simplifiedTypes.MessageModelConnection, ModelElementDiagramId, modelid string) *simplifiedTypes.MessageModelVisualConnection {
	con1 := simplifiedTypes.MessageModelVisualConnection{
		Id:                         simplifiedFunctions.NewId(),
		Identification:             BPMNconn.Identification,
		ModelId:                    modelid,
		ConnectionId:               SDconn.Id,
		ModelElementDiagramId:      ModelElementDiagramId,
		ModelVisualElementIdSource: SDconn.ModelElementIdSource,
		ModelVisualElementIdTarget: SDconn.ModelElementIdTarget,
		ConnectionType:             SDconn.ConnectionType,
		Visible:                    BPMNconn.Visible,
		//Visible: !BPMNconn.Visible,
		Path:    BPMNconn.Path,
		Payload: nil,
	}

	//fmt.Println("create conn step 6 ===============================>")
	cons, errs := simplifiedFunctions.SaveModelVisualConnection(wsClient, simplifiedTypes.ModelVisualConnection, con1)
	if errs != nil {
		fmt.Println("Error creating Source element", errs)
		return nil
	}

	//fmt.Println("create conn step 7 ===============================>")
	results := &simplifiedTypes.MessageModelVisualConnection{}
	for _, c := range *cons {
		if c.ConnectionId == SDconn.Id {
			results = &c
		}
	}
	return results
}

//func CustomCreateVisual(wsClient *client.WebSocketClient, targetModelId, targetDiagramId, sourceTargetElId, sourceTargetElementType string, startOriginalEl *simplifiedTypes.MessageModelElement, visualOriginalEls *[]simplifiedTypes.MessageModelVisualElement) {
//	visualStartOriginalEl := &simplifiedTypes.MessageModelVisualElement{}
//	if len(*visualOriginalEls) > 0 {
//		for _, ve := range *visualOriginalEls {
//			if startOriginalEl.Id == ve.ModelElementId {
//				visualStartOriginalEl = &ve
//			}
//		}
//	}
//
//	CreateTargetVisualElement(wsClient, visualStartOriginalEl, targetModelId, targetDiagramId, sourceTargetElId, sourceTargetElementType)
//}

func MyFunction(wsClient *client.WebSocketClient, msg *simplifiedTypes.Message) (result []simplifiedTypes.Message, error *errors.Error) {
	met := simplifiedTypes.MessageModelElement{
		Identification:    "My Element",
		ModelId:           "",
		ModelFolderId:     "",
		NotationElementId: "",
		ElementType:       "MyNotationElement",
		AttributeValue:    nil,
		ElementKind:       simplifiedTypes.ElementKindElement,
	}
	// Create a answer payload
	buffer, err := json.Marshal(met)
	if err != nil {
		log.Println(err)
		return
	}
	// Create a request message containing the payload
	message := simplifiedTypes.Message{Payload: buffer, MessageType: simplifiedTypes.MessageTypeGOT, ContentType: simplifiedTypes.ModelElement}

	// Send the message to the server
	return []simplifiedTypes.Message{message}, nil
}

func TransformOriginalToTarget2(wsClient *client.WebSocketClient, msg *simplifiedTypes.Message) (result []simplifiedTypes.Message, error *errors.Error) {
	originalModel := &simplifiedTypes.MessageModel{}
	err := json.Unmarshal(msg.Payload, &originalModel)
	if err != nil {
		return nil, errors.New(err)
	}
	//Set translated DatabaseModel name. Uses original DEMO model identification + "-to-VISI-TransCon".
	//The ModelConnections are added to the original model from which they are translated.
	//This will enable to keep track of which element, and which connection has been translated into what element or what connection in the new model.
	//translationConnection := originalModel.Identification + "-to-VISI-TransCon"
	//Create New VISI model, set Identification as DEMO model identification + "to-VISI".
	//This model will be used to store the to be translated model. It is given a random "Id", uses the original model identification and appends "-to-VISI" to it.
	//It also uses the same "repositoryId" as the original model.
	found := false
	visiModel := simplifiedTypes.MessageModel{}
	mdls, errs := simplifiedFunctions.GetModel(wsClient, simplifiedTypes.ModelsByUser, wsClient.Who[0].UserId)
	if errs != nil || mdls == nil {
		return nil, errs
	}
	for _, mdl := range *mdls {
		if mdl.Identification == originalModel.Identification+"-to-VISI" {
			fmt.Println(mdl.Identification + " " + mdl.Id)
			visiModel = mdl
			found = true
			break
		}
	}
	if !found {
		visiModel = simplifiedTypes.MessageModel{
			Id:             simplifiedFunctions.NewId(),
			Identification: originalModel.Identification + "-to-VISI",
			RepositoryId:   originalModel.RepositoryId,
		}
		mdls, errs := simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.Model, &visiModel)
		if errs != nil || mdls == nil {
			return nil, errs
		}
		visiModel = (*mdls)[0]
	}
	//
	//	//Delete present translations
	//	//TODO: Wordt hier nog niet verwijderd
	//	for i, mdl := range *mdls {
	//		if mdl.Identification == visiModel.Identification {
	//			visiModel = (*mdls)[i]
	//			found = true
	//			break
	//		}
	//	}
	//	if !found {
	//		//Save the new model. Uses the uc(user credentials) to check if this user is allowed to do this procedure.
	//		//In order to be able to reference to this model and store the created elements, connections, and attributes, the model will have to exist, otherwise it would return nil, and crash the translation.
	//		//TODO:Check genoeg why?
	//		var mdlSave *[]simplifiedTypes.MessageModel
	//		mdlSave, errs = simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.Model, visiModel)
	//		//_, _, errs = simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.Model, visiModel)
	//		//In case the saving of the model does not succeed. Create panic with error message.
	//		//TODO:Michiel Should not be able to return updated right? should only be able to receive saved. Create check for this.
	//		if len(*mdlSave) < 1 || errs != nil {
	//			return nil, errs
	//		}
	//	}
	//
	//	//Set folderId to the util.NilId (00000000-0000-0000-0000-000000000000) to initiate the folderId as an uuid.
	//	folderId := ""
	//	//Retrieve modelFolders used in original model.
	//	sourceFolders, errs := simplifiedFunctions.GetModelFolder(wsClient, simplifiedTypes.ModelFoldersByModel, originalModel.Id)
	//	//Check for errors while retrieving the source modelFolders. If error occurs halt the program and show error.
	//	if errs != nil {
	//		fmt.Errorf("Problem retrieving modelFolders from modelIdent:" + visiModel.Identification)
	//		return nil, errs
	//	}
	//	//Check if there is a folder called "Elements" in the sourcefolder. Use this ID if this is the case, otherwise generate new ID for folder.
	//	if len(*sourceFolders) > 0 {
	//		for _, folder := range *sourceFolders {
	//			if folder.Identification == "Elements" {
	//				folderId = folder.Id
	//			}
	//		}
	//		if folderId == "" {
	//			//If the folder "Elements" does not exist in the source folder create a new id for the folder
	//			//TODO: Check if this should happen, or should throw exception, or just keep util.nilID
	//			folderId = simplifiedFunctions.NewId()
	//		}
	//	}

	//Start of Transform Request.
	transformRequest := simplifiedTypes.TransformModelRequest{
		//		FromModel:           originalModel.Id,
		//		ToModel:             visiModel.Id,
		//		FolderDestiny:       folderId,
		//		TransformConnection: translationConnection,
		//		// 			<!--!!E.TransformOriginalToTarget.Request!!-->
		//		// 			<!--!!S.TransformOriginalToTarget.ElementToElement!!-->
		//		//
		//		//									ELEMENT TO ELEMENT
		//		//
		//		FromElementTypeToElementType: []simplifiedTypes.FromElementTypeToElementType{
		//			// 			<!--!!E.TransformOriginalToTarget.ElementToElement!!-->
		//			// 			<!--!!S.TransformOriginalToTarget.RoleTypes!!-->
		//			//
		//			//			Elementary Actor Roles to RoleTypes
		//			//				For each Elementary Actor Role in demo creates a RoleType in VISI
		//			{
		//				FromElementType: demo3.DemoEar,
		//				ToElementType:   visi16.VisiRoleType,
		//				AttributeFromTypeToTypeForElement: []simplifiedTypes.AttributeFromTypeToTypeForElement{
		//					//	give the RoleType element the attribute "description". Use the attribute "name" from the demo EAR as value.
		//					{
		//						FromAttributeIdentification: "attribute.name",
		//						ToAttributeIdentification:   "attribute.description",
		//					},
		//					//  Use the demo EAR element identification as element identification for RoleType.
		//					//  remove leading EAR identification (AR) and replace it with RT (only for readability)
		//					{
		//						FromAttributeIdentification: "element.identification",
		//						ToAttributeIdentification:   "element.identification",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddRtPrefix,
		//					},
		//					//	give the RoleType element the attribute "id". Use the EAR identification as value. (Adding the attribute id is for readability, xml would support random identifiers)
		//					{
		//						FromAttributeIdentification: "element.identification",
		//						ToAttributeIdentification:   "attribute.id",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddRtPrefix,
		//					},
		//				},
		//			},
		//			// 			<!--!!E.TransformOriginalToTarget.RoleTypes!!-->
		//			// 			<!--!!S.TransformOriginalToTarget.TransactionTypes!!-->
		//			//
		//			//			TransactionKinds to TransactionTypes
		//			//				For each TransactionKind in DEMO create a TransactionType in VISI
		//			{
		//				FromElementType: demo3.DemoTk,
		//				ToElementType:   visi16.VisiTransactionType,
		//				AttributeFromTypeToTypeForElement: []simplifiedTypes.AttributeFromTypeToTypeForElement{
		//					//  Use the demo TransactionKind element identification as element identification for the VISI TransactionType.
		//					//  remove leading TransactionKind identification (TK) and replace it with TT (only for readability)
		//					{
		//						FromAttributeIdentification: "element.identification",
		//						ToAttributeIdentification:   "element.identification",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddTtPrefix,
		//					},
		//					//	give the TransactionType element the attribute "description". Use the attribute "name" from the DEMO TransactionKind as value.
		//					{
		//						FromAttributeIdentification: "attribute.name",
		//						ToAttributeIdentification:   "attribute.description",
		//					},
		//					//	give the TransactionType element the attribute "id". Use the TransactionKind identification as value. (Adding the attribute id is for readability, xml would support random identifiers)
		//					{
		//						FromAttributeIdentification: "element.identification",
		//						ToAttributeIdentification:   "attribute.id",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddTtPrefix,
		//					},
		//				},
		//			},
		//			// 			<!--!!E.TransformOriginalToTarget.TransactionTypes!!-->
		//			// 			<!--!!S.TransformOriginalToTarget.Groups!!-->
		//			//
		//			//			TransactionKinds to GroupTypes
		//			//				For each TransactionKind in DEMO create a GroupType in VISI
		//			{
		//				FromElementType: demo3.DemoTk,
		//				ToElementType:   visi16.VisiGroupType,
		//				AttributeFromTypeToTypeForElement: []simplifiedTypes.AttributeFromTypeToTypeForElement{
		//					//  Use the demo TransactionKind element identification as element identification for the VISI GroupType.
		//					//  remove leading TransactionKind identification (TK) and replace it with GT (only for readability)
		//					{
		//						FromAttributeIdentification: "element.identification",
		//						ToAttributeIdentification:   "element.identification",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddGtPrefix,
		//					},
		//					//	give the GroupType element the attribute "description". Use the attribute "name" from the DEMO TransactionKind as value.
		//					//  Add "GroupType-" before name, to distinct from TransactionTypes.
		//					//TODO Michiel: hier is een hoofdletter gebruikt bij "Name" hierboven is het "name". checken
		//					{
		//						FromAttributeIdentification: "attribute.Name",
		//						ToAttributeIdentification:   "attribute.description",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddGrTpPrefix,
		//					},
		//					//	give the GroupType element the attribute "id". Use the TransactionKind identification as value. (Adding the attribute id is for readability, xml would support random identifiers)
		//					{
		//						FromAttributeIdentification: "element.identification",
		//						ToAttributeIdentification:   "attribute.id",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddGtPrefix,
		//					},
		//				},
		//			},
		//			// 			<!--!!E.TransformOriginalToTarget.Groups!!-->
		//			// 			<!--!!S.TransformOriginalToTarget.Mitts!!-->
		//			//
		//			//					Transaction Process Step Kind(TPSK) to MessageInTransactionType(Mitt)
		//			//						For each TPSK in DEMO create a Mitt in VISI
		//			{
		//				FromElementType: demo3.DemoTpsk,
		//				ToElementType:   visi16.VisiMessageInTransactionType,
		//				AttributeFromTypeToTypeForElement: []simplifiedTypes.AttributeFromTypeToTypeForElement{
		//					//  Use the demo TPSK attribute "StepKind" as element identification for the VISI Mitt.
		//					//  Add "Mitt-" in front of this identification to indicate this is a Mitt (only for readability)
		//					{
		//						FromAttributeIdentification: "attribute.StepKind",
		//						ToAttributeIdentification:   "element.identification",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddMittPrefix,
		//					},
		//					//  Use the demo TPSK attribute "StepKind" add it to the attribute "id" for the VISI Mitt.
		//					//  Add "Mitt-" in front of it for readability
		//					{
		//						FromAttributeIdentification: "attribute.StepKind",
		//						ToAttributeIdentification:   "attribute.id",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddMittPrefix,
		//					},
		//					// Checks if the mitt should have the attribute "firstMessage" value true or false. By checking if the tpsk is either the original starting point of all transactions, or it's a starting point of a sub-transaction.
		//					// TODO:Michiel: the attribute name does not constitute to the value of the target attribute. should it be renamed?
		//					{
		//						FromAttributeIdentification: "attribute.name",
		//						ToAttributeIdentification:   "attribute.firstMessage",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fIsFirstMessage,
		//					},
		//					// Checks if the mitt should have the attribute "initiatorToExecutor" value true or false. By checking if the step kind of the TPSK belongs to the initiator or the executor.
		//					{
		//						FromAttributeIdentification: "attribute.name",
		//						ToAttributeIdentification:   "attribute.initiatorToExecutor",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fIsInitiatorToExecutor,
		//					},
		//				},
		//			},
		//			// 			<!--!!E.TransformOriginalToTarget.Mitts!!-->
		//			// 			<!--!!S.TransformOriginalToTarget.TransactionPhases!!-->
		//			//
		//			//					Transaction Process Step Kind (TPSK) to TransactionPhaseType
		//			//						For each (DISTINCT) TPSK type in DEMO create a TransactionPhaseType in VISI
		//			//						A demo model can contain duplicate tpsk stepkinds. For instance: multiple "requests". In VISI this would be the same phase and this phase would be reused.
		//			//						Therefore, the "CreateOnePerSameIdentification" is set to true. This will prevent the generation of multiple TransactionPhases with the same meaning.
		//			{
		//				FromElementType:                demo3.DemoTpsk,
		//				ToElementType:                  visi16.VisiTransactionPhaseType,
		//				CreateOnePerSameIdentification: true,
		//				AttributeFromTypeToTypeForElement: []simplifiedTypes.AttributeFromTypeToTypeForElement{
		//					//	Use the value of the attribute "Abbreviation" from the TPSK for the attribute "code" in the TransactionPhaseType.
		//					{
		//						FromAttributeIdentification: "attribute.Abbreviation",
		//						ToAttributeIdentification:   "attribute.code",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddAbbreviationToCode,
		//					},
		//					// Use the Tpsk identification as value for the attribute "description" for the TransactionPhaseType.
		//					{
		//						FromAttributeIdentification: "element.identification",
		//						ToAttributeIdentification:   "attribute.description",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddStepKindToDescription,
		//					},
		//					// Use the stepkind naming as the transactionPhaseType identification. Add "TPhT-" in front of the stepkind value for readability
		//					{
		//						FromAttributeIdentification: "attribute.StepKind",
		//						ToAttributeIdentification:   "element.identification",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddStepKindToIdent,
		//					},
		//					// Use the StepKind attribute as the "id" attribute for the created transactionPhaseType. Add "TPhT-" in front of the stepkind value for readability. Adding the attribute Id is mainly to improve readability of a model viewing the XML structure.
		//					{
		//						FromAttributeIdentification: "attribute.StepKind",
		//						ToAttributeIdentification:   "attribute.id",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddStepKindToIdent,
		//					},
		//				},
		//			},
		//			// 			<!--!!E.TransformOriginalToTarget.TransactionPhases!!-->
		//		},
		//
		//		// 			<!--!!S.TransformOriginalToTarget.ElementToConnection!!-->
		//		//
		//		//												ELEMENT TO CONNECTION
		//		//
		//		//  List of Element to Connection Transformations.
		//		FromElementTypeToConnectionType: []simplifiedTypes.FromElementTypeToConnectionType{
		//			// 			<!--!!E.TransformOriginalToTarget.ElementToConnection!!-->
		//			// 			<!--!!S.TransformOriginalToTarget.PhasesInMitt!!-->
		//			//
		//			//			From Transaction Process Step Kind (TPSK) to TransactionPhaseType in Mitt
		//			//				Give a Mitt a TransactionPhaseType.
		//			//Create the connection VisiPhaseInMittConnection, which connects TransactionPhaseType to Mitt.
		//			{
		//				FromElementType:  demo3.DemoTpsk,
		//				ToConnectionType: visi16.VisiPhaseInMittConnection,
		//				//FromFunction:     nil,
		//				AttributeFromTypeToTypeForElement: []simplifiedTypes.AttributeFromTypeToTypeForElement{
		//					//use Mitt id as sourceId for the PhaseInMittConnection. Where the Tpsk Element Id TODO:Uitleg uitbreiden, komt even niet naar boven hoe ik het moet uitleggen.
		//					//Get the newly created TransactionPhase id and use it as sourceId for the to be created connection (VisiPhaseInMittConnection).
		//					//Get the object belonging to the sourceId of the connection.
		//					//Get the connection.sourceId  with sourceId of the TPSK
		//					{
		//						FromAttributeIdentification: "connection.sourceId",
		//						ToAttributeIdentification:   "connection.sourceId",
		//						FromType:                    demo3.DemoTpsk,
		//						ToType:                      visi16.VisiMessageInTransactionType, //Mitt
		//					},
		//					//TODO:Uitleggen, inverse van bovenstaande uitleggen.
		//					{
		//						FromAttributeIdentification: "connection.targetId",
		//						ToAttributeIdentification:   "connection.targetId",
		//						FromType:                    demo3.DemoTpsk,
		//						ToType:                      visi16.VisiTransactionPhaseType, //TransactionPhaseType
		//					},
		//					//add connection identification to the newly made connection fill with "Phase in" (mainly for readability)
		//					{
		//						//TODO: Element.ElementType klopt niet... hier wordt eigenlijk puur context toegevoegd.
		//						FromAttributeIdentification: "element.elementType",
		//						ToAttributeIdentification:   "connection.identification",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddPhaseInMitt,
		//					},
		//				},
		//			},
		//			// 			<!--!!E.TransformOriginalToTarget.PhasesInMitt!!-->
		//		},
		//
		//		// 			<!--!!S.TransformOriginalToTarget.ConnectionToElement!!-->
		//		//
		//		//												CONNECTION TO ELEMENT
		//		//
		//		// List of Connection to Element transformations
		//		FromConnectionTypeToElementType: []simplifiedTypes.FromConnectionTypeToElementType{
		//			// 			<!--!!E.TransformOriginalToTarget.ConnectionToElement!!-->
		//			// 			<!--!!S.TransformOriginalToTarget.ConditionType!!-->
		//
		//			//			From TPSK wait link to MessageInTransactionTypeCondition
		//			{
		//				FromConnectionType: demo3.DemoTpskWaitLink,
		//				ToElementType:      visi16.VisiMessageInTransactionTypeCondition,
		//				AttributeFromTypeToTypeForConnection: []simplifiedTypes.AttributeFromTypeToTypeForConnection{
		//					//Use Tpsk wait link connection identification for MessageInTransactionTypeCondition identification
		//					{
		//						FromAttributeIdentification: "connection.identification",
		//						ToAttributeIdentification:   "element.identification",
		//					},
		//					// Waitlink identification format: TK02/stepkind/TK01/da  . Where TK01/da is the TPSK waiting on the completion of TPSK TK02/ac.
		//					// The sourceId points to the
		//					// The targetId points to the
		//					{
		//						FromAttributeIdentification: "connection.sourceId",
		//						ToAttributeIdentification:   "attribute.id",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddSendAfter,
		//					},
		//				},
		//			},
		//			// 			<!--!!E.TransformOriginalToTarget.ConditionType!!-->
		//		},
		//
		//		// 			<!--!!S.TransformOriginalToTarget.ConnectionToConnection!!-->
		//		//
		//		//												CONNECTION TO CONNECTION
		//		//
		//		// List of Connection to Connection transformations
		//		FromConnectionTypeToConnectionType: []simplifiedTypes.FromConnectionTypeToConnectionType{
		//			// 			<!--!!E.TransformOriginalToTarget.ConnectionToConnection!!-->
		//			// 			<!--!!S.TransformOriginalToTarget.InitiatorToTransaction!!-->
		//			//
		//			//					Initiator to Initiator in TransactionType
		//			{
		//				FromConnectionType: demo3.DemoInitE,
		//				ToConnectionType:   visi16.VisiInitiatorInTransactionConnection,
		//				AttributeFromTypeToTypeForConnection: []simplifiedTypes.AttributeFromTypeToTypeForConnection{
		//					// In DEMO the DemoInitE connection orientation points from EAR to TK.
		//					// In VISI the connection VisiInitiatorInTransactionConnection starts in TransactionType and points towards the RoleType.
		//					// EAR has become RoleType, and TK has become TransactionType and therefore a change in direction of the connection is enough.
		//					// The sourceId points to the Elementary Actor Role
		//					// The targetId points to the RoleType
		//					{
		//						FromAttributeIdentification: "connection.sourceId",
		//						ToAttributeIdentification:   "connection.targetId",
		//						FromType:                    demo3.DemoEar,
		//						ToType:                      visi16.VisiRoleType,
		//					},
		//					// The sourceId points to the TransactionKind
		//					// The targetId points to the TransactionType
		//					{
		//						FromAttributeIdentification: "connection.targetId",
		//						ToAttributeIdentification:   "connection.sourceId",
		//						FromType:                    demo3.DemoTk,
		//						ToType:                      visi16.VisiTransactionType,
		//					},
		//					// The following step is to give the connection a clear identification. It is set to: "Initiator off" to clarify the direction of the connection. (mainly for readability)
		//					{
		//						FromAttributeIdentification: "connection.connectionType",
		//						ToAttributeIdentification:   "connection.identification",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddInitiatorOff,
		//					},
		//				},
		//			},
		//			// 			<!--!!E.TransformOriginalToTarget.InitiatorToTransaction!!-->
		//			// 			<!--!!S.TransformOriginalToTarget.ExecutorToTransaction!!-->
		//			//					Add Executor to transactiontype
		//			{
		//				FromConnectionType: demo3.DemoExecE,
		//				ToConnectionType:   visi16.VisiExecutorInTransactionConnection,
		//				AttributeFromTypeToTypeForConnection: []simplifiedTypes.AttributeFromTypeToTypeForConnection{
		//					// In DEMO the DemoExecE connection orientation points from EAR to TK.
		//					// In VISI the connection VisiExecutorInTransactionConnection starts in TransactionType and points towards the RoleType.
		//					// EAR has become RoleType, and TK has become TransactionType and therefore a change in direction of the connection is enough.
		//					// The sourceId points to the Elementary Actor Role
		//					// The targetId points to the RoleType
		//					{
		//						FromAttributeIdentification: "connection.targetId",
		//						ToAttributeIdentification:   "connection.targetId",
		//						FromType:                    demo3.DemoEar,
		//						ToType:                      visi16.VisiRoleType,
		//					},
		//					{
		//						FromAttributeIdentification: "connection.sourceId",
		//						ToAttributeIdentification:   "connection.sourceId",
		//						FromType:                    demo3.DemoTk,
		//						ToType:                      visi16.VisiTransactionType,
		//					},
		//					{
		//						FromAttributeIdentification: "connection.connectionType",
		//						ToAttributeIdentification:   "connection.identification",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddExecutorOff,
		//					},
		//				},
		//			},
		//			// 			<!--!!E.TransformOriginalToTarget.ExecutorToTransaction!!-->
		//			// 			<!--!!S.TransformOriginalToTarget.TransactionToMitt!!-->
		//			//					Add Transaction to MITT
		//			{
		//				FromConnectionType: demo3.DemoTkTpsk,
		//				ToConnectionType:   visi16.VisiTransactionInMittConnection,
		//				AttributeFromTypeToTypeForConnection: []simplifiedTypes.AttributeFromTypeToTypeForConnection{
		//					{
		//						FromAttributeIdentification: "connection.sourceId",
		//						ToAttributeIdentification:   "connection.targetId",
		//						FromType:                    demo3.DemoTk,
		//						ToType:                      visi16.VisiTransactionType,
		//					},
		//					{
		//						FromAttributeIdentification: "connection.targetId",
		//						ToAttributeIdentification:   "connection.sourceId",
		//						FromType:                    demo3.DemoTpsk,
		//						ToType:                      visi16.VisiMessageInTransactionType,
		//					},
		//					{
		//						FromAttributeIdentification: "connection.connectionType",
		//						ToAttributeIdentification:   "connection.identification",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddTransactionIn,
		//					},
		//				},
		//			},
		//			// 			<!--!!E.TransformOriginalToTarget.TransactionToMitt!!-->
		//			// 			<!--!!S.TransformOriginalToTarget.GroupToMitt!!-->
		//			//						Add group to mitts
		//			{
		//				FromConnectionType: demo3.DemoTkTpsk,
		//				ToConnectionType:   visi16.VisiGroupInMittConnection,
		//				AttributeFromTypeToTypeForConnection: []simplifiedTypes.AttributeFromTypeToTypeForConnection{
		//					{
		//						FromAttributeIdentification: "connection.sourceId",
		//						ToAttributeIdentification:   "connection.targetId",
		//						FromType:                    demo3.DemoTk,
		//						ToType:                      visi16.VisiGroupType,
		//					},
		//					{
		//						FromAttributeIdentification: "connection.targetId",
		//						ToAttributeIdentification:   "connection.sourceId",
		//						FromType:                    demo3.DemoTpsk,
		//						ToType:                      visi16.VisiMessageInTransactionType,
		//					},
		//					{
		//						FromAttributeIdentification: "connection.connectionType",
		//						ToAttributeIdentification:   "connection.identification",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddGroupIn,
		//					},
		//				},
		//			},
		//			// 			<!--!!E.TransformOriginalToTarget.GroupToMitt!!-->
		//			// 			<!--!!S.TransformOriginalToTarget.PreviousMittToMitt.TpskIn!!-->
		//			//					Add Previous Mitt to Mitt
		//			{
		//				FromConnectionType: demo3.DemoTpskTpskIn,
		//				ToConnectionType:   visi16.VisiPreviousInMittConnection,
		//				AttributeFromTypeToTypeForConnection: []simplifiedTypes.AttributeFromTypeToTypeForConnection{
		//					{
		//						FromAttributeIdentification: "connection.sourceId",
		//						ToAttributeIdentification:   "connection.targetId",
		//						FromType:                    demo3.DemoTpsk,
		//						ToType:                      visi16.VisiMessageInTransactionType,
		//					},
		//					{
		//						FromAttributeIdentification: "connection.targetId",
		//						ToAttributeIdentification:   "connection.sourceId",
		//						FromType:                    demo3.DemoTpsk,
		//						ToType:                      visi16.VisiMessageInTransactionType,
		//					},
		//					{
		//						FromAttributeIdentification: "connection.connectionType",
		//						ToAttributeIdentification:   "connection.identification",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddPreviousIn,
		//					},
		//				},
		//			},
		//			// 			<!--!!E.TransformOriginalToTarget.PreviousMittToMitt.TpskIn!!-->
		//			// 			<!--!!S.TransformOriginalToTarget.PreviousMittToMitt.TpskCallLink!!-->
		//			{
		//				FromConnectionType: demo3.DemoTpskCallLink,
		//				ToConnectionType:   visi16.VisiPreviousInMittConnection,
		//				AttributeFromTypeToTypeForConnection: []simplifiedTypes.AttributeFromTypeToTypeForConnection{
		//					{
		//						FromAttributeIdentification: "connection.sourceId",
		//						ToAttributeIdentification:   "connection.targetId",
		//						FromType:                    demo3.DemoTpsk,
		//						ToType:                      visi16.VisiMessageInTransactionType,
		//					},
		//					{
		//						FromAttributeIdentification: "connection.targetId",
		//						ToAttributeIdentification:   "connection.sourceId",
		//						FromType:                    demo3.DemoTpsk,
		//						ToType:                      visi16.VisiMessageInTransactionType,
		//					},
		//					{
		//						FromAttributeIdentification: "connection.connectionType",
		//						ToAttributeIdentification:   "connection.identification",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddPreviousIn,
		//					},
		//				},
		//			},
		//			// 			<!--!!E.TransformOriginalToTarget.PreviousMittToMitt.TpskCallLink!!-->
		//			// 			<!--!!S.TransformOriginalToTarget.ConditionToMitt!!-->
		//			//						Add Condition to Mitt
		//			{
		//				FromConnectionType: demo3.DemoTpskWaitLink,
		//				ToConnectionType:   visi16.VisiConditionInMittConnection,
		//				AttributeFromTypeToTypeForConnection: []simplifiedTypes.AttributeFromTypeToTypeForConnection{
		//					{
		//						FromAttributeIdentification: "connection.targetId",
		//						ToAttributeIdentification:   "connection.sourceId",
		//						FromType:                    demo3.DemoTpsk,
		//						ToType:                      visi16.VisiMessageInTransactionType,
		//					},
		//					{
		//						FromAttributeIdentification: "connection.targetId",
		//						ToAttributeIdentification:   "connection.targetId",
		//						FromType:                    demo3.DemoTpskWaitLink,
		//						ToType:                      visi16.VisiMessageInTransactionTypeCondition,
		//					},
		//					{
		//						FromAttributeIdentification: "connection.connectionType",
		//						ToAttributeIdentification:   "connection.identification",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddConditionIn,
		//					},
		//				},
		//			},
		//			// 			<!--!!E.TransformOriginalToTarget.ConditionToMitt!!-->
		//			// 			<!--!!S.TransformOriginalToTarget.SendAfterMittToCondition!!-->
		//			//					Add Mitt ref to sendafter in condition
		//			{
		//				FromConnectionType: demo3.DemoTpskWaitLink,
		//				ToConnectionType:   visi16.VisiSendAfterInMittConditionConnection,
		//				AttributeFromTypeToTypeForConnection: []simplifiedTypes.AttributeFromTypeToTypeForConnection{
		//					{
		//						FromAttributeIdentification: "connection.sourceId",
		//						ToAttributeIdentification:   "connection.sourceId",
		//						FromType:                    demo3.DemoTpskWaitLink,
		//						ToType:                      visi16.VisiMessageInTransactionTypeCondition,
		//					},
		//					{
		//						FromAttributeIdentification: "connection.sourceId",
		//						ToAttributeIdentification:   "connection.targetId",
		//						FromType:                    demo3.DemoTpsk,
		//						ToType:                      visi16.VisiMessageInTransactionType,
		//					},
		//					{
		//						FromAttributeIdentification: "sendAfter",
		//						ToAttributeIdentification:   "connection.identification",
		//						FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddSendAfterIn,
		//					},
		//				},
		//			},
		//			// 			<!--!!E.TransformOriginalToTarget.SendAfterMittToCondition!!-->
		//		},
	}

	errs = simplifiedFunctions.Execute(wsClient, simplifiedTypes.ModelTransform, transformRequest) //transformation.TransformModelRequest(&transformRequest)

	data, err := json.Marshal(visiModel)
	if err != nil {
		return nil, errors.New(err)
	}
	result = append(result, simplifiedTypes.Message{MessageId: msg.MessageId, MessageType: simplifiedTypes.MessageTypeExecuted, ContentType: simplifiedTypes.Model, Payload: data})
	return
}

func Unpack(msg *simplifiedTypes.Message) (trans *simplifiedTypes.MessageTransform, error *errors.Error) {
	tfm := &simplifiedTypes.MessageTransform{}
	err := json.Unmarshal(msg.Payload, tfm)
	if err != nil {
		return nil, errors.New(err)
	}
	return tfm, nil
}
func Pack(msg *simplifiedTypes.Message, val string) (result []simplifiedTypes.Message, error *errors.Error) {
	data, err := json.Marshal(val)
	if err != nil {
		return nil, errors.New(err)
	}
	result = append(result, simplifiedTypes.Message{
		MessageId:   msg.MessageId,
		ContentType: simplifiedTypes.ModelText,
		Payload:     data,
	})
	return result, nil
}

func AddRtPrefix(wsClient *client.WebSocketClient, msg *simplifiedTypes.Message) (result []simplifiedTypes.Message, error *errors.Error) {
	tfm, errs := Unpack(msg)
	if errs != nil {
		return nil, errs
	}

	val := "RT" + strings.TrimPrefix(tfm.Met.Identification, "AR")

	result, errs = Pack(msg, val)
	if errs != nil {
		return nil, errs
	}
	return
}
