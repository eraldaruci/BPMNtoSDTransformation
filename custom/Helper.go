package custom

import (
	"encoding/base64"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"gitlab.com/teec2/simplified/components/serverextender/client"
	"gitlab.com/teec2/simplified/components/serverextender/simplifiedFunctions"
	"gitlab.com/teec2/simplified/components/serverextender/simplifiedTypes"
	"strconv"
	"strings"
	"time"
)

func GetModelByName(wsClient *client.WebSocketClient, name string) (model simplifiedTypes.MessageModel) {
	models, _ := simplifiedFunctions.GetModel(wsClient, simplifiedTypes.Models, "")
	for _, m := range *models {
		if m.Identification == name {
			return m
		}
	}
	return
}

func GetIdentificationByID(wsClient *client.WebSocketClient, id string) (ident string) {
	mes, errs2 := simplifiedFunctions.GetModelElement(wsClient, simplifiedTypes.ModelElements, "")
	if errs2 != nil {
		fmt.Println("Error getting model elements: ", errs2)
	}
	for _, me := range *mes {
		if me.Id == id {
			return me.Identification
		}
	}
	nes, errs3 := simplifiedFunctions.GetNotationConnection(wsClient, simplifiedTypes.NotationConnections, "")
	if errs3 != nil {
		fmt.Println("Error getting model Notation Connections: ", errs3)
	}
	for _, ne := range *nes {
		if ne.Id == id {
			return ne.Identification
		}
	}
	connections, errs4 := simplifiedFunctions.GetModelConnection(wsClient, simplifiedTypes.ModelConnections, "")
	if errs4 != nil {
		fmt.Println("Error getting model Connections: ", errs4)
	}
	for _, connection := range *connections {
		if connection.Id == id {
			return connection.Identification
		}
	}
	models, errs5 := simplifiedFunctions.GetModel(wsClient, simplifiedTypes.Models, "")
	if errs5 != nil {
		fmt.Println("Error getting models: ", errs5)
	}
	for _, model := range *models {
		if model.Id == id {
			return model.Identification
		}
	}
	return "could not find Identification for ID: " + id
}

// GetDemoElementTranslatedFromVisiElement
// return the DEMO ElementId which was the origin of the given VISI element
func GetDemoElementTranslatedFromVisiElement(wsClient *client.WebSocketClient, visiElement *simplifiedTypes.MessageModelElement, translationModel string) (demoElement []*simplifiedTypes.MessageModelElement) {
	connections, errs := simplifiedFunctions.GetModelConnection(wsClient, simplifiedTypes.ModelConnections, "")
	if errs != nil {
		panic("GetModelConnections Failed")
	}
	for _, connection := range *connections {
		if connection.ConnectionType == translationModel {
			if connection.ModelElementIdSource == visiElement.Id {
				demoElm, _ := simplifiedFunctions.GetModelElement(wsClient, simplifiedTypes.ModelElement, connection.ModelElementIdTarget)
				demoElement = append(demoElement, &(*demoElm)[0])
				break
			}
		}
	}
	return
}

// GetVisiElementTranslatedFromDemoElement
// return the VISI ElementId which is the result of the translation of a given DEMO element
func GetVisiElementTranslatedFromDemoElement(wsClient *client.WebSocketClient, demoElement *simplifiedTypes.MessageModelElement, translationModel string) (visiElement []*simplifiedTypes.MessageModelElement) {
	connections, errs := simplifiedFunctions.GetModelConnection(wsClient, simplifiedTypes.ModelConnections, "")
	if errs != nil {
		panic("GetModelConnections Failed")
	}
	for _, connection := range *connections {
		if connection.ConnectionType == translationModel {
			if connection.ModelElementIdSource == demoElement.Id {
				visiElm, _ := simplifiedFunctions.GetModelElement(wsClient, simplifiedTypes.ModelElement, connection.ModelElementIdTarget)
				visiElement = append(visiElement, &(*visiElm)[0])
				break
			}
		}
	}
	return
}

