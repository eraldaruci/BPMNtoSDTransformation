package custom

import (
	"fmt"
	"gitlab.com/teec2/simplified/components/serverextender/client"
	"gitlab.com/teec2/simplified/components/serverextender/simplifiedFunctions"
	"gitlab.com/teec2/simplified/components/serverextender/simplifiedTypes"
)

func GetDiagramId(wsClient *client.WebSocketClient, id string) *simplifiedTypes.MessageModelElement {
	diagram := &simplifiedTypes.MessageModelElement{}
	mets, errs := simplifiedFunctions.GetModelElement(wsClient, simplifiedTypes.ModelElementsByModel, id)
	if errs != nil {
		fmt.Println("Error getting ModelElements", errs)
		//return
	}
	if len(*mets) > 0 {
		for _, met := range *mets {
			if met.ElementKind == 1 {
				diagram = &met
				break
			}
		}
		//return
	}
	//fmt.Println("Diagram details ---", diagram.ElementKind)
	return diagram
}
func compiledSourceElements(startElement, sourceElement *simplifiedTypes.MessageModelElement, StartAttr, SourceAttr *simplifiedTypes.MessageModelAttribute) simplifiedTypes.FromElementTypeToElementType {
	return simplifiedTypes.FromElementTypeToElementType{
		FromElementType: StartEvent,
		ToElementType:   Source,
		AttributeFromTypeToTypeForElement: []simplifiedTypes.AttributeFromTypeToTypeForElement{
			//	give the RoleType element the attribute "description". Use the attribute "name" from the demo EAR as value.
			{
				FromAttributeIdentification: "element.identification",
				ToAttributeIdentification:   "element.identification",
			},
			//  Use the demo EAR element identification as element identification for RoleType.
			//  remove leading EAR identification (AR) and replace it with RT (only for readability)
			//{
			//	FromAttributeIdentification: startElement.Identification,
			//	ToAttributeIdentification:   sourceElement.Identification,
			//	FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddRtPrefix,
			//},
			//	give the RoleType element the attribute "id". Use the EAR identification as value. (Adding the attribute id is for readability, xml would support random identifiers)
			//{
			//	FromAttributeIdentification: "element.identification",
			//	ToAttributeIdentification:   "attribute.id",
			//	FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddRtPrefix,
			//},
		},
	}
}

func CompiledSinkElement(EndEventEl, SinkElement *simplifiedTypes.MessageModelElement, EndEventAttr, SinkAttr *simplifiedTypes.MessageModelAttribute) simplifiedTypes.FromElementTypeToElementType {
	return simplifiedTypes.FromElementTypeToElementType{
		FromElementType: EndEvent,
		ToElementType:   Sink,
		AttributeFromTypeToTypeForElement: []simplifiedTypes.AttributeFromTypeToTypeForElement{
			//  Use the demo TransactionKind element identification as element identification for the VISI TransactionType.
			//  remove leading TransactionKind identification (TK) and replace it with TT (only for readability)
			//{
			//	FromAttributeIdentification: EndEventEl.Identification,
			//	ToAttributeIdentification:   SinkElement.Identification,
			//	FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddTtPrefix,
			//},
			//	give the TransactionType element the attribute "description". Use the attribute "name" from the DEMO TransactionKind as value.
			{
				FromAttributeIdentification: "element.identification",
				ToAttributeIdentification:   "element.identification",
			},
			////	give the TransactionType element the attribute "id". Use the TransactionKind identification as value. (Adding the attribute id is for readability, xml would support random identifiers)
			//{
			//	FromAttributeIdentification: "element.identification",
			//	ToAttributeIdentification:   "attribute.id",
			//	FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddTtPrefix,
			//},
		},
	}
}
func CompiledConnections(co *simplifiedTypes.MessageModelConnection, connection1 *simplifiedTypes.MessageModelConnection) simplifiedTypes.FromConnectionTypeToConnectionType {
	return simplifiedTypes.FromConnectionTypeToConnectionType{
		FromConnectionType: co.ConnectionType,
		ToConnectionType:   connection1.ConnectionType,
		AttributeFromTypeToTypeForConnection: []simplifiedTypes.AttributeFromTypeToTypeForConnection{
			{
				FromAttributeIdentification: co.ModelElementIdSource,
				ToAttributeIdentification:   connection1.ModelElementIdSource,
				//FromType:                    BPMN.SequenceFlow202,
				//ToType:                      visi16.VisiRoleType,
			},
			// The sourceId points to the TransactionKind
			// The targetId points to the TransactionType
			{
				FromAttributeIdentification: co.ModelElementIdTarget,
				ToAttributeIdentification:   connection1.ModelElementIdTarget,
				//FromType:                    demo3.DemoTk,
				//ToType:                      visi16.VisiTransactionType,
			},
			// The following step is to give the connection a clear identification. It is set to: "Initiator off" to clarify the direction of the connection. (mainly for readability)
			{
				FromAttributeIdentification: co.Identification,
				ToAttributeIdentification:   connection1.Identification,
				FromFunctionName:            ExtMyModule + simplifiedTypes.ModuleFunctionSep + fAddInitiatorOff,
			},
		},
	}
}
