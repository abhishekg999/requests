from flask import Flask, render_template, session
from flask_sock import Sock
import os


from db.client import allocateDefaultBin, getResponses

app = Flask(__name__)
app.secret_key = os.urandom(32)
app.config['SOCK_SERVER_OPTIONS'] = {'ping_interval': 25}

sock = Sock(app)

@sock.route("/ws")
def handler(ws):
    for message in getResponses(session['bin']).listen():
        if message["type"] != "message":
            continue
        ws.send(message["data"].decode("utf-8"))
    
@app.route("/")
def index():
    if "bin" not in session:
        session["bin"] = allocateDefaultBin()
    bin = session["bin"]
    return render_template("index.html", bin=bin)