func GetAttributeValueForElement(wsClient *client.WebSocketClient, Id string, attribute string) (res interface{}) {
	attribute = strings.ToLower(attribute)
	attributes, _ := simplifiedFunctions.GetModelAttribute(wsClient, simplifiedTypes.ModelAttributesByModelElement, Id)
	if len(*attributes) == 0 {
		attributes, _ = simplifiedFunctions.GetModelAttribute(wsClient, simplifiedTypes.ModelAttributesByModelConnection, Id)
		if len(*attributes) == 0 {
			fmt.Println("\"Helper.go - GetAttributeValueForElement - Element or Connection could not be found in cache: ")
		}
	}
	for _, mat := range *attributes {
		if strings.ToLower(mat.Identification) == attribute {
			//if mat.ValueString != nil && *mat.ValueString != "" {
			if mat.Value != "" {
				res = mat.Value
				//fmt.Println("valueString: ", valueString)
				return
			} else {
				fmt.Println("Helper.go:GetAttributeValueForElement:\t", mat.Identification, "has no value for: ", attribute)
				return
			}
		}
	}

	//fmt.Println("Helper.go - GetAttributeValueForElement - There are no Attributes for element: ", element[0].Identification, elementId)

	return
}

func SetAttributeValueForElement(wsClient *client.WebSocketClient, elementId string, attributeValue interface{}, attribute string) {
	if PrintFunctionNames {
		fmt.Println("Helper.go - SetAttributeValueForElement()")
	}
	var mdls *[]simplifiedTypes.MessageModel
	modelAttributesForElement, errs := simplifiedFunctions.GetModelAttribute(wsClient, simplifiedTypes.ModelAttributesByModelElement, elementId)
	if errs != nil {
		panic("GetModelAttributesByModelElement Failed")
		//fmt.Println("\"Helper.go - GetModelAttributesByModelElement - Attributes could not be loaded from cache")
		//return
	}
	if len(*modelAttributesForElement) > 0 {
		for _, mat := range *modelAttributesForElement {
			if strings.ToLower(mat.Identification) == strings.ToLower(attribute) {
				//attribute already exists for this element.
				var attrMet *[]simplifiedTypes.MessageModelAttribute
				attributeValueType := fmt.Sprintf("%T", attributeValue)
				if attributeValueType != mat.ElementType {
					//changing ValueType. Remove old attributeValue, but reuse attributeID and all to prevent having both attributes.
					matNew := simplifiedTypes.MessageModelAttribute{
						Id:                  mat.Id,
						Identification:      mat.Identification,
						ModelId:             mat.ModelId,
						NotationAttributeId: mat.NotationAttributeId,
						FolderId:            mat.FolderId,
						ElementId:           mat.ElementId,
						ConnectionId:        mat.ConnectionId,
						ConnectionBeginId:   mat.ConnectionBeginId,
						ConnectionEndId:     mat.ConnectionEndId,
					}
					switch t := attributeValue.(type) {
					case string:
						strValue := t
						matNew.Value = strValue
						matNew.ElementType = "string"
					case int:
						intValue := t
						matNew.Value = strconv.Itoa(intValue)
						matNew.ElementType = "integer"
					case time.Time:
						timeValue := t
						matNew.Value = timeValue.String()
						matNew.ElementType = "time"
					case float64:
						doubleValue := t
						matNew.Value = strconv.FormatFloat(doubleValue, 'E', -1, -1)
						matNew.ElementType = "double"
					case []byte:
						binaryValue := t
						matNew.Value = string(binaryValue)
						matNew.ElementType = "double"
					case []string:
						//For subtransactions in transactiontype
						uuidValue := t
						if len(t) > 0 {
							matNew.Value = uuidValue[0]
							matNew.ElementType = "[]string"
						}
					}
					*attrMet = append(*attrMet, matNew)
				} else {
					at := make(map[string]interface{})
					at[attribute] = attributeValue
					neas, errs2 := simplifiedFunctions.GetNotationAttribute(wsClient, simplifiedTypes.NotationAttributesByNotationElement, elementId)
					if errs2 != nil {
						panic("GetNotationAttributesByNotationElement failed")
					}
					element, errs2 := simplifiedFunctions.GetModelElement(wsClient, simplifiedTypes.ModelElement, elementId)
					if errs2 != nil {
						panic("GetModelElement failed")
					}
					attrMet, _ = simplifiedFunctions.CreateAttribute(wsClient, simplifiedTypes.ModelAttributesByModelElement, &(*element)[0], neas, &at)
				}
				mdls, errs = simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.ModelAttribute, (*attrMet)[0])
				if errs != nil {
					panic("SaveModelAttribute Failed")
				}
				break
			}
		}
	}
	modelAttributesForElement, errs = simplifiedFunctions.GetModelAttribute(wsClient, simplifiedTypes.ModelAttributesByModelElement, elementId)
	if len(*mdls) > 0 {
		element, errs := simplifiedFunctions.GetModelElement(wsClient, simplifiedTypes.ModelElement, elementId)
		neas, errs2 := simplifiedFunctions.GetNotationAttribute(wsClient, simplifiedTypes.NotationAttributesByNotationElement, (*element)[0].NotationElementId)
		if errs2 != nil {
			panic("GetNotationAttributesByNotationElement failed")
		}
		at := make(map[string]interface{})
		at[attribute] = attributeValue
		attrMet, _ := simplifiedFunctions.CreateAttribute(wsClient, simplifiedTypes.ModelAttributesByModelElement, &(*element)[0], neas, &at)
		mdls, errs = simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.ModelAttribute, (*attrMet)[0])
		if len(*mdls) > 0 || errs != nil {
			panic("SaveModelAttribute Failed")
		}
	}
}

