package gameserver

import (
	"fmt"
	"l2gogameserver/gameserver/clientpackets"
	"l2gogameserver/gameserver/models"
	"log"
)

type clientPacketList struct {
	Id   byte   `json:"id"`
	Name string `json:"name"`
}

var clientPacket = []clientPacketList{{
	Id:   0,
	Name: "Logout",
}, {
	Id:   1,
	Name: "AttackRequest",
}, {
	Id:   3,
	Name: "ReqStartPledgeWar",
}, {
	Id:   4,
	Name: "ReqReplyStartPledgeWar",
}, {
	Id:   5,
	Name: "ReqStopPledgeWar",
}, {
	Id:   6,
	Name: "ReqReplyStopPledgeWar",
}, {
	Id:   7,
	Name: "ReqSurrenderPledgeWar",
}, {
	Id:   8,
	Name: "ReqReplySurrenderPledgeWar",
}, {
	Id:   9,
	Name: "ReqSetPledgeCrest",
}, {
	Id:   11,
	Name: "RequestGiveNickName",
}, {
	Id:   12,
	Name: "CharacterCreate",
}, {
	Id:   13,
	Name: "CharacterDelete",
}, {
	Id:   14,
	Name: "ProtocolVersion",
}, {
	Id:   15,
	Name: "MoveBackwardToLocation",
}, {
	Id:   17,
	Name: "EnterWorld",
}, {
	Id:   18,
	Name: "CharSelected",
}, {
	Id:   19,
	Name: "NewCharacter",
}, {
	Id:   20,
	Name: "RequestItemList",
}, {
	Id:   22,
	Name: "RequestUnEquipItem",
}, {
	Id:   23,
	Name: "RequestDropItem",
}, {
	Id:   25,
	Name: "UseItem",
}, {
	Id:   26,
	Name: "TradeRequest",
}, {
	Id:   27,
	Name: "AddTradeItem",
}, {
	Id:   28,
	Name: "TradeDone",
}, {
	Id:   31,
	Name: "Action",
}, {
	Id:   34,
	Name: "RequestLinkHtml",
}, {
	Id:   35,
	Name: "ReqBypassToServer",
}, {
	Id:   36,
	Name: "ReqBBSwrite",
}, {
	Id:   38,
	Name: "ReqJoinPledge",
}, {
	Id:   39,
	Name: "ReqAnswerJoinPledge",
}, {
	Id:   40,
	Name: "ReqWithdrawalPledge",
}, {
	Id:   41,
	Name: "ReqOustPledgeMember",
}, {
	Id:   43,
	Name: "ReqAuthLogin",
}, {
	Id:   44,
	Name: "ReqGetItemFromPet",
}, {
	Id:   46,
	Name: "ReqAllyInfo",
}, {
	Id:   47,
	Name: "ReqCrystallizeItem",
}, {
	Id:   48,
	Name: "ReqPrivateStoreManageSell",
}, {
	Id:   49,
	Name: "SetPrivateStoreListSell",
}, {
	Id:   50,
	Name: "AttackRequest",
}, {
	Id:   52,
	Name: "RequestSocialAction",
}, {
	Id:   53,
	Name: "ChangeMoveType2",
}, {
	Id:   54,
	Name: "ChangeWaitType2",
}, {
	Id:   55,
	Name: "RequestSellItem",
}, {
	Id:   57,
	Name: "RequestMagicSkillUse",
}, {
	Id:   58,
	Name: "Appearing",
}, {
	Id:   59,
	Name: "SendWareHouseDepositList",
}, {
	Id:   60,
	Name: "SendWareHouseWithDrawList",
}, {
	Id:   61,
	Name: "RequestShortCutReg",
}, {
	Id:   63,
	Name: "RequestShortCutDel",
}, {
	Id:   64,
	Name: "RequestBuyItem",
}, {
	Id:   66,
	Name: "RequestJoinParty",
}, {
	Id:   67,
	Name: "RequestAnswerJoinParty",
}, {
	Id:   68,
	Name: "RequestWithDrawalParty",
}, {
	Id:   69,
	Name: "RequestOustPartyMember",
}, {
	Id:   71,
	Name: "CannotMoveAnymore",
}, {
	Id:   72,
	Name: "RequestTargetCancel",
}, {
	Id:   73,
	Name: "Say2",
}, {
	Id:   77,
	Name: "RequestPledgeMemberList",
}, {
	Id:   79,
	Name: "DummyPacket",
}, {
	Id:   80,
	Name: "RequestSkillList",
}, {
	Id:   82,
	Name: "MoveWithDelta",
}, {
	Id:   83,
	Name: "RequestGetOnVehicle",
}, {
	Id:   84,
	Name: "RequestGetOffVehicle",
}, {
	Id:   85,
	Name: "AnswerTradeRequest",
}, {
	Id:   86,
	Name: "RequestActionUse",
}, {
	Id:   87,
	Name: "RequestRestart",
}, {
	Id:   88,
	Name: "RequestSiegeInfo",
}, {
	Id:   89,
	Name: "ValidatePosition",
}, {
	Id:   91,
	Name: "StartRotating",
}, {
	Id:   92,
	Name: "FinishRotating",
}, {
	Id:   94,
	Name: "RequestShowBoard",
}, {
	Id:   95,
	Name: "RequestEnchantItem",
}, {
	Id:   96,
	Name: "RequestDestroyItem",
}, {
	Id:   98,
	Name: "RequestQuestList",
}, {
	Id:   99,
	Name: "RequestQuestAbort",
}, {
	Id:   101,
	Name: "RequestPledgeInfo",
}, {
	Id:   102,
	Name: "RequestPledgeExtendedInfo",
}, {
	Id:   103,
	Name: "RequestPledgeCrest",
}, {
	Id:   107,
	Name: "RequestSendFriendMsg",
}, {
	Id:   108,
	Name: "RequestShowMiniMap",
}, {
	Id:   110,
	Name: "RequestRecordInfo",
}, {
	Id:   111,
	Name: "RequestHennaEquip",
}, {
	Id:   112,
	Name: "RequestHennaRemoveList",
}, {
	Id:   113,
	Name: "RequestHennaItemRemoveInfo",
}, {
	Id:   114,
	Name: "RequestHennaRemove",
}, {
	Id:   115,
	Name: "RequestAcquireSkillInfo",
}, {
	Id:   116,
	Name: "SendBypassBuildCmd",
}, {
	Id:   117,
	Name: "ReqMoveToLocationInVehicle",
}, {
	Id:   118,
	Name: "CannotMoveAnymoreInVehicle",
}, {
	Id:   119,
	Name: "RequestFriendInvite",
}, {
	Id:   120,
	Name: "RequestAnswerFriendInvite",
}, {
	Id:   121,
	Name: "RequestFriendList",
}, {
	Id:   122,
	Name: "RequestFriendDel",
}, {
	Id:   123,
	Name: "CharacterRestore",
}, {
	Id:   124,
	Name: "RequestAcquireSkill",
}, {
	Id:   125,
	Name: "RequestRestartPoint",
}, {
	Id:   126,
	Name: "RequestGMCommand",
}, {
	Id:   127,
	Name: "RequestPartyMatchConfig",
}, {
	Id:   128,
	Name: "RequestPartyMatchList",
}, {
	Id:   129,
	Name: "RequestPartyMatchDetail",
}, {
	Id:   131,
	Name: "RequestPrivateStoreBuy",
}, {
	Id:   133,
	Name: "RequestTutorialLinkHtml",
}, {
	Id:   134,
	Name: "RequestTutorialPassCmdToServer",
}, {
	Id:   135,
	Name: "RequestTutorialQuestionMark",
}, {
	Id:   136,
	Name: "RequestTutorialClientEvent",
}, {
	Id:   137,
	Name: "RequestPetition",
}, {
	Id:   138,
	Name: "RequestPetitionCancel",
}, {
	Id:   139,
	Name: "RequestGMList",
}, {
	Id:   140,
	Name: "RequestJoinAlly",
}, {
	Id:   141,
	Name: "RequestAnswerJoinAlly",
}, {
	Id:   142,
	Name: "AllyLeave",
}, {
	Id:   143,
	Name: "AllyDismiss",
}, {
	Id:   144,
	Name: "RequestDismissAlly",
}, {
	Id:   145,
	Name: "RequestSetAllyCrest",
}, {
	Id:   146,
	Name: "RequestAllyCrest",
}, {
	Id:   147,
	Name: "RequestChangePetName",
}, {
	Id:   148,
	Name: "RequestPetUseItem",
}, {
	Id:   149,
	Name: "RequestGiveItemToPet",
}, {
	Id:   150,
	Name: "ReqPrivateStoreQuitSell",
}, {
	Id:   151,
	Name: "SetPrivateStoreMsgSell",
}, {
	Id:   152,
	Name: "RequestPetGetItem",
}, {
	Id:   153,
	Name: "ReqPrivateStoreManageBuy",
}, {
	Id:   154,
	Name: "SetPrivateStoreListBuy",
}, {
	Id:   156,
	Name: "ReqPrivateStoreQuitBuy",
}, {
	Id:   157,
	Name: "SetPrivateStoreMsgBuy",
}, {
	Id:   159,
	Name: "RequestPrivateStoreSell",
}, {
	Id:   166,
	Name: "RequestSkillCoolTime",
}, {
	Id:   167,
	Name: "ReqPackageSendableItemList",
}, {
	Id:   168,
	Name: "RequestPackageSend",
}, {
	Id:   169,
	Name: "RequestBlock",
}, {
	Id:   170,
	Name: "RequestSiegeInfo",
}, {
	Id:   171,
	Name: "RequestSiegeAttackerList",
}, {
	Id:   172,
	Name: "RequestSiegeDefenderList",
}, {
	Id:   173,
	Name: "RequestJoinSiege",
}, {
	Id:   174,
	Name: "ReqConfirmSiegeWaitingList",
}, {
	Id:   176,
	Name: "MultiSellChoose",
}, {
	Id:   177,
	Name: "NetPing",
}, {
	Id:   179,
	Name: "RequestUserCommand",
}, {
	Id:   180,
	Name: "SnoopQuit",
}, {
	Id:   181,
	Name: "RequestRecipeBookOpen",
}, {
	Id:   182,
	Name: "RequestRecipeBookDestroy",
}, {
	Id:   183,
	Name: "RequestRecipeItemMakeInfo",
}, {
	Id:   184,
	Name: "RequestRecipeItemMakeSelf",
}, {
	Id:   186,
	Name: "RequestRecipeShopMessageSet",
}, {
	Id:   187,
	Name: "RequestRecipeShopListSet",
}, {
	Id:   188,
	Name: "RequestRecipeShopManageQuit",
}, {
	Id:   190,
	Name: "RequestRecipeShopMakeInfo",
}, {
	Id:   191,
	Name: "RequestRecipeShopMakeItem",
}, {
	Id:   192,
	Name: "RequestRecipeShopManagePrev",
}, {
	Id:   193,
	Name: "ObserverReturn",
}, {
	Id:   194,
	Name: "RequestEvaluate",
}, {
	Id:   195,
	Name: "RequestHennaList",
}, {
	Id:   196,
	Name: "RequestHennaItemInfo",
}, {
	Id:   197,
	Name: "RequestBuySeed",
}, {
	Id:   198,
	Name: "DlgAnswer",
}, {
	Id:   199,
	Name: "RequestPreviewItem",
}, {
	Id:   200,
	Name: "RequestSSQStatus",
}, {
	Id:   203,
	Name: "GameGuardReply",
}, {
	Id:   204,
	Name: "RequestPledgePower",
}, {
	Id:   205,
	Name: "RequestMakeMacro",
}, {
	Id:   206,
	Name: "RequestDeleteMacro",
}, {
	Id:   207,
	Name: "RequestBuyProcure",
}}

