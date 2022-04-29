//Обработка админского ALT+G

package clientpackets

import (
	"l2gogameserver/gameserver/admin"
	"l2gogameserver/gameserver/interfaces"
	"l2gogameserver/packets"
	"strings"
)

func SendBypassBuildCmd(data []byte, client interfaces.ReciverAndSender) {
	var packet = packets.NewReader(data)
	commandArr := strings.Fields(packet.ReadString())
	admin.Command(client, commandArr)
}
