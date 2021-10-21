from flask import Flask
from flask import request
from flask import render_template
from flask import jsonify
from redis import Redis
from services.product_event_handler import emit_product_order

app = Flask(__name__)
redis = Redis(host='redis', port=6379)

@app.route('/student')
def student():
	return render_template("student.html"), 201

@app.route('/professor')
def professor():
	return render_template("professor.html"), 201

@app.route('/cleardb')
def cleardb():
	redis.flushdb()
	return ("Dd Flushed") , 201


@app.route('/create', methods=['POST'])
def create():
	name = request.form['name']
	weight = request.form['weight']
	redis.mset({name: weight})
	return render_template("professorsub.html", name=name), 201

@app.route('/list', methods=['POST'])
def list():
	html = """\
	<table cellspacing="10">
	<tr><th> </th> <th>Subject</th><th>Weightage</th><th>Enroll to Subject</th></tr>"""
	names = redis.keys('*')
	for name in names:
	   name = name.decode('utf-8')
	   html = html + "<tr>"
	   html = html + "<td> <form id=\"" + str(name) + "\"><input type=\"hidden\" name=\"id\" value=\"" + str(name) + "\" method=\"post\" action=\"enroll\"/></form></td> "
	   html = html + "<td> <input form=\"" + str(name) + "\" type=\"text\" name=\"" + str(name) + "\" value=\"" + str(name) + "\" disabled /> </td>"
	   html = html + "<td> <input form=\"" + str(name) + "\" type=\"text\" name=\"" + str(redis.get(name).decode('utf-8')) + "\" value=\"" + str(redis.get(name).decode('utf-8')) + "\" disabled /> </td>"
	   html = html + "<td> <input form=\"" + str(name) + "\" type=\"submit\" value=\"Enroll\" formmethod=\"post\" formaction=\"/enroll\" /> </td> "
	   html = html + "</tr>"
  
	return html, 201 

@app.route('/enroll', methods=['POST'])
def buy():
	name = request.form['id']
	weight = redis.get(name)
	emit_product_order(name)
	return render_template("thankyou.html", name=name), 201
