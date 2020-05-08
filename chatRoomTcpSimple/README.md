# chatRoom
读写分离设计：

![Image text](https://github.com/xiaoma-nmg/chatRoom/blob/master/chatRoomTcpSimple/Image/test.jpg)

每个客户端用对应的 goroutine 处理读写

全局的channel负责消息的广播
