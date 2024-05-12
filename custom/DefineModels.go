package custom

import (
	"encoding/json"
	"fmt"
	"gitlab.com/teec2/simplified/components/serverextender/client"
	"gitlab.com/teec2/simplified/components/serverextender/simplifiedFunctions"
	"gitlab.com/teec2/simplified/components/serverextender/simplifiedTypes"
	"strconv"
)

func DefineModels(wsClient *client.WebSocketClient) {
	//log.Println("----------------------------------------->Starting the DefineModels...")
	//startTime := time.Now()
	//Define model names
	const originalName = "BPMNmodel"
	const addName = "-to-SD"
	if wsClient == nil || len(wsClient.Who) < 1 {
		fmt.Println("Not correct connection to the server")
		return
	}
	// Get all the created models
	mdls, errs := simplifiedFunctions.GetModel(wsClient, simplifiedTypes.ModelsByUser, wsClient.Who[0].UserId)
	if errs != nil || mdls == nil {
		fmt.Println("Not correct connection to the Model", errs)
		return
	}
	//Create empty model
	//emptyModel := simplifiedTypes.MessageModel{
	//	Id:             simplifiedFunctions.NewId(),
	//	Identification: originalName,
	//	RepositoryId:   (*mdls)[0].RepositoryId,
	//}
	//
	//demoModels, errs := simplifiedFunctions.SaveModel(wsClient, simplifiedTypes.Model, emptyModel)
	//if errs != nil || len(*demoModels) < 1 {
	//	fmt.Println("Error")
	//	return
	//}
	//demoModel := (*demoModels)[0]
	//mfr1 := simplifiedTypes.MessageModelFolder{
	//	XMLName:        xml.Name{},
	//	Id:             simplifiedFunctions.NewId(),
	//	Identification: "Elements",
	//	ModelId:        demoModel.Id,
	//}
	//mfrs, _ := simplifiedFunctions.SaveModelFolder(wsClient, simplifiedTypes.ModelFolder, &mfr1)
	//mfr1 = (*mfrs)[0]
	//
	//return

	//Delete present translations
	//for _, x := range *mdls {
	//	if x.Identification == originalName || x.Identification == originalName+addName {
	//		fmt.Println(x.Identification + " " + x.Id)
	//		deleted, errs := simplifiedFunctions.DeleteModel(wsClient, simplifiedTypes.Model, x)
	//		if errs != nil || !deleted {
	//			fmt.Println("Error while deleting present translations", errs)
	//			//return
	//		}
	//	}
	//}
	//return

	//Delete present translations
	for _, x := range *mdls {
		if x.Identification == originalName+addName {
			//fmt.Println(x.Identification + " " + x.Id)
			deleted, errs := simplifiedFunctions.DeleteModel(wsClient, simplifiedTypes.Model, x)
			if errs != nil || !deleted {
				fmt.Println("Error while deleting present translations", errs)
				//return
			}
		}
	}

	// Get BPMN model
	sourceModel := simplifiedTypes.MessageModel{}
	for _, mdl := range *mdls {
		if mdl.Identification == sourceModel.Identification || mdl.Identification == originalName {
			sourceModel = mdl
			break
		}
	}

	data, err := json.Marshal(sourceModel)
	if err != nil {
		fmt.Println("Error marshalling BPMN model")
		return
	}
	msgMdl := &simplifiedTypes.Message{
		MessageId:       simplifiedFunctions.NextMessageId(),
		MessageType:     simplifiedTypes.MessageTypeGET,
		ContentType:     simplifiedTypes.Model,
		ProtocolVersion: simplifiedTypes.Version01,
		Payload:         data,
	}
	//</editor-fold>
	fmt.Println("Before Transformation")
	result, errs := TransformOriginalToTarget(wsClient, msgMdl)
	if errs != nil {
		fmt.Println("Error transforming Source (BPMN) model to Target (SFD) model")
		return
	}
	mdlRet := &simplifiedTypes.MessageModel{}
	err = json.Unmarshal(result[0].Payload, mdlRet)
	if err != nil {
		fmt.Println("Error unmarshalling Message Model")
		return
	}
	targetModel := simplifiedTypes.MessageModel{}
	mdls, errs = simplifiedFunctions.GetModel(wsClient, simplifiedTypes.ModelsByUser, wsClient.Who[0].UserId)
	if errs != nil || mdls == nil {
		fmt.Println("Error getting mdls")
		return
	}

	//Check if BPMN-to-SD exists, if yes, assign the new transformed model to that.
	for _, mdl := range *mdls {
		if mdl.Identification == sourceModel.Identification+"-to-SD" {
			//fmt.Println("TransformOriginalToTarget: SFD model details", mdl.Identification+" "+mdl.Id)
			targetModel = mdl
			break
		}
	}

	//Get BPMN folders
	originalFolders, errs := simplifiedFunctions.GetModelFolder(wsClient, simplifiedTypes.ModelFoldersByModel, sourceModel.Id)
	if errs != nil {
		fmt.Println("Problem retrieving Folders from SFD model:", targetModel.Identification)
		return
	}
	//Get SFD folders
	targetFolders, errs := simplifiedFunctions.GetModelFolder(wsClient, simplifiedTypes.ModelFoldersByModel, targetModel.Id)
	if errs != nil {
		fmt.Println("Problem retrieving Folders from SFD model:", targetModel.Identification)
		return
	}

	mets, errs := simplifiedFunctions.GetModelElement(wsClient, simplifiedTypes.ModelElementsByModel, targetModel.Id)
	if errs != nil {
		fmt.Println("Error getting SFD Model Elements", errs)
		//return
	}

	conns, errs := simplifiedFunctions.GetModelConnection(wsClient, simplifiedTypes.ModelConnectionsByModel, sourceModel.Id)
	if errs != nil {
		fmt.Println("Error getting SFD Model Connections", errs)
		//return
	}

	mets1, errs := simplifiedFunctions.GetModelElement(wsClient, simplifiedTypes.ModelElementsByModel, sourceModel.Id)
	if errs != nil {
		fmt.Println("Error getting BPMN Model Elements", errs)

	}

	if len(*targetFolders) > 0 {
		for _, folder := range *targetFolders {

			originalFolder := GetOriginalFolder(folder.Identification, originalFolders)
			sourceDiagram := GetOriginalDiagram(originalFolder.Id, mets1)
			targetDiagram := CreateTargetDiagram(wsClient, targetModel, folder.Id, sourceDiagram.Identification)

			if len(*mets) > 0 {
				for _, e := range *mets {
					//Create visual elements for SFD elements created by transform request
					// mets -> SFD model elements
					// mets1 -> BPMN model elements

					//Get Visual Notations for Elements of transformed diagram
					diagramTypeT := []string{"SFD100"}
					nvts, e3 := simplifiedFunctions.GetNotationVisualElementByElementDiagram(wsClient, e.ElementType, diagramTypeT)
					if e3 != nil {
						fmt.Println("Error getting SFD Notation Visual Element", e3)
					}
					visualSourceEls, er := simplifiedFunctions.GetModelVisualElement(wsClient, simplifiedTypes.ModelVisualElementsByDiagram, sourceDiagram.Id)
					if er != nil {
						fmt.Println("Error getting BPMN Visual Model Elements", er)
						//return
					}
					for _, e1 := range *mets1 {
						if e.ElementType == ActiveStock && e.Identification == "Active "+e1.Identification {
							//visualActivityEl, _ := simplifiedFunctions.GetModelVisualElement(wsClient, simplifiedTypes.ModelElement, e1.Id)
							visualActivityEl := &simplifiedTypes.MessageModelVisualElement{}
							for _, v := range *visualSourceEls {
								if v.ModelElementId == e1.Id {
									visualActivityEl = &v
									break
								}
							}
							_ = CreateVisualActiveStock(wsClient, e.Id, targetModel.Id, targetDiagram.Id, (*nvts)[0].Id, e.ElementType, *visualActivityEl)
							break
						} else if e.ElementType == FinishedStock && e.Identification == "Finished "+e1.Identification {
							//visualActivityEl, _ := simplifiedFunctions.GetModelVisualElement(wsClient, simplifiedTypes.ModelElement, e1.Id)
							visualActivityEl := &simplifiedTypes.MessageModelVisualElement{}
							for _, v := range *visualSourceEls {
								if v.ModelElementId == e1.Id {
									visualActivityEl = &v
									break
								}
							}
							_ = CreateVisualFinishedStock(wsClient, e.Id, targetModel.Id, targetDiagram.Id, (*nvts)[0].Id, e.ElementType, *visualActivityEl)
							break
						} else if e.ElementType == FinishedStock && e.Identification == e1.Identification {
							//visualActivityEl, _ := simplifiedFunctions.GetModelVisualElement(wsClient, simplifiedTypes.ModelElement, e1.Id)
							visualGatewayEl := &simplifiedTypes.MessageModelVisualElement{}
							for _, v := range *visualSourceEls {
								if v.ModelElementId == e1.Id {
									visualGatewayEl = &v
									break
								}
							}
							_ = CreateVisualTargetEl(wsClient, e.Id, targetModel.Id, targetDiagram.Id, (*nvts)[0].Id, e.ElementType, *visualGatewayEl)
							break
						} else if e.ElementType == Sink && e.Identification == e1.Identification {
							//visualEndEl, _ := simplifiedFunctions.GetModelVisualElement(wsClient, simplifiedTypes.ModelElement, e1.Id)
							visualEndEl := &simplifiedTypes.MessageModelVisualElement{}
							for _, v := range *visualSourceEls {
								if v.ModelElementId == e1.Id {
									visualEndEl = &v
									break
								}
							}
							_ = CreateVisualTargetEl(wsClient, e.Id, targetModel.Id, targetDiagram.Id, (*nvts)[0].Id, e.ElementType, *visualEndEl)
							break
						} else if e.ElementType == Source && e.Identification == e1.Identification {
							//visualStartEl, _ := simplifiedFunctions.GetModelVisualElement(wsClient, simplifiedTypes.ModelElement, e1.Id)
							visualStartEl := &simplifiedTypes.MessageModelVisualElement{}
							for _, v := range *visualSourceEls {
								if v.ModelElementId == e1.Id {
									visualStartEl = &v
									break
								}
							}
							_ = CreateVisualTargetEl(wsClient, e.Id, targetModel.Id, targetDiagram.Id, (*nvts)[0].Id, e.ElementType, *visualStartEl)
							break
						}
					}

				}

				// create Visual connections for Connection between ActiveStock & FinishedStock
				conns1, errs1 := simplifiedFunctions.GetModelConnection(wsClient, simplifiedTypes.ModelConnectionsByModel, targetModel.Id)
				if errs1 != nil {
					fmt.Println("Error getting SFD Model Connections", errs)
				}

				visualTargetEls, er := simplifiedFunctions.GetModelVisualElement(wsClient, simplifiedTypes.ModelVisualElementsByDiagram, targetDiagram.Id)
				if er != nil {
					fmt.Println("Error getting SFD Visual Model Elements", er)
				}
				if len(*conns1) > 0 {
					for _, c1 := range *conns1 {

						//fmt.Println(" getting SFD Model Connections", conns1)
						visualHeadEl := &simplifiedTypes.MessageModelVisualElement{}
						for _, v := range *visualTargetEls {
							if v.ModelElementId == c1.ModelElementIdSource {
								visualHeadEl = &v
								break
							}
						}
						visualTailEl := &simplifiedTypes.MessageModelVisualElement{}
						for _, v := range *visualTargetEls {
							if v.ModelElementId == c1.ModelElementIdTarget {
								visualTailEl = &v
								break
							}
						}

						x1 := visualHeadEl.PositionX + 550
						x2 := visualTailEl.PositionX - 100
						y1 := visualHeadEl.PositionY + 25
						y2 := visualTailEl.PositionY + 25
						vcpath := strconv.Itoa(x1) + "," + strconv.Itoa(y1) + "," + strconv.Itoa(x2) + "," + strconv.Itoa(y2)
						//fmt.Println("SFD Connection Path", vcpath)

						_ = CreateSFDVisualConnection(wsClient, &c1, targetDiagram.Id, targetModel.Id, visualHeadEl, visualTailEl, vcpath)

					}
				}

				connectionCount := 0
				if len(*conns) > 0 {
					//Get all BPMN connections
					for _, c := range *conns {
						//fmt.Println("BPMN connections", c)

						//Get the BPMN element at the head of the connection
						cHeadEl, e := simplifiedFunctions.GetModelElement(wsClient, simplifiedTypes.ModelElement, c.ModelElementIdSource)
						if e != nil {
							fmt.Println("Error getting BPMN Connection Head Element", e)
						}
						//Loop through the transformed SFD elements, to get the SFD element that cHead was transformed into
						transformedHead := &simplifiedTypes.MessageModelElement{}
						if len(*mets) > 0 {
							for _, el := range *mets {
								if (*cHeadEl)[0].Identification == el.Identification || "Finished "+(*cHeadEl)[0].Identification == el.Identification {
									transformedHead = &el
									//fmt.Println("Chead SFD", &el)
									break
								}
							}
						}

						//Get the BPMN element at the tail of the connection
						cTailEl, e1 := simplifiedFunctions.GetModelElement(wsClient, simplifiedTypes.ModelElement, c.ModelElementIdTarget)
						if e1 != nil {
							fmt.Println("Error getting BPMN Connection Tail Element", e1)
						}

						//Loop through the transformed SFD elements, to get the SFD element that cTail was transformed into
						transformedTail := &simplifiedTypes.MessageModelElement{}
						if len(*mets) > 0 {
							for _, el := range *mets {
								if (*cTailEl)[0].Identification == el.Identification || "Active "+(*cTailEl)[0].Identification == el.Identification {
									transformedTail = &el
									//fmt.Println(" Ctail SFD", &el)
									break

								}
							}
						}

						//Create new SFD connection with cHeadEl and cTailEl (SFD) elements
						conn2 := CreateConnectionBetweenElements(wsClient, c.Identification, transformedHead, transformedTail, Flow, &targetModel)
						//_, er = simplifiedFunctions.SaveModelConnection(wsClient, simplifiedTypes.ModelConnection, conn2)
						cnvts, e3 := simplifiedFunctions.GetNotationConnection(wsClient, simplifiedTypes.ModelConnection, conn2.Id)
						if e3 != nil {
							fmt.Println("Error getting SFD Notation Visual Connection", e3)
						}
						conn2.NotationConnectionId = (*cnvts)[0].Id

						connectionCount++

						//Now create visuals
						//Visual Elements
						visualHeadEl := &simplifiedTypes.MessageModelVisualElement{}
						for _, v := range *visualTargetEls {
							if v.ModelElementId == transformedHead.Id {
								visualHeadEl = &v
								break
							}
						}

						visualTailEl := &simplifiedTypes.MessageModelVisualElement{}
						for _, v := range *visualTargetEls {
							if v.ModelElementId == transformedTail.Id {
								visualTailEl = &v
								break
							}
						}

						//Visual Connections
						//Adjusting x, y coordinates of the connection path
						x1 := visualHeadEl.PositionX + 150
						x2 := visualTailEl.PositionX - 100
						y1 := visualHeadEl.PositionY + 25
						y2 := visualTailEl.PositionY + 25
						vcpath := strconv.Itoa(x1) + "," + strconv.Itoa(y1) + "," + strconv.Itoa(x2) + "," + strconv.Itoa(y2)

						_ = CreateSFDVisualConnection(wsClient, &conn2, targetDiagram.Id, targetModel.Id, visualHeadEl, visualTailEl, vcpath)
					}

				}

			}

		}

	}
}
