package cmdproc

import "time"

func (r *CmdProcessor) medSetCommand(userID int64, key, name, comment string) []CmdResponse {
	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) medSetTemplateCommand(userID int64, key string) []CmdResponse {
	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) medDelCommand(userID int64, key string) []CmdResponse {
	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) medListCommand(userID int64) []CmdResponse {
	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) medIndicatorSetCommand(
	userID int64,
	ts time.Time,
	medKey string,
	value float64,
) []CmdResponse {
	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) medIndicatorDelCommand(userID int64, ts time.Time, medKey string) []CmdResponse {
	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) medIndicatorReportCommand(userID int64, tsFrom, tsTo time.Time) []CmdResponse {
	return NewSingleCmdResponse(MsgOK)
}