func getNamePacket(id byte) string {
	for _, p := range clientPacket {
		if p.Id == id {
			return p.Name
		}
	}
	return "NotFind" + string(id)
}

// loop клиента в ожидании входящих пакетов
func (g *GameServer) handler(client *models.Client) {
	for {
		opcode, data, err := client.Receive()

		if err != nil {
			fmt.Println(err)
			fmt.Println("Коннект закрыт")
			break // todo  return ?
		}
		log.Println("Client->Server: #", opcode, getNamePacket(opcode))
		switch opcode {
		case 14:
			pkg := clientpackets.ProtocolVersion(data, client)
			client.SSend(pkg)
		case 43:
			pkg := clientpackets.AuthLogin(data, client)
			client.SSend(pkg)
		case 19:
			pkg := clientpackets.RequestNewCharacter(client, data)
			client.SSend(pkg)
		case 12:
			pkg := clientpackets.CharacterCreate(data, client)
			client.SSend(pkg)
		case 18:
			pkg := clientpackets.CharSelected(data, client)
			client.SSend(pkg)
			g.addOnlineChar(client.CurrentChar)
			go g.ChannelListener(client)

		case 208:
			if len(data) >= 2 {
				switch data[0] {
				case 1:
					pkg := clientpackets.RequestManorList(client, data)
					client.SSend(pkg)
				case 54:
					pkg := clientpackets.RequestGoToLobby(client, data)
					client.SSend(pkg)
				case 13:
					pkg := clientpackets.RequestAutoSoulShot(data, client)
					client.SSend(pkg)
				case 36:
					clientpackets.RequestSaveInventoryOrder(client, data)
				default:
					log.Println("Не реализованный пакет: ", data[0], getNamePacket(data[0]))
				}
			}

		case 86:
			if len(data) >= 2 {
				log.Println(data[0])
				switch data[0] {
				case 0: //посадить персонажа на жопу
					pkg0 := clientpackets.ChangeWaitType(client)
					client.SSend(pkg0)
				case 10: //Продажа в личном лавке
					pkg := clientpackets.PrivateStoreManageListSell(client)
					client.SSend(pkg)
				}

			}

		case 193:
			pkg := clientpackets.RequestObserverEnd(client, data)
			client.SSend(pkg)
		case 108:
			pkg := clientpackets.RequestShowMiniMap(client, data)
			client.SSend(pkg)
		case 17:
			pkg := clientpackets.RequestEnterWorld(client, data)
			client.SSend(pkg)
			g.BroadCastUserInfoInRadius(client, 2000)
			g.GetCharInfoAboutCharactersInRadius(client, 2000)
		case 166:
			pkg := clientpackets.RequestSkillCoolTime(client, data)
			client.SSend(pkg)
		case 15:
			pkg := clientpackets.MoveBackwardToLocation(client, data)
			g.Checkaem(client, pkg)

		case 73:
			say := clientpackets.Say(client, data)
			g.BroadCastChat(client, say)
		case 89:
			pkg := clientpackets.ValidationPosition(data, client.CurrentChar)
			client.SSend(pkg)
		case 31:
			pkg := clientpackets.Action(data, client)
			client.SSend(pkg)
		case 72:
			pkg := clientpackets.RequestTargetCancel(data, client)
			client.SSend(pkg)
		case 1:
			pkg := clientpackets.Attack(data, client)
			client.SSend(pkg)
		case 25:
			pkg := clientpackets.UseItem(client, data)
			client.SSend(pkg)
		case 87:
			pkg := clientpackets.RequestRestart(data, client)
			client.SSend(pkg)
		case 57:
			pkg := clientpackets.RequestMagicSkillUse(data, client)
			client.SSend(pkg)
		case 61:
			pkg := clientpackets.RequestShortCutReg(data, client)
			client.SSend(pkg)
		case 63:
			pkg := clientpackets.RequestShortCutDel(data, client)
			client.SSend(pkg)
		case 80:
			pkg := clientpackets.RequestSkillList(client, data)
			client.SSend(pkg)
		case 20:
			pkg := clientpackets.RequestItemList(client, data)
			client.SSend(pkg)
		case 205:
			pkg := clientpackets.RequestMakeMacro(client, data)
			client.SSend(pkg)
		default:
			log.Println("Not Found case with opcode: ", opcode)
		}

	}
}
