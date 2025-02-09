package cmdproc

func (r *CmdProcessor) processMaintenance(cmdParts []string, userID int64) []CmdResponse {
	return NewSingleCmdResponse(MsgErrNotImplemented)
}