func GetAttributeIdForElement(wsClient *client.WebSocketClient, Id string, attribute string) (attributeId string) {
	attribute = strings.ToLower(attribute)
	attributes, _ := simplifiedFunctions.GetModelAttribute(wsClient, simplifiedTypes.ModelAttributesByModelElement, Id)
	if len(*attributes) == 0 {
		attributes, _ = simplifiedFunctions.GetModelAttribute(wsClient, simplifiedTypes.ModelAttributesByModelConnection, Id)
		if len(*attributes) == 0 {

		}
	}
	for _, attr := range *attributes {
		if strings.ToLower(attr.Identification) == attribute {
			return attr.Id
		}
	}
	return
}

func GetTranslatedFrom(wsClient *client.WebSocketClient, id string, typ string) (from []string) {
	connections, _ := simplifiedFunctions.GetModelConnection(wsClient, simplifiedTypes.ModelConnections, "")
	var from2 []string
	for _, connection := range *connections {
		if connection.ConnectionType == typ && connection.ModelElementIdTarget == id {
			from2 = append(from2, connection.ModelElementIdSource)
		}
	}
	from = from2
	return
}

func GetTranslatedInto(wsClient *client.WebSocketClient, id string, typ string) (into []string) {
	connections, _ := simplifiedFunctions.GetModelConnection(wsClient, simplifiedTypes.ModelConnections, "")
	var into2 []string
	for _, connection := range *connections {
		if connection.ConnectionType == typ && connection.ModelElementIdSource == id {
			into2 = append(into2, connection.ModelElementIdTarget)
		}
	}
	into = into2
	return
}

func CreateModel(wsClient *client.WebSocketClient, modelName string) (model *simplifiedTypes.MessageModel, errs *errors.Error) {
	if len(modelName) > 0 {
		m := simplifiedTypes.MessageModel{
			Id:             simplifiedFunctions.NewId(),
			Identification: modelName,
			RepositoryId:   "", //util.DummyRepositoryId1Demo,
		}
		mdls, errs := simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.Model, m)
		if len(*mdls) > 0 || errs != nil {
			panic("SaveModel Failed")
		}
		time.Sleep(10 * time.Millisecond)
		model = &m
	}
	return
}

func GetModelElementsOfElementType(wsClient *client.WebSocketClient, elementType string, model *simplifiedTypes.MessageModel) (elementsOfType []simplifiedTypes.MessageModelElement) {
	elements, errs := simplifiedFunctions.GetModelElement(wsClient, simplifiedTypes.ModelElementsByModel, model.Id)
	if errs != nil {
		return
	}
	for _, element := range *elements {
		if element.ElementType == elementType {
			elementsOfType = append(elementsOfType, element)
		}
	}
	return
}

func GetModelConnectionsOfConnectionType(wsClient *client.WebSocketClient, connectionType string, model *simplifiedTypes.MessageModel) (connectionsOfType []simplifiedTypes.MessageModelConnection) {
	connections, errs := simplifiedFunctions.GetModelConnection(wsClient, simplifiedTypes.ModelConnectionsByModel, model.Id)
	if errs != nil {
		return
	}
	for _, connection := range *connections {
		if connection.ConnectionType == connectionType {
			connectionsOfType = append(connectionsOfType, connection)
		}
	}
	return
}

