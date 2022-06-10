package handlers

import (
	"l2gogameserver/data/logger"
	"l2gogameserver/gameserver"
	"l2gogameserver/gameserver/broadcast"
	"l2gogameserver/gameserver/clientpackets"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/gameserver/listeners"
	"l2gogameserver/gameserver/models/chat"
	"l2gogameserver/gameserver/models/say/command"
	"l2gogameserver/packets"
)

// Handler loop клиента в ожидании входящих пакетов
func Handler(client interfaces.ReciverAndSender) {
	for {
		opcode, data, err := client.Receive()
		//defer kickClient(client)
		if err != nil {
			if client.Player().ObjectID() != 0 {
				gameserver.CharOffline(client)
			}
			return
			//break // todo  return ?
		}
		if packets.GetNamePacket(opcode) == "ValidatePosition" && packets.GetNamePacket(opcode) == "MoveBackwardToLocation" {
			logger.Info.Println("Client->Server: #", opcode, packets.GetNamePacket(opcode))
		}

		switch opcode {
		case 0:
			clientpackets.Logout(client, data)
			return
		case 116: //SendBypassBuildCmd
			clientpackets.SendBypassBuildCmd(data, client)
		case 26: //Запрос другому персонажу на желание торговать
			clientpackets.TradeRequest(data, client)
		case 85: //AnswerTradeRequest (если пользователь отвечает Да/Нет на предложение торговли)
			clientpackets.AnswerTradeRequest(data, client)
		case 27: //AddTradeItem
			clientpackets.AddTradeItem(data, client)
		case 28: //tradeDone
			clientpackets.TradeDone(data, client)
		case 13:
			// CharacterDelete
		case 35:
			clientpackets.BypassToServer(data, client)
		case 96:
			clientpackets.DestroyItem(data, client)
		case 14:
			clientpackets.ProtocolVersion(client, data)
		case 43:
			clientpackets.AuthLogin(data, client)
		case 19:
			clientpackets.RequestNewCharacter(client, data)
		case 12:
			clientpackets.CharacterCreate(data, client)
		case 18:
			clientpackets.CharSelected(data, client)
			gameserver.AddOnlineChar(client.Player())
		case 208:
			if len(data) >= 2 {
				switch data[0] {
				case 1:
					clientpackets.RequestManorList(client, data)
				case 54:
					clientpackets.RequestGoToLobby(client, data)
				case 13:
					clientpackets.RequestAutoSoulShot(data, client)
				case 36:
					clientpackets.RequestSaveInventoryOrder(client, data)
				default:
					logger.Info.Println("Не реализованный пакет: ", data[0], packets.GetNamePacket(data[0]))
				}
			}

		case 86:
			if len(data) >= 2 {
				logger.Info.Println(data[0])
				switch data[0] {
				case 0: //посадить персонажа на жопу
					clientpackets.ChangeWaitType(client)

				}

			}

		case 23:
			_, item, modify := clientpackets.DropItem(client, data)
			pkgInventoryUpdate := clientpackets.InventoryUpdate(client, item, modify)
			client.EncryptAndSend(pkgInventoryUpdate)

		case 193:
			clientpackets.RequestObserverEnd(client, data)
		case 108:
			clientpackets.RequestShowMiniMap(client, data)
		case 17:
			clientpackets.RequestEnterWorld(client, data)
			broadcast.BroadCastUserInfoInRadius(client, 2000)
			broadcast.SendCharInfoAboutCharactersInRadius(client, 2000)
			go listeners.StartClientListener(client)
		case 166:
			clientpackets.RequestSkillCoolTime(client, data)
		case 15:

			pkg := clientpackets.MoveBackwardToLocation(client, data)
			broadcast.Checkaem(client, pkg)

		case 73:
			say := clientpackets.Say(client, data)
			if say.Type == chat.SpecialCommand {
				command.ExistCommand(say.Text, client)
			} else {
				broadcast.BroadCastChat(client, say)
			}
		case 89:
			clientpackets.ValidationPosition(data, client.Player())
			//broadcast.Checkaem(client, pkg)
		case 31:
			pkg := clientpackets.Action(data, client)
			if pkg != nil {
				broadcast.Checkaem(client, *pkg)
			}
		case 72:
			clientpackets.RequestTargetCancel(data, client)
		case 114:
			clientpackets.MoveToPawn(client, data)
		case 1:
			clientpackets.Attack(data, client)
		case 25:
			clientpackets.UseItem(client, data)
		case 87:
			clientpackets.RequestRestart(data, client)
			break
		case 57:
			clientpackets.RequestMagicSkillUse(data, client)
		case 61:
			clientpackets.RequestShortCutReg(data, client)
		case 63:
			clientpackets.RequestShortCutDel(data, client)
		case 80:
			clientpackets.RequestSkillList(client, data)
		case 20:
			clientpackets.RequestItemList(client, data)
		case 205:
			clientpackets.RequestMakeMacro(client, data)
		default:
			logger.Info.Println("Not Found case with opcode: ", opcode, packets.GetNamePacket(opcode))
		}

	}
}
