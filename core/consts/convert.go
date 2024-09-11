package consts

import "github.com/tiny-sky/Tdtm/proto"

func ConvertBranchActionToGrpc(action BranchAction) proto.Action {
	switch action {
	case Try:
		return proto.Action_TRY
	case Confirm:
		return proto.Action_CONFIRM
	case Cancel:
		return proto.Action_CANCEL
	case Normal:
		return proto.Action_NORMAL
	case Compensation:
		return proto.Action_COMPENSATION
	default:
	}
	return proto.Action_UN_KNOW_TRANSACTION_TYPE
}

func ConvertTranTypeToGrpc(tranType TransactionType) proto.TranType {
	switch tranType {
	case TCC:
		return proto.TranType_TCC
	case SAGA:
		return proto.TranType_SAGE
	default:
	}
	return proto.TranType_UN_KNOW
}
