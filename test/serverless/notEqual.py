import flask,json
from flask import request

#创建一个服务，把当前这个python文件当做一个服务
server = flask.Flask(__name__)

#server.route()可以将普通函数转变为服务　登录接口的路径、请求方式
@server.route('/',methods=['post'])
def main():
    return "not equal!"

if __name__== '__main__':
    server.run(debug=True,port = 8888,host='0.0.0.0')