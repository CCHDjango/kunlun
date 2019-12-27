'''
数据处理相关的方法
开发人：陈常鸿
创建时间 ： 2019-12-25
最后一次修改 ： 2019-12-25

注意事项：

'''
import csv
import json
try:
    import pymongo
except ModuleNotFoundError:
    pass

class DataContrl:
    def __init__(self):
        self.mgoURL=''                          # mongoDB数据库的连接地址
        self.mgoPwd=''                          # mongoDB数据库的登陆密码，可为空
        self.localSavePath=''                   # 保存到本地文件的路径

    def settingMgo(self,address,port,password=''):
        '''
        function : 设置mongoDB数据库
        param address : 数据库连接地址 type : string
        param port : 数据库的端口 type : int
        param passwrod : 连接数据库的密码 type : string
        '''
        assert isinstance(address,str)
        assert isinstance(port,int)
        assert isinstance(password,str)
        self.mgoURL=':'.join([address,str(port)])
        self.mgoURL=password

    def settingCSV(self,path,head):
        '''
        function : 设置保存成CSV变量
        param path : 保存成CSV的路径
        param head : CSV列头
        '''
        assert isinstance(path,str)
        self.localSavePath=path