package consts

var SmsBizType_name = map[int32]string{
	0:  "Login",
	1:  "Register",
	2:  "Edit_Password",
	3:  "Forget_Password",
	4:  "Upgrade_Permit",
	5:  "Upgrade_Audite",
	6:  "Upgrade_Success",
	7:  "Upgrade_Omit",
	8:  "Modify_Mobile",
	9:  "Upgrade_Reward_Frozen",
	10: "Upgrade_Reward_Unfreeze",
	11: "Wallet_Set_Pay_Pass",
	12: "Virtual_Money_OK",
	13: "Transferts_Shift_To",
	14: "Withdraw_Money_Remit",
	15: "Withdraw_Money_Failed",
	16: "Virtual_Money_Freeze",
}

const (
	SmsBizType_Login                   = 0
	SmsBizType_Register                = 1
	SmsBizType_Edit_Password           = 2
	SmsBizType_Forget_Password         = 3
	SmsBizType_Upgrade_Permit          = 4
	SmsBizType_Upgrade_Audite          = 5
	SmsBizType_Upgrade_Success         = 6
	SmsBizType_Upgrade_Omit            = 7
	SmsBizType_Modify_Mobile           = 8
	SmsBizType_Upgrade_Reward_Frozen   = 9
	SmsBizType_Upgrade_Reward_Unfreeze = 10
	SmsBizType_Wallet_Set_Pay_Pass     = 11
	SmsBizType_Virtual_Money_OK        = 12
	SmsBizType_Transferts_Shift_To     = 13
	SmsBizType_Withdraw_Money_Remit    = 14
	SmsBizType_Withdraw_Money_Failed   = 15
	SmsBizType_Virtual_Money_Freeze    = 16
)

const (
	Enabled_No  = 0
	Enabled_Yes = 1
)

const (
	SmsMode_Code   = 0
	SmsMode_Notice = 1
)
