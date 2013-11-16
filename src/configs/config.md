游戏配置文件:
    "#"表示注释，必须出现在行头，即当前行如果使用“#”开头，则表明此行为注释行；
    每一个配置键必须唯一(不能为“=”)，如果重复出现后面的键将覆盖前面键；
    “=”表示配置赋值，等号右边如果有多行，需要使用“"”号进行包含，如果内容中出现“"”符号，请使用“\"”进行转义
    如果单行值，可以不需要使用“""”
    
# game.config
# 游戏服务器配置：
# game.welcome:登陆游戏发送的欢迎信息
# game.host:游戏监听的IP地址
# game.port:游戏监听端口
#####################################################################################################
game.welcome="
    你好，这个用来测试多行配置！
    谢谢！
"
game.host="localhost"
game.port=1234