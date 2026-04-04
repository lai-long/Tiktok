package dto

type User struct {
	ID        string `protobuf:"bytes,1,opt,name=ID,proto3" form:"ID" json:"id" query:"ID"`
	Username  string `protobuf:"bytes,2,opt,name=Username,proto3" binding:"required" form:"username" json:"username" query:"Username"`
	Password  string `protobuf:"bytes,3,opt,name=Password,proto3" binding:"required" form:"password" json:"password" query:"Password"`
	AvatarURL string `protobuf:"bytes,4,opt,name=AvatarURL,proto3" form:"AvatarURL" json:"avatar_url" query:"AvatarURL"`
	CreatedAt string `protobuf:"bytes,5,opt,name=CreatedAt,proto3" form:"CreatedAt" json:"created_at" query:"CreatedAt"`
	UpdatedAt string `protobuf:"bytes,6,opt,name=UpdatedAt,proto3" form:"UpdatedAt" json:"updated_at" query:"UpdatedAt"`
	DeletedAt string `protobuf:"bytes,7,opt,name=DeletedAt,proto3" form:"DeletedAt" json:"deleted_at" query:"DeletedAt"`
}
