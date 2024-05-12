package custom

import (
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"gitlab.com/teec2/simplified/components/serverextender/client"
	"gitlab.com/teec2/simplified/components/serverextender/simplifiedFunctions"
	"gitlab.com/teec2/simplified/components/serverextender/simplifiedTypes"
)

// <!--!!S.TransformOriginalToTarget!!-->

// TransformOriginalToTarget translates a BPMN model(simplifiedTypes.MessageModel) into a SFD model(simplifiedTypes.MessageModel).
//
//	 Creates SFD model(simplifiedTypes.MessageModel) and returns the model as simplifiedTypes.MessageModel
//		Also saves the translations in the new SFD Model as connections with connectionType translationModel, which is returned as string
//			input:
//					bpmnModel *simplifiedTypes.MessageModel, bpmnModel to be used
//					uc *simplifiedTypes.UserContextInfo, user credentials
//			returns:
//					translationConnections string,
//					SFDmodel simplifiedTypes.MessageModel
func TransformOriginalToTarget(wsClient *client.WebSocketClient, msg *simplifiedTypes.Message) (result []simplifiedTypes.Message, error *errors.Error) {
	originalModel := &simplifiedTypes.MessageModel{}
	err := json.Unmarshal(msg.Payload, &originalModel)
	if err != nil {
		return nil, errors.New(err)
	}
	//Set translated DatabaseModel name. Uses original BPMN model identification + "-to-SD-TransCon". The ModelConnections are added to the original model from which they are translated.
	//This will enable to keep track of which element, and which connection has been translated into what element or what connection in the new model.
	translationConnection := originalModel.Identification + "-to-SD-TransCon"

	//Create New SFD model, set Identification as BPMN model identification + "to-SFD". This model will be used to store the to be translated model. It is given a random "Id", uses the original model identification and appends "-to-SD" to it.
	//It also uses the same "repositoryId" as the original model.
	found := false
	targetModel := simplifiedTypes.MessageModel{}
	mdls, errs := simplifiedFunctions.GetModel(wsClient, simplifiedTypes.ModelsByUser, wsClient.Who[0].UserId)
	if errs != nil || mdls == nil {
		return nil, errs
	}

	//Check if BPMN-to-SD exists, if yes, assign the new transformed model to that. If not, create a new model BPMN-to-SD
	for _, mdl := range *mdls {
		if mdl.Identification == originalModel.Identification+"-to-SD" {
			//fmt.Println("TransformOriginalToTarget: SFD model details", mdl.Identification+" "+mdl.Id)
			targetModel = mdl
			found = true
			break
		}
	}
	if !found {
		targetModel = simplifiedTypes.MessageModel{
			Id:              simplifiedFunctions.NewId(),
			Identification:  originalModel.Identification + "-to-SD",
			RepositoryId:    originalModel.RepositoryId,
			VersionObjectId: originalModel.RepositoryId,
		}
		mdls, errs = simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.Model, &targetModel)
		if errs != nil || mdls == nil {
			return nil, errs
		}
		targetModel = (*mdls)[0]
	}
	//return
	//Delete present translations
	for i, mdl := range *mdls {
		if mdl.Identification == targetModel.Identification {
			targetModel = (*mdls)[i]
			found = true
			break
		}
	}
	if !found {
		//Save the new model. Uses the uc(user credentials) to check if this user is allowed to do this procedure.
		//In order to be able to reference to this model and store the created elements, connections, and attributes, the model will have to exist, otherwise it would return nil, and crash the translation.
		var mdlSave *[]simplifiedTypes.MessageModel
		mdlSave, errs = simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.Model, targetModel)
		//_, _, errs = simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.Model, targetModel)
		//In case the saving of the model does not succeed. Create panic with error message.
		if len(*mdlSave) < 1 || errs != nil {
			return nil, errs
		}
	}

	//Get originalFolders in original model.
	originalFolders, errs := simplifiedFunctions.GetModelFolder(wsClient, simplifiedTypes.ModelFoldersByModel, originalModel.Id)
	//Check for errors while retrieving the source modelFolders. If error occurs halt the program and show error.
	if errs != nil {
		fmt.Println("Problem retrieving modelFolders from BPMN model Id:", originalModel.Id)
		return nil, errs
	}

	//Create a new folder in SD model
	if len(*originalFolders) > 0 {
		for _, folder := range *originalFolders {
			CreateSDFolder(wsClient, targetModel, folder.Identification)
		}
	}
	//Get SFD folders
	targetFolders, errs := simplifiedFunctions.GetModelFolder(wsClient, simplifiedTypes.ModelFoldersByModel, targetModel.Id)
	if errs != nil {
		fmt.Println("Problem retrieving Folders from SFD model:", targetModel.Identification)
		return nil, errs
	}

	var T_BPMN_FROM_Activity = []simplifiedTypes.ChainType{
		{
			Id:         simplifiedTypes.NewId(),      // Unique identifier of the modeling element, generated for each transformation
			ObjectType: Activity,                     //Describes the type of the source element (e.g. Activity, StartEvent etc.)
			ObjectKind: simplifiedTypes.ModelElement, // ObjectKind, can be ModelElement or ModelConnection
			Attributes: []simplifiedTypes.AttributeRef{ //List of Attributes to save, so they can be used in the target tree
				{
					Id:                      simplifiedTypes.NewId(),  // Unique identifier of the modeling element
					FunctionName:            "",                       // optional: used to perform further manipulations to object attributes
					AttributeIdentification: "element.Identification", //define identification to be passed to target element
					VariableExpression:      "Ident1",
				},
			},
			Next: nil, //Any further element/connection is defined here, otherwise null
		},
	}
	var T_BPMN_FROM_StartEvent = []simplifiedTypes.ChainType{
		{
			Id: simplifiedTypes.NewId(),
			//Object Type
			ObjectType: StartEvent,
			// ObjectKind, can be ModelElement or ModelConnection
			ObjectKind: simplifiedTypes.ModelElement,
			//List of Attributes to save so they can be used in the target tree
			Attributes: []simplifiedTypes.AttributeRef{
				{
					Id:                      simplifiedTypes.NewId(),
					FunctionName:            "",
					AttributeIdentification: "element.Identification",
					VariableExpression:      "Ident1",
				},
			},
			Next: nil,
		},
	}
	var T_BPMN_FROM_SplitG = []simplifiedTypes.ChainType{
		{
			Id:         simplifiedTypes.NewId(),
			ObjectType: SplitGateway,
			ObjectKind: simplifiedTypes.ModelElement,
			Attributes: []simplifiedTypes.AttributeRef{
				{
					Id:                      simplifiedTypes.NewId(),
					FunctionName:            "",
					AttributeIdentification: "element.Identification",
					VariableExpression:      "Ident1",
				},
			},
			Next: nil,
		},
	}
	var T_BPMN_FROM_JoinG = []simplifiedTypes.ChainType{
		{
			Id:         simplifiedTypes.NewId(),
			ObjectType: JoinGateway,
			ObjectKind: simplifiedTypes.ModelElement,
			Attributes: []simplifiedTypes.AttributeRef{
				{
					Id:                      simplifiedTypes.NewId(),
					FunctionName:            "",
					AttributeIdentification: "element.Identification",
					VariableExpression:      "Ident1",
				},
			},
			Next: nil,
		},
	}
	var T_BPMN_FROM_EndEvent = []simplifiedTypes.ChainType{
		{
			Id:         simplifiedTypes.NewId(),
			ObjectType: EndEvent,
			ObjectKind: simplifiedTypes.ModelElement,
			Attributes: []simplifiedTypes.AttributeRef{
				{
					Id:                      simplifiedTypes.NewId(),
					FunctionName:            "",
					AttributeIdentification: "element.Identification",
					VariableExpression:      "Ident1",
				},
			},
			Next: nil,
		},
	}

	var T_SD_TO_AStock_FStock = []simplifiedTypes.ChainType{
		{
			Id:         simplifiedTypes.NewId(),
			ObjectType: ActiveStock,
			ObjectKind: simplifiedTypes.ModelElement,
			Attributes: []simplifiedTypes.AttributeRef{
				{Id: simplifiedTypes.NewId(),
					FunctionName:            "",
					AttributeIdentification: "element.Identification",
					VariableExpression:      "Active {Ident1}", //Modify attribute identification by adding string "Active"
				}},
			Next: []simplifiedTypes.ChainType{ //Define any further element/connection on the pattern
				{Id: simplifiedTypes.NewId(),
					ObjectType: Flow,
					ObjectKind: simplifiedTypes.ModelConnection,
					Direction:  simplifiedTypes.ConnectionDirectionIdToSource, //defines the relation to predecessor ChainType
					Next: []simplifiedTypes.ChainType{
						{Id: simplifiedTypes.NewId(),
							ObjectType: FinishedStock,
							ObjectKind: simplifiedTypes.ModelElement,
							Direction:  simplifiedTypes.ConnectionDirectionTargetToId,
							//In target chains attributes are being filled
							Attributes: []simplifiedTypes.AttributeRef{
								{Id: simplifiedTypes.NewId(),
									FunctionName:            ExtMyFunction + simplifiedTypes.ModuleFunctionSep + Active,
									AttributeIdentification: "element.Identification",
									VariableExpression:      "Finished {Ident1}",
								},
							},
							Next: nil,
						}}}}}}
	var T_SD_TO_Source = []simplifiedTypes.ChainType{
		{
			Id:         simplifiedTypes.NewId(),
			ObjectType: Source,
			ObjectKind: simplifiedTypes.ModelElement,
			Attributes: []simplifiedTypes.AttributeRef{
				{
					Id:                      simplifiedTypes.NewId(),
					FunctionName:            "",
					AttributeIdentification: "element.Identification",
					VariableExpression:      "{Ident1}",
				},
			},
			Next: nil,
		},
	}
	var T_SD_TO_Stock = []simplifiedTypes.ChainType{
		{
			Id:         simplifiedTypes.NewId(),
			ObjectType: FinishedStock,
			ObjectKind: simplifiedTypes.ModelElement,
			Attributes: []simplifiedTypes.AttributeRef{
				{
					Id:                      simplifiedTypes.NewId(),
					FunctionName:            "",
					AttributeIdentification: "element.Identification",
					VariableExpression:      "{Ident1}",
				},
			},
			Next: nil,
		},
	}
	var T_SD_TO_Sink = []simplifiedTypes.ChainType{
		{
			Id:         simplifiedTypes.NewId(),
			ObjectType: Sink,
			ObjectKind: simplifiedTypes.ModelElement,
			Attributes: []simplifiedTypes.AttributeRef{
				{
					Id:                      simplifiedTypes.NewId(),
					FunctionName:            "",
					AttributeIdentification: "element.Identification",
					VariableExpression:      "{Ident1}",
				},
			},
			Next: nil,
		},
	}

	folder := (*targetFolders)[0]
	transformRequest := simplifiedTypes.TransformModelRequest{
		FromModel:           originalModel.Id,
		ToModel:             targetModel.Id,
		FolderDestiny:       folder.Id,
		TransformConnection: translationConnection,
		FromChainTypeToChainType: []simplifiedTypes.FromChainTypeToChainType{
			{
				FromChainType:                  T_BPMN_FROM_Activity,
				ToChainType:                    T_SD_TO_AStock_FStock,
				CreateOnePerSameIdentification: true,
			},
			{
				FromChainType:                  T_BPMN_FROM_StartEvent,
				ToChainType:                    T_SD_TO_Source,
				CreateOnePerSameIdentification: true,
			},
			{
				FromChainType:                  T_BPMN_FROM_SplitG,
				ToChainType:                    T_SD_TO_Stock,
				CreateOnePerSameIdentification: true,
			},
			{
				FromChainType:                  T_BPMN_FROM_JoinG,
				ToChainType:                    T_SD_TO_Stock,
				CreateOnePerSameIdentification: true,
			},
			{
				FromChainType:                  T_BPMN_FROM_EndEvent,
				ToChainType:                    T_SD_TO_Sink,
				CreateOnePerSameIdentification: true,
			},
		},
		//FromConnectionTypeToConnectionType: []simplifiedTypes.FromConnectionTypeToConnectionType{
		//	{
		//		FromConnectionType: SequenceFlow,
		//		ToConnectionType:   Flow,
		//		AttributeFromTypeToTypeForConnection: []simplifiedTypes.AttributeFromTypeToTypeForConnection{
		//			{
		//				FromType: SequenceFlow,
		//				ToType:   Flow,
		//			},
		//		},
		//	},
		//},
	}
	errs = simplifiedFunctions.Execute(wsClient, simplifiedTypes.ModelTransform, transformRequest)
	data, err := json.Marshal(targetModel)
	if err != nil {
		return nil, errors.New(err)
	}
	//myString := string(data[:])
	//fmt.Println("-----> ", myString)
	result = append(result, simplifiedTypes.Message{MessageId: msg.MessageId, MessageType: simplifiedTypes.MessageTypeExecuted, ContentType: simplifiedTypes.Model, Payload: data})
	return
}
