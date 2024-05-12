package custom

import (
	"gitlab.com/teec2/simplified/components/serverextender/simplifiedTypes"
)

func Configuration() (result []simplifiedTypes.MessageConfiguration) {
	// Extra configurations can be customised.
	result = append(result, simplifiedTypes.MessageConfiguration{
		Server:         ExtMyServer,
		Identification: ExtMyServerId,
		Action:         simplifiedTypes.ActionExecute,
		Module:         ExtMyModule,
		Function:       ExtMyFunction,
		Context: []simplifiedTypes.MessageMenu{
			simplifiedTypes.MessageMenu{
				Identification:                "",
				ParentIdentification:          ExtMyModule,
				MenuType:                      simplifiedTypes.MenuRepository,
				DisplayText:                   "Original to Target",
				DisplayFolder:                 "Transform",
				NotationIdentification:        "SD",
				NotationVersion:               "1.0",
				NotationElementIdentification: "model",
				NotationElementVersion:        "",
				Parameters:                    "{model.Id}",
			},
		},
	})

	result = append(result, simplifiedTypes.MessageConfiguration{
		Server:         ExtMyServer,
		Identification: ExtMyServerId + OriginalToTargetTransformModel,
		Action:         simplifiedTypes.ActionExecute,
		Module:         ExtMyModule,
		Function:       OriginalToTargetTransformModel,
		Context: []simplifiedTypes.MessageMenu{
			simplifiedTypes.MessageMenu{
				Identification:                "",
				ParentIdentification:          ExtMyModule,
				MenuType:                      simplifiedTypes.MenuContext,
				DisplayText:                   "Transform Original to Target",
				DisplayFolder:                 "Generate",
				NotationIdentification:        "SD",
				NotationVersion:               "1.0",
				NotationElementIdentification: "TransactionKind",
				NotationElementVersion:        "3.7",
				Parameters:                    "{element.Id}",
			},
		},
	})
	result = append(result, simplifiedTypes.MessageConfiguration{
		Server:         ExtMyServer,
		Identification: ExtMyServerId + fAddRtPrefix,
		Action:         simplifiedTypes.ActionGet,
		Module:         ExtMyModule,
		Function:       fAddRtPrefix,
	})
	//result = append(result, simplifiedTypes.MessageConfiguration{
	//	Server:         ExtMyServer,
	//	Identification: ExtMyServerId + fAddTtPrefix,
	//	Action:         simplifiedTypes.ActionGet,
	//	Module:         ExtMyModule,
	//	Function:       fAddTtPrefix,
	//})
	//result = append(result, simplifiedTypes.MessageConfiguration{
	//	Server:         ExtMyServer,
	//	Identification: ExtMyServerId + fAddGtPrefix,
	//	Action:         simplifiedTypes.ActionGet,
	//	Module:         ExtMyModule,
	//	Function:       fAddGtPrefix,
	//})
	//result = append(result, simplifiedTypes.MessageConfiguration{
	//	Server:         ExtMyServer,
	//	Identification: ExtMyServerId + fAddGrTpPrefix,
	//	Action:         simplifiedTypes.ActionGet,
	//	Module:         ExtMyModule,
	//	Function:       fAddGrTpPrefix,
	//})
	//result = append(result, simplifiedTypes.MessageConfiguration{
	//	Server:         ExtMyServer,
	//	Identification: ExtMyServerId + fAddMittPrefix,
	//	Action:         simplifiedTypes.ActionGet,
	//	Module:         ExtMyModule,
	//	Function:       fAddMittPrefix,
	//})
	//result = append(result, simplifiedTypes.MessageConfiguration{
	//	Server:         ExtMyServer,
	//	Identification: ExtMyServerId + fIsFirstMessage,
	//	Action:         simplifiedTypes.ActionGet,
	//	Module:         ExtMyModule,
	//	Function:       fIsFirstMessage,
	//})
	//result = append(result, simplifiedTypes.MessageConfiguration{
	//	Server:         ExtMyServer,
	//	Identification: ExtMyServerId + fIsInitiatorToExecutor,
	//	Action:         simplifiedTypes.ActionGet,
	//	Module:         ExtMyModule,
	//	Function:       fIsInitiatorToExecutor,
	//})
	//result = append(result, simplifiedTypes.MessageConfiguration{
	//	Server:         ExtMyServer,
	//	Identification: ExtMyServerId + fAddAbbreviationToCode,
	//	Action:         simplifiedTypes.ActionGet,
	//	Module:         ExtMyModule,
	//	Function:       fAddAbbreviationToCode,
	//})
	//result = append(result, simplifiedTypes.MessageConfiguration{
	//	Server:         ExtMyServer,
	//	Identification: ExtMyServerId + fAddStepKindToDescription,
	//	Action:         simplifiedTypes.ActionGet,
	//	Module:         ExtMyModule,
	//	Function:       fAddStepKindToDescription,
	//})
	//result = append(result, simplifiedTypes.MessageConfiguration{
	//	Server:         ExtMyServer,
	//	Identification: ExtMyServerId + fAddStepKindToIdent,
	//	Action:         simplifiedTypes.ActionGet,
	//	Module:         ExtMyModule,
	//	Function:       fAddStepKindToIdent,
	//})
	//result = append(result, simplifiedTypes.MessageConfiguration{
	//	Server:         ExtMyServer,
	//	Identification: ExtMyServerId + fAddPhaseInMitt,
	//	Action:         simplifiedTypes.ActionGet,
	//	Module:         ExtMyModule,
	//	Function:       fAddPhaseInMitt,
	//})
	//result = append(result, simplifiedTypes.MessageConfiguration{
	//	Server:         ExtMyServer,
	//	Identification: ExtMyServerId + fAddSendAfter,
	//	Action:         simplifiedTypes.ActionGet,
	//	Module:         ExtMyModule,
	//	Function:       fAddSendAfter,
	//})
	//result = append(result, simplifiedTypes.MessageConfiguration{
	//	Server:         ExtMyServer,
	//	Identification: ExtMyServerId + fAddInitiatorOff,
	//	Action:         simplifiedTypes.ActionGet,
	//	Module:         ExtMyModule,
	//	Function:       fAddInitiatorOff,
	//})
	//result = append(result, simplifiedTypes.MessageConfiguration{
	//	Server:         ExtMyServer,
	//	Identification: ExtMyServerId + fAddExecutorOff,
	//	Action:         simplifiedTypes.ActionGet,
	//	Module:         ExtMyModule,
	//	Function:       fAddExecutorOff,
	//})
	//result = append(result, simplifiedTypes.MessageConfiguration{
	//	Server:         ExtMyServer,
	//	Identification: ExtMyServerId + fAddTransactionIn,
	//	Action:         simplifiedTypes.ActionGet,
	//	Module:         ExtMyModule,
	//	Function:       fAddTransactionIn,
	//})
	//result = append(result, simplifiedTypes.MessageConfiguration{
	//	Server:         ExtMyServer,
	//	Identification: ExtMyServerId + fAddGroupIn,
	//	Action:         simplifiedTypes.ActionGet,
	//	Module:         ExtMyModule,
	//	Function:       fAddGroupIn,
	//})
	//result = append(result, simplifiedTypes.MessageConfiguration{
	//	Server:         ExtMyServer,
	//	Identification: ExtMyServerId + fAddPreviousIn,
	//	Action:         simplifiedTypes.ActionGet,
	//	Module:         ExtMyModule,
	//	Function:       fAddPreviousIn,
	//})
	//result = append(result, simplifiedTypes.MessageConfiguration{
	//	Server:         ExtMyServer,
	//	Identification: ExtMyServerId + fAddConditionIn,
	//	Action:         simplifiedTypes.ActionGet,
	//	Module:         ExtMyModule,
	//	Function:       fAddConditionIn,
	//})
	//result = append(result, simplifiedTypes.MessageConfiguration{
	//	Server:         ExtMyServer,
	//	Identification: ExtMyServerId + fAddSendAfterIn,
	//	Action:         simplifiedTypes.ActionGet,
	//	Module:         ExtMyModule,
	//	Function:       fAddSendAfterIn,
	//})
	return
}