func DuplicateModel(wsClient *client.WebSocketClient, original *simplifiedTypes.MessageModel, duplicate *simplifiedTypes.MessageModel) (success bool, errs *errors.Error) {
	var elem []simplifiedTypes.MessageModelElement
	var attr []simplifiedTypes.MessageModelAttribute
	var conn []simplifiedTypes.MessageModelConnection
	translationMap := make(map[string]string)
	//first copy the elements and their attributes
	elements, _ := simplifiedFunctions.GetModelElement(wsClient, simplifiedTypes.ModelElementsByModel, original.Id)
	attributes, _ := simplifiedFunctions.GetModelAttribute(wsClient, simplifiedTypes.ModelAttributesByModel, original.Id)
	for _, element := range *elements {
		met := simplifiedTypes.MessageModelElement{
			Id:             simplifiedFunctions.NewId(),
			Identification: element.Identification,
			ModelElementId: element.ModelElementId,
			ModelId:        duplicate.Id,
			//ModelFolderId:     util.NilId,
			ModelFolderId:     element.ModelFolderId,
			NotationElementId: element.NotationElementId,
			ElementType:       element.ElementType,
			AttributeValue:    element.AttributeValue,
			ElementKind:       element.ElementKind,
		}
		translationMap[element.Id] = met.Id
		elem = append(elem, met)
	}
	for _, met := range elem {
		mdls, errs := simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.ModelElement, met)
		if len(*mdls) > 0 || errs != nil {
			panic("SaveModelElement Failed")
		}
	}

	connections, _ := simplifiedFunctions.GetModelConnection(wsClient, simplifiedTypes.ModelConnectionsByModel, original.Id)
	for _, connection := range *connections {
		mcn := simplifiedTypes.MessageModelConnection{
			Id:                   simplifiedFunctions.NewId(),
			Identification:       connection.Identification,
			ConnectionType:       connection.ConnectionType,
			ModelElementIdSource: translationMap[connection.ModelElementIdSource],
			ModelElementIdTarget: translationMap[connection.ModelElementIdTarget],
			ModelId:              duplicate.Id,
			NotationConnectionId: connection.NotationConnectionId,
		}
		conn = append(conn, mcn)
	}
	for _, met := range conn {
		mdls, errs := simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.ModelConnection, met)
		if len(*mdls) > 0 || errs != nil {
			panic("SaveModelConnection Failed 1")
		}
		if errs != nil {
			return false, errs
		}
	}
	for _, element := range *elements {
		at := make(map[string]interface{})
		duplicateElement, _ := simplifiedFunctions.GetModelElement(wsClient, simplifiedTypes.ModelElement, translationMap[element.Id])
		if len(*duplicateElement) > 0 {
			neas, _ := simplifiedFunctions.GetNotationAttribute(wsClient, simplifiedTypes.NotationAttributesByNotationElement, (*duplicateElement)[0].NotationElementId)
			for _, attribute := range *attributes {
				if attribute.ElementId == element.Id {
					if attribute.Value != "" {
						at[attribute.Identification] = attribute.Value
					}
				}
			}
			attrMet, _ := simplifiedFunctions.CreateAttribute(wsClient, simplifiedTypes.ModelAttribute, &(*duplicateElement)[0], neas, &at)
			attr = append(attr, *attrMet...)
		}
	}
	for _, met := range attr {
		mdls, errs := simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.ModelAttribute, met)
		if len(*mdls) > 0 || errs != nil {
			panic("SaveModelAttribute Failed")
		}
	}
	return true, errs
}

// GetConnectionTargetElementsOfTypeForElement returns []simplifiedTypes.MessageModelElement get all DatabaseModelConnections
// checks if the given element is the source element, and if the connectionType is the same as the requested one
// returns the target of the connection, or multiple if they exist
func GetConnectionTargetElementsOfTypeForElement(wsClient *client.WebSocketClient, connectionType string, source *simplifiedTypes.MessageModelElement) (connectionTargets []*simplifiedTypes.MessageModelElement) {
	connectionTargets = nil
	connections, _ := simplifiedFunctions.GetModelConnection(wsClient, simplifiedTypes.ModelConnectionsByModel, source.ModelId)
	if len(*connections) > 0 {
		for _, connection := range *connections {
			if connection.ConnectionType == connectionType && connection.ModelElementIdSource == source.Id {
				newElement, _ := simplifiedFunctions.GetModelElement(wsClient, simplifiedTypes.ModelElement, connection.ModelElementIdTarget)
				if len(*newElement) > 0 {
					connectionTargets = append(connectionTargets, &(*newElement)[0])
				}
			}
		}
	}
	//fmt.Println(connectionTargets)
	//i := connectionTargets
	//i = i
	return
}

