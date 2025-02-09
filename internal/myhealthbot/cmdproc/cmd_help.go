package cmdproc

func (r *CmdProcessor) processHelp(userID int64) []CmdResponse {
	return NewSingleCmdResponse(MsgErrNotImplemented)
}
