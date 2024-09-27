from flask import Flask, send_from_directory, session
from flask_sock import Sock
import os


from db.client import allocateDefaultBin, getResponses

app = Flask(__name__, static_folder="static/assets", static_url_path="/assets")
app.secret_key = os.urandom(32)
app.config['SOCK_SERVER_OPTIONS'] = {'ping_interval': 25}

sock = Sock(app)

@sock.route("/ws")
def handler(ws):
    try:
        for message in getResponses(session['bin']).listen():
            if message["type"] != "message":
                continue
            ws.send(message["data"].decode("utf-8"))
    except Exception as e:
        print(e)
        pass
    finally:
        ws.close()
        
@app.route("/")
def index():
    if "bin" not in session:
        session["bin"] = allocateDefaultBin()
    return send_from_directory("static", "index.html")


@app.route("/api/bin")
def me():
    return {"bin": session.get("bin", None)}