// GetConnectionTargetIdOfTypeForElement returns []string get all DatabaseModelConnections
// checks if the given element is the source element, and if the connectionType is the same as the requested one
// returns the targets as []string of the connection, or multiple if they exist
func GetConnectionTargetIdOfTypeForElement(wsClient *client.WebSocketClient, connectionType string, source *simplifiedTypes.MessageModelElement) (connectionTargetIds []string) {
	connections, _ := simplifiedFunctions.GetModelConnection(wsClient, simplifiedTypes.ModelConnections, "")
	if len(*connections) > 0 {
		for _, connection := range *connections {
			if connection.ConnectionType == connectionType && connection.ModelElementIdSource == source.Id {
				connectionTargetIds = append(connectionTargetIds, connection.ModelElementIdTarget)
			}
		}
	}
	return
}

func CreateColumnWidthForPrint(input string, desiredWidth int) (result string) {
	result = input
	for len(result) < desiredWidth {
		result = result + " "
	}
	return
}

func CreateConnectionBetweenElements(wsClient *client.WebSocketClient, identification string, source *simplifiedTypes.MessageModelElement, target *simplifiedTypes.MessageModelElement, connectionType string, model *simplifiedTypes.MessageModel) (connection simplifiedTypes.MessageModelConnection) {
	x := uuid.New()
	//uuid := x.String()
	uuidBytes, _ := x.MarshalBinary()
	base64UUID := base64.StdEncoding.EncodeToString(uuidBytes)

	connection = simplifiedTypes.MessageModelConnection{
		Id:                   base64UUID,
		Identification:       identification,
		ConnectionType:       connectionType,
		ModelElementIdSource: source.Id,
		ModelElementIdTarget: target.Id,
		ModelId:              model.Id,
		NotationConnectionId: "",
		VersionObjectId:      model.Id,
	}
	_, errs := simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.ModelConnection, connection)
	if errs != nil {
		panic("SaveModelConnection Failed 2")
	}
	return
}

func SaveAttributeTranslationConnection(wsClient *client.WebSocketClient, source *simplifiedTypes.MessageModelAttribute, target *simplifiedTypes.MessageModelAttribute, connectionType string) {
	attributeTranslationConnection := simplifiedTypes.MessageModelConnection{
		Id:                   simplifiedFunctions.NewId(),
		Identification:       "AttrTrans - " + source.Identification + "/" + target.Identification,
		ConnectionType:       connectionType,
		ModelElementIdSource: source.Id,
		ModelElementIdTarget: target.Id,
		ModelId:              source.ModelId,
		NotationConnectionId: "",
	}
	mdls, errs := simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.ModelConnection, attributeTranslationConnection)
	if len(*mdls) > 0 || errs != nil {
		panic("SaveModelConnection Failed 3")
	}
}

func ChangeModelIdentification(wsClient *client.WebSocketClient, newIdentification string, model *simplifiedTypes.MessageModel) {
	newModel := simplifiedTypes.MessageModel{
		Id:             model.Id,
		Identification: newIdentification,
		RepositoryId:   model.RepositoryId,
		Created:        model.Created,
		Updated:        model.Updated,
	}
	mdls, errs := simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.Model, newModel)
	if len(*mdls) > 0 || errs != nil {
		panic("SaveModel Failed")
	}
}

func ChangeModelTransformationConnection(wsClient *client.WebSocketClient, newTransformationConnection string, oldTransformationConnection string, model *simplifiedTypes.MessageModel) {
	connections, _ := simplifiedFunctions.GetModelConnection(wsClient, simplifiedTypes.ModelConnectionsByModel, model.Id)
	var newCons []simplifiedTypes.MessageModelConnection
	for _, connection := range *connections {
		if connection.ConnectionType == oldTransformationConnection {
			newCon := simplifiedTypes.MessageModelConnection{
				Id:                   connection.Id,
				Identification:       connection.Identification,
				ConnectionType:       newTransformationConnection,
				ModelElementIdSource: connection.ModelElementIdSource,
				ModelElementIdTarget: connection.ModelElementIdTarget,
				ModelId:              connection.ModelId,
				NotationConnectionId: connection.NotationConnectionId,
				Created:              connection.Created,
				Updated:              connection.Updated,
			}
			newCons = append(newCons, newCon)
		}
	}
	if len(newCons) > 0 {
		for _, newCon := range newCons {
			mdls, errs := simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.ModelConnection, newCon)
			if len(*mdls) > 0 || errs != nil {
				panic("SaveModelConnection Failed 4")
			}
		}
	}
}

