 参考: https://www.jianshu.com/p/ea656dc9b037     enum   oneof   map  package JSON

enum Location {
    SHANGHAI = 0;
    BEIJING = 1;
    GUANGZHOU = 2;
}

//如果你有一些字段同时最多只有一个能被设置，可以使用oneof关键字来实现
message SampleMessage {
    oneof test_oneof {
        string name = 4;
        SubMessage sub_message = 9;
    }
}

protoc -I ../routeguide --go_out=plugins=grpc:../routeguide ../routeguide/route_guide.proto protoc
I 参数：指定import路径，可以指定多个-I参数，编译时按顺序查找，不指定时默认查找当前目录
--go_out ：golang编译支持，支持以下参数
plugins=plugin1+plugin2 - 指定插件，目前只支持grpc，即：plugins=grpc

protoc -I  .   --go_out=plugins=grpc:. XfcUser.proto



