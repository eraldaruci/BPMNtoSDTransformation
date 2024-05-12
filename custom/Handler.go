package custom

import (
	"fmt"
	"github.com/go-errors/errors"
	"gitlab.com/teec2/simplified/components/serverextender/client"
	"gitlab.com/teec2/simplified/components/serverextender/simplifiedTypes"
)

func HandleMessage(wsClient *client.WebSocketClient, msg *simplifiedTypes.Message) (result []simplifiedTypes.Message) {
	var err *errors.Error
	if module, function, ok := client.GetMessageModuleFunction(msg); ok {
		if module == ExtMyModule {
			fmt.Println("Handle", function)
			switch function {
			case OriginalToTargetTransformModel:
				result, err = TransformOriginalToTarget(wsClient, msg)
			case fAddRtPrefix:
				result, err = AddRtPrefix(wsClient, msg)
				//case fAddTtPrefix:
				//	result, err = AddTtPrefix(wsClient, msg)
				//case fAddGtPrefix:
				//	result, err = AddGtPrefix(wsClient, msg)
				//case fAddGrTpPrefix:
				//	result, err = AddGrTpPrefix(wsClient, msg)
				//case fAddMittPrefix:
				//	result, err = AddMittPrefix(wsClient, msg)
				//case fIsFirstMessage:
				//	result, err = IsFirstMessage(wsClient, msg)
				//case fIsInitiatorToExecutor:
				//	result, err = IsInitiatorToExecutor(wsClient, msg)
				//case fAddAbbreviationToCode:
				//	result, err = AddAbbreviationToCode(wsClient, msg)
				//case fAddStepKindToDescription:
				//	result, err = AddStepKindToDescription(wsClient, msg)
				//case fAddStepKindToIdent:
				//	result, err = AddStepKindToIdent(wsClient, msg)
				//case fAddPhaseInMitt:
				//	result, err = AddPhaseInMitt(wsClient, msg)
				//case fAddSendAfter:
				//	result, err = AddSendAfter(wsClient, msg)
				//case fAddInitiatorOff:
				//	result, err = AddInitiatorOff(wsClient, msg)
				//case fAddExecutorOff:
				//	result, err = AddExecutorOff(wsClient, msg)
				//case fAddTransactionIn:
				//	result, err = AddTransactionIn(wsClient, msg)
				//case fAddGroupIn:
				//	result, err = AddGroupIn(wsClient, msg)
				//case fAddPreviousIn:
				//	result, err = AddPreviousIn(wsClient, msg)
				//case fAddConditionIn:
				//	result, err = AddConditionIn(wsClient, msg)
				//case fAddSendAfterIn:
				//	result, err = AddSendAfterIn(wsClient, msg)
			}
			err = nil
			if err != nil {
				result = append(result, simplifiedTypes.Message{MessageId: msg.MessageId, ErrorNumber: 999, ErrorMessage: err.Error()})
			}
		}
	}
	return
}