func SaveTranslationConnection(wsClient *client.WebSocketClient, sourceId string, targetId string, ident string, transformConnection string, newId string) {
	sourceType := ""
	sourceModel := ""
	var sourceElement *[]simplifiedTypes.MessageModelElement
	var sourceConnection *[]simplifiedTypes.MessageModelConnection
	var sourceAttribute *[]simplifiedTypes.MessageModelAttribute
	sourceElement, _ = simplifiedFunctions.GetModelElement(wsClient, simplifiedTypes.ModelElement, sourceId)
	if len(*sourceElement) == 0 {
		sourceConnection, _ = simplifiedFunctions.GetModelConnection(wsClient, simplifiedTypes.ModelConnection, sourceId)
		if len(*sourceConnection) == 0 {
			sourceAttribute, _ = simplifiedFunctions.GetModelAttribute(wsClient, simplifiedTypes.ModelAttribute, sourceId)
			if len(*sourceAttribute) == 0 {
				panic("SavingTranslationConnection Failed. SourceId does not belong to an element, connection or attribute")
			} else {
				sourceModel = (*sourceAttribute)[0].ModelId
				sourceType = (*sourceAttribute)[0].Identification
			}
		} else {
			sourceModel = (*sourceConnection)[0].ModelId
			sourceType = (*sourceConnection)[0].ConnectionType
		}
	} else {
		sourceModel = (*sourceElement)[0].ModelId
		sourceType = (*sourceElement)[0].ElementType
	}
	targetType := ""
	//targetModel := util.NilId
	var targetElement *[]simplifiedTypes.MessageModelElement
	var targetConnection *[]simplifiedTypes.MessageModelConnection
	var targetAttribute *[]simplifiedTypes.MessageModelAttribute
	targetElement, _ = simplifiedFunctions.GetModelElement(wsClient, simplifiedTypes.ModelElement, targetId)
	if len(*targetElement) == 0 {
		targetConnection, _ = simplifiedFunctions.GetModelConnection(wsClient, simplifiedTypes.ModelConnection, targetId)
		if len(*targetConnection) == 0 {
			targetAttribute, _ = simplifiedFunctions.GetModelAttribute(wsClient, simplifiedTypes.ModelAttribute, targetId)
			if len(*targetAttribute) == 0 {
				//laten returnen, later nog een keer proberen als het target element wel gemaakt is.
				fmt.Println("sourceId: ", sourceId, "sourceType: ", sourceType, "sourceElement: ", sourceElement, "sourceConnection: ", sourceConnection, "sourceAttribute: ", sourceAttribute)
				fmt.Println("targetId: ", targetId, "targetType: ", targetType, "targetElement: ", targetElement, "targetConnection: ", targetConnection, "targetAttribute: ", targetAttribute)
				panic("SavingTranslationConnection Failed. TargetId does not belong to an element, connection or attribute")
			} else {
				//targetModel = targetAttribute[0].ModelId
				targetType = (*targetAttribute)[0].Identification
			}
		} else {
			//targetModel = targetConnection[0].ModelId
			targetType = (*targetConnection)[0].ConnectionType
		}
	} else {
		//targetModel = targetElement[0].ModelId
		targetType = (*targetElement)[0].ElementType
	}
	//if sourceModel == targetModel {
	//	panic("SavingTranslationConnection Failed. SourceModel is the same as TargetModel")
	//}
	if sourceType == "" || targetType == "" {
		panic("")
	}
	mcn := simplifiedTypes.MessageModelConnection{
		Id: newId,
		//Id:             uuid.NewV4(),
		Identification: ident + " - transCon",
		//Identification:       sourceType + " - " + targetType + " - transCon",
		ConnectionType:       transformConnection,
		ModelElementIdSource: sourceId,
		ModelElementIdTarget: targetId,
		ModelId:              sourceModel,
		NotationConnectionId: "",
	}
	mdls, errs := simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.ModelConnection, mcn)
	if len(*mdls) > 0 || errs != nil {
		panic("SaveModelConnection failed")
	}
}
