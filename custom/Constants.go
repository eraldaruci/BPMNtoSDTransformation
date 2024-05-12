package custom

const (
	ExtMyServer   = "OriginalToTargetServerTEEC2"
	ExtMyModule   = "OriginalToTarget"
	ExtMyServerId = "OriginalToTargetServerTEEC2"
	ExtMyFunction = "MyFunction"
)

const (
	DEMO3ToVISI16Module            = "DEMO3ToVISI16"
	OriginalToTargetTransformModel = "TransformOriginalToTarget"
	PrintFunctionNames             = false
)

// BPMN
const (
	BPMN         = "BPMN"
	StartEvent   = "StartEvent200"
	EndEvent     = "EndEvent200"
	Activity     = "Activity200"
	SplitGateway = "SplitGateway200"
	JoinGateway  = "JoinGateway200"
	SequenceFlow = "SequenceFlow200"
)

// SFD
const (
	SFD           = "SFD100"
	ActiveStock   = "ActiveStock100"
	FinishedStock = "FinishedStock100"
	Conflux       = "Conflux100"
	Source        = "Source100"
	Sink          = "Sink100"
	Converter     = "Converter100"
	Flow          = "Flow100"
	Active        = "_Active"
	Finished      = "_Finished"
)

// Functions
const (
	fAddRtPrefix              = "AddRtPrefix"
	fAddTtPrefix              = "AddTtPrefix"
	fAddGtPrefix              = "AddGtPrefix"
	fAddMittPrefix            = "AddMittPrefix"
	fIsFirstMessage           = "IsFirstMessage"
	fIsInitiatorToExecutor    = "IsInitiatorToExecutor"
	fAddAbbreviationToCode    = "AddAbbreviationToCode"
	fAddStepKindToDescription = "AddStepKindToDescription"
	fAddStepKindToIdent       = "AddStepKindToIdent"
	fAddPhaseInMitt           = "AddPhaseInMitt"
	fAddSendAfter             = "AddSendAfter"
	fAddInitiatorOff          = "AddInitiatorOff"
	fAddExecutorOff           = "AddExecutorOff"
	fAddTransactionIn         = "AddTransactionIn"
	fAddGroupIn               = "AddGroupIn"
	fAddPreviousIn            = "AddPreviousIn"
	fAddConditionIn           = "AddConditionIn"
	fAddSendAfterIn           = "AddSendAfterIn"
	fAddGrTpPrefix            = "AddGrTpPrefix"
)
