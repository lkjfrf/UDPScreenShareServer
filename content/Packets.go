package content

import "encoding/json"

const (
	Error = iota
	DBSignin
	PlayerLogout
	ChannelEnter
	NearPlayerUpdate

	PlayerMove //5
	PlayerActionEvent
	OtherPlayerMove
	PlayerLogin
	OtherPlayerSpawnInfo

	OtherPlayerDestroyInfo //10
	OtherInfo
	Voice
	RoomUserList
	RoomListUpdate

	Permission //15
	KickFromRoom
	MicToggle
	NoticeWrite
	NoticeContent

	NoticeList //20
	NoticeDelete
	NoticeModify
	ChannelCreate
	ChannelDelete

	CalenderRequest //25
	ChannelWidgetInfo
	NormalChat
	PrivateChat
	NoticeChat

	Questions // 30
	Invite
	InviteUserList
	CostumeSet
	UpdateCostume

	OtherUpdateCostume // 35
	HeartBeat
	AllFriendList
	SearchAddFriendList
	SearchDeleteFriendList

	AddFriend // 40
	DeleteFriend
	RequestAddFriend
	RequestDeleteFriend
	SpawnAvatar

	SaveFile //45
	CancelQuestion
	ModifyIntroduce
	FileList
	AccpetQuestion

	QuestionList //50
	ESaveShareData
	EPlaySaveShareData
	EEnterFileComplete
	EUploadComplete

	EScreenDataControlling //55
	RecvFileStatus
	EChannelCreateAfterEnter
	ERemoveFile
	EPlaceFBX

	ETestPlayerLogin //60
	EPlaceLoopingMP4
	EScreenShare
	EDownloadPPTtoPDF
	EGroupChat

	EGroupCreate //65
	EGroupActive
	EGroupUserListUpdate
	ERequestGroupList
	ERequestGroupUserList

	ESaveGroupAlarm //70
	EScreenShareToggle
	EScreenShareView

	Max
)

func JsonStrToStruct[T any](jsonstr string) T {
	var data T
	json.Unmarshal([]byte(jsonstr), &data)
	return data
}

type SR_ScreenShare struct {
	Id     string
	Status int32
	Size   int32
	Width  int32
	Height int32
	Data   []uint16
}

type S_ChannelEnter struct {
	Id          string
	ChannelNum  int32
	ChannelType int32 // 0: Auditorium, 1: Convention, 2: VirtualOffice, 3: VirtualGallery, 4: Plaza
}

type S_ScreenShareToggle struct {
	IsOn bool
}

type S_ScreenShareView struct {
	ViewTarget string
	IsOn       bool
}

type R_ScreenShareView struct {
	IsHasViewer bool
}

type S_PlayerLogout struct {
	Id string
}
