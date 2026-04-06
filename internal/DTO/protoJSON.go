package DTO

import "google.golang.org/protobuf/encoding/protojson"

// ProtoJSON - настройки protojson для HTTP ответа.
// UseProtoNames=true => имена полей в JSON будут snake_case как в proto (access_token),
// а не lowerCamelCase (accessToken).
//
// EmitUnpopulated=false => не шлем "нулевые" значения (пустые строки/0/false),
// ответ компактнее и не путает клиента (но иногда надо наоборот - решается продуктом).
var ProtoJSON = protojson.MarshalOptions{
	UseProtoNames:   true,  // snake_case для proto полей
	EmitUnpopulated: false, // не мусорим нулевыми
}
